package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/ory/x/sqlcon"
	"github.com/pkg/errors"
)

type oldTableName string

const (
	MigrationTableName                     = "hydra_schema_migrations"
	clientMigrationTableName  oldTableName = "hydra_client_migration"
	jwkMigrationTableName     oldTableName = "hydra_jwk_migration"
	consentMigrationTableName oldTableName = "hydra_oauth2_authentication_consent_migration"
	oauth2MigrationTableName  oldTableName = "hydra_oauth2_migration"
)

// this type is copied from sql-migrate to remove the dependency
type OldMigrationRecord struct {
	Id        string    `db:"id"`
	AppliedAt time.Time `db:"applied_at"`
}

func getMigrationRecords(c *pop.Connection, tableName oldTableName) ([]OldMigrationRecord, error) {
	var records []OldMigrationRecord

	/* #nosec G201 TableName is static */
	err := sqlcon.HandleError(
		c.RawQuery(fmt.Sprintf("SELECT * FROM %s", tableName)).
			Eager().
			All(&records))

	return records, err
}

func migrateOldMigrationTables(c *pop.Connection) error {
	if err := c.RawQuery(fmt.Sprintf("SELECT * FROM %s", clientMigrationTableName)).Exec(); err != nil {
		// assume there are no old migration tables => done
		return nil
	}

	return sqlcon.HandleError(c.Transaction(func(tx *pop.Connection) error {
		if err := pop.CreateSchemaMigrations(c); err != nil {
			return errors.WithStack(err)
		}

		// in this order the migrations only depend on already done ones
		for i, table := range []oldTableName{clientMigrationTableName, jwkMigrationTableName, consentMigrationTableName, oauth2MigrationTableName} {
			// get old migrations
			migrations, err := getMigrationRecords(tx, table)
			if err != nil {
				return err
			}

			// translate migrations
			for _, m := range migrations {
				if m.AppliedAt.Before(time.Now()) {
					// the migration was run already -> set it as run for fizz
					// fizz standard version pattern: YYYYMMDDhhmmss
					migrationNumber, err := strconv.ParseInt(m.Id, 10, 0)
					if err != nil {
						return errors.WithStack(err)
					}

					/* #nosec G201 - i is static (0..3) and migrationNumber is from the database */
					if err := tx.RawQuery(
						fmt.Sprintf("INSERT INTO schema_migration (version) VALUES (2019%02d%08d)", i+1, migrationNumber)).
						Exec(); err != nil {
						return sqlcon.HandleError(err)
					}
				}
			}

			// delete old migration table
			if err := tx.RawQuery(fmt.Sprintf("DROP TABLE %s", table)).Exec(); err != nil {
				return sqlcon.HandleError(err)
			}
		}

		return nil
	}))
}
