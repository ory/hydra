package configuration

import (
	"net/url"
	"time"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/ory/x/tracing"
)

type Provider interface {
	ServesHTTPS() bool

	//HashSignature() bool
	IsUsingJWTAsAccessTokens() bool
	WellKnownKeys(include ...string) []string
	InsecureRedirects() []string

	CORSEnabled(iface string) bool
	CORSOptions(iface string) cors.Options

	SubjectTypesSupported() []string
	ConsentURL() *url.URL
	ErrorURL() *url.URL
	PublicURL() *url.URL
	IssuerURL() *url.URL
	OAuth2AuthURL() string
	OAuth2ClientRegistrationURL() *url.URL
	AllowTLSTerminationFrom() []string
	AccessTokenStrategy() string
	SubjectIdentifierAlgorithmSalt() string
	OIDCDiscoverySupportedScope() []string
	OIDCDiscoverySupportedClaims() []string
	OIDCDiscoveryUserinfoEndpoint() string
	ShareOAuth2Debug() bool
	DSN() string
	BCryptCost() int
	DataSourcePlugin() string
	DefaultClientScope() []string
	AdminListenOn() string
	AdminDisableHealthAccessLog() bool
	PublicListenOn() string
	PublicDisableHealthAccessLog() bool
	ConsentRequestMaxAge() time.Duration
	AccessTokenLifespan() time.Duration
	RefreshTokenLifespan() time.Duration
	IDTokenLifespan() time.Duration
	AuthCodeLifespan() time.Duration
	ScopeStrategy() string
	TracingServiceName() string
	TracingProvider() string
	TracingJaegerConfig() *tracing.JaegerConfig
	GetCookieSecrets() [][]byte
	GetRotatedSystemSecrets() [][]byte
	GetSystemSecret() []byte
	LogoutRedirectURL() *url.URL
	LoginURL() *url.URL
	LogoutURL() *url.URL
	PKCEEnforced() bool
}

func MustValidate(l logrus.FieldLogger, p Provider) {
	if p.ServesHTTPS() {
		if p.IssuerURL().String() == "" {
			l.Fatalf(`Configuration key "%s" must be set unless flag "--dangerous-force-http" is set. To find out more, use "hydra help serve".`, ViperKeyIssuerURL)
		}

		if p.IssuerURL().Scheme != "https" {
			l.Fatalf(`Scheme from configuration key "%s" must be "https" unless --dangerous-force-http is passed but got scheme in value "%s" is "%s". To find out more, use "hydra help serve".`, ViperKeyIssuerURL, p.IssuerURL().String(), p.IssuerURL().Scheme)
		}

		if len(p.InsecureRedirects()) > 0 {
			l.Fatal(`Flag --dangerous-allow-insecure-redirect-urls can only be used in combination with flag --dangerous-force-http`)
		}
	}
}

func urlRoot(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}
