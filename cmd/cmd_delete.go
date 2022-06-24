package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewDeleteCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}
	cmd.AddCommand(NewDeleteClientCmd(root))
	cliclient.RegisterClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
