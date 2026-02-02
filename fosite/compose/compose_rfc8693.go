// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/rfc8693"
)

// RFC8693TokenExchangeFactory creates an OAuth2 Token Exchange (RFC 8693) handler.
func RFC8693TokenExchangeFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &rfc8693.Handler{
		Strategy: strategy.(oauth2.AccessTokenStrategyProvider),
		Storage:  storage.(oauth2.AccessTokenStorageProvider),
		Config:   config,
	}
}
