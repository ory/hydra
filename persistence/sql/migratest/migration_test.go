package migratest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/gofrs/uuid"
	"github.com/instana/testify/assert"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/configx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/flow"
	testhelpersuuid "github.com/ory/hydra/internal/testhelpers/uuid"
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

func snapshotFor(paths ...string) *cupaloy.Config {
	return cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(true),
		cupaloy.SnapshotFileExtension(".json"),
		cupaloy.SnapshotSubdirectory(filepath.Join(paths...)),
	)
}

func CompareWithFixture(t *testing.T, actual interface{}, prefix string, id string) {
	s := snapshotFor("fixtures", prefix)
	actualJSON, err := json.MarshalIndent(actual, "", "  ")
	require.NoError(t, err)
	assert.NoError(t, s.SnapshotWithName(id, actualJSON))
}

func TestMigrations2(t *testing.T) {
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

	var test = func(db string, c *pop.Connection) func(t *testing.T) {
		return func(t *testing.T) {
			ctx := context.Background()
			x.CleanSQLPop(t, c)
			url := c.URL()

			// workaround for https://github.com/gobuffalo/pop/issues/538
			if db == "mysql" {
				url = "mysql://" + url
			} else if db == "sqlite" {
				url = "sqlite://" + url
			}

			d := driver.New(
				ctx,
				driver.WithOptions(configx.WithValue(config.KeyDSN, url)),
				driver.DisablePreloading(),
				driver.DisableValidation(),
			)

			tm := popx.NewTestMigrator(t, c, os.DirFS("../migrations"), os.DirFS("./testdata"), d.Logger())
			require.NoError(t, tm.Up(ctx))

			t.Run("suite=fixtures", func(t *testing.T) {
				t.Run("case=hydra_client", func(t *testing.T) {
					cs := []client.Client{}
					require.NoError(t, c.All(&cs))
					if db == "cockroach" {
						require.Equal(t, 4, len(cs))
					} else {
						require.Equal(t, 16, len(cs))
					}
					for _, c := range cs {
						require.False(t, c.CreatedAt.IsZero())
						require.False(t, c.UpdatedAt.IsZero())
						c.CreatedAt = time.Time{} // Some CreatedAt and UpdatedAt values are generated during migrations so we zero them in the fixtures
						c.UpdatedAt = time.Time{}
						CompareWithFixture(t, c, "hydra_client", c.OutfacingID)
						testhelpersuuid.AssertUUID(t, &c.ID)
					}
				})

				t.Run("case=hydra_jwk", func(t *testing.T) {
					js := []jwk.SQLData{}
					require.NoError(t, c.All(&js))
					if db == "cockroach" {
						require.Equal(t, 3, len(js))
					} else {
						require.Equal(t, 6, len(js))
					}
					for _, j := range js {
						testhelpersuuid.AssertUUID(t, &j.ID)
						j.ID = uuid.Nil // Some IDs are generated at migration time so we zero them in the fixtures
						require.False(t, j.CreatedAt.IsZero())
						j.CreatedAt = time.Time{}
						CompareWithFixture(t, j, "hydra_jwk", j.KID)
					}
				})

				flows := []flow.Flow{}
				require.NoError(t, c.All(&flows))
				if db == "cockroach" {
					require.Equal(t, 3, len(flows))
				} else {
					require.Equal(t, 14, len(flows))
				}

				t.Run("case=hydra_oauth2_flow", func(t *testing.T) {
					for _, f := range flows {
						fixturizeFlow(t, &f)
						CompareWithFixture(t, f, "hydra_oauth2_flow", f.ID)
					}
				})

				t.Run("case=hydra_oauth2_authentication_session", func(t *testing.T) {
					ss := []consent.LoginSession{}
					c.All(&ss)
					if db == "cockroach" {
						require.Equal(t, 3, len(ss))
					} else {
						require.Equal(t, 14, len(ss))
					}

					for _, s := range ss {
						s.AuthenticatedAt = sqlxx.NullTime(time.Time{})
						CompareWithFixture(t, s, "hydra_oauth2_authentication_session", s.ID)
					}
				})

				t.Run("case=hydra_oauth2_obfuscated_authentication_session", func(t *testing.T) {
					ss := []consent.ForcedObfuscatedLoginSession{}
					c.All(&ss)
					if db == "cockroach" {
						require.Equal(t, 3, len(ss))
					} else {
						require.Equal(t, 13, len(ss))
					}

					for _, s := range ss {
						CompareWithFixture(t, s, "hydra_oauth2_obfuscated_authentication_session", fmt.Sprintf("%s_%s", s.Subject, s.ClientID))
					}
				})

				t.Run("case=hydra_oauth2_logout_request", func(t *testing.T) {
					lrs := []consent.LogoutRequest{}
					c.All(&lrs)
					if db == "cockroach" {
						require.Equal(t, 3, len(lrs))
					} else {
						require.Equal(t, 6, len(lrs))
					}

					for _, s := range lrs {
						s.Client = nil
						CompareWithFixture(t, s, "hydra_oauth2_logout_request", s.ID)
					}
				})

				t.Run("case=hydra_oauth2_jti_blacklist", func(t *testing.T) {
					bjtis := []oauth2.BlacklistedJTI{}
					c.All(&bjtis)
					require.Equal(t, 1, len(bjtis))
					for _, bjti := range bjtis {
						bjti.Expiry = time.Time{}
						CompareWithFixture(t, bjti, "hydra_oauth2_jti_blacklist", bjti.ID)
					}
				})

				t.Run("case=hydra_oauth2_access", func(t *testing.T) {
					as := []sql.OAuth2RequestSQL{}
					c.RawQuery("SELECT * FROM hydra_oauth2_access").All(&as)
					if db == "cockroach" {
						require.Equal(t, 5, len(as))
					} else {
						require.Equal(t, 13, len(as))
					}

					for _, a := range as {
						require.False(t, a.RequestedAt.IsZero())
						a.RequestedAt = time.Time{}
						require.NotZero(t, a.Client)
						a.Client = ""
						CompareWithFixture(t, a, "hydra_oauth2_access", a.ID)
					}
				})

				t.Run("case=hydra_oauth2_refresh", func(t *testing.T) {
					rs := []sql.OAuth2RequestSQL{}
					c.RawQuery("SELECT * FROM hydra_oauth2_refresh").All(&rs)
					if db == "cockroach" {
						require.Equal(t, 5, len(rs))
					} else {
						require.Equal(t, 13, len(rs))
					}

					for _, a := range rs {
						require.False(t, a.RequestedAt.IsZero())
						a.RequestedAt = time.Time{}
						require.NotZero(t, a.Client)
						a.Client = ""
						CompareWithFixture(t, a, "hydra_oauth2_refresh", a.ID)
					}
				})

				t.Run("case=hydra_oauth2_code", func(t *testing.T) {
					cs := []sql.OAuth2RequestSQL{}
					c.RawQuery("SELECT * FROM hydra_oauth2_code").All(&cs)
					if db == "cockroach" {
						require.Equal(t, 5, len(cs))
					} else {
						require.Equal(t, 13, len(cs))
					}

					for _, a := range cs {
						require.False(t, a.RequestedAt.IsZero())
						a.RequestedAt = time.Time{}
						require.NotZero(t, a.Client)
						a.Client = ""
						CompareWithFixture(t, a, "hydra_oauth2_code", a.ID)
					}
				})

				t.Run("case=hydra_oauth2_oidc", func(t *testing.T) {
					os := []sql.OAuth2RequestSQL{}
					c.RawQuery("SELECT * FROM hydra_oauth2_oidc").All(&os)
					if db == "cockroach" {
						require.Equal(t, 5, len(os))
					} else {
						require.Equal(t, 13, len(os))
					}

					for _, a := range os {
						require.False(t, a.RequestedAt.IsZero())
						a.RequestedAt = time.Time{}
						require.NotZero(t, a.Client)
						a.Client = ""
						CompareWithFixture(t, a, "hydra_oauth2_oidc", a.ID)
					}
				})

				t.Run("case=hydra_oauth2_pkce", func(t *testing.T) {
					ps := []sql.OAuth2RequestSQL{}
					c.RawQuery("SELECT * FROM hydra_oauth2_pkce").All(&ps)
					if db == "cockroach" {
						require.Equal(t, 5, len(ps))
					} else {
						require.Equal(t, 11, len(ps))
					}

					for _, a := range ps {
						require.False(t, a.RequestedAt.IsZero())
						a.RequestedAt = time.Time{}
						require.NotZero(t, a.Client)
						a.Client = ""
						CompareWithFixture(t, a, "hydra_oauth2_pkce", a.ID)
					}
				})

				t.Run("case=networks", func(t *testing.T) {
					ns := []networkx.Network{}
					c.RawQuery("SELECT * FROM networks").All(&ns)
					require.Equal(t, 1, len(ns))
					for _, n := range ns {
						require.NotZero(t, n.CreatedAt)
						require.NotZero(t, n.UpdatedAt)
						CompareWithFixture(t, n, "networks", n.ID.String())
					}
				})

			})
		}
	}

	for db, c := range connections {
		t.Run(fmt.Sprintf("database=%s", db), test(db, c))
		x.CleanSQLPop(t, c)
		require.NoError(t, c.Close())
	}
}
