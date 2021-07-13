// +build legacy_migration_test

package jwk_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/hydra/internal/fizzmigrate/client"
	migrateJWK "github.com/ory/hydra/internal/fizzmigrate/jwk"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(client.Migrations[dbal.DriverMySQL].Box.List()) == len(client.Migrations[dbal.DriverPostgreSQL].Box.List()))

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{migrateJWK.Migrations},
		migratest.MigrationSchemas{dbal.FindMatchingTestMigrations("migrations/sql/tests/", migrateJWK.Migrations, migrateJWK.AssetNames(), migrateJWK.Asset)},
		x.CleanSQL,
		x.CleanSQL,
		func(t *testing.T, dbName string, db *sqlx.DB, k, m, steps int) {
			t.Run(fmt.Sprintf("poll=%d", k), func(t *testing.T) {
				if dbName == "cockroach" {
					k += 3
				}
				conf := internal.NewConfigurationWithDefaults()
				reg := internal.NewRegistrySQLFromDB(conf, db)

				sid := fmt.Sprintf("%d-sid", k+1)
				m := jwk.NewSQLManager(db, reg)
				_, err := m.GetKeySet(context.TODO(), sid)
				require.Error(t, err, "malformed ciphertext")
			})
		},
	)
}
