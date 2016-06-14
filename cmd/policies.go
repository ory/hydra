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
	var dry bool
	c.Dry = &dry

	RootCmd.AddCommand(policiesCmd)
	policiesCmd.PersistentFlags().BoolVar(c.Dry, "dry", false, "do not execute the command but show the corresponding curl command instead")
}
