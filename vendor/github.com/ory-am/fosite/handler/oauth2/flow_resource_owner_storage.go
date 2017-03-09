package oauth2

import (
	"golang.org/x/net/context"
)

type ResourceOwnerPasswordCredentialsGrantStorage interface {
	Authenticate(ctx context.Context, name string, secret string) error
	AccessTokenStorage
	RefreshTokenStorage
}
