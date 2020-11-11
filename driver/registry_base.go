package driver

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ory/hydra/persistence"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ory/hydra/metrics/prometheus"
	"github.com/ory/x/logrusx"

	"github.com/gorilla/sessions"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/herodot"
	"github.com/ory/x/healthx"
	"github.com/ory/x/resilience"
	"github.com/ory/x/tracing"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
)

type RegistryBase struct {
	l            *logrusx.Logger
	al           *logrusx.Logger
	C            configuration.Provider
	ch           *client.Handler
	fh           fosite.Hasher
	kh           *jwk.Handler
	cv           *client.Validator
	hh           *healthx.Handler
	kg           map[string]jwk.KeyGenerator
	kc           *jwk.AEAD
	cs           sessions.Store
	csPrev       [][]byte
	cos          consent.Strategy
	writer       herodot.Writer
	fsc          fosite.ScopeStrategy
	atjs         jwk.JWTStrategy
	idtjs        jwk.JWTStrategy
	fscPrev      string
	fos          *openid.DefaultStrategy
	forv         *openid.OpenIDConnectRequestValidator
	fop          fosite.OAuth2Provider
	coh          *consent.Handler
	oah          *oauth2.Handler
	sia          map[string]consent.SubjectIdentifierAlgorithm
	trc          *tracing.Tracer
	pmm          *prometheus.MetricsManager
	oa2mw        func(h http.Handler) http.Handler
	o2mc         *foauth2.HMACSHAStrategy
	buildVersion string
	buildHash    string
	buildDate    string
	r            Registry
	persister    persistence.Persister
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
		m.oa2mw = OAuth2AwareCORSMiddleware("public", m.r, m.C)
	}
	return m.oa2mw
}

func (m *RegistryBase) RegisterRoutes(admin *x.RouterAdmin, public *x.RouterPublic) {
	m.HealthHandler().SetRoutes(admin.Router, true)

	public.GET(healthx.AliveCheckPath, m.HealthHandler().Alive)
	public.GET(healthx.ReadyCheckPath, m.HealthHandler().Ready(false))

	admin.Handler("GET", prometheus.MetricsPrometheusPath, promhttp.Handler())

	m.ConsentHandler().SetRoutes(admin)
	m.KeyHandler().SetRoutes(admin, public, m.OAuth2AwareMiddleware())
	m.ClientHandler().SetRoutes(admin)
	m.OAuth2Handler().SetRoutes(admin, public, m.OAuth2AwareMiddleware())
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

func (m *RegistryBase) WithConfig(c configuration.Provider) Registry {
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
		m.l = logrusx.New("ORY Hydra", m.BuildVersion())
	}
	return m.l
}

func (m *RegistryBase) AuditLogger() *logrusx.Logger {
	if m.al == nil {
		m.al = logrusx.NewAudit("ORY Hydra", m.BuildVersion())
	}
	return m.al
}

func (m *RegistryBase) ClientHasher() fosite.Hasher {
	if m.fh == nil {
		if m.Tracer().IsLoaded() {
			m.fh = &tracing.TracedBCrypt{WorkFactor: m.C.BCryptCost()}
		} else {
			m.fh = x.NewBCrypt(m.C)
		}
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
		m.cv = client.NewValidator(m.C)
	}
	return m.cv
}

func (m *RegistryBase) KeyHandler() *jwk.Handler {
	if m.kh == nil {
		m.kh = jwk.NewHandler(m.r, m.C)
	}
	return m.kh
}
func (m *RegistryBase) HealthHandler() *healthx.Handler {
	if m.hh == nil {
		m.hh = healthx.NewHandler(m.Writer(), m.buildVersion, healthx.ReadyCheckers{
			"database": m.r.Ping,
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

func (m *RegistryBase) KeyGenerators() map[string]jwk.KeyGenerator {
	if m.kg == nil {
		m.kg = map[string]jwk.KeyGenerator{
			"RS256": &jwk.RS256Generator{},
			"ES512": &jwk.ECDSA512Generator{},
			"HS256": &jwk.HS256Generator{},
			"HS512": &jwk.HS512Generator{},
		}
	}
	return m.kg
}

func (m *RegistryBase) KeyCipher() *jwk.AEAD {
	if m.kc == nil {
		m.kc = jwk.NewAEAD(m.C)
	}
	return m.kc
}

func (m *RegistryBase) CookieStore() sessions.Store {
	if m.cs == nil {
		m.cs = sessions.NewCookieStore(m.C.GetCookieSecrets()...)
		m.csPrev = m.C.GetCookieSecrets()
	}
	return m.cs
}

func (m *RegistryBase) oAuth2Config() *compose.Config {
	return &compose.Config{
		AccessTokenLifespan:            m.C.AccessTokenLifespan(),
		RefreshTokenLifespan:           m.C.RefreshTokenLifespan(),
		AuthorizeCodeLifespan:          m.C.AuthCodeLifespan(),
		IDTokenLifespan:                m.C.IDTokenLifespan(),
		IDTokenIssuer:                  m.C.IssuerURL().String(),
		HashCost:                       m.C.BCryptCost(),
		ScopeStrategy:                  m.ScopeStrategy(),
		SendDebugMessagesToClients:     m.C.ShareOAuth2Debug(),
		UseLegacyErrorFormat:           m.C.OAuth2LegacyErrors(),
		EnforcePKCE:                    m.C.PKCEEnforced(),
		EnforcePKCEForPublicClients:    m.C.EnforcePKCEForPublicClients(),
		EnablePKCEPlainChallengeMethod: false,
		TokenURL:                       urlx.AppendPaths(m.C.PublicURL(), oauth2.TokenPath).String(),
		RedirectSecureChecker:          x.IsRedirectURISecure(m.C),
	}
}

func (m *RegistryBase) OAuth2HMACStrategy() *foauth2.HMACSHAStrategy {
	if m.o2mc == nil {
		m.o2mc = compose.NewOAuth2HMACStrategy(m.oAuth2Config(), m.C.GetSystemSecret(), m.C.GetRotatedSystemSecrets())
	}
	return m.o2mc
}

func (m *RegistryBase) OAuth2Provider() fosite.OAuth2Provider {
	if m.fop == nil {
		fc := m.oAuth2Config()
		oidcStrategy := &openid.DefaultStrategy{
			JWTStrategy: m.OpenIDJWTStrategy(),
			Expiry:      m.C.IDTokenLifespan(),
			Issuer:      m.C.IssuerURL().String(),
		}

		var coreStrategy foauth2.CoreStrategy
		hmacStrategy := m.OAuth2HMACStrategy()

		switch ats := strings.ToLower(m.C.AccessTokenStrategy()); ats {
		case "jwt":
			coreStrategy = &foauth2.DefaultJWTStrategy{
				JWTStrategy:     m.AccessTokenJWTStrategy(),
				HMACSHAStrategy: hmacStrategy,
			}
		case "opaque":
			coreStrategy = hmacStrategy
		default:
			m.Logger().Fatalf(`Environment variable OAUTH2_ACCESS_TOKEN_STRATEGY is set to "%s" but only "opaque" and "jwt" are valid values.`, ats)
		}

		return compose.Compose(
			fc,
			m.r.OAuth2Storage(),
			&compose.CommonStrategy{
				CoreStrategy:               coreStrategy,
				OpenIDConnectTokenStrategy: oidcStrategy,
				JWTStrategy:                m.OpenIDJWTStrategy(),
			},
			m.ClientHasher(),
			compose.OAuth2AuthorizeExplicitFactory,
			compose.OAuth2AuthorizeImplicitFactory,
			compose.OAuth2ClientCredentialsGrantFactory,
			compose.OAuth2RefreshTokenGrantFactory,
			compose.OpenIDConnectExplicitFactory,
			compose.OpenIDConnectHybridFactory,
			compose.OpenIDConnectImplicitFactory,
			compose.OpenIDConnectRefreshFactory,
			compose.OAuth2TokenRevocationFactory,
			compose.OAuth2TokenIntrospectionFactory,
			compose.OAuth2PKCEFactory,
		)
	}
	return m.fop
}

func (m *RegistryBase) ScopeStrategy() fosite.ScopeStrategy {
	if m.fsc == nil {
		if m.C.ScopeStrategy() == "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY" {
			m.Logger().Warn(`Using deprecated hierarchical scope strategy, consider upgrading to "wildcard"" or "exact"".`)
			m.fsc = fosite.HierarchicScopeStrategy
		} else if strings.ToLower(m.C.ScopeStrategy()) == "exact" {
			m.fsc = fosite.ExactScopeStrategy
		} else {
			m.fsc = fosite.WildcardScopeStrategy
		}
		m.fscPrev = m.C.ScopeStrategy()
	}
	return m.fsc
}

func (m *RegistryBase) newKeyStrategy(key string) (s jwk.JWTStrategy) {
	if err := jwk.EnsureAsymmetricKeypairExists(context.Background(), m.r, new(jwk.RS256Generator), key); err != nil {
		m.Logger().WithError(err).Fatalf(`Could not ensure that signing keys for "%s" exists. If you are running against a persistent SQL database this is most likely because your "secrets.system" ("SECRETS_SYSTEM" environment variable) is not set or changed. When running with an SQL database backend you need to make sure that the secret is set and stays the same, unless when doing key rotation. This may also happen when you forget to run "hydra migrate sql"..`, key)
	}

	if err := resilience.Retry(m.Logger(), time.Second*15, time.Minute*15, func() (err error) {
		s, err = jwk.NewRS256JWTStrategy(m.r, func() string {
			return key
		})
		return err
	}); err != nil {
		m.Logger().WithError(err).Fatalf("Unable to initialize JSON Web Token strategy.")
	}

	return s
}

func (m *RegistryBase) AccessTokenJWTStrategy() jwk.JWTStrategy {
	if m.atjs == nil {
		m.atjs = m.newKeyStrategy(x.OAuth2JWTKeyName)
	}
	return m.atjs
}

func (m *RegistryBase) OpenIDJWTStrategy() jwk.JWTStrategy {
	if m.idtjs == nil {
		m.idtjs = m.newKeyStrategy(x.OpenIDConnectKeyName)
	}
	return m.idtjs
}

func (m *RegistryBase) FositeOpenIDDefaultStrategy() *openid.DefaultStrategy {
	if m.fos == nil {
		m.fos = &openid.DefaultStrategy{
			JWTStrategy: m.OpenIDJWTStrategy(),
		}
	}
	return m.fos
}

func (m *RegistryBase) OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator {
	if m.forv == nil {
		m.forv = openid.NewOpenIDConnectRequestValidator([]string{"login", "none", "consent"}, m.FositeOpenIDDefaultStrategy())
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

func (m *RegistryBase) SubjectIdentifierAlgorithm() map[string]consent.SubjectIdentifierAlgorithm {
	if m.sia == nil {
		m.sia = map[string]consent.SubjectIdentifierAlgorithm{}
		for _, t := range m.C.SubjectTypesSupported() {
			switch t {
			case "public":
				m.sia["public"] = consent.NewSubjectIdentifierAlgorithmPublic()
			case "pairwise":
				m.sia["pairwise"] = consent.NewSubjectIdentifierAlgorithmPairwise([]byte(m.C.SubjectIdentifierAlgorithmSalt()))
			}
		}
	}
	return m.sia
}

func (m *RegistryBase) Tracer() *tracing.Tracer {
	if m.trc == nil {
		m.trc = &tracing.Tracer{
			ServiceName:  m.C.TracingServiceName(),
			JaegerConfig: m.C.TracingJaegerConfig(),
			ZipkinConfig: m.C.TracingZipkinConfig(),
			Provider:     m.C.TracingProvider(),
			Logger:       m.Logger(),
		}

		if err := m.trc.Setup(); err != nil {
			m.Logger().WithError(err).Fatalf("Unable to initialize Tracer.")
		}
	}

	return m.trc
}

func (m *RegistryBase) PrometheusManager() *prometheus.MetricsManager {
	if m.pmm == nil {
		m.pmm = prometheus.NewMetricsManager(m.buildVersion, m.buildHash, m.buildDate)
	}
	return m.pmm
}

func (m *RegistryBase) Persister() persistence.Persister {
	return m.persister
}

func (m *RegistryBase) Config() configuration.Provider {
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
