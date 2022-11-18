// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
)

func TestCreateJWKS(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewCreateJWKSCmd()
	reg := setup(t, c)

	t.Run("case=creates successfully", func(t *testing.T) {
		set := uuid.Must(uuid.NewV4()).String()
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, set, "--use", "enc", "--alg", "ES256"))
		assert.Len(t, actual.Get("keys.0").Array(), 1, "%s", actual.Raw)
		assert.NotEmpty(t, actual.Get("keys.0.kid").Array(), "%s", actual.Raw)
		assert.Equal(t, "ES256", actual.Get("keys.0.alg").String(), "%s", actual.Raw)

		expected, err := reg.KeyManager().GetKeySet(ctx, set)
		require.NoError(t, err)
		assert.Equal(t, expected.Keys[0].KeyID, actual.Get("keys.0.kid").String())
	})
}
