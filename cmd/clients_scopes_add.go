package cmd

import (
	"github.com/spf13/cobra"
)

// policyActionsAddCmd represents the add command
var clientsScopesAddCmd = &cobra.Command{
	Use:   "add <client> <scope> [<scope>...]",
	Short: "Add scopes to the client",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra clients scopes add my-client newscope`,
	Run: cmdHandler.Clients.AddScopeToClient,
}

func init() {
	clientsScopesCmd.AddCommand(clientsScopesAddCmd)
}
