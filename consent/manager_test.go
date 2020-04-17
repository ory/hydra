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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent_test

import (
	"context"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"
)

var regs = make(map[string]driver.Registry)

func getRegistry(t *testing.T, url string) driver.Registry {
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
	return getRegistry(t, url)
}

func connectToPG(t *testing.T) driver.Registry {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	x.CleanSQLPop(t, c)
	return getRegistry(t, c.URL())
}

func connectToCRDB(t *testing.T) driver.Registry {
	c := dockertest.ConnectToTestCockroachDBPop(t)
	x.CleanSQLPop(t, c)
	url := c.URL()
	if !strings.HasPrefix(url, "cockroach") {
		url = "cockroach://" + strings.Split(url, "://")[1]
	}
	return getRegistry(t, url)
}

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Hour)
	regs["memory"] = internal.NewRegistry(conf)

	if !testing.Short() {
		regs["postgres"] = connectToPG(t)
		regs["mysql"] = connectToMySQL(t)
		regs["cockroach"] = connectToCRDB(t)
	}

	for k, m := range regs {
		t.Run("manager="+k, ManagerTests(m.ConsentManager(), m.ClientManager(), m.OAuth2Storage()))
	}

	for _, m := range regs {
		if mm, ok := m.ConsentManager().(*SQLManager); ok {
			x.CleanSQL(t, mm.DB)
		}
	}
}
