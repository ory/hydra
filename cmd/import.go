package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewImportCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "Import resources",
	}
	cmd.AddCommand(NewImportClientCmd(root))
	cliclient.RegisterClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
