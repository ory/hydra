package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewRevokeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke resources",
	}
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
