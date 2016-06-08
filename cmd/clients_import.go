package cmd

import (
	"github.com/spf13/cobra"
)

var clientsImportCmd = &cobra.Command{
	Use:   "import <path/to/file.json> [<path/to/other/file.json>...]",
	Short: "Import clients from JSON files",
	Run:   cmdHandler.Clients.ImportClients,
}

func init() {
	clientsCmd.AddCommand(clientsImportCmd)
}
