//go:build sqlite
// +build sqlite

package driver

import (
	"strings"
)

func (m *RegistrySQL) CanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == "sqlite" || scheme == "sqlite3" || m.alwaysCanHandle(dsn)
}
