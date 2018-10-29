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

package jwk_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/jwk"
	"github.com/ory/x/sqlcon/dockertest"
)

var createJWKMigrations = []*migrate.Migration{
	{
		Id: "1-data",
		Up: []string{
			`INSERT INTO hydra_jwk (sid, kid, version, keydata) VALUES ('1-sid', '1-kid', 0, 'some-key')`,
		},
		Down: []string{
			`DELETE FROM hydra_jwk WHERE sid='1-sid'`,
		},
	},
	{
		Id: "2-data",
		Up: []string{
			`INSERT INTO hydra_jwk (sid, kid, version, keydata, created_at) VALUES ('2-sid', '2-kid', 0, 'some-key', NOW())`,
		},
		Down: []string{
			`DELETE FROM hydra_jwk WHERE sid='2-sid'`,
		},
	},
	{
		Id: "3-data",
		Up: []string{
			`INSERT INTO hydra_jwk (sid, kid, version, keydata, created_at) VALUES ('3-sid', '3-kid', 0, 'some-key', NOW())`,
		},
		Down: []string{
			`DELETE FROM hydra_jwk WHERE sid='3-sid'`,
		},
	},
	{
		Id: "4-data",
		Up: []string{
			`INSERT INTO hydra_jwk (sid, kid, version, keydata, created_at) VALUES ('4-sid', '4-kid', 0, 'some-key', NOW())`,
		},
		Down: []string{
			`DELETE FROM hydra_jwk WHERE sid='4-sid'`,
		},
	},
}

var migrations = map[string]*migrate.MemoryMigrationSource{
	"mysql": {
		Migrations: []*migrate.Migration{
			{Id: "0-data-0", Up: []string{"DROP TABLE IF EXISTS hydra_jwk"}},
			{Id: "0-data-1", Up: []string{"DROP TABLE IF EXISTS hydra_jwk_migration"}},
			jwk.Migrations["mysql"].Migrations[0],
			createJWKMigrations[0],
			jwk.Migrations["mysql"].Migrations[1],
			createJWKMigrations[1],
			jwk.Migrations["mysql"].Migrations[2],
			createJWKMigrations[2],
			jwk.Migrations["mysql"].Migrations[3],
			createJWKMigrations[3],
		},
	},
	"postgres": {
		Migrations: []*migrate.Migration{
			{Id: "0-data-0", Up: []string{"DROP TABLE IF EXISTS hydra_jwk"}},
			{Id: "0-data-1", Up: []string{"DROP TABLE IF EXISTS hydra_jwk_migration"}},
			jwk.Migrations["postgres"].Migrations[0],
			createJWKMigrations[0],
			jwk.Migrations["postgres"].Migrations[1],
			createJWKMigrations[1],
			jwk.Migrations["postgres"].Migrations[2],
			createJWKMigrations[2],
			jwk.Migrations["postgres"].Migrations[3],
			createJWKMigrations[3],
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
			migrate.SetTable("hydra_jwk_migration_integration")
			for step := range migrations[k].Migrations {
				t.Run(fmt.Sprintf("step=%d", step), func(t *testing.T) {
					n, err := migrate.ExecMax(db.DB, db.DriverName(), migrations[k], migrate.Up, 1)
					require.NoError(t, err)
					require.Equal(t, n, 1)
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
