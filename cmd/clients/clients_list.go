package clients

import (
	"github.com/ory/hydra/cmd"
	"github.com/spf13/cobra"
)

var clientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List OAuth 2.0 Clients",
	Long: `This command list an OAuth 2.0 Clients.

Example:
  hydra clients list`,
	Run: cmd.cmdHandler.Clients.ListClients,
}

func init() {
	clientsListCmd.Flags().Int("limit", 20, "The maximum amount of policies returned.")
	clientsListCmd.Flags().Int("page", 1, "The number of page.")
}
