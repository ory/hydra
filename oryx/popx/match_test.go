// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseMigrationFilenameSQLUp(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.up.sql")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "all")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "sql")
	r.Equal(m.Autocommit, false)
}

func Test_ParseMigrationFilenameSQLUpPostgres(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.pg.up.sql")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "postgres")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "sql")
	r.Equal(m.Autocommit, false)
}

func Test_ParseMigrationFilenameSQLUpAutocommit(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.autocommit.up.sql")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "all")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "sql")
	r.Equal(m.Autocommit, true)
}

func Test_ParseMigrationFilenameSQLDownAutocommit(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.mysql.autocommit.down.sql")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "mysql")
	r.Equal(m.Direction, "down")
	r.Equal(m.Type, "sql")
	r.Equal(m.Autocommit, true)
}
