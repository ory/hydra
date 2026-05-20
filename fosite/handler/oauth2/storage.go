// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
)

// AuthorizeCodeStorage handles storage requests related to authorization codes.
//
// The /token-side methods receive both the original opaque authorization code
// and its strategy-derived signature. The default SQL persister keys rows by
// signature.
type AuthorizeCodeStorage interface {
	// CreateAuthorizeCodeSession stores the authorization request keyed by signature.
	CreateAuthorizeCodeSession(ctx context.Context, signature string, request fosite.Requester) (err error)

	// GetAuthorizeCodeSession hydrates the session for an authorization code and returns
	// the authorization request. code is the original opaque code string from the client;
	// signature is the strategy-derived lookup key. Implementations that key sessions in a
	// database use signature; implementations that decode session state from the code
	// itself (e.g., AEAD-encoded codes) use code.
	//
	// If the authorization code has been invalidated with `InvalidateAuthorizeCodeSession`, this
	// method should return the ErrInvalidatedAuthorizeCode error.
	//
	// Make sure to also return the fosite.Requester value when returning the fosite.ErrInvalidatedAuthorizeCode error!
	GetAuthorizeCodeSession(ctx context.Context, code, signature string, session fosite.Session) (request fosite.Requester, err error)

	// InvalidateAuthorizeCodeSession is called when an authorize code is being used. The state of the authorization
	// code should be set to invalid and consecutive requests to GetAuthorizeCodeSession should return the
	// ErrInvalidatedAuthorizeCode error.
	//
	// code is the original opaque code; signature is the lookup key. Implementations that
	// record redemption state derived from the code's contents (e.g., a replay-marker row
	// whose fields come from the AEAD payload) use code; implementations that mutate a row
	// by PK use signature.
	InvalidateAuthorizeCodeSession(ctx context.Context, code, signature string) (err error)
}
type AuthorizeCodeStorageProvider interface {
	AuthorizeCodeStorage() AuthorizeCodeStorage
}

type AccessTokenStorage interface {
	CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error)

	GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error)

	DeleteAccessTokenSession(ctx context.Context, signature string) (err error)
}

type AccessTokenStorageProvider interface {
	AccessTokenStorage() AccessTokenStorage
}

type RefreshTokenStorage interface {
	CreateRefreshTokenSession(ctx context.Context, signature string, accessSignature string, request fosite.Requester) (err error)

	GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error)

	DeleteRefreshTokenSession(ctx context.Context, signature string) (err error)

	RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) (err error)
}

type RefreshTokenStorageProvider interface {
	RefreshTokenStorage() RefreshTokenStorage
}
