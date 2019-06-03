package consent

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/x/metricsx"
)

type MetricsManager struct {
	Manager
	metricsx.Observer
	ConsentRequestsHandled  metricsx.CounterVec
	ConsentRequestScopes    metricsx.CounterVec
	ConsentRequestAudiences metricsx.CounterVec
	LoginRequestsHandled    metricsx.CounterVec
	LogoutRequests          metricsx.CounterVec
	AccessTokens            prometheus.Gauge
	RefreshTokens           prometheus.Gauge
}

func (m *MetricsManager) HandleConsentRequest(ctx context.Context, challenge string, r *HandledConsentRequest) (*ConsentRequest, error) {
	c, err := m.Manager.HandleConsentRequest(ctx, challenge, r)
	if r.Error != nil {
		// RejectConsentRequest was called
		m.ConsentRequestsHandled.With(prometheus.Labels{
			"rejected": "true",
		}).Inc()
	} else {
		for _, v := range r.GrantedScope {
			m.ConsentRequestScopes.With(prometheus.Labels{
				"scope": v,
				"type":  "granted",
			}).Inc()
		}

		for _, v := range r.GrantedAudience {
			m.ConsentRequestScopes.With(prometheus.Labels{
				"scope": v,
				"type":  "granted",
			}).Inc()
		}
	}

	if val := r.ConsentRequest; val != nil {
		for _, v := range val.RequestedScope {
			m.ConsentRequestScopes.With(prometheus.Labels{
				"scope": v,
				"type":  "requested",
			}).Inc()
		}
		for _, v := range val.RequestedAudience {
			m.ConsentRequestAudiences.With(prometheus.Labels{
				"audience": v,
				"type":     "requested",
			}).Inc()
		}
	}

	m.ConsentRequestsHandled.With(prometheus.Labels{
		"rejected": "false",
	}).Inc()
	return c, err
}

func (m *MetricsManager) HandleLoginRequest(ctx context.Context, challenge string, r *HandledLoginRequest) (*LoginRequest, error) {
	l, err := m.Manager.HandleLoginRequest(ctx, challenge, r)
	if r.Error != nil {
		m.LoginRequestsHandled.With(prometheus.Labels{
			"remember": strconv.FormatBool(r.Remember),
			"rejected": "true",
		}).Inc()
	} else {
		m.LoginRequestsHandled.With(prometheus.Labels{
			"remember": strconv.FormatBool(r.Remember),
			"rejected": "false",
		}).Inc()
	}
	return l, err
}

func (m *MetricsManager) AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	r, err := m.Manager.AcceptLogoutRequest(ctx, challenge)
	m.LogoutRequests.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()

	return r, err
}

func (m *MetricsManager) Describe(c chan<- *prometheus.Desc) {
	m.ConsentRequestsHandled.Describe(c)
	m.ConsentRequestScopes.Describe(c)
	m.ConsentRequestAudiences.Describe(c)
	m.LoginRequestsHandled.Describe(c)
	m.LogoutRequests.Describe(c)
}

func (m *MetricsManager) Collect(c chan<- prometheus.Metric) {
	m.ConsentRequestsHandled.Collect(c)
	m.ConsentRequestScopes.Collect(c)
	m.ConsentRequestAudiences.Collect(c)
	m.LoginRequestsHandled.Collect(c)
	m.LogoutRequests.Collect(c)
}

func (m *MetricsManager) Observe() error {
	accessTokens, err := m.CountAccessTokens()
	if err != nil {
		return err
	}

	m.AccessTokens.Set(float64(accessTokens))

	refreshTokens, err := m.CountRefreshTokens()
	if err != nil {
		return err
	}

	m.RefreshTokens.Set(float64(refreshTokens))

	return nil
}

func WithMetrics(m Manager) *MetricsManager {
	return &MetricsManager{
		Manager: m,
		ConsentRequestsHandled: metricsx.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "requests_handled",
			Help:      "",
		}, []string{"rejected"}),
		ConsentRequestScopes: metricsx.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "request_scopes",
			Help:      "tracks the number of consent requests submitted per scope",
		}, []string{"scope", "type"}),
		ConsentRequestAudiences: metricsx.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "consent",
			Name:      "request_audiences",
			Help:      "tracks the number of consent requests submitted per audience",
		}, []string{"audience", "type"}),
		LoginRequestsHandled: metricsx.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "login",
			Name:      "requests",
			Help:      "",
		}, []string{"remember", "rejected"}),
		LogoutRequests: metricsx.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hydra",
			Subsystem: "logout",
			Name:      "requests",
			Help:      "incremented when logout request is made",
		}, []string{"error"}),
		AccessTokens: metricsx.NewGauge(prometheus.GaugeOpts{
			Namespace: "hydra",
			Subsystem: "access",
			Name:      "tokens_total",
		}),
		RefreshTokens: metricsx.NewGauge(prometheus.GaugeOpts{
			Namespace: "hydra",
			Subsystem: "refresh",
			Name:      "tokens_total",
		}),
	}
}
