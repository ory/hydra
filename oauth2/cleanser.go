package oauth2

import "context"

type Cleanser interface {
	CleanseTokens(ctx context.Context) error
}
