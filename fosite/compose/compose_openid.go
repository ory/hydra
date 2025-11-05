// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

// OpenIDConnectExplicitFactory creates an OpenID Connect explicit ("authorize code flow") grant handler.
//
// **Important note:** You must add this handler *after* you have added an OAuth2 authorize code handler!
func OpenIDConnectExplicitFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &openid.ExplicitHandler{
		Storage: storage.(openid.OpenIDConnectRequestStorageProvider),
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategyProvider),
		},
		OpenIDConnectRequestValidator: openid.NewOpenIDConnectRequestValidator(strategy.(jwt.Signer), config),
		Config:                        config,
	}
}

// OpenIDConnectRefreshFactory creates a handler for refreshing openid connect tokens.
//
// **Important note:** You must add this handler *after* you have added an OAuth2 authorize code handler!
func OpenIDConnectRefreshFactory(config fosite.Configurator, _ fosite.Storage, strategy interface{}) interface{} {
	return &openid.OpenIDConnectRefreshHandler{
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategyProvider),
		},
		Config: config,
	}
}

// OpenIDConnectImplicitFactory creates an OpenID Connect implicit ("implicit flow") grant handler.
//
// **Important note:** You must add this handler *after* you have added an OAuth2 authorize code handler!
func OpenIDConnectImplicitFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &openid.OpenIDConnectImplicitHandler{
		AuthorizeImplicitGrantTypeHandler: &oauth2.AuthorizeImplicitGrantHandler{
			Strategy: strategy.(oauth2.AccessTokenStrategyProvider),
			Storage:  storage.(oauth2.AccessTokenStorageProvider),
			Config:   config,
		},
		Config: config,
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategyProvider),
		},
		OpenIDConnectRequestValidator: openid.NewOpenIDConnectRequestValidator(strategy.(jwt.Signer), config),
	}
}

// OpenIDConnectHybridFactory creates an OpenID Connect hybrid grant handler.
//
// **Important note:** You must add this handler *after* you have added an OAuth2 authorize code handler!
func OpenIDConnectHybridFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &openid.OpenIDConnectHybridHandler{
		AuthorizeExplicitGrantHandler: &oauth2.AuthorizeExplicitGrantHandler{
			Strategy: strategy.(interface {
				oauth2.AuthorizeCodeStrategyProvider
				oauth2.AccessTokenStrategyProvider
				oauth2.RefreshTokenStrategyProvider
			}),
			Storage: storage.(interface {
				oauth2.AuthorizeCodeStorageProvider
				oauth2.AccessTokenStorageProvider
				oauth2.RefreshTokenStorageProvider
				oauth2.TokenRevocationStorageProvider
			}),
			Config: config,
		},
		Config: config,
		AuthorizeImplicitGrantHandler: &oauth2.AuthorizeImplicitGrantHandler{
			Strategy: strategy.(oauth2.AccessTokenStrategyProvider),
			Storage:  storage.(oauth2.AccessTokenStorageProvider),
			Config:   config,
		},
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategyProvider),
		},
		OpenIDConnectRequestStorage:   storage.(openid.OpenIDConnectRequestStorageProvider),
		OpenIDConnectRequestValidator: openid.NewOpenIDConnectRequestValidator(strategy.(jwt.Signer), config),
	}
}

// OpenIDConnectDeviceFactory creates an OpenID Connect device ("device code flow") grant handler.
//
// **Important note:** You must add this handler *after* you have added an OAuth2 device authorization handler!
func OpenIDConnectDeviceFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &openid.OpenIDConnectDeviceHandler{
		Storage: storage.(openid.OpenIDConnectRequestStorageProvider),
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategyProvider),
		},
		Strategy: strategy.(rfc8628.DeviceCodeStrategyProvider),
		Config:   config,
	}
}
