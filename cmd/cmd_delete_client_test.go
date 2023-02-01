// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/assertx"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/sqlcon"
)

func TestDeleteClient(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewDeleteClientCmd()
	reg := setup(t, c)

	t.Run("case=deletes client", func(t *testing.T) {
		expected := createClient(t, reg, nil)
		stdout := cmdx.ExecNoErr(t, c, expected.GetID())
		assert.Equal(t, fmt.Sprintf(`"%s"`, expected.GetID()), strings.TrimSpace(stdout))

		_, err := reg.ClientManager().GetClient(ctx, expected.GetID())
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})

	t.Run("case=deletes multiple clients", func(t *testing.T) {
		expected1 := createClient(t, reg, nil)
		expected2 := createClient(t, reg, nil)
		assertx.EqualAsJSON(t, []string{expected1.GetID(), expected2.GetID()}, json.RawMessage(cmdx.ExecNoErr(t, c, expected1.GetID(), expected2.GetID())))

		_, err := reg.ClientManager().GetClient(ctx, expected1.GetID())
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)

		_, err = reg.ClientManager().GetClient(ctx, expected2.GetID())
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})

	t.Run("case=one client deletion fails", func(t *testing.T) {
		expected := createClient(t, reg, nil)
		stdout, stderr, err := cmdx.Exec(t, c, nil, "i-do-not-exist", expected.GetID())
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf(`"%s"`, expected.GetID()), strings.TrimSpace(stdout))
		snapshotx.SnapshotT(t, stderr)

		_, err = reg.ClientManager().GetClient(ctx, expected.GetID())
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})
}
