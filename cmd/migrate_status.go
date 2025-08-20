// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/popx"

	"github.com/ory/hydra/v2/cmd/cli"
	"github.com/ory/hydra/v2/driver"
)

func NewMigrateStatusCmd(dOpts []driver.OptionsModifier) *cobra.Command {
	cmd := popx.RegisterMigrateStatusFlags(&cobra.Command{
		Use:        "status",
		Deprecated: "Please use `hydra migrate sql status` instead.",
		Short:      "Get the current migration status",
		RunE:       cli.NewHandler(dOpts).Migration.MigrateStatus,
	})
	cmd.PersistentFlags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	return cmd
}
