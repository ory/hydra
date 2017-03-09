package oauth2

import (
	"golang.org/x/net/context"
)

// TokenRevocationStorage provides the storage implementation
// as specified in: https://tools.ietf.org/html/rfc7009
type TokenRevocationStorage interface {
	RefreshTokenStorage
	AccessTokenStorage

	// RevokeRefreshToken revokes a refresh token as specified in:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	// If the particular
	// token is a refresh token and the authorization server supports the
	// revocation of access tokens, then the authorization server SHOULD
	// also invalidate all access tokens based on the same authorization
	// grant (see Implementation Note).
	RevokeRefreshToken(ctx context.Context, requestID string) error

	// RevokeAccessToken revokes an access token as specified in:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	// If the token passed to the request
	// is an access token, the server MAY revoke the respective refresh
	// token as well.
	RevokeAccessToken(ctx context.Context, requestID string) error
}
