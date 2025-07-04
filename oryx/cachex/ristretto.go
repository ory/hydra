// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cachex

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// RistrettoCollector collects Ristretto cache metrics.
type RistrettoCollector struct {
	prefix      string
	metricsFunc func() *ristretto.Metrics
}

// NewRistrettoCollector creates a new RistrettoCollector.
//
// To use this collector, you need to register it with a Prometheus registry:
//
//	func main() {
//		cache, _ := ristretto.NewCache(&ristretto.Config{
//			NumCounters: 1e7,
//			MaxCost:     1 << 30,
//			BufferItems: 64,
//		})
//		collector := NewRistrettoCollector("prefix_", func() *ristretto.Metrics {
//			return cache.Metrics
//		})
//		prometheus.MustRegister(collector)
//	}
func NewRistrettoCollector(prefix string, metricsFunc func() *ristretto.Metrics) *RistrettoCollector {
	return &RistrettoCollector{
		prefix:      prefix,
		metricsFunc: metricsFunc,
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector.
func (c *RistrettoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(c.prefix+"ristretto_hits", "Total number of cache hits", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_misses", "Total number of cache misses", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_ratio", "Cache hit ratio", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_keys_added", "Total number of keys added to the cache", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_cost_added", "Total cost of keys added to the cache", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_keys_evicted", "Total number of keys evicted from the cache", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_cost_evicted", "Total cost of keys evicted from the cache", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_sets_dropped", "Total number of sets dropped", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_sets_rejected", "Total number of sets rejected", nil, nil)
	ch <- prometheus.NewDesc(c.prefix+"ristretto_gets_kept", "Total number of gets kept", nil, nil)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *RistrettoCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.metricsFunc()
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_hits", "Total number of cache hits", nil, nil), prometheus.GaugeValue, float64(metrics.Hits()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_misses", "Total number of cache misses", nil, nil), prometheus.GaugeValue, float64(metrics.Misses()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_ratio", "Cache hit ratio", nil, nil), prometheus.GaugeValue, metrics.Ratio())
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_keys_added", "Total number of keys added to the cache", nil, nil), prometheus.GaugeValue, float64(metrics.KeysAdded()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_cost_added", "Total cost of keys added to the cache", nil, nil), prometheus.GaugeValue, float64(metrics.CostAdded()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_keys_evicted", "Total number of keys evicted from the cache", nil, nil), prometheus.GaugeValue, float64(metrics.KeysEvicted()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_cost_evicted", "Total cost of keys evicted from the cache", nil, nil), prometheus.GaugeValue, float64(metrics.CostEvicted()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_sets_dropped", "Total number of sets dropped", nil, nil), prometheus.GaugeValue, float64(metrics.SetsDropped()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_sets_rejected", "Total number of sets rejected", nil, nil), prometheus.GaugeValue, float64(metrics.SetsRejected()))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(c.prefix+"ristretto_gets_kept", "Total number of gets kept", nil, nil), prometheus.GaugeValue, float64(metrics.GetsKept()))
}
