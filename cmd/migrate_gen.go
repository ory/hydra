// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cli"
)

func NewMigrateGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen </source/path> </target/path>",
		Short: "Generate migration files from migration templates",
		Run:   cli.NewHandler(nil, nil, nil).Migration.MigrateGen,
	}
	cmd.Flags().StringSlice("dialects", []string{"sqlite", "cockroach", "mysql", "postgres"}, "Expect migrations for these dialects and no others to be either explicitly defined, or to have a generic fallback. \"\" disables dialect validation.")
	return cmd
}
