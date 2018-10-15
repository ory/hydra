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
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/x/sqlcon/dockertest"

	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/consent"
	"github.com/ory/hydra/oauth2"
)

var m sync.Mutex
var clientManager = client.NewMemoryManager(&fosite.BCrypt{WorkFactor: 8})
var fositeManager = oauth2.NewFositeMemoryStore(clientManager, time.Hour)
var managers = map[string]Manager{
	"memory": NewMemoryManager(fositeManager),
}

func connectToPostgres(managers map[string]Manager, c client.Manager) {
	db, err := dockertest.ConnectToTestPostgreSQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return
	}

	s := NewSQLManager(db, c, fositeManager)
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return
	}

	m.Lock()
	managers["postgres"] = s
	m.Unlock()
}

func connectToMySQL(managers map[string]Manager, c client.Manager) {
	db, err := dockertest.ConnectToTestMySQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return
	}

	s := NewSQLManager(db, c, fositeManager)
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create mysql schema: %v", err)
		return
	}

	m.Lock()
	managers["mysql"] = s
	m.Unlock()
}

func TestMain(m *testing.M) {
	runner := dockertest.Register()

	flag.Parse()
	if !testing.Short() {
		dockertest.Parallel([]func(){
			func() {
				connectToPostgres(managers, clientManager)
			}, func() {
				connectToMySQL(managers, clientManager)
			},
		})
	}

	runner.Exit(m.Run())
}

func TestManagers(t *testing.T) {
	for k, m := range managers {
		t.Run("manager="+k, ManagerTests(m, clientManager, fositeManager))
	}
}
