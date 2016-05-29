package cmd

import (
	"github.com/spf13/cobra"
)

// policyResourcesDeleteCmd represents the remove command
var policyResourcesDeleteCmd = &cobra.Command{
	Use:   "delete <policy> <resource> [<resource>...]",
	Short: "Remove resources from the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies resources remove my-policy some-item-123 some-item-<[234|345]>`,
	Run: cmdHandler.Policies.RemoveResourceFromPolicy,
}

func init() {
	policiesResourcesCmd.AddCommand(policyResourcesDeleteCmd)
}
