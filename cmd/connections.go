package cmd

import (
	"github.com/spf13/cobra"
)

// connectionsCmd represents the connections command
var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Manage SSO connections",
	Long: `With SSO connections, an identity can be associated with a social login provider like
Google, Twitter, or any other SSO provider.`,
}

func init() {
	RootCmd.AddCommand(connectionsCmd)
	connectionsCmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
}
