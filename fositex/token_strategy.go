// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fositex

import (
	"context"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/hydra/driver/config"
)

var _ foauth2.CoreStrategy = (*TokenStrategy)(nil)

// TokenStrategy uses the correct token strategy (jwt, opaque) depending on the configuration.
type TokenStrategy struct {
	c    *config.DefaultProvider
	hmac *foauth2.HMACSHAStrategy
	jwt  *foauth2.DefaultJWTStrategy
}

// NewTokenStrategy returns a new TokenStrategy.
func NewTokenStrategy(c *config.DefaultProvider, hmac *foauth2.HMACSHAStrategy, jwt *foauth2.DefaultJWTStrategy) *TokenStrategy {
	return &TokenStrategy{c: c, hmac: hmac, jwt: jwt}
}

// gs returns the configured strategy.
func (t TokenStrategy) gs(ctx context.Context) foauth2.CoreStrategy {
	switch ats := t.c.AccessTokenStrategy(ctx); ats {
	case config.AccessTokenJWTStrategy:
		return t.jwt
	}
	return t.hmac
}

func (t TokenStrategy) AccessTokenSignature(ctx context.Context, token string) string {
	return t.gs(ctx).AccessTokenSignature(ctx, token)
}

func (t TokenStrategy) GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return t.gs(ctx).GenerateAccessToken(ctx, requester)
}

func (t TokenStrategy) ValidateAccessToken(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx).ValidateAccessToken(ctx, requester, token)
}

func (t TokenStrategy) RefreshTokenSignature(ctx context.Context, token string) string {
	return t.gs(ctx).RefreshTokenSignature(ctx, token)
}

func (t TokenStrategy) GenerateRefreshToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return t.gs(ctx).GenerateRefreshToken(ctx, requester)
}

func (t TokenStrategy) ValidateRefreshToken(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx).ValidateRefreshToken(ctx, requester, token)
}

func (t TokenStrategy) AuthorizeCodeSignature(ctx context.Context, token string) string {
	return t.gs(ctx).AuthorizeCodeSignature(ctx, token)
}

func (t TokenStrategy) GenerateAuthorizeCode(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return t.gs(ctx).GenerateAuthorizeCode(ctx, requester)
}

func (t TokenStrategy) ValidateAuthorizeCode(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx).ValidateAuthorizeCode(ctx, requester, token)
}
