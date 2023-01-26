// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"context"
	"testing"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
)

func TestIntrospectToken(t *testing.T) {
	c := cmd.NewIntrospectTokenCmd()
	public, admin, reg := setupRoutes(t, c)
	require.NoError(t, c.Flags().Set(cmdx.FlagEndpoint, admin.URL))

	expected := createClientCredentialsClient(t, reg)
	cc := clientcredentials.Config{
		ClientID:     expected.GetID(),
		ClientSecret: expected.Secret,
		TokenURL:     public.URL + "/oauth2/token",
		Scopes:       []string{},
	}

	t.Run("case=checks valid token", func(t *testing.T) {
		token, err := cc.Token(context.Background())
		require.NoError(t, err)

		actual := gjson.Parse(cmdx.ExecNoErr(t, c, token.AccessToken))
		assert.Equal(t, expected.GetID(), actual.Get("sub").String())
		assert.Equal(t, expected.GetID(), actual.Get("client_id").String())
		assert.True(t, actual.Get("active").Bool())
	})

	t.Run("case=checks invalid token", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "invalid-token"))
		assert.Empty(t, actual.Get("sub").String())
		assert.Empty(t, actual.Get("client_id").String())
		assert.False(t, actual.Get("active").Bool())
	})
}
