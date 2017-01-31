package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <group> <member> [<member>...]",
	Short: "Add members to a warden group",
	Long: `This command adds members to a warden group.

Example:
  hydra groups members add my-group peter julia
`,
	Run: cmdHandler.Groups.AddMembers,
}

func init() {
	groupsMembersCmd.AddCommand(addCmd)
}
