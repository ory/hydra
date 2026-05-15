// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package landlockx

import (
	"path/filepath"
	"strings"

	"github.com/ory/x/sqlxx"
)

// SQLiteDirFromDSN returns the directory containing the SQLite database
// referenced by the DSN. Granting the directory once is simpler and more
// robust than enumerating every sibling: SQLite manages -journal in
// rollback mode and -wal/-shm in WAL mode, plus transient -mj-XXXXX
// rollback masters during multi-database checkpoints — the latter use
// random suffixes that a per-file rule cannot cover.
//
// Returns an empty string for non-SQLite DSNs and for the in-memory
// SQLite DSN.
func SQLiteDirFromDSN(dsn string) string {
	scheme, path, err := sqlxx.ExtractSchemeFromDSN(dsn)
	if err != nil || (scheme != "sqlite" && scheme != "sqlite3") {
		return ""
	}
	path, _, _ = strings.Cut(path, "?")
	if path == "" || strings.Contains(path, ":memory:") {
		return ""
	}
	return filepath.Dir(path)
}
