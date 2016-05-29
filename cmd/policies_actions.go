package cmd

import (
	"github.com/spf13/cobra"
)

// policyActionsCmd represents the actions command
var policyActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Manage a policy's action matches",
}

func init() {
	policiesCmd.AddCommand(policyActionsCmd)
}
