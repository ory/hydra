package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewCreateCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create resources",
	}
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
