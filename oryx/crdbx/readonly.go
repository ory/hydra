// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package crdbx

import (
	"github.com/ory/pop/v6"

	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
)

// SetTransactionReadOnly sets the transaction to read only for CockroachDB.
func SetTransactionReadOnly(c *pop.Connection) error {
	if c.Dialect.Name() != dbal.DriverCockroachDB {
		// Only CockroachDB supports this.
		return nil
	}

	return sqlcon.HandleError(c.RawQuery("SET TRANSACTION READ ONLY").Exec())
}
