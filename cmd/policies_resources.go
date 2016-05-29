package cmd

import (
	"github.com/spf13/cobra"
)

// policiesResourcesCmd represents the resources command
var policiesResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Manage a policy's resource matches",
}

func init() {
	policiesCmd.AddCommand(policiesResourcesCmd)
}
