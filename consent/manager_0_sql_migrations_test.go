package consent

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

var createMigrations = map[string]*migrate.PackrMigrationSource{
	dbal.DriverMySQL:      dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/tests"}),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/tests"}),
}

func TestMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(client.Migrations[dbal.DriverMySQL].Box.List()) == len(client.Migrations[dbal.DriverPostgreSQL].Box.List()))

	var clients []client.Client
	for k := range client.Migrations[dbal.DriverMySQL].Box.List() {
		clients = append(clients, client.Client{ClientID: fmt.Sprintf("%d-client", k+1)})
	}

	var cm = &client.MemoryManager{Clients: clients}

	var clean = func(t *testing.T, db *sqlx.DB) {
		_, err := db.Exec("DROP TABLE IF EXISTS hydra_oauth2_consent_request")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_request")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_session")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_consent_request_handled")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_request_handled")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_obfuscated_authentication_session")
		require.NoError(t, err)
	}

	migratest.RunPackrMigrationTests(
		t,
		migrations,
		createMigrations,
		clean, clean,
		func(t *testing.T, db *sqlx.DB, k int) {
			t.Run(fmt.Sprintf("poll=%d", k), func(t *testing.T) {
				kk := k + 1
				s := &SQLManager{db: db, c: cm}
				_, err := s.GetAuthenticationRequest(context.TODO(), fmt.Sprintf("%d-challenge", kk))
				require.NoError(t, err)
				_, err = s.GetAuthenticationSession(context.TODO(), fmt.Sprintf("%d-auth", kk))
				require.NoError(t, err)
				_, err = s.GetConsentRequest(context.TODO(), fmt.Sprintf("%d-challenge", kk))
				require.NoError(t, err)
				if k > 1 {
					_, err = s.GetForcedObfuscatedAuthenticationSession(context.TODO(), fmt.Sprintf("%d-client", kk), fmt.Sprintf("%d-obfuscated", kk))
					require.NoError(t, err)
				}
			})
		},
	)
}
