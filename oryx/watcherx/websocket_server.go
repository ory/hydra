// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/ory/herodot"
)

type (
	eventChannelSlice struct {
		sync.Mutex
		cs []EventChannel
	}
	websocketWatcher struct {
		wsWriteLock      sync.Mutex
		wsReadLock       sync.Mutex
		wsClientChannels eventChannelSlice
	}
)

const (
	messageSendNow     = "send values now"
	messageSendNowDone = "done sending %d values"
)

func WatchAndServeWS(ctx context.Context, u *url.URL, writer herodot.Writer) (http.HandlerFunc, error) {
	c := make(EventChannel)
	watcher, err := Watch(ctx, u, c)
	if err != nil {
		return nil, err
	}
	w := &websocketWatcher{
		wsClientChannels: eventChannelSlice{},
	}
	go w.broadcaster(ctx, c)
	return w.serveWS(ctx, writer, watcher), nil
}

func (ww *websocketWatcher) broadcaster(ctx context.Context, c EventChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-c:
			ww.wsClientChannels.Lock()
			for _, cc := range ww.wsClientChannels.cs {
				cc <- e
			}
			ww.wsClientChannels.Unlock()
		}
	}
}

func (ww *websocketWatcher) readWebsocket(ws *websocket.Conn, c chan<- struct{}, watcher Watcher) {
	for {
		// blocking call to ReadMessage that waits for a close message
		ww.wsReadLock.Lock()
		_, msg, err := ws.ReadMessage()
		ww.wsReadLock.Unlock()

		switch errTyped := err.(type) {
		case nil:
			if string(msg) == messageSendNow {
				done, err := watcher.DispatchNow()
				if err != nil {
					// we cant do much about this error
					ww.wsWriteLock.Lock()
					_ = ws.WriteJSON(&ErrorEvent{
						error:  err,
						source: "",
					})
					ww.wsWriteLock.Unlock()
				}

				go func() {
					eventsSend := <-done

					ww.wsWriteLock.Lock()
					defer ww.wsWriteLock.Unlock()

					// we cant do much about this error
					_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(messageSendNowDone, eventsSend)))
				}()
			}
		case *websocket.CloseError:
			if errTyped.Code == websocket.CloseNormalClosure {
				close(c)
				return
			}
		case *net.OpError:
			if errTyped.Op == "read" && strings.Contains(errTyped.Err.Error(), "closed") {
				// the context got canceled and therefore the connection closed
				close(c)
				return
			}
		default:
			// some other unexpected error, best we can do is return
			return
		}
	}
}

func (ww *websocketWatcher) serveWS(ctx context.Context, writer herodot.Writer, watcher Watcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := (&websocket.Upgrader{
			ReadBufferSize:  256, // the only message we expect is the close message
			WriteBufferSize: 1024,
		}).Upgrade(w, r, nil)
		if err != nil {
			writer.WriteError(w, r, err)
			return
		}

		// make channel and register it at broadcaster
		c := make(EventChannel)
		ww.wsClientChannels.Lock()
		ww.wsClientChannels.cs = append(ww.wsClientChannels.cs, c)
		ww.wsClientChannels.Unlock()

		wsClosed := make(chan struct{})
		go ww.readWebsocket(ws, wsClosed, watcher)

		defer func() {
			// attempt to close the websocket
			// ignore errors as we are closing everything anyway
			ww.wsWriteLock.Lock()
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server context canceled"))
			ww.wsWriteLock.Unlock()

			_ = ws.Close()

			ww.wsClientChannels.Lock()
			for i, cc := range ww.wsClientChannels.cs {
				if c == cc {
					ww.wsClientChannels.cs[i] = ww.wsClientChannels.cs[len(ww.wsClientChannels.cs)-1]
					ww.wsClientChannels.cs[len(ww.wsClientChannels.cs)-1] = nil
					ww.wsClientChannels.cs = ww.wsClientChannels.cs[:len(ww.wsClientChannels.cs)-1]
				}
			}
			ww.wsClientChannels.Unlock()
			close(c)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-wsClosed:
				return
			case e, ok := <-c:
				if !ok {
					return
				}

				ww.wsWriteLock.Lock()
				err := ws.WriteJSON(e)
				ww.wsWriteLock.Unlock()

				if err != nil {
					return
				}
			}
		}
	}
}
