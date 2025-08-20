// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"bytes"
	"cmp"
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	_ "embed"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
)

//go:embed stub/jwk.json
var stubJsonWebKeySet []byte

func TestImportJWKS(t *testing.T) {
	c := cmd.NewKeysImportCmd()
	_ = setup(t, c)

	t.Run("case=imports without alg fails", func(t *testing.T) {
		result := cmdx.ExecExpectedErr(t, c, uuid.Must(uuid.NewV4()).String(), "stub/ecdh.key")
		assert.Contains(t, result, "Flag `--alg` is required when imported key does not define the `alg` field itself.")
	})

	for _, tc := range [][2]string{
		{"ES256", "stub/ecdh.key"},
		{"ES384", "stub/ecdh.key"},
		{"ES256", "stub/ecdh.pub"},
		{"RS256", "stub/rsa.key"},
		{"RS256", "stub/rsa.pub"},
		{"RS512", "stub/rsa.pub"},
		{"", "stub/jwk.json"},
	} {
		t.Run("case=imports encoded "+tc[1]+" key from file", func(t *testing.T) {
			args := []string{uuid.Must(uuid.NewV4()).String(), tc[1]}
			if len(tc[0]) > 0 {
				args = append(args, "--alg", tc[0])
			}

			actual := gjson.Parse(cmdx.ExecNoErr(t, c, args...))
			assert.Len(t, actual.Get("keys.0").Array(), 1, "%s", actual.Raw)
			assert.NotEmpty(t, actual.Get("keys.0.kid").String(), "%s", actual.Raw)
			assert.NotEmpty(t, cmp.Or(actual.Get("keys.0.x").String(), actual.Get("keys.0.n").String()), "%s", actual.Raw)
			assert.Equal(t, cmp.Or(tc[0], "RS256"), actual.Get("keys.0.alg").String(), "%s", actual.Raw)
		})
	}

	t.Run("case=imports JWK key from STDIN", func(t *testing.T) {
		stdin := bytes.NewBuffer(stubJsonWebKeySet)

		stdout, stderr, err := cmdx.Exec(t, c, stdin, uuid.Must(uuid.NewV4()).String())
		require.NoError(t, err, stderr)

		snapshotx.SnapshotT(t, json.RawMessage(stdout), snapshotx.ExceptNestedKeys("set", "kid"))
	})
}
