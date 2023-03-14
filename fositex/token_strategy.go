// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fositex

import (
	"context"
	"strings"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	fdevice "github.com/ory/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
)

var _ foauth2.CoreStrategy = (*TokenStrategy)(nil)

// TokenStrategy uses the correct token strategy (jwt, opaque) depending on the configuration.
type TokenStrategy struct {
	c       *config.DefaultProvider
	hmac    *foauth2.HMACSHAStrategy
	devHmac *fdevice.DefaultDeviceStrategy
	jwt     *foauth2.DefaultJWTStrategy
}

// NewTokenStrategy returns a new TokenStrategy.
func NewTokenStrategy(c *config.DefaultProvider, hmac *foauth2.HMACSHAStrategy, jwt *foauth2.DefaultJWTStrategy) *TokenStrategy {
	return &TokenStrategy{c: c, hmac: hmac, jwt: jwt}
}

// gs returns the configured strategy.
func (t TokenStrategy) gs(ctx context.Context, additionalSources ...config.AccessTokenStrategySource) foauth2.CoreStrategy {
	switch ats := t.c.AccessTokenStrategy(ctx, additionalSources...); ats {
	case config.AccessTokenJWTStrategy:
		return t.jwt
	}
	return t.hmac
}

func (t TokenStrategy) AccessTokenSignature(_ context.Context, token string) string {
	return genericSignature(token)
}

func (t TokenStrategy) GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return t.gs(ctx, withRequester(requester)).GenerateAccessToken(ctx, requester)
}

func (t TokenStrategy) ValidateAccessToken(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx, withRequester(requester)).ValidateAccessToken(ctx, requester, token)
}

func (t TokenStrategy) RefreshTokenSignature(ctx context.Context, token string) string {
	return t.gs(ctx).RefreshTokenSignature(ctx, token)
}

func (t TokenStrategy) GenerateRefreshToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return t.gs(ctx, withRequester(requester)).GenerateRefreshToken(ctx, requester)
}

func (t TokenStrategy) ValidateRefreshToken(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx, withRequester(requester)).ValidateRefreshToken(ctx, requester, token)
}

func (t TokenStrategy) AuthorizeCodeSignature(ctx context.Context, token string) string {
	return t.gs(ctx).AuthorizeCodeSignature(ctx, token)
}

func (t TokenStrategy) GenerateAuthorizeCode(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return t.gs(ctx, withRequester(requester)).GenerateAuthorizeCode(ctx, requester)
}

func (t TokenStrategy) ValidateAuthorizeCode(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx, withRequester(requester)).ValidateAuthorizeCode(ctx, requester, token)
}

func withRequester(requester fosite.Requester) config.AccessTokenStrategySource {
	return client.AccessTokenStrategySource(requester.GetClient())
}

func genericSignature(token string) string {
	switch parts := strings.Split(token, "."); len(parts) {
	case 2:
		return parts[1]
	case 3:
		return parts[2]
	default:
		return ""
	}
}

func (t TokenStrategy) DeviceCodeSignature(ctx context.Context, token string) (signature string, err error) {
	return t.devHmac.DeviceCodeSignature(ctx, token)
}

func (t *TokenStrategy) GenerateDeviceCode(ctx context.Context) (token string, signature string, err error) {
	return t.devHmac.GenerateDeviceCode(ctx)
}

func (t *TokenStrategy) ValidateDeviceCode(ctx context.Context, r fosite.Requester, code string) (err error) {
	return t.devHmac.ValidateDeviceCode(ctx, r, code)
}

func (t TokenStrategy) UserCodeSignature(ctx context.Context, token string) (signature string, err error) {
	return t.devHmac.UserCodeSignature(ctx, token)
}

func (t *TokenStrategy) GenerateUserCode(ctx context.Context) (token string, signature string, err error) {
	return t.devHmac.GenerateUserCode(ctx)
}

func (t *TokenStrategy) ValidateUserCode(context context.Context, r fosite.Requester, code string) (err error) {
	return t.devHmac.ValidateUserCode(context, r, code)
}
