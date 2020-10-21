package clients

import (
	"context"
	"errors"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/sqlcon"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteCmd(t *testing.T) {
	runCase := newTestSuite(t, newDeleteCmd)

	runCase("deletes a client by ID", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		c := &client.Client{
			ID: "test-cid",
		}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))

		stdOut := cmdx.ExecNoErr(t, cmd, c.ID)

		assert.Equal(t, c.ID + "\n", stdOut)

		// check if client is deleted in the manager
		_, err := reg.ClientManager().GetClient(context.Background(), c.ID)
		assert.True(t, errors.Is(err, sqlcon.ErrNoRows))
	})

	runCase("returns error when the client ID is unknown", func(t *testing.T, cmd *cobra.Command, reg driver.Registry) {
		stdErr := cmdx.ExecExpectedErr(t, cmd, "undefined client")

		t.Error(stdErr)
	})
}
