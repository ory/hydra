// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package crdbx

import (
	"github.com/ory/pop/v6"

	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
)

// SetTransactionReadOnly sets the current transaction to read only for
// CockroachDB and PostgreSQL, which both let the server skip write-path
// bookkeeping and reject accidental writes. Must be called before the
// transaction runs its first query. MySQL's SET TRANSACTION applies to the
// next transaction rather than the current one, and SQLite has no
// equivalent, so every other dialect is a no-op.
func SetTransactionReadOnly(c *pop.Connection) error {
	switch c.Dialect.Name() {
	case dbal.DriverCockroachDB, dbal.DriverPostgreSQL:
		return sqlcon.HandleError(c.RawQuery("SET TRANSACTION READ ONLY").Exec())
	default:
		return nil
	}
}
