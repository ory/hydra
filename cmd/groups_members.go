package cmd

import (
	"github.com/spf13/cobra"
)

var groupsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage warden group members",
}

func init() {
	groupsCmd.AddCommand(groupsMembersCmd)
}
