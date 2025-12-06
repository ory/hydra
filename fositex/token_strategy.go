// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fositex

import (
	"context"
	"strings"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	foauth2 "github.com/ory/hydra/v2/fosite/handler/oauth2"
)

var _ foauth2.CoreStrategy = (*TokenStrategy)(nil)

type (
	// TokenStrategy uses the correct token strategy (jwt, opaque) depending on the configuration.
	TokenStrategy struct {
		d tokenStrategyDependencies
	}
	tokenStrategyDependencies interface {
		OAuth2HMACStrategy() foauth2.CoreStrategy
		OAuth2JWTStrategy() foauth2.AccessTokenStrategy
		config.Provider
	}
)

// NewTokenStrategy returns a new TokenStrategy.
func NewTokenStrategy(d tokenStrategyDependencies) *TokenStrategy { return &TokenStrategy{d: d} }

// gs returns the configured strategy.
func (t TokenStrategy) gs(ctx context.Context, additionalSources ...config.AccessTokenStrategySource) foauth2.AccessTokenStrategy {
	switch ats := t.d.Config().AccessTokenStrategy(ctx, additionalSources...); ats {
	case config.AccessTokenJWTStrategy:
		return t.d.OAuth2JWTStrategy()
	}
	return t.d.OAuth2HMACStrategy()
}

func (t TokenStrategy) AccessTokenSignature(_ context.Context, token string) string {
	return genericSignature(token)
}

func (t TokenStrategy) GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token, signature string, err error) {
	return t.gs(ctx, withRequester(requester)).GenerateAccessToken(ctx, requester)
}

func (t TokenStrategy) ValidateAccessToken(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.gs(ctx, withRequester(requester)).ValidateAccessToken(ctx, requester, token)
}

func (t TokenStrategy) RefreshTokenSignature(ctx context.Context, token string) string {
	return t.d.OAuth2HMACStrategy().RefreshTokenSignature(ctx, token)
}

func (t TokenStrategy) GenerateRefreshToken(ctx context.Context, requester fosite.Requester) (token, signature string, err error) {
	return t.d.OAuth2HMACStrategy().GenerateRefreshToken(ctx, requester)
}

func (t TokenStrategy) ValidateRefreshToken(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.d.OAuth2HMACStrategy().ValidateRefreshToken(ctx, requester, token)
}

func (t TokenStrategy) AuthorizeCodeSignature(ctx context.Context, token string) string {
	return t.d.OAuth2HMACStrategy().AuthorizeCodeSignature(ctx, token)
}

func (t TokenStrategy) GenerateAuthorizeCode(ctx context.Context, requester fosite.Requester) (token, signature string, err error) {
	return t.d.OAuth2HMACStrategy().GenerateAuthorizeCode(ctx, requester)
}

func (t TokenStrategy) ValidateAuthorizeCode(ctx context.Context, requester fosite.Requester, token string) (err error) {
	return t.d.OAuth2HMACStrategy().ValidateAuthorizeCode(ctx, requester, token)
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
