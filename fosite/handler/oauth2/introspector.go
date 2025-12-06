// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/x/errorsx"
)

var _ fosite.TokenIntrospector = (*CoreValidator)(nil)

type CoreValidator struct {
	Storage interface {
		AccessTokenStorageProvider
		RefreshTokenStorageProvider
	}
	Strategy interface {
		AccessTokenStrategyProvider
		RefreshTokenStrategyProvider
	}
	Config interface {
		fosite.ScopeStrategyProvider
		fosite.DisableRefreshTokenValidationProvider
	}
}

func (c *CoreValidator) IntrospectToken(ctx context.Context, token string, tokenUse fosite.TokenUse, accessRequest fosite.AccessRequester, scopes []string) (fosite.TokenUse, error) {
	if c.Config.GetDisableRefreshTokenValidation(ctx) {
		if err := c.introspectAccessToken(ctx, token, accessRequest, scopes); err != nil {
			return "", err
		}
		return fosite.AccessToken, nil
	}

	var err error
	switch tokenUse {
	case fosite.RefreshToken:
		if err = c.introspectRefreshToken(ctx, token, accessRequest, scopes); err == nil {
			return fosite.RefreshToken, nil
		} else if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
			return fosite.AccessToken, nil
		}
		return "", err
	}

	if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
		return fosite.AccessToken, nil
	} else if err := c.introspectRefreshToken(ctx, token, accessRequest, scopes); err == nil {
		return fosite.RefreshToken, nil
	}

	return "", err
}

func matchScopes(ss fosite.ScopeStrategy, granted, scopes []string) error {
	for _, scope := range scopes {
		if scope == "" {
			continue
		}

		if !ss(granted, scope) {
			return errorsx.WithStack(fosite.ErrInvalidScope.WithHintf("The request scope '%s' has not been granted or is not allowed to be requested.", scope))
		}
	}

	return nil
}

func (c *CoreValidator) introspectAccessToken(ctx context.Context, token string, accessRequest fosite.AccessRequester, scopes []string) error {
	sig := c.Strategy.AccessTokenStrategy().AccessTokenSignature(ctx, token)
	or, err := c.Storage.AccessTokenStorage().GetAccessTokenSession(ctx, sig, accessRequest.GetSession())
	if err != nil {
		return errorsx.WithStack(fosite.ErrRequestUnauthorized.WithWrap(err).WithDebug(err.Error()))
	} else if err := c.Strategy.AccessTokenStrategy().ValidateAccessToken(ctx, or, token); err != nil {
		return err
	}

	if err := matchScopes(c.Config.GetScopeStrategy(ctx), or.GetGrantedScopes(), scopes); err != nil {
		return err
	}

	accessRequest.Merge(or)
	return nil
}

func (c *CoreValidator) introspectRefreshToken(ctx context.Context, token string, accessRequest fosite.AccessRequester, scopes []string) error {
	sig := c.Strategy.RefreshTokenStrategy().RefreshTokenSignature(ctx, token)
	or, err := c.Storage.RefreshTokenStorage().GetRefreshTokenSession(ctx, sig, accessRequest.GetSession())

	if err != nil {
		return errorsx.WithStack(fosite.ErrRequestUnauthorized.WithWrap(err).WithDebug(err.Error()))
	} else if err := c.Strategy.RefreshTokenStrategy().ValidateRefreshToken(ctx, or, token); err != nil {
		return err
	}

	if err := matchScopes(c.Config.GetScopeStrategy(ctx), or.GetGrantedScopes(), scopes); err != nil {
		return err
	}

	accessRequest.Merge(or)
	return nil
}
