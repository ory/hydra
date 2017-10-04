package oauth2

import (
	"context"
	"time"
)

type Cleanser interface {
	CleanseTokens(ctx context.Context, before time.Time) error
}
