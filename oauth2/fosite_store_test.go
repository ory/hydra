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
	"flag"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/sqlcon/dockertest"
)

type managerTestSetup struct {
	f  pkg.FositeStorer
	cl client.Manager
	co consent.Manager
}

var fositeStores = map[string]managerTestSetup{}
var clientManager = &client.MemoryManager{
	Clients: []client.Client{{ClientID: "foobar"}},
	Hasher:  &fosite.BCrypt{},
}
var fm = NewFositeMemoryStore(clientManager, time.Hour)
var databases = make(map[string]*sqlx.DB)
var m sync.Mutex

func init() {
	fositeStores["memory"] = managerTestSetup{
		f:  fm,
		cl: clientManager,
		co: consent.NewMemoryManager(fm),
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	runner := dockertest.Register()
	runner.Exit(m.Run())
}

func connectToPG(t *testing.T) {
	db, err := dockertest.ConnectToTestPostgreSQL()
	require.NoError(t, err)
	t.Logf("Cleaning postgres db...")
	cleanDB(t, db)
	t.Logf("Cleaned postgres db")

	c := client.NewSQLManager(db, &fosite.BCrypt{WorkFactor: 8})
	_, err = c.CreateSchemas()
	require.NoError(t, err)

	cm := consent.NewSQLManager(db, c, nil)
	_, err = cm.CreateSchemas()
	require.NoError(t, err)

	s := NewFositeSQLStore(c, db, logrus.New(), time.Hour, false)
	_, err = s.CreateSchemas()
	require.NoError(t, err)

	m.Lock()
	databases["postgres"] = db
	fositeStores["postgres"] = managerTestSetup{
		f:  s,
		co: cm,
		cl: c,
	}
	m.Unlock()
}

func connectToMySQL(t *testing.T) {
	db, err := dockertest.ConnectToTestMySQL()
	require.NoError(t, err)
	t.Logf("Cleaning mysql db...")
	cleanDB(t, db)
	t.Logf("Cleaned mysql db")

	c := client.NewSQLManager(db, &fosite.BCrypt{WorkFactor: 8})
	_, err = c.CreateSchemas()
	require.NoError(t, err)

	cm := consent.NewSQLManager(db, c, nil)
	_, err = cm.CreateSchemas()
	require.NoError(t, err)

	s := NewFositeSQLStore(c, db, logrus.New(), time.Hour, false)
	_, err = s.CreateSchemas()
	require.NoError(t, err)

	m.Lock()
	databases["mysql"] = db
	fositeStores["mysql"] = managerTestSetup{
		f:  s,
		co: cm,
		cl: c,
	}
	m.Unlock()
}

func TestManagers(t *testing.T) {
	if !testing.Short() {
		dockertest.Parallel([]func(){
			func() {
				connectToPG(t)
			},
			func() {
				connectToMySQL(t)
			},
		})
	}

	for k, store := range fositeStores {
		if k != "memory" {
			t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteAuthorizeCodes/db=%s", k), testHelperUniqueConstraints(store, k))
		}
		t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteAuthorizeCodes/db=%s", k), testHelperCreateGetDeleteAuthorizeCodes(store))
		t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteAccessTokenSession/db=%s", k), testHelperCreateGetDeleteAccessTokenSession(store))
		t.Run(fmt.Sprintf("case=testHelperNilAccessToken/db=%s", k), testHelperNilAccessToken(store))
		t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteOpenIDConnectSession/db=%s", k), testHelperCreateGetDeleteOpenIDConnectSession(store))
		t.Run(fmt.Sprintf("case=testHelperCreateGetDeleteRefreshTokenSession/db=%s", k), testHelperCreateGetDeleteRefreshTokenSession(store))
		t.Run(fmt.Sprintf("case=testHelperRevokeRefreshToken/db=%s", k), testHelperRevokeRefreshToken(store))
		t.Run(fmt.Sprintf("case=testHelperCreateGetDeletePKCERequestSession/db=%s", k), testHelperCreateGetDeletePKCERequestSession(store))
		t.Run(fmt.Sprintf("case=testHelperFlushTokens/db=%s", k), testHelperFlushTokens(store, time.Hour))
	}

	for _, m := range databases {
		cleanDB(t, m)
	}
}
