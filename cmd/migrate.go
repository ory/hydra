// Copyright © 2022 Ory Corp

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/configx"
)

func NewMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Various migration helpers",
	}
	configx.RegisterFlags(cmd.PersistentFlags())
	return cmd
}
