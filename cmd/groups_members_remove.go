package cmd

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <group> <member> [<member>...]",
	Short: "Remove members from a warden group",
	Long: `This command removes members from a warden group.

Example:
  hydra groups members remove my-group peter julia
`,
	Run: cmdHandler.Groups.RemoveMembers,
}

func init() {
	groupsMembersCmd.AddCommand(removeCmd)
}
