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

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
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

func TestMain(m *testing.M) {
	flag.Parse()
	runner := dockertest.Register()
	runner.Exit(m.Run())
}

func connectToPG(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestPostgreSQL()
	require.NoError(t, err)
	x.CleanSQL(t, db)
	return db
}

func connectToMySQL(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestMySQL()
	require.NoError(t, err)
	x.CleanSQL(t, db)
	return db
}

func connectToCRDB(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestCockroachDB()
	require.NoError(t, err)
	x.CleanSQL(t, db)
	return db
}

func connectSQL(t *testing.T, conf *configuration.ViperProvider, dbName string, db *sqlx.DB) driver.Registry {
	x.CleanSQL(t, db)
	reg := internal.NewRegistrySQL(conf, db)
	_, err := reg.CreateSchemas(dbName)
	require.NoError(t, err)
	return reg
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
				var p, m, c *sqlx.DB
				dockertest.Parallel([]func(){
					func() {
						p = connectToPG(t)
					},
					func() {
						m = connectToMySQL(t)
					},
					func() {
						c = connectToCRDB(t)
					},
				})
				registries["postgres"] = connectSQL(t, conf, "postgres", p)
				registries["mysql"] = connectSQL(t, conf, "mysql", m)
				registries["cockroach"] = connectSQL(t, conf, "cockroach", c)
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
