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
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"
)

var registries = make(map[string]driver.Registry)

func getManager(t *testing.T, url string) driver.Registry {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyDSN, url)
	reg, err := driver.NewRegistry(conf)
	require.NoError(t, err)
	require.NoError(t, reg.Init())
	require.NoError(t, reg.Persister().MigrateUp(context.Background()))
	return reg
}

func connectToMySQL(t *testing.T) driver.Registry {
	c := dockertest.ConnectToTestMySQLPop(t)
	x.CleanSQLPop(t, c)
	url := c.URL()
	if !strings.HasPrefix(url, "mysql") {
		url = "mysql://" + url
	}
	return getManager(t, url)
}

func connectToPG(t *testing.T) driver.Registry {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	x.CleanSQLPop(t, c)
	return getManager(t, c.URL())
}

func connectToCRDB(t *testing.T) driver.Registry {
	c := dockertest.ConnectToTestCockroachDBPop(t)
	x.CleanSQLPop(t, c)
	url := c.URL()
	if !strings.HasPrefix(url, "cockroach") {
		url = "cockroach://" + strings.Split(url, "://")[1]
	}
	return getManager(t, url)
}

func TestManagers(t *testing.T) {
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
			conf := internal.NewConfigurationWithDefaults()
			viper.Set(configuration.ViperKeyEncryptSessionData, tc.enableSessionEncrypted)
			reg := internal.NewRegistry(conf)

			require.NoError(t, reg.ClientManager().CreateClient(context.Background(), &client.Client{ClientID: "foobar"})) // this is a workaround because the client is not being created for memory store by test helpers.
			registries["memory"] = reg

			if !testing.Short() {
				registries["postgres"] = connectToPG(t)
				registries["mysql"] = connectToMySQL(t)
				registries["cockroach"] = connectToCRDB(t)
			}

			for k, store := range registries {
				TestHelperRunner(t, store, k)
			}

			for _, m := range registries {
				if mm, ok := m.(*driver.RegistrySQL); ok {
					x.CleanSQL(t, mm.DB())
				}
			}
		})

	}
}
