// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk_test

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/hydra/v2/internal/testhelpers"
	. "github.com/ory/hydra/v2/jwk"
)

func TestJWTStrategy(t *testing.T) {
	for _, alg := range []string{"RS256", "ES256", "ES512"} {
		t.Run("case="+alg, func(t *testing.T) {
			reg := testhelpers.NewRegistryMemory(t)
			m := reg.KeyManager()

			_, err := m.GenerateAndPersistKeySet(context.Background(), "foo-set", "foo", alg, "sig")
			require.NoError(t, err)

			s := NewDefaultJWTSigner(reg, "foo-set")
			a, b, err := s.Generate(context.Background(), jwt.MapClaims{"foo": "bar"}, &jwt.Headers{})
			require.NoError(t, err)
			assert.NotEmpty(t, a)
			assert.NotEmpty(t, b)

			token, err := base64.RawStdEncoding.DecodeString(strings.Split(a, ".")[0])
			require.NoError(t, err)
			assert.Equal(t, alg, gjson.GetBytes(token, "alg").String())

			_, err = s.Validate(context.Background(), a)
			require.NoError(t, err)

			kidFoo, err := s.GetPublicKeyID(context.Background())
			assert.NoError(t, err)

			_, err = m.GenerateAndPersistKeySet(context.Background(), "foo-set", "bar", alg, "sig")
			require.NoError(t, err)

			a, b, err = s.Generate(context.Background(), jwt.MapClaims{"foo": "bar"}, &jwt.Headers{})
			require.NoError(t, err)
			assert.NotEmpty(t, a)
			assert.NotEmpty(t, b)

			token, err = base64.RawStdEncoding.DecodeString(strings.Split(a, ".")[0])
			require.NoError(t, err)
			assert.Equal(t, alg, gjson.GetBytes(token, "alg").String())

			_, err = s.Validate(context.Background(), a)
			require.NoError(t, err)

			kidBar, err := s.GetPublicKeyID(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, "foo", kidFoo)
			assert.Equal(t, "bar", kidBar)
		})
	}
}
