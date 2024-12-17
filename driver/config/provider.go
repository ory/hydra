// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/hasherx"

	"github.com/gofrs/uuid"

	"github.com/ory/x/otelx"

	"github.com/ory/hydra/v2/spec"
	"github.com/ory/x/dbal"

	"github.com/ory/x/configx"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/urlx"
)

const (
	KeyRoot                                      = ""
	HSMEnabled                                   = "hsm.enabled"
	HSMLibraryPath                               = "hsm.library"
	HSMPin                                       = "hsm.pin"
	HSMSlotNumber                                = "hsm.slot"
	HSMKeySetPrefix                              = "hsm.key_set_prefix"
	HSMTokenLabel                                = "hsm.token_label" // #nosec G101
	KeyWellKnownKeys                             = "webfinger.jwks.broadcast_keys"
	KeyOAuth2ClientRegistrationURL               = "webfinger.oidc_discovery.client_registration_url"
	KeyOAuth2TokenURL                            = "webfinger.oidc_discovery.token_url" // #nosec G101
	KeyOAuth2AuthURL                             = "webfinger.oidc_discovery.auth_url"
	KeyVerifiableCredentialsURL                  = "webfinger.oidc_discovery.verifiable_credentials_url" // #nosec G101
	KeyJWKSURL                                   = "webfinger.oidc_discovery.jwks_url"
	KeyOIDCDiscoverySupportedClaims              = "webfinger.oidc_discovery.supported_claims"
	KeyOIDCDiscoverySupportedScope               = "webfinger.oidc_discovery.supported_scope"
	KeyOIDCDiscoveryUserinfoEndpoint             = "webfinger.oidc_discovery.userinfo_url"
	KeySubjectTypesSupported                     = "oidc.subject_identifiers.supported_types"
	KeyDefaultClientScope                        = "oidc.dynamic_client_registration.default_scope"
	KeyDSN                                       = "dsn"
	KeyClientHTTPNoPrivateIPRanges               = "clients.http.disallow_private_ip_ranges"
	KeyClientHTTPPrivateIPExceptionURLs          = "clients.http.private_ip_exception_urls"
	KeyHasherAlgorithm                           = "oauth2.hashers.algorithm"
	KeyBCryptCost                                = "oauth2.hashers.bcrypt.cost"
	KeyPBKDF2Iterations                          = "oauth2.hashers.pbkdf2.iterations"
	KeyEncryptSessionData                        = "oauth2.session.encrypt_at_rest"
	KeyCookieSameSiteMode                        = "serve.cookies.same_site_mode"
	KeyCookieSameSiteLegacyWorkaround            = "serve.cookies.same_site_legacy_workaround"
	KeyCookieDomain                              = "serve.cookies.domain"
	KeyCookieSecure                              = "serve.cookies.secure"
	KeyCookieLoginCSRFName                       = "serve.cookies.names.login_csrf"
	KeyCookieConsentCSRFName                     = "serve.cookies.names.consent_csrf"
	KeyCookieSessionName                         = "serve.cookies.names.session"
	KeyCookieSessionPath                         = "serve.cookies.paths.session"
	KeyConsentRequestMaxAge                      = "ttl.login_consent_request"
	KeyAccessTokenLifespan                       = "ttl.access_token"  // #nosec G101
	KeyRefreshTokenLifespan                      = "ttl.refresh_token" // #nosec G101
	KeyVerifiableCredentialsNonceLifespan        = "ttl.vc_nonce"      // #nosec G101
	KeyIDTokenLifespan                           = "ttl.id_token"      // #nosec G101
	KeyAuthCodeLifespan                          = "ttl.auth_code"
	KeyScopeStrategy                             = "strategies.scope"
	KeyGetCookieSecrets                          = "secrets.cookie"
	KeyGetSystemSecret                           = "secrets.system"
	KeyLogoutRedirectURL                         = "urls.post_logout_redirect"
	KeyLoginURL                                  = "urls.login"
	KeyRegistrationURL                           = "urls.registration"
	KeyLogoutURL                                 = "urls.logout"
	KeyConsentURL                                = "urls.consent"
	KeyErrorURL                                  = "urls.error"
	KeyPublicURL                                 = "urls.self.public"
	KeyAdminURL                                  = "urls.self.admin"
	KeyIssuerURL                                 = "urls.self.issuer"
	KeyIdentityProviderAdminURL                  = "urls.identity_provider.url"
	KeyIdentityProviderPublicURL                 = "urls.identity_provider.publicUrl"
	KeyIdentityProviderHeaders                   = "urls.identity_provider.headers"
	KeyAccessTokenStrategy                       = "strategies.access_token"
	KeyJWTScopeClaimStrategy                     = "strategies.jwt.scope_claim"
	KeyDBIgnoreUnknownTableColumns               = "db.ignore_unknown_table_columns"
	KeySubjectIdentifierAlgorithmSalt            = "oidc.subject_identifiers.pairwise.salt"
	KeyPublicAllowDynamicRegistration            = "oidc.dynamic_client_registration.enabled"
	KeyPKCEEnforced                              = "oauth2.pkce.enforced"
	KeyPKCEEnforcedForPublicClients              = "oauth2.pkce.enforced_for_public_clients"
	KeyLogLevel                                  = "log.level"
	KeyCGroupsV1AutoMaxProcsEnabled              = "cgroups.v1.auto_max_procs_enabled"
	KeyGrantAllClientCredentialsScopesPerDefault = "oauth2.client_credentials.default_grant_allowed_scope" // #nosec G101
	KeyExposeOAuth2Debug                         = "oauth2.expose_internal_errors"
	KeyExcludeNotBeforeClaim                     = "oauth2.exclude_not_before_claim"
	KeyAllowedTopLevelClaims                     = "oauth2.allowed_top_level_claims"
	KeyMirrorTopLevelClaims                      = "oauth2.mirror_top_level_claims"
	KeyRefreshTokenRotationGracePeriod           = "oauth2.grant.refresh_token.rotation_grace_period" // #nosec G101
	KeyOAuth2GrantJWTIDOptional                  = "oauth2.grant.jwt.jti_optional"
	KeyOAuth2GrantJWTIssuedDateOptional          = "oauth2.grant.jwt.iat_optional"
	KeyOAuth2GrantJWTMaxDuration                 = "oauth2.grant.jwt.max_ttl"
	KeyRefreshTokenHook                          = "oauth2.refresh_token_hook" // #nosec G101
	KeyTokenHook                                 = "oauth2.token_hook"         // #nosec G101
	KeyDevelopmentMode                           = "dev"
)

const DSNMemory = "memory"

var (
	_ hasherx.PBKDF2Configurator = (*DefaultProvider)(nil)
	_ hasherx.BCryptConfigurator = (*DefaultProvider)(nil)
)

type DefaultProvider struct {
	l *logrusx.Logger
	p *configx.Provider
	c contextx.Contextualizer
}

func (p *DefaultProvider) GetHasherAlgorithm(ctx context.Context) x.HashAlgorithm {
	switch strings.ToLower(p.getProvider(ctx).String(KeyHasherAlgorithm)) {
	case x.HashAlgorithmBCrypt.String():
		return x.HashAlgorithmBCrypt
	case x.HashAlgorithmPBKDF2.String():
		fallthrough
	default:
		return x.HashAlgorithmPBKDF2
	}
}

func (p *DefaultProvider) HasherBcryptConfig(ctx context.Context) *hasherx.BCryptConfig {
	var cost uint32
	costInt := int64(p.GetBCryptCost(ctx))
	if costInt < 0 {
		cost = 10
	} else if costInt > math.MaxUint32 {
		cost = math.MaxUint32
	} else {
		cost = uint32(costInt)
	}
	return &hasherx.BCryptConfig{
		Cost: cost,
	}
}

func (p *DefaultProvider) HasherPBKDF2Config(ctx context.Context) *hasherx.PBKDF2Config {
	var iters uint32
	itersInt := p.getProvider(ctx).Int64(KeyPBKDF2Iterations)
	if itersInt < 1 {
		iters = 1
	} else if int64(itersInt) > math.MaxUint32 {
		iters = math.MaxUint32
	} else {
		iters = uint32(itersInt)
	}

	return &hasherx.PBKDF2Config{
		Algorithm:  "sha256",
		Iterations: iters,
		SaltLength: 16,
		KeyLength:  32,
	}
}

func MustNew(ctx context.Context, l *logrusx.Logger, opts ...configx.OptionModifier) *DefaultProvider {
	p, err := New(ctx, l, opts...)
	if err != nil {
		l.WithError(err).Fatalf("Unable to load config.")
	}
	return p
}

func (p *DefaultProvider) getProvider(ctx context.Context) *configx.Provider {
	return p.c.Config(ctx, p.p)
}

func New(ctx context.Context, l *logrusx.Logger, opts ...configx.OptionModifier) (*DefaultProvider, error) {
	opts = append(
		[]configx.OptionModifier{
			configx.WithStderrValidationReporter(),
			configx.OmitKeysFromTracing("dsn", "secrets.system", "secrets.cookie"),
			configx.WithImmutables("log", "serve", "dsn", "profiling"),
			configx.WithLogrusWatcher(l),
		}, opts...,
	)

	p, err := configx.New(ctx, spec.ConfigValidationSchema, opts...)
	if err != nil {
		return nil, err
	}
	return NewCustom(l, p, &contextx.Default{}), nil
}

func NewCustom(l *logrusx.Logger, p *configx.Provider, ctxt contextx.Contextualizer) *DefaultProvider {
	l.UseConfig(p)
	return &DefaultProvider{l: l, p: p, c: ctxt}
}

func (p *DefaultProvider) Set(ctx context.Context, key string, value interface{}) error {
	return p.getProvider(ctx).Set(key, value)
}

func (p *DefaultProvider) MustSet(ctx context.Context, key string, value interface{}) {
	if err := p.Set(ctx, key, value); err != nil {
		p.l.WithError(err).Fatalf("Unable to set \"%s\" to \"%s\".", key, value)
	}
}

func (p *DefaultProvider) Delete(ctx context.Context, key string) {
	p.getProvider(ctx).Delete(key)
}

func (p *DefaultProvider) Source(ctx context.Context) *configx.Provider {
	return p.getProvider(ctx)
}

func (p *DefaultProvider) IsDevelopmentMode(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyDevelopmentMode)
}

func (p *DefaultProvider) WellKnownKeys(ctx context.Context, include ...string) []string {
	include = append(include, x.OAuth2JWTKeyName, x.OpenIDConnectKeyName)
	return stringslice.Unique(append(p.getProvider(ctx).Strings(KeyWellKnownKeys), include...))
}

func (p *DefaultProvider) ClientHTTPNoPrivateIPRanges() bool {
	return p.getProvider(contextx.RootContext).Bool(KeyClientHTTPNoPrivateIPRanges)
}

func (p *DefaultProvider) ClientHTTPPrivateIPExceptionURLs() []string {
	return p.getProvider(contextx.RootContext).Strings(KeyClientHTTPPrivateIPExceptionURLs)
}

func (p *DefaultProvider) AllowedTopLevelClaims(ctx context.Context) []string {
	return stringslice.Unique(p.getProvider(ctx).Strings(KeyAllowedTopLevelClaims))
}

func (p *DefaultProvider) MirrorTopLevelClaims(ctx context.Context) bool {
	return p.getProvider(ctx).BoolF(KeyMirrorTopLevelClaims, true)
}

func (p *DefaultProvider) SubjectTypesSupported(ctx context.Context, additionalSources ...AccessTokenStrategySource) []string {
	types := stringslice.Filter(
		p.getProvider(ctx).StringsF(KeySubjectTypesSupported, []string{"public"}),
		func(s string) bool {
			return !(s == "public" || s == "pairwise")
		},
	)

	if len(types) == 0 {
		types = []string{"public"}
	}

	if stringslice.Has(types, "pairwise") {
		if p.AccessTokenStrategy(ctx, additionalSources...) == AccessTokenJWTStrategy {
			p.l.Warn(`The pairwise subject identifier algorithm is not supported by the JWT OAuth 2.0 Access Token Strategy and is thus being disabled. Please remove "pairwise" from oidc.subject_identifiers.supported_types" (e.g. oidc.subject_identifiers.supported_types=public) or set strategies.access_token to "opaque".`)
			types = stringslice.Filter(
				types,
				func(s string) bool {
					return !(s == "public")
				},
			)
		} else if len(p.SubjectIdentifierAlgorithmSalt(ctx)) < 8 {
			p.l.Fatalf(
				`The pairwise subject identifier algorithm was set but length of oidc.subject_identifier.salt is too small (%d < 8), please set oidc.subject_identifiers.pairwise.salt to a random string with 8 characters or more.`,
				len(p.SubjectIdentifierAlgorithmSalt(ctx)),
			)
		}
	}

	return types
}

func (p *DefaultProvider) DefaultClientScope(ctx context.Context) []string {
	return p.getProvider(ctx).StringsF(
		KeyDefaultClientScope,
		[]string{"offline_access", "offline", "openid"},
	)
}

func (p *DefaultProvider) DSN() string {
	dsn := p.getProvider(contextx.RootContext).String(KeyDSN)

	if dsn == DSNMemory {
		return dbal.NewSQLiteInMemoryDatabase(uuid.Must(uuid.NewV4()).String())
	}

	if len(dsn) > 0 {
		return dsn
	}

	p.l.Fatal("dsn must be set")
	return ""
}

func (p *DefaultProvider) EncryptSessionData(ctx context.Context) bool {
	return p.getProvider(ctx).BoolF(KeyEncryptSessionData, true)
}

func (p *DefaultProvider) ExcludeNotBeforeClaim(ctx context.Context) bool {
	return p.getProvider(ctx).BoolF(KeyExcludeNotBeforeClaim, false)
}

func (p *DefaultProvider) CookieSecure(ctx context.Context) bool {
	if !p.IsDevelopmentMode(ctx) {
		return true
	}
	return p.getProvider(ctx).BoolF(KeyCookieSecure, false)
}

func (p *DefaultProvider) CookieSameSiteMode(ctx context.Context) http.SameSite {
	sameSiteModeStr := p.getProvider(ctx).String(KeyCookieSameSiteMode)
	switch strings.ToLower(sameSiteModeStr) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		if p.IssuerURL(ctx).Scheme != "https" {
			// SameSite=None can only be set for HTTPS issuers.
			return http.SameSiteLaxMode
		}
		return http.SameSiteNoneMode
	default:
		if p.IsDevelopmentMode(ctx) {
			return http.SameSiteLaxMode
		}
		return http.SameSiteDefaultMode
	}
}

func (p *DefaultProvider) PublicAllowDynamicRegistration(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyPublicAllowDynamicRegistration)
}

func (p *DefaultProvider) CookieSameSiteLegacyWorkaround(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyCookieSameSiteLegacyWorkaround)
}

func (p *DefaultProvider) ConsentRequestMaxAge(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyConsentRequestMaxAge, time.Minute*30)
}

func (p *DefaultProvider) Tracing() *otelx.Config {
	return p.getProvider(contextx.RootContext).TracingConfig("Ory Hydra")
}

func (p *DefaultProvider) GetCookieSecrets(ctx context.Context) ([][]byte, error) {
	secrets := p.getProvider(ctx).Strings(KeyGetCookieSecrets)
	if len(secrets) == 0 {
		secret, err := p.GetGlobalSecret(ctx)
		if err != nil {
			return nil, err
		}
		return [][]byte{secret}, nil
	}

	bs := make([][]byte, len(secrets))
	for k := range secrets {
		bs[k] = []byte(secrets[k])
	}
	return bs, nil
}

func (p *DefaultProvider) LogoutRedirectURL(ctx context.Context) *url.URL {
	return urlRoot(
		p.getProvider(ctx).RequestURIF(
			KeyLogoutRedirectURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/logout/callback"),
		),
	)
}

func (p *DefaultProvider) publicFallbackURL(ctx context.Context, path string) *url.URL {
	if len(p.PublicURL(ctx).String()) > 0 {
		return urlx.AppendPaths(p.PublicURL(ctx), path)
	}
	return p.fallbackURL(ctx, path, p.host(PublicInterface), p.port(PublicInterface))
}

func (p *DefaultProvider) fallbackURL(ctx context.Context, path string, host string, port int) *url.URL {
	var u url.URL
	u.Scheme = "http"
	if tls := p.TLS(ctx, PublicInterface); tls.Enabled() || !p.IsDevelopmentMode(ctx) {
		u.Scheme = "https"
	}
	if host == "" {
		u.Host = fmt.Sprintf("%s:%d", "localhost", port)
	}
	u.Path = path
	return &u
}

func (p *DefaultProvider) LoginURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).URIF(KeyLoginURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/login")))
}

func (p *DefaultProvider) RegistrationURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).URIF(KeyRegistrationURL, p.LoginURL(ctx)))
}

func (p *DefaultProvider) LogoutURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).RequestURIF(KeyLogoutURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/logout")))
}

func (p *DefaultProvider) ConsentURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).URIF(KeyConsentURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/consent")))
}

func (p *DefaultProvider) ErrorURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).RequestURIF(KeyErrorURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/error")))
}

func (p *DefaultProvider) PublicURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).RequestURIF(KeyPublicURL, p.IssuerURL(ctx)))
}

func (p *DefaultProvider) AdminURL(ctx context.Context) *url.URL {
	return urlRoot(
		p.getProvider(ctx).RequestURIF(
			KeyAdminURL, p.fallbackURL(ctx, "/", p.host(AdminInterface), p.port(AdminInterface)),
		),
	)
}

func (p *DefaultProvider) IssuerURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(
		KeyIssuerURL, p.fallbackURL(ctx, "/", p.host(PublicInterface), p.port(PublicInterface)),
	)
}

func (p *DefaultProvider) KratosAdminURL(ctx context.Context) (*url.URL, bool) {
	u := p.getProvider(ctx).RequestURIF(KeyIdentityProviderAdminURL, nil)

	return u, u != nil
}
func (p *DefaultProvider) KratosPublicURL(ctx context.Context) (*url.URL, bool) {
	u := p.getProvider(ctx).RequestURIF(KeyIdentityProviderPublicURL, nil)

	return u, u != nil
}

func (p *DefaultProvider) KratosRequestHeader(ctx context.Context) http.Header {
	hh := map[string]string{}
	if err := p.getProvider(ctx).Unmarshal(KeyIdentityProviderHeaders, &hh); err != nil {
		p.l.WithError(errors.WithStack(err)).
			Errorf("Configuration value from key %s could not be decoded.", KeyIdentityProviderHeaders)
		return nil
	}

	h := make(http.Header)
	for k, v := range hh {
		h.Set(k, v)
	}

	return h
}

func (p *DefaultProvider) OAuth2ClientRegistrationURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyOAuth2ClientRegistrationURL, new(url.URL))
}

func (p *DefaultProvider) OAuth2TokenURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyOAuth2TokenURL, urlx.AppendPaths(p.PublicURL(ctx), "/oauth2/token"))
}

func (p *DefaultProvider) OAuth2AuthURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyOAuth2AuthURL, urlx.AppendPaths(p.PublicURL(ctx), "/oauth2/auth"))
}

func (p *DefaultProvider) JWKSURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyJWKSURL, urlx.AppendPaths(p.IssuerURL(ctx), "/.well-known/jwks.json"))
}

func (p *DefaultProvider) CredentialsEndpointURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyVerifiableCredentialsURL, urlx.AppendPaths(p.PublicURL(ctx), "/credentials"))
}

type AccessTokenStrategySource interface {
	GetAccessTokenStrategy() AccessTokenStrategyType
}

func (p *DefaultProvider) AccessTokenStrategy(ctx context.Context, additionalSources ...AccessTokenStrategySource) AccessTokenStrategyType {
	for _, src := range additionalSources {
		if src == nil {
			continue
		}
		if strategy := src.GetAccessTokenStrategy(); strategy != "" {
			return strategy
		}
	}
	s, err := ToAccessTokenStrategyType(p.getProvider(ctx).String(KeyAccessTokenStrategy))
	if err != nil {
		p.l.WithError(err).Warn("Key `strategies.access_token` contains an invalid value, falling back to `opaque` strategy.")
		return AccessTokenDefaultStrategy
	}

	return s
}

type (
	Auth struct {
		Type   string     `json:"type"`
		Config AuthConfig `json:"config"`
	}
	AuthConfig struct {
		In    string `json:"in"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	HookConfig struct {
		URL  string `json:"url"`
		Auth *Auth  `json:"auth"`
	}
)

func (p *DefaultProvider) getHookConfig(ctx context.Context, key string) *HookConfig {
	if p.getProvider(ctx).String(key) == "" {
		return nil
	}

	if hookURL := p.getProvider(ctx).RequestURIF(key, nil); hookURL != nil {
		return &HookConfig{
			URL: hookURL.String(),
		}
	}

	var hookConfig *HookConfig
	if err := p.getProvider(ctx).Unmarshal(key, &hookConfig); err != nil {
		p.l.WithError(errors.WithStack(err)).
			Errorf("Configuration value from key %s could not be decoded.", key)
		return nil
	}
	if hookConfig == nil {
		return nil
	}

	// validate URL by parsing it
	u, err := url.ParseRequestURI(hookConfig.URL)
	if err != nil {
		p.l.WithError(errors.WithStack(err)).
			Errorf("Configuration value from key %s could not be decoded.", key)
		return nil
	}
	hookConfig.URL = u.String()

	return hookConfig
}

func (p *DefaultProvider) TokenHookConfig(ctx context.Context) *HookConfig {
	return p.getHookConfig(ctx, KeyTokenHook)
}

func (p *DefaultProvider) TokenRefreshHookConfig(ctx context.Context) *HookConfig {
	return p.getHookConfig(ctx, KeyRefreshTokenHook)
}

func (p *DefaultProvider) DbIgnoreUnknownTableColumns() bool {
	return p.p.Bool(KeyDBIgnoreUnknownTableColumns)
}

func (p *DefaultProvider) SubjectIdentifierAlgorithmSalt(ctx context.Context) string {
	return p.getProvider(ctx).String(KeySubjectIdentifierAlgorithmSalt)
}

func (p *DefaultProvider) OIDCDiscoverySupportedClaims(ctx context.Context) []string {
	return stringslice.Unique(
		append(
			[]string{"sub"},
			p.getProvider(ctx).Strings(KeyOIDCDiscoverySupportedClaims)...,
		),
	)
}

func (p *DefaultProvider) OIDCDiscoverySupportedScope(ctx context.Context) []string {
	return stringslice.Unique(
		append(
			[]string{"offline_access", "offline", "openid"},
			p.getProvider(ctx).Strings(KeyOIDCDiscoverySupportedScope)...,
		),
	)
}

func (p *DefaultProvider) OIDCDiscoveryUserinfoEndpoint(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(
		KeyOIDCDiscoveryUserinfoEndpoint, urlx.AppendPaths(p.PublicURL(ctx), "/userinfo"),
	)
}

func (p *DefaultProvider) GetSendDebugMessagesToClients(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyExposeOAuth2Debug)
}

func (p *DefaultProvider) GetEnforcePKCE(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyPKCEEnforced)
}

func (p *DefaultProvider) GetEnforcePKCEForPublicClients(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyPKCEEnforcedForPublicClients)
}

func (p *DefaultProvider) CGroupsV1AutoMaxProcsEnabled() bool {
	return p.getProvider(contextx.RootContext).Bool(KeyCGroupsV1AutoMaxProcsEnabled)
}

func (p *DefaultProvider) GrantAllClientCredentialsScopesPerDefault(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyGrantAllClientCredentialsScopesPerDefault)
}

func (p *DefaultProvider) HSMEnabled() bool {
	return p.getProvider(contextx.RootContext).Bool(HSMEnabled)
}

func (p *DefaultProvider) HSMLibraryPath() string {
	return p.getProvider(contextx.RootContext).String(HSMLibraryPath)
}

func (p *DefaultProvider) HSMSlotNumber() *int {
	n := p.getProvider(contextx.RootContext).Int(HSMSlotNumber)
	return &n
}

func (p *DefaultProvider) HSMPin() string {
	return p.getProvider(contextx.RootContext).String(HSMPin)
}

func (p *DefaultProvider) HSMTokenLabel() string {
	return p.getProvider(contextx.RootContext).String(HSMTokenLabel)
}

func (p *DefaultProvider) GetGrantTypeJWTBearerIDOptional(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyOAuth2GrantJWTIDOptional)
}

func (p *DefaultProvider) HSMKeySetPrefix() string {
	return p.getProvider(contextx.RootContext).String(HSMKeySetPrefix)
}

func (p *DefaultProvider) GetGrantTypeJWTBearerIssuedDateOptional(ctx context.Context) bool {
	return p.getProvider(ctx).Bool(KeyOAuth2GrantJWTIssuedDateOptional)
}

func (p *DefaultProvider) GetJWTMaxDuration(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyOAuth2GrantJWTMaxDuration, time.Hour*24*30)
}

func (p *DefaultProvider) CookieDomain(ctx context.Context) string {
	return p.getProvider(ctx).String(KeyCookieDomain)
}

func (p *DefaultProvider) SessionCookiePath(ctx context.Context) string {
	return p.getProvider(ctx).StringF(KeyCookieSessionPath, "/")
}

func (p *DefaultProvider) CookieNameLoginCSRF(ctx context.Context) string {
	return p.cookieSuffix(ctx, KeyCookieLoginCSRFName)
}

func (p *DefaultProvider) CookieNameConsentCSRF(ctx context.Context) string {
	return p.cookieSuffix(ctx, KeyCookieConsentCSRFName)
}

func (p *DefaultProvider) SessionCookieName(ctx context.Context) string {
	return p.cookieSuffix(ctx, KeyCookieSessionName)
}

func (p *DefaultProvider) cookieSuffix(ctx context.Context, key string) string {
	var suffix string
	if p.IsDevelopmentMode(ctx) {
		suffix = "_dev"
	}

	return p.getProvider(ctx).String(key) + suffix
}

func (p *DefaultProvider) RefreshTokenRotationGracePeriod(ctx context.Context) time.Duration {
	gracePeriod := p.getProvider(ctx).DurationF(KeyRefreshTokenRotationGracePeriod, 0)
	if gracePeriod > time.Minute*5 {
		return time.Minute * 5
	}
	return gracePeriod
}
