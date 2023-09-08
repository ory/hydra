// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
)

func TestPerformClientCredentialsGrant(t *testing.T) {
	c := cmd.NewPerformClientCredentialsCmd()
	public, _, reg := setupRoutes(t, c)
	require.NoError(t, c.Flags().Set(cmdx.FlagEndpoint, public.URL))

	expected := createClientCredentialsClient(t, reg)
	t.Run("case=exchanges for access token", func(t *testing.T) {
		result := cmdx.ExecNoErr(t, c, "--client-id", expected.GetID(), "--client-secret", expected.Secret)
		actual := gjson.Parse(result)
		assert.Equal(t, "bearer", actual.Get("token_type").String(), result)
		assert.NotEmpty(t, actual.Get("access_token").String(), result)
		assert.NotEmpty(t, actual.Get("expiry").String(), result)
		assert.Empty(t, actual.Get("refresh_token").String(), result)
		assert.Empty(t, actual.Get("id_token").String(), result)
	})
}
