package pkg

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/fosite/handler/openid"
	"golang.org/x/net/context"
)

type FositeStorer interface {
	oauth2.AccessTokenStorage
	fosite.Storage
	oauth2.AuthorizeCodeGrantStorage
	oauth2.RefreshTokenGrantStorage
	oauth2.ImplicitGrantStorage
	openid.OpenIDConnectRequestStorage

	RevokeRefreshToken(ctx context.Context, requestID string) error

	// RevokeAccessToken revokes an access token as specified in:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	// If the token passed to the request
	// is an access token, the server MAY revoke the respective refresh
	// token as well.
	RevokeAccessToken(ctx context.Context, requestID string) error
}
