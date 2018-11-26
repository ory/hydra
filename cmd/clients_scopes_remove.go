package cmd

import (
	"github.com/spf13/cobra"
)

// policyActionsAddCmd represents the add command
var clientsScopesRemoveCmd = &cobra.Command{
	Use:   "remove <client> <scope> [<scope>...]",
	Short: "Remove scopes from the client",
	Run:   cmdHandler.Clients.RemoveScopeFromClient,
}

func init() {
	clientsScopesCmd.AddCommand(clientsScopesRemoveCmd)
}
