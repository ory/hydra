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

package integration

import (
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/ladon"
	lsql "github.com/ory/ladon/manager/sql"
	"github.com/ory/sqlcon/dockertest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSQLSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
		return
	}

	var testGenerator = &jwk.RS256Generator{}
	ks, _ := testGenerator.Generate("foo", "sig")
	p1 := ks.Key("private:foo")
	r := fosite.NewRequest()
	r.ID = "foo"

	db, err := dockertest.ConnectToTestPostgreSQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	cm := &client.SQLManager{DB: db, Hasher: &fosite.BCrypt{}}
	jm := jwk.SQLManager{DB: db, Cipher: &jwk.AEAD{Key: []byte("11111111111111111111111111111111")}}
	om := oauth2.FositeSQLStore{Manager: cm, DB: db, L: logrus.New()}
	crm := consent.NewSQLManager(db, nil)
	pm := lsql.NewSQLManager(db, nil)

	_, err = pm.CreateSchemas("", "hydra_policy_migration")
	require.NoError(t, err)
	_, err = cm.CreateSchemas()
	require.NoError(t, err)
	require.NoError(t, err)
	_, err = jm.CreateSchemas()
	require.NoError(t, err)
	_, err = om.CreateSchemas()
	require.NoError(t, err)
	_, err = crm.CreateSchemas()
	require.NoError(t, err)

	require.NoError(t, jm.AddKey("integration-test-foo", jwk.First(p1)))
	require.NoError(t, pm.Create(&ladon.DefaultPolicy{ID: "integration-test-foo", Resources: []string{"foo"}, Actions: []string{"bar"}, Subjects: []string{"baz"}, Effect: "allow"}))
	require.NoError(t, cm.CreateClient(&client.Client{ID: "integration-test-foo"}))
	require.NoError(t, crm.CreateAuthenticationSession(&consent.AuthenticationSession{
		ID:              "foo",
		AuthenticatedAt: time.Now(),
		Subject:         "bar",
	}))
	require.NoError(t, om.CreateAccessTokenSession(nil, "asdfasdf", r))
}
