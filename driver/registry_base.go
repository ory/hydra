package driver

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ory/hydra/metrics/prometheus"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/serverx"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/healthx"
	"github.com/ory/x/resilience"
	"github.com/ory/x/tracing"
	"github.com/ory/x/urlx"
)

const (
	MetricsPrometheusPath = "/metrics/prometheus"
)

type RegistryBase struct {
	l            logrus.FieldLogger
	c            configuration.Provider
	cm           client.Manager
	ch           *client.Handler
	fh           fosite.Hasher
	kh           *jwk.Handler
	cv           *client.Validator
	hh           *healthx.Handler
	kg           map[string]jwk.KeyGenerator
	km           jwk.Manager
	kc           *jwk.AEAD
	cs           sessions.Store
	csPrev       [][]byte
	com          consent.Manager
	cos          consent.Strategy
	writer       herodot.Writer
	fs           x.FositeStorer
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
		m.oa2mw = OAuth2AwareCORSMiddleware("public", m.r, m.c)
	}
	return m.oa2mw
}

func (m *RegistryBase) RegisterRoutes(admin *x.RouterAdmin, public *x.RouterPublic) {
	m.HealthHandler().SetRoutes(admin.Router, true)

	public.GET(healthx.AliveCheckPath, m.HealthHandler().Alive)
	public.GET(healthx.ReadyCheckPath, m.HealthHandler().Ready(false))

	admin.Handler("GET", MetricsPrometheusPath, promhttp.Handler())

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
	m.c = c
	return m.r
}

func (m *RegistryBase) Writer() herodot.Writer {
	if m.writer == nil {
		h := herodot.NewJSONWriter(m.Logger())
		h.ErrorEnhancer = serverx.ErrorEnhancerRFC6749
		m.writer = h
	}
	return m.writer
}

func (m *RegistryBase) WithLogger(l logrus.FieldLogger) Registry {
	m.l = l
	return m.r
}

func (m *RegistryBase) Logger() logrus.FieldLogger {
	if m.l == nil {
		m.l = logrusx.New()
	}
	return m.l
}

func (m *RegistryBase) ClientHasher() fosite.Hasher {
	if m.fh == nil {
		if m.Tracer().IsLoaded() {
			m.fh = &tracing.TracedBCrypt{WorkFactor: m.c.BCryptCost()}
		} else {
			m.fh = x.NewBCrypt(m.c)
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
		m.cv = client.NewValidator(m.c)
	}
	return m.cv
}

func (m *RegistryBase) KeyHandler() *jwk.Handler {
	if m.kh == nil {
		m.kh = jwk.NewHandler(m.r, m.c)
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
		m.cos = consent.NewStrategy(m.r, m.c)
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
		m.kc = jwk.NewAEAD(m.c)
	}
	return m.kc
}

func (m *RegistryBase) CookieStore() sessions.Store {
	if m.cs == nil {
		m.cs = sessions.NewCookieStore(m.c.GetCookieSecrets()...)
		m.csPrev = m.c.GetCookieSecrets()
	}
	return m.cs
}

func (m *RegistryBase) oAuth2Config() *compose.Config {
	return &compose.Config{
		AccessTokenLifespan:            m.c.AccessTokenLifespan(),
		RefreshTokenLifespan:           m.c.RefreshTokenLifespan(),
		AuthorizeCodeLifespan:          m.c.AuthCodeLifespan(),
		IDTokenLifespan:                m.c.IDTokenLifespan(),
		IDTokenIssuer:                  m.c.IssuerURL().String(),
		HashCost:                       m.c.BCryptCost(),
		ScopeStrategy:                  m.ScopeStrategy(),
		SendDebugMessagesToClients:     m.c.ShareOAuth2Debug(),
		EnforcePKCE:                    true,
		EnablePKCEPlainChallengeMethod: false,
		TokenURL:                       urlx.AppendPaths(m.c.PublicURL(), oauth2.TokenPath).String(),
		RedirectSecureChecker:          x.IsRedirectURISecure(m.c),
	}
}

func (m *RegistryBase) OAuth2HMACStrategy() *foauth2.HMACSHAStrategy {
	if m.o2mc == nil {
		m.o2mc = compose.NewOAuth2HMACStrategy(m.oAuth2Config(), m.c.GetSystemSecret(), m.c.GetRotatedSystemSecrets())
	}
	return m.o2mc
}

func (m *RegistryBase) OAuth2Provider() fosite.OAuth2Provider {
	if m.fop == nil {
		fc := m.oAuth2Config()
		oidcStrategy := &openid.DefaultStrategy{
			JWTStrategy: m.OpenIDJWTStrategy(),
			Expiry:      m.c.IDTokenLifespan(),
			Issuer:      m.c.IssuerURL().String(),
		}

		var coreStrategy foauth2.CoreStrategy
		hmacStrategy := m.OAuth2HMACStrategy()

		switch ats := strings.ToLower(m.c.AccessTokenStrategy()); ats {
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
		if m.c.ScopeStrategy() == "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY" {
			m.Logger().Warn(`Using deprecated hierarchical scope strategy, consider upgrading to "wildcard"" or "exact"".`)
			m.fsc = fosite.HierarchicScopeStrategy
		} else if strings.ToLower(m.c.ScopeStrategy()) == "exact" {
			m.fsc = fosite.ExactScopeStrategy
		} else {
			m.fsc = fosite.WildcardScopeStrategy
		}
		m.fscPrev = m.c.ScopeStrategy()
	}
	return m.fsc
}

func (m *RegistryBase) newKeyStrategy(key string) (s jwk.JWTStrategy) {
	if err := jwk.EnsureAsymmetricKeypairExists(context.Background(), m.r, new(jwk.RS256Generator), key); err != nil {
		m.Logger().WithError(err).Fatalf(`Could not ensure that signing keys for "%s" exists. This can happen if you forget to run "hydra migrate sql", set the wrong "secrets.system" or forget to set "secrets.system" entirely.`, key)
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
		m.coh = consent.NewHandler(m.r, m.c)
	}
	return m.coh
}

func (m *RegistryBase) OAuth2Handler() *oauth2.Handler {
	if m.oah == nil {
		m.oah = oauth2.NewHandler(m.r, m.c)
	}
	return m.oah
}

func (m *RegistryBase) SubjectIdentifierAlgorithm() map[string]consent.SubjectIdentifierAlgorithm {
	if m.sia == nil {
		m.sia = map[string]consent.SubjectIdentifierAlgorithm{}
		for _, t := range m.c.SubjectTypesSupported() {
			switch t {
			case "public":
				m.sia["public"] = consent.NewSubjectIdentifierAlgorithmPublic()
			case "pairwise":
				m.sia["pairwise"] = consent.NewSubjectIdentifierAlgorithmPairwise([]byte(m.c.SubjectIdentifierAlgorithmSalt()))
			}
		}
	}
	return m.sia
}

func (m *RegistryBase) Tracer() *tracing.Tracer {
	if m.trc == nil {
		m.trc = &tracing.Tracer{
			ServiceName:  m.c.TracingServiceName(),
			JaegerConfig: m.c.TracingJaegerConfig(),
			Provider:     m.c.TracingProvider(),
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
