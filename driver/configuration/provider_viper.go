package configuration

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/cors"

	"github.com/ory/x/logrusx"

	"github.com/ory/viper"

	"github.com/ory/x/corsx"
	"github.com/ory/x/stringsx"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/tracing"
	"github.com/ory/x/urlx"
	"github.com/ory/x/viperx"

	"github.com/ory/hydra/x"
)

type ViperProvider struct {
	l                 *logrusx.Logger
	ss                [][]byte
	generatedSecret   []byte
	forcedHTTP        bool
	insecureRedirects []string
}

const (
	ViperKeyWellKnownKeys                 = "webfinger.jwks.broadcast_keys"
	ViperKeyOAuth2ClientRegistrationURL   = "webfinger.oidc_discovery.client_registration_url"
	ViperKLeyOAuth2TokenURL               = "webfinger.oidc_discovery.token_url" // #nosec G101
	ViperKLeyOAuth2AuthURL                = "webfinger.oidc_discovery.auth_url"
	ViperKeyJWKSURL                       = "webfinger.oidc_discovery.jwks_url"
	ViperKeyOIDCDiscoverySupportedClaims  = "webfinger.oidc_discovery.supported_claims"
	ViperKeyOIDCDiscoverySupportedScope   = "webfinger.oidc_discovery.supported_scope"
	ViperKeyOIDCDiscoveryUserinfoEndpoint = "webfinger.oidc_discovery.userinfo_url"

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

const DefaultSQLiteMemoryDSN = "sqlite://:memory:?_fk=true"

const DSNMemory = "memory"

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func NewViperProvider(l *logrusx.Logger, forcedHTTP bool, insecureRedirects []string) Provider {
	if insecureRedirects == nil {
		insecureRedirects = []string{}
	}
	return &ViperProvider{
		l:                 l,
		forcedHTTP:        forcedHTTP,
		insecureRedirects: insecureRedirects,
	}
}

func (v *ViperProvider) getAddress(address string, port int) string {
	if strings.HasPrefix(address, "unix:") {
		return address
	}
	return fmt.Sprintf("%s:%d", address, port)
}

func (v *ViperProvider) InsecureRedirects() []string {
	return v.insecureRedirects
}

func (v *ViperProvider) WellKnownKeys(include ...string) []string {
	if v.AccessTokenStrategy() == "jwt" {
		include = append(include, x.OAuth2JWTKeyName)
	}

	include = append(include, x.OpenIDConnectKeyName)
	return stringslice.Unique(append(viperx.GetStringSlice(v.l, ViperKeyWellKnownKeys, []string{}), include...))
}

func (v *ViperProvider) ServesHTTPS() bool {
	return !v.forcedHTTP
}

func (v *ViperProvider) IsUsingJWTAsAccessTokens() bool {
	return v.AccessTokenStrategy() != "opaque"
}

func (v *ViperProvider) SubjectTypesSupported() []string {
	types := stringslice.Filter(
		viperx.GetStringSlice(v.l,
			ViperKeySubjectTypesSupported,
			[]string{"public"},
			"oidc.subject_identifiers.enabled",
			"OIDC_SUBJECT_TYPES_SUPPORTED",
		),
		func(s string) bool {
			return !(s == "public" || s == "pairwise")
		},
	)

	if len(types) == 0 {
		types = []string{"public"}
	}

	if stringslice.Has(types, "pairwise") {
		if v.AccessTokenStrategy() == "jwt" {
			v.l.Warn(`The pairwise subject identifier algorithm is not supported by the JWT OAuth 2.0 Access Token Strategy and is thus being disabled. Please remove "pairwise" from oidc.subject_identifiers.enable" (e.g. oidc.subject_identifiers.enable=public) or set strategies.access_token to "opaque".`)
			types = stringslice.Filter(types,
				func(s string) bool {
					return !(s == "public")
				},
			)
		} else if len(v.SubjectIdentifierAlgorithmSalt()) < 8 {
			v.l.Fatalf(`The pairwise subject identifier algorithm was set but length of oidc.subject_identifier.salt is too small (%d < 8), please set oidc.subject_identifiers.pairwise.salt to a random string with 8 characters or more.`, len(v.SubjectIdentifierAlgorithmSalt()))
		}
	}

	return types
}

func (v *ViperProvider) DefaultClientScope() []string {
	return viperx.GetStringSlice(v.l,
		ViperKeyDefaultClientScope,
		[]string{"offline_access", "offline", "openid"},
		"OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE",
	)
}

func (v *ViperProvider) CORSEnabled(iface string) bool {
	return corsx.IsEnabled(v.l, "serve."+iface)
}

func (v *ViperProvider) CORSOptions(iface string) cors.Options {
	return cors.Options{
		AllowedOrigins:     viperx.GetStringSlice(v.l, "serve."+iface+".cors.allowed_origins", []string{}, "CORS_ALLOWED_ORIGINS"),
		AllowedMethods:     viperx.GetStringSlice(v.l, "serve."+iface+".cors.allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, "CORS_ALLOWED_METHODS"),
		AllowedHeaders:     viperx.GetStringSlice(v.l, "serve."+iface+".cors.allowed_headers", []string{"Authorization", "Content-Type"}, "CORS_ALLOWED_HEADERS"),
		ExposedHeaders:     viperx.GetStringSlice(v.l, "serve."+iface+".cors.exposed_headers", []string{"Content-Type"}, "CORS_EXPOSED_HEADERS"),
		AllowCredentials:   viperx.GetBool(v.l, "serve."+iface+".cors.allow_credentials", true, "CORS_ALLOWED_CREDENTIALS"),
		OptionsPassthrough: viperx.GetBool(v.l, "serve."+iface+".cors.options_passthrough", false),
		MaxAge:             viperx.GetInt(v.l, "serve."+iface+".cors.max_age", 0, "CORS_MAX_AGE"),
		Debug:              viperx.GetBool(v.l, "serve."+iface+".cors.debug", false, "CORS_DEBUG"),
	}
}

func (v *ViperProvider) DSN() string {
	dsn := viperx.GetString(v.l, ViperKeyDSN, "", "DATABASE_URL")

	if dsn == DSNMemory {
		return DefaultSQLiteMemoryDSN
	}

	if len(dsn) > 0 {
		return dsn
	}

	v.l.Fatal("dsn must be set")
	return ""
}

func (v *ViperProvider) EncryptSessionData() bool {
	return viperx.GetBool(v.l, ViperKeyEncryptSessionData, true)
}

func (v *ViperProvider) DataSourcePlugin() string {
	return viperx.GetString(v.l, ViperKeyDSN, "", "DATABASE_URL")
}

func (v *ViperProvider) BCryptCost() int {
	return viperx.GetInt(v.l, ViperKeyBCryptCost, 10, "BCRYPT_COST")
}

func (v *ViperProvider) AdminListenOn() string {
	host := viperx.GetString(v.l, ViperKeyAdminListenOnHost, "", "ADMIN_HOST")
	port := viperx.GetInt(v.l, ViperKeyAdminListenOnPort, 4445, "ADMIN_PORT")
	return v.getAddress(host, port)
}

func (v *ViperProvider) AdminDisableHealthAccessLog() bool {
	return viperx.GetBool(v.l, ViperKeyAdminDisableHealthAccessLog, false)
}

func (v *ViperProvider) PublicListenOn() string {
	return v.getAddress(v.publicHost(), v.publicPort())
}

func (v *ViperProvider) PublicDisableHealthAccessLog() bool {
	return viperx.GetBool(v.l, ViperKeyPublicDisableHealthAccessLog, false)
}

func (v *ViperProvider) publicHost() string {
	return viperx.GetString(v.l, ViperKeyPublicListenOnHost, "", "PUBLIC_HOST")
}

func (v *ViperProvider) publicPort() int {
	return viperx.GetInt(v.l, ViperKeyPublicListenOnPort, 4444, "PUBLIC_PORT")
}

func (v *ViperProvider) PublicSocketPermission() *UnixPermission {
	return &UnixPermission{
		Owner: viperx.GetString(v.l, ViperKeyPublicSocketOwner, ""),
		Group: viperx.GetString(v.l, ViperKeyPublicSocketGroup, ""),
		Mode:  os.FileMode(viperx.GetInt(v.l, ViperKeyPublicSocketMode, 0755)),
	}
}

func (v *ViperProvider) adminHost() string {
	return viperx.GetString(v.l, ViperKeyAdminListenOnHost, "", "ADMIN_HOST")
}

func (v *ViperProvider) adminPort() int {
	return viperx.GetInt(v.l, ViperKeyAdminListenOnPort, 4445, "ADMIN_PORT")
}

func (v *ViperProvider) AdminSocketPermission() *UnixPermission {
	return &UnixPermission{
		Owner: viperx.GetString(v.l, ViperKeyAdminSocketOwner, ""),
		Group: viperx.GetString(v.l, ViperKeyAdminSocketGroup, ""),
		Mode:  os.FileMode(viperx.GetInt(v.l, ViperKeyAdminSocketMode, 0755)),
	}
}

func (v *ViperProvider) CookieSameSiteMode() http.SameSite {
	sameSiteModeStr := viperx.GetString(v.l, ViperKeyCookieSameSiteMode, "", "COOKIE_SAME_SITE_MODE")
	switch strings.ToLower(sameSiteModeStr) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		if v.forcedHTTP {
			return http.SameSiteLaxMode
		}
		return http.SameSiteNoneMode
	default:
		if v.forcedHTTP {
			return http.SameSiteLaxMode
		}
		return http.SameSiteDefaultMode
	}
}

func (v *ViperProvider) CookieSameSiteLegacyWorkaround() bool {
	return viperx.GetBool(v.l, ViperKeyCookieSameSiteLegacyWorkaround, false, "COOKIE_SAME_SITE_LEGACY_WORKAROUND")
}

func (v *ViperProvider) ConsentRequestMaxAge() time.Duration {
	return viperx.GetDuration(v.l, ViperKeyConsentRequestMaxAge, time.Minute*30, "LOGIN_CONSENT_REQUEST_LIFESPAN")
}

func (v *ViperProvider) AccessTokenLifespan() time.Duration {
	return viperx.GetDuration(v.l, ViperKeyAccessTokenLifespan, time.Hour, "ACCESS_TOKEN_LIFESPAN")
}

func (v *ViperProvider) RefreshTokenLifespan() time.Duration {
	return viperx.GetDuration(v.l, ViperKeyRefreshTokenLifespan, time.Hour*720, "REFRESH_TOKEN_LIFESPAN")
}

func (v *ViperProvider) IDTokenLifespan() time.Duration {
	return viperx.GetDuration(v.l, ViperKeyIDTokenLifespan, time.Hour, "ID_TOKEN_LIFESPAN")
}

func (v *ViperProvider) AuthCodeLifespan() time.Duration {
	return viperx.GetDuration(v.l, ViperKeyAuthCodeLifespan, time.Minute*10, "AUTH_CODE_LIFESPAN")
}

func (v *ViperProvider) ScopeStrategy() string {
	return viperx.GetString(v.l, ViperKeyScopeStrategy, "", "SCOPE_STRATEGY")
}

func (v *ViperProvider) TracingServiceName() string {
	return viperx.GetString(v.l, "tracing.service_name", "ORY Hydra")
}

func (v *ViperProvider) TracingProvider() string {
	return viperx.GetString(v.l, "tracing.provider", "", "TRACING_PROVIDER")
}

func (v *ViperProvider) TracingJaegerConfig() *tracing.JaegerConfig {
	return &tracing.JaegerConfig{
		LocalAgentHostPort: viperx.GetString(v.l, "tracing.providers.jaeger.local_agent_address", "", "TRACING_PROVIDER_JAEGER_LOCAL_AGENT_ADDRESS"),
		SamplerType:        viperx.GetString(v.l, "tracing.providers.jaeger.sampling.type", "const", "TRACING_PROVIDER_JAEGER_SAMPLING_TYPE"),
		SamplerValue:       viperx.GetFloat64(v.l, "tracing.providers.jaeger.sampling.value", float64(1), "TRACING_PROVIDER_JAEGER_SAMPLING_VALUE"),
		SamplerServerURL:   viperx.GetString(v.l, "tracing.providers.jaeger.sampling.server_url", "", "TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL"),
		Propagation: stringsx.Coalesce(
			viper.GetString("JAEGER_PROPAGATION"), // Standard Jaeger client config
			viperx.GetString(v.l, "tracing.providers.jaeger.propagation", "", "TRACING_PROVIDER_JAEGER_PROPAGATION"),
		),
	}
}

func (v *ViperProvider) TracingZipkinConfig() *tracing.ZipkinConfig {
	return &tracing.ZipkinConfig{
		ServerURL: viperx.GetString(v.l, "tracing.providers.zipkin.server_url", ""),
	}
}

func (v *ViperProvider) GetCookieSecrets() [][]byte {
	return [][]byte{
		[]byte(viperx.GetString(v.l, ViperKeyGetCookieSecrets, string(v.GetSystemSecret()), "COOKIE_SECRET")),
	}
}

func (v *ViperProvider) GetRotatedSystemSecrets() [][]byte {
	secrets := viperx.GetStringSlice(v.l, ViperKeyGetSystemSecret, []string{})

	if len(secrets) < 2 {
		return nil
	}

	var rotated [][]byte
	for _, secret := range secrets[1:] {
		rotated = append(rotated, x.HashStringSecret(secret))
	}

	return rotated
}

func (v *ViperProvider) GetSystemSecret() []byte {
	secrets := viperx.GetStringSlice(v.l, ViperKeyGetSystemSecret, []string{}, "SYSTEM_SECRET")

	if len(secrets) == 0 {
		if v.generatedSecret != nil {
			return v.generatedSecret
		}

		v.l.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secret, err := x.GenerateSecret(32)
		cmdx.Must(err, "Could not generate secret: %s", err)

		v.l.Warnf("Generated secret: %s", secret)
		v.generatedSecret = x.HashByteSecret(secret)

		v.l.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		return x.HashByteSecret(secret)
	}

	secret := secrets[0]
	if len(secret) >= 16 {
		return x.HashStringSecret(secret)
	}

	v.l.Fatalf("System secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
	return nil
}

func (v *ViperProvider) LogoutRedirectURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyLogoutRedirectURL, v.publicFallbackURL("oauth2/fallbacks/logout/callback"), "OAUTH2_LOGOUT_REDIRECT_URL")))
}

func (v *ViperProvider) adminFallbackURL(path string) string {
	return v.fallbackURL(path, v.adminHost(), v.adminPort())

}

func (v *ViperProvider) publicFallbackURL(path string) string {
	if len(v.IssuerURL().String()) > 0 {
		return urlx.AppendPaths(v.IssuerURL(), path).String()
	}

	return v.fallbackURL(path, v.publicHost(), v.publicPort())
}

func (v *ViperProvider) fallbackURL(path string, host string, port int) string {
	var u url.URL
	u.Scheme = "https"
	if !v.ServesHTTPS() {
		u.Scheme = "http"
	}
	if host == "" {
		u.Host = fmt.Sprintf("%s:%d", "localhost", port)
	}
	u.Path = path
	return u.String()
}

func (v *ViperProvider) LoginURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyLoginURL, v.publicFallbackURL("oauth2/fallbacks/login"), "OAUTH2_LOGIN_URL")))
}

func (v *ViperProvider) LogoutURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyLogoutURL, v.publicFallbackURL("oauth2/fallbacks/logout"))))
}

func (v *ViperProvider) ConsentURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyConsentURL, v.publicFallbackURL("oauth2/fallbacks/consent"), "OAUTH2_CONSENT_URL")))
}

func (v *ViperProvider) ErrorURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyErrorURL, v.publicFallbackURL("oauth2/fallbacks/error"), "OAUTH2_ERROR_URL")))
}

func (v *ViperProvider) PublicURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyPublicURL, v.publicFallbackURL("/"))))
}

func (v *ViperProvider) IssuerURL() *url.URL {
	return urlRoot(urlx.ParseOrFatal(v.l, strings.TrimRight(viperx.GetString(v.l, ViperKeyIssuerURL, v.fallbackURL("/", v.publicHost(), v.publicPort()), "OAUTH2_ISSUER_URL", "ISSUER", "ISSUER_URL"), "/")+"/"))
}

func (v *ViperProvider) OAuth2ClientRegistrationURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyOAuth2ClientRegistrationURL, "", "OAUTH2_CLIENT_REGISTRATION_URL"))
}

func (v *ViperProvider) OAuth2TokenURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKLeyOAuth2TokenURL, urlx.AppendPaths(v.IssuerURL(), "/oauth2/token").String()))
}

func (v *ViperProvider) OAuth2AuthURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKLeyOAuth2AuthURL, urlx.AppendPaths(v.IssuerURL(), "/oauth2/auth").String()))
}

func (v *ViperProvider) JWKSURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, ViperKeyJWKSURL, urlx.AppendPaths(v.IssuerURL(), "/.well-known/jwks.json").String()))
}

func (v *ViperProvider) AllowTLSTerminationFrom() []string {
	return viperx.GetStringSlice(v.l, ViperKeyAllowTLSTerminationFrom, []string{}, "HTTPS_ALLOW_TERMINATION_FROM")
}

func (v *ViperProvider) AccessTokenStrategy() string {
	return strings.ToLower(viperx.GetString(v.l, ViperKeyAccessTokenStrategy, "opaque", "OAUTH2_ACCESS_TOKEN_STRATEGY"))
}

func (v *ViperProvider) SubjectIdentifierAlgorithmSalt() string {
	return viperx.GetString(v.l, ViperKeySubjectIdentifierAlgorithmSalt, "", "OIDC_SUBJECT_TYPE_PAIRWISE_SALT")
}

func (v *ViperProvider) OIDCDiscoverySupportedClaims() []string {
	return stringslice.Unique(
		append(
			[]string{"sub"},
			viperx.GetStringSlice(v.l, ViperKeyOIDCDiscoverySupportedClaims, []string{}, "OIDC_DISCOVERY_CLAIMS_SUPPORTED")...,
		),
	)
}

func (v *ViperProvider) OIDCDiscoverySupportedScope() []string {
	return stringslice.Unique(
		append(
			[]string{"offline_access", "offline", "openid"},
			viperx.GetStringSlice(v.l, ViperKeyOIDCDiscoverySupportedScope, []string{}, "OIDC_DISCOVERY_SCOPES_SUPPORTED")...,
		),
	)
}

func (v *ViperProvider) OIDCDiscoveryUserinfoEndpoint() string {
	return viperx.GetString(v.l, ViperKeyOIDCDiscoveryUserinfoEndpoint, urlx.AppendPaths(v.PublicURL(), "/userinfo").String(), "OIDC_DISCOVERY_USERINFO_ENDPOINT")
}

func (v *ViperProvider) ShareOAuth2Debug() bool {
	return viperx.GetBool(v.l, ViperKeyExposeOAuth2Debug, false, "OAUTH2_SHARE_ERROR_DEBUG")
}

func (v *ViperProvider) OAuth2LegacyErrors() bool {
	return viperx.GetBool(v.l, ViperKeyOAuth2LegacyErrors, false)
}

func (v *ViperProvider) PKCEEnforced() bool {
	return viperx.GetBool(v.l, ViperKeyPKCEEnforced, false, "OAUTH2_PKCE_ENFORCED")
}

func (v *ViperProvider) EnforcePKCEForPublicClients() bool {
	return viperx.GetBool(v.l, ViperKeyPKCEEnforcedForPublicClients, false, "OAUTH2_PKCE_ENFORCED_FOR_PUBLIC_CLIENTS")
}

func (v *ViperProvider) CGroupsV1AutoMaxProcsEnabled() bool {
	return viperx.GetBool(v.l, ViperKeyCGroupsV1AutoMaxProcsEnabled, false)
}

func (v *ViperProvider) GrantAllClientCredentialsScopesPerDefault() bool {
	return viperx.GetBool(v.l, ViperKeyGrantAllClientCredentialsScopesPerDefault, false)
}
