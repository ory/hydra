package consent

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricConsentRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "requests",
		Help:      "incremented when a request is sent to consent.AcceptConsentRequest",
	}, []string{"error"})
	metricConsentRequestScopes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "request_scopes",
		Help:      "tracks the number of consent requests submitted per scope",
	}, []string{"scope", "type"})
	metricConsentRequestAudiences = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "request_audiences",
		Help:      "tracks the number of consent requests submitted per audience",
	}, []string{"audience", "type"})
	metricConsentRequestsRejected = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "requests_rejected",
		Help:      "incremented when consent.RejectConsentRequest is successful",
	}, []string{"error"})

	metricLoginRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "login",
		Name:      "requests",
		Help:      "incremented when a request is sent to consent.AcceptLoginRequest",
	}, []string{"remember", "error"})
	metricLoginRequestsRejected = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "login",
		Name:      "requests_rejected",
		Help:      "incremented when a request to consent.RejectLoginRequest is successful",
	}, []string{"error"})

	metricAccessTokens = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "access",
		Name:      "tokens_total",
	})
	metricAccessTokensRevoked = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "access",
		Name:      "tokens_revoked",
	})
	metricRefreshTokens = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "refresh",
		Name:      "tokens_created",
	})
	metricRefreshTokensRevoked = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "hydra",
		Subsystem: "refresh",
		Name:      "tokens_total",
	})
)

func collectCounters(c chan<- prometheus.Metric) {
	metricConsentRequests.Collect(c)
	metricConsentRequestScopes.Collect(c)
	metricConsentRequestAudiences.Collect(c)
	metricConsentRequestsRejected.Collect(c)
	metricLoginRequests.Collect(c)
	metricLoginRequestsRejected.Collect(c)

	metricRefreshTokensRevoked.Collect(c)
	metricAccessTokensRevoked.Collect(c)
}

func describeCounters(c chan<- *prometheus.Desc) {
	metricConsentRequests.Describe(c)
	metricConsentRequestScopes.Describe(c)
	metricConsentRequestAudiences.Describe(c)
	metricConsentRequestsRejected.Describe(c)
	metricLoginRequests.Describe(c)
	metricLoginRequestsRejected.Describe(c)

	metricRefreshTokensRevoked.Describe(c)
	metricAccessTokensRevoked.Describe(c)
}

func recordAcceptConsentRequest(
	cr *HandledConsentRequest,
	err error,
) {
	metricConsentRequests.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()

	for _, v := range cr.GrantedScope {
		metricConsentRequestScopes.With(prometheus.Labels{
			"scope": v,
			"type":  "granted",
		}).Inc()
	}

	for _, v := range cr.GrantedAudience {
		metricConsentRequestAudiences.With(prometheus.Labels{
			"audience": v,
			"type":     "granted",
		}).Inc()
	}

	if cr.ConsentRequest != nil {
		for _, v := range cr.ConsentRequest.RequestedScope {
			metricConsentRequestScopes.With(prometheus.Labels{
				"scope": v,
				"type":  "requested",
			}).Inc()
		}
		for _, v := range cr.ConsentRequest.RequestedAudience {
			metricConsentRequestAudiences.With(prometheus.Labels{
				"audience": v,
				"type":     "requested",
			}).Inc()
		}
	}
}

func recordAcceptLoginRequest(lr *HandledLoginRequest, err error) {
	metricLoginRequests.With(prometheus.Labels{
		"error":    strconv.FormatBool(err != nil),
		"remember": strconv.FormatBool(lr.Remember),
	}).Inc()
}

func recordRejectLoginRequest(err error) {
	metricLoginRequestsRejected.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()
}

func recordRejectConsentRequest(err error) {
	metricConsentRequestsRejected.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()
}
