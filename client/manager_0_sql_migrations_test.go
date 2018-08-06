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

package client_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/sqlcon/dockertest"
	"github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var createClientMigrations = []*migrate.Migration{
	{
		Id: "1-data",
		Up: []string{
			`INSERT INTO hydra_client (id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, public) VALUES ('1-data', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', true)`,
		},
		Down: []string{
			`DELETE FROM hydra_client WHERE id='1-data'`,
		},
	},
	{
		Id: "2-data",
		Up: []string{
			`INSERT INTO hydra_client (id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, public, client_secret_expires_at) VALUES ('2-data', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', true, 0)`,
		},
		Down: []string{
			`DELETE FROM hydra_client WHERE id='2-data'`,
		},
	},
	{
		Id: "3-data",
		Up: []string{
			`INSERT INTO hydra_client (id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, public, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg) VALUES ('3-data', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', true, 0, 'http://sector', '{"keys": []}', 'http://jwks', 'client_secret', 'http://uri1|http://uri2', 'rs256', 'rs526')`,
		},
		Down: []string{
			`DELETE FROM hydra_client WHERE id='3-data'`,
		},
	},
	{
		Id: "4-data",
		Up: []string{
			`INSERT INTO hydra_client (id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, public, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg) VALUES ('4-data', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', true, 0, 'http://sector', '{"keys": []}', 'http://jwks', 'client_secret', 'http://uri1|http://uri2', 'rs256', 'rs526')`,
		},
		Down: []string{
			`DELETE FROM hydra_client WHERE id='4-data'`,
		},
	},
	{
		Id: "5-data",
		Up: []string{
			`INSERT INTO hydra_client (id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg) VALUES ('5-data', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', 0, 'http://sector', '{"keys": []}', 'http://jwks', 'none', 'http://uri1|http://uri2', 'rs256', 'rs526')`,
		},
		Down: []string{
			`DELETE FROM hydra_client WHERE id='5-data'`,
		},
	},
	{
		Id: "6-data",
		Up: []string{
			`INSERT INTO hydra_client (id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg, subject_type) VALUES ('6-data', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', 0, 'http://sector', '{"keys": []}', 'http://jwks', 'none', 'http://uri1|http://uri2', 'rs256', 'rs526', 'public')`,
		},
		Down: []string{
			`DELETE FROM hydra_client WHERE id='6-data'`,
		},
	},
}

var migrations = map[string]*migrate.MemoryMigrationSource{
	"mysql": {
		Migrations: []*migrate.Migration{
			{Id: "0-data", Up: []string{"DROP TABLE IF EXISTS hydra_client"}},
			client.Migrations["mysql"].Migrations[0],
			createClientMigrations[0],
			client.Migrations["mysql"].Migrations[1],
			createClientMigrations[1],
			client.Migrations["mysql"].Migrations[2],
			createClientMigrations[2],
			client.Migrations["mysql"].Migrations[3],
			createClientMigrations[3],
			client.Migrations["mysql"].Migrations[4],
			createClientMigrations[4],
			client.Migrations["mysql"].Migrations[5],
			createClientMigrations[5],
		},
	},
	"postgres": {
		Migrations: []*migrate.Migration{
			{Id: "0-data", Up: []string{"DROP TABLE IF EXISTS hydra_client"}},
			client.Migrations["postgres"].Migrations[0],
			createClientMigrations[0],
			client.Migrations["postgres"].Migrations[1],
			createClientMigrations[1],
			client.Migrations["postgres"].Migrations[2],
			createClientMigrations[2],
			client.Migrations["postgres"].Migrations[3],
			createClientMigrations[3],
			client.Migrations["postgres"].Migrations[4],
			createClientMigrations[4],
			client.Migrations["postgres"].Migrations[5],
			createClientMigrations[5],
		},
	},
}

func TestMigrations(t *testing.T) {
	var m sync.Mutex
	var dbs = map[string]*sqlx.DB{}
	if testing.Short() {
		return
	}

	dockertest.Parallel([]func(){
		func() {
			db, err := dockertest.ConnectToTestPostgreSQL()
			if err != nil {
				log.Fatalf("Could not connect to database: %v", err)
			}
			m.Lock()
			dbs["postgres"] = db
			m.Unlock()
		},
		func() {
			db, err := dockertest.ConnectToTestMySQL()
			if err != nil {
				log.Fatalf("Could not connect to database: %v", err)
			}
			m.Lock()
			dbs["mysql"] = db
			m.Unlock()
		},
	})

	for k, db := range dbs {
		t.Run(fmt.Sprintf("database=%s", k), func(t *testing.T) {
			migrate.SetTable("hydra_client_migration_integration")
			for step := range migrations[k].Migrations {
				t.Run(fmt.Sprintf("step=%d", step), func(t *testing.T) {
					n, err := migrate.ExecMax(db.DB, db.DriverName(), migrations[k], migrate.Up, 1)
					require.NoError(t, err)
					require.Equal(t, n, 1)
				})
			}

			for _, key := range []string{"1-data", "2-data", "3-data", "4-data", "5-data"} {
				t.Run("client="+key, func(t *testing.T) {
					s := &client.SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}
					c, err := s.GetConcreteClient(key)
					require.NoError(t, err)
					assert.EqualValues(t, c.GetID(), key)
				})
			}

			for step := range migrations[k].Migrations {
				t.Run(fmt.Sprintf("step=%d", step), func(t *testing.T) {
					n, err := migrate.ExecMax(db.DB, db.DriverName(), migrations[k], migrate.Down, 1)
					require.NoError(t, err)
					require.Equal(t, n, 1)
				})
			}
		})
	}
}
