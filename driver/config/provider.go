// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"crypto/sha512"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/spec"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/hasherx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/randx"
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
	KeyOAuth2DeviceAuthorisationURL              = "webfinger.oidc_discovery.device_authorization_url"
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
	KeyCookieDeviceCSRFName                      = "serve.cookies.names.device_csrf"
	KeyCookieConsentCSRFName                     = "serve.cookies.names.consent_csrf"
	KeyCookieSessionName                         = "serve.cookies.names.session"
	KeyCookieSessionPath                         = "serve.cookies.paths.session"
	KeyConsentRequestMaxAge                      = "ttl.login_consent_request"
	KeyAccessTokenLifespan                       = "ttl.access_token"  // #nosec G101
	KeyRefreshTokenLifespan                      = "ttl.refresh_token" // #nosec G101
	KeyVerifiableCredentialsNonceLifespan        = "ttl.vc_nonce"      // #nosec G101
	KeyIDTokenLifespan                           = "ttl.id_token"      // #nosec G101
	KeyAuthCodeLifespan                          = "ttl.auth_code"
	KeyDeviceAndUserCodeLifespan                 = "ttl.device_user_code"
	KeyAuthenticationSessionLifespan             = "ttl.authentication_session"
	KeyScopeStrategy                             = "strategies.scope"
	KeyGetCookieSecrets                          = "secrets.cookie"
	KeyGetSystemSecret                           = "secrets.system"
	KeyPaginationSecrets                         = "secrets.pagination"
	KeyLogoutRedirectURL                         = "urls.post_logout_redirect"
	KeyLoginURL                                  = "urls.login"
	KeyRegistrationURL                           = "urls.registration"
	KeyLogoutURL                                 = "urls.logout"
	KeyConsentURL                                = "urls.consent"
	KeyErrorURL                                  = "urls.error"
	KeyDeviceVerificationURL                     = "urls.device.verification"
	KeyDeviceDoneURL                             = "urls.device.success"
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
	KeyDeviceAuthTokenPollingInterval            = "oauth2.device_authorization.token_polling_interval" // #nosec G101
	KeyDeviceAuthUserCodeEntropyPreset           = "oauth2.device_authorization.user_code.entropy_preset"
	KeyDeviceAuthUserCodeLength                  = "oauth2.device_authorization.user_code.length"
	KeyDeviceAuthUserCodeCharacterSet            = "oauth2.device_authorization.user_code.character_set"
	KeyPKCEEnforced                              = "oauth2.pkce.enforced"
	KeyPKCEEnforcedForPublicClients              = "oauth2.pkce.enforced_for_public_clients"
	KeyLogLevel                                  = "log.level"
	KeyCGroupsV1AutoMaxProcsEnabled              = "cgroups.v1.auto_max_procs_enabled"
	KeyGrantAllClientCredentialsScopesPerDefault = "oauth2.client_credentials.default_grant_allowed_scope" // #nosec G101
	KeyExposeOAuth2Debug                         = "oauth2.expose_internal_errors"
	KeyExcludeNotBeforeClaim                     = "oauth2.exclude_not_before_claim"
	KeyAllowedTopLevelClaims                     = "oauth2.allowed_top_level_claims"
	KeyMirrorTopLevelClaims                      = "oauth2.mirror_top_level_claims"
	KeyRefreshTokenRotationGracePeriod           = "oauth2.grant.refresh_token.rotation_grace_period"      // #nosec G101
	KeyRefreshTokenRotationGraceReuseCount       = "oauth2.grant.refresh_token.rotation_grace_reuse_count" // #nosec G101
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

func (p *DefaultProvider) GetHasherAlgorithm(ctx context.Context) string {
	return strings.ToLower(p.getProvider(ctx).String(KeyHasherAlgorithm))
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

func MustNew(t testing.TB, l *logrusx.Logger, opts ...configx.OptionModifier) *DefaultProvider {
	ctxt := contextx.NewTestConfigProvider(spec.ConfigValidationSchema, opts...)
	p, err := New(t.Context(), l, ctxt, opts...)
	require.NoError(t, err)
	return p
}

func (p *DefaultProvider) getProvider(ctx context.Context) *configx.Provider {
	return p.c.Config(ctx, p.p)
}

func New(ctx context.Context, l *logrusx.Logger, ctxt contextx.Contextualizer, opts ...configx.OptionModifier) (*DefaultProvider, error) {
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
	return NewCustom(l, p, ctxt), nil
}

func NewCustom(l *logrusx.Logger, p *configx.Provider, ctxt contextx.Contextualizer) *DefaultProvider {
	l.UseConfig(p)
	return &DefaultProvider{l: l, p: p, c: ctxt}
}

// Deprecated: use context-based test setters
func (p *DefaultProvider) Set(ctx context.Context, key string, value interface{}) error {
	return p.getProvider(ctx).Set(key, value)
}

// Deprecated: use context-based test setters
func (p *DefaultProvider) MustSet(ctx context.Context, key string, value interface{}) {
	if err := p.Set(ctx, key, value); err != nil {
		p.l.WithError(err).Fatalf("Unable to set \"%s\" to \"%s\".", key, value)
	}
}

// Deprecated: use context-based test setters
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
	public, pairwise := false, false
	for _, t := range p.getProvider(ctx).StringsF(KeySubjectTypesSupported, []string{"public"}) {
		switch t {
		case "public":
			public = true
		case "pairwise":
			pairwise = true
		}
	}

	// when neither public nor pairwise are set, force public
	public = public || !pairwise

	if pairwise {
		if p.AccessTokenStrategy(ctx, additionalSources...) == AccessTokenJWTStrategy {
			p.l.Warn(`The pairwise subject identifier algorithm is not supported by the JWT OAuth 2.0 Access Token Strategy and is thus being disabled. Please remove "pairwise" from oidc.subject_identifiers.supported_types" (e.g. oidc.subject_identifiers.supported_types=public) or set strategies.access_token to "opaque".`)
			pairwise = false
		} else if len(p.SubjectIdentifierAlgorithmSalt(ctx)) < 8 {
			p.l.Fatalf(
				`The pairwise subject identifier algorithm was set but length of oidc.subject_identifier.salt is too small (%d < 8), please set oidc.subject_identifiers.pairwise.salt to a random string with 8 characters or more.`,
				len(p.SubjectIdentifierAlgorithmSalt(ctx)),
			)
		}
	}

	types := make([]string, 0, 2)
	if public {
		types = append(types, "public")
	}
	if pairwise {
		types = append(types, "pairwise")
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
	if publicURL := p.PublicURL(ctx); len(publicURL.String()) > 0 {
		return urlx.AppendPaths(publicURL, path)
	}
	return p.fallbackURL(ctx, path, p.ServePublic(ctx))
}

func (p *DefaultProvider) fallbackURL(ctx context.Context, path string, serve *configx.Serve) *url.URL {
	u := url.URL{
		Scheme: "http",
		Host:   serve.GetAddress(),
	}
	if serve.TLS.Enabled || !p.IsDevelopmentMode(ctx) {
		u.Scheme = "https"
	}
	if serve.Host == "" {
		u.Host = fmt.Sprintf("%s:%d", "localhost", serve.Port)
	}
	u.Path = path
	return &u
}

// GetDeviceAndUserCodeLifespan returns the device_code and user_code lifespan. Defaults to 15 minutes.
func (p *DefaultProvider) GetDeviceAndUserCodeLifespan(ctx context.Context) time.Duration {
	return p.p.DurationF(KeyDeviceAndUserCodeLifespan, time.Minute*15)
}

// GetAuthenticationSessionLifespan returns the authentication_session lifespan.
func (p *DefaultProvider) GetAuthenticationSessionLifespan(ctx context.Context) time.Duration {
	lifespan := p.p.Duration(KeyAuthenticationSessionLifespan)
	if lifespan > time.Hour*24*180 {
		return time.Hour * 24 * 180
	}
	return lifespan
}

// GetDeviceAuthTokenPollingInterval returns device grant token endpoint polling interval. Defaults to 5 seconds.
func (p *DefaultProvider) GetDeviceAuthTokenPollingInterval(ctx context.Context) time.Duration {
	return p.p.DurationF(KeyDeviceAuthTokenPollingInterval, time.Second*5)
}

func (p *DefaultProvider) userCodeEntropyPreset(t string) (int, []rune) {
	switch t {
	default:
		p.l.Errorf(`invalid user code entropy preset %q, allowed values are "high", "medium", or "low"`, t)
		fallthrough
	case "high":
		return 8, randx.AlphaNumNoAmbiguous
	case "medium":
		return 8, randx.AlphaUpper
	case "low":
		return 9, randx.Numeric
	}
}

// GetUserCodeLength returns configured user_code length
func (p *DefaultProvider) GetUserCodeLength(ctx context.Context) int {
	if l := p.getProvider(ctx).Int(KeyDeviceAuthUserCodeLength); l > 0 {
		return l
	}
	k := p.getProvider(ctx).StringF(KeyDeviceAuthUserCodeEntropyPreset, "high")
	l, _ := p.userCodeEntropyPreset(k)
	return l
}

// GetUserCodeSymbols returns configured user_code allowed symbols
func (p *DefaultProvider) GetUserCodeSymbols(ctx context.Context) []rune {
	if s := p.getProvider(ctx).String(KeyDeviceAuthUserCodeCharacterSet); s != "" {
		return []rune(s)
	}
	k := p.getProvider(ctx).StringF(KeyDeviceAuthUserCodeEntropyPreset, "high")
	_, s := p.userCodeEntropyPreset(k)
	return s
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

// DeviceVerificationURL returns user_code verification page URL. Defaults to "oauth2/fallbacks/device".
func (p *DefaultProvider) DeviceVerificationURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).URIF(KeyDeviceVerificationURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/device")))
}

// DeviceDoneURL returns the post device authorization URL. Defaults to "oauth2/fallbacks/device/done".
func (p *DefaultProvider) DeviceDoneURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).RequestURIF(KeyDeviceDoneURL, p.publicFallbackURL(ctx, "oauth2/fallbacks/device/done")))
}

func (p *DefaultProvider) PublicURL(ctx context.Context) *url.URL {
	return urlRoot(p.getProvider(ctx).RequestURIF(KeyPublicURL, p.IssuerURL(ctx)))
}

func (p *DefaultProvider) AdminURL(ctx context.Context) *url.URL {
	return urlRoot(
		p.getProvider(ctx).RequestURIF(
			KeyAdminURL, p.fallbackURL(ctx, "/", p.ServeAdmin(ctx)),
		),
	)
}

func (p *DefaultProvider) IssuerURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyIssuerURL, p.fallbackURL(ctx, "/", p.ServePublic(ctx)))
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

// OAuth2DeviceAuthorisationURL returns device authorization endpoint. Defaults to "/oauth2/device/auth".
func (p *DefaultProvider) OAuth2DeviceAuthorisationURL(ctx context.Context) *url.URL {
	return p.getProvider(ctx).RequestURIF(KeyOAuth2DeviceAuthorisationURL, urlx.AppendPaths(p.PublicURL(ctx), "/oauth2/device/auth"))
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

// CookieNameDeviceCSRF returns the device CSRF cookie name.
func (p *DefaultProvider) CookieNameDeviceCSRF(ctx context.Context) string {
	return p.cookieSuffix(ctx, KeyCookieDeviceCSRFName)
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

type GracefulRefreshTokenRotation struct {
	Period time.Duration
	Count  int32
}

func (p *DefaultProvider) GracefulRefreshTokenRotation(ctx context.Context) (cfg GracefulRefreshTokenRotation) {
	//nolint:gosec
	cfg.Count = int32(x.Clamp(p.getProvider(ctx).IntF(KeyRefreshTokenRotationGraceReuseCount, 0), 0, math.MaxInt32))

	// The maximum value is 5 minutes, unless also a reuse count is configured, in
	// which case the maximum is 180 days
	maxPeriod := 5 * time.Minute
	if cfg.Count > 0 {
		maxPeriod = 180 * 24 * time.Hour
	}
	cfg.Period = x.Clamp(p.getProvider(ctx).DurationF(KeyRefreshTokenRotationGracePeriod, 0), 0, maxPeriod)

	return
}

func (p *DefaultProvider) GetPaginationEncryptionKeys(ctx context.Context) [][32]byte {
	secrets := p.getProvider(ctx).Strings(KeyPaginationSecrets)
	if len(secrets) == 0 {
		secrets = p.getProvider(ctx).Strings(KeyGetSystemSecret)
	}

	hashed := make([][32]byte, len(secrets))
	for i := range secrets {
		hashed[i] = sha512.Sum512_256([]byte(secrets[i]))
	}
	return hashed
}
