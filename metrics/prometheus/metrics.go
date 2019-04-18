package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics contains the metrics that will be recorded.
// Rather than create metrics ad-hoc, we prefer having a central list of metrics,
// as per the Prometheus "best practices" in their documentation.
// See: https://prometheus.io/docs/practices/instrumentation/#avoid-missing-metrics
//
// The labels associated with metrics should be exclusively descriptive.
// For example, one might think that this is an applicable demonstration of "hydra_consent_requests_accepted":
//   hydra_consent_requests_accepted{scope="scope:1"} 1
//   hydra_consent_requests_accepted{scope="scope:2"} 1
// when a consent request requests []string{"scope:1", "scope:2"}.
// This is not good; if we were to query Prometheus for `sum(hydra_consent_requests_accepted)`, we would get 2, when in reality,
// only 1 consent request was accepted. Instead, separate data like "scopes" into a more descriptive time series:
//   hydra_consent_requests_accepted{scopes="2"} 1
//   hydra_consent_requests_accepted_scopes{scope="scope:1"} 1
//   hydra_consent_requests_accepted_scopes{scope="scope:2"} 1
//
// Some of these values will be populated by HTTP middleware, while others
// will be populated elsewhere throughout the project
type Metrics struct {
	HTTPRequestCount      *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPResponseSizeBytes *prometheus.HistogramVec

	ConsentRequests         *prometheus.CounterVec
	ConsentRequestScopes    *prometheus.CounterVec
	ConsentRequestAudiences *prometheus.CounterVec
	ConsentRequestsRejected *prometheus.CounterVec

	LoginRequests *prometheus.CounterVec

	AccessTokensIssued   *prometheus.CounterVec
	AccessTokensRevoked  *prometheus.CounterVec
	RefreshTokensIssued  *prometheus.CounterVec
	RefreshTokensRevoked *prometheus.CounterVec
	IDTokensIssued       *prometheus.CounterVec

	Grants  *prometheus.GaugeVec
	Clients *prometheus.GaugeVec
	Keys    *prometheus.GaugeVec
}

// NewMetrics registers our metrics and their labels to the prometheus Registry.
func NewMetrics(version, hash, date string) *Metrics {
	pm := &Metrics{
		HTTPRequestCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "",
		}, []string{"code", "method"}),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "A histogram of latencies for requests",
				Buckets: []float64{.05, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method"},
		),
		HTTPResponseSizeBytes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "A histogram of response sizes for requests",
				Buckets: []float64{200, 400, 800, 1600},
			},
			[]string{},
		),
		ConsentRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "requests",
			Help:      "incremented when a request is sent to consent.AcceptConsentRequest",
		}, []string{"scopes", "audiences", "error"}),
		ConsentRequestScopes: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "request_scopes",
			Help:      "tracks the number of consent requests submitted per scope",
		}, []string{"scope"}),
		ConsentRequestAudiences: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "request_audiences",
			Help:      "tracks the number of consent requests submitted per audience",
		}, []string{"audience"}),
		ConsentRequestsRejected: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "requests_rejected",
			Help:      "incremented when consent.RejectConsentRequest is successful",
		}, []string{}),
	}

	prometheus.MustRegister(
		pm.HTTPRequestCount,
		pm.HTTPRequestDuration,
		pm.HTTPResponseSizeBytes,
		pm.ConsentRequests,
		pm.ConsentRequestScopes,
		pm.ConsentRequestAudiences,
		pm.ConsentRequestsRejected,
	)
	return pm
}

// IncrementCounterWithLabels increments a single metric with multiple label combinations by 1.
func IncrementCounterWithLabels(c prometheus.CounterVec, labels map[string][]string) {
	IncreaseCounterWithLabels(c, labels, float64(1))
}

// IncreaseCounterWithLabels increases a single metric with multiple label combinations by "value". Be cautious of using an oversized map in this function. High degrees of cardinality can slow down Prometheus queries significantly.
// Examples:
//    IncreaseCounterWithLabels(c, map[string][]string{"accepted": []string{"true"}, "scope": []string{"scope:1", "scope:2"}}, 4)
// will produce:
//    metric_name{accepted="true", scope="scope:1"}    4
//    metric_name{accepted="true", scope="scope:2"}    4
// ---
//    IncreaseCounterWithLabels(c, map[string][]string{"accepted": []string{"true"}, "scope": []string{"scope:1", "scope:2"}, "other_value": []string{"other", "label", "value"}}, 4)
// will produce:
//    metric_name{accepted="true", scope="scope:1", other_value="other"}    4
//    metric_name{accepted="true", scope="scope:2", other_value="other"}    4
//    metric_name{accepted="true", scope="scope:1", other_value="label"}    4
//    metric_name{accepted="true", scope="scope:2", other_value="label"}    4
//    metric_name{accepted="true", scope="scope:1", other_value="value"}    4
//    metric_name{accepted="true", scope="scope:2", other_value="value"}    4
func IncreaseCounterWithLabels(c prometheus.CounterVec, labels map[string][]string, value float64) {
	for k, values := range labels {
		for _, v := range values {
			c.With(prometheus.Labels{
				k: v,
			}).Add(value)
		}
	}
}

// SetGaugeWithLabels increases a single metric with multiple label combinations by "value".
func SetGaugeWithLabels(c prometheus.GaugeVec, labels map[string][]string, value float64) {
	for k, values := range labels {
		for _, v := range values {
			c.With(prometheus.Labels{
				k: v,
			}).Set(value)
		}
	}
}
