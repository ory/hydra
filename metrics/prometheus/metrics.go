package prometheus

// Metrics prototypes
// Example:
// 	Counter      *prometheus.CounterVec
// 	ResponseTime *prometheus.HistogramVec
type Metrics struct{}

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
func NewMetrics(version, hash, buildTime string) *Metrics {
	pm := &Metrics{}
	return pm
}
