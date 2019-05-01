package consent

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ConsentRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "requests",
		Help:      "incremented when a request is sent to consent.AcceptConsentRequest",
	}, []string{"error"})
	ConsentRequestScopes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "request_scopes",
		Help:      "tracks the number of consent requests submitted per scope",
	}, []string{"scope", "type"})
	ConsentRequestAudiences = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "request_audiences",
		Help:      "tracks the number of consent requests submitted per audience",
	}, []string{"audience", "type"})
	ConsentRequestsRejected = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "consent",
		Name:      "requests_rejected",
		Help:      "incremented when consent.RejectConsentRequest is successful",
	}, []string{"error"})

	LoginRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "login",
		Name:      "requests",
		Help:      "incremented when a request is sent to consent.AcceptLoginRequest",
	}, []string{"remember", "error"})
	LoginRequestsRejected = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "hydra",
		Subsystem: "login",
		Name:      "requests_rejected",
		Help:      "incremented when a request to consent.RejectLoginRequest is successful",
	}, []string{"error"})
)

func collectCounters(c chan<- prometheus.Metric) {
	ConsentRequests.Collect(c)
	ConsentRequestScopes.Collect(c)
	ConsentRequestAudiences.Collect(c)
	ConsentRequestsRejected.Collect(c)
	LoginRequests.Collect(c)
	LoginRequestsRejected.Collect(c)
}

func describeCounters(c chan<- *prometheus.Desc) {
	ConsentRequests.Describe(c)
	ConsentRequestScopes.Describe(c)
	ConsentRequestAudiences.Describe(c)
	ConsentRequestsRejected.Describe(c)
	LoginRequests.Describe(c)
	LoginRequestsRejected.Describe(c)
}

func recordAcceptConsentRequest(
	cr *HandledConsentRequest,
	err error,
) {
	ConsentRequests.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()

	for _, v := range cr.GrantedScope {
		ConsentRequestScopes.With(prometheus.Labels{
			"scope": v,
			"type":  "granted",
		}).Inc()
	}

	for _, v := range cr.GrantedAudience {
		ConsentRequestAudiences.With(prometheus.Labels{
			"audience": v,
			"type":     "granted",
		}).Inc()
	}

	if cr.ConsentRequest != nil {
		for _, v := range cr.ConsentRequest.RequestedScope {
			ConsentRequestScopes.With(prometheus.Labels{
				"scope": v,
				"type":  "requested",
			}).Inc()
		}
		for _, v := range cr.ConsentRequest.RequestedAudience {
			ConsentRequestAudiences.With(prometheus.Labels{
				"audience": v,
				"type":     "requested",
			}).Inc()
		}
	}
}

func recordAcceptLoginRequest(lr *HandledLoginRequest, err error) {
	LoginRequests.With(prometheus.Labels{
		"error":    strconv.FormatBool(err != nil),
		"remember": strconv.FormatBool(lr.Remember),
	}).Inc()
}

func recordRejectLoginRequest(err error) {
	LoginRequestsRejected.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()
}

func recordRejectConsentRequest(err error) {
	ConsentRequestsRejected.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()
}
