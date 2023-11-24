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

func TestGetClient(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewGetClientsCmd()
	reg := setup(t, c)

	expected := createClient(t, reg, nil)
	t.Run("case=gets client", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, expected.GetID()))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.Empty(t, actual.Get("client_secret").String())

		expected, err := reg.ClientManager().GetClient(ctx, actual.Get("client_id").String())
		require.NoError(t, err)

		assert.Equal(t, expected.GetID(), actual.Get("client_id").String())
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=gets multiple clients", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, expected.GetID(), expected.ID))
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})
}
