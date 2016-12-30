package cmd

import (
	"github.com/spf13/cobra"
)

// policiesActionsRemoveCmd represents the remove command
var policiesActionsRemoveCmd = &cobra.Command{
	Use:   "remove <policy> <subject> [<subject>...]",
	Short: "Remove actions from the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies actions remove my-policy create delete <[get|update]>`,
	Run: cmdHandler.Policies.RemoveActionFromPolicy,
}

func init() {
	policyActionsCmd.AddCommand(policiesActionsRemoveCmd)
}
