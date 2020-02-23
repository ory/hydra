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

package client_test

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/client"
	"github.com/ory/hydra/internal"
	"github.com/ory/x/sqlcon/dockertest"
)

var clientManagers = map[string]Manager{}
var m sync.Mutex

func TestMain(m *testing.M) {
	flag.Parse()
	runner := dockertest.Register()
	runner.Exit(m.Run())
}

func connectToMySQL() {
	db, err := dockertest.ConnectToTestMySQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistrySQL(conf, db)

	m.Lock()
	clientManagers["mysql"] = reg.ClientManager()
	m.Unlock()
}

func connectToPG() {
	db, err := dockertest.ConnectToTestPostgreSQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistrySQL(conf, db)

	m.Lock()
	clientManagers["postgres"] = reg.ClientManager()
	m.Unlock()
}

func connectToCRDB() {
	db, err := dockertest.ConnectToTestCockroachDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistrySQL(conf, db)

	m.Lock()
	clientManagers["cockroach"] = reg.ClientManager()
	m.Unlock()
}

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	clientManagers["memory"] = reg.ClientManager()

	if !testing.Short() {
		dockertest.Parallel([]func(){
			connectToPG,
			connectToMySQL,
			connectToCRDB,
		})
	}

	t.Log("Creating schemas...")
	for k, m := range clientManagers {
		s, ok := m.(*SQLManager)
		if ok {
			CleanTestDB(t, s.DB)
			x, err := s.CreateSchemas(k)
			if err != nil {
				t.Fatal("Could not create schemas", err.Error())
			} else {
				t.Logf("Schemas created. Rows affected: %+v", x)
			}
			require.NoError(t, err)
		}

		t.Run("case=create-get-update-delete", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperCreateGetUpdateDeleteClient(k, m))
		})

		t.Run("case=autogenerate-key", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperClientAutoGenerateKey(k, m))
		})

		t.Run("case=auth-client", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperClientAuthenticate(k, m))
		})

		if ok {
			CleanTestDB(t, s.DB)
		}
	}
}
