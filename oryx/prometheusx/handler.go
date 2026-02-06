// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	MetricsPrometheusPath = "/metrics/prometheus"
)

type muxrouter interface {
	GET(path string, handle http.HandlerFunc)
}

// SetMuxRoutes registers the prometheus handler.
func SetMuxRoutes(mux muxrouter) {
	mux.GET(MetricsPrometheusPath, promhttp.Handler().ServeHTTP)
}
