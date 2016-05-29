package cmd

import (
	"github.com/spf13/cobra"
)

// policiesCreateCmd represents the create command
var policiesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new policy",
	Long: `To create a policy, either specify the files flag or pass arguments to create it directly from the CLI.

Example
  hydra policies create -f [policy-a.json,policy-b.json]
  hydra policies create -s [peter,max] -r [blog,users] -a [post,ban] --allow`,
	Run: cmdHandler.Policies.CreatePolicy,
}

func init() {
	policiesCmd.AddCommand(policiesCreateCmd)

	policiesCreateCmd.Flags().StringSliceP("files", "f", []string{}, "A list of paths to JSON encoded policy files")
	policiesCreateCmd.Flags().StringP("id", "i", "", "The policy's id")
	policiesCreateCmd.Flags().StringSliceP("description", "d", []string{}, "The policy's description")
	policiesCreateCmd.Flags().StringSliceP("resources", "r", []string{}, "A list of resource regex strings this policy will match to (required)")
	policiesCreateCmd.Flags().StringSliceP("subjects", "s", []string{}, "A list of subject regex strings this policy will match to (required)")
	policiesCreateCmd.Flags().StringSliceP("actions", "a", []string{}, "A list of action regex strings this policy will match to (required)")
	policiesCreateCmd.Flags().Bool("allow", false, "A list of action regex strings this policy will match to")
}
