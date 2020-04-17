package consent

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"

	"github.com/ory/x/dbal"
)

var Migrations = map[string]*dbal.PackrMigrationSource{
	"mysql": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/mysql",
	}, true),
	"postgres": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/postgres",
	}, true),
	"cockroach": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/cockroach",
	}, true),
}

type migrator struct {
	DB *sqlx.DB
}

func NewMigrator(db *sqlx.DB) *migrator {
	return &migrator{db}
}

func (m *migrator) PlanMigration(dbName string) ([]*migrate.PlannedMigration, error) {
	migrate.SetTable("hydra_oauth2_authentication_consent_migration")
	plan, _, err := migrate.PlanMigration(m.DB.DB, dbal.Canonicalize(m.DB.DriverName()), Migrations[dbName], migrate.Up, 0)
	return plan, errors.WithStack(err)
}

func (m *migrator) CreateSchemas(dbName string) (int, error) {
	migrate.SetTable("hydra_oauth2_authentication_consent_migration")
	n, err := migrate.Exec(m.DB.DB, dbal.Canonicalize(m.DB.DriverName()), Migrations[dbName], migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}
