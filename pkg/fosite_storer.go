package pkg

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
)

type FositeStorer interface {
	fosite.Storage
	oauth2.CoreStorage
	openid.OpenIDConnectRequestStorage

	RevokeRefreshToken(ctx context.Context, requestID string) error

	RevokeAccessToken(ctx context.Context, requestID string) error
}
