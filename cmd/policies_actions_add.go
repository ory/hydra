package cmd

import (
	"github.com/spf13/cobra"
)

// policyActionsAddCmd represents the add command
var policyActionsAddCmd = &cobra.Command{
	Use:   "add <policy> <subject> [<subject>...]",
	Short: "Add actions to the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies actions add my-policy create delete <[get|update]>`,
	Run: cmdHandler.Policies.AddActionToPolicy,
}

func init() {
	policyActionsCmd.AddCommand(policyActionsAddCmd)
}
