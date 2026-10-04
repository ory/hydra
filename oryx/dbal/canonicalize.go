// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dbal

const (
	DriverMySQL       = "mysql"
	DriverCockroachDB = "cockroach"
	DriverPostgreSQL  = "postgres"
	// DriverYugabyteDB is the dialect name YugabyteDB connections report. It is
	// a commercial (Ory Enterprise License) dialect that speaks the PostgreSQL
	// wire protocol and SQL dialect, so most SQL paths treat it like postgres.
	DriverYugabyteDB = "yugabyte"
)

// IsPostgresCompatible reports whether dialect speaks the PostgreSQL SQL
// dialect. It does not imply support for CockroachDB-specific syntax or
// features.
func IsPostgresCompatible(dialect string) bool {
	switch dialect {
	case DriverPostgreSQL, DriverCockroachDB, DriverYugabyteDB:
		return true
	default:
		return false
	}
}
