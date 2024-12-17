// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/ory/hydra/v2/consent/test"

	"github.com/go-jose/go-jose/v3"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/networkx"

	"github.com/ory/hydra/v2/jwk"

	"github.com/ory/hydra/v2/driver"
)

func init() {
	pop.SetNowFunc(func() time.Time {
		return time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func testRegistry(t *testing.T, ctx context.Context, k string, t1 driver.Registry, t2 driver.Registry) {
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

	t.Run("parallel-boundary", func(t *testing.T) {
		t.Run("package=consent/janitor="+k, testhelpers.JanitorTests(t1, "t1", parallel))
		t.Run("package=consent/janitor="+k, testhelpers.JanitorTests(t2, "t2", parallel))
	})

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
			t.Run("case=errors", trust.TestHelperGrantManagerErrors(t1.GrantManager(), t1.KeyManager(), parallel))
			t.Run("case=errors", trust.TestHelperGrantManagerErrors(t2.GrantManager(), t2.KeyManager(), parallel))
		})
	})
}

func TestManagersNextGen(t *testing.T) {
	regs := map[string]driver.Registry{
		"memory": testhelpers.NewRegistrySQLFromURL(t, dbal.NewSQLiteTestDatabase(t), true, &contextx.Default{}),
	}

	if !testing.Short() {
		regs["postgres"], regs["mysql"], regs["cockroach"], _ = testhelpers.ConnectDatabases(t, true, &contextx.Default{})
	}

	ctx := context.Background()
	networks := make([]uuid.UUID, 5)
	for k := range networks {
		nid := uuid.Must(uuid.NewV4())
		for k := range regs {
			require.NoError(t, regs[k].Persister().Connection(ctx).Create(&networkx.Network{ID: nid}))
		}
		networks[k] = nid
	}

	for k := range regs {
		regs[k].WithContextualizer(new(contextx.TestContextualizer))
	}

	for k := range regs {
		k := k
		t.Run("database="+k, func(t *testing.T) {
			t.Parallel()
			client.TestHelperCreateGetUpdateDeleteClientNext(t, regs[k].Persister(), networks)
		})
	}
}

func TestManagers(t *testing.T) {
	ctx := context.TODO()
	t1registries := map[string]driver.Registry{
		"memory": testhelpers.NewRegistrySQLFromURL(t, dbal.NewSQLiteTestDatabase(t), true, &contextx.Default{}),
	}

	t2registries := map[string]driver.Registry{
		"memory": testhelpers.NewRegistrySQLFromURL(t, dbal.NewSQLiteTestDatabase(t), false, &contextx.Default{}),
	}

	if !testing.Short() {
		t2registries["postgres"], t2registries["mysql"], t2registries["cockroach"], _ = testhelpers.ConnectDatabases(t, false, &contextx.Default{})
		t1registries["postgres"], t1registries["mysql"], t1registries["cockroach"], _ = testhelpers.ConnectDatabases(t, true, &contextx.Default{})
	}

	network1NID, _ := uuid.NewV4()
	network2NID, _ := uuid.NewV4()

	for k, t1 := range t1registries {
		t2 := t2registries[k]
		require.NoError(t, t1.Persister().Connection(ctx).Create(&networkx.Network{ID: network1NID}))
		require.NoError(t, t2.Persister().Connection(ctx).Create(&networkx.Network{ID: network2NID}))
		t1.WithContextualizer(&contextx.Static{NID: network1NID, C: t1.Config().Source(context.Background())})
		t2.WithContextualizer(&contextx.Static{NID: network2NID, C: t2.Config().Source(context.Background())})
		t.Run("parallel-boundary", func(t *testing.T) { testRegistry(t, ctx, k, t1, t2) })
	}

	for k, t1 := range t1registries {
		t2 := t2registries[k]
		t2.WithContextualizer(&contextx.Static{NID: uuid.Nil, C: t2.Config().Source(context.Background())})

		if !t1.Config().HSMEnabled() { // We don't support NID isolation with HSM at the moment
			t.Run("package=jwk/manager="+k+"/case=nid",
				jwk.TestHelperNID(t1.KeyManager(), t2.KeyManager()),
			)
		}
		t.Run("package=consent/manager="+k+"/case=nid",
			test.TestHelperNID(t1, t1.ConsentManager(), t2.ConsentManager()),
		)
	}
}
