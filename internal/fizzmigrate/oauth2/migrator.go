package oauth2

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/ory/x/errorsx"

	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
)

var Migrations = map[string]*dbal.PackrMigrationSource{
	dbal.DriverMySQL: dbal.NewMustPackerMigrationSource(logrusx.New("", ""), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/mysql",
	}, true),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrusx.New("", ""), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/postgres",
	}, true),
	dbal.DriverCockroachDB: dbal.NewMustPackerMigrationSource(logrusx.New("", ""), AssetNames(), Asset, []string{
		"migrations/sql/cockroach",
	}, true),
}

type migrator struct {
	DB *sqlx.DB
}

func NewMigrator(db *sqlx.DB) *migrator {
	return &migrator{db}
}

func (s *migrator) PlanMigration(dbName string) ([]*migrate.PlannedMigration, error) {
	migrate.SetTable("hydra_oauth2_migration")
	plan, _, err := migrate.PlanMigration(s.DB.DB, dbal.Canonicalize(s.DB.DriverName()), Migrations[dbName], migrate.Up, 0)
	return plan, errorsx.WithStack(err)
}

func (s *migrator) CreateSchemas(dbName string) (int, error) {
	migrate.SetTable("hydra_oauth2_migration")
	n, err := migrate.Exec(s.DB.DB, dbal.Canonicalize(s.DB.DriverName()), Migrations[dbName], migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *migrator) CreateMaxSchemas(dbName string, steps int) (int, error) {
	migrate.SetTable("hydra_oauth2_migration")
	n, err := migrate.ExecMax(m.DB.DB, dbal.Canonicalize(m.DB.DriverName()), Migrations[dbName], migrate.Up, steps)
	return n, errorsx.WithStack(err)
}
