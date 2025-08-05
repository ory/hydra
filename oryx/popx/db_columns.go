// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
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

func WritableDBColumnNames[T any]() []string {
	var names []string
	for _, c := range (&pop.Model{Value: new(T)}).Columns().Writeable().Cols {
		names = append(names, c.Name)
	}
	return names
}

func DBColumnsExcluding[T any](quoter Quoter, exclude ...string) string {
	cols := (&pop.Model{Value: new(T)}).Columns()
	for _, e := range exclude {
		cols.Remove(e)
	}
	return cols.QuotedString(quoter)
}

type (
	PrefixQuoter struct {
		Prefix string
		Quoter Quoter
	}
	Quoter interface {
		Quote(key string) string
	}
)

func (pq *PrefixQuoter) Quote(key string) string {
	return pq.Quoter.Quote(pq.Prefix + key)
}
