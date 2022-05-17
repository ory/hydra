package driver

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ory/hydra/fositex"
	ctxx "github.com/ory/x/contextx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/otelx"
	"net"
	"net/http"
	"time"

	"github.com/ory/hydra/hsm"

	prometheus "github.com/ory/x/prometheusx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/hydra/x/oauth2cors"
	"github.com/ory/x/contextx"

	"github.com/ory/hydra/persistence"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ory/x/logrusx"

	"github.com/gorilla/sessions"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/healthx"
	"github.com/ory/x/resilience"
	"github.com/ory/x/tracing"
)

var (
	_ contextx.Provider = (*RegistryBase)(nil)
)

type RegistryBase struct {
	l            *logrusx.Logger
	al           *logrusx.Logger
	C            *config.DefaultProvider
	ch           *client.Handler
	fh           fosite.Hasher
	jwtGrantH    *trust.Handler
	jwtGrantV    *trust.GrantValidator
	kh           *jwk.Handler
	cv           *client.Validator
	ctxer        contextx.Contextualizer
	hh           *healthx.Handler
	kc           *jwk.AEAD
	cs           sessions.Store
	csPrev       [][]byte
	cos          consent.Strategy
	writer       herodot.Writer
	fsc          fosite.ScopeStrategy
	atjs         jwk.JWTSigner
	idtjs        jwk.JWTSigner
	hsm          hsm.Context
	fscPrev      string
	forv         *openid.OpenIDConnectRequestValidator
	fop          fosite.OAuth2Provider
	coh          *consent.Handler
	oah          *oauth2.Handler
	sia          map[string]consent.SubjectIdentifierAlgorithm
	trc          *otelx.Tracer
	pmm          *prometheus.MetricsManager
	oa2mw        func(h http.Handler) http.Handler
	o2mc         *foauth2.HMACSHAStrategy
	o2jwt        *foauth2.DefaultJWTStrategy
	arhs         []oauth2.AccessRequestHook
	buildVersion string
	buildHash    string
	buildDate    string
	r            Registry
	persister    persistence.Persister
	jfs          fosite.JWKSFetcherStrategy
	oc           fosite.Configurator
	oidcs        jwk.JWTSigner
	ats          jwk.JWTSigner
	hmacs        *foauth2.HMACSHAStrategy
	fc           *fositex.Config
}

func (m *RegistryBase) GetJWKSFetcherStrategy(ctx context.Context) fosite.JWKSFetcherStrategy {
	if m.jfs == nil {
		m.jfs = fosite.NewDefaultJWKSFetcherStrategy(fosite.JWKSFetcherWithHTTPClient(m.HTTPClient(ctx)))
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

func (m *RegistryBase) OAuth2AwareMiddleware(ctx context.Context) func(h http.Handler) http.Handler {
	if m.oa2mw == nil {
		m.oa2mw = oauth2cors.Middleware(ctx, m.r)
	}
	return m.oa2mw
}

func (m *RegistryBase) RegisterRoutes(ctx context.Context, admin *x.RouterAdmin, public *x.RouterPublic) {
	m.HealthHandler().SetHealthRoutes(admin.Router, true)
	m.HealthHandler().SetVersionRoutes(admin.Router)

	m.HealthHandler().SetHealthRoutes(public.Router, false)

	admin.Handler("GET", prometheus.MetricsPrometheusPath, promhttp.Handler())

	m.ConsentHandler().SetRoutes(admin)
	m.KeyHandler().SetRoutes(admin, public, m.OAuth2AwareMiddleware(ctx))
	m.ClientHandler().SetRoutes(admin, public)
	m.OAuth2Handler().SetRoutes(admin, public, m.OAuth2AwareMiddleware(ctx))
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
	m.C = c
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

func (m *RegistryBase) Logger() *logrusx.Logger {
	if m.l == nil {
		m.l = logrusx.New("Ory Hydra", m.BuildVersion())
	}
	return m.l
}

func (m *RegistryBase) AuditLogger() *logrusx.Logger {
	if m.al == nil {
		m.al = logrusx.NewAudit("Ory Hydra", m.BuildVersion())
		m.al.UseConfig(m.Config().Source(ctxx.RootContext))
	}
	return m.al
}

func (m *RegistryBase) ClientHasher() fosite.Hasher {
	if m.fh == nil {
		m.fh = &tracing.TracedBCrypt{GetWorkFactor: m.Config().GetBCryptCost}
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
				status, err := m.r.Persister().MigrationStatus(r.Context())
				if err != nil {
					return err
				}

				if status.HasPending() {
					err := errors.Errorf("migrations have not yet been fully applied: %+v", status)
					m.Logger().WithField("status", fmt.Sprintf("%+v", status)).WithError(err).Warn("Instance is not yet ready because migrations have not yet been fully applied.")
					return err
				}

				return nil
			},
		})
	}

	return m.hh
}

func (m *RegistryBase) ConsentStrategy() consent.Strategy {
	if m.cos == nil {
		m.cos = consent.NewStrategy(m.r, m.C)
	}
	return m.cos
}

func (m *RegistryBase) KeyCipher() *jwk.AEAD {
	if m.kc == nil {
		m.kc = jwk.NewAEAD(m.C)
	}
	return m.kc
}

func (m *RegistryBase) CookieStore(ctx context.Context) sessions.Store {
	if m.cs == nil {
		cs := sessions.NewCookieStore(m.C.GetCookieSecrets(ctx)...)
		// CookieStore MaxAge is set to 86400 * 30 by default. This prevents secure cookies retrieval with expiration > 30 days.
		// MaxAge(0) disables internal MaxAge check by SecureCookie, see:
		//
		// https://github.com/ory/hydra/pull/2488#discussion_r618992698
		cs.MaxAge(0)

		m.cs = cs
		m.csPrev = m.C.GetCookieSecrets(ctx)
	}
	return m.cs
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

	if m.Config().ClientHTTPNoPrivateIPRanges(ctx) {
		opts = append(opts, httpx.ResilientClientDisallowInternalIPs())
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

	m.oidcs = m.newKeyStrategy(x.OpenIDConnectKeyName)
	return m.oidcs
}

func (m *RegistryBase) AccessTokenJWTStrategy() jwk.JWTSigner {
	if m.ats != nil {
		return m.ats
	}

	m.ats = m.newKeyStrategy(x.OAuth2JWTKeyName)
	return m.ats
}

func (m *RegistryBase) OAuth2HMACStrategy() *foauth2.HMACSHAStrategy {
	if m.hmacs != nil {
		return m.hmacs
	}

	m.hmacs = compose.NewOAuth2HMACStrategy(m.OAuth2Config())
	return m.hmacs
}

func (m *RegistryBase) OAuth2Config() *fositex.Config {
	if m.fc != nil {
		return m.fc
	}

	m.fc = fositex.NewConfig(m.r)
	return m.fc
}

func (m *RegistryBase) OAuth2ProviderConfig() fosite.Configurator {
	if m.oc != nil {
		return m.oc
	}

	conf := m.OAuth2Config()
	hmacAtStrategy := m.OAuth2HMACStrategy()
	oidcSigner := m.OpenIDJWTStrategy()
	atSigner := m.AccessTokenJWTStrategy()
	jwtAtStrategy := &foauth2.DefaultJWTStrategy{
		Signer:          atSigner,
		HMACSHAStrategy: hmacAtStrategy,
		Config:          conf,
	}

	conf.LoadDefaultHanlders(&compose.CommonStrategy{
		CoreStrategy: fositex.NewTokenStrategy(m.Config(), hmacAtStrategy, &foauth2.DefaultJWTStrategy{
			Signer:          jwtAtStrategy,
			HMACSHAStrategy: hmacAtStrategy,
			Config:          conf,
		}),
		OpenIDConnectTokenStrategy: &openid.DefaultStrategy{
			Config: conf,
			Signer: oidcSigner,
		},
		Signer: oidcSigner,
	})

	m.oc = conf
	return m.oc
}

func (m *RegistryBase) newKeyStrategy(key string) (s jwk.JWTSigner) {
	if err := jwk.EnsureAsymmetricKeypairExists(context.Background(), m.r, "RS256", key); err != nil {
		var netError net.Error
		if errors.As(err, &netError) {
			m.Logger().WithError(err).Fatalf(`Could not ensure that signing keys for "%s" exists. A network error occurred, see error for specific details.`, key)
			return
		}

		m.Logger().WithError(err).Fatalf(`Could not ensure that signing keys for "%s" exists. If you are running against a persistent SQL database this is most likely because your "secrets.system" ("SECRETS_SYSTEM" environment variable) is not set or changed. When running with an SQL database backend you need to make sure that the secret is set and stays the same, unless when doing key rotation. This may also happen when you forget to run "hydra migrate sql"..`, key)
	}

	if err := resilience.Retry(m.Logger(), time.Second*15, time.Minute*15, func() (err error) {
		s, err = jwk.NewDefaultJWTSigner(*m.C, m.r, key)
		return err
	}); err != nil {
		m.Logger().WithError(err).Fatalf("Unable to initialize JSON Web Token strategy.")
	}

	return s
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
		m.coh = consent.NewHandler(m.r, m.C)
	}
	return m.coh
}

func (m *RegistryBase) OAuth2Handler() *oauth2.Handler {
	if m.oah == nil {
		m.oah = oauth2.NewHandler(m.r, m.C)
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

func (m *RegistryBase) Tracer(ctx context.Context) *otelx.Tracer {
	if m.trc == nil {
		t, err := otelx.New("Ory Hydra", m.l, m.C.Tracing(ctx))
		if err != nil {
			m.Logger().WithError(err).Error("Unable to initialize Tracer.")
		} else {
			m.trc = t
		}
	}

	return m.trc
}

func (m *RegistryBase) PrometheusManager() *prometheus.MetricsManager {
	if m.pmm == nil {
		m.pmm = prometheus.NewMetricsManager("hydra", m.buildVersion, m.buildHash, m.buildDate)
	}
	return m.pmm
}

func (m *RegistryBase) Persister() persistence.Persister {
	return m.persister
}

// Config returns the configuration for the given context. It may or may not be the same as the global configuration.
func (m *RegistryBase) Config() *config.DefaultProvider {
	return m.C
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
			oauth2.RefreshTokenHook(m.C),
		}
	}
	return m.arhs
}

func (m *RegistryBase) WithHsmContext(h hsm.Context) {
	m.hsm = h
}

func (m *RegistryBase) HsmContext() hsm.Context {
	if m.hsm == nil {
		m.hsm = hsm.NewContext(m.C, m.l)
	}
	return m.hsm
}

func (m *RegistrySQL) ClientAuthenticator() x.ClientAuthenticator {
	return m.OAuth2Provider().(*fosite.Fosite)
}
