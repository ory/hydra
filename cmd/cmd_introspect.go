package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewIntrospectCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "introspect",
		Short: "Introspect resources",
	}
	cmd.AddCommand(NewIntrospectTokenCmd(root))
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
