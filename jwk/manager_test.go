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
	"context"
	"fmt"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/viper"
	"strings"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"
)

var managers = map[string]Manager{}

var m sync.Mutex
var testGenerator = &RS256Generator{}

func getManager(t *testing.T, url string) Manager {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyDSN, url)
	reg, err := driver.NewRegistry(conf)
	require.NoError(t, err)
	require.NoError(t, reg.Init())
	require.NoError(t, reg.Persister().MigrateUp(context.Background()))
	return reg.KeyManager()
}

func connectToMySQL(t *testing.T) Manager {
	c := dockertest.ConnectToTestMySQLPop(t)
	x.CleanSQLPop(t, c)
	url := c.URL()
	if !strings.HasPrefix(url, "mysql") {
		url = "mysql://" + url
	}
	return getManager(t, url)
}

func connectToPG(t *testing.T) Manager {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	x.CleanSQLPop(t, c)
	return getManager(t, c.URL())
}

func connectToCRDB(t *testing.T) Manager {
	c := dockertest.ConnectToTestCockroachDBPop(t)
	x.CleanSQLPop(t, c)
	url := c.URL()
	if !strings.HasPrefix(url, "cockroach") {
		url = "cockroach://" + strings.Split(url, "://")[1]
	}
	return getManager(t, url)
}

func TestManager(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)
	managers["memory"] = reg.KeyManager()

	if !testing.Short() {
		managers["postgres"] = connectToPG(t)
		managers["mysql"] = connectToMySQL(t)
		managers["cockroach"] = connectToCRDB(t)
	}

	t.Run("TestManagerKey", func(t *testing.T) {
		ks, err := testGenerator.Generate("TestManagerKey", "sig")
		require.NoError(t, err)

		for name, m := range managers {
			t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKey(m, ks, "TestManagerKey"))
		}
	})

	t.Run("TestManagerKeySet", func(t *testing.T) {
		ks, err := testGenerator.Generate("TestManagerKeySet", "sig")
		require.NoError(t, err)
		ks.Key("private")

		for name, m := range managers {
			t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKeySet(m, ks, "TestManagerKeySet"))
		}
	})
}
