package jwk

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/x/metricsx"
)

type MetricManager struct {
	Manager
	metricsx.Observer
	JWKs prometheus.Gauge
}

func (m *MetricManager) Describe(c chan<- *prometheus.Desc) {
	m.JWKs.Describe(c)
}

func (m *MetricManager) Collect(c chan<- prometheus.Metric) {
	m.JWKs.Collect(c)
}

func (m *MetricManager) Observe() error {
	n, err := m.CountJWKs(context.Background())
	if err != nil {
		return err
	}

	m.JWKs.Set(float64(n))
	return nil
}

func WithMetrics(m Manager) *MetricManager {
	return &MetricManager{
		Manager: m,
		JWKs: metricsx.NewGauge(prometheus.GaugeOpts{
			Namespace: "hydra",
			Subsystem: "jwks",
			Name:      "total",
			Help:      "The number of JWKs issued",
		}),
	}
}
