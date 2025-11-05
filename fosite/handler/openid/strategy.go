// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"
	"time"

	"github.com/ory/hydra/v2/fosite"
)

type OpenIDConnectTokenStrategy interface {
	GenerateIDToken(ctx context.Context, lifespan time.Duration, requester fosite.Requester) (token string, err error)
}

type OpenIDConnectTokenStrategyProvider interface {
	OpenIDConnectTokenStrategy() OpenIDConnectTokenStrategy
}
