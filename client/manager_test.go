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
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/x"
	"github.com/ory/viper"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"

	. "github.com/ory/hydra/client"
	"github.com/ory/hydra/internal"
	"github.com/ory/x/sqlcon/dockertest"
)

func getManager(t *testing.T, url string) Manager {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyDSN, url)
	reg, err := driver.NewRegistry(conf)
	require.NoError(t, err)
	require.NoError(t, reg.Init())
	require.NoError(t, reg.Persister().MigrateUp(context.Background()))
	return reg.ClientManager()
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

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	clientManagers := map[string]Manager{
		"memory": reg.ClientManager(),
	}

	if !testing.Short() {
		clientManagers["postgres"] = connectToPG(t)
		clientManagers["mysql"] = connectToMySQL(t)
		clientManagers["cockroach"] = connectToCRDB(t)
	}

	t.Log("Creating schemas...")
	for k, m := range clientManagers {
		t.Run("case=create-get-update-delete", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperCreateGetUpdateDeleteClient(k, m))
		})

		t.Run("case=autogenerate-key", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperClientAutoGenerateKey(k, m))
		})

		t.Run("case=auth-client", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperClientAuthenticate(k, m))
		})
	}
}
