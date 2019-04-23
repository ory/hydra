package consent

import (
	"strconv"

	"github.com/ory/hydra/metrics/prometheus"

	prom "github.com/prometheus/client_golang/prometheus"
)

func recordAcceptConsentRequest(
	m *prometheus.MetricsManager,
	cr *HandledConsentRequest,
	err error,
) {
	m.PrometheusMetrics.ConsentRequests.With(prom.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()

	for _, v := range cr.GrantedScope {
		m.PrometheusMetrics.ConsentRequestScopes.With(prom.Labels{
			"scope": v,
			"type":  "granted",
		}).Inc()
	}

	for _, v := range cr.GrantedAudience {
		m.PrometheusMetrics.ConsentRequestAudiences.With(prom.Labels{
			"audience": v,
			"type":     "granted",
		}).Inc()
	}

	if cr.ConsentRequest != nil {
		for _, v := range cr.ConsentRequest.RequestedScope {
			m.PrometheusMetrics.ConsentRequestScopes.With(prom.Labels{
				"scope": v,
				"type":  "requested",
			}).Inc()
		}
		for _, v := range cr.ConsentRequest.RequestedAudience {
			m.PrometheusMetrics.ConsentRequestAudiences.With(prom.Labels{
				"audience": v,
				"type":     "requested",
			}).Inc()
		}
	}
}

func recordAcceptLoginRequest(
	m *prometheus.MetricsManager,
	lr *HandledLoginRequest,
	err error,
) {
	m.PrometheusMetrics.LoginRequests.With(prom.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()

	m.PrometheusMetrics.LoginRequestSubjects.With(prom.Labels{
		"subject": lr.Subject,
	}).Inc()
}

// func recordRejectConsentRequest(
// 	m *prometheus.MetricManager,
// 	cr *handledConsentRequest,
// ) {
//
// }
