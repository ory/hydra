package oauth2

import "context"

type Revocator interface {
	RevokeToken(ctx context.Context, token string) error
}
