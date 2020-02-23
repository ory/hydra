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

package jwk_test

import (
	"fmt"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"
)

var managers = map[string]Manager{}

var m sync.Mutex
var testGenerator = &RS256Generator{}

func TestMain(m *testing.M) {
	runner := dockertest.Register()

	runner.Exit(m.Run())
}

func connectToPG(t *testing.T) *SQLManager {
	db, err := dockertest.ConnectToTestPostgreSQL()
	require.NoError(t, err)
	x.CleanSQL(t, db)
	return internal.NewRegistrySQL(internal.NewConfigurationWithDefaults(), db).KeyManager().(*SQLManager)
}

func connectToMySQL(t *testing.T) *SQLManager {
	db, err := dockertest.ConnectToTestMySQL()
	require.NoError(t, err)
	x.CleanSQL(t, db)
	return internal.NewRegistrySQL(internal.NewConfigurationWithDefaults(), db).KeyManager().(*SQLManager)
}

func connectToCRDB(t *testing.T) *SQLManager {
	db, err := dockertest.ConnectToTestCockroachDB()
	require.NoError(t, err)
	x.CleanSQL(t, db)
	return internal.NewRegistrySQL(internal.NewConfigurationWithDefaults(), db).KeyManager().(*SQLManager)
}

func TestManager(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)
	managers["memory"] = reg.KeyManager()

	if !testing.Short() {
		dockertest.Parallel([]func(){
			func() {
				m.Lock()
				managers["postgres"] = connectToPG(t)
				m.Unlock()
			},
			func() {
				m.Lock()
				managers["mysql"] = connectToMySQL(t)
				m.Unlock()
			},
			func() {
				m.Lock()
				managers["cockroach"] = connectToCRDB(t)
				m.Unlock()
			},
		})
	}

	t.Run("TestManagerKey", func(t *testing.T) {
		ks, err := testGenerator.Generate("TestManagerKey", "sig")
		require.NoError(t, err)

		for name, m := range managers {
			if m, ok := m.(*SQLManager); ok {
				n, err := m.CreateSchemas(name)
				require.NoError(t, err)
				t.Logf("Applied %d migrations to %s", n, name)
			}
			t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKey(m, ks, "TestManagerKey"))
		}
	})

	t.Run("TestManagerKeySet", func(t *testing.T) {
		ks, err := testGenerator.Generate("TestManagerKeySet", "sig")
		require.NoError(t, err)
		ks.Key("private")

		for name, m := range managers {
			if m, ok := m.(*SQLManager); ok {
				n, err := m.CreateSchemas(name)
				require.NoError(t, err)
				t.Logf("Applied %d migrations to %s", n, name)
			}
			t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKeySet(m, ks, "TestManagerKeySet"))
		}
	})
}
