package jwk

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricJWKs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "jwks",
		Name:      "total",
		Help:      "The number of JWKs issued",
	}, []string{})
)
