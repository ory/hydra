package cmd

import (
	"github.com/spf13/cobra"
)

// policyActionsAddCmd represents the add command
var clientsScopesCmd = &cobra.Command{
	Use:   "scopes",
	Short: "Manage a client's scopes",
}

func init() {
	clientsCmd.AddCommand(clientsScopesCmd)
}
