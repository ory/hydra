// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/snapshotx"
)

func TestRevokeToken(t *testing.T) {
	c := cmd.NewRevokeTokenCmd()
	public, _, reg := setupRoutes(t, c)
	require.NoError(t, c.Flags().Set(cmdx.FlagEndpoint, public.URL))

	expected := createClientCredentialsClient(t, reg)
	cc := clientcredentials.Config{
		ClientID: expected.GetID(), ClientSecret: expected.Secret,
		TokenURL: public.URL + "/oauth2/token",
	}

	t.Run("case=revokes valid token but without client credentials", func(t *testing.T) {
		token, err := cc.Token(context.Background())
		require.NoError(t, err)

		snapshotx.SnapshotT(t, cmdx.ExecExpectedErr(t, c, token.AccessToken))
	})

	t.Run("case=revokes valid token", func(t *testing.T) {
		token, err := cc.Token(context.Background())
		require.NoError(t, err)

		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "--client-id", expected.GetID(), "--client-secret", expected.Secret, token.AccessToken))
		assert.Equal(t, token.AccessToken, actual.String())
	})

	t.Run("case=checks invalid token", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "invalid-token"))
		assert.Equal(t, "invalid-token", actual.String())
	})
}
