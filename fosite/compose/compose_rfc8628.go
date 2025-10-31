// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package compose provides various objects which can be used to
// instantiate OAuth2Providers with different functionality.
package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
)

// RFC8628DeviceFactory creates an OAuth2 device code grant ("Device Authorization Grant") handler and registers
// a user code, device code, access token and a refresh token validator.
func RFC8628DeviceFactory(config fosite.Configurator, storage interface{}, strategy interface{}) interface{} {
	return &rfc8628.DeviceAuthHandler{
		Strategy: strategy.(rfc8628.RFC8628CodeStrategy),
		Storage:  storage.(rfc8628.RFC8628CoreStorage),
		Config:   config,
	}
}

// RFC8628DeviceAuthorizationTokenFactory creates an OAuth2 device authorization grant ("Device Authorization Grant") handler and registers
// an access token, refresh token and authorize code validator.
func RFC8628DeviceAuthorizationTokenFactory(config fosite.Configurator, storage interface{}, strategy interface{}) interface{} {
	return &rfc8628.DeviceCodeTokenEndpointHandler{
		DeviceRateLimitStrategy: strategy.(rfc8628.DeviceRateLimitStrategy),
		DeviceCodeStrategy:      strategy.(rfc8628.DeviceCodeStrategy),
		UserCodeStrategy:        strategy.(rfc8628.UserCodeStrategy),
		AccessTokenStrategy:     strategy.(oauth2.AccessTokenStrategy),
		RefreshTokenStrategy:    strategy.(oauth2.RefreshTokenStrategy),
		CoreStorage:             storage.(rfc8628.RFC8628CoreStorage),
		TokenRevocationStorage:  storage.(oauth2.TokenRevocationStorage),
		Config:                  config,
	}
}
