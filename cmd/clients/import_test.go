package clients

import (
	"github.com/ory/hydra/driver"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path"
	"testing"
)

func TestImportCmd(t *testing.T) {
	runCase := newTestSuite(t, newImportCmd)

	runCase("imports client from file", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		inputClient := map[string]interface{}{
			"client_id":   "some id",
			"client_name": "Test Client Name",
		}
		fn := path.Join(t.TempDir(), "client.json")

		require.NoError(t,
			ioutil.WriteFile(
				fn,
				[]byte(requireMarshaledJSON(t, inputClient)),
				0600,
			))

		stdOut := cmdx.ExecNoErr(t, cmd, fn)

		assertPartialClient(t, inputClient, stdOut, reg)
	})
}
