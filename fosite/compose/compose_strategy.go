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

type CommonStrategyProvider struct {
	CoreStrategy      oauth2.CoreStrategy
	AccessTokenStrat  oauth2.AccessTokenStrategy
	DeviceStrategy    *rfc8628.DefaultDeviceStrategy
	OIDCTokenStrategy openid.OpenIDConnectTokenStrategy
	jwt.Signer
}

var _ oauth2.AuthorizeCodeStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) AuthorizeCodeStrategy() oauth2.AuthorizeCodeStrategy {
	return s.CoreStrategy
}

var _ oauth2.AccessTokenStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) AccessTokenStrategy() oauth2.AccessTokenStrategy {
	if s.AccessTokenStrat != nil {
		return s.AccessTokenStrat
	}
	return s.CoreStrategy
}

var _ oauth2.RefreshTokenStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) RefreshTokenStrategy() oauth2.RefreshTokenStrategy {
	return s.CoreStrategy
}

var _ openid.OpenIDConnectTokenStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) OpenIDConnectTokenStrategy() openid.OpenIDConnectTokenStrategy {
	return s.OIDCTokenStrategy
}

var _ rfc8628.DeviceRateLimitStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) DeviceRateLimitStrategy() rfc8628.DeviceRateLimitStrategy {
	return s.DeviceStrategy
}

var _ rfc8628.DeviceCodeStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) DeviceCodeStrategy() rfc8628.DeviceCodeStrategy {
	return s.DeviceStrategy
}

var _ rfc8628.UserCodeStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) UserCodeStrategy() rfc8628.UserCodeStrategy {
	return s.DeviceStrategy
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

func NewOAuth2JWTStrategy(keyGetter func(context.Context) (interface{}, error), config fosite.Configurator) *oauth2.DefaultJWTStrategy {
	return &oauth2.DefaultJWTStrategy{
		Signer: &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		Config: config,
	}
}

func NewOpenIDConnectStrategy(keyGetter func(context.Context) (interface{}, error), config fosite.Configurator) *openid.DefaultStrategy {
	return &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		Config: config,
	}
}

func NewDeviceStrategy(config fosite.Configurator) *rfc8628.DefaultDeviceStrategy {
	return &rfc8628.DefaultDeviceStrategy{
		Enigma: &hmac.HMACStrategy{Config: config},
		Config: config,
	}
}
