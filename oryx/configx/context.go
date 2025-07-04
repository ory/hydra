// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import "context"

type contextKey int

const configContextKey contextKey = iota + 1

func ContextWithConfigOptions(ctx context.Context, opts ...OptionModifier) context.Context {
	return context.WithValue(ctx, configContextKey, opts)
}

func ConfigOptionsFromContext(ctx context.Context) []OptionModifier {
	opts, ok := ctx.Value(configContextKey).([]OptionModifier)
	if !ok {
		return []OptionModifier{}
	}
	return opts
}
