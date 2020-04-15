package migratest

import (
	"context"
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	for db, connect := range map[string]func(*testing.T) *pop.Connection{
		"postgres":  dockertest.ConnectToTestPostgreSQLPop,
		"mysql":     dockertest.ConnectToTestMySQLPop,
		"cockroach": dockertest.ConnectToTestCockroachDBPop,
	} {
		t.Run(fmt.Sprintf("database=%s", db), func(t *testing.T) {
			t.Parallel()
			c := connect(t)
			x.CleanSQLPop(t, c)
			require.NoError(t, os.Setenv("DSN", c.URL()))
			d := driver.NewDefaultDriver(logrus.New(), true, []string{}, "", "", "", false)
			var dbx *sqlx.DB
			require.NoError(t,
				dbal.Connect(c.URL(), logrus.New(), func() error {
					return nil
				}, func(db *sqlx.DB) error {
					dbx = db
					return nil
				}),
			)

			tm := NewTestMigrator(t, c, "../migrations", "./testdata")
			require.NoError(t, tm.Up())

			var lastClient *client.Client
			for i := 1; i <= 14; i++ {
				// skip cockroach assertions until migration 13
				if db == "cockroach" && i < 13 {
					continue
				}
				t.Run(fmt.Sprintf("case=client migration %d", i), func(t *testing.T) {
					expected := expectedClient(i)
					actual := &client.Client{}
					require.NoError(t, c.Find(actual, expected.ClientID))
					assertEqualClients(t, expected, actual)
					lastClient = actual
				})
			}

			for i := 1; i <= 4; i++ {
				// skip cockroach assertions until migration 4
				if db == "cockroach" && i < 4 {
					continue
				}
				t.Run(fmt.Sprintf("case=jwk migration %d", i), func(t *testing.T) {
					expected := expectedJWK(i)
					actual := &jwk.SQLData{}
					require.NoError(t, c.Where("pk = ?", expected.PK).First(actual))
					assertEqualJWKs(t, expected, actual)
				})
			}

			for i := 1; i <= 14; i++ {
				if db == "cockroach" && i < 12 {
					continue
				}
				t.Run(fmt.Sprintf("case=consent migration %d", i), func(t *testing.T) {
					ecr, elr, els, ehcr, ehlr, efols, elor := expectedConsent(i)

					acr, err := d.Registry().ConsentManager().GetConsentRequest(context.Background(), ecr.Challenge)
					require.NoError(t, err)
					assertEqualConsentRequests(t, ecr, acr)

					alr, err := d.Registry().ConsentManager().GetLoginRequest(context.Background(), elr.Challenge)
					require.NoError(t, err)
					assertEqualLoginRequests(t, elr, alr)

					als := &consent.LoginSession{}
					require.NoError(t, c.Find(als, els.ID))
					assertEqualLoginSessions(t, els, als)

					ahcr := &consent.HandledConsentRequest{}
					require.NoError(t, c.Q().Where("challenge = ?", ehcr.Challenge).First(ahcr))
					assertEqualHandledConsentRequests(t, ehcr, ahcr)

					ahlr := &consent.HandledLoginRequest{}
					require.NoError(t, c.Q().Where("challenge = ?", ehlr.Challenge).First(ahlr))
					assertEqualHandledLoginRequests(t, ehlr, ahlr)

					if efols != nil {
						afols, err := d.Registry().ConsentManager().GetForcedObfuscatedLoginSession(context.Background(), lastClient.ClientID, efols.SubjectObfuscated)
						require.NoError(t, err)
						assertEqualForcedObfucscatedLoginSessions(t, efols, afols)
					}

					if elor != nil {
						alor := &consent.LogoutRequest{}
						require.NoError(t, dbx.Get(alor, dbx.Rebind("select * from hydra_oauth2_logout_request where challenge = ?"), elor.Challenge))
						assertEqualLogoutRequests(t, elor, alor)
					}
				})
			}

			// TODO this is very stupid and should be replaced as soon the manager uses pop
			// necessary because the manager does not provide any way to access the data
			for i := 1; i <= 11; i++ {
				if db == "cockroach" && i < 9 {
					continue
				}

				tables := []string{"hydra_oauth2_access", "hydra_oauth2_refresh", "hydra_oauth2_code", "hydra_oauth2_oidc"}
				if i >= 3 {
					tables = append(tables, "hydra_oauth2_pkce")
				}
				ed, ebjti := expectedOauth2(i)
				ad := &oauth2.SQLData{}
				for _, table := range tables {
					require.NoError(t, dbx.Get(ad, dbx.Rebind(fmt.Sprintf("select * from %s where signature = ?", table)), ed.Signature), "table: %s\n%+v", table, ed)
					assertEqualOauth2Data(t, ed, ad)
				}

				if i >= 11 {
					abjti := &oauth2.BlacklistedJTI{}
					require.NoError(t, dbx.Get(abjti, dbx.Rebind("select * from hydra_oauth2_jti_blacklist where signature = ?"), ebjti.Signature))
					assertEqualOauth2BlacklistedJTIs(t, ebjti, abjti)
				}
			}

			x.CleanSQLPop(t, c)
			require.NoError(t, c.Close())
		})
	}
}
