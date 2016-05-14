package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <id> [<id>...]",
	Short: "Remove a SSO connection",
	Long: `Example:
  hydra connections delete 4adb79ab-f89d-4445-ab01-ff670e51cefa`,
	Run: cmdHandler.Connections.DeleteConnection,
}

func init() {
	connectionsCmd.AddCommand(deleteCmd)
}
