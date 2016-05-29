package cmd

import (
	"github.com/spf13/cobra"
)

// policiesGetCmd represents the delete command
var policiesGetCmd = &cobra.Command{
	Use:   "get <policy>",
	Short: "View a policy",
	Run:   cmdHandler.Policies.GetPolicy,
}

func init() {
	policiesCmd.AddCommand(policiesGetCmd)
}
