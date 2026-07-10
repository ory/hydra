// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package crdbx

import (
	"github.com/ory/pop/v6"
	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
)

// SetTransactionReadCommitted lowers the isolation of the current transaction
// to READ COMMITTED on CockroachDB. Multi-statement read transactions at
// SERIALIZABLE must refresh their whole read set when a concurrent write
// bumps the transaction timestamp, and surface RETRY_SERIALIZABLE to the
// client when the refresh fails; under READ COMMITTED every statement reads
// at its own timestamp and those conflicts are absorbed server-side. Must be
// called before the transaction runs its first query. Postgres already runs
// statements at READ COMMITTED by default, and MySQL's SET TRANSACTION
// applies to the next transaction rather than the current one, so every
// other dialect is a no-op.
func SetTransactionReadCommitted(c *pop.Connection) error {
	if c.Dialect.Name() != dbal.DriverCockroachDB {
		// Only CockroachDB supports and needs this.
		return nil
	}

	return sqlcon.HandleError(c.RawQuery("SET TRANSACTION ISOLATION LEVEL READ COMMITTED").Exec())
}
