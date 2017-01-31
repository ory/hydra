package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a warden group",
	Long: `This command deletes a warden group.

Example:
  hydra groups delete my-group
`,
	Run: cmdHandler.Groups.DeleteGroup,
}

func init() {
	groupsCmd.AddCommand(deleteCmd)

}
