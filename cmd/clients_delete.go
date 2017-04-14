package cmd

import (
	"github.com/spf13/cobra"
)

// clientsDeleteCmd represents the delete command
var clientsDeleteCmd = &cobra.Command{
	Use:   "delete <id> [<id>...]",
	Short: "Delete an OAuth2 client",
	Run:   cmdHandler.Clients.DeleteClient,
}

func init() {
	clientsCmd.AddCommand(clientsDeleteCmd)
}
