package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ory/x/dbal"

	"github.com/ory/hydra/spec"

	"github.com/ory/x/configx"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/x"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/tracing"
	"github.com/ory/x/urlx"
)

const (
	KeyRoot                                      = ""
	HsmEnabled                                   = "hsm.enabled"
	HsmLibraryPath                               = "hsm.library"
	HsmPin                                       = "hsm.pin"
	HsmSlotNumber                                = "hsm.slot"
	HsmKeySetPrefix                              = "hsm.key_set_prefix"
	HsmTokenLabel                                = "hsm.token_label" // #nosec G101
	KeyWellKnownKeys                             = "webfinger.jwks.broadcast_keys"
	KeyOAuth2ClientRegistrationURL               = "webfinger.oidc_discovery.client_registration_url"
	KeyOAuth2TokenURL                            = "webfinger.oidc_discovery.token_url" // #nosec G101
	KeyOAuth2AuthURL                             = "webfinger.oidc_discovery.auth_url"
	KeyJWKSURL                                   = "webfinger.oidc_discovery.jwks_url"
	KeyOIDCDiscoverySupportedClaims              = "webfinger.oidc_discovery.supported_claims"
	KeyOIDCDiscoverySupportedScope               = "webfinger.oidc_discovery.supported_scope"
	KeyOIDCDiscoveryUserinfoEndpoint             = "webfinger.oidc_discovery.userinfo_url"
	KeySubjectTypesSupported                     = "oidc.subject_identifiers.supported_types"
	KeyDefaultClientScope                        = "oidc.dynamic_client_registration.default_scope"
	KeyDSN                                       = "dsn"
	KeyBCryptCost                                = "oauth2.hashers.bcrypt.cost"
	KeyEncryptSessionData                        = "oauth2.session.encrypt_at_rest"
	KeyCookieSameSiteMode                        = "serve.cookies.same_site_mode"
	KeyCookieSameSiteLegacyWorkaround            = "serve.cookies.same_site_legacy_workaround"
	KeyConsentRequestMaxAge                      = "ttl.login_consent_request"
	KeyAccessTokenLifespan                       = "ttl.access_token"  // #nosec G101
	KeyRefreshTokenLifespan                      = "ttl.refresh_token" // #nosec G101
	KeyIDTokenLifespan                           = "ttl.id_token"      // #nosec G101
	KeyAuthCodeLifespan                          = "ttl.auth_code"
	KeyScopeStrategy                             = "strategies.scope"
	KeyGetCookieSecrets                          = "secrets.cookie"
	KeyGetSystemSecret                           = "secrets.system"
	KeyLogoutRedirectURL                         = "urls.post_logout_redirect"
	KeyLoginURL                                  = "urls.login"
	KeyLogoutURL                                 = "urls.logout"
	KeyConsentURL                                = "urls.consent"
	KeyErrorURL                                  = "urls.error"
	KeyPublicURL                                 = "urls.self.public"
	KeyIssuerURL                                 = "urls.self.issuer"
	KeyAccessTokenStrategy                       = "strategies.access_token"
	KeySubjectIdentifierAlgorithmSalt            = "oidc.subject_identifiers.pairwise.salt"
	KeyPublicAllowDynamicRegistration            = "oidc.dynamic_client_registration.enabled"
	KeyPKCEEnforced                              = "oauth2.pkce.enforced"
	KeyPKCEEnforcedForPublicClients              = "oauth2.pkce.enforced_for_public_clients"
	KeyLogLevel                                  = "log.level"
	KeyCGroupsV1AutoMaxProcsEnabled              = "cgroups.v1.auto_max_procs_enabled"
	KeyGrantAllClientCredentialsScopesPerDefault = "oauth2.client_credentials.default_grant_allowed_scope" // #nosec G101
	KeyExposeOAuth2Debug                         = "oauth2.expose_internal_errors"
	KeyOAuth2LegacyErrors                        = "oauth2.include_legacy_error_fields"
	KeyExcludeNotBeforeClaim                     = "oauth2.exclude_not_before_claim"
	KeyAllowedTopLevelClaims                     = "oauth2.allowed_top_level_claims"
	KeyOAuth2GrantJWTIDOptional                  = "oauth2.grant.jwt.jti_optional"
	KeyOAuth2GrantJWTIssuedDateOptional          = "oauth2.grant.jwt.iat_optional"
	KeyOAuth2GrantJWTMaxDuration                 = "oauth2.grant.jwt.max_ttl"
	KeyRefreshTokenHookURL                       = "oauth2.refresh_token_hook" // #nosec G101
)

const DSNMemory = "memory"

type Provider struct {
	l               *logrusx.Logger
	generatedSecret []byte
	p               *configx.Provider
}

func MustNew(ctx context.Context, l *logrusx.Logger, opts ...configx.OptionModifier) *Provider {
	p, err := New(ctx, l, opts...)
	if err != nil {
		l.WithError(err).Fatalf("Unable to load config.")
	}
	return p
}

func New(ctx context.Context, l *logrusx.Logger, opts ...configx.OptionModifier) (*Provider, error) {
	opts = append([]configx.OptionModifier{
		configx.WithStderrValidationReporter(),
		configx.OmitKeysFromTracing("dsn", "secrets.system", "secrets.cookie"),
		configx.WithImmutables("log", "serve", "dsn", "profiling"),
		configx.WithLogrusWatcher(l),
	}, opts...)

	p, err := configx.New(ctx, spec.ConfigValidationSchema, opts...)
	if err != nil {
		return nil, err
	}

	l.UseConfig(p)
	return &Provider{l: l, p: p}, nil
}

func (p *Provider) Set(key string, value interface{}) error {
	return p.p.Set(key, value)
}

func (p *Provider) MustSet(key string, value interface{}) {
	if err := p.Set(key, value); err != nil {
		p.l.WithError(err).Fatalf("Unable to set \"%s\" to \"%s\".", key, value)
	}
}

func (p *Provider) Source() *configx.Provider {
	return p.p
}

func (p *Provider) InsecureRedirects() []string {
	return p.p.Strings("dangerous-allow-insecure-redirect-urls")
}

func (p *Provider) WellKnownKeys(include ...string) []string {
	if p.AccessTokenStrategy() == "jwt" {
		include = append(include, x.OAuth2JWTKeyName)
	}

	include = append(include, x.OpenIDConnectKeyName)
	return stringslice.Unique(append(p.p.Strings(KeyWellKnownKeys), include...))
}

func (p *Provider) IsUsingJWTAsAccessTokens() bool {
	return p.AccessTokenStrategy() != "opaque"
}

func (p *Provider) AllowedTopLevelClaims() []string {
	return stringslice.Unique(p.p.Strings(KeyAllowedTopLevelClaims))
}

func (p *Provider) SubjectTypesSupported() []string {
	types := stringslice.Filter(
		p.p.StringsF(KeySubjectTypesSupported, []string{"public"}),
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

func (p *Provider) DefaultClientScope() []string {
	return p.p.StringsF(
		KeyDefaultClientScope,
		[]string{"offline_access", "offline", "openid"},
	)
}

func (p *Provider) DSN() string {
	dsn := p.p.String(KeyDSN)

	if dsn == DSNMemory {
		return dbal.SQLiteInMemory
	}

	if len(dsn) > 0 {
		return dsn
	}

	p.l.Fatal("dsn must be set")
	return ""
}

func (p *Provider) EncryptSessionData() bool {
	return p.p.BoolF(KeyEncryptSessionData, true)
}

func (p *Provider) ExcludeNotBeforeClaim() bool {
	return p.p.BoolF(KeyExcludeNotBeforeClaim, false)
}

func (p *Provider) DataSourcePlugin() string {
	return p.p.String(KeyDSN)
}

func (p *Provider) BCryptCost() int {
	return p.p.IntF(KeyBCryptCost, 10)
}

func (p *Provider) CookieSameSiteMode() http.SameSite {
	sameSiteModeStr := p.p.String(KeyCookieSameSiteMode)
	switch strings.ToLower(sameSiteModeStr) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		if tls := p.TLS(PublicInterface); !tls.Enabled() {
			return http.SameSiteLaxMode
		}
		return http.SameSiteNoneMode
	default:
		if tls := p.TLS(PublicInterface); !tls.Enabled() {
			return http.SameSiteLaxMode
		}
		return http.SameSiteDefaultMode
	}
}

func (p *Provider) PublicAllowDynamicRegistration() bool {
	return p.p.Bool(KeyPublicAllowDynamicRegistration)
}

func (p *Provider) CookieSameSiteLegacyWorkaround() bool {
	return p.p.Bool(KeyCookieSameSiteLegacyWorkaround)
}

func (p *Provider) ConsentRequestMaxAge() time.Duration {
	return p.p.DurationF(KeyConsentRequestMaxAge, time.Minute*30)
}

func (p *Provider) AccessTokenLifespan() time.Duration {
	return p.p.DurationF(KeyAccessTokenLifespan, time.Hour)
}

func (p *Provider) RefreshTokenLifespan() time.Duration {
	return p.p.DurationF(KeyRefreshTokenLifespan, time.Hour*720)
}

func (p *Provider) IDTokenLifespan() time.Duration {
	return p.p.DurationF(KeyIDTokenLifespan, time.Hour)
}

func (p *Provider) AuthCodeLifespan() time.Duration {
	return p.p.DurationF(KeyAuthCodeLifespan, time.Minute*10)
}

func (p *Provider) ScopeStrategy() string {
	return p.p.String(KeyScopeStrategy)
}

func (p *Provider) Tracing() *tracing.Config {
	return p.p.TracingConfig("Ory Hydra")
}

func (p *Provider) GetCookieSecrets() [][]byte {
	secrets := p.p.Strings(KeyGetCookieSecrets)
	if len(secrets) == 0 {
		return [][]byte{p.GetSystemSecret()}
	}

	bs := make([][]byte, len(secrets))
	for k := range secrets {
		bs[k] = []byte(secrets[k])
	}
	return bs
}

func (p *Provider) GetRotatedSystemSecrets() [][]byte {
	secrets := p.p.Strings(KeyGetSystemSecret)

	if len(secrets) < 2 {
		return nil
	}

	var rotated [][]byte
	for _, secret := range secrets[1:] {
		rotated = append(rotated, x.HashStringSecret(secret))
	}

	return rotated
}

func (p *Provider) GetSystemSecret() []byte {
	secrets := p.p.Strings(KeyGetSystemSecret)

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

func (p *Provider) LogoutRedirectURL() *url.URL {
	return urlRoot(p.p.RequestURIF(KeyLogoutRedirectURL, p.publicFallbackURL("oauth2/fallbacks/logout/callback")))
}

func (p *Provider) publicFallbackURL(path string) *url.URL {
	if len(p.PublicURL().String()) > 0 {
		return urlx.AppendPaths(p.PublicURL(), path)
	}
	return p.fallbackURL(path, p.host(PublicInterface), p.port(PublicInterface))
}

func (p *Provider) fallbackURL(path string, host string, port int) *url.URL {
	var u url.URL
	u.Scheme = "http"
	if tls := p.TLS(PublicInterface); tls.Enabled() {
		u.Scheme = "https"
	}
	if host == "" {
		u.Host = fmt.Sprintf("%s:%d", "localhost", port)
	}
	u.Path = path
	return &u
}

func (p *Provider) LoginURL() *url.URL {
	return urlRoot(p.p.URIF(KeyLoginURL, p.publicFallbackURL("oauth2/fallbacks/login")))
}

func (p *Provider) LogoutURL() *url.URL {
	return urlRoot(p.p.RequestURIF(KeyLogoutURL, p.publicFallbackURL("oauth2/fallbacks/logout")))
}

func (p *Provider) ConsentURL() *url.URL {
	return urlRoot(p.p.URIF(KeyConsentURL, p.publicFallbackURL("oauth2/fallbacks/consent")))
}

func (p *Provider) ErrorURL() *url.URL {
	return urlRoot(p.p.RequestURIF(KeyErrorURL, p.publicFallbackURL("oauth2/fallbacks/error")))
}

func (p *Provider) PublicURL() *url.URL {
	return urlRoot(p.p.RequestURIF(KeyPublicURL, p.IssuerURL()))
}

func (p *Provider) IssuerURL() *url.URL {
	issuerURL := p.p.RequestURIF(KeyIssuerURL, p.fallbackURL("/", p.host(PublicInterface), p.port(PublicInterface)))
	issuerURL.Path = strings.TrimRight(issuerURL.Path, "/") + "/"
	return urlRoot(issuerURL)
}

func (p *Provider) OAuth2ClientRegistrationURL() *url.URL {
	return p.p.RequestURIF(KeyOAuth2ClientRegistrationURL, new(url.URL))
}

func (p *Provider) OAuth2TokenURL() *url.URL {
	return p.p.RequestURIF(KeyOAuth2TokenURL, urlx.AppendPaths(p.PublicURL(), "/oauth2/token"))
}

func (p *Provider) OAuth2AuthURL() *url.URL {
	return p.p.RequestURIF(KeyOAuth2AuthURL, urlx.AppendPaths(p.PublicURL(), "/oauth2/auth"))
}

func (p *Provider) JWKSURL() *url.URL {
	return p.p.RequestURIF(KeyJWKSURL, urlx.AppendPaths(p.IssuerURL(), "/.well-known/jwks.json"))
}

func (p *Provider) TokenRefreshHookURL() *url.URL {
	return p.p.URIF(KeyRefreshTokenHookURL, nil)
}

func (p *Provider) AccessTokenStrategy() string {
	return strings.ToLower(p.p.StringF(KeyAccessTokenStrategy, "opaque"))
}

func (p *Provider) SubjectIdentifierAlgorithmSalt() string {
	return p.p.String(KeySubjectIdentifierAlgorithmSalt)
}

func (p *Provider) OIDCDiscoverySupportedClaims() []string {
	return stringslice.Unique(
		append(
			[]string{"sub"},
			p.p.Strings(KeyOIDCDiscoverySupportedClaims)...,
		),
	)
}

func (p *Provider) OIDCDiscoverySupportedScope() []string {
	return stringslice.Unique(
		append(
			[]string{"offline_access", "offline", "openid"},
			p.p.Strings(KeyOIDCDiscoverySupportedScope)...,
		),
	)
}

func (p *Provider) OIDCDiscoveryUserinfoEndpoint() *url.URL {
	return p.p.RequestURIF(KeyOIDCDiscoveryUserinfoEndpoint, urlx.AppendPaths(p.PublicURL(), "/userinfo"))
}

func (p *Provider) ShareOAuth2Debug() bool {
	return p.p.Bool(KeyExposeOAuth2Debug)
}

func (p *Provider) OAuth2LegacyErrors() bool {
	return p.p.Bool(KeyOAuth2LegacyErrors)
}

func (p *Provider) PKCEEnforced() bool {
	return p.p.Bool(KeyPKCEEnforced)
}

func (p *Provider) EnforcePKCEForPublicClients() bool {
	return p.p.Bool(KeyPKCEEnforcedForPublicClients)
}

func (p *Provider) CGroupsV1AutoMaxProcsEnabled() bool {
	return p.p.Bool(KeyCGroupsV1AutoMaxProcsEnabled)
}

func (p *Provider) GrantAllClientCredentialsScopesPerDefault() bool {
	return p.p.Bool(KeyGrantAllClientCredentialsScopesPerDefault)
}

func (p *Provider) HsmEnabled() bool {
	return p.p.Bool(HsmEnabled)
}

func (p *Provider) HsmLibraryPath() string {
	return p.p.String(HsmLibraryPath)
}

func (p *Provider) HsmSlotNumber() *int {
	n := p.p.Int(HsmSlotNumber)
	return &n
}

func (p *Provider) HsmPin() string {
	return p.p.String(HsmPin)
}

func (p *Provider) HsmTokenLabel() string {
	return p.p.String(HsmTokenLabel)
}

func (p *Provider) HsmKeySetPrefix() string {
	return p.p.String(HsmKeySetPrefix)
}

func (p *Provider) GrantTypeJWTBearerIDOptional() bool {
	return p.p.Bool(KeyOAuth2GrantJWTIDOptional)
}

func (p *Provider) GrantTypeJWTBearerIssuedDateOptional() bool {
	return p.p.Bool(KeyOAuth2GrantJWTIssuedDateOptional)
}

func (p *Provider) GrantTypeJWTBearerMaxDuration() time.Duration {
	return p.p.DurationF(KeyOAuth2GrantJWTMaxDuration, time.Hour*24*30)
}
