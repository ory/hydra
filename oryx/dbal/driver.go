// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	drivers = make([]func() Driver, 0)
	dmtx    sync.Mutex

	// ErrNoResponsibleDriverFound is returned when no driver was found for the provided DSN.
	ErrNoResponsibleDriverFound = errors.New("dsn value requested an unknown driver")
	ErrSQLiteSupportMissing     = errors.New(`the DSN connection string looks like a SQLite connection, but SQLite support was not built into the binary. Please check if you have downloaded the correct binary or are using the correct Docker Image. Binary archives and Docker Images indicate SQLite support by appending the -sqlite suffix`)
)

// Driver represents a driver
type Driver interface {
	// CanHandle returns true if the driver is capable of handling the given DSN or false otherwise.
	CanHandle(dsn string) bool

	// Ping returns nil if the driver has connectivity and is healthy or an error otherwise.
	Ping() error
	PingContext(context.Context) error
}

// RegisterDriver registers a driver
func RegisterDriver(d func() Driver) {
	dmtx.Lock()
	drivers = append(drivers, d)
	dmtx.Unlock()
}

// GetDriverFor returns a driver for the given DSN or ErrNoResponsibleDriverFound if no driver was found.
func GetDriverFor(dsn string) (Driver, error) {
	for _, f := range drivers {
		driver := f()
		if driver.CanHandle(dsn) {
			return driver, nil
		}
	}

	if IsSQLite(dsn) {
		return nil, ErrSQLiteSupportMissing
	}

	return nil, ErrNoResponsibleDriverFound
}

// IsSQLite returns true if the connection is a SQLite string.
func IsSQLite(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == "sqlite" || scheme == "sqlite3"
}
