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

func testRegistry(t *testing.T, k string, t1, t2 *driver.RegistrySQL) {
	t.Run("package=client/manager="+k, func(t *testing.T) {
		t.Run("case=create-get-update-delete", client.TestHelperCreateGetUpdateDeleteClient(k, t1.Persister().Connection(context.Background()), t1.ClientManager(), t2.ClientManager()))

		t.Run("case=autogenerate-key", client.TestHelperClientAutoGenerateKey(k, t1.ClientManager()))

		t.Run("case=auth-client", client.TestHelperClientAuthenticate(k, t1.ClientManager()))

		t.Run("case=update-two-clients", client.TestHelperUpdateTwoClients(k, t1.ClientManager()))
	})

	parallel := true
	if k == "memory" || k == "mysql" || k == "cockroach" { // TODO enable parallel tests for cockroach once we configure the cockroach integration test server to support retry
		parallel = false
	}

	t.Run("package=consent/manager="+k, test.ManagerTests(t1, t1.ConsentManager(), t1.ClientManager(), t1.OAuth2Storage(), "t1", parallel))
	t.Run("package=consent/manager="+k, test.ManagerTests(t2, t2.ConsentManager(), t2.ClientManager(), t2.OAuth2Storage(), "t2", parallel))

	t.Run("package=jwk/manager="+k, func(t *testing.T) {
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
					ks, err := jwk.GenerateJWK(context.Background(), jose.SignatureAlgorithm(tc.alg), kid.String(), "sig")
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

	t.Run("package=grant/trust/manager="+k, func(t *testing.T) {
		t.Run("parallel-boundary", func(t *testing.T) {
			t.Run("case=create-get-delete/network=t1", trust.TestHelperGrantManagerCreateGetDeleteGrant(t1.GrantManager(), t1.KeyManager(), parallel))
			t.Run("case=create-get-delete/network=t2", trust.TestHelperGrantManagerCreateGetDeleteGrant(t2.GrantManager(), t2.KeyManager(), parallel))
		})
		t.Run("parallel-boundary", func(t *testing.T) {
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
	regs1["memory"] = testhelpers.NewRegistrySQLFromURL(t, sqlite, true,
		driver.DisableValidation(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
	regs2["memory"] = testhelpers.NewRegistrySQLFromURL(t, sqlite, false,
		driver.DisableValidation(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
	regsInvalid["memory"] = testhelpers.NewRegistrySQLFromURL(t, sqlite, false,
		driver.DisableValidation(),
		driver.SkipNetworkInit(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))

	if !testing.Short() {
		pg, mysql, crdb := testhelpers.ConnectDatabasesURLs(t)
		regs1["postgres"] = testhelpers.NewRegistrySQLFromURL(t, pg, true,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
		regs2["postgres"] = testhelpers.NewRegistrySQLFromURL(t, pg, false,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
		regsInvalid["postgres"] = testhelpers.NewRegistrySQLFromURL(t, pg, false,
			driver.DisableValidation(),
			driver.SkipNetworkInit(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
		regs1["mysql"] = testhelpers.NewRegistrySQLFromURL(t, mysql, true,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
		regs2["mysql"] = testhelpers.NewRegistrySQLFromURL(t, mysql, false,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
		regsInvalid["mysql"] = testhelpers.NewRegistrySQLFromURL(t, mysql, false,
			driver.DisableValidation(),
			driver.SkipNetworkInit(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
		regs1["cockroach"] = testhelpers.NewRegistrySQLFromURL(t, crdb, true,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network1NID})))
		regs2["cockroach"] = testhelpers.NewRegistrySQLFromURL(t, crdb, false,
			driver.DisableValidation(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: network2NID})))
		regsInvalid["cockroach"] = testhelpers.NewRegistrySQLFromURL(t, crdb, false,
			driver.DisableValidation(),
			driver.SkipNetworkInit(),
			driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: invalidNID})))
	}

	for k, r1 := range regs1 {
		r2 := regs2[k]
		rInv := regsInvalid[k]

		require.NoError(t, r1.Persister().Connection(t.Context()).Create(&networkx.Network{ID: network1NID}))
		require.NoError(t, r1.Persister().Connection(t.Context()).Create(&networkx.Network{ID: network2NID}))

		require.Equal(t, network1NID, r1.Persister().NetworkID(t.Context()))
		require.Equal(t, network2NID, r2.Persister().NetworkID(t.Context()))
		require.Equal(t, invalidNID, rInv.Persister().NetworkID(t.Context()))

		t.Run("parallel-boundary", func(t *testing.T) { testRegistry(t, k, r1, r2) })

		if k == "memory" {
			// The following tests rely on foreign key constraints, which some of them are not correctly created in the SQLite schema.
			continue
		}

		if !r1.Config().HSMEnabled() && k != "memory" { // We don't support NID isolation with HSM at the moment
			t.Run("package=jwk/manager="+k+"/case=nid",
				jwk.TestHelperNID(r1.KeyManager(), rInv.KeyManager()),
			)
		}
		t.Run("package=consent/manager="+k+"/case=nid",
			test.TestHelperNID(r1.ConsentManager(), rInv.ConsentManager()),
		)
	}
}
