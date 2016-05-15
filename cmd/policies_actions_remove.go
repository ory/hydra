package cmd

import (
	"github.com/spf13/cobra"
)

// policiesActionsDeleteCmd represents the remove command
var policiesActionsDeleteCmd = &cobra.Command{
	Use:   "delete <policy> <subject> [<subject>...]",
	Short: "Remove actions from the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies actions delete my-policy create delete <[get|update]>`,
	Run: cmdHandler.Policies.RemoveActionFromPolicy,
}

func init() {
	policyActionsCmd.AddCommand(policiesActionsDeleteCmd)
}
