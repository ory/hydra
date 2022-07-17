package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewListCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List resources",
	}
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
