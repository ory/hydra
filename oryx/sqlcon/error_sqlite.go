// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlcon

import (
	"strings"

	"github.com/pkg/errors"
	"modernc.org/sqlite"
	sqlite3lib "modernc.org/sqlite/lib"
)

// handleSqlite handles the error iff (if and only if) it is an sqlite error
func handleSqlite(e *sqlite.Error, err error) error {
	// Code() returns the full extended error code.
	switch e.Code() {
	case sqlite3lib.SQLITE_CONSTRAINT_UNIQUE,
		sqlite3lib.SQLITE_CONSTRAINT_PRIMARYKEY:
		return errors.WithStack(ErrUniqueViolation().WithWrap(err))
	case sqlite3lib.SQLITE_ERROR:
		if strings.Contains(e.Error(), "no such table") {
			return errors.WithStack(ErrNoSuchTable().WithWrap(err))
		}
	case sqlite3lib.SQLITE_LOCKED,
		sqlite3lib.SQLITE_BUSY,
		sqlite3lib.SQLITE_BUSY_RECOVERY,
		sqlite3lib.SQLITE_BUSY_SNAPSHOT,
		sqlite3lib.SQLITE_BUSY_TIMEOUT:
		return errors.WithStack(ErrConcurrentUpdate().WithWrap(err))
	}

	if strings.HasPrefix(e.Error(), "no such column:") ||
		strings.Contains(e.Error(), "has no column named") {
		return errors.WithStack(ErrNoSuchColumn().WithWrap(err))
	}

	return errors.WithStack(err)
}
