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
	"flag"
	"sync"
	"testing"
	"time"

	"github.com/ory/hydra/x"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/consent"
	"github.com/ory/x/sqlcon/dockertest"
)

var m sync.Mutex
var regs = make(map[string]driver.Registry)

func connectToPostgres(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestPostgreSQL()
	require.NoError(t, err)
	t.Logf("Cleaning postgres db...")
	x.CleanSQL(t, db)
	t.Logf("Cleaned postgres db")
	return db
}

func connectToMySQL(t *testing.T) *sqlx.DB {
	db, err := dockertest.ConnectToTestMySQL()
	require.NoError(t, err)
	t.Logf("Cleaning mysql db...")
	x.CleanSQL(t, db)
	t.Logf("Cleaned mysql db")
	return db
}

func TestMain(m *testing.M) {
	flag.Parse()
	runner := dockertest.Register()
	runner.Exit(m.Run())
}

func createSQL(db *sqlx.DB) driver.Registry {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistrySQL(conf, db)
	if _, err := reg.CreateSchemas(); err != nil {
		panic(err)
	}

	return reg
}

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Hour)
	regs["memory"] = internal.NewRegistry(conf)

	if !testing.Short() {
		var p, m *sqlx.DB
		dockertest.Parallel([]func(){
			func() {
				p = connectToPostgres(t)
			}, func() {
				m = connectToMySQL(t)
			},
		})
		regs["postgres"] = createSQL(p)
		regs["mysql"] = createSQL(m)
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
