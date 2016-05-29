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
}
