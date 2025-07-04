// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/sirupsen/logrus"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/popx"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"
)

type MockPersistenceProvider struct {
	c  *pop.Connection
	mb *popx.MigrationBox
}

func (m *MockPersistenceProvider) MigrateDown(ctx context.Context, i int) error {
	return m.mb.Down(ctx, i)
}

func (m *MockPersistenceProvider) Connection(ctx context.Context) *pop.Connection {
	return m.c
}

func (m *MockPersistenceProvider) MigrationStatus(ctx context.Context) (popx.MigrationStatuses, error) {
	return m.mb.Status(ctx)
}

func (m *MockPersistenceProvider) MigrateUp(ctx context.Context) error {
	return m.mb.Up(ctx)
}

func NewMockPersistenceProvider(
	c *pop.Connection,
	mb *popx.MigrationBox,
) *MockPersistenceProvider {
	return &MockPersistenceProvider{c: c, mb: mb}
}

func TestMigrateSQLUp(t *testing.T) {
	ctx := context.Background()
	c, err := pop.NewConnection(&pop.ConnectionDetails{
		URL: dbal.NewSQLiteTestDatabase(t),
	})
	require.NoError(t, err)
	require.NoError(t, c.Open())

	migrator := popx.NewMigrator(c, logrusx.New("", "", logrusx.ForceLevel(logrus.DebugLevel)), nil, 0)
	mb, err := popx.NewMigrationBox(transactionalMigrations, migrator)
	require.NoError(t, err)

	p := NewMockPersistenceProvider(c, mb)
	newCmd := func() *cobra.Command {

		cmd := &cobra.Command{Use: ""}
		cmd.AddCommand(popx.RegisterMigrateSQLUpFlags(&cobra.Command{
			Use:  "up <database-url>",
			Args: cobra.RangeArgs(0, 1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return popx.MigrateSQLUp(cmd, p)
			}}))
		cmd.AddCommand(popx.RegisterMigrateSQLDownFlags(&cobra.Command{
			Use:  "down <database-url>",
			Args: cobra.RangeArgs(0, 1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return popx.MigrateSQLDown(cmd, p)
			}}))
		cmd.AddCommand(popx.RegisterMigrateStatusFlags(&cobra.Command{
			Use:  "status <database-url>",
			Args: cobra.RangeArgs(0, 1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return popx.MigrateStatus(cmd, p)
			}}))
		return cmd
	}

	run := func(t *testing.T, cmd *cobra.Command, stdIn io.Reader, args ...string) {
		t.Helper()
		stdout, stderr, err := cmdx.ExecCtx(ctx, newCmd(), stdIn, args...)
		require.NoError(t, err, stdout, stderr)

		cupaloy.New(
			cupaloy.CreateNewAutomatically(true),
			cupaloy.FailOnUpdate(true),
			cupaloy.SnapshotFileExtension(".txt"),
		).SnapshotT(t, fmt.Sprintf("stdout: %s\nstderr: %s", stdout, stderr))
	}

	t.Run("status pre", func(t *testing.T) {
		run(t, newCmd(), nil, "status")
	})

	t.Run("migrate up", func(t *testing.T) {
		run(t, newCmd(), nil, "up", "-y")
	})

	t.Run("status migrated", func(t *testing.T) {
		run(t, newCmd(), nil, "status")
	})

	t.Run("migrate down four steps", func(t *testing.T) {
		run(t, newCmd(), nil, "down", "-y", "--steps", "4")
	})

	t.Run("status two steps rolled back", func(t *testing.T) {
		run(t, newCmd(), nil, "status")
	})

	t.Run("migrate down but no steps", func(t *testing.T) {
		stdout, stderr, err := cmdx.ExecCtx(ctx, newCmd(), nil, "down", "-y")
		require.Error(t, err)

		cupaloy.New(
			cupaloy.CreateNewAutomatically(true),
			cupaloy.FailOnUpdate(true),
			cupaloy.SnapshotFileExtension(".txt"),
		).SnapshotT(t, fmt.Sprintf("stdout: %s\nstderr: %s", stdout, stderr))
	})

	t.Run("migrate down but do not confirm", func(t *testing.T) {
		run(t, newCmd(), bytes.NewBufferString("n\n"), "down", "--steps", "2")
	})

	t.Run("migrate down two steps", func(t *testing.T) {
		run(t, newCmd(), bytes.NewBufferString("y\n"), "down", "--steps", "2")
	})

	t.Run("status two versions rolled back", func(t *testing.T) {
		run(t, newCmd(), nil, "status")
	})

	t.Run("migrate rollbacks up again", func(t *testing.T) {
		run(t, newCmd(), bytes.NewBufferString("y\n"), "up")
	})

	t.Run("final status", func(t *testing.T) {
		run(t, newCmd(), nil, "status")
	})

	t.Run("migrate rollbacks up without confirm", func(t *testing.T) {
		run(t, newCmd(), bytes.NewBufferString("n\n"), "up")
	})
}
