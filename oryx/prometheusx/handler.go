// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	MetricsPrometheusPath = "/metrics/prometheus"
)

type muxrouter interface {
	Handle(path string, handle http.Handler)
}

// Handler returns the metrics handler. It mirrors promhttp.Handler() but
// additionally serves the OpenMetrics format when the scraper negotiates it,
// because that is the only exposition format that carries exemplars (trace
// links on counters and histograms). It is safe to call multiple times.
func Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}),
	)
}

// SetMuxRoutes registers the prometheus handler.
func SetMuxRoutes(mux muxrouter) {
	mux.Handle("GET "+MetricsPrometheusPath, Handler())
}
