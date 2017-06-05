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
	groupsCmd.PersistentFlags().Bool("fake-tls-termination", false, `fake tls termination by adding "X-Forwarded-Proto: https"" to http headers`)
	groupsCmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
}
