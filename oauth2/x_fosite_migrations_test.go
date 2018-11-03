package oauth2_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

var createMigrations = map[string]*migrate.PackrMigrationSource{
	dbal.DriverMySQL:      dbal.NewMustPackerMigrationSource(logrus.New(), oauth2.AssetNames(), oauth2.Asset, []string{"migrations/sql/tests"}),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrus.New(), oauth2.AssetNames(), oauth2.Asset, []string{"migrations/sql/tests"}),
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

	var cm = &client.MemoryManager{Clients: clients}

	var clean = func(t *testing.T, db *sqlx.DB) {
		_, err := db.Exec("DROP TABLE IF EXISTS hydra_oauth2_access")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_refresh")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_code")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_oidc")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_pkce")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_oauth2_migration")
		require.NoError(t, err)
	}

	migratest.RunPackrMigrationTests(
		t,
		oauth2.Migrations,
		createMigrations,
		clean, clean,
		func(t *testing.T, db *sqlx.DB, k int) {
			t.Run(fmt.Sprintf("poll=%d", k), func(t *testing.T) {
				sig := fmt.Sprintf("%d-sig", k+1)
				s := oauth2.NewFositeSQLStore(cm, db, logrus.New(), time.Minute, false)
				_, err := s.GetAccessTokenSession(context.Background(), sig, oauth2.NewSession(""))
				require.NoError(t, err)
				_, err = s.GetRefreshTokenSession(context.Background(), sig, oauth2.NewSession(""))
				require.NoError(t, err)
				_, err = s.GetAuthorizeCodeSession(context.Background(), sig, oauth2.NewSession(""))
				require.NoError(t, err)
				_, err = s.GetOpenIDConnectSession(context.Background(), sig, &fosite.Request{Session: oauth2.NewSession("")})
				require.NoError(t, err)
				if k > 2 {
					_, err = s.GetPKCERequestSession(context.Background(), sig, oauth2.NewSession(""))
					require.NoError(t, err)
				}
			})
		},
	)
}
