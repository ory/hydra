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

package jwk

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/x/dbal"
	"github.com/ory/x/dbal/migratest"
)

var createMigrations = map[string]*dbal.PackrMigrationSource{
	dbal.DriverMySQL:      dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/tests"}, true),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{"migrations/sql/tests"}, true),
}

func TestXXMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	require.True(t, len(client.Migrations[dbal.DriverMySQL].Box.List()) == len(client.Migrations[dbal.DriverPostgreSQL].Box.List()))

	var clean = func(t *testing.T, db *sqlx.DB) {
		_, err := db.Exec("DROP TABLE IF EXISTS hydra_jwk")
		require.NoError(t, err)
		_, err = db.Exec("DROP TABLE IF EXISTS hydra_jwk_migration")
		require.NoError(t, err)
	}

	migratest.RunPackrMigrationTests(
		t,
		migratest.MigrationSchemas{migrations},
		migratest.MigrationSchemas{createMigrations},
		clean, clean,
		func(t *testing.T, db *sqlx.DB, k, m, steps int) {
			t.Run(fmt.Sprintf("poll=%d", k), func(t *testing.T) {
				sid := fmt.Sprintf("%d-sid", k+1)
				m := NewSQLManager(db, []byte("01234567890123456789012345678912"))
				_, err := m.GetKeySet(context.TODO(), sid)
				require.Error(t, err, "malformed ciphertext")
			})
		},
	)
}
