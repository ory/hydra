package client

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricClients = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "clients",
		Name:      "total",
		Help:      "The current number of clients",
	})
)
