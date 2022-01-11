package migratest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/configx"

	"github.com/ory/hydra/persistence/sql"

	"github.com/ory/x/popx"

	"github.com/ory/x/sqlcon/dockertest"

	"github.com/ory/hydra/driver/config"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
)

func TestMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	connections := make(map[string]*pop.Connection, 3)
	dockertest.Parallel([]func(){
		func() {
			connections["postgres"] = dockertest.ConnectToTestPostgreSQLPop(t)
		},
		func() {
			connections["mysql"] = dockertest.ConnectToTestMySQLPop(t)
		},
		func() {
			connections["cockroach"] = dockertest.ConnectToTestCockroachDBPop(t)
		},
	})

	for db, c := range connections {
		t.Run(fmt.Sprintf("database=%s", db), func(t *testing.T) {
			x.CleanSQLPop(t, c)
			url := c.URL()

			// workaround for https://github.com/gobuffalo/pop/issues/538
			if db == "mysql" {
				url = "mysql://" + url
			}

			d := driver.New(
				context.Background(),
				driver.WithOptions(configx.WithValue(config.KeyDSN, url)),
				driver.DisablePreloading(),
				driver.DisableValidation(),
			)

			tm := popx.NewTestMigrator(t, c, os.DirFS("../migrations"), os.DirFS("./testdata"), d.Logger())
			require.NoError(t, tm.Up(context.Background()))

			var lastClient *client.Client
			for i := 1; i <= 14; i++ {
				// skip cockroach assertions until migration 13
				if db == "cockroach" && i < 13 {
					continue
				}
				t.Run(fmt.Sprintf("case=client migration %d", i), func(t *testing.T) {
					expected := expectedClient(i)
					actual := &client.Client{}
					require.NoError(t, c.Where("id = ?", expected.OutfacingID).First(actual))
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
					require.NoError(t, c.Where("pk = ?", expected.ID).First(actual))
					assertEqualJWKs(t, expected, actual)
				})
			}

			for i := 1; i <= 14; i++ {
				if db == "cockroach" && i < 12 {
					continue
				}
				t.Run(fmt.Sprintf("case=consent migration %d", i), func(t *testing.T) {
					ecr, elr, els, ehcr, ehlr, efols, elor := expectedConsent(i)

					acr, err := d.ConsentManager().GetConsentRequest(context.Background(), ecr.ID)
					require.NoError(t, err)
					assertEqualConsentRequests(t, ecr, acr)

					alr, err := d.ConsentManager().GetLoginRequest(context.Background(), elr.ID)
					require.NoError(t, err)
					assertEqualLoginRequests(t, elr, alr)

					als := &consent.LoginSession{}
					require.NoError(t, c.Find(als, els.ID))
					assertEqualLoginSessions(t, els, als)

					ahcr := &consent.HandledConsentRequest{}
					require.NoError(t, c.Q().Where("challenge = ?", ehcr.ID).First(ahcr))
					require.NoError(t, ehcr.AfterFind(c))
					assertEqualHandledConsentRequests(t, ehcr, ahcr)

					ahlr := &consent.HandledLoginRequest{}
					require.NoError(t, c.Q().Where("challenge = ?", ehlr.ID).First(ahlr))
					assertEqualHandledLoginRequests(t, ehlr, ahlr)

					if efols != nil {
						afols, err := d.ConsentManager().GetForcedObfuscatedLoginSession(context.Background(), lastClient.OutfacingID, efols.SubjectObfuscated)
						require.NoError(t, err)
						assertEqualForcedObfucscatedLoginSessions(t, efols, afols)
					}

					if elor != nil {
						alor := &consent.LogoutRequest{}
						require.NoError(t, d.Persister().Connection(context.Background()).RawQuery("select * from hydra_oauth2_logout_request where challenge = ?", elor.ID).First(alor))
						alor.Client = nil
						assertEqualLogoutRequests(t, elor, alor)
					}
				})
			}

			// TODO https://github.com/ory/hydra/issues/1815
			// this is very stupid and should be replaced as soon the manager uses pop
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
				ad := &sql.OAuth2RequestSQL{}
				for _, table := range tables {
					require.NoError(t, d.Persister().Connection(context.Background()).RawQuery(fmt.Sprintf("select * from %s where signature = ?", table), ed.ID).First(ad))
					assertEqualOauth2Data(t, ed, ad)
				}

				if i >= 11 {
					abjti := &oauth2.BlacklistedJTI{}
					require.NoError(t, c.Where("signature = ?", ebjti.ID).First(abjti))
					assertEqualOauth2BlacklistedJTIs(t, ebjti, abjti)
				}
			}

			x.CleanSQLPop(t, c)
			require.NoError(t, c.Close())
		})
	}
}
