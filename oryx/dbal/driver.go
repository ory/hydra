// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"context"
	"errors"
	"strings"
)

var (
	// ErrNoResponsibleDriverFound is returned when no driver was found for the provided DSN.
	ErrNoResponsibleDriverFound = errors.New("dsn value requested an unknown driver")
	ErrSQLiteSupportMissing     = errors.New(`the DSN connection string looks like a SQLite connection, but SQLite support was not built into the binary. Please check if you have downloaded the correct binary or are using the correct Docker Image. Binary archives and Docker Images indicate SQLite support by appending the -sqlite suffix`)
)

// Driver represents a driver
type Driver interface {
	// Ping returns nil if the driver has connectivity and is healthy or an error otherwise.
	Ping() error
	PingContext(context.Context) error
}

// IsSQLite returns true if the connection is a SQLite string.
func IsSQLite(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == "sqlite" || scheme == "sqlite3"
}
