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
func RFC8628DeviceFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &rfc8628.DeviceAuthHandler{
		Strategy: strategy.(interface {
			rfc8628.DeviceRateLimitStrategyProvider
			rfc8628.DeviceCodeStrategyProvider
			rfc8628.UserCodeStrategyProvider
		}),
		Storage: storage.(interface {
			rfc8628.DeviceAuthStorageProvider
			oauth2.AccessTokenStorageProvider
			oauth2.RefreshTokenStorageProvider
		}),
		Config: config,
	}
}

// RFC8628DeviceAuthorizationTokenFactory creates an OAuth2 device authorization grant ("Device Authorization Grant") handler and registers
// an access token, refresh token and authorize code validator.
func RFC8628DeviceAuthorizationTokenFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &rfc8628.DeviceCodeTokenEndpointHandler{
		Strategy: strategy.(interface {
			rfc8628.DeviceRateLimitStrategyProvider
			rfc8628.DeviceCodeStrategyProvider
			rfc8628.UserCodeStrategyProvider
			oauth2.AccessTokenStrategyProvider
			oauth2.RefreshTokenStrategyProvider
		}),
		Storage: storage.(interface {
			rfc8628.DeviceAuthStorageProvider
			oauth2.AccessTokenStorageProvider
			oauth2.RefreshTokenStorageProvider
			oauth2.TokenRevocationStorageProvider
		}),
		Config: config,
	}
}
