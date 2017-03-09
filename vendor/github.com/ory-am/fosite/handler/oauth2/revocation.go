package oauth2

import (
	"github.com/ory-am/fosite"
	"golang.org/x/net/context"
)

type TokenRevocationHandler struct {
	TokenRevocationStorage TokenRevocationStorage
	RefreshTokenStrategy   RefreshTokenStrategy
	AccessTokenStrategy    AccessTokenStrategy
}

// RevokeToken implements https://tools.ietf.org/html/rfc7009#section-2.1
// The token type hint indicates which token type check should be performed first.
func (r *TokenRevocationHandler) RevokeToken(ctx context.Context, token string, tokenType fosite.TokenType) error {
	discoveryFuncs := []func() (request fosite.Requester, err error){
		func() (request fosite.Requester, err error) {
			// Refresh token
			signature := r.RefreshTokenStrategy.RefreshTokenSignature(token)
			return r.TokenRevocationStorage.GetRefreshTokenSession(ctx, signature, nil)
		},
		func() (request fosite.Requester, err error) {
			// Access token
			signature := r.AccessTokenStrategy.AccessTokenSignature(token)
			return r.TokenRevocationStorage.GetAccessTokenSession(ctx, signature, nil)
		},
	}

	// Token type hinting
	if tokenType == fosite.AccessToken {
		discoveryFuncs[0], discoveryFuncs[1] = discoveryFuncs[1], discoveryFuncs[0]
	}

	var ar fosite.Requester
	var err error
	if ar, err = discoveryFuncs[0](); err != nil {
		ar, err = discoveryFuncs[1]()
	}
	if err != nil {
		return err
	}

	requestID := ar.GetID()
	r.TokenRevocationStorage.RevokeRefreshToken(ctx, requestID)
	r.TokenRevocationStorage.RevokeAccessToken(ctx, requestID)

	return nil
}
