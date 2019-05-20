package consent

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ory/hydra/client"
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

func (m *MetricsManager) CreateConsentRequest(ctx context.Context, req *ConsentRequest) error {
	return m.Manager.CreateConsentRequest(ctx, req)
}

func (m *MetricsManager) GetConsentRequest(ctx context.Context, challenge string) (*ConsentRequest, error) {
	return m.Manager.GetConsentRequest(ctx, challenge)
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

func (m *MetricsManager) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	return m.Manager.RevokeSubjectConsentSession(ctx, user)
}

func (m *MetricsManager) RevokeSubjectClientConsentSession(ctx context.Context, user string, client string) error {
	return m.Manager.RevokeSubjectClientConsentSession(ctx, user, client)
}

func (m *MetricsManager) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*HandledConsentRequest, error) {
	return m.Manager.VerifyAndInvalidateConsentRequest(ctx, verifier)
}

func (m *MetricsManager) FindGrantedAndRememberedConsentRequests(ctx context.Context, client string, user string) ([]HandledConsentRequest, error) {
	return m.Manager.FindGrantedAndRememberedConsentRequests(ctx, client, user)
}

func (m *MetricsManager) FindSubjectsGrantedConsentRequests(ctx context.Context, user string, limit int, offset int) ([]HandledConsentRequest, error) {
	return m.Manager.FindSubjectsGrantedConsentRequests(ctx, user, limit, offset)
}

func (m *MetricsManager) CountSubjectsGrantedConsentRequests(ctx context.Context, user string) (int, error) {
	return m.Manager.CountSubjectsGrantedConsentRequests(ctx, user)
}

func (m *MetricsManager) GetRememberedLoginSession(ctx context.Context, id string) (*LoginSession, error) {
	return m.Manager.GetRememberedLoginSession(ctx, id)
}

func (m *MetricsManager) CreateLoginSession(ctx context.Context, session *LoginSession) error {
	return m.Manager.CreateLoginSession(ctx, session)
}

func (m *MetricsManager) DeleteLoginSession(ctx context.Context, id string) error {
	return m.Manager.DeleteLoginSession(ctx, id)
}

func (m *MetricsManager) RevokeSubjectLoginSession(ctx context.Context, user string) error {
	return m.Manager.RevokeSubjectLoginSession(ctx, user)
}

func (m *MetricsManager) ConfirmLoginSession(ctx context.Context, id string, subject string, remember bool) error {
	return m.Manager.ConfirmLoginSession(ctx, id, subject, remember)
}

func (m *MetricsManager) CreateLoginRequest(ctx context.Context, req *LoginRequest) error {
	return m.Manager.CreateLoginRequest(ctx, req)
}

func (m *MetricsManager) GetLoginRequest(ctx context.Context, challenge string) (*LoginRequest, error) {
	return m.Manager.GetLoginRequest(ctx, challenge)
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

func (m *MetricsManager) VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*HandledLoginRequest, error) {
	return m.Manager.VerifyAndInvalidateLoginRequest(ctx, verifier)
}

func (m *MetricsManager) CreateForcedObfuscatedLoginSession(ctx context.Context, session *ForcedObfuscatedLoginSession) error {
	return m.Manager.CreateForcedObfuscatedLoginSession(ctx, session)
}

func (m *MetricsManager) GetForcedObfuscatedLoginSession(ctx context.Context, client string, obfuscated string) (*ForcedObfuscatedLoginSession, error) {
	return m.Manager.GetForcedObfuscatedLoginSession(ctx, client, obfuscated)
}

func (m *MetricsManager) ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject string) ([]client.Client, error) {
	return m.Manager.ListUserAuthenticatedClientsWithFrontChannelLogout(ctx, subject)
}

func (m *MetricsManager) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject string) ([]client.Client, error) {
	return m.Manager.ListUserAuthenticatedClientsWithBackChannelLogout(ctx, subject)
}

func (m *MetricsManager) CreateLogoutRequest(ctx context.Context, request *LogoutRequest) error {
	return m.Manager.CreateLogoutRequest(ctx, request)
}

func (m *MetricsManager) GetLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	return m.Manager.GetLogoutRequest(ctx, challenge)
}

func (m *MetricsManager) AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	r, err := m.Manager.AcceptLogoutRequest(ctx, challenge)
	m.LogoutRequests.With(prometheus.Labels{
		"error": strconv.FormatBool(err != nil),
	}).Inc()

	return r, err
}

func (m *MetricsManager) RejectLogoutRequest(ctx context.Context, challenge string) error {
	return m.Manager.RejectLogoutRequest(ctx, challenge)
}

func (m *MetricsManager) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*LogoutRequest, error) {
	return m.Manager.VerifyAndInvalidateLogoutRequest(ctx, verifier)
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
