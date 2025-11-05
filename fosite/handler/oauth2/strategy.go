// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
)

type CoreStrategy interface {
	AuthorizeCodeStrategyProvider
	AccessTokenStrategyProvider
	RefreshTokenStrategyProvider
}

type AuthorizeCodeStrategy interface {
	AuthorizeCodeSignature(ctx context.Context, token string) string
	GenerateAuthorizeCode(ctx context.Context, requester fosite.Requester) (token string, signature string, err error)
	ValidateAuthorizeCode(ctx context.Context, requester fosite.Requester, token string) (err error)
}
type AuthorizeCodeStrategyProvider interface {
	AuthorizeCodeStrategy() AuthorizeCodeStrategy
}
type AccessTokenStrategy interface {
	AccessTokenSignature(ctx context.Context, token string) string
	GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error)
	ValidateAccessToken(ctx context.Context, requester fosite.Requester, token string) (err error)
}
type AccessTokenStrategyProvider interface {
	AccessTokenStrategy() AccessTokenStrategy
}

type RefreshTokenStrategy interface {
	RefreshTokenSignature(ctx context.Context, token string) string
	GenerateRefreshToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error)
	ValidateRefreshToken(ctx context.Context, requester fosite.Requester, token string) (err error)
}
type RefreshTokenStrategyProvider interface {
	RefreshTokenStrategy() RefreshTokenStrategy
}
