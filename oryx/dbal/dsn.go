// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"fmt"
	"os"
	"regexp"
)

const (
	// SQLiteInMemory is a DNS string for SQLite in-memory database.
	//
	// DEPRECATED: Do not use this DSN string as it can cause flaky tests
	// due to the way SQL connection pooling works. Please use NewSQLiteTestDatabase instead.
	SQLiteInMemory = "sqlite://file::memory:?_fk=true"
	// SQLiteSharedInMemory is a DNS string for SQLite in-memory database in shared mode.
	//
	// DEPRECATED: Do not use this DSN string as it can cause flaky tests
	// due to the way SQL connection pooling works. Please use NewSQLiteTestDatabase instead.
	SQLiteSharedInMemory = "sqlite://file::memory:?_fk=true&cache=shared"
)

var dsnRegex = regexp.MustCompile(`^(sqlite://file:(?:.+)\?((\w+=\w+)(&\w+=\w+)*)?(&?mode=memory)(&\w+=\w+)*)$|(?:sqlite://(file:)?:memory:(?:\?\w+=\w+)?(?:&\w+=\w+)*)|^(?:(?::memory:)|(?:memory))$`)

// IsMemorySQLite returns true if a given DSN string is pointing to a SQLite database.
//
// SQLite can be written in different styles depending on the use case
// - just in memory
// - shared connection
// - shared but unique in the same process
// see: https://sqlite.org/inmemorydb.html
func IsMemorySQLite(dsn string) bool {
	return dsnRegex.MatchString(dsn)
}

// NewSharedUniqueInMemorySQLiteDatabase creates a new unique SQLite database
// which is shared amongst all callers and identified by an individual file name.
//
// DEPRECATED: Please use NewSQLiteTestDatabase instead.
func NewSharedUniqueInMemorySQLiteDatabase() (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "unique-sqlite-db-*")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sqlite://file:%s/db.sqlite?_fk=true&mode=memory&cache=shared", dir), nil
}

// NewSQLiteTestDatabase creates a new unique SQLite database
// which is shared amongst all callers and identified by an individual file name.
func NewSQLiteTestDatabase(t interface {
	TempDir() string
}) string {
	return NewSQLiteInMemoryDatabase(t.TempDir())
}

// NewSQLiteInMemoryDatabase creates a new unique SQLite database
// which is shared amongst all callers and identified by an individual file name.
func NewSQLiteInMemoryDatabase(name string) string {
	return fmt.Sprintf("sqlite://file:%s?_fk=true&mode=memory&cache=shared", name)
}
