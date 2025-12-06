// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package otelx

import (
	"cmp"
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var WithDefaultFilters otelhttp.Option = otelhttp.WithFilter(func(r *http.Request) bool {
	return !(strings.HasPrefix(r.URL.Path, "/health") ||
		strings.HasPrefix(r.URL.Path, "/admin/health") ||
		strings.HasPrefix(r.URL.Path, "/metrics") ||
		strings.HasPrefix(r.URL.Path, "/admin/metrics"))
})

type contextKey int

const callbackContextKey contextKey = iota

func SpanNameRecorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			cb, _ := r.Context().Value(callbackContextKey).(func(string))
			if cb == nil {
				return
			}
			if r.Pattern != "" {
				cb(r.Pattern)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func SpanNameRecorderNegroniFunc(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		cb, _ := r.Context().Value(callbackContextKey).(func(string))
		if cb == nil {
			return
		}
		if r.Pattern != "" {
			cb(r.Pattern)
		}
	}()
	next(w, r)
}

func NewMiddleware(next http.Handler, operation string, opts ...otelhttp.Option) http.Handler {
	myOpts := []otelhttp.Option{
		WithDefaultFilters,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return cmp.Or(r.Pattern, operation, r.Method+" "+r.URL.Path)
		}),
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback := func(s string) {
			r.Pattern = cmp.Or(r.Pattern, s)
		}
		ctx := context.WithValue(r.Context(), callbackContextKey, callback)
		r2 := r.WithContext(ctx)
		next.ServeHTTP(w, r2)
		r.Pattern = cmp.Or(r2.Pattern, r.Pattern) // best-effort in case callback never is called
	})
	return otelhttp.NewHandler(handler, operation, append(myOpts, opts...)...)
}
