package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

// Prometheus contains the metric registry and HTTP middleware for tracking HTTP requests
type Prometheus struct {
	Registry *prometheus.Registry

	HTTPRequestCount      *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPResponseSizeBytes *prometheus.HistogramVec
}

// NewMetrics registers our metrics and their labels to the prometheus Registry.
func NewPrometheus(registry *prometheus.Registry, collectors ...prometheus.Collector) *Prometheus {
	pm := &Prometheus{
		Registry: registry,
		HTTPRequestCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "",
		}, []string{"code", "method"}),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "A histogram of latencies for requests",
				Buckets: []float64{.05, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method"},
		),
		HTTPResponseSizeBytes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "A histogram of response sizes for requests",
				Buckets: []float64{200, 400, 800, 1600},
			},
			[]string{},
		),
	}

	if pm.Registry != nil {
		pm.Registry.MustRegister(
			prometheus.NewGoCollector(),
			pm.HTTPRequestCount,
			pm.HTTPRequestDuration,
			pm.HTTPResponseSizeBytes,
		)
		for _, v := range collectors {
			pm.Registry.MustRegister(v)
		}
	}
	return pm
}

func (p *Prometheus) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nrw := negroni.NewResponseWriter(rw)
		next(nrw, r)
	})

	sizeHandler := promhttp.InstrumentHandlerResponseSize(
		p.HTTPResponseSizeBytes,
		nextHandler,
	)

	durationHandler := promhttp.InstrumentHandlerDuration(
		p.HTTPRequestDuration,
		sizeHandler,
	)

	countHandler := promhttp.InstrumentHandlerCounter(
		p.HTTPRequestCount,
		durationHandler,
	)

	countHandler(rw, r)
}
