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

func TestUpdateClient(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewUpdateClientCmd()
	reg := setup(t, c)

	original := createClient(t, reg, nil)
	t.Run("case=creates successfully", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "--grant-type", "implicit", original.GetID()))
		expected, err := reg.ClientManager().GetClient(ctx, actual.Get("client_id").String())
		require.NoError(t, err)

		assert.Equal(t, expected.GetID(), actual.Get("client_id").String())
		assert.Equal(t, "implicit", actual.Get("grant_types").Array()[0].String())
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=supports encryption", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c,
			original.GetID(),
			"--secret", "some-userset-secret",
			"--pgp-key", base64EncodedPGPPublicKey(t),
		))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.NotEmpty(t, actual.Get("client_secret").String())

		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})
}
