// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
)

func TestCreateClient(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewCreateClientsCommand()
	reg := setup(t, c)

	t.Run("case=creates successfully", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.NotEmpty(t, actual.Get("client_secret").String())

		expected, err := reg.ClientManager().GetClient(ctx, actual.Get("client_id").String())
		require.NoError(t, err)

		assert.Equal(t, expected.GetID(), actual.Get("client_id").String())
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=supports setting flags", func(t *testing.T) {
		useSecret := "some-userset-secret"
		actual := gjson.Parse(cmdx.ExecNoErr(t, c,
			"--secret", useSecret,
			"--metadata", `{"foo":"bar"}`,
			"--audience", "https://www.ory.sh/audience1",
			"--audience", "https://www.ory.sh/audience2",
		))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.Equal(t, useSecret, actual.Get("client_secret").String())

		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=supports encryption", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c,
			"--secret", "some-userset-secret",
			"--pgp-key", base64EncodedPGPPublicKey(t),
		))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.NotEmpty(t, actual.Get("client_secret").String())

		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})
}
