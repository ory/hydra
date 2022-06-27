package cmd_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
)

func TestPerformClientCredentialsGrant(t *testing.T) {
	c := cmd.NewPerformClientCredentialsCmd(new(cobra.Command))
	public, _, reg := setupRoutes(t, c)
	require.NoError(t, c.Flags().Set(cmdx.FlagEndpoint, public.URL))

	expected := createClient(t, reg, &client.Client{
		GrantTypes:              []string{"client_credentials"},
		TokenEndpointAuthMethod: "client_secret_post",
		Secret:                  uuid.Must(uuid.NewV4()).String()},
	)

	t.Run("case=exchanges for access token", func(t *testing.T) {
		result := cmdx.ExecNoErr(t, c, "--client-id", expected.ID.String(), "--client-secret", expected.Secret)
		actual := gjson.Parse(result)
		assert.Equal(t, "bearer", actual.Get("token_type").String(), result)
		assert.NotEmpty(t, actual.Get("access_token").String(), result)
		assert.NotEmpty(t, actual.Get("expiry").String(), result)
		assert.Empty(t, actual.Get("refresh_token").String(), result)
		assert.Empty(t, actual.Get("id_token").String(), result)
	})
}
