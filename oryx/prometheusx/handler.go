// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ory/herodot"
)

const (
	MetricsPrometheusPath = "/metrics/prometheus"
)

// Handler handles HTTP requests to health and version endpoints.
type Handler struct {
	H             herodot.Writer
	VersionString string
}

// NewHandler instantiates a handler.
func NewHandler(
	h herodot.Writer,
	version string,
) *Handler {
	return &Handler{
		H:             h,
		VersionString: version,
	}
}

type router interface {
	GET(path string, handle httprouter.Handle)
}

// SetRoutes registers this handler's routes.
func (h *Handler) SetRoutes(r router) {
	r.GET(MetricsPrometheusPath, h.Metrics)
}

type muxrouter interface {
	GET(path string, handle http.HandlerFunc)
}

// SetMuxRoutes registers this handler's routes on a ServeMux.
func (h *Handler) SetMuxRoutes(mux muxrouter) {
	mux.GET(MetricsPrometheusPath, promhttp.Handler().ServeHTTP)
}

// Metrics outputs prometheus metrics
//
// swagger:route GET /metrics/prometheus metadata prometheus
//
// Get snapshot metrics from the service. If you're using k8s, you can then add annotations to
// your deployment like so:
//
// ```
// metadata:
//
//	annotations:
//	  prometheus.io/port: "4434"
//	    prometheus.io/path: "/metrics/prometheus"
//
// ```
//
//	Produces:
//	- plain/text
//
//	Responses:
//	  200: emptyResponse
func (h *Handler) Metrics(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	promhttp.Handler().ServeHTTP(rw, r)
}
