// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/ory/x/cmdx"
	"github.com/ory/x/configx"
	"github.com/ory/x/servicelocatorx"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/driver"
)

func NewMigrateStatusCmd(slOpts []servicelocatorx.Option, dOpts []driver.OptionsModifier, cOpts []configx.OptionModifier) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get the current migration status",
		RunE:  cli.NewHandler(slOpts, dOpts, cOpts).Migration.MigrateStatus,
	}

	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	cmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	cmd.Flags().Bool("block", false, "Block until all migrations have been applied")

	return cmd
}
