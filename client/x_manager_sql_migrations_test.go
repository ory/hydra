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

package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

var createMigrations = map[string]*dbal.PackrMigrationSource{
	dbal.DriverMySQL:      dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/tests"}, true),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/tests"}, true),
}

func CleanTestDB(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS hydra_client_migration")
	require.NoError(t, err)
	_, err = db.Exec("DROP TABLE IF EXISTS hydra_client")
	require.NoError(t, err)
}

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(Migrations[dbal.DriverMySQL].Box.List()) == len(Migrations[dbal.DriverPostgreSQL].Box.List()))

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{Migrations},
		migratest.MigrationSchemas{createMigrations},
		CleanTestDB, CleanTestDB,
		func(t *testing.T, db *sqlx.DB, _, step, steps int) {
			id := fmt.Sprintf("%d-data", step+1)
			t.Run("poll="+id, func(t *testing.T) {
				s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}
				c, err := s.GetConcreteClient(context.TODO(), id)
				require.NoError(t, err)
				assert.EqualValues(t, c.GetID(), id)
			})
		},
	)
}
