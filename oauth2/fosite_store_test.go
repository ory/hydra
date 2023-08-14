// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal"
	. "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/contextx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/sqlcon/dockertest"
)

func TestMain(m *testing.M) {
	flag.Parse()

	defer dockertest.KillAllTestDatabases()
	m.Run()
}

var registries = make(map[string]driver.Registry)
var cleanRegistries = func(t *testing.T) {
	registries["memory"] = internal.NewRegistryMemory(t, internal.NewConfigurationWithDefaults(), &contextx.Default{})
}

// returns clean registries that can safely be used for one test
// to reuse call cleanRegistries
func setupRegistries(t *testing.T) {
	if len(registries) == 0 && !testing.Short() {
		// first time called and sql tests
		var cleanSQL func(*testing.T)
		registries["postgres"], registries["mysql"], registries["cockroach"], cleanSQL = internal.ConnectDatabases(t, true, &contextx.Default{})
		cleanMem := cleanRegistries
		cleanMem(t)
		cleanRegistries = func(t *testing.T) {
			cleanMem(t)
			cleanSQL(t)
		}
	} else {
		// reset all/init mem
		cleanRegistries(t)
	}
}

func TestManagers(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name                   string
		enableSessionEncrypted bool
	}{
		{
			name:                   "DisableSessionEncrypted",
			enableSessionEncrypted: false,
		},
		{
			name:                   "EnableSessionEncrypted",
			enableSessionEncrypted: true,
		},
	}
	for _, tc := range tests {
		t.Run("suite="+tc.name, func(t *testing.T) {
			setupRegistries(t)

			require.NoError(t, registries["memory"].ClientManager().CreateClient(context.Background(), &client.Client{LegacyClientID: "foobar"})) // this is a workaround because the client is not being created for memory store by test helpers.

			for k, store := range registries {
				net := &networkx.Network{}
				require.NoError(t, store.Persister().Connection(context.Background()).First(net))
				store.Config().MustSet(ctx, config.KeyEncryptSessionData, tc.enableSessionEncrypted)
				store.WithContextualizer(&contextx.Static{NID: net.ID, C: store.Config().Source(ctx)})
				TestHelperRunner(t, store, k)
			}
		})

	}
}
