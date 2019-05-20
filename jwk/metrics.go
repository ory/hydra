package jwk

import (
	"context"

	"github.com/ory/x/metricsx"
	"github.com/prometheus/client_golang/prometheus"
	jose "gopkg.in/square/go-jose.v2"
)

type MetricManager struct {
	Manager
	metricsx.Observer
	JWKs prometheus.Gauge
}

func (m *MetricManager) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	return m.Manager.AddKey(ctx, set, key)
}

func (m *MetricManager) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	return m.Manager.AddKeySet(ctx, set, keys)
}

func (m *MetricManager) GetKey(ctx context.Context, set string, kid string) (*jose.JSONWebKeySet, error) {
	return m.Manager.GetKey(ctx, set, kid)
}

func (m *MetricManager) GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error) {
	return m.Manager.GetKeySet(ctx, set)
}

func (m *MetricManager) DeleteKey(ctx context.Context, set string, kid string) error {
	return m.Manager.DeleteKey(ctx, set, kid)
}

func (m *MetricManager) DeleteKeySet(ctx context.Context, set string) error {
	return m.Manager.DeleteKeySet(ctx, set)
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
