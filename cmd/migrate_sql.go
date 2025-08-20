// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/popx"

	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/driver"
)

func NewMigrateSQLCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "sql [database_url]",
		Deprecated: "Please use `hydra migrate sql up` instead.",
		Short:      "Perform SQL migrations",
		Long: `Run this command on a fresh SQL installation and when you upgrade Hydra to a new minor version. For example,
upgrading Hydra 0.7.0 to 0.8.0 requires running this command.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

You can read in the database URL using the -e flag, for example:
	export DSN=...
	hydra migrate sql up -e

### WARNING ###

Before running this command on an existing database, create a back up!`,
		RunE: cli.NewHandler(dOpts).Migration.MigrateSQLUp,
	}

	cmd.Flags().BoolP("yes", "y", false, "If set all confirmation requests are accepted without user interaction.")
	cmd.PersistentFlags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")

	cmd.AddCommand(newMigrateSQLDownCmd(dOpts))
	cmd.AddCommand(newMigrateSQLUpCmd(dOpts))
	cmd.AddCommand(newMigrateSQLStatusCmd(dOpts))

	return cmd
}

func newMigrateSQLDownCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	return popx.NewMigrateSQLDownCmd(cli.NewHandler(dOpts).Migration.MigrateSQLDown)
}

func newMigrateSQLStatusCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	return popx.NewMigrateSQLStatusCmd(cli.NewHandler(dOpts).Migration.MigrateStatus)
}

func newMigrateSQLUpCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	return popx.NewMigrateSQLUpCmd(cli.NewHandler(dOpts).Migration.MigrateSQLUp)
}
