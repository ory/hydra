// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/fsx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/popx"
	"github.com/ory/x/sqlcon"
)

//go:embed migrations/*.sql
var Migrations embed.FS

var SilenceMigrations = false

type (
	MigrationManager struct {
		d               migrationDependencies
		conn            *pop.Connection
		extraMigrations []fs.FS
		goMigrations    []popx.Migration

		// cached values
		mb  *popx.MigrationBox
		mbs popx.MigrationStatuses
	}
	migrationDependencies interface {
		x.RegistryLogger
	}
)

func NewMigrationManager(c *pop.Connection, d migrationDependencies, extraMigrations []fs.FS, goMigrations []popx.Migration) *MigrationManager {
	return &MigrationManager{
		d:               d,
		conn:            c,
		extraMigrations: extraMigrations,
		goMigrations:    goMigrations,
	}
}

func (m *MigrationManager) migrationBox() (_ *popx.MigrationBox, err error) {
	if m.mb == nil {
		logger := m.d.Logger()
		if SilenceMigrations {
			inner, _ := test.NewNullLogger()
			logger = logrusx.New("hydra", "", logrusx.UseLogger(inner))
		}
		m.mb, err = popx.NewMigrationBox(
			fsx.Merge(append([]fs.FS{Migrations}, m.extraMigrations...)...),
			m.conn, logger,
			popx.WithGoMigrations(m.goMigrations))
		if err != nil {
			return nil, err
		}
	}
	return m.mb, nil
}

func (m *MigrationManager) MigrationStatus(ctx context.Context) (popx.MigrationStatuses, error) {
	if m.mbs != nil {
		return m.mbs, nil
	}

	mb, err := m.migrationBox()
	if err != nil {
		return nil, err
	}
	status, err := mb.Status(ctx)
	if err != nil {
		return nil, err
	}

	if !status.HasPending() {
		m.mbs = status
	}

	return status, nil
}

func (m *MigrationManager) MigrateDown(ctx context.Context, steps int) error {
	mb, err := m.migrationBox()
	if err != nil {
		return err
	}
	return mb.Down(ctx, steps)
}

func (m *MigrationManager) MigrateUp(ctx context.Context) error {
	if err := m.migrateOldMigrationTables(); err != nil {
		return err
	}
	mb, err := m.migrationBox()
	if err != nil {
		return err
	}
	return mb.Up(ctx)
}

func (m *MigrationManager) PrepareMigration(_ context.Context) error {
	return m.migrateOldMigrationTables()
}

type oldTableName string

const (
	clientMigrationTableName  oldTableName = "hydra_client_migration"
	jwkMigrationTableName     oldTableName = "hydra_jwk_migration"
	consentMigrationTableName oldTableName = "hydra_oauth2_authentication_consent_migration"
	oauth2MigrationTableName  oldTableName = "hydra_oauth2_migration"
)

// this function is idempotent
func (m *MigrationManager) migrateOldMigrationTables() error {
	if err := m.conn.RawQuery(fmt.Sprintf("SELECT * FROM %s", clientMigrationTableName)).Exec(); err != nil {
		// assume there are no old migration tables => done
		return nil
	}

	if err := pop.CreateSchemaMigrations(m.conn); err != nil {
		return errors.WithStack(err)
	}

	// in this order the migrations only depend on already done ones
	for i, table := range []oldTableName{clientMigrationTableName, jwkMigrationTableName, consentMigrationTableName, oauth2MigrationTableName} {
		// If table does not exist, we will skip it. Previously, we created a stub table here which
		// caused the cached statements to fail, see:
		//
		// https://github.com/flynn/flynn/pull/2306/files
		// https://github.com/jackc/pgx/issues/110
		// https://github.com/flynn/flynn/issues/2235
		// get old migrations
		var migrations []string

		/* #nosec G201 table is static */
		if err := m.conn.RawQuery(fmt.Sprintf("SELECT id FROM %s", table)).All(&migrations); err != nil {
			if strings.Contains(err.Error(), string(table)) {
				continue
			}
			return err
		}

		// translate migrations
		for _, migration := range migrations {
			// mark the migration as run for fizz
			// fizz standard version pattern: YYYYMMDDhhmmss
			migrationNumber, err := strconv.ParseInt(migration, 10, 0)
			if err != nil {
				return errors.WithStack(err)
			}

			/* #nosec G201 - i is static (0..3) and migrationNumber is from the database */
			if err := m.conn.RawQuery(
				fmt.Sprintf("INSERT INTO schema_migration (version) VALUES ('2019%02d%08d')", i+1, migrationNumber)).
				Exec(); err != nil {
				return errors.WithStack(err)
			}
		}

		// delete old migration table
		if err := m.conn.RawQuery(fmt.Sprintf("DROP TABLE %s", table)).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}
	}

	return nil
}
