package clients

import (
	"bytes"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/x"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path"
	"testing"
)

func TestCreateCmd(t *testing.T) {
	runCase := newTestSuite(t, newCreateCmd)

	runCase("creates a client from flags", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		redirectURI := "https://test1.com"
		require.NoError(t, cmd.Flags().Set(FlagCallbacks, redirectURI))

		stdOut := cmdx.ExecNoErr(t, cmd)

		assertPartialClient(
			t,
			map[string]interface{}{
				"redirect_uris": []string{redirectURI},
			},
			stdOut,
			reg,
		)
	})

	runCase("creates a client from STD_IN", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		inputClient := map[string]interface{}{
			"redirect_uris": []string{"https://test2.com"},
			"client_name":   "test client",
		}

		stdOut, stdErr, err := cmdx.Exec(t, cmd, bytes.NewBufferString(requireMarshaledJSON(t, inputClient)), "-")
		require.NoError(t, err, stdOut, stdErr)
		require.Len(t, stdErr, 0)

		assertPartialClient(t, inputClient, stdOut, reg)
	})

	runCase("creates a client from file", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		require.NoError(t, cmd.Flags().Set(cmdx.FlagFormat, string(cmdx.FormatJSON)))

		dir := t.TempDir()
		fn := path.Join(dir, "client.json")
		inputClient := map[string]interface{}{
			"redirect_uris": []string{"https://test3.com"},
			"client_name":   "other test client",
		}
		require.NoError(t, ioutil.WriteFile(fn, []byte(requireMarshaledJSON(t, inputClient)), 0600), fn)

		stdOut := cmdx.ExecNoErr(t, cmd, fn)

		assertPartialClient(t, inputClient, stdOut, reg)
	})

	runCase("flags overwrite file input", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		actualClientName := "not test client foo"
		require.NoError(t, cmd.Flags().Set(FlagClientName, actualClientName))

		inputClient := map[string]interface{}{
			"redirect_uris": []string{"https://test4.com"},
			"client_name":   "test client foo",
		}

		stdOut, stdErr, err := cmdx.Exec(t, cmd, bytes.NewBufferString(requireMarshaledJSON(t, inputClient)), "-")
		require.NoError(t, err)
		require.Len(t, stdErr, 0)

		inputClient["client_name"] = actualClientName
		assertPartialClient(t, inputClient, stdOut, reg)
	})

	runCase("warns about secrets in through flag", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		clientSecret, err := x.GenerateSecret(26)
		require.NoError(t, err)
		inputClient := map[string]interface{}{
			"client_secret": string(clientSecret),
		}
		require.NoError(t, cmd.Flags().Set(FlagSecret, string(clientSecret)))

		stdOut, stdErr, err := cmdx.Exec(t, cmd, nil)
		require.NoError(t, err)
		assert.Contains(t, stdErr, "secret might leak")

		assertPartialClient(t, inputClient, stdOut, reg)
	})
}
