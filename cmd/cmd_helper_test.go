package cmd_test

import (
	"context"
	"encoding/base64"
	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/snapshotx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
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
	snapshotx.ExceptNestedKeys("client_id"),
	snapshotx.ExceptNestedKeys("registration_access_token"),
	snapshotx.ExceptNestedKeys("registration_client_uri"),
	snapshotx.ExceptNestedKeys("client_secret"),
	snapshotx.ExceptNestedKeys("created_at"),
	snapshotx.ExceptNestedKeys("updated_at"),
}
