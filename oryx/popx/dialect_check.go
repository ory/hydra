// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"
	"github.com/ory/x/dbal"
)

// VerifyDialect makes sure the declared SQL dialect matches the actual
// database server. It distinguishes PostgreSQL, CockroachDB, and YugabyteDB:
// all three speak the Postgres wire protocol, so an operator can easily point
// a service configured for one at a server running another. That mismatch
// silently selects the wrong dialect-specific migrations and produces broken
// schemas.
//
// Other dialects (MySQL, SQLite) are skipped.
func VerifyDialect(ctx context.Context, conn *pop.Connection) error {
	declared := conn.Dialect.Name()
	if declared != namePostgres && declared != dbal.DriverCockroachDB && declared != dbal.DriverYugabyteDB {
		return nil
	}

	var version string
	if err := conn.WithContext(ctx).RawQuery("SELECT version() AS version").First(&version); err != nil {
		return errors.Wrap(err, "could not query database version to verify dialect")
	}
	return checkDialect(declared, version, pop.DialectSupported(dbal.DriverYugabyteDB))
}

const (
	namePostgres = "postgres"

	// yugabyteVersionMarker identifies a YugabyteDB server in the output of
	// SELECT version() (e.g. "PostgreSQL 11.2-YB-2025.2.4.0-b122 ...").
	yugabyteVersionMarker = "-YB-"
)

func checkDialect(declared, version string, yugabyteDialectSupported bool) error {
	detectedCockroach := strings.Contains(version, "CockroachDB")
	detectedYugabyte := strings.Contains(version, yugabyteVersionMarker)
	switch {
	case declared == dbal.DriverCockroachDB && !detectedCockroach:
		return errors.Errorf(
			"DSN scheme cockroach:// declares a CockroachDB database but the server is not CockroachDB. Server reported: %q. Use a postgres:// DSN, or point the service at a CockroachDB cluster.",
			firstLine(version),
		)
	case declared == dbal.DriverYugabyteDB && !detectedYugabyte:
		return errors.Errorf(
			"DSN scheme yugabyte:// declares a YugabyteDB database but the server is not YugabyteDB. Server reported: %q. Use a DSN scheme matching the actual server, or point the service at a YugabyteDB cluster.",
			firstLine(version),
		)
	case declared == namePostgres && detectedCockroach:
		return errors.Errorf(
			"DSN scheme postgres:// declares a PostgreSQL database but the server is CockroachDB. Replace the scheme with cockroach:// so that the service picks the correct migrations and SQL dialect. Server reported: %q",
			firstLine(version),
		)
	case declared == namePostgres && detectedYugabyte && yugabyteDialectSupported:
		return errors.Errorf(
			"DSN scheme postgres:// declares a PostgreSQL database but the server is YugabyteDB. Replace the scheme with yugabyte:// (supported in Ory Enterprise License builds) so that the service picks the correct migrations. Server reported: %q",
			firstLine(version),
		)
	}
	return nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
