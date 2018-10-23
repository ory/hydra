package tracing

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/urfave/negroni"
)

func (t *Tracer) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var span opentracing.Span
	opName := r.URL.Path

	// It's very possible that Hydra is fronted by a proxy which could have initiated a trace.
	// If so, we should attempt to join it.
	remoteContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	if err != nil {
		span = opentracing.StartSpan(opName)
	} else {
		span = opentracing.StartSpan(opName, opentracing.ChildOf(remoteContext))
	}

	defer span.Finish()

	r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))

	next(rw, r)

	ext.HTTPMethod.Set(span, r.Method)
	if negroniWriter, ok := rw.(negroni.ResponseWriter); ok {
		statusCode := uint16(negroniWriter.Status())
		if statusCode >= 400 {
			ext.Error.Set(span, true)
		}
		ext.HTTPStatusCode.Set(span, statusCode)
	}
}
