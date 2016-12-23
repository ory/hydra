package cmd

import (
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find <subject>",
	Short: "Find all groups a subject belongs to",
	Long: `This command find all groups a subject belongs to.

Example:
  hydra groups find peter
`,
	Run: cmdHandler.Groups.FindGroups,
}

func init() {
	groupsCmd.AddCommand(findCmd)
}
