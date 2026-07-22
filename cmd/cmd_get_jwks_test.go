// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
)

func TestGetJWKS(t *testing.T) {
	t.Parallel()

	c := cmd.NewGetJWKSCmd()
	reg := setup(t, c)

	set := uuid.Must(uuid.NewV4()).String()
	_ = createJWK(t, reg, set, "RS256")

	t.Run("case=gets jwks", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, set))
		assert.NotEmpty(t, actual.Get("kid").String(), actual.Raw)

		expected, err := reg.KeyManager().GetKeySet(t.Context(), set)
		require.NoError(t, err)

		assert.Equal(t, expected.Keys[0].KeyID, actual.Get("kid").String())
		assert.Equal(t, set, actual.Get("set").String(), actual.Raw)
	})

	t.Run("case=gets all keys of a set with the set name", func(t *testing.T) {
		if reg.Config().HSMEnabled() {
			t.Skip("Skipping test. When Hardware Security Module is enabled, generating a key set replaces any existing keys in the set, so a set never holds more than one key.")
		}
		multiSet := uuid.Must(uuid.NewV4()).String()
		k1 := createJWK(t, reg, multiSet, "ES256")
		k2 := createJWK(t, reg, multiSet, "ES256")

		actual := gjson.Parse(cmdx.ExecNoErr(t, c, multiSet))
		assert.Equal(t, multiSet, actual.Get("set").String(), actual.Raw)
		require.Len(t, actual.Get("keys").Array(), 2, actual.Raw)

		var kids []string
		for _, key := range actual.Get("keys").Array() {
			assert.Equal(t, multiSet, key.Get("set").String(), actual.Raw)
			kids = append(kids, key.Get("kid").String())
		}
		assert.ElementsMatch(t, []string{k1.KeyID, k2.KeyID}, kids)
	})

	t.Run("case=gets keys from multiple sets with per-key set names", func(t *testing.T) {
		otherSet := uuid.Must(uuid.NewV4()).String()
		otherKey := createJWK(t, reg, otherSet, "ES256")

		actual := gjson.Parse(cmdx.ExecNoErr(t, c, set, otherSet))
		assert.False(t, actual.Get("set").Exists(), "no single top-level set name exists for keys merged from multiple sets: %s", actual.Raw)
		require.Len(t, actual.Get("keys").Array(), 2, actual.Raw)

		setsByKid := make(map[string]string)
		for _, key := range actual.Get("keys").Array() {
			setsByKid[key.Get("kid").String()] = key.Get("set").String()
		}
		expected, err := reg.KeyManager().GetKeySet(t.Context(), set)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{
			expected.Keys[0].KeyID: set,
			otherKey.KeyID:         otherSet,
		}, setsByKid, actual.Raw)
	})

	t.Run("case=gets jwks public", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, set, "--public"))

		expected, err := reg.KeyManager().GetKeySet(t.Context(), set)
		require.NoError(t, err)

		assert.Equal(t, expected.Keys[0].KeyID, actual.Get("kid").String())
		assert.Equal(t, set, actual.Get("set").String(), actual.Raw)

		assert.NotEmptyf(t, actual.Get("kid").String(), "Expected kid to be set but got: %s", actual.Raw)
		assert.Empty(t, actual.Get("p").String(), "public key should not contain private key components: %s", actual.Raw)
	})

	t.Run("case=table output contains the set name", func(t *testing.T) {
		stdout := cmdx.ExecNoErr(t, c, set, "--format", "table")
		assert.Contains(t, stdout, set)
	})
}
