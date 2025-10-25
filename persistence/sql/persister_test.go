// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent/test"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/pop/v6"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/networkx"
	"github.com/ory/x/servicelocatorx"
)

func init() {
	pop.SetNowFunc(func() time.Time {
		return time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func testRegistry(t *testing.T, db string, t1, t2 *driver.RegistrySQL) {
	// TODO enable parallel tests for cockroach once we configure the cockroach integration test server to support retry
	var parallel bool
	switch db {
	case "memory", "mysql", "cockroach":
		parallel = false
	default:
		parallel = true
	}

	t.Run("client", func(t *testing.T) {
		if parallel {
			// currently not possible as we have a lot of side-effects on listing of the clients between this and other tests
			// t.Parallel()
		}

		t.Run("case=create-get-update-delete", client.TestHelperCreateGetUpdateDeleteClient(t1.ClientManager(), t2.ClientManager()))

		t.Run("case=autogenerate-key", client.TestHelperClientAutoGenerateKey(t1.ClientManager()))

		t.Run("case=auth-client", client.TestHelperClientAuthenticate(t1.ClientManager()))

		t.Run("case=update-two-clients", client.TestHelperUpdateTwoClients(t1.ClientManager()))
	})

	for _, reg := range []*driver.RegistrySQL{t1, t2} {
		t.Run("consent", func(t *testing.T) {
			if parallel {
				t.Parallel()
			}
			test.ConsentManagerTests(t, reg, reg.ConsentManager(), reg.LoginManager(), reg.ClientManager(), reg.OAuth2Storage())
		})

		t.Run("login", func(t *testing.T) {
			if parallel {
				t.Parallel()
			}
			test.LoginManagerTest(t, reg, reg.LoginManager())
		})

		t.Run("obfuscated subject", func(t *testing.T) {
			if parallel {
				t.Parallel()
			}
			test.ObfuscatedSubjectManagerTest(t, reg, reg.ObfuscatedSubjectManager(), reg.ClientManager())
		})

		t.Run("logout", func(t *testing.T) {
			if parallel {
				t.Parallel()
			}
			test.LogoutManagerTest(t, reg.LogoutManager(), reg.ClientManager())
		})
	}

	t.Run("jwk", func(t *testing.T) {
		for _, tc := range []struct {
			alg  string
			skip bool
		}{
			{alg: "RS256", skip: false},
			{alg: "ES256", skip: false},
			{alg: "ES512", skip: false},
			{alg: "HS256", skip: true},
			{alg: "HS512", skip: true},
			{alg: "EdDSA", skip: t1.Config().HSMEnabled()},
		} {
			t.Run("key_generator="+tc.alg, func(t *testing.T) {
				if tc.skip {
					t.Skipf("Skipping test. Not applicable for alg: %s", tc.alg)
				}
				if t1.Config().HSMEnabled() {
					t.Run("TestManagerGenerateAndPersistKeySet", jwk.TestHelperManagerGenerateAndPersistKeySet(t1.KeyManager(), tc.alg, false))
					// We don't support NID isolation with HSM at the moment
					// t.Run("TestManagerGenerateAndPersistKeySet", jwk.TestHelperManagerNIDIsolationKeySet(t1.KeyManager(), t2.KeyManager(), tc.alg))
				} else {
					kid, err := uuid.NewV4()
					require.NoError(t, err)
					ks, err := jwk.GenerateJWK(jose.SignatureAlgorithm(tc.alg), kid.String(), "sig")
					require.NoError(t, err)
					t.Run("TestManagerKey", jwk.TestHelperManagerKey(t1.KeyManager(), tc.alg, ks, kid.String()))
					t.Run("Parallel", func(t *testing.T) {
						t.Run("TestManagerKeySet", jwk.TestHelperManagerKeySet(t1.KeyManager(), tc.alg, ks, kid.String(), parallel))
						t.Run("TestManagerKeySet", jwk.TestHelperManagerKeySet(t2.KeyManager(), tc.alg, ks, kid.String(), parallel))
					})
					t.Run("Parallel", func(t *testing.T) {
						t.Run("TestManagerGenerateAndPersistKeySet", jwk.TestHelperManagerGenerateAndPersistKeySet(t1.KeyManager(), tc.alg, parallel))
						t.Run("TestManagerGenerateAndPersistKeySet", jwk.TestHelperManagerGenerateAndPersistKeySet(t2.KeyManager(), tc.alg, parallel))
					})
				}
			})
		}

		t.Run("TestManagerGenerateAndPersistKeySetWithUnsupportedKeyGenerator", func(t *testing.T) {
			_, err := t1.KeyManager().GenerateAndPersistKeySet(context.TODO(), "foo", "bar", "UNKNOWN", "sig")
			require.Error(t, err)
			assert.IsType(t, errors.WithStack(jwk.ErrUnsupportedKeyAlgorithm), err)
		})
	})

	t.Run("trust", func(t *testing.T) {
		t.Run("parallel boundary", func(t *testing.T) {
			t.Run("case=create-get-delete/network=t1", trust.TestHelperGrantManagerCreateGetDeleteGrant(t1.GrantManager(), t1.KeyManager(), parallel))
			t.Run("case=create-get-delete/network=t2", trust.TestHelperGrantManagerCreateGetDeleteGrant(t2.GrantManager(), t2.KeyManager(), parallel))
		})
		t.Run("parallel boundary", func(t *testing.T) {
			t.Run("case=errors", trust.TestHelperGrantManagerErrors(t1.GrantManager(), t1.KeyManager()))
			t.Run("case=errors", trust.TestHelperGrantManagerErrors(t2.GrantManager(), t2.KeyManager()))
		})
	})
}

func TestManagersNextGen(t *testing.T) {
	t.Parallel()

	regs := testhelpers.ConnectDatabases(t, true, driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.TestContextualizer{})))

	ctx := t.Context()
	networks := make([]uuid.UUID, 5)
	for k := range networks {
		nid := uuid.Must(uuid.NewV4())
		for k := range regs {
			require.NoError(t, regs[k].Persister().Connection(ctx).Create(&networkx.Network{ID: nid}))
		}
		networks[k] = nid
	}

	for k := range regs {
		t.Run("database="+k, func(t *testing.T) {
			t.Parallel()
			client.TestHelperCreateGetUpdateDeleteClientNext(t, regs[k].Persister(), networks)
		})
	}
}

func TestManagers(t *testing.T) {
	t.Parallel()

	regs1, regs2, regsInvalid := make(map[string]*driver.RegistrySQL), make(map[string]*driver.RegistrySQL), make(map[string]*driver.RegistrySQL)
	network1NID, network2NID, invalidNID := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())

	sqlite := dbal.NewSQLiteTestDatabase(t)
	var err error
	regs1["memory"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), sqlite, true,
		driver.DisableValidation(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
	require.NoError(t, err)
	regs2["memory"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), sqlite, false,
		driver.DisableValidation(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
	require.NoError(t, err)
	regsInvalid["memory"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), sqlite, false,
		driver.DisableValidation(),
		driver.SkipNetworkInit(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
	require.NoError(t, err)

	if !testing.Short() {
		pg, mysql, crdb := testhelpers.ConnectDatabasesURLs(t)
		eg, ctx := errgroup.WithContext(t.Context())
		eg.Go(func() (err error) {
			regs1["postgres"], err = testhelpers.NewRegistrySQLFromURL(ctx, pg, true,
				driver.DisableValidation(),
				driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
			return
		})
		eg.Go(func() (err error) {
			regs1["mysql"], err = testhelpers.NewRegistrySQLFromURL(ctx, mysql, true,
				driver.DisableValidation(),
				driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
			return
		})
		eg.Go(func() (err error) {
			regs1["cockroach"], err = testhelpers.NewRegistrySQLFromURL(ctx, crdb, true,
				driver.DisableValidation(),
				driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
			return
		})
		require.NoError(t, eg.Wait())

		regs2["postgres"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), pg, false,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
		require.NoError(t, err)
		regsInvalid["postgres"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), pg, false,
			driver.DisableValidation(),
			driver.SkipNetworkInit(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
		require.NoError(t, err)
		regs2["mysql"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), mysql, false,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
		require.NoError(t, err)
		regsInvalid["mysql"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), mysql, false,
			driver.DisableValidation(),
			driver.SkipNetworkInit(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
		require.NoError(t, err)
		regs2["cockroach"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), crdb, false,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
		require.NoError(t, err)
		regsInvalid["cockroach"], err = testhelpers.NewRegistrySQLFromURL(t.Context(), crdb, false,
			driver.DisableValidation(),
			driver.SkipNetworkInit(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
		require.NoError(t, err)
	}

	for db, r1 := range regs1 {
		t.Run(db, func(t *testing.T) {
			t.Parallel()

			r2 := regs2[db]
			rInv := regsInvalid[db]

			r1.Persister()
			require.NoError(t, r1.Persister().Connection(t.Context()).Create(&networkx.Network{ID: network1NID}))
			require.NoError(t, r1.Persister().Connection(t.Context()).Create(&networkx.Network{ID: network2NID}))

			require.Equal(t, network1NID, r1.Persister().NetworkID(t.Context()))
			require.Equal(t, network2NID, r2.Persister().NetworkID(t.Context()))
			require.Equal(t, invalidNID, rInv.Persister().NetworkID(t.Context()))

			t.Run("parallel boundary", func(t *testing.T) { testRegistry(t, db, r1, r2) })

			if db == "memory" {
				// The following tests rely on foreign key constraints, which some of them are not correctly created in the SQLite schema.
				return
			}

			// if !r1.Config().HSMEnabled() {
			t.Run("jwk nid",
				jwk.TestHelperNID(r1.KeyManager(), rInv.KeyManager()),
			)
			// }
			t.Run("login nid",
				test.LoginNIDTest(r1.Persister(), rInv.Persister()),
			)
		})
	}
}
