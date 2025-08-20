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
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/ory/pop/v6"
	"github.com/ory/x/fsx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/popx"
	"github.com/ory/x/sqlcon"
)

//go:embed migrations/*.sql
var Migrations embed.FS

var SilenceMigrations = false

func (p *Persister) migrationBox() (*popx.MigrationBox, error) {
	logger := p.r.Logger()
	if SilenceMigrations {
		inner, _ := test.NewNullLogger()
		logger = logrusx.New("hydra", "", logrusx.UseLogger(inner))
	}
	return popx.NewMigrationBox(
		fsx.Merge(append([]fs.FS{Migrations}, p.extraMigrations...)...),
		p.conn, logger,
		popx.WithGoMigrations(p.goMigrations))
}

func (p *Persister) MigrationStatus(ctx context.Context) (popx.MigrationStatuses, error) {
	if p.mbs != nil {
		return p.mbs, nil
	}

	mb, err := p.migrationBox()
	if err != nil {
		return nil, err
	}
	status, err := mb.Status(ctx)
	if err != nil {
		return nil, err
	}

	if !status.HasPending() {
		p.mbs = status
	}

	return status, nil
}

func (p *Persister) MigrateDown(ctx context.Context, steps int) error {
	mb, err := p.migrationBox()
	if err != nil {
		return err
	}
	return mb.Down(ctx, steps)
}

func (p *Persister) MigrateUp(ctx context.Context) error {
	if err := p.migrateOldMigrationTables(); err != nil {
		return err
	}
	mb, err := p.migrationBox()
	if err != nil {
		return err
	}
	return mb.Up(ctx)
}

func (p *Persister) PrepareMigration(_ context.Context) error {
	return p.migrateOldMigrationTables()
}

type oldTableName string

const (
	clientMigrationTableName  oldTableName = "hydra_client_migration"
	jwkMigrationTableName     oldTableName = "hydra_jwk_migration"
	consentMigrationTableName oldTableName = "hydra_oauth2_authentication_consent_migration"
	oauth2MigrationTableName  oldTableName = "hydra_oauth2_migration"
)

// this type is copied from sql-migrate to remove the dependency
type OldMigrationRecord struct {
	ID        string    `db:"id"`
	AppliedAt time.Time `db:"applied_at"`
}

// this function is idempotent
func (p *Persister) migrateOldMigrationTables() error {
	if err := p.conn.RawQuery(fmt.Sprintf("SELECT * FROM %s", clientMigrationTableName)).Exec(); err != nil {
		// assume there are no old migration tables => done
		return nil
	}

	if err := pop.CreateSchemaMigrations(p.conn); err != nil {
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
		var migrations []OldMigrationRecord

		/* #nosec G201 table is static */
		if err := p.conn.RawQuery(fmt.Sprintf("SELECT * FROM %s", table)).All(&migrations); err != nil {
			if strings.Contains(err.Error(), string(table)) {
				continue
			}
			return err
		}

		// translate migrations
		for _, m := range migrations {
			// mark the migration as run for fizz
			// fizz standard version pattern: YYYYMMDDhhmmss
			migrationNumber, err := strconv.ParseInt(m.ID, 10, 0)
			if err != nil {
				return errors.WithStack(err)
			}

			/* #nosec G201 - i is static (0..3) and migrationNumber is from the database */
			if err := p.conn.RawQuery(
				fmt.Sprintf("INSERT INTO schema_migration (version) VALUES ('2019%02d%08d')", i+1, migrationNumber)).
				Exec(); err != nil {
				return errors.WithStack(err)
			}
		}

		// delete old migration table
		if err := p.conn.RawQuery(fmt.Sprintf("DROP TABLE %s", table)).Exec(); err != nil {
			return sqlcon.HandleError(err)
		}
	}

	return nil
}
