// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx_test

import (
	"context"
	"embed"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"

	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	. "github.com/ory/x/popx"
)

//go:embed stub/migrations/transactional/*.sql
var transactionalMigrations embed.FS

func TestMigratorUpgradingFromStart(t *testing.T) {
	ctx := context.Background()

	c, err := pop.NewConnection(&pop.ConnectionDetails{
		URL: dbal.NewSQLiteTestDatabase(t),
	})
	require.NoError(t, err)
	require.NoError(t, c.Open())

	l := logrusx.New("", "", logrusx.ForceLevel(logrus.DebugLevel))
	transactional, err := NewMigrationBox(transactionalMigrations, NewMigrator(c, l, nil, 0))
	require.NoError(t, err)
	status, err := transactional.Status(ctx)
	require.NoError(t, err)
	assert.True(t, status.HasPending())

	applied, err := transactional.UpTo(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, applied)

	status, err = transactional.Status(ctx)
	require.NoError(t, err)
	assert.True(t, status.HasPending())
	assert.Equal(t, Applied, status[0].State)
	assert.Equal(t, Pending, status[1].State)

	require.NoError(t, transactional.Up(ctx))

	status, err = transactional.Status(ctx)
	require.NoError(t, err)
	assert.False(t, status.HasPending())

	// Are all the tables here?
	var rows []string
	require.NoError(t, c.RawQuery("SELECT name FROM sqlite_master WHERE type='table'").All(&rows))

	assert.ElementsMatch(t, rows, []string{"schema_migration", "identities", "identity_credential_types",
		"identity_credentials", "identity_credential_identifiers", "selfservice_login_flows", "selfservice_login_flow_methods",
		"selfservice_registration_flows", "selfservice_registration_flow_methods", "selfservice_errors", "courier_messages",
		"selfservice_settings_flow_methods", "continuity_containers", "identity_recovery_addresses",
		"selfservice_recovery_flows", "selfservice_recovery_flow_methods", "selfservice_settings_flows", "sessions",
		"selfservice_verification_flow_methods", "selfservice_verification_flows", "identity_verification_tokens",
		"identity_recovery_tokens", "identity_verifiable_addresses"})

	require.NoError(t, transactional.Down(ctx, -1))
}

func TestMigratorSanitizeMigrationTableName(t *testing.T) {
	ctx := context.Background()

	c, err := pop.NewConnection(&pop.ConnectionDetails{
		URL: dbal.NewSQLiteTestDatabase(t) + "&migration_table_name=injection--",
	})
	require.NoError(t, err)
	require.NoError(t, c.Open())

	l := logrusx.New("", "", logrusx.ForceLevel(logrus.DebugLevel))
	transactional, err := NewMigrationBox(transactionalMigrations, NewMigrator(c, l, nil, 0))
	require.NoError(t, err)
	status, err := transactional.Status(ctx)
	require.NoError(t, err)
	require.True(t, status.HasPending())

	require.NoError(t, transactional.Up(ctx))

	status, err = transactional.Status(ctx)
	require.NoError(t, err)
	require.False(t, status.HasPending())

	require.NoError(t, transactional.Down(ctx, -1))
}
