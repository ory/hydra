package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/ory/x/dbal"

	"github.com/markbates/pkger"

	"github.com/ory/x/configx"

	"github.com/rs/cors"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/x"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/tracing"
	"github.com/ory/x/urlx"
)

const (
	ViperKeyWellKnownKeys                             = "webfinger.jwks.broadcast_keys"
	ViperKeyOAuth2ClientRegistrationURL               = "webfinger.oidc_discovery.client_registration_url"
	ViperKLeyOAuth2TokenURL                           = "webfinger.oidc_discovery.token_url" // #nosec G101
	ViperKLeyOAuth2AuthURL                            = "webfinger.oidc_discovery.auth_url"
	ViperKeyJWKSURL                                   = "webfinger.oidc_discovery.jwks_url"
	ViperKeyOIDCDiscoverySupportedClaims              = "webfinger.oidc_discovery.supported_claims"
	ViperKeyOIDCDiscoverySupportedScope               = "webfinger.oidc_discovery.supported_scope"
	ViperKeyOIDCDiscoveryUserinfoEndpoint             = "webfinger.oidc_discovery.userinfo_url"
	ViperKeySubjectTypesSupported                     = "oidc.subject_identifiers.supported_types"
	ViperKeyDefaultClientScope                        = "oidc.dynamic_client_registration.default_scope"
	ViperKeyDSN                                       = "dsn"
	ViperKeyBCryptCost                                = "oauth2.hashers.bcrypt.cost"
	ViperKeyEncryptSessionData                        = "oauth2.session.encrypt_at_rest"
	ViperKeyAdminListenOnHost                         = "serve.admin.host"
	ViperKeyAdminListenOnPort                         = "serve.admin.port"
	ViperKeyAdminSocketOwner                          = "serve.admin.socket.owner"
	ViperKeyAdminSocketGroup                          = "serve.admin.socket.group"
	ViperKeyAdminSocketMode                           = "serve.admin.socket.mode"
	ViperKeyAdminDisableHealthAccessLog               = "serve.admin.access_log.disable_for_health"
	ViperKeyPublicListenOnHost                        = "serve.public.host"
	ViperKeyPublicListenOnPort                        = "serve.public.port"
	ViperKeyPublicSocketOwner                         = "serve.public.socket.owner"
	ViperKeyPublicSocketGroup                         = "serve.public.socket.group"
	ViperKeyPublicSocketMode                          = "serve.public.socket.mode"
	ViperKeyPublicDisableHealthAccessLog              = "serve.public.access_log.disable_for_health"
	ViperKeyCookieSameSiteMode                        = "serve.cookies.same_site_mode"
	ViperKeyCookieSameSiteLegacyWorkaround            = "serve.cookies.same_site_legacy_workaround"
	ViperKeyConsentRequestMaxAge                      = "ttl.login_consent_request"
	ViperKeyAccessTokenLifespan                       = "ttl.access_token"  // #nosec G101
	ViperKeyRefreshTokenLifespan                      = "ttl.refresh_token" // #nosec G101
	ViperKeyIDTokenLifespan                           = "ttl.id_token"      // #nosec G101
	ViperKeyAuthCodeLifespan                          = "ttl.auth_code"
	ViperKeyScopeStrategy                             = "strategies.scope"
	ViperKeyGetCookieSecrets                          = "secrets.cookie"
	ViperKeyGetSystemSecret                           = "secrets.system"
	ViperKeyLogoutRedirectURL                         = "urls.post_logout_redirect"
	ViperKeyLoginURL                                  = "urls.login"
	ViperKeyLogoutURL                                 = "urls.logout"
	ViperKeyConsentURL                                = "urls.consent"
	ViperKeyErrorURL                                  = "urls.error"
	ViperKeyPublicURL                                 = "urls.self.public"
	ViperKeyIssuerURL                                 = "urls.self.issuer"
	ViperKeyAllowTLSTerminationFrom                   = "serve.tls.allow_termination_from"
	ViperKeyAccessTokenStrategy                       = "strategies.access_token"
	ViperKeySubjectIdentifierAlgorithmSalt            = "oidc.subject_identifiers.pairwise.salt"
	ViperKeyPKCEEnforced                              = "oauth2.pkce.enforced"
	ViperKeyPKCEEnforcedForPublicClients              = "oauth2.pkce.enforced_for_public_clients"
	ViperKeyLogLevel                                  = "log.level"
	ViperKeyCGroupsV1AutoMaxProcsEnabled              = "cgroups.v1.auto_max_procs_enabled"
	ViperKeyGrantAllClientCredentialsScopesPerDefault = "oauth2.client_credentials.default_grant_allowed_scope"
	ViperKeyExposeOAuth2Debug                         = "oauth2.expose_internal_errors"
	ViperKeyOAuth2LegacyErrors                        = "oauth2.include_legacy_error_fields"
)

const DSNMemory = "memory"

type ViperProvider struct {
	l               *logrusx.Logger
	ss              [][]byte
	generatedSecret []byte
	p               *configx.Provider
}

func MustNew(flags *pflag.FlagSet, l *logrusx.Logger) *ViperProvider {
	p, err := New(flags, l)
	if err != nil {
		l.WithError(err).Fatalf("Unable to load config.")
	}
	return p
}

func New(flags *pflag.FlagSet, l *logrusx.Logger) (*ViperProvider, error) {
	f, err := pkger.Open("/.schema/config.schema.json")
	if err != nil {
		return nil, err
	}

	schema, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	p, err := configx.New(
		schema,
		flags,
		configx.WithStderrValidationReporter(),
		configx.OmitKeysFromTracing([]string{"dsn","secrets.system","secrets.cookie"}),
		configx.WithImmutables([]string{"log", "serve", "dsn", "profiling"}),
		configx.WithLogrusWatcher(l),
	)
	if err != nil {
		return nil, err
	}

	return &ViperProvider{l: l, p: p}, nil
}

func (p *ViperProvider) Set(key string, value interface{}) {
	p.p.Set(key, value)
}

func (p *ViperProvider) Source() *configx.Provider {
	return p.p
}

func (p *ViperProvider) cors(prefix string) (cors.Options, bool) {
	return p.p.CORS(prefix, cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
}

func (p *ViperProvider) CORS(iface string) (cors.Options, bool) {
	switch iface {
	case "admin":
		return p.AdminCORS()
	case "public":
		return p.PublicCORS()
	default:
		panic(fmt.Sprintf("Received unexpected CORS interface: %s", iface))
	}
}

func (p *ViperProvider) PublicCORS() (cors.Options, bool) {
	return p.cors("serve.public")
}

func (p *ViperProvider) AdminCORS() (cors.Options, bool) {
	return p.cors("serve.admin")
}

func (p *ViperProvider) getAddress(address string, port int) string {
	if strings.HasPrefix(address, "unix:") {
		return address
	}
	return fmt.Sprintf("%s:%d", address, port)
}

func (p *ViperProvider) InsecureRedirects() []string {
	return p.p.Strings("dangerous-allow-insecure-redirect-urls")
}

func (p *ViperProvider) WellKnownKeys(include ...string) []string {
	if p.AccessTokenStrategy() == "jwt" {
		include = append(include, x.OAuth2JWTKeyName)
	}

	include = append(include, x.OpenIDConnectKeyName)
	return stringslice.Unique(append(p.p.Strings(ViperKeyWellKnownKeys), include...))
}

func (p *ViperProvider) ServesHTTPS() bool {
	return !p.forcedHTTP()
}

func (p *ViperProvider) IsUsingJWTAsAccessTokens() bool {
	return p.AccessTokenStrategy() != "opaque"
}

func (p *ViperProvider) SubjectTypesSupported() []string {
	types := stringslice.Filter(
		p.p.StringsF(ViperKeySubjectTypesSupported, []string{"public"}),
		func(s string) bool {
			return !(s == "public" || s == "pairwise")
		},
	)

	if len(types) == 0 {
		types = []string{"public"}
	}

	if stringslice.Has(types, "pairwise") {
		if p.AccessTokenStrategy() == "jwt" {
			p.l.Warn(`The pairwise subject identifier algorithm is not supported by the JWT OAuth 2.0 Access Token Strategy and is thus being disabled. Please remove "pairwise" from oidc.subject_identifiers.supported_types" (e.g. oidc.subject_identifiers.supported_types=public) or set strategies.access_token to "opaque".`)
			types = stringslice.Filter(types,
				func(s string) bool {
					return !(s == "public")
				},
			)
		} else if len(p.SubjectIdentifierAlgorithmSalt()) < 8 {
			p.l.Fatalf(`The pairwise subject identifier algorithm was set but length of oidc.subject_identifier.salt is too small (%d < 8), please set oidc.subject_identifiers.pairwise.salt to a random string with 8 characters or more.`, len(p.SubjectIdentifierAlgorithmSalt()))
		}
	}

	return types
}

func (p *ViperProvider) DefaultClientScope() []string {
	return p.p.StringsF(
		ViperKeyDefaultClientScope,
		[]string{"offline_access", "offline", "openid"},
	)
}

func (p *ViperProvider) DSN() string {
	dsn := p.p.String(ViperKeyDSN)

	if dsn == DSNMemory {
		return dbal.InMemoryDSN
	}

	if len(dsn) > 0 {
		return dsn
	}

	p.l.Fatal("dsn must be set")
	return ""
}

func (p *ViperProvider) EncryptSessionData() bool {
	return p.p.BoolF(ViperKeyEncryptSessionData, true)
}

func (p *ViperProvider) DataSourcePlugin() string {
	return p.p.String(ViperKeyDSN)
}

func (p *ViperProvider) BCryptCost() int {
	return p.p.IntF(ViperKeyBCryptCost, 10)
}

func (p *ViperProvider) AdminListenOn() string {
	host := p.p.String(ViperKeyAdminListenOnHost)
	port := p.p.IntF(ViperKeyAdminListenOnPort, 4445)
	return p.getAddress(host, port)
}

func (p *ViperProvider) AdminDisableHealthAccessLog() bool {
	return p.p.Bool(ViperKeyAdminDisableHealthAccessLog)
}

func (p *ViperProvider) PublicListenOn() string {
	return p.getAddress(p.publicHost(), p.publicPort())
}

func (p *ViperProvider) PublicDisableHealthAccessLog() bool {
	return p.p.Bool(ViperKeyPublicDisableHealthAccessLog)
}

func (p *ViperProvider) publicHost() string {
	return p.p.String(ViperKeyPublicListenOnHost)
}

func (p *ViperProvider) publicPort() int {
	return p.p.IntF(ViperKeyPublicListenOnPort, 4444)
}

func (p *ViperProvider) PublicSocketPermission() *UnixPermission {
	return &UnixPermission{
		Owner: p.p.String(ViperKeyPublicSocketOwner),
		Group: p.p.String(ViperKeyPublicSocketGroup),
		Mode:  os.FileMode(p.p.IntF(ViperKeyPublicSocketMode, 0755)),
	}
}

func (p *ViperProvider) adminHost() string {
	return p.p.String(ViperKeyAdminListenOnHost)
}

func (p *ViperProvider) adminPort() int {
	return p.p.IntF(ViperKeyAdminListenOnPort, 4445)
}

func (p *ViperProvider) AdminSocketPermission() *UnixPermission {
	return &UnixPermission{
		Owner: p.p.String(ViperKeyAdminSocketOwner),
		Group: p.p.String(ViperKeyAdminSocketGroup),
		Mode:  os.FileMode(p.p.IntF(ViperKeyAdminSocketMode, 0755)),
	}
}

func (p *ViperProvider) forcedHTTP() bool {
	return p.p.Bool("dangerous-force-http")
}

func (p *ViperProvider) CookieSameSiteMode() http.SameSite {
	sameSiteModeStr := p.p.String(ViperKeyCookieSameSiteMode)
	switch strings.ToLower(sameSiteModeStr) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		if p.forcedHTTP() {
			return http.SameSiteLaxMode
		}
		return http.SameSiteNoneMode
	default:
		if p.forcedHTTP() {
			return http.SameSiteLaxMode
		}
		return http.SameSiteDefaultMode
	}
}

func (p *ViperProvider) CookieSameSiteLegacyWorkaround() bool {
	return p.p.Bool(ViperKeyCookieSameSiteLegacyWorkaround)
}

func (p *ViperProvider) ConsentRequestMaxAge() time.Duration {
	return p.p.DurationF(ViperKeyConsentRequestMaxAge, time.Minute*30)
}

func (p *ViperProvider) AccessTokenLifespan() time.Duration {
	return p.p.DurationF(ViperKeyAccessTokenLifespan, time.Hour)
}

func (p *ViperProvider) RefreshTokenLifespan() time.Duration {
	return p.p.DurationF(ViperKeyRefreshTokenLifespan, time.Hour*720)
}

func (p *ViperProvider) IDTokenLifespan() time.Duration {
	return p.p.DurationF(ViperKeyIDTokenLifespan, time.Hour)
}

func (p *ViperProvider) AuthCodeLifespan() time.Duration {
	return p.p.DurationF(ViperKeyAuthCodeLifespan, time.Minute*10)
}

func (p *ViperProvider) ScopeStrategy() string {
	return p.p.String(ViperKeyScopeStrategy)
}

func (p *ViperProvider) Tracing() *tracing.Config {
	return p.p.TracingConfig("ORY Hydra")
}

func (p *ViperProvider) GetCookieSecrets() [][]byte {
	secrets := p.p.Strings(ViperKeyGetCookieSecrets)
	if len(secrets) == 0 {
		return [][]byte{p.GetSystemSecret()}
	}

	bs := make([][]byte, len(secrets))
	for k := range secrets {
		bs[k] = []byte(secrets[k])
	}
	return bs
}

func (p *ViperProvider) GetRotatedSystemSecrets() [][]byte {
	secrets := p.p.Strings(ViperKeyGetSystemSecret)

	if len(secrets) < 2 {
		return nil
	}

	var rotated [][]byte
	for _, secret := range secrets[1:] {
		rotated = append(rotated, x.HashStringSecret(secret))
	}

	return rotated
}

func (p *ViperProvider) GetSystemSecret() []byte {
	secrets := p.p.Strings(ViperKeyGetSystemSecret)

	if len(secrets) == 0 {
		if p.generatedSecret != nil {
			return p.generatedSecret
		}

		p.l.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secret, err := x.GenerateSecret(32)
		cmdx.Must(err, "Could not generate secret: %s", err)

		p.l.Warnf("Generated secret: %s", secret)
		p.generatedSecret = x.HashByteSecret(secret)

		p.l.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		return x.HashByteSecret(secret)
	}

	secret := secrets[0]
	if len(secret) >= 16 {
		return x.HashStringSecret(secret)
	}

	p.l.Fatalf("System secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
	return nil
}

func (p *ViperProvider) LogoutRedirectURL() *url.URL {
	return urlRoot(p.p.RequestURIF(ViperKeyLogoutRedirectURL, p.publicFallbackURL("oauth2/fallbacks/logout/callback")))
}

func (p *ViperProvider) adminFallbackURL(path string) *url.URL {
	return p.fallbackURL(path, p.adminHost(), p.adminPort())

}

func (p *ViperProvider) publicFallbackURL(path string) *url.URL {
	if len(p.IssuerURL().String()) > 0 {
		return urlx.AppendPaths(p.IssuerURL(), path)
	}

	return p.fallbackURL(path, p.publicHost(), p.publicPort())
}

func (p *ViperProvider) fallbackURL(path string, host string, port int) *url.URL {
	var u url.URL
	u.Scheme = "https"
	if !p.ServesHTTPS() {
		u.Scheme = "http"
	}
	if host == "" {
		u.Host = fmt.Sprintf("%s:%d", "localhost", port)
	}
	u.Path = path
	return &u
}

func (p *ViperProvider) LoginURL() *url.URL {
	return urlRoot(p.p.RequestURIF(ViperKeyLoginURL, p.publicFallbackURL("oauth2/fallbacks/login")))
}

func (p *ViperProvider) LogoutURL() *url.URL {
	return urlRoot(p.p.RequestURIF(ViperKeyLogoutURL, p.publicFallbackURL("oauth2/fallbacks/logout")))
}

func (p *ViperProvider) ConsentURL() *url.URL {
	return urlRoot(p.p.RequestURIF(ViperKeyConsentURL, p.publicFallbackURL("oauth2/fallbacks/consent")))
}

func (p *ViperProvider) ErrorURL() *url.URL {
	return urlRoot(p.p.RequestURIF(ViperKeyErrorURL, p.publicFallbackURL("oauth2/fallbacks/error")))
}

func (p *ViperProvider) PublicURL() *url.URL {
	return urlRoot(p.p.RequestURIF(ViperKeyPublicURL, p.publicFallbackURL("/")))
}

func (p *ViperProvider) IssuerURL() *url.URL {
	issuerURL := p.p.RequestURIF(ViperKeyIssuerURL, p.fallbackURL("/", p.publicHost(), p.publicPort()))
	issuerURL.Path = strings.TrimRight(issuerURL.Path, "/") + "/"
	return urlRoot(issuerURL)
}

func (p *ViperProvider) OAuth2ClientRegistrationURL() *url.URL {
	return p.p.RequestURIF(ViperKeyOAuth2ClientRegistrationURL, new(url.URL))
}

func (p *ViperProvider) OAuth2TokenURL() *url.URL {
	return p.p.RequestURIF(ViperKLeyOAuth2TokenURL, urlx.AppendPaths(p.IssuerURL(), "/oauth2/token"))
}

func (p *ViperProvider) OAuth2AuthURL() *url.URL {
	return p.p.RequestURIF(ViperKLeyOAuth2AuthURL, urlx.AppendPaths(p.IssuerURL(), "/oauth2/auth"))
}

func (p *ViperProvider) JWKSURL() *url.URL {
	return p.p.RequestURIF(ViperKeyJWKSURL, urlx.AppendPaths(p.IssuerURL(), "/.well-known/jwks.json"))
}

func (p *ViperProvider) AllowTLSTerminationFrom() []string {
	return p.p.Strings(ViperKeyAllowTLSTerminationFrom)
}

func (p *ViperProvider) AccessTokenStrategy() string {
	return strings.ToLower(p.p.StringF(ViperKeyAccessTokenStrategy, "opaque"))
}

func (p *ViperProvider) SubjectIdentifierAlgorithmSalt() string {
	return p.p.String(ViperKeySubjectIdentifierAlgorithmSalt)
}

func (p *ViperProvider) OIDCDiscoverySupportedClaims() []string {
	return stringslice.Unique(
		append(
			[]string{"sub"},
			p.p.Strings(ViperKeyOIDCDiscoverySupportedClaims)...,
		),
	)
}

func (p *ViperProvider) OIDCDiscoverySupportedScope() []string {
	return stringslice.Unique(
		append(
			[]string{"offline_access", "offline", "openid"},
			p.p.Strings(ViperKeyOIDCDiscoverySupportedScope)...,
		),
	)
}

func (p *ViperProvider) OIDCDiscoveryUserinfoEndpoint() *url.URL {
	return p.p.RequestURIF(ViperKeyOIDCDiscoveryUserinfoEndpoint, urlx.AppendPaths(p.PublicURL(), "/userinfo"))
}

func (p *ViperProvider) ShareOAuth2Debug() bool {
	return p.p.Bool(ViperKeyExposeOAuth2Debug)
}

func (p *ViperProvider) OAuth2LegacyErrors() bool {
	return p.p.Bool(ViperKeyOAuth2LegacyErrors)
}

func (p *ViperProvider) PKCEEnforced() bool {
	return p.p.Bool(ViperKeyPKCEEnforced)
}

func (p *ViperProvider) EnforcePKCEForPublicClients() bool {
	return p.p.Bool(ViperKeyPKCEEnforcedForPublicClients)
}

func (p *ViperProvider) CGroupsV1AutoMaxProcsEnabled() bool {
	return p.p.Bool(ViperKeyCGroupsV1AutoMaxProcsEnabled)
}

func (p *ViperProvider) GrantAllClientCredentialsScopesPerDefault() bool {
	return p.p.Bool(ViperKeyGrantAllClientCredentialsScopesPerDefault)
}
