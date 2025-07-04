// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import (
	"context"
	"net/http"

	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// ContextualizedMiddleware is a context-aware CORS middleware. It allows hot-reloading CORS configuration using
// the HTTP request context.
//
//	n := negroni.New()
//	n.UseFunc(ContextualizedMiddleware(func(context.Context) (opts cors.Options, enabled bool) {
//	  panic("implement me")
//	})
//	// ...
//
// Deprecated: because this is not really practical to use, you should use CheckOrigin as the cors.Options.AllowOriginRequestFunc instead.
func ContextualizedMiddleware(provider func(context.Context) (opts cors.Options, enabled bool)) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		options, enabled := provider(r.Context())
		if !enabled {
			next(rw, r)
			return
		}

		cors.New(options).Handler(next).ServeHTTP(rw, r)
	}
}
