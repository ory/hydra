package cmd

import (
	"github.com/spf13/cobra"
)

// policiesDeleteCmd represents the delete command
var policiesDeleteCmd = &cobra.Command{
	Use:   "delete <policy>",
	Short: "Delete a policy",
	Run:   cmdHandler.Policies.DeletePolicy,
}

func init() {
	policiesCmd.AddCommand(policiesDeleteCmd)
}
