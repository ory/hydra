// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ory/x/configx"
)

// contextKey is a value for use with context.WithValue.
type contextKey int

const (
	// contextConfig is the key for the config in the context.
	contextConfig contextKey = iota + 1
)

// ErrNoConfigInContext is returned when no config is found in the context.
var ErrNoConfigInContext = errors.New("configuration provider not found in context")

// WithConfig returns a new context with the given configuration provider.
func WithConfig(ctx context.Context, p *configx.Provider) context.Context {
	return context.WithValue(ctx, contextConfig, p)
}

// ConfigFromContext returns the configuration provider from the context or an error if no
// configuration provider is found in the context.
func ConfigFromContext(ctx context.Context) (*configx.Provider, error) {
	if p, ok := ctx.Value(contextConfig).(*configx.Provider); ok {
		return p, nil
	}
	return nil, ErrNoConfigInContext
}

// MustConfigFromContext returns the configuration provider from the context or panics if no
// configuration provider is found in the context.
func MustConfigFromContext(ctx context.Context) *configx.Provider {
	p, err := ConfigFromContext(ctx)
	if err != nil {
		panic(err)
	}
	return p
}
