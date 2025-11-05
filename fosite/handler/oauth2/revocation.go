// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
)

var _ fosite.RevocationHandler = (*TokenRevocationHandler)(nil)

type TokenRevocationHandler struct {
	Storage interface {
		TokenRevocationStorageProvider
		AccessTokenStorageProvider
		RefreshTokenStorageProvider
	}
	Strategy interface {
		AccessTokenStrategyProvider
		RefreshTokenStrategyProvider
	}
}

// RevokeToken implements https://tools.ietf.org/html/rfc7009#section-2.1
// The token type hint indicates which token type check should be performed first.
func (r *TokenRevocationHandler) RevokeToken(ctx context.Context, token string, tokenType fosite.TokenType, client fosite.Client) error {
	discoveryFuncs := []func() (request fosite.Requester, err error){
		func() (request fosite.Requester, err error) {
			// Refresh token
			signature := r.Strategy.RefreshTokenStrategy().RefreshTokenSignature(ctx, token)
			return r.Storage.RefreshTokenStorage().GetRefreshTokenSession(ctx, signature, nil)
		},
		func() (request fosite.Requester, err error) {
			// Access token
			signature := r.Strategy.AccessTokenStrategy().AccessTokenSignature(ctx, token)
			return r.Storage.AccessTokenStorage().GetAccessTokenSession(ctx, signature, nil)
		},
	}

	// Token type hinting
	if tokenType == fosite.AccessToken {
		discoveryFuncs[0], discoveryFuncs[1] = discoveryFuncs[1], discoveryFuncs[0]
	}

	var ar fosite.Requester
	var err1, err2 error
	if ar, err1 = discoveryFuncs[0](); err1 != nil {
		ar, err2 = discoveryFuncs[1]()
	}
	// err2 can only be not nil if first err1 was not nil
	if err2 != nil {
		return storeErrorsToRevocationError(err1, err2)
	}

	if ar.GetClient().GetID() != client.GetID() {
		return errorsx.WithStack(fosite.ErrUnauthorizedClient)
	}

	requestID := ar.GetID()
	err1 = r.Storage.TokenRevocationStorage().RevokeRefreshToken(ctx, requestID)
	err2 = r.Storage.TokenRevocationStorage().RevokeAccessToken(ctx, requestID)

	return storeErrorsToRevocationError(err1, err2)
}

func storeErrorsToRevocationError(err1, err2 error) error {
	// both errors are fosite.ErrNotFound and fosite.ErrInactiveToken or nil <=> the token is revoked
	if (errors.Is(err1, fosite.ErrNotFound) || errors.Is(err1, fosite.ErrInactiveToken) || err1 == nil) &&
		(errors.Is(err2, fosite.ErrNotFound) || errors.Is(err2, fosite.ErrInactiveToken) || err2 == nil) {
		return nil
	}

	// there was an unexpected error => the token may still exist and the client should retry later
	return errorsx.WithStack(fosite.ErrTemporarilyUnavailable)
}
