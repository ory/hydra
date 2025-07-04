// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ory/x/logrusx"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/herodot"
	"github.com/ory/x/urlx"
)

func TestWatchWebsocket(t *testing.T) {
	t.Run("case=forwards events", func(t *testing.T) {
		ctx, c, dir, cancel := setup(t)
		defer cancel()

		hook := &test.Hook{}
		l := logrusx.New("", "", logrusx.WithHook(hook))

		fn := filepath.Join(dir, "some.file")
		f, err := os.Create(fn) //#nosec:G304
		require.NoError(t, err)

		url, err := urlx.Parse("file://" + fn)
		require.NoError(t, err)
		t.Log(url)
		handler, err := WatchAndServeWS(ctx, url, herodot.NewJSONWriter(l))
		require.NoError(t, err)
		s := httptest.NewServer(handler)
		defer s.Close()

		u := urlx.ParseOrPanic("ws" + strings.TrimPrefix(s.URL, "http"))
		_, err = WatchWebsocket(ctx, u, c)
		require.NoError(t, err)

		_, err = fmt.Fprint(f, "content here")
		require.NoError(t, err)
		require.NoError(t, f.Close())
		assertChange(t, <-c, "content here", u.String()+fn)

		require.NoError(t, os.Remove(fn))
		assertRemove(t, <-c, u.String()+fn)

		assert.Len(t, hook.Entries, 0, "%+v", hook.Entries)
	})

	t.Run("case=client closes itself on context cancel", func(t *testing.T) {
		ctx1, c, dir, cancel1 := setup(t)
		defer cancel1()

		hook := &test.Hook{}
		l := logrusx.New("", "", logrusx.WithHook(hook))

		fn := filepath.Join(dir, "some.file")

		handler, err := WatchAndServeWS(ctx1, urlx.ParseOrPanic("file://"+fn), herodot.NewJSONWriter(l))
		require.NoError(t, err)
		s := httptest.NewServer(handler)
		defer s.Close()

		ctx2, cancel2 := context.WithCancel(context.Background())
		u := urlx.ParseOrPanic("ws" + strings.TrimPrefix(s.URL, "http"))
		_, err = WatchWebsocket(ctx2, u, c)
		require.NoError(t, err)

		cancel2()

		e, ok := <-c
		assert.False(t, ok, "%#v", e)

		assert.Len(t, hook.Entries, 0, "%+v", hook.Entries)
	})

	t.Run("case=quits client watcher when server connection is closed", func(t *testing.T) {
		ctxClient, c, dir, cancel := setup(t)
		defer cancel()

		hook := &test.Hook{}
		l := logrusx.New("", "", logrusx.WithHook(hook))

		fn := filepath.Join(dir, "some.file")

		ctxServe, cancelServe := context.WithCancel(context.Background())
		handler, err := WatchAndServeWS(ctxServe, urlx.ParseOrPanic("file://"+fn), herodot.NewJSONWriter(l))
		require.NoError(t, err)
		s := httptest.NewServer(handler)
		defer s.Close()

		u := urlx.ParseOrPanic("ws" + strings.TrimPrefix(s.URL, "http"))
		_, err = WatchWebsocket(ctxClient, u, c)
		require.NoError(t, err)

		cancelServe()

		e, ok := <-c
		assert.False(t, ok, "%#v", e)

		assert.Len(t, hook.Entries, 0, "%+v", hook.Entries)
	})

	t.Run("case=successive watching works after client connection is closed", func(t *testing.T) {
		ctxServer, c, dir, cancel := setup(t)
		defer cancel()

		hook := &test.Hook{}
		l := logrusx.New("", "", logrusx.WithHook(hook))

		fn := filepath.Join(dir, "some.file")

		handler, err := WatchAndServeWS(ctxServer, urlx.ParseOrPanic("file://"+fn), herodot.NewJSONWriter(l))
		require.NoError(t, err)
		s := httptest.NewServer(handler)
		defer s.Close()

		ctxClient1, cancelClient1 := context.WithCancel(context.Background())
		defer cancelClient1()
		u := urlx.ParseOrPanic("ws" + strings.TrimPrefix(s.URL, "http"))
		_, err = WatchWebsocket(ctxClient1, u, c)
		require.NoError(t, err)

		cancelClient1()

		_, ok := <-c
		assert.False(t, ok)

		ctxClient2, cancelClient2 := context.WithCancel(context.Background())
		defer cancelClient2()
		c2 := make(EventChannel)
		_, err = WatchWebsocket(ctxClient2, u, c2)
		require.NoError(t, err)

		f, err := os.Create(fn) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		assertChange(t, <-c2, "", u.String()+fn)

		assert.Len(t, hook.Entries, 0, "%+v", hook.Entries)
	})

	t.Run("case=broadcasts to multiple client connections", func(t *testing.T) {
		ctxServer, c1, dir, cancel := setup(t)
		defer cancel()

		hook := &test.Hook{}
		l := logrusx.New("", "", logrusx.WithHook(hook))

		fn := filepath.Join(dir, "some.file")

		handler, err := WatchAndServeWS(ctxServer, urlx.ParseOrPanic("file://"+fn), herodot.NewJSONWriter(l))
		require.NoError(t, err)
		s := httptest.NewServer(handler)
		defer s.Close()

		ctxClient1, cancelClient1 := context.WithCancel(context.Background())
		defer cancelClient1()

		u := urlx.ParseOrPanic("ws" + strings.TrimPrefix(s.URL, "http"))
		_, err = WatchWebsocket(ctxClient1, u, c1)
		require.NoError(t, err)

		ctxClient2, cancelClient2 := context.WithCancel(context.Background())
		defer cancelClient2()
		c2 := make(EventChannel)
		_, err = WatchWebsocket(ctxClient2, u, c2)
		require.NoError(t, err)

		f, err := os.Create(fn) //#nosec:G304
		require.NoError(t, err)
		require.NoError(t, f.Close())

		assertChange(t, <-c1, "", u.String()+fn)
		assertChange(t, <-c2, "", u.String()+fn)

		assert.Len(t, hook.Entries, 0, "%+v", hook.Entries)
	})

	t.Run("case=sends event when requested", func(t *testing.T) {
		ctxServer, _, dir, cancel := setup(t)
		defer cancel()

		// buffered channel to allow usage of DispatchNow().done
		c := make(EventChannel, 1)

		hook := &test.Hook{}
		l := logrusx.New("", "", logrusx.WithHook(hook))

		fn := filepath.Join(dir, "some.file")
		initialContent := "initial content"
		require.NoError(t, os.WriteFile(fn, []byte(initialContent), 0600))

		handler, err := WatchAndServeWS(ctxServer, urlx.ParseOrPanic("file://"+fn), herodot.NewJSONWriter(l))
		require.NoError(t, err)
		s := httptest.NewServer(handler)
		defer s.Close()

		ctxClient, cancelClient := context.WithCancel(context.Background())
		defer cancelClient()

		u := urlx.ParseOrPanic("ws" + strings.TrimPrefix(s.URL, "http"))
		d, err := WatchWebsocket(ctxClient, u, c)
		require.NoError(t, err)
		done, err := d.DispatchNow()
		require.NoError(t, err)

		// wait for d.DispatchNow to be done
		select {
		case <-time.After(time.Second):
			t.Logf("Waiting for done timed out. %+v", <-c)
			t.FailNow()
		case eventsSend := <-done:
			assert.Equal(t, 1, eventsSend)
		}

		assertChange(t, <-c, initialContent, u.String()+fn)

		assert.Len(t, hook.Entries, 0, "%+v", hook.Entries)
	})
}
