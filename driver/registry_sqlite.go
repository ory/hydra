// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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
