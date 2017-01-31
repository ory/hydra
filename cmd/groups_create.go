package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <id>",
	Short: "Create a warden group",
	Long: `This command creates a warden group.

Example:
  hydra groups create my-group
`,
	Run: cmdHandler.Groups.CreateGroup,
}

func init() {
	groupsCmd.AddCommand(createCmd)
}
