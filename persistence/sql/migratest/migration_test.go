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
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/fatih/structs"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/flow"
	testhelpersuuid "github.com/ory/hydra/v2/internal/testhelpers/uuid"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/pop/v6"
	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/popx"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/ory/x/sqlxx"
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
	connections := make(map[string]*pop.Connection, 4)

	{
		c, err := pop.NewConnection(&pop.ConnectionDetails{URL: dbal.NewSQLiteTestDatabase(t)})
		require.NoError(t, err)
		require.NoError(t, c.Open())
		connections["sqlite"] = c
	}

	if !testing.Short() {
		wg := sync.WaitGroup{}
		for db, dsn := range map[string]string{
			"postgres":  dockertest.RunTestPostgreSQL(t),
			"mysql":     dockertest.RunTestMySQL(t),
			"cockroach": dockertest.RunTestCockroachDBWithVersion(t, "latest-v25.1"),
		} {
			wg.Add(1)
			go func() {
				defer wg.Done()

				dbName := "testdb" + strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
				t.Logf("using %s database %q", db, dbName)

				require.EventuallyWithT(t, func(t *assert.CollectT) {
					c, err := pop.NewConnection(&pop.ConnectionDetails{URL: dsn})
					require.NoError(t, err)
					require.NoError(t, c.Open())
					require.NoError(t, c.RawQuery("CREATE DATABASE "+dbName).Exec())
					dsn = regexp.MustCompile(`/[a-z0-9]+\?`).ReplaceAllString(dsn, "/"+dbName+"?")
					require.NoError(t, c.Close())

					c, err = pop.NewConnection(&pop.ConnectionDetails{URL: dsn})
					require.NoError(t, err)
					require.NoError(t, c.Open())
					connections[db] = c
				}, 20*time.Second, 100*time.Millisecond)
				t.Cleanup(func() {
					connections[db].Close()
				})
			}()
		}
		wg.Wait()
	}

	for db, c := range connections {
		t.Run("database="+db, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

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
					require.Len(t, cs, 20)
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
					require.Len(t, js, 7)
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
				require.Len(t, flows, 18)

				t.Run("case=hydra_oauth2_flow", func(t *testing.T) {
					for _, f := range flows {
						fixturizeFlow(t, &f)
						CompareWithFixture(t, f, "hydra_oauth2_flow", f.ID)
					}
				})

				t.Run("case=hydra_oauth2_authentication_session", func(t *testing.T) {
					ss := []flow.LoginSession{}
					require.NoError(t, c.All(&ss))
					require.Len(t, ss, 17)

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
					require.Len(t, ss, 13)

					for _, s := range ss {
						testhelpersuuid.AssertUUID(t, s.NID)
						s.NID = uuid.Nil
						CompareWithFixture(t, s, "hydra_oauth2_obfuscated_authentication_session", fmt.Sprintf("%s_%s", s.Subject, s.ClientID))
					}
				})

				t.Run("case=hydra_oauth2_logout_request", func(t *testing.T) {
					lrs := []flow.LogoutRequest{}
					require.NoError(t, c.All(&lrs))
					require.Len(t, lrs, 7)

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
					require.Len(t, bjtis, 1)
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
					require.Len(t, as, 13)

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
					rs := []sql.OAuth2RefreshTable{}
					require.NoError(t, c.All(&rs))
					require.Len(t, rs, 14)

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
					require.Len(t, cs, 13)

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
					require.Len(t, os, 13)

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
					require.Len(t, ps, 11)

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

				t.Run("case=hydra_oauth2_device_auth_codes", func(t *testing.T) {
					rs := []sql.DeviceRequestSQL{}
					require.NoError(t, c.All(&rs))
					require.Len(t, rs, 1)

					for _, r := range rs {
						testhelpersuuid.AssertUUID(t, r.NID)
						r.NID = uuid.Nil
						CompareWithFixture(t, r, "hydra_oauth2_device_auth_codes", r.ID)
					}
				})

				t.Run("case=networks", func(t *testing.T) {
					ns := []networkx.Network{}
					require.NoError(t, c.RawQuery("SELECT * FROM networks").All(&ns))
					require.Len(t, ns, 1)
					for _, n := range ns {
						testhelpersuuid.AssertUUID(t, n.ID)
						require.NotZero(t, n.CreatedAt)
						require.NotZero(t, n.UpdatedAt)
					}
				})
			})
		})
	}
}
