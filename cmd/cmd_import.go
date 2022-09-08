package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewImportCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "Import resources",
	}
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
