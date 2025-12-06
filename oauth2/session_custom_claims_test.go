// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestCustomClaimsInSession(t *testing.T) {
	t.Parallel()

	session := Session{DefaultSession: &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Subject: "alice",
			Issuer:  "hydra.localhost",
		},
		Headers: new(jwt.Headers),
		Subject: "alice",
	}}

	for _, tc := range []struct {
		name                                        string
		extra, expectedClaims                       map[string]any
		allowedTopLevelClaims, expectNotSet         []string
		mirrorTopLevelClaims, excludeNotBeforeClaim bool
	}{{
		name:  "no custom claims",
		extra: map[string]any{},
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
		},
		expectNotSet: []string{"ext"},
	}, {
		name:                  "top level mirrored",
		extra:                 map[string]any{"foo": "bar"},
		allowedTopLevelClaims: []string{"foo"},
		mirrorTopLevelClaims:  true,
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
			"foo": "bar",
			"ext": map[string]any{"foo": "bar"},
		},
	}, {
		name: "top level mirrored with reserved",
		extra: map[string]any{
			"foo": "bar",
			"iss": "hydra.remote",
			"sub": "another-alice",
		},
		allowedTopLevelClaims: []string{"foo", "iss", "sub"},
		mirrorTopLevelClaims:  true,
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
			"foo": "bar",
			"ext": map[string]any{
				"foo": "bar",
				"iss": "hydra.remote",
				"sub": "another-alice",
			},
		},
	}, {
		name: "with disallowed top level mirrored",
		extra: map[string]any{
			"foo": "bar",
			"baz": "qux",
		},
		allowedTopLevelClaims: []string{"foo"},
		mirrorTopLevelClaims:  true,
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
			"foo": "bar",
			"ext": map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
		},
		expectNotSet: []string{"baz"},
	}, {
		name:                  "mirrored top level claims with other keys",
		extra:                 map[string]any{"foo": "bar"},
		allowedTopLevelClaims: []string{"foo", "bar"},
		mirrorTopLevelClaims:  true,
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
			"foo": "bar",
			"ext": map[string]any{"foo": "bar"},
		},
		expectNotSet: []string{"bar"},
	}, {
		name:                  "disabled mirror top level claims",
		extra:                 map[string]any{"foo": "bar"},
		allowedTopLevelClaims: []string{"foo"},
		mirrorTopLevelClaims:  false,
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
			"foo": "bar",
		},
		expectNotSet: []string{"ext"},
	}, {
		name:                  "exclude not before claim",
		extra:                 map[string]any{},
		excludeNotBeforeClaim: true,
		expectedClaims: map[string]any{
			"sub": "alice",
			"iss": "hydra.localhost",
		},
		expectNotSet: []string{"nbf"},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			sess := session
			sess.Extra = tc.extra
			sess.AllowedTopLevelClaims = tc.allowedTopLevelClaims
			sess.MirrorTopLevelClaims = tc.mirrorTopLevelClaims
			sess.ExcludeNotBeforeClaim = tc.excludeNotBeforeClaim

			claims := sess.GetJWTClaims().ToMapClaims()
			assert.Subset(t, claims, tc.expectedClaims)
			for _, key := range tc.expectNotSet {
				assert.NotContains(t, claims, key)
			}
		})
	}
}
