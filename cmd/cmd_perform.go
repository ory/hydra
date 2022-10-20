package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewPerformCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "perform",
		Short: "Perform OAuth 2.0 Flows",
	}
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
