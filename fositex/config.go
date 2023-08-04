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
	"github.com/ory/x/urlx"
)

type configDependencies interface {
	config.Provider
	persistence.Provider
	x.HTTPClientProvider
	GetJWKSFetcherStrategy() fosite.JWKSFetcherStrategy
	ClientHasher() fosite.Hasher
}

type factory func(config fosite.Configurator, storage interface{}, strategy interface{}) interface{}

type Config struct {
	deps configDependencies

	authorizeEndpointHandlers  fosite.AuthorizeEndpointHandlers
	tokenEndpointHandlers      fosite.TokenEndpointHandlers
	tokenIntrospectionHandlers fosite.TokenIntrospectionHandlers
	revocationHandlers         fosite.RevocationHandlers

	*config.DefaultProvider
}

var defaultResponseModeHandler = fosite.NewDefaultResponseModeHandler()
var defaultFactories = []factory{
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
}

func NewConfig(deps configDependencies) *Config {
	c := &Config{
		deps:            deps,
		DefaultProvider: deps.Config(),
	}
	return c
}

func (c *Config) LoadDefaultHandlers(strategy interface{}) {
	for _, factory := range defaultFactories {
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
	}
}

func (c *Config) GetJWKSFetcherStrategy(ctx context.Context) fosite.JWKSFetcherStrategy {
	return c.deps.GetJWKSFetcherStrategy()
}

func (c *Config) GetHTTPClient(ctx context.Context) *retryablehttp.Client {
	return c.deps.HTTPClient(ctx)
}

func (c *Config) GetAuthorizeEndpointHandlers(ctx context.Context) fosite.AuthorizeEndpointHandlers {
	return c.authorizeEndpointHandlers
}

func (c *Config) GetTokenEndpointHandlers(ctx context.Context) fosite.TokenEndpointHandlers {
	return c.tokenEndpointHandlers
}

func (c *Config) GetTokenIntrospectionHandlers(ctx context.Context) (r fosite.TokenIntrospectionHandlers) {
	return c.tokenIntrospectionHandlers
}

func (c *Config) GetRevocationHandlers(ctx context.Context) fosite.RevocationHandlers {
	return c.revocationHandlers
}

func (c *Config) GetGrantTypeJWTBearerCanSkipClientAuth(ctx context.Context) bool {
	return false
}

func (c *Config) GetAudienceStrategy(ctx context.Context) fosite.AudienceMatchingStrategy {
	return fosite.DefaultAudienceMatchingStrategy
}

func (c *Config) GetOmitRedirectScopeParam(ctx context.Context) bool {
	return false
}

func (c *Config) GetSanitationWhiteList(ctx context.Context) []string {
	return []string{"code", "redirect_uri"}
}

func (c *Config) GetEnablePKCEPlainChallengeMethod(ctx context.Context) bool {
	return false
}

func (c *Config) GetDisableRefreshTokenValidation(ctx context.Context) bool {
	return false
}

func (c *Config) GetRefreshTokenScopes(ctx context.Context) []string {
	return []string{"offline", "offline_access"}
}

func (c *Config) GetMinParameterEntropy(_ context.Context) int {
	return fosite.MinParameterEntropy
}

func (c *Config) GetClientAuthenticationStrategy(ctx context.Context) fosite.ClientAuthenticationStrategy {
	// Fosite falls back to the default fosite.Fosite.DefaultClientAuthenticationStrategy when this is nil.
	return nil
}

func (c *Config) GetResponseModeHandlerExtension(ctx context.Context) fosite.ResponseModeHandler {
	return defaultResponseModeHandler
}

func (c *Config) GetSendDebugMessagesToClients(ctx context.Context) bool {
	return c.deps.Config().GetSendDebugMessagesToClients(ctx)
}

func (c *Config) GetMessageCatalog(ctx context.Context) i18n.MessageCatalog {
	// Fosite falls back to the default messages when this is nil.
	return nil
}

func (c *Config) GetSecretsHasher(ctx context.Context) fosite.Hasher {
	return c.deps.ClientHasher()
}

func (c *Config) GetTokenEntropy(ctx context.Context) int {
	return 32
}

func (c *Config) GetHMACHasher(ctx context.Context) func() hash.Hash {
	return sha512.New512_256
}

func (c *Config) GetIDTokenIssuer(ctx context.Context) string {
	return c.deps.Config().IssuerURL(ctx).String()
}

func (c *Config) GetAllowedPrompts(ctx context.Context) []string {
	return []string{"login", "none", "consent"}
}

func (c *Config) GetRedirectSecureChecker(ctx context.Context) func(context.Context, *url.URL) bool {
	return x.IsRedirectURISecure(c.deps.Config())
}

func (c *Config) GetAccessTokenIssuer(ctx context.Context) string {
	return c.deps.Config().IssuerURL(ctx).String()
}

func (c *Config) GetJWTScopeField(ctx context.Context) jwt.JWTScopeFieldEnum {
	return c.deps.Config().GetJWTScopeField(ctx)
}

func (c *Config) GetFormPostHTMLTemplate(ctx context.Context) *template.Template {
	return fosite.DefaultFormPostTemplate
}

func (c *Config) GetTokenURL(ctx context.Context) string {
	return urlx.AppendPaths(c.deps.Config().PublicURL(ctx), oauth2.TokenPath).String()
}
