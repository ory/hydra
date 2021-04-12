package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cli"
)

func NewClientsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List OAuth 2.0 Clients",
		Long: `This command list an OAuth 2.0 Clients.

Example:
  hydra clients list`,
		Run: cli.NewHandler().Clients.ListClients,
	}
	cmd.Flags().Int("limit", 20, "The maximum amount of policies returned.")
	cmd.Flags().Int("page", 1, "The number of page.")
	return cmd
}
