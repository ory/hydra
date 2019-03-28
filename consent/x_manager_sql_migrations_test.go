package consent_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/hydra/internal"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

var createMigrations = map[string]*dbal.PackrMigrationSource{
	dbal.DriverMySQL:      dbal.NewMustPackerMigrationSource(logrus.New(), consent.AssetNames(), consent.Asset, []string{"migrations/sql/tests"}, true),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrus.New(), consent.AssetNames(), consent.Asset, []string{"migrations/sql/tests"}, true),
}

func cleanDB(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS hydra_oauth2_access")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_refresh")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_code")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_oidc")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_pkce")
	t.Logf("Unable to execute clean up query: %s", err)

	// hydra_oauth2_consent_request_handled depends on hydra_oauth2_consent_request
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_consent_request_handled")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_consent_request")
	t.Logf("Unable to execute clean up query: %s", err)

	// hydra_oauth2_authentication_request_handled depends on hydra_oauth2_authentication_request
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_request_handled")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_request")
	t.Logf("Unable to execute clean up query: %s", err)

	// everything depends on hydra_oauth2_authentication_session
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_session")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_obfuscated_authentication_session")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_client")
	t.Logf("Unable to execute clean up query: %s", err)

	// clean up migration tables
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_authentication_consent_migration")
	t.Logf("Unable to execute clean up query: %s", err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_client_migration")
	t.Logf("Unable to execute clean up query: %s", err)

	_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_migration")
	t.Logf("Unable to execute clean up query: %s", err)
}

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(client.Migrations[dbal.DriverMySQL].Box.List()) == len(client.Migrations[dbal.DriverPostgreSQL].Box.List()))

	var clients []client.Client
	for k := range client.Migrations[dbal.DriverMySQL].Box.List() {
		clients = append(clients, client.Client{ClientID: fmt.Sprintf("%d-client", k+1)})
	}

	var clean = cleanDB

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{client.Migrations, consent.Migrations},
		migratest.MigrationSchemas{nil, createMigrations},
		clean, clean,
		func(t *testing.T, db *sqlx.DB, sk, step, steps int) {
			if sk == 0 {
				t.Skip("Nothing to do...")
				return
			}

			t.Run(fmt.Sprintf("poll=%d", step), func(t *testing.T) {
				conf := internal.NewConfigurationWithDefaults()
				reg := internal.NewRegistrySQL(conf, db)

				kk := step + 1
				if kk <= 2 {
					t.Skip("Skipping the first two entries were deleted in migration 7.sql login_session_id is not defined")
					return
				}

				s := consent.NewSQLManager(db, reg)
				_, err := s.GetAuthenticationRequest(context.TODO(), fmt.Sprintf("%d-challenge", kk))
				require.NoError(t, err, "%d-challenge", kk)
				_, err = s.GetAuthenticationSession(context.TODO(), fmt.Sprintf("%d-login-session-id", kk))
				require.NoError(t, err, "%d-login-session-id", kk)
				_, err = s.GetConsentRequest(context.TODO(), fmt.Sprintf("%d-challenge", kk))
				require.NoError(t, err, "%d-challenge", kk)
				if step > 1 {
					_, err = s.GetForcedObfuscatedAuthenticationSession(context.TODO(), fmt.Sprintf("%d-client", kk), fmt.Sprintf("%d-obfuscated", kk))
					require.NoError(t, err, "%d-client", kk)
				}
			})
		},
	)
}
