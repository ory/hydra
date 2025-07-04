// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
)

func TestContextualizedMiddleware(t *testing.T) {
	createServer := func(t *testing.T, cb func(ctx context.Context) (cors.Options, bool)) *httptest.Server {
		n := negroni.New()
		n.UseFunc(ContextualizedMiddleware(cb))
		n.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			_, _ = rw.Write([]byte("ok"))
		})
		ts := httptest.NewServer(n)
		t.Cleanup(ts.Close)
		return ts
	}

	fetchCORS := func(t *testing.T, origin string, ts *httptest.Server) http.Header {
		req, err := http.NewRequest("OPTIONS", ts.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Origin", origin)
		req.Header.Set("Access-Control-Request-Method", "DELETE")
		req.Header.Set("Access-Control-Request-Headers", "")
		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer res.Body.Close()
		return res.Header
	}

	t.Run("switches enabled on and off", func(t *testing.T) {
		var enabled bool
		var origins []string
		ts := createServer(t, func(ctx context.Context) (cors.Options, bool) {
			return cors.Options{
				AllowedMethods: []string{"OPTIONS", "DELETE"},
				AllowedOrigins: origins,
				Debug:          true,
			}, enabled
		})

		origins = append(origins, "http://localhost:8080")
		actual := fetchCORS(t, "http://localhost:8080", ts)
		assert.Empty(t, actual.Get("Access-Control-Allow-Origin"))

		enabled = true
		actual = fetchCORS(t, "http://localhost:8080", ts)
		assert.Equal(t, "http://localhost:8080", actual.Get("Access-Control-Allow-Origin"), actual)

		enabled = false
		actual = fetchCORS(t, "http://localhost:8080", ts)
		assert.Empty(t, actual.Get("Access-Control-Allow-Origin"))

		enabled = true
		origins = []string{"http://localhost:9090"}
		actual = fetchCORS(t, "http://localhost:8080", ts)
		assert.Empty(t, actual.Get("Access-Control-Allow-Origin"))

		actual = fetchCORS(t, "http://localhost:9090", ts)
		assert.Equal(t, "http://localhost:9090", actual.Get("Access-Control-Allow-Origin"), actual)
	})
}
