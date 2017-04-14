package cmd

import (
	"github.com/spf13/cobra"
)

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage warden groups",
}

func init() {
	RootCmd.AddCommand(groupsCmd)
}
