package x

import (
	"net/http"

	"github.com/ory/x/otelx"
)

// TraceHandler wraps otelx.NewHandler, passing the URL path as the span name.
func TraceHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otelx.NewHandler(h, r.URL.Path).ServeHTTP(w, r)
	})
}
