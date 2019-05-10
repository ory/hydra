package oauth2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{client.Migrations, consent.Migrations, oauth2.Migrations},
		migratest.MigrationSchemas{nil, nil, dbal.FindMatchingTestMigrations("migrations/sql/tests/", oauth2.Migrations, oauth2.AssetNames(), oauth2.Asset)},
		x.CleanSQL,
		x.CleanSQL,
		func(t *testing.T, dbName string, db *sqlx.DB, m, k, steps int) {
			t.Run(fmt.Sprintf("poll=%d", k), func(t *testing.T) {
				conf := internal.NewConfigurationWithDefaults()
				reg := internal.NewRegistrySQL(conf, db)

				if m != 2 {
					t.Skip("Skipping polling unless it's the last migration schema")
					return
				}

				s := reg.OAuth2Storage().(*oauth2.FositeSQLStore)
				if dbName == "cockroach" {
					k += 8
				}
				sig := fmt.Sprintf("%d-sig", k+1)

				if k < 8 {
					// With migration 8, all previous test data has been removed because the client is non-existent.
					_, err := s.GetAccessTokenSession(context.Background(), sig, oauth2.NewSession(""))
					require.Error(t, err)
					return
				}
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
