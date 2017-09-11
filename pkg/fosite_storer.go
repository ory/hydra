package pkg

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
)

type FositeStorer interface {
	oauth2.AccessTokenStorage
	fosite.Storage
	oauth2.AuthorizeCodeGrantStorage
	oauth2.RefreshTokenGrantStorage
	openid.OpenIDConnectRequestStorage

	RevokeRefreshToken(ctx context.Context, requestID string) error

	// RevokeAccessToken revokes an access token as specified in:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	// If the token passed to the request
	// is an access token, the server MAY revoke the respective refresh
	// token as well.
	RevokeAccessToken(ctx context.Context, requestID string) error
}
