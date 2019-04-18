package consent

import (
	"strconv"

	"github.com/ory/hydra/metrics/prometheus"

	prom "github.com/prometheus/client_golang/prometheus"
)

func recordAcceptConsentRequest(
	m *prometheus.MetricsManager,
	cr *HandledConsentRequest,
	err *error,
) {
	m.PrometheusMetrics.ConsentRequests.With(prom.Labels{
		"scopes":    strconv.Itoa(len(cr.GrantedScope)),
		"audiences": strconv.Itoa(len(cr.GrantedAudience)),
		"error":     strconv.FormatBool(err != nil),
	}).Inc()
}
