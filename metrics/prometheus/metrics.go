package prometheus

import "github.com/prometheus/client_golang/prometheus"

const (
	MetricsPrometheusPath = "/metrics/prometheus"
)

// Metrics prototypes
// Example:
// 	Counter      *prometheus.CounterVec
// 	ResponseTime *prometheus.HistogramVec
type Metrics struct {
	ResponseTime *prometheus.HistogramVec
}

// Method for creation new custom Prometheus  metrics
// Example:
// 	pm := &Metrics{
//		Counter: prometheus.NewCounterVec(
//			prometheus.CounterOpts{
//				Name:        "servicename_requests_total",
//				Help:        "Description",
//				ConstLabels: map[string]string{
//					"version":   version,
//					"hash":      hash,
//					"buildTime": buildTime,
//				},
//			},
//			[]string{"endpoint"},
//		),
//		ResponseTime: prometheus.NewHistogramVec(
//			prometheus.HistogramOpts{
//				Name:        "servicename_response_time_seconds",
//				Help:        "Description",
//				ConstLabels: map[string]string{
//					"version":   version,
//					"hash":      hash,
//					"buildTime": buildTime,
//				},
//			},
//			[]string{"endpoint"},
//		),
//	}
//	prometheus.Register(pm.Counter)
//  prometheus.Register(pm.ResponseTime)
func NewMetrics(version, hash, date string) *Metrics {
	pm := &Metrics{
		ResponseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "hydra_response_time_seconds",
				Help: "Description",
				ConstLabels: map[string]string{
					"version":   version,
					"hash":      hash,
					"buildTime": date,
				},
			},
			[]string{"endpoint"},
		),
	}
	err := prometheus.Register(pm.ResponseTime)

	if err != nil {
		panic(err)
	}
	return pm
}
