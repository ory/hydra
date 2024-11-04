// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/tidwall/sjson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
)

func TestUpdateClient(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewUpdateClientCmd()
	reg := setup(t, c)

	original := createClient(t, reg, nil)
	t.Run("case=creates successfully", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "--grant-type", "implicit", original.GetID()))
		expected, err := reg.ClientManager().GetClient(ctx, actual.Get("client_id").Str)
		require.NoError(t, err)

		assert.Equal(t, expected.GetID(), actual.Get("client_id").Str)
		assert.Equal(t, "implicit", actual.Get("grant_types").Array()[0].Str)
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=supports encryption", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c,
			original.GetID(),
			"--secret", "some-userset-secret",
			"--pgp-key", base64EncodedPGPPublicKey(t),
		))
		assert.Equal(t, original.ID, actual.Get("client_id").Str)
		assert.NotEmpty(t, actual.Get("client_secret").Str)
		assert.NotEqual(t, original.Secret, actual.Get("client_secret").Str)

		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=updates from file", func(t *testing.T) {
		original, err := reg.ClientManager().GetConcreteClient(ctx, original.GetID())
		require.NoError(t, err)

		raw, err := json.Marshal(original)
		require.NoError(t, err)

		t.Run("file=stdin", func(t *testing.T) {
			raw, err = sjson.SetBytes(raw, "client_name", "updated through file stdin")
			require.NoError(t, err)

			stdout, stderr, err := cmdx.Exec(t, c, bytes.NewReader(raw), original.GetID(), "--file", "-")
			require.NoError(t, err, stderr)

			actual := gjson.Parse(stdout)
			assert.Equal(t, original.ID, actual.Get("client_id").Str)
			assert.Equal(t, "updated through file stdin", actual.Get("client_name").Str)

			snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
		})

		t.Run("file=from disk", func(t *testing.T) {
			raw, err = sjson.SetBytes(raw, "client_name", "updated through file from disk")
			require.NoError(t, err)

			fn := writeTempFile(t, json.RawMessage(raw))

			stdout, stderr, err := cmdx.Exec(t, c, nil, original.GetID(), "--file", fn)
			require.NoError(t, err, stderr)

			actual := gjson.Parse(stdout)
			assert.Equal(t, original.ID, actual.Get("client_id").Str)
			assert.Equal(t, "updated through file from disk", actual.Get("client_name").Str)

			snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
		})
	})
}
