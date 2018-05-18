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
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/sqlcon/dockertest"
	"github.com/sirupsen/logrus"
)

var clientManagers = map[string]pkg.FositeStorer{}
var clientManager = &client.MemoryManager{
	Clients: []client.Client{{ID: "foobar"}},
	Hasher:  &fosite.BCrypt{},
}

func init() {
	clientManagers["memory"] = &FositeMemoryStore{
		AuthorizeCodes:      make(map[string]fosite.Requester),
		IDSessions:          make(map[string]fosite.Requester),
		AccessTokens:        make(map[string]fosite.Requester),
		RefreshTokens:       make(map[string]fosite.Requester),
		AccessTokenLifespan: time.Hour,
		PKCES:               make(map[string]fosite.Requester),
	}
}

func TestMain(m *testing.M) {
	runner := dockertest.Register()

	flag.Parse()
	if !testing.Short() {
		dockertest.Parallel([]func(){
			connectToPG,
			connectToMySQL,
		})
	}

	runner.Exit(m.Run())
}

func connectToPG() {
	db, err := dockertest.ConnectToTestPostgreSQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	s := &FositeSQLStore{DB: db, Manager: clientManager, L: logrus.New(), AccessTokenLifespan: time.Hour}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func connectToMySQL() {
	db, err := dockertest.ConnectToTestMySQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	s := &FositeSQLStore{DB: db, Manager: clientManager, L: logrus.New(), AccessTokenLifespan: time.Hour}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

func TestCreateGetDeleteAuthorizeCodes(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteAuthorizeCodes(m))
	}
}

func TestCreateGetDeleteAccessTokenSession(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteAccessTokenSession(m))
	}
}

func TestCreateGetDeleteOpenIDConnectSession(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteOpenIDConnectSession(m))
	}
}

func TestCreateGetDeleteRefreshTokenSession(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteRefreshTokenSession(m))
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperRevokeRefreshToken(m))
	}
}

func TestPKCEReuqest(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeletePKCERequestSession(m))
	}
}

func TestFlushAccessTokens(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperFlushTokens(m, time.Hour))
	}
}
