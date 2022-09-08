/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2_test

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/x/contextx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/sqlcon/dockertest"
)

func TestMain(m *testing.M) {
	flag.Parse()

	runner := dockertest.Register()
	runner.Exit(m.Run())
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
	ctx := context.TODO()
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
		t.Run(tc.name, func(t *testing.T) {
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
