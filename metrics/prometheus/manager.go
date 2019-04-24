package prometheus

import (
	"log"
	"net/http"

	// "github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

type MetricsManager struct {
	PrometheusMetrics *Metrics
}

func NewMetricsManager(version, hash, buildTime string, r *prometheus.Registry) *MetricsManager {
	return &MetricsManager{
		PrometheusMetrics: NewMetrics(version, hash, buildTime, r),
	}
}

// Main middleware method to collect metrics for Prometheus.
func (pmm *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nrw := negroni.NewResponseWriter(rw)
		log.Println(r.RequestURI)
		next(nrw, r)
	})

	sizeHandler := promhttp.InstrumentHandlerResponseSize(
		pmm.PrometheusMetrics.HTTPResponseSizeBytes,
		nextHandler,
	)
	durationHandler := promhttp.InstrumentHandlerDuration(
		pmm.PrometheusMetrics.HTTPRequestDuration,
		sizeHandler,
	)

	countHandler := promhttp.InstrumentHandlerCounter(
		pmm.PrometheusMetrics.HTTPRequestCount,
		durationHandler,
	)

	countHandler(rw, r)
}
