// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fositex

import (
	"context"
	"testing"

	"github.com/ory/hydra/v2/fosite/token/hmac"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
)

// Test that the generic signature function implements the same signature as the
// HMAC and JWT strategies.
func TestAccessTokenSignature(t *testing.T) {
	ctx := context.Background()

	t.Run("strategy=DefaultJWTStrategy", func(t *testing.T) {
		strategy := new(oauth2.DefaultJWTStrategy)
		for _, tc := range []struct{ token string }{
			{""},
			{"foo"},
			// tokens with two parts will be handled by the HMAC strategy
			{"foo.bar.baz"},
			{"foo.bar.baz.qux"},
		} {
			t.Run("case="+tc.token, func(t *testing.T) {
				assert.Equal(t,
					strategy.AccessTokenSignature(ctx, tc.token),
					genericSignature(tc.token))
			})
		}
	})
	t.Run("strategy=HMACStrategy", func(t *testing.T) {
		strategy := oauth2.NewHMACSHAStrategy(&hmac.HMACStrategy{}, nil)
		for _, tc := range []struct{ token string }{
			{""},
			{"foo"},
			{"foo.bar"},
			// tokens with three parts will be handled by the JWT strategy
			{"foo.bar.baz.qux"},
		} {
			t.Run("case="+tc.token, func(t *testing.T) {
				assert.Equal(t,
					strategy.AccessTokenSignature(ctx, tc.token),
					genericSignature(tc.token))
			})
		}
	})
}
