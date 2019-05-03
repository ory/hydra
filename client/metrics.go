package client

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Clients = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "clients",
		Name:      "sum",
		Help:      "The current number of clients",
	})
)
