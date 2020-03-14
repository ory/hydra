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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(Migrations[dbal.DriverMySQL].Box.List()) == len(Migrations[dbal.DriverPostgreSQL].Box.List()))

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{Migrations},
		migratest.MigrationSchemas{dbal.FindMatchingTestMigrations("migrations/sql/tests/", Migrations, AssetNames(), Asset)},
		x.CleanSQL, x.CleanSQL,
		func(t *testing.T, dbName string, db *sqlx.DB, _, step, steps int) {
			if dbName == "cockroach" {
				step += 12
			}
			id := fmt.Sprintf("%d-data", step+1)
			t.Run("poll="+id, func(t *testing.T) {
				s := NewSQLManager(db, nil)
				c, err := s.GetConcreteClient(context.TODO(), id)
				require.NoError(t, err)
				assert.EqualValues(t, c.GetID(), id)
			})
		},
	)
}
