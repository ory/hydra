package clients

import (
	"context"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/driver"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"testing"
)

func TestListCmd(t *testing.T) {
	runCase := newTestSuite(t, newListCmd)

	runCase("lists all clients", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		clients := []client.Client{
			{
				ID: "first client",
				Name: "John C",
			},
			{
				ID: "second client",
				Name: "John B",
			},
		}
		for _, c := range clients {
			require.NoError(t, reg.ClientManager().CreateClient(context.Background(), &c))
		}

		stdOut := cmdx.ExecNoErr(t, cmd)

		assert.Equal(t, 2, len(gjson.Parse(stdOut).Array()))
		assert.Equal(t, clients[0].ID, gjson.Get(stdOut, "0.client_id").Str)
		assert.Equal(t, clients[1].ID, gjson.Get(stdOut, "1.client_id").Str)
		assert.Equal(t, clients[0].Name, gjson.Get(stdOut, "0.client_name").Str)
		assert.Equal(t, clients[1].Name, gjson.Get(stdOut, "1.client_name").Str)
	})

	runCase("applies pagination", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		clients := []client.Client{
			{ ID: "a" },
			{ ID: "b" },
			{ ID: "c" },
			{ ID: "d" },
			{ ID: "e" },
		}
		for _, c := range clients {
			require.NoError(t, reg.ClientManager().CreateClient(context.Background(), &c))
		}

		require.NoError(t, cmd.Flags().Set(cli.FlagPage, "2"))
		require.NoError(t, cmd.Flags().Set(cli.FlagLimit, "2"))

		stdOut := cmdx.ExecNoErr(t, cmd)
		assert.Equal(t, 2, len(gjson.Parse(stdOut).Array()))
		assert.Equal(t, clients[2].ID, gjson.Get(stdOut, "0.client_id").Str)
		assert.Equal(t, clients[3].ID, gjson.Get(stdOut, "1.client_id").Str)
	})
}
