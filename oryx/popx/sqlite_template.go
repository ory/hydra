// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"
	"github.com/ory/x/logrusx"
)

// sqliteFilePath returns the absolute file path for an on-disk SQLite DSN.
// SQLite memory databases and non-SQLite DSNs return ("", false).
func sqliteFilePath(dsn string) (string, bool) {
	rest := dsn
	if scheme, afterScheme, ok := strings.Cut(dsn, "://"); ok {
		if scheme != "sqlite" && scheme != "sqlite3" {
			return "", false
		}
		rest = afterScheme
	}
	if strings.Contains(rest, ":memory:") || strings.Contains(rest, "mode=memory") {
		return "", false
	}
	rest = strings.TrimPrefix(rest, "file:")
	if i := strings.IndexByte(rest, '?'); i >= 0 {
		rest = rest[:i]
	}
	if !filepath.IsAbs(rest) {
		return "", false
	}
	return rest, true
}

// PrepareTestSQLiteDatabase migrates a fresh on-disk SQLite test database
// before the application opens it. MigrationBox's shared template cache makes
// this a cheap template restore after the first migration set has been built.
//
// This is useful when an application normally applies multiple migration sets
// in separate passes: pass their merged filesystem here so the cached template
// represents the final combined schema. Non-SQLite DSNs, memory databases, and
// existing database files are left untouched.
func PrepareTestSQLiteDatabase(
	ctx context.Context,
	dsn string,
	migrations fs.FS,
	logger *logrusx.Logger,
	opts ...MigrationBoxOption,
) error {
	if !testing.Testing() {
		return nil
	}
	path, ok := sqliteFilePath(dsn)
	if !ok {
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return errors.WithStack(err)
	}

	c, err := pop.NewConnection(&pop.ConnectionDetails{URL: dsn})
	if err != nil {
		return errors.WithStack(err)
	}
	if err := c.Open(); err != nil {
		return errors.WithStack(err)
	}

	mb, err := NewMigrationBox(migrations, c, logger, opts...)
	if err == nil {
		err = mb.Up(ctx)
	}
	closeErr := c.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(closeErr)
}
