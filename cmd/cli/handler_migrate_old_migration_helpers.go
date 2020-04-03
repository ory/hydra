package cli

import (
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/ory/x/sqlcon"
	migrate "github.com/rubenv/sql-migrate"
	"time"
)

type oldTableName string

const (
	MigrationTableName                     = "hydra_schema_migrations"
	clientMigrationTableName  oldTableName = "hydra_client_migration"
	jwkMigrationTableName     oldTableName = "hydra_jwk_migration"
	consentMigrationTableName oldTableName = "hydra_oauth2_authentication_consent_migration"
	oauth2MigrationTableName  oldTableName = "hydra_oauth2_migration"
)

func getMigrationRecords(c *pop.Connection, tableName oldTableName) ([]migrate.MigrationRecord, error) {
	var records []migrate.MigrationRecord

	/* #nosec G201 TableName is static */
	err := sqlcon.HandleError(
		c.RawQuery(fmt.Sprintf("SELECT * FROM %s", tableName)).
			Eager().
			All(&records))

	return records, err
}

func migrateOldMigrationTables(c *pop.Connection) error {
	var migrations []migrate.MigrationRecord
	err := c.Transaction(func(tx *pop.Connection) error {
		// in this order the migrations only depend on already done ones
		for _, table := range []oldTableName{clientMigrationTableName, jwkMigrationTableName, consentMigrationTableName, oauth2MigrationTableName} {
			ms, err := getMigrationRecords(tx, table)
			if err != nil {

				return err
			}
			migrations = append(migrations, ms...)
		}

		for i, m := range migrations {
			if m.AppliedAt.Before(time.Now()) {
				// the migration was run already -> set it as run for fizz
				// fizz standard version pattern: YYYYMMDDhhmmss
				err := tx.RawQuery(fmt.Sprintf("INSERT INTO %s (version) VALUES (2019%010d)", MigrationTableName, i+1)).Exec()
				if err != nil {
					return sqlcon.HandleError(err)
				}
			}
		}

		return nil
	})
	if err != nil {
		return sqlcon.HandleError(err)
	}

}

func hasOldMigrationTables() {

}
