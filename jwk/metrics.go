package jwk

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	JWKs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "jwks",
		Name:      "sum",
		Help:      "The number of JWKs issued",
	}, []string{})
)
