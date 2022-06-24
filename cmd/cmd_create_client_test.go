package cmd_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/cmd"
	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/snapshotx"
)

func base64EncodedPGPPublicKey(t *testing.T) string {
	t.Helper()
	gpgPublicKey, err := ioutil.ReadFile("../test/stub/pgp.pub")
	if err != nil {
		t.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(gpgPublicKey)
}

func setup(t *testing.T, cmd *cobra.Command) driver.Registry {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	reg := internal.NewMockedRegistry(t, &contextx.Default{})
	_, admin := testhelpers.NewOAuth2Server(ctx, t, reg)

	cliclient.RegisterClientFlags(cmd.Flags())
	cmdx.RegisterFormatFlags(cmd.Flags())
	require.NoError(t, cmd.Flags().Set(cliclient.FlagEndpoint, admin.URL))
	require.NoError(t, cmd.Flags().Set(cmdx.FlagFormat, string(cmdx.FormatJSON)))
	return reg
}

var snapshotExcludedClientFields = []snapshotx.ExceptOpt{
	snapshotx.ExceptPaths("client_id"),
	snapshotx.ExceptPaths("registration_access_token"),
	snapshotx.ExceptPaths("registration_client_uri"),
	snapshotx.ExceptPaths("client_secret"),
	snapshotx.ExceptPaths("created_at"),
	snapshotx.ExceptPaths("updated_at"),
}

func TestCreateClient(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewCreateClientsCommand(new(cobra.Command))
	reg := setup(t, c)

	t.Run("case=creates successfully", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "create", "client"))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.NotEmpty(t, actual.Get("client_secret").String())

		expected, err := reg.ClientManager().GetClient(ctx, actual.Get("client_id").String())
		require.NoError(t, err)

		assert.Equal(t, expected.GetID(), actual.Get("client_id").String())
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=supports setting flags", func(t *testing.T) {
		useSecret := "some-userset-secret"
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "create", "client",
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
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, "create", "client",
			"--secret", "some-userset-secret",
			"--pgp-key", base64EncodedPGPPublicKey(t),
		))
		assert.NotEmpty(t, actual.Get("client_id").String())
		assert.NotEmpty(t, actual.Get("client_secret").String())

		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})
}
