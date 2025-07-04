// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func WatchWebsocket(ctx context.Context, u *url.URL, c EventChannel) (Watcher, error) {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	wsClosed := make(chan struct{})
	go cleanupOnDone(ctx, conn, c, wsClosed)

	d := newDispatcher()

	go forwardWebsocketEvents(conn, c, u, wsClosed, d.done)

	go forwardDispatchNow(ctx, conn, c, d.trigger, u.String())

	return d, nil
}

func cleanupOnDone(ctx context.Context, conn *websocket.Conn, c EventChannel, wsClosed <-chan struct{}) {
	// wait for one of the events to occur
	select {
	case <-ctx.Done():
	case <-wsClosed:
	}

	// clean up channel
	close(c)
	// attempt to close the websocket
	// ignore errors as we are closing everything anyway
	_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "context canceled by server"))
	_ = conn.Close()
}

func forwardWebsocketEvents(ws *websocket.Conn, c EventChannel, u *url.URL, wsClosed chan<- struct{}, sendNowDone chan<- int) {
	serverURL := source(u.String())

	defer func() {
		// this triggers the cleanupOnDone subroutine
		close(wsClosed)
	}()

	for {
		// receive messages, this call is blocking
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == websocket.CloseNormalClosure {
				return
			}
			// assuming the connection got closed through context canceling
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "read" && strings.Contains(opErr.Err.Error(), "closed") {
				return
			}
			c <- &ErrorEvent{
				error:  errors.WithStack(err),
				source: serverURL,
			}
			return
		}

		var eventsSend int
		_, err = fmt.Sscanf(string(msg), messageSendNowDone, &eventsSend)
		if err == nil {
			sendNowDone <- eventsSend
			continue
		}

		e, err := unmarshalEvent(msg)
		if err != nil {
			c <- &ErrorEvent{
				error:  err,
				source: serverURL,
			}
			continue
		}
		localURL := *u
		localURL.Path = e.Source()
		e.setSource(localURL.String())
		c <- e
	}
}

func forwardDispatchNow(ctx context.Context, ws *websocket.Conn, c EventChannel, sendNow <-chan struct{}, serverURL string) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-sendNow:
			if !ok {
				return
			}

			if err := ws.WriteMessage(websocket.TextMessage, []byte(messageSendNow)); err != nil {
				c <- &ErrorEvent{
					source: source(serverURL),
					error:  err,
				}
			}
		}
	}
}
