package prometheus

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type MetricsManager struct {
	prometheusMetrics *Metrics
	routers []*httprouter.Router
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
	start := time.Now()
	next(rw, r)

	// looking for a match in one of registered routers
	matched := false
	for _, router := range pmm.routers {
		handler, _, _ := router.Lookup(r.Method, r.URL.Path)
		if handler != nil {
			matched = true
			break
		}
	}

	if matched {
		pmm.prometheusMetrics.ResponseTime.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())
	} else {
		pmm.prometheusMetrics.ResponseTime.WithLabelValues("{unmatched}").Observe(time.Since(start).Seconds())
	}
}

func (pmm *MetricsManager) RegisterRouter(router *httprouter.Router) {
	pmm.routers = append(pmm.routers, router)
}
