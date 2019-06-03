package client

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/x/metricsx"
)

type Metrics struct {
	prometheus.Collector
	metricsx.Observer
	Manager

	Clients              prometheus.Gauge
	ClientsCreated       prometheus.Counter
	ClientsAuthenticated metricsx.CounterVec
}

func (m *Metrics) Authenticate(ctx context.Context, id string, secret []byte) (*Client, error) {
	c, err := m.Manager.Authenticate(ctx, id, secret)
	if err == nil {
		m.ClientsAuthenticated.With(prometheus.Labels{
			"client": strings.ToLower(c.Name),
		}).Inc()
	}

	return c, err
}

func (m *Metrics) CreateClient(ctx context.Context, c *Client) error {
	err := m.Manager.CreateClient(ctx, c)
	if err == nil {
		m.ClientsCreated.Inc()
	}

	return err
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent. The
// descriptor of each sent metric is one of those returned by Describe
// (unless the Collector is unchecked, see above). Returned metrics that
// share the same descriptor must differ in their variable label
// values.
//
// This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. Blocking occurs at the expense
// of total performance of rendering all registered metrics. Ideally,
// Collector implementations support concurrent readers.
func (m *Metrics) Collect(c chan<- prometheus.Metric) {
	m.Clients.Collect(c)
	m.ClientsCreated.Collect(c)
	m.ClientsAuthenticated.Collect(c)
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent. The sent descriptors fulfill the
// consistency and uniqueness requirements described in the Desc
// documentation.
//
// It is valid if one and the same Collector sends duplicate
// descriptors. Those duplicates are simply ignored. However, two
// different Collectors must not send duplicate descriptors.
//
// Sending no descriptor at all marks the Collector as “unchecked”,
// i.e. no checks will be performed at registration time, and the
// Collector may yield any Metric it sees fit in its Collect method.
//
// This method idempotently sends the same descriptors throughout the
// lifetime of the Collector. It may be called concurrently and
// therefore must be implemented in a concurrency safe way.
//
// If a Collector encounters an error while executing this method, it
// must send an invalid descriptor (created with NewInvalidDesc) to
// signal the error to the registry.
func (m *Metrics) Describe(c chan<- *prometheus.Desc) {
	m.Clients.Describe(c)
	m.ClientsCreated.Describe(c)
	m.ClientsAuthenticated.Describe(c)
}

func (m *Metrics) Observe() error {
	n, err := m.CountClients(context.Background())
	if err != nil {
		return err
	}
	m.Clients.Set(float64(n))

	return nil
}

func WithMetrics(m Manager) *Metrics {
	return &Metrics{
		Manager: m,
		Clients: metricsx.NewGauge(prometheus.GaugeOpts{
			Namespace: "hydra",
			Subsystem: "clients",
			Name:      "total",
			Help:      "The current number of clients",
		}),
		ClientsCreated: metricsx.NewCounter(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "clients",
			Name:      "created_sum",
			Help:      "The running total of clients created",
		}),
		ClientsAuthenticated: metricsx.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "clients",
			Name:      "authenticated_sum",
			Help:      "The running total of successfully authenticated clients",
		}, []string{"client"}),
	}
}
