// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/token/hmac"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

type CommonStrategy struct {
	oauth2.CoreStrategy
	rfc8628.RFC8628CodeStrategy
	openid.OpenIDConnectTokenStrategy
	jwt.Signer
}

type HMACSHAStrategyConfigurator interface {
	fosite.AccessTokenLifespanProvider
	fosite.RefreshTokenLifespanProvider
	fosite.AuthorizeCodeLifespanProvider
	fosite.TokenEntropyProvider
	fosite.GlobalSecretProvider
	fosite.RotatedGlobalSecretsProvider
	fosite.HMACHashingProvider
	fosite.DeviceAndUserCodeLifespanProvider
}

func NewOAuth2HMACStrategy(config HMACSHAStrategyConfigurator) *oauth2.HMACSHAStrategy {
	return oauth2.NewHMACSHAStrategy(&hmac.HMACStrategy{Config: config}, config)
}

func NewOAuth2JWTStrategy(keyGetter func(context.Context) (interface{}, error), strategy oauth2.CoreStrategy, config fosite.Configurator) *oauth2.DefaultJWTStrategy {
	return &oauth2.DefaultJWTStrategy{
		Signer:          &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		HMACSHAStrategy: strategy,
		Config:          config,
	}
}

func NewOpenIDConnectStrategy(keyGetter func(context.Context) (interface{}, error), config fosite.Configurator) *openid.DefaultStrategy {
	return &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		Config: config,
	}
}

// Create a new device strategy
func NewDeviceStrategy(config fosite.Configurator) *rfc8628.DefaultDeviceStrategy {
	return &rfc8628.DefaultDeviceStrategy{
		Enigma: &hmac.HMACStrategy{Config: config},
		Config: config,
	}
}
