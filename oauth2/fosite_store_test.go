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
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/configuration"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/internal"

	. "github.com/ory/hydra/oauth2"
	"github.com/ory/x/sqlcon/dockertest"
)

var registries = make(map[string]driver.Registry)

var m sync.Mutex

func TestMain(m *testing.M) {
	flag.Parse()
	runner := dockertest.Register()
	runner.Exit(m.Run())
}

func connectToPG(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestPostgreSQL()
	require.NoError(t, err)
	t.Logf("Cleaning postgres db...")
	cleanDB(t, db)
	t.Logf("Cleaned postgres db")

	return db
}

func connectToMySQL(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestMySQL()
	require.NoError(t, err)
	t.Logf("Cleaning mysql db...")
	cleanDB(t, db)
	t.Logf("Cleaned mysql db")

	return db
}

func connectSQL(t *testing.T, conf *configuration.ViperProvider, db *sqlx.DB) driver.Registry {
	reg := internal.NewRegistrySQL(conf, db)
	_, err := reg.ClientManager().(*client.SQLManager).CreateSchemas()
	require.NoError(t, err)
	_, err = reg.ConsentManager().(*consent.SQLManager).CreateSchemas()
	require.NoError(t, err)
	_, err = reg.OAuth2Storage().(*FositeSQLStore).CreateSchemas()
	require.NoError(t, err)
	return reg
}

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	require.NoError(t, reg.ClientManager().CreateClient(context.Background(), &client.Client{ClientID: "foobar"})) // this is a workaround because the client is not being created for memory store by test helpers.
	registries["memory"] = reg

	if !testing.Short() {
		dockertest.Parallel([]func(){
			func() {
				m.Lock()

				registries["postgres"] = connectSQL(t, conf, connectToPG(t))
				m.Unlock()
			},
			func() {
				m.Lock()
				registries["mysql"] = connectSQL(t, conf, connectToMySQL(t))
				m.Unlock()
			},
		})
	}

	for k, store := range registries {
		TestHelperRunner(t, store, k)
	}

	for _, m := range registries {
		if mm, ok := m.(*driver.RegistrySQL); ok {
			cleanDB(t, mm.DB())
		}
	}
}
