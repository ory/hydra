// +build legacy_migration_test

package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/hydra/client"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(Migrations[dbal.DriverMySQL].Box.List()) == len(Migrations[dbal.DriverPostgreSQL].Box.List()))

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{Migrations},
		migratest.MigrationSchemas{dbal.FindMatchingTestMigrations("migrations/sql/tests/", Migrations, AssetNames(), Asset)},
		x.CleanSQL, x.CleanSQL,
		func(t *testing.T, dbName string, db *sqlx.DB, _, step, steps int) {
			if dbName == "cockroach" {
				step += 12
			}
			id := fmt.Sprintf("%d-data", step+1)
			t.Run("poll="+id, func(t *testing.T) {
				s := client.NewSQLManager(db, nil)
				c, err := s.GetConcreteClient(context.TODO(), id)
				require.NoError(t, err)
				assert.EqualValues(t, c.GetID(), id)
			})
		},
	)
}
