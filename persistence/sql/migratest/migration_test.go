// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package migratest

import (
	"context"
	stdsql "database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/ory/x/contextx"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/fatih/structs"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/v2/flow"
	testhelpersuuid "github.com/ory/hydra/v2/internal/testhelpers/uuid"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/x/popx"

	"github.com/ory/x/sqlcon/dockertest"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
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

func TestMigrations(t *testing.T) {
	connections := make(map[string]*pop.Connection, 1)

	if testing.Short() {
		reg := testhelpers.NewMockedRegistry(t, &contextx.Default{})
		require.NoError(t, reg.Persister().MigrateUp(context.Background()))
		c := reg.Persister().Connection(context.Background())
		connections["sqlite"] = c
	}

	if !testing.Short() {
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
	}

	var test = func(db string, c *pop.Connection) func(t *testing.T) {
		return func(t *testing.T) {
			ctx := context.Background()
			x.CleanSQLPop(t, c)

			l := logrusx.New("", "", logrusx.ForceLevel(logrus.DebugLevel))

			tm, err := popx.NewMigrationBox(
				sql.Migrations,
				popx.NewMigrator(c, l, nil, 1*time.Minute),
				popx.WithTestdata(t, os.DirFS("./testdata")))
			require.NoError(t, err)
			require.NoError(t, tm.Up(ctx))

			t.Run("suite=fixtures", func(t *testing.T) {
				t.Run("case=hydra_client", func(t *testing.T) {
					cs := []client.Client{}
					require.NoError(t, c.All(&cs))
					require.Equal(t, 19, len(cs))
					for _, c := range cs {
						require.False(t, c.CreatedAt.IsZero())
						require.False(t, c.UpdatedAt.IsZero())
						c.CreatedAt = time.Time{} // Some CreatedAt and UpdatedAt values are generated during migrations so we zero them in the fixtures
						c.UpdatedAt = time.Time{}
						testhelpersuuid.AssertUUID(t, c.NID)
						testhelpersuuid.AssertUUID(t, c.PK.String)
						c.NID = uuid.Nil
						c.PK = stdsql.NullString{}
						CompareWithFixture(t, structs.Map(c), "hydra_client", c.ID)
					}
				})

				t.Run("case=hydra_jwk", func(t *testing.T) {
					js := []jwk.SQLData{}
					require.NoError(t, c.All(&js))
					require.Equal(t, 7, len(js))
					for _, j := range js {
						testhelpersuuid.AssertUUID(t, j.ID)
						testhelpersuuid.AssertUUID(t, j.NID)
						j.ID = uuid.Nil // Some IDs are generated at migration time so we zero them in the fixtures
						j.NID = uuid.Nil
						require.False(t, j.CreatedAt.IsZero())
						j.CreatedAt = time.Time{}
						CompareWithFixture(t, j, "hydra_jwk", j.KID)
					}
				})

				flows := []flow.Flow{}
				require.NoError(t, c.All(&flows))
				require.Equal(t, 17, len(flows))

				t.Run("case=hydra_oauth2_flow", func(t *testing.T) {
					for _, f := range flows {
						fixturizeFlow(t, &f)
						CompareWithFixture(t, f, "hydra_oauth2_flow", f.ID)
					}
				})

				t.Run("case=hydra_oauth2_authentication_session", func(t *testing.T) {
					ss := []flow.LoginSession{}
					require.NoError(t, c.All(&ss))
					require.Equal(t, 17, len(ss))

					for _, s := range ss {
						testhelpersuuid.AssertUUID(t, s.NID)
						s.NID = uuid.Nil
						s.AuthenticatedAt = sqlxx.NullTime(time.Time{})
						CompareWithFixture(t, s, "hydra_oauth2_authentication_session", s.ID)
					}
				})

				t.Run("case=hydra_oauth2_obfuscated_authentication_session", func(t *testing.T) {
					ss := []consent.ForcedObfuscatedLoginSession{}
					require.NoError(t, c.All(&ss))
					require.Equal(t, 13, len(ss))

					for _, s := range ss {
						testhelpersuuid.AssertUUID(t, s.NID)
						s.NID = uuid.Nil
						CompareWithFixture(t, s, "hydra_oauth2_obfuscated_authentication_session", fmt.Sprintf("%s_%s", s.Subject, s.ClientID))
					}
				})

				t.Run("case=hydra_oauth2_logout_request", func(t *testing.T) {
					lrs := []flow.LogoutRequest{}
					require.NoError(t, c.All(&lrs))
					require.Equal(t, 7, len(lrs))

					for _, s := range lrs {
						testhelpersuuid.AssertUUID(t, s.NID)
						s.NID = uuid.Nil
						s.Client = nil
						CompareWithFixture(t, s, "hydra_oauth2_logout_request", s.ID)
					}
				})

				t.Run("case=hydra_oauth2_jti_blacklist", func(t *testing.T) {
					bjtis := []oauth2.BlacklistedJTI{}
					require.NoError(t, c.All(&bjtis))
					require.Equal(t, 1, len(bjtis))
					for _, bjti := range bjtis {
						testhelpersuuid.AssertUUID(t, bjti.NID)
						bjti.NID = uuid.Nil
						bjti.Expiry = time.Time{}
						CompareWithFixture(t, bjti, "hydra_oauth2_jti_blacklist", bjti.ID)
					}
				})

				t.Run("case=hydra_oauth2_access", func(t *testing.T) {
					as := []sql.OAuth2RequestSQL{}
					require.NoError(t, c.RawQuery("SELECT * FROM hydra_oauth2_access").All(&as))
					require.Equal(t, 13, len(as))

					for _, a := range as {
						testhelpersuuid.AssertUUID(t, a.NID)
						a.NID = uuid.Nil
						require.False(t, a.RequestedAt.IsZero())
						a.RequestedAt = time.Time{}
						require.NotZero(t, a.Client)
						a.Client = ""
						CompareWithFixture(t, a, "hydra_oauth2_access", a.ID)
					}
				})

				t.Run("case=hydra_oauth2_refresh", func(t *testing.T) {
					rs := []sql.OAuth2RequestSQL{}
					require.NoError(t, c.RawQuery(`SELECT signature, nid, request_id, challenge_id, requested_at, client_id, scope, granted_scope, requested_audience, granted_audience, form_data, subject, active, session_data, expires_at	FROM hydra_oauth2_refresh`).All(&rs))
					require.Equal(t, 13, len(rs))

					for _, r := range rs {
						testhelpersuuid.AssertUUID(t, r.NID)
						r.NID = uuid.Nil
						require.False(t, r.RequestedAt.IsZero())
						r.RequestedAt = time.Time{}
						require.NotZero(t, r.Client)
						r.Client = ""
						CompareWithFixture(t, r, "hydra_oauth2_refresh", r.ID)
					}
				})

				t.Run("case=hydra_oauth2_code", func(t *testing.T) {
					cs := []sql.OAuth2RequestSQL{}
					require.NoError(t, c.RawQuery("SELECT * FROM hydra_oauth2_code").All(&cs))
					require.Equal(t, 13, len(cs))

					for _, c := range cs {
						testhelpersuuid.AssertUUID(t, c.NID)
						c.NID = uuid.Nil
						require.False(t, c.RequestedAt.IsZero())
						c.RequestedAt = time.Time{}
						require.NotZero(t, c.Client)
						c.Client = ""
						CompareWithFixture(t, c, "hydra_oauth2_code", c.ID)
					}
				})

				t.Run("case=hydra_oauth2_oidc", func(t *testing.T) {
					os := []sql.OAuth2RequestSQL{}
					require.NoError(t, c.RawQuery("SELECT * FROM hydra_oauth2_oidc").All(&os))
					require.Equal(t, 13, len(os))

					for _, o := range os {
						testhelpersuuid.AssertUUID(t, o.NID)
						o.NID = uuid.Nil
						require.False(t, o.RequestedAt.IsZero())
						o.RequestedAt = time.Time{}
						require.NotZero(t, o.Client)
						o.Client = ""
						CompareWithFixture(t, o, "hydra_oauth2_oidc", o.ID)
					}
				})

				t.Run("case=hydra_oauth2_pkce", func(t *testing.T) {
					ps := []sql.OAuth2RequestSQL{}
					require.NoError(t, c.RawQuery("SELECT * FROM hydra_oauth2_pkce").All(&ps))
					require.Equal(t, 11, len(ps))

					for _, p := range ps {
						testhelpersuuid.AssertUUID(t, p.NID)
						p.NID = uuid.Nil
						require.False(t, p.RequestedAt.IsZero())
						p.RequestedAt = time.Time{}
						require.NotZero(t, p.Client)
						p.Client = ""
						CompareWithFixture(t, p, "hydra_oauth2_pkce", p.ID)
					}
				})

				t.Run("case=networks", func(t *testing.T) {
					ns := []networkx.Network{}
					require.NoError(t, c.RawQuery("SELECT * FROM networks").All(&ns))
					require.Equal(t, 1, len(ns))
					for _, n := range ns {
						testhelpersuuid.AssertUUID(t, n.ID)
						require.NotZero(t, n.CreatedAt)
						require.NotZero(t, n.UpdatedAt)
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
