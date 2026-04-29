// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package otelx

import (
	"cmp"
	"net/http"
	"strings"

	"github.com/ory/x/httprouterx"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var withDefaultFilters = otelhttp.WithFilter(func(r *http.Request) bool {
	return !(strings.HasPrefix(r.URL.Path, "/health") ||
		strings.HasPrefix(r.URL.Path, "/admin/health") ||
		strings.HasPrefix(r.URL.Path, "/metrics") ||
		strings.HasPrefix(r.URL.Path, "/admin/metrics"))
})

func NewMiddleware(next http.Handler, operation string, opts ...otelhttp.Option) http.Handler {
	myOpts := []otelhttp.Option{
		withDefaultFilters,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return cmp.Or(r.Pattern, operation, r.Method+" "+r.URL.Path)
		}),
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback := func(r2 *http.Request) {
			r.Pattern = cmp.Or(r.Pattern, r2.Pattern)
		}
		r2 := httprouterx.WithAfterMatchHook(r, callback)
		next.ServeHTTP(w, r2)
		r.Pattern = cmp.Or(r2.Pattern, r.Pattern) // best-effort in case callback never is called
	})
	return otelhttp.NewHandler(handler, operation, append(myOpts, opts...)...)
}
