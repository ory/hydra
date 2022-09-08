package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewRevokeCmd(root *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke resources",
	}
	cmd.AddCommand(NewRevokeTokenCmd(root))
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
