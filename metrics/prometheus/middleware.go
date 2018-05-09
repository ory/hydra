package prometheus

import (
	"net/http"
)

type MetricsManager struct {
	prometheusMetrics *Metrics
}

func NewMetricsManager(version, hash, buildTime string) *MetricsManager {
	return &MetricsManager{
		prometheusMetrics: NewMetrics(version, hash, buildTime),
	}
}

// Main middleware method to collect metrics for Prometheus.
// Example:
//	start := time.Now()
//	next(rw, r)
//	Request counter metric
//	pmm.prometheusMetrics.Counter.WithLabelValues(r.URL.Path).Inc()
//	Response time metric
//	pmm.prometheusMetrics.ResponseTime.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())
func (pmm *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}
