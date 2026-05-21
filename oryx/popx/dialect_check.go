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
// database server. It only distinguishes PostgreSQL from CockroachDB,
// because CockroachDB speaks the Postgres wire protocol and advertises a
// "postgres://" DSN — an operator can easily configure a service as
// Postgres against a CockroachDB cluster. That mismatch silently selects
// the wrong dialect-specific migrations and produces broken schemas.
//
// Other dialects (MySQL, SQLite) are skipped.
func VerifyDialect(ctx context.Context, conn *pop.Connection) error {
	declared := conn.Dialect.Name()
	if declared != namePostgres && declared != dbal.DriverCockroachDB {
		return nil
	}

	var version string
	if err := conn.WithContext(ctx).RawQuery("SELECT version() AS version").First(&version); err != nil {
		return errors.Wrap(err, "could not query database version to verify dialect")
	}
	return checkDialect(declared, version)
}

const namePostgres = "postgres"

func checkDialect(declared, version string) error {
	detectedCockroach := strings.Contains(version, "CockroachDB")
	switch {
	case declared == dbal.DriverCockroachDB && !detectedCockroach:
		return errors.Errorf(
			"DSN scheme cockroach:// declares a CockroachDB database but the server is not CockroachDB. Server reported: %q. Use a postgres:// DSN, or point the service at a CockroachDB cluster.",
			firstLine(version),
		)
	case declared == namePostgres && detectedCockroach:
		return errors.Errorf(
			"DSN scheme postgres:// declares a PostgreSQL database but the server is CockroachDB. Replace the scheme with cockroach:// so that the service picks the correct migrations and SQL dialect. Server reported: %q",
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
