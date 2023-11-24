// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/rfc8628"
	"github.com/ory/herodot"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fositex"
	"github.com/ory/hydra/v2/hsm"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/hydra/v2/x/oauth2cors"
	"github.com/ory/x/contextx"
	"github.com/ory/x/healthx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/popx"
	prometheus "github.com/ory/x/prometheusx"
)

var (
	_ contextx.Provider = (*RegistryBase)(nil)
)

type RegistryBase struct {
	l               *logrusx.Logger
	al              *logrusx.Logger
	conf            *config.DefaultProvider
	ch              *client.Handler
	fh              fosite.Hasher
	jwtGrantH       *trust.Handler
	jwtGrantV       *trust.GrantValidator
	kh              *jwk.Handler
	cv              *client.Validator
	ctxer           contextx.Contextualizer
	hh              *healthx.Handler
	migrationStatus *popx.MigrationStatuses
	kc              *aead.AESGCM
	flowc           *aead.XChaCha20Poly1305
	cos             consent.Strategy
	writer          herodot.Writer
	hsm             hsm.Context
	forv            *openid.OpenIDConnectRequestValidator
	fop             fosite.OAuth2Provider
	coh             *consent.Handler
	oah             *oauth2.Handler
	sia             map[string]consent.SubjectIdentifierAlgorithm
	trc             *otelx.Tracer
	tracerWrapper   func(*otelx.Tracer) *otelx.Tracer
	pmm             *prometheus.MetricsManager
	oa2mw           func(h http.Handler) http.Handler
	arhs            []oauth2.AccessRequestHook
	buildVersion    string
	buildHash       string
	buildDate       string
	r               Registry
	persister       persistence.Persister
	jfs             fosite.JWKSFetcherStrategy
	oc              fosite.Configurator
	oidcs           jwk.JWTSigner
	ats             jwk.JWTSigner
	hmacs           *foauth2.HMACSHAStrategy
	devHmac         rfc8628.RFC8628CodeStrategy
	fc              *fositex.Config
	publicCORS      *cors.Cors
	kratos          kratos.Client
	fositeFactories []fositex.Factory
}

func (m *RegistryBase) GetJWKSFetcherStrategy() fosite.JWKSFetcherStrategy {
	if m.jfs == nil {
		m.jfs = fosite.NewDefaultJWKSFetcherStrategy(fosite.JWKSFetcherWithHTTPClientSource(func(ctx context.Context) *retryablehttp.Client {
			return m.HTTPClient(ctx)
		}))
	}
	return m.jfs
}

func (m *RegistryBase) WithContextualizer(ctxer contextx.Contextualizer) Registry {
	m.ctxer = ctxer
	return m.r
}

func (m *RegistryBase) Contextualizer() contextx.Contextualizer {
	if m.ctxer == nil {
		panic("registry Contextualizer not set")
	}
	return m.ctxer
}

func (m *RegistryBase) with(r Registry) *RegistryBase {
	m.r = r
	return m
}

func (m *RegistryBase) WithBuildInfo(version, hash, date string) Registry {
	m.buildVersion = version
	m.buildHash = hash
	m.buildDate = date
	return m.r
}

func (m *RegistryBase) OAuth2AwareMiddleware() func(h http.Handler) http.Handler {
	if m.oa2mw == nil {
		m.oa2mw = oauth2cors.Middleware(m.r)
	}
	return m.oa2mw
}

func (m *RegistryBase) addPublicCORSOnHandler(ctx context.Context) func(http.Handler) http.Handler {
	corsConfig, corsEnabled := m.Config().CORS(ctx, config.PublicInterface)
	if !corsEnabled {
		return func(h http.Handler) http.Handler {
			return h
		}
	}
	if m.publicCORS == nil {
		m.publicCORS = cors.New(corsConfig)
	}
	return func(h http.Handler) http.Handler {
		return m.publicCORS.Handler(h)
	}
}

func (m *RegistryBase) RegisterRoutes(ctx context.Context, admin *httprouterx.RouterAdmin, public *httprouterx.RouterPublic) {
	m.HealthHandler().SetHealthRoutes(admin.Router, true)
	m.HealthHandler().SetVersionRoutes(admin.Router)

	m.HealthHandler().SetHealthRoutes(public.Router, false, healthx.WithMiddleware(m.addPublicCORSOnHandler(ctx)))

	admin.Handler("GET", prometheus.MetricsPrometheusPath, promhttp.Handler())

	m.ConsentHandler().SetRoutes(admin)
	m.KeyHandler().SetRoutes(admin, public, m.OAuth2AwareMiddleware())
	m.ClientHandler().SetRoutes(admin, public)
	m.OAuth2Handler().SetRoutes(admin, public, m.OAuth2AwareMiddleware())
	m.JWTGrantHandler().SetRoutes(admin)
}

func (m *RegistryBase) BuildVersion() string {
	return m.buildVersion
}

func (m *RegistryBase) BuildDate() string {
	return m.buildDate
}

func (m *RegistryBase) BuildHash() string {
	return m.buildHash
}

func (m *RegistryBase) WithConfig(c *config.DefaultProvider) Registry {
	m.conf = c
	return m.r
}

func (m *RegistryBase) Writer() herodot.Writer {
	if m.writer == nil {
		h := herodot.NewJSONWriter(m.Logger())
		h.ErrorEnhancer = x.ErrorEnhancer
		m.writer = h
	}
	return m.writer
}

func (m *RegistryBase) WithLogger(l *logrusx.Logger) Registry {
	m.l = l
	return m.r
}

func (m *RegistryBase) WithTracer(t trace.Tracer) Registry {
	m.trc = new(otelx.Tracer).WithOTLP(t)
	return m.r
}

func (m *RegistryBase) WithTracerWrapper(wrapper TracerWrapper) Registry {
	m.tracerWrapper = wrapper
	return m.r
}

func (m *RegistryBase) WithKratos(k kratos.Client) Registry {
	m.kratos = k
	return m.r
}

func (m *RegistryBase) Logger() *logrusx.Logger {
	if m.l == nil {
		m.l = logrusx.New("Ory Hydra", m.BuildVersion())
	}
	return m.l
}

func (m *RegistryBase) AuditLogger() *logrusx.Logger {
	if m.al == nil {
		m.al = logrusx.NewAudit("Ory Hydra", m.BuildVersion())
		m.al.UseConfig(m.Config().Source(contextx.RootContext))
	}
	return m.al
}

func (m *RegistryBase) ClientHasher() fosite.Hasher {
	if m.fh == nil {
		m.fh = x.NewHasher(m.Config())
	}
	return m.fh
}

func (m *RegistryBase) ClientHandler() *client.Handler {
	if m.ch == nil {
		m.ch = client.NewHandler(m.r)
	}
	return m.ch
}

func (m *RegistryBase) ClientValidator() *client.Validator {
	if m.cv == nil {
		m.cv = client.NewValidator(m.r)
	}
	return m.cv
}

func (m *RegistryBase) KeyHandler() *jwk.Handler {
	if m.kh == nil {
		m.kh = jwk.NewHandler(m.r)
	}
	return m.kh
}

func (m *RegistryBase) JWTGrantHandler() *trust.Handler {
	if m.jwtGrantH == nil {
		m.jwtGrantH = trust.NewHandler(m.r)
	}
	return m.jwtGrantH
}

func (m *RegistryBase) GrantValidator() *trust.GrantValidator {
	if m.jwtGrantV == nil {
		m.jwtGrantV = trust.NewGrantValidator()
	}
	return m.jwtGrantV
}

func (m *RegistryBase) HealthHandler() *healthx.Handler {
	if m.hh == nil {
		m.hh = healthx.NewHandler(m.Writer(), m.buildVersion, healthx.ReadyCheckers{
			"database": func(_ *http.Request) error {
				return m.r.Ping()
			},
			"migrations": func(r *http.Request) error {
				if m.migrationStatus != nil && !m.migrationStatus.HasPending() {
					return nil
				}

				status, err := m.r.Persister().MigrationStatus(r.Context())
				if err != nil {
					return err
				}

				if status.HasPending() {
					err := errors.Errorf("migrations have not yet been fully applied: %+v", status)
					m.Logger().WithField("status", fmt.Sprintf("%+v", status)).WithError(err).Warn("Instance is not yet ready because migrations have not yet been fully applied.")
					return err
				}

				m.migrationStatus = &status
				return nil
			},
		})
	}

	return m.hh
}

func (m *RegistryBase) ConsentStrategy() consent.Strategy {
	if m.cos == nil {
		m.cos = consent.NewStrategy(m.r, m.Config())
	}
	return m.cos
}

func (m *RegistryBase) KeyCipher() *aead.AESGCM {
	if m.kc == nil {
		m.kc = aead.NewAESGCM(m.Config())
	}
	return m.kc
}

func (m *RegistryBase) FlowCipher() *aead.XChaCha20Poly1305 {
	if m.flowc == nil {
		m.flowc = aead.NewXChaCha20Poly1305(m.Config())
	}
	return m.flowc
}

func (m *RegistryBase) CookieStore(ctx context.Context) (sessions.Store, error) {
	var keys [][]byte
	secrets, err := m.conf.GetCookieSecrets(ctx)
	if err != nil {
		return nil, err
	}

	for _, k := range secrets {
		encrypt := sha256.Sum256(k)
		keys = append(keys, k, encrypt[:])
	}

	cs := sessions.NewCookieStore(keys...)
	cs.Options.Secure = m.Config().CookieSecure(ctx)
	cs.Options.HttpOnly = true

	// CookieStore MaxAge is set to 86400 * 30 by default. This prevents secure cookies retrieval with expiration > 30 days.
	// MaxAge(0) disables internal MaxAge check by SecureCookie, see:
	//
	// https://github.com/ory/hydra/pull/2488#discussion_r618992698
	cs.MaxAge(0)

	if domain := m.Config().CookieDomain(ctx); domain != "" {
		cs.Options.Domain = domain
	}

	cs.Options.Path = "/"
	if sameSite := m.Config().CookieSameSiteMode(ctx); sameSite != 0 {
		cs.Options.SameSite = sameSite
	}

	return cs, nil
}

func (m *RegistryBase) HTTPClient(ctx context.Context, opts ...httpx.ResilientOptions) *retryablehttp.Client {
	opts = append(opts,
		httpx.ResilientClientWithLogger(m.Logger()),
		httpx.ResilientClientWithMaxRetry(2),
		httpx.ResilientClientWithConnectionTimeout(30*time.Second))

	tracer := m.Tracer(ctx)
	if tracer.IsLoaded() {
		opts = append(opts, httpx.ResilientClientWithTracer(tracer.Tracer()))
	}

	if m.Config().ClientHTTPNoPrivateIPRanges() {
		opts = append(
			opts,
			httpx.ResilientClientDisallowInternalIPs(),
			httpx.ResilientClientAllowInternalIPRequestsTo(m.Config().ClientHTTPPrivateIPExceptionURLs()...),
		)
	}
	return httpx.NewResilientClient(opts...)
}

func (m *RegistryBase) OAuth2Provider() fosite.OAuth2Provider {
	if m.fop != nil {
		return m.fop
	}

	m.fop = fosite.NewOAuth2Provider(m.r.OAuth2Storage(), m.OAuth2ProviderConfig())
	return m.fop
}

func (m *RegistryBase) OpenIDJWTStrategy() jwk.JWTSigner {
	if m.oidcs != nil {
		return m.oidcs
	}

	m.oidcs = jwk.NewDefaultJWTSigner(m.Config(), m.r, x.OpenIDConnectKeyName)
	return m.oidcs
}

func (m *RegistryBase) AccessTokenJWTStrategy() jwk.JWTSigner {
	if m.ats != nil {
		return m.ats
	}

	m.ats = jwk.NewDefaultJWTSigner(m.Config(), m.r, x.OAuth2JWTKeyName)
	return m.ats
}

func (m *RegistryBase) OAuth2HMACStrategy() *foauth2.HMACSHAStrategy {
	if m.hmacs != nil {
		return m.hmacs
	}

	m.hmacs = compose.NewOAuth2HMACStrategy(m.OAuth2Config())
	return m.hmacs
}

func (m *RegistryBase) RFC8628HMACStrategy() rfc8628.RFC8628CodeStrategy {
	if m.devHmac != nil {
		return m.devHmac
	}

	m.devHmac = compose.NewDeviceStrategy(m.OAuth2Config())
	return m.devHmac
}

func (m *RegistryBase) OAuth2Config() *fositex.Config {
	if m.fc != nil {
		return m.fc
	}

	m.fc = fositex.NewConfig(m.r)
	return m.fc
}

func (m *RegistryBase) ExtraFositeFactories() []fositex.Factory {
	return m.fositeFactories
}

func (m *RegistryBase) WithExtraFositeFactories(f []fositex.Factory) Registry {
	m.fositeFactories = f

	return m.r
}

func (m *RegistryBase) OAuth2ProviderConfig() fosite.Configurator {
	if m.oc != nil {
		return m.oc
	}

	conf := m.OAuth2Config()
	hmacAtStrategy := m.OAuth2HMACStrategy()
	devHmacAtStrategy := m.RFC8628HMACStrategy()
	oidcSigner := m.OpenIDJWTStrategy()
	atSigner := m.AccessTokenJWTStrategy()
	jwtAtStrategy := &foauth2.DefaultJWTStrategy{
		Signer:          atSigner,
		HMACSHAStrategy: hmacAtStrategy,
		Config:          conf,
	}

	conf.LoadDefaultHandlers(&compose.CommonStrategy{
		CoreStrategy: fositex.NewTokenStrategy(m.Config(), hmacAtStrategy, &foauth2.DefaultJWTStrategy{
			Signer:          jwtAtStrategy,
			HMACSHAStrategy: hmacAtStrategy,
			Config:          conf,
		}),
		RFC8628CodeStrategy: devHmacAtStrategy,
		OpenIDConnectTokenStrategy: &openid.DefaultStrategy{
			Config: conf,
			Signer: oidcSigner,
		},
		Signer: oidcSigner,
	})

	m.oc = conf
	return m.oc
}

func (m *RegistryBase) OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator {
	if m.forv == nil {
		m.forv = openid.NewOpenIDConnectRequestValidator(&openid.DefaultStrategy{
			Config: m.OAuth2ProviderConfig(),
			Signer: m.OpenIDJWTStrategy(),
		}, m.OAuth2ProviderConfig())
	}
	return m.forv
}

func (m *RegistryBase) AudienceStrategy() fosite.AudienceMatchingStrategy {
	return fosite.DefaultAudienceMatchingStrategy
}

func (m *RegistryBase) ConsentHandler() *consent.Handler {
	if m.coh == nil {
		m.coh = consent.NewHandler(m.r, m.Config())
	}
	return m.coh
}

func (m *RegistryBase) OAuth2Handler() *oauth2.Handler {
	if m.oah == nil {
		m.oah = oauth2.NewHandler(m.r, m.Config())
	}
	return m.oah
}

func (m *RegistryBase) SubjectIdentifierAlgorithm(ctx context.Context) map[string]consent.SubjectIdentifierAlgorithm {
	if m.sia == nil {
		m.sia = map[string]consent.SubjectIdentifierAlgorithm{}
		for _, t := range m.Config().SubjectTypesSupported(ctx) {
			switch t {
			case "public":
				m.sia["public"] = consent.NewSubjectIdentifierAlgorithmPublic()
			case "pairwise":
				m.sia["pairwise"] = consent.NewSubjectIdentifierAlgorithmPairwise([]byte(m.Config().SubjectIdentifierAlgorithmSalt(ctx)))
			}
		}
	}
	return m.sia
}

func (m *RegistryBase) Tracer(_ context.Context) *otelx.Tracer {
	if m.trc == nil {
		t, err := otelx.New("Ory Hydra", m.l, m.conf.Tracing())
		if err != nil {
			m.Logger().WithError(err).Error("Unable to initialize Tracer.")
		} else {
			// Wrap the tracer if required
			if m.tracerWrapper != nil {
				t = m.tracerWrapper(t)
			}

			m.trc = t
		}
	}
	if m.trc.Tracer() == nil {
		m.trc = otelx.NewNoop(m.l, m.Config().Tracing())
	}

	return m.trc
}

func (m *RegistryBase) PrometheusManager() *prometheus.MetricsManager {
	if m.pmm == nil {
		m.pmm = prometheus.NewMetricsManagerWithPrefix("hydra", prometheus.HTTPMetrics, m.buildVersion, m.buildHash, m.buildDate)
	}
	return m.pmm
}

func (m *RegistryBase) Persister() persistence.Persister {
	return m.persister
}

// Config returns the configuration for the given context. It may or may not be the same as the global configuration.
func (m *RegistryBase) Config() *config.DefaultProvider {
	return m.conf
}

// WithOAuth2Provider forces an oauth2 provider which is only used for testing.
func (m *RegistryBase) WithOAuth2Provider(f fosite.OAuth2Provider) {
	m.fop = f
}

// WithConsentStrategy forces a consent strategy which is only used for testing.
func (m *RegistryBase) WithConsentStrategy(c consent.Strategy) {
	m.cos = c
}

func (m *RegistryBase) AccessRequestHooks() []oauth2.AccessRequestHook {
	if m.arhs == nil {
		m.arhs = []oauth2.AccessRequestHook{
			oauth2.RefreshTokenHook(m),
			oauth2.TokenHook(m),
		}
	}
	return m.arhs
}

func (m *RegistryBase) WithHsmContext(h hsm.Context) {
	m.hsm = h
}

func (m *RegistryBase) HSMContext() hsm.Context {
	if m.hsm == nil {
		m.hsm = hsm.NewContext(m.Config(), m.l)
	}
	return m.hsm
}

func (m *RegistrySQL) ClientAuthenticator() x.ClientAuthenticator {
	return m.OAuth2Provider().(*fosite.Fosite)
}

func (m *RegistryBase) Kratos() kratos.Client {
	if m.kratos == nil {
		m.kratos = kratos.New(m)
	}
	return m.kratos
}
