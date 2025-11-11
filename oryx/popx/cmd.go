// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/flagx"
)

type MigrationProvider interface {
	MigrationStatus(context.Context) (MigrationStatuses, error)
	MigrateUp(context.Context) error
	MigrateDown(context.Context, int) error
}

type MigrationPreparer interface {
	PrepareMigration(context.Context) error
}

func RegisterMigrateSQLUpFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP("yes", "y", false, "If set all confirmation requests are accepted without user interaction.")
	return cmd
}

func NewMigrateSQLUpCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return RegisterMigrateSQLUpFlags(&cobra.Command{
		Use:   "up [database_url]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Apply all pending SQL migrations",
		Long: `This command applies all pending SQL migrations for Ory {{ title .Root.Name }}.

:::warning

Before running this command, create a backup of your database. This command can be destructive as it may apply changes that cannot be easily reverted. Run this command close to the SQL instance (same VPC / same machine).

:::

It is recommended to review the migrations before running them. You can do this by running the command without the --yes flag:

	DSN=... {{ .CommandPath }} -e`,
		Example: `Apply all pending migrations:
	DSN=... {{ .CommandPath }} -e

Apply all pending migrations:
	DSN=... {{ .CommandPath }} -e --yes`,
		RunE: runE,
	})
}

func MigrateSQLUp(cmd *cobra.Command, p MigrationProvider) (err error) {
	// convert migration tables
	if prep, ok := p.(MigrationPreparer); ok {
		if err := prep.PrepareMigration(cmd.Context()); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not convert the migration table:\n%+v\n", err)
			return cmdx.FailSilently(cmd)
		}
	}

	// print migration status
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "The migration plan is as follows:")

	// print migration status
	status, err := p.MigrationStatus(cmd.Context())
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not get the migration status:\n%+v\n", errorsx.WithStack(err))
		return cmdx.FailSilently(cmd)
	}
	cmdx.PrintTable(cmd, status)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nThe SQL statements to be executed from top to bottom are:\n\n")
	for i := range status {
		if status[i].State == Pending {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ %s - %s ------------\n", status[i].Version, status[i].Name)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", status[i].ContentUp)
		}
	}

	if !flagx.MustGetBool(cmd, "yes") {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "To skip the next question use flag --yes (at your own risk).")
		if !cmdx.AskForConfirmation("Do you wish to execute this migration plan?", cmd.InOrStdin(), cmd.OutOrStdout()) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ WARNING ------------\n")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Migration aborted.")
			return nil
		}
	}

	// apply migrations
	if err := p.MigrateUp(cmd.Context()); err != nil {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ ERROR ------------\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not apply migrations:\n%+v\n", errorsx.WithStack(err))
		return cmdx.FailSilently(cmd)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ SUCCESS ------------\n")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Successfully applied migrations!")
	return nil
}

func RegisterMigrateSQLDownFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP("yes", "y", false, "If set all confirmation requests are accepted without user interaction.")
	cmd.Flags().Int("steps", 0, "The number of migrations to roll back.")
	return cmd
}

func NewMigrateSQLDownCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return RegisterMigrateSQLDownFlags(&cobra.Command{
		Use:   "down [database_url]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Rollback the last applied SQL migrations",
		Long: `This command rolls back the last applied SQL migrations for Ory {{ title .Root.Name }}.

:::warning

Before running this command, create a backup of your database. This command can be destructive as it may revert changes made by previous migrations. Run this command close to the SQL instance (same VPC / same machine).

:::

It is recommended to review the migrations before running them. You can do this by running the command without the --yes flag:

	DSN=... {{ .CommandPath }} -e`,
		Example: `See the current migration status:
	DSN=... {{ .CommandPath }} -e

Rollback the last 10 migrations:
	{{ .CommandPath }} $DSN --steps 10

Rollback the last 10 migrations without confirmation:
	DSN=... {{ .CommandPath }} -e --yes --steps 10`,
		RunE: runE,
	})
}

func MigrateSQLDown(cmd *cobra.Command, p MigrationProvider) (err error) {
	steps := flagx.MustGetInt(cmd, "steps")
	if steps < 0 {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Flag --steps must be larger than 0.")
		return cmdx.FailSilently(cmd)
	}

	// convert migration tables
	if prep, ok := p.(MigrationPreparer); ok {
		if err := prep.PrepareMigration(cmd.Context()); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not convert the migration table:\n%+v\n", err)
			return cmdx.FailSilently(cmd)
		}
	}

	status, err := p.MigrationStatus(cmd.Context())
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not get the migration status:\n%+v\n", errorsx.WithStack(err))
		return cmdx.FailSilently(cmd)
	}

	// Now we need to rollback the last `steps` migrations that have a status of "Applied":
	var count int
	var rollingBack int
	for i := len(status) - 1; i >= 0; i-- {
		if status[i].State == Applied {
			count++
			if steps > 0 && count <= steps {
				status[i].State = "Rollback"
				rollingBack++
			}
		}
	}

	// print migration status
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "The migration plan is as follows:")
	cmdx.PrintTable(cmd, status)

	if rollingBack < 1 {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "")
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "There are apparently no migrations to roll back.")
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Please provide the --steps argument with a value larger than 0.")
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "")
		return cmdx.FailSilently(cmd)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nThe SQL statements to be executed from top to bottom are:\n\n")

	for i := len(status) - 1; i >= 0; i-- {
		if status[i].State == "Rollback" {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ %s - %s ------------\n", status[i].Version, status[i].Name)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", status[i].ContentDown)
		}
	}

	if !flagx.MustGetBool(cmd, "yes") {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "To skip the next question use flag --yes (at your own risk).")
		if !cmdx.AskForConfirmation("Do you wish to execute this migration plan?", cmd.InOrStdin(), cmd.OutOrStdout()) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ WARNING ------------\n")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Migration aborted.")
			return nil
		}
	}

	// apply migrations
	if err := p.MigrateDown(cmd.Context(), rollingBack); err != nil {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ ERROR ------------\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not apply migrations:\n%+v\n", errorsx.WithStack(err))
		return cmdx.FailSilently(cmd)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "------------ SUCCESS ------------\n")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Successfully applied migrations!")
	return nil
}

func RegisterMigrateStatusFlags(cmd *cobra.Command) *cobra.Command {
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	cmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	cmd.Flags().Bool("block", false, "Block until all migrations have been applied")
	return cmd
}

func NewMigrateSQLStatusCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return RegisterMigrateStatusFlags(&cobra.Command{
		Use:   "status [database_url]",
		Short: "Display the current migration status",
		Long: `This command shows the current migration status for Ory {{ title .Root.Name }}.

You can use this command to check which migrations have been applied and which are pending.

To block until all migrations are applied, use the --block flag:

	DSN=... {{ .CommandPath }} -e --block`,
		Example: `See the current migration status:
	DSN=... {{ .CommandPath }} -e

Block until all migrations are applied:
	DSN=... {{ .CommandPath }} -e --block`,
		RunE: runE,
	})
}

func MigrateStatus(cmd *cobra.Command, p MigrationProvider) (err error) {
	block := flagx.MustGetBool(cmd, "block")
	ctx := cmd.Context()
	s, err := p.MigrationStatus(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not get migration status: %+v\n", err)
		return cmdx.FailSilently(cmd)
	}

	for block && s.HasPending() {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Waiting for migrations to finish...\n")
		for _, m := range s {
			if m.State == Pending {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), " - %s\n", m.Name)
			}
		}
		time.Sleep(time.Second)
		s, err = p.MigrationStatus(ctx)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not get migration status: %+v\n", err)
			return cmdx.FailSilently(cmd)
		}
	}

	cmdx.PrintTable(cmd, s)
	return nil
}
