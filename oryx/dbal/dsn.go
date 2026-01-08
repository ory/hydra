// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

var sqliteMemoryRegexp = regexp.MustCompile(`^sqlite://file:.+\?.*&?mode=memory($|&.*)|sqlite://(file:)?:memory:\?.*$|^(:memory:|memory)$`)

// IsMemorySQLite returns true if a given DSN string is pointing to a SQLite database.
//
// SQLite can be written in different styles depending on the use case
// - just in memory
// - shared connection
// - shared but unique in the same process
// see: https://sqlite.org/inmemorydb.html
func IsMemorySQLite(dsn string) bool { return sqliteMemoryRegexp.MatchString(dsn) }

// NewSQLiteTestDatabase creates a new, unique SQLite database
// which is shared amongst all callers and identified by an individual file name.
// The database file is created in the system's temporary directory, and not actively
// removed to allow debugging in case of test failures.
func NewSQLiteTestDatabase(t testing.TB) string {
	fn, err := os.MkdirTemp("", "sqlite-test-db-*")
	require.NoError(t, err)
	return NewSQLiteDatabase(fn)
}

// NewSQLiteInMemoryDatabase creates a new in-memory, unique SQLite database
// which is shared amongst all callers and identified by an individual file name.
func NewSQLiteInMemoryDatabase(name string) string {
	return fmt.Sprintf("sqlite://file:%s?_fk=true&mode=memory&cache=shared&_busy_timeout=100000", name)
}

// NewSQLiteDatabase creates a new on-disk, unique SQLite database
// which is shared amongst all callers and identified by an individual file name.
// This is sometimes necessary over a in-memory database, for example when multiple tests/goroutines run in parallel
// and access the same table.
// This would result in a locking error from SQLite when running in-memory.
// Additionally, shared cache mode is deprecated and discouraged, and the problem is better solved with the WAL,
// according to official docs.
func NewSQLiteDatabase(name string) string {
	return fmt.Sprintf("sqlite://file:%s/db.sqlite?_fk=true&_journal_mode=WAL&_busy_timeout=100000", name)
}
