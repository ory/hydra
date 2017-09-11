package cmd

import (
	"github.com/spf13/cobra"
)

var clientsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a client by id",
	Run:   cmdHandler.Clients.GetClient,
}

func init() {
	clientsCmd.AddCommand(clientsGetCmd)
}
