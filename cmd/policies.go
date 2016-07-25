package cmd

import (
	"github.com/spf13/cobra"
)

// policiesCmd represents the policies command
var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "Manage access control policies",
}

func init() {
	RootCmd.AddCommand(policiesCmd)
	policiesCmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
}
