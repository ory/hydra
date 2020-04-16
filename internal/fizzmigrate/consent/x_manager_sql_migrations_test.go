package consent_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
)

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

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{client.Migrations, consent.Migrations},
		migratest.MigrationSchemas{nil, dbal.FindMatchingTestMigrations("migrations/sql/tests/", consent.Migrations, consent.AssetNames(), consent.Asset)},
		x.CleanSQL, x.CleanSQL,
		func(t *testing.T, dbName string, db *sqlx.DB, sk, step, steps int) {
			if sk == 0 {
				t.Skip("Nothing to do...")
				return
			}

			t.Run(fmt.Sprintf("poll=%d", step), func(t *testing.T) {
				conf := internal.NewConfigurationWithDefaults()
				reg := internal.NewRegistrySQL(conf, db)

				kk := step + 1
				if dbName == "cockroach" {
					kk += 11
				}
				if kk <= 2 {
					t.Skip("Skipping the first two entries were deleted in migration 7.sql login_session_id is not defined")
					return
				}

				s := consent.NewSQLManager(db, reg)
				_, err := s.GetLoginRequest(context.TODO(), fmt.Sprintf("%d-challenge", kk))
				require.NoError(t, err, "%d-challenge", kk)
				_, err = s.GetRememberedLoginSession(context.TODO(), fmt.Sprintf("%d-login-session-id", kk))
				require.NoError(t, err, "%d-login-session-id", kk)
				_, err = s.GetConsentRequest(context.TODO(), fmt.Sprintf("%d-challenge", kk))
				require.NoError(t, err, "%d-challenge", kk)

				rs, err := s.FindGrantedAndRememberedConsentRequests(context.TODO(), fmt.Sprintf("%d-client", kk), fmt.Sprintf("%d-subject", kk))
				require.NoError(t, err, "%d-challenge %d-subject", kk, kk)
				require.True(t, len(rs) > 0)

				rs, err = s.FindSubjectsGrantedConsentRequests(context.TODO(), fmt.Sprintf("%d-subject", kk), 1, 0)
				require.NoError(t, err, "%d-challenge %d-subject", kk, kk)
				require.True(t, len(rs) > 0)

				if step > 1 {
					_, err = s.GetForcedObfuscatedLoginSession(context.TODO(), fmt.Sprintf("%d-client", kk), fmt.Sprintf("%d-obfuscated", kk))
					require.NoError(t, err, "%d-client", kk)
				}
			})
		},
	)
}
