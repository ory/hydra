// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"strings"
)

// IsSQLite returns true if the connection is a SQLite string.
func IsSQLite(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == "sqlite" || scheme == "sqlite3"
}
