// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/rfc7523"
)

// RFC7523AssertionGrantFactory creates an OAuth2 Authorize JWT Grant (using JWTs as Authorization Grants) handler
// and registers an access token, refresh token and authorize code validator.
func RFC7523AssertionGrantFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &rfc7523.Handler{
		Strategy: strategy.(oauth2.AccessTokenStrategyProvider),
		Storage: storage.(interface {
			oauth2.AccessTokenStorageProvider
			rfc7523.RFC7523KeyStorageProvider
		}),
		Config: config,
	}
}
