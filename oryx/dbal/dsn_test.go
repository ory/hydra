// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMemorySQLite(t *testing.T) {
	testCases := map[string]bool{
		SQLiteInMemory:               true,
		SQLiteSharedInMemory:         true,
		"memory":                     true,
		":memory:":                   true,
		"sqlite://:memory:?_fk=true": true,
		"sqlite://file:uniquedb:?_fk=true&mode=memory":              true,
		"sqlite://file:uniquedb:?_fk=true&mode=memory&cache=shared": true,
		"sqlite://file:uniquedb:?_fk=true&cache=shared&mode=memory": true,
		"sqlite://file:uniquedb:?mode=memory":                       true,
		"sqlite://file:::uniquedb:?_fk=true&mode=memory":            true,
		"sqlite://file:memdb1?mode=memory&cache=shared":             true,
		"sqlite://file:uniquedb:?_fk=true&cache=shared":             false,
		"sqlite://":                                            false,
		"sqlite://file":                                        false,
		"sqlite://file:::":                                     false,
		"sqlite://?_fk=true&mode=memory":                       false,
		"sqlite://?_fk=true&cache=shared":                      false,
		"sqlite://file::?_fk=true":                             false,
		"sqlite://file:::?_fk=true":                            false,
		"postgresql://username:secret@localhost:5432/database": false,
	}

	for dsn, expected := range testCases {
		t.Run("dsn="+dsn, func(t *testing.T) {
			assert.Equal(t, expected, IsMemorySQLite(dsn))
		})
	}
}
