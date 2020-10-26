package clients

import (
	"bytes"
	"context"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUpdateCmd(t *testing.T) {
	const updatedName = "updated name"

	runCase := newTestSuite(t, newUpdateCmd)

	newClient := func(t *testing.T, reg driver.Registry) *client.Client {
		c := &client.Client{
			ID:   "test client ID",
			Name: "initial name",
		}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))

		return c
	}

	runCase("updates client from flags", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		c := newClient(t, reg)

		require.NoError(t, cmd.Flags().Set(FlagClientName, updatedName))

		stdOut := cmdx.ExecNoErr(t, cmd, c.ID)

		assertPartialClient(t, map[string]interface{}{
			"client_id":   c.ID,
			"client_name": updatedName,
		}, stdOut, reg)
	})

	runCase("updates client from file", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		c := newClient(t, reg)

		fn := filepath.Join(t.TempDir(), "client.json")
		fileClient := map[string]interface{}{
			"client_id": c.ID,
			"client_name": updatedName,
		}
		require.NoError(t, ioutil.WriteFile(fn, []byte(requireMarshaledJSON(t, fileClient)), 0600))

		stdOut := cmdx.ExecNoErr(t, cmd, c.ID, fn)

		assertPartialClient(t, fileClient, stdOut, reg)
	})

	runCase("updates client from STD_IN", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		c := newClient(t, reg)

		inputClient := map[string]interface{}{
			"client_id": c.ID,
			"client_name": updatedName,
		}
		stdIn := bytes.NewBufferString(requireMarshaledJSON(t, inputClient))

		stdOut, stdErr, err := cmdx.Exec(t, cmd, stdIn, c.ID, "-")
		require.NoError(t, err)
		require.Len(t, stdErr, 0)

		assertPartialClient(t, inputClient, stdOut, reg)
	})

	runCase("updates client from file and flags", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		c := newClient(t, reg)

		inputClient := map[string]interface{}{
			"client_id": c.ID,
			"client_name": updatedName,
		}
		stdIn := bytes.NewBufferString(requireMarshaledJSON(t, inputClient))
		flagName := "updated name from flag"
		require.NoError(t, cmd.Flags().Set(FlagClientName, flagName))

		// set the actually expected name
		inputClient["client_name"] = flagName

		stdOut, stdErr, err := cmdx.Exec(t, cmd, stdIn, c.ID, "-")
		require.NoError(t, err)
		require.Len(t, stdErr, 0)

		assertPartialClient(t, inputClient, stdOut, reg)
	})
}
