// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"

	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var header = &Headers{
	Extra: map[string]interface{}{
		"foo": "bar",
	},
}

func TestHash(t *testing.T) {
	for k, tc := range []struct {
		d        string
		strategy Signer
	}{
		{
			d: "RS256",
			strategy: &DefaultSigner{GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return gen.MustRSAKey(), nil
			}},
		},
		{
			d: "ES256",
			strategy: &DefaultSigner{GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return gen.MustES256Key(), nil
			}},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/strategy=%s", k, tc.d), func(t *testing.T) {
			in := []byte("foo")
			out, err := tc.strategy.Hash(context.TODO(), in)
			assert.NoError(t, err)
			assert.NotEqual(t, in, out)
		})
	}
}

func TestAssign(t *testing.T) {
	for k, c := range [][]map[string]interface{}{
		{
			{"foo": "bar"},
			{"baz": "bar"},
			{"foo": "bar", "baz": "bar"},
		},
		{
			{"foo": "bar"},
			{"foo": "baz"},
			{"foo": "bar"},
		},
		{
			{},
			{"foo": "baz"},
			{"foo": "baz"},
		},
		{
			{"foo": "bar"},
			{"foo": "baz", "bar": "baz"},
			{"foo": "bar", "bar": "baz"},
		},
	} {
		assert.EqualValues(t, c[2], assign(c[0], c[1]), "Case %d", k)
	}
}

func TestGenerateJWT(t *testing.T) {
	var key interface{} = gen.MustRSAKey()
	for k, tc := range []struct {
		d        string
		strategy Signer
		resetKey func(strategy Signer)
	}{
		{
			d: "DefaultSigner",
			strategy: &DefaultSigner{
				GetPrivateKey: func(_ context.Context) (interface{}, error) {
					return key, nil
				},
			},
			resetKey: func(strategy Signer) {
				key = gen.MustRSAKey()
			},
		},
		{
			d: "ES256JWTStrategy",
			strategy: &DefaultSigner{
				GetPrivateKey: func(_ context.Context) (interface{}, error) {
					return key, nil
				},
			},
			resetKey: func(strategy Signer) {
				key = &jose.JSONWebKey{
					KeyID:     "test-id",
					Key:       gen.MustES521Key(),
					Algorithm: "ES512",
				}
			},
		},
		{
			d: "ES256JWTStrategy",
			strategy: &DefaultSigner{
				GetPrivateKey: func(_ context.Context) (interface{}, error) {
					return key, nil
				},
			},
			resetKey: func(strategy Signer) {
				key = gen.MustES256Key()
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/strategy=%s", k, tc.d), func(t *testing.T) {
			claims := &JWTClaims{
				ExpiresAt: time.Now().UTC().Add(time.Hour),
			}

			token, sig, err := tc.strategy.Generate(context.TODO(), claims.ToMapClaims(), header)
			require.NoError(t, err)
			require.NotEmpty(t, token)
			require.NotEmpty(t, sig)

			decoded, err := tc.strategy.Decode(context.TODO(), token)
			require.NoError(t, err)
			require.NotNil(t, decoded)

			if k, ok := key.(*jose.JSONWebKey); ok && k.KeyID != "" {
				require.Equal(t, k.KeyID, decoded.Header["kid"])
			}

			sig, err = tc.strategy.Validate(context.TODO(), token)
			require.NoError(t, err)

			sig, err = tc.strategy.Validate(context.TODO(), token+"."+"0123456789")
			require.Error(t, err)

			partToken := strings.Split(token, ".")[2]

			sig, err = tc.strategy.Validate(context.TODO(), partToken)
			require.Error(t, err)

			// Reset private key
			tc.resetKey(tc.strategy)

			// Lets validate the exp claim
			claims = &JWTClaims{
				ExpiresAt: time.Now().UTC().Add(-time.Hour),
			}
			token, sig, err = tc.strategy.Generate(context.TODO(), claims.ToMapClaims(), header)
			require.NoError(t, err)
			require.NotNil(t, token)

			sig, err = tc.strategy.Validate(context.TODO(), token)
			require.Error(t, err)

			// Lets validate the nbf claim
			claims = &JWTClaims{
				NotBefore: time.Now().UTC().Add(time.Hour),
			}
			token, sig, err = tc.strategy.Generate(context.TODO(), claims.ToMapClaims(), header)
			require.NoError(t, err)
			require.NotNil(t, token)
			//t.Logf("%s.%s", token, sig)
			sig, err = tc.strategy.Validate(context.TODO(), token)
			require.Error(t, err)
			require.Empty(t, sig, "%s", err)
		})
	}
}

func TestValidateSignatureRejectsJWT(t *testing.T) {
	for k, tc := range []struct {
		d        string
		strategy Signer
	}{
		{
			d: "RS256",
			strategy: &DefaultSigner{GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return gen.MustRSAKey(), nil
			},
			},
		},
		{
			d: "ES256",
			strategy: &DefaultSigner{
				GetPrivateKey: func(_ context.Context) (interface{}, error) {
					return gen.MustES256Key(), nil
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/strategy=%s", k, tc.d), func(t *testing.T) {
			for k, c := range []string{
				"",
				" ",
				"foo.bar",
				"foo.",
				".foo",
			} {
				_, err := tc.strategy.Validate(context.TODO(), c)
				assert.Error(t, err)
				t.Logf("Passed test case %d", k)
			}
		})
	}
}
