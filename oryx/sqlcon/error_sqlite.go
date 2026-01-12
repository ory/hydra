// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build sqlite
// +build sqlite

package sqlcon

import (
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// handleSqlite handles the error iff (if and only if) it is an sqlite error
func handleSqlite(err error) error {
	if e := new(sqlite3.Error); errors.As(err, e) {
		switch e.ExtendedCode {
		case sqlite3.ErrConstraintUnique:
			fallthrough
		case sqlite3.ErrConstraintPrimaryKey:
			return errors.WithStack(ErrUniqueViolation.WithWrap(err))

		}

		switch e.Code {
		case sqlite3.ErrError:
			if strings.Contains(err.Error(), "no such table") {
				return errors.WithStack(ErrNoSuchTable.WithWrap(err))
			}
		case sqlite3.ErrLocked:
			return errors.WithStack(ErrConcurrentUpdate.WithWrap(err))
		}

		if strings.HasPrefix(e.Error(), "no such column:") ||
			strings.Contains(e.Error(), "has no column named") {
			return errors.WithStack(ErrNoSuchColumn.WithWrap(err))
		}

		return errors.WithStack(err)
	}

	return nil
}
