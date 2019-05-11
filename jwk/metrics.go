package jwk

import (
	"github.com/ory/x/metricsx"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricJWKs = metricsx.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "jwks",
		Name:      "total",
		Help:      "The number of JWKs issued",
	}, []string{})
)
