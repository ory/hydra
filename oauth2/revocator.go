package oauth2

import "golang.org/x/net/context"

type Revocator interface {
	RevokeToken(ctx context.Context, token string) error
}
