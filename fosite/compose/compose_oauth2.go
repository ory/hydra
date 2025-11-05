// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

// OAuth2AuthorizeExplicitFactory creates an OAuth2 authorize code grant ("authorize explicit flow") handler and registers
// an access token, refresh token and authorize code validator.
func OAuth2AuthorizeExplicitFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.AuthorizeExplicitGrantHandler{
		Strategy: strategy.(interface {
			oauth2.AccessTokenStrategyProvider
			oauth2.RefreshTokenStrategyProvider
			oauth2.AuthorizeCodeStrategyProvider
		}),
		Storage: storage.(interface {
			oauth2.AuthorizeCodeStorageProvider
			oauth2.AccessTokenStorageProvider
			oauth2.RefreshTokenStorageProvider
			oauth2.TokenRevocationStorageProvider
		}),
		Config: config,
	}
}

// OAuth2ClientCredentialsGrantFactory creates an OAuth2 client credentials grant handler and registers
// an access token, refresh token and authorize code validator.
func OAuth2ClientCredentialsGrantFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.ClientCredentialsGrantHandler{
		Strategy: strategy.(oauth2.AccessTokenStrategyProvider),
		Storage:  storage.(oauth2.AccessTokenStorageProvider),
		Config:   config,
	}
}

// OAuth2RefreshTokenGrantFactory creates an OAuth2 refresh grant handler and registers
// an access token, refresh token and authorize code validator.nmj
func OAuth2RefreshTokenGrantFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.RefreshTokenGrantHandler{
		Strategy: strategy.(interface {
			oauth2.AccessTokenStrategyProvider
			oauth2.RefreshTokenStrategyProvider
		}),
		Storage: storage.(interface {
			oauth2.AccessTokenStorageProvider
			oauth2.RefreshTokenStorageProvider
			oauth2.TokenRevocationStorageProvider
		}),
		Config: config,
	}
}

// OAuth2AuthorizeImplicitFactory creates an OAuth2 implicit grant ("authorize implicit flow") handler and registers
// an access token, refresh token and authorize code validator.
func OAuth2AuthorizeImplicitFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.AuthorizeImplicitGrantHandler{
		Strategy: strategy.(oauth2.AccessTokenStrategyProvider),
		Storage:  storage.(oauth2.AccessTokenStorageProvider),
		Config:   config,
	}
}

// OAuth2ResourceOwnerPasswordCredentialsFactory creates an OAuth2 resource owner password credentials grant handler and registers
// an access token, refresh token and authorize code validator.
//
// Deprecated: This factory is deprecated as a means to communicate that the ROPC grant type is widely discouraged and
// is at the time of this writing going to be omitted in the OAuth 2.1 spec. For more information on why this grant type
// is discouraged see: https://www.scottbrady91.com/oauth/why-the-resource-owner-password-credentials-grant-type-is-not-authentication-nor-suitable-for-modern-applications
func OAuth2ResourceOwnerPasswordCredentialsFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.ResourceOwnerPasswordCredentialsGrantHandler{
		Strategy: strategy.(interface {
			oauth2.AccessTokenStrategyProvider
			oauth2.RefreshTokenStrategyProvider
		}),
		Storage: storage.(oauth2.ResourceOwnerPasswordCredentialsGrantStorage),
		Config:  config,
	}
}

// OAuth2TokenRevocationFactory creates an OAuth2 token revocation handler.
func OAuth2TokenRevocationFactory(_ fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.TokenRevocationHandler{
		Strategy: strategy.(interface {
			oauth2.AccessTokenStrategyProvider
			oauth2.RefreshTokenStrategyProvider
		}),
		Storage: storage.(interface {
			oauth2.AccessTokenStorageProvider
			oauth2.RefreshTokenStorageProvider
			oauth2.TokenRevocationStorageProvider
		}),
	}
}

// OAuth2TokenIntrospectionFactory creates an OAuth2 token introspection handler and registers
// an access token and refresh token validator.
func OAuth2TokenIntrospectionFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.CoreValidator{
		Strategy: strategy.(interface {
			oauth2.AccessTokenStrategyProvider
			oauth2.RefreshTokenStrategyProvider
		}),
		Storage: storage.(interface {
			oauth2.AccessTokenStorageProvider
			oauth2.RefreshTokenStorageProvider
		}),
		Config: config,
	}
}

// OAuth2StatelessJWTIntrospectionFactory creates an OAuth2 token introspection handler and
// registers an access token validator. This can only be used to validate JWTs and does so
// statelessly, meaning it uses only the data available in the JWT itself, and does not access the
// storage implementation at all.
//
// Due to the stateless nature of this factory, THE BUILT-IN REVOCATION MECHANISMS WILL NOT WORK.
// If you need revocation, you can validate JWTs statefully, using the other factories.
func OAuth2StatelessJWTIntrospectionFactory(config fosite.Configurator, _ fosite.Storage, strategy interface{}) interface{} {
	return &oauth2.StatelessJWTValidator{
		Signer: strategy.(jwt.Signer),
		Config: config,
	}
}
