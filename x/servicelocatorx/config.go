package servicelocatorx

import (
	"context"
	"github.com/ory/hydra/driver/config"
)

type key int

const (
	keyConfig key = iota + 1
)

// ContextWithConfig returns a new context with the provided config.
func ContextWithConfig(ctx context.Context, c *config.DefaultProvider) context.Context {
	return context.WithValue(ctx, keyConfig, c)
}

// ConfigFromContext returns the config from the context.
func ConfigFromContext(ctx context.Context, fallback *config.DefaultProvider) *config.DefaultProvider {
	if c, ok := ctx.Value(keyConfig).(*config.DefaultProvider); ok {
		return c
	}
	return fallback
}
