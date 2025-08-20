// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fositex

import (
	"context"
	"crypto/sha512"
	"hash"
	"html/template"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/i18n"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/urlx"
)

type configDependencies interface {
	config.Provider
	persistence.Provider
	x.HTTPClientProvider
	ClientHasher() fosite.Hasher
	ExtraFositeFactories() []Factory
}

type Factory func(config fosite.Configurator, storage interface{}, strategy interface{}) interface{}

type Config struct {
	deps configDependencies

	authorizeEndpointHandlers  fosite.AuthorizeEndpointHandlers
	tokenEndpointHandlers      fosite.TokenEndpointHandlers
	tokenIntrospectionHandlers fosite.TokenIntrospectionHandlers
	revocationHandlers         fosite.RevocationHandlers
	deviceEndpointHandlers     fosite.DeviceEndpointHandlers
	jwksFetcherStrategy        fosite.JWKSFetcherStrategy

	*config.DefaultProvider
}

var defaultResponseModeHandler = fosite.NewDefaultResponseModeHandler()
var defaultFactories = []Factory{
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
	compose.RFC7523AssertionGrantFactory,
	compose.OIDCUserinfoVerifiableCredentialFactory,
	compose.RFC8628DeviceFactory,
	compose.RFC8628DeviceAuthorizationTokenFactory,
	compose.OpenIDConnectDeviceFactory,
}

func NewConfig(deps configDependencies) *Config {
	return &Config{
		deps:            deps,
		DefaultProvider: deps.Config(),
	}
}

func (c *Config) LoadDefaultHandlers(strategy interface{}) {
	factories := append(defaultFactories, c.deps.ExtraFositeFactories()...)
	for _, factory := range factories {
		res := factory(c, c.deps.Persister(), strategy)
		if ah, ok := res.(fosite.AuthorizeEndpointHandler); ok {
			c.authorizeEndpointHandlers.Append(ah)
		}
		if th, ok := res.(fosite.TokenEndpointHandler); ok {
			c.tokenEndpointHandlers.Append(th)
		}
		if tv, ok := res.(fosite.TokenIntrospector); ok {
			c.tokenIntrospectionHandlers.Append(tv)
		}
		if rh, ok := res.(fosite.RevocationHandler); ok {
			c.revocationHandlers.Append(rh)
		}
		if dh, ok := res.(fosite.DeviceEndpointHandler); ok {
			c.deviceEndpointHandlers.Append(dh)
		}
	}
}

func (c *Config) GetJWKSFetcherStrategy(context.Context) fosite.JWKSFetcherStrategy {
	if c.jwksFetcherStrategy == nil {
		c.jwksFetcherStrategy = fosite.NewDefaultJWKSFetcherStrategy(fosite.JWKSFetcherWithHTTPClientSource(
			func(ctx context.Context) *retryablehttp.Client { return c.deps.HTTPClient(ctx) },
		))
	}
	return c.jwksFetcherStrategy
}

func (c *Config) GetHTTPClient(ctx context.Context) *retryablehttp.Client {
	return c.deps.HTTPClient(ctx)
}

func (c *Config) GetAuthorizeEndpointHandlers(context.Context) fosite.AuthorizeEndpointHandlers {
	return c.authorizeEndpointHandlers
}

func (c *Config) GetTokenEndpointHandlers(context.Context) fosite.TokenEndpointHandlers {
	return c.tokenEndpointHandlers
}

func (c *Config) GetTokenIntrospectionHandlers(context.Context) (r fosite.TokenIntrospectionHandlers) {
	return c.tokenIntrospectionHandlers
}

func (c *Config) GetRevocationHandlers(context.Context) fosite.RevocationHandlers {
	return c.revocationHandlers
}

// GetDeviceEndpointHandlers returns the deviceEndpointHandlers
func (c *Config) GetDeviceEndpointHandlers(context.Context) fosite.DeviceEndpointHandlers {
	return c.deviceEndpointHandlers
}

func (c *Config) GetGrantTypeJWTBearerCanSkipClientAuth(context.Context) bool {
	return false
}

func (c *Config) GetAudienceStrategy(context.Context) fosite.AudienceMatchingStrategy {
	return fosite.DefaultAudienceMatchingStrategy
}

func (c *Config) GetOmitRedirectScopeParam(context.Context) bool {
	return false
}

func (c *Config) GetSanitationWhiteList(context.Context) []string {
	return []string{"code", "redirect_uri"}
}

func (c *Config) GetEnablePKCEPlainChallengeMethod(context.Context) bool {
	return false
}

func (c *Config) GetDisableRefreshTokenValidation(context.Context) bool {
	return false
}

func (c *Config) GetRefreshTokenScopes(context.Context) []string {
	return []string{"offline", "offline_access"}
}

func (c *Config) GetMinParameterEntropy(_ context.Context) int {
	return fosite.MinParameterEntropy
}

func (c *Config) GetClientAuthenticationStrategy(context.Context) fosite.ClientAuthenticationStrategy {
	// Fosite falls back to the default fosite.Fosite.DefaultClientAuthenticationStrategy when this is nil.
	return nil
}

func (c *Config) GetResponseModeHandlerExtension(context.Context) fosite.ResponseModeHandler {
	return defaultResponseModeHandler
}

func (c *Config) GetSendDebugMessagesToClients(ctx context.Context) bool {
	return c.deps.Config().GetSendDebugMessagesToClients(ctx)
}

func (c *Config) GetMessageCatalog(context.Context) i18n.MessageCatalog {
	// Fosite falls back to the default messages when this is nil.
	return nil
}

func (c *Config) GetSecretsHasher(context.Context) fosite.Hasher {
	return c.deps.ClientHasher()
}

func (c *Config) GetTokenEntropy(context.Context) int {
	return 32
}

func (c *Config) GetHMACHasher(context.Context) func() hash.Hash {
	return sha512.New512_256
}

func (c *Config) GetIDTokenIssuer(ctx context.Context) string {
	return c.deps.Config().IssuerURL(ctx).String()
}

func (c *Config) GetAllowedPrompts(context.Context) []string {
	return []string{"login", "none", "consent", "registration"}
}

func (c *Config) GetRedirectSecureChecker(context.Context) func(context.Context, *url.URL) bool {
	return x.IsRedirectURISecure(c.deps.Config())
}

func (c *Config) GetAccessTokenIssuer(ctx context.Context) string {
	return c.deps.Config().IssuerURL(ctx).String()
}

func (c *Config) GetJWTScopeField(ctx context.Context) jwt.JWTScopeFieldEnum {
	return c.deps.Config().GetJWTScopeField(ctx)
}

func (c *Config) GetFormPostHTMLTemplate(context.Context) *template.Template {
	return fosite.DefaultFormPostTemplate
}

func (c *Config) GetTokenURLs(ctx context.Context) []string {
	return stringslice.Unique([]string{
		c.OAuth2TokenURL(ctx).String(),
		urlx.AppendPaths(c.deps.Config().PublicURL(ctx), oauth2.TokenPath).String(),
	})
}

// GetDeviceVerificationURL returns the device verification url
func (c *Config) GetDeviceVerificationURL(ctx context.Context) string {
	return urlx.AppendPaths(c.deps.Config().PublicURL(ctx), oauth2.DeviceVerificationPath).String()
}
