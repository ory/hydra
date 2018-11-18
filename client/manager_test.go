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
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	. "github.com/ory/hydra/client"
	"github.com/ory/x/sqlcon/dockertest"
)

var clientManagers = map[string]Manager{}
var m sync.Mutex

func init() {
	clientManagers["memory"] = NewMemoryManager(&fosite.BCrypt{})
}

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

	s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}
	m.Lock()
	clientManagers["mysql"] = s
	m.Unlock()
}

func connectToPG() {
	db, err := dockertest.ConnectToTestPostgreSQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}
	m.Lock()
	clientManagers["postgres"] = s
	m.Unlock()
}

func TestManagers(t *testing.T) {
	if !testing.Short() {
		dockertest.Parallel([]func(){
			connectToPG,
			connectToMySQL,
		})
	}

	for k, m := range clientManagers {
		s, ok := m.(*SQLManager)
		if ok {
			CleanTestDB(t, s.DB)
			_, err := s.CreateSchemas()
			require.NoError(t, err)
		}

		t.Run("case=create-get-delete", func(t *testing.T) {
			t.Run(fmt.Sprintf("db=%s", k), TestHelperCreateGetDeleteClient(k, m))
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
