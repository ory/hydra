package clients

import (
	"context"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCmd(t *testing.T) {
	runCase := newTestSuite(t, newGetCmd)

	runCase("gets a client by ID", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		c := &client.Client{
			ID: "some id",
		}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))

		stdOut := cmdx.ExecNoErr(t, cmd, c.ID)

		assertPartialClient(t, map[string]interface{}{"client_id": c.ID}, stdOut, reg)
	})

	runCase("fails when ID is unkown", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		stdErr := cmdx.ExecExpectedErr(t, cmd, "unknown ID")

		assert.Contains(t, stdErr, "401")
	})
}
