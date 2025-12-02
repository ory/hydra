// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"fmt"

	"github.com/ory/pop/v6"
)

func DBColumns[T any](quoter Quoter) string {
	return (&pop.Model{Value: new(T)}).Columns().QuotedString(quoter)
}

// IndexHint returns the table name including the index hint, if the database
// supports it.
func IndexHint(conn *pop.Connection, table string, index string) string {
	if conn.Dialect.Name() == "cockroach" {
		return table + "@" + index
	}
	return table
}

func DBColumnsExcluding[T any](quoter Quoter, exclude ...string) string {
	cols := (&pop.Model{Value: new(T)}).Columns()
	for _, e := range exclude {
		cols.Remove(e)
	}
	return cols.QuotedString(quoter)
}

type (
	AliasQuoter struct {
		Alias  string
		Quoter Quoter
	}
	Quoter interface {
		Quote(key string) string
	}
)

func (pq *AliasQuoter) Quote(key string) string {
	return fmt.Sprintf("%s.%s", pq.Quoter.Quote(pq.Alias), pq.Quoter.Quote(key))
}
