package cmd

import (
	"github.com/spf13/cobra"
)

// policyResourcesAddCmd represents the add command
var policyResourcesAddCmd = &cobra.Command{
	Use:   "add <policy> <subject> [<subject>...]",
	Short: "Add subjects to the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies resources add my-policy some-item-123 some-item-<[234|345]>`,
	Run: cmdHandler.Policies.AddResourceToPolicy,
}

func init() {
	policiesResourcesCmd.AddCommand(policyResourcesAddCmd)
}
