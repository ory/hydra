// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal/gen"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestIntrospectJWT(t *testing.T) {
	rsaKey := gen.MustRSAKey()

	signer := &jwt.DefaultSigner{GetPrivateKey: func(_ context.Context) (interface{}, error) { return rsaKey, nil }}
	strat := &oauth2.DefaultJWTStrategy{
		Signer: signer,
		Config: &fosite.Config{},
	}

	v := &oauth2.StatelessJWTValidator{
		Signer: signer,
		Config: &fosite.Config{
			ScopeStrategy: fosite.HierarchicScopeStrategy,
		},
	}
	for k, c := range []struct {
		description string
		token       func(t *testing.T) string
		expectErr   error
		scopes      []string
	}{
		{
			description: "should fail because jwt is expired",
			token: func(t *testing.T) string {
				tok := jwtExpiredCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(t.Context(), tok)
				require.NoError(t, err)
				return token
			},
			expectErr: fosite.ErrTokenExpired,
		},
		{
			description: "should pass because scope was granted",
			token: func(t *testing.T) string {
				tok := jwtValidCase(fosite.AccessToken)
				tok.GrantedScope = []string{"foo", "bar"}
				token, _, err := strat.GenerateAccessToken(t.Context(), tok)
				require.NoError(t, err)
				return token
			},
			scopes: []string{"foo"},
		},
		{
			description: "should fail because scope was not granted",
			token: func(t *testing.T) string {
				tok := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(t.Context(), tok)
				require.NoError(t, err)
				return token
			},
			scopes:    []string{"foo"},
			expectErr: fosite.ErrInvalidScope,
		},
		{
			description: "should fail because signature is invalid",
			token: func(t *testing.T) string {
				tok := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(t.Context(), tok)
				require.NoError(t, err)
				parts := strings.Split(token, ".")
				require.Len(t, parts, 3, "%s - %v", token, parts)
				dec, err := base64.RawURLEncoding.DecodeString(parts[1])
				require.NoError(t, err)
				s := strings.ReplaceAll(string(dec), "peter", "piper")
				parts[1] = base64.RawURLEncoding.EncodeToString([]byte(s))
				return strings.Join(parts, ".")
			},
			expectErr: fosite.ErrTokenSignatureMismatch,
		},
		{
			description: "should pass",
			token: func(t *testing.T) string {
				tok := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(t.Context(), tok)
				require.NoError(t, err)
				return token
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d:%v", k, c.description), func(t *testing.T) {
			if c.scopes == nil {
				c.scopes = []string{}
			}

			areq := fosite.NewAccessRequest(nil)
			_, err := v.IntrospectToken(t.Context(), c.token(t), fosite.AccessToken, areq, c.scopes)

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, "peter", areq.Session.GetSubject())
			}
		})
	}
}

func BenchmarkIntrospectJWT(b *testing.B) {
	key := gen.MustRSAKey()
	signer := &jwt.DefaultSigner{GetPrivateKey: func(_ context.Context) (interface{}, error) { return key, nil }}
	strat := &oauth2.DefaultJWTStrategy{Signer: signer, Config: &fosite.Config{}}

	v := &oauth2.StatelessJWTValidator{Signer: signer}

	tok := jwtValidCase(fosite.AccessToken)
	token, _, err := strat.GenerateAccessToken(b.Context(), tok)
	assert.NoError(b, err)
	areq := fosite.NewAccessRequest(nil)

	for n := 0; n < b.N; n++ {
		_, err = v.IntrospectToken(b.Context(), token, fosite.AccessToken, areq, []string{})
	}

	assert.NoError(b, err)
}
