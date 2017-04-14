package cmd

import (
	"github.com/spf13/cobra"
)

// policyResourcesRemoveCmd represents the remove command
var policyResourcesRemoveCmd = &cobra.Command{
	Use:   "remove <policy> <resource> [<resource>...]",
	Short: "Remove resources from the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies resources remove my-policy some-item-123 some-item-<[234|345]>`,
	Run: cmdHandler.Policies.RemoveResourceFromPolicy,
}

func init() {
	policiesResourcesCmd.AddCommand(policyResourcesRemoveCmd)
}
