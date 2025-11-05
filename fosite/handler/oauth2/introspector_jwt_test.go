// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestIntrospectJWT(t *testing.T) {
	rsaKey := gen.MustRSAKey()
	strat := &oauth2.DefaultJWTStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return rsaKey, nil
			},
		},
		Config: &fosite.Config{},
	}

	v := &oauth2.StatelessJWTValidator{
		Signer: strat,
		Config: &fosite.Config{
			ScopeStrategy: fosite.HierarchicScopeStrategy,
		},
	}
	for k, c := range []struct {
		description string
		token       func() string
		expectErr   error
		scopes      []string
	}{
		{
			description: "should fail because jwt is expired",
			token: func() string {
				jwt := jwtExpiredCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(context.Background(), jwt)
				assert.NoError(t, err)
				return token
			},
			expectErr: fosite.ErrTokenExpired,
		},
		{
			description: "should pass because scope was granted",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				jwt.GrantedScope = []string{"foo", "bar"}
				token, _, err := strat.GenerateAccessToken(context.Background(), jwt)
				assert.NoError(t, err)
				return token
			},
			scopes: []string{"foo"},
		},
		{
			description: "should fail because scope was not granted",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(context.Background(), jwt)
				assert.NoError(t, err)
				return token
			},
			scopes:    []string{"foo"},
			expectErr: fosite.ErrInvalidScope,
		},
		{
			description: "should fail because signature is invalid",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(context.Background(), jwt)
				assert.NoError(t, err)
				parts := strings.Split(token, ".")
				require.Len(t, parts, 3, "%s - %v", token, parts)
				dec, err := base64.RawURLEncoding.DecodeString(parts[1])
				assert.NoError(t, err)
				s := strings.Replace(string(dec), "peter", "piper", -1)
				parts[1] = base64.RawURLEncoding.EncodeToString([]byte(s))
				return strings.Join(parts, ".")
			},
			expectErr: fosite.ErrTokenSignatureMismatch,
		},
		{
			description: "should pass",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(context.Background(), jwt)
				assert.NoError(t, err)
				return token
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d:%v", k, c.description), func(t *testing.T) {
			if c.scopes == nil {
				c.scopes = []string{}
			}

			areq := fosite.NewAccessRequest(nil)
			_, err := v.IntrospectToken(context.Background(), c.token(), fosite.AccessToken, areq, c.scopes)

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
	strat := &oauth2.DefaultJWTStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return gen.MustRSAKey(), nil
			},
		},
		Config: &fosite.Config{},
	}

	v := &oauth2.StatelessJWTValidator{
		Signer: strat,
	}

	jwt := jwtValidCase(fosite.AccessToken)
	token, _, err := strat.GenerateAccessToken(context.Background(), jwt)
	assert.NoError(b, err)
	areq := fosite.NewAccessRequest(nil)

	for n := 0; n < b.N; n++ {
		_, err = v.IntrospectToken(context.Background(), token, fosite.AccessToken, areq, []string{})
	}

	assert.NoError(b, err)
}
