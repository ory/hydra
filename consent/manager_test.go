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

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/sqlcon/dockertest"
)

type managerRegistry struct {
	consent Manager
	client  client.Manager
	fosite  pkg.FositeStorer
}

var m sync.Mutex
var clientManager = client.NewMemoryManager(&fosite.BCrypt{WorkFactor: 8})
var fositeManager = oauth2.NewFositeMemoryStore(clientManager, time.Hour)
var managers = map[string]managerRegistry{
	"memory": {
		consent: NewMemoryManager(fositeManager),
		client:  clientManager,
		fosite:  fositeManager,
	},
}

func connectToPostgres(t *testing.T, managers map[string]managerRegistry) {
	db, err := dockertest.ConnectToTestPostgreSQL()
	require.NoError(t, err)
	t.Logf("Cleaning postgres db...")
	cleanDB(t, db)
	t.Logf("Cleaned postgres db")

	c := client.NewSQLManager(db, &fosite.BCrypt{WorkFactor: 8})
	d, err := c.CreateSchemas()
	require.NoError(t, err)
	t.Logf("Migrated %d postgres schemas", d)

	fositeManager := oauth2.NewFositeMemoryStore(c, time.Hour)

	s := NewSQLManager(db, c, fositeManager)
	d, err = s.CreateSchemas()
	t.Logf("Migrated %d postgres schemas", d)
	require.NoError(t, err)

	m.Lock()
	managers["postgres"] = managerRegistry{
		consent: s,
		client:  c,
		fosite:  fositeManager,
	}
	m.Unlock()
}

func connectToMySQL(t *testing.T, managers map[string]managerRegistry) {
	db, err := dockertest.ConnectToTestMySQL()
	require.NoError(t, err)
	t.Logf("Cleaning mysql db...")
	cleanDB(t, db)
	t.Logf("Cleaned mysql db")

	c := client.NewSQLManager(db, &fosite.BCrypt{WorkFactor: 8})
	d, err := c.CreateSchemas()
	require.NoError(t, err)
	t.Logf("Migrated %d mysql schemas", d)

	fositeManager := oauth2.NewFositeMemoryStore(c, time.Hour)

	s := NewSQLManager(db, c, fositeManager)
	d, err = s.CreateSchemas()
	t.Logf("Migrated %d mysql schemas", d)
	require.NoError(t, err)

	m.Lock()
	managers["mysql"] = managerRegistry{
		consent: s,
		client:  c,
		fosite:  fositeManager,
	}
	m.Unlock()
}

func TestMain(m *testing.M) {
	flag.Parse()
	runner := dockertest.Register()
	runner.Exit(m.Run())
}

func TestManagers(t *testing.T) {
	if !testing.Short() {
		dockertest.Parallel([]func(){
			func() {
				connectToPostgres(t, managers)
			}, func() {
				connectToMySQL(t, managers)
			},
		})
	}

	for k, m := range managers {
		t.Run("manager="+k, ManagerTests(m.consent, m.client, m.fosite))
	}

	for _, m := range managers {
		if mm, ok := m.consent.(*SQLManager); ok {
			cleanDB(t, mm.DB)
		}
	}
}
