// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"fmt"
	"strings"

	"github.com/ory/pop/v6"
)

type Order int

const (
	OrderAscending Order = iota
	OrderDescending
)

func (o Order) extract() (string, string) {
	switch o {
	case OrderAscending:
		return ">", "ASC"
	case OrderDescending:
		return "<", "DESC"
	default:
		panic(fmt.Sprintf("keyset pagination: unknown order %d", o))
	}
}

// Paginate returns a function that paginates a pop.Query.
// Usage:
//
//	q := c.Where("foo = ?", foo).Scope(keysetpagination.Paginate[MyItemType](paginator))
func Paginate[I any](p *Paginator) pop.ScopeFunc {
	model := pop.Model{Value: *new(I)}
	tableName := model.Alias()
	return func(q *pop.Query) *pop.Query {
		quoteAndContextualize := func(name string) string {
			quote := q.Connection.Dialect.Quote
			return quote(tableName) + "." + quote(name)
		}
		where, args, order := BuildWhereAndOrder(p.PageToken().Columns(), quoteAndContextualize)
		// IMPORTANT: Ensures correct query logic by grouping conditions.
		// Without parentheses, `WHERE otherCond AND pageCond1 OR pageCond2` would be
		// evaluated as `(otherCond = ? AND pageCond1) OR pageCond2`, potentially returning
		// rows that do not match `otherCond`.
		// We fix it by forcing the query to be: `WHERE otherCond AND (paginationCond1 OR paginationCond2)`.
		where = "(" + where + ")"

		return q.
			Where(where, args...).
			Order(order).
			Limit(p.Size() + 1)
	}
}

func BuildWhereAndOrder(columns []Column, quote func(string) string) (string, []any, string) {
	var whereBuilder, orderByBuilder, prevEqual strings.Builder
	args := make([]any, 0, len(columns)*(len(columns)+1)/2)
	prevEqualArgs := make([]any, 0, len(columns))

	whereBuilder.WriteRune('(')

	for i, part := range columns {
		column := quote(part.Name)
		sign, keyword := part.Order.extract()

		// Build query
		if i > 0 {
			whereBuilder.WriteString(") OR (")
		}
		whereBuilder.WriteString(prevEqual.String())
		if prevEqual.Len() > 0 {
			whereBuilder.WriteString(" AND ")
		}
		whereBuilder.WriteString(fmt.Sprintf("%s %s ?", column, sign))

		// Build orderBy
		if i > 0 {
			orderByBuilder.WriteString(", ")
		}
		orderByBuilder.WriteString(column + " " + keyword)

		// Update prevEqual
		if i > 0 {
			prevEqual.WriteString(" AND ")
		}
		prevEqual.WriteString(fmt.Sprintf("%s = ?", column))
		prevEqualArgs = append(prevEqualArgs, part.Value)
		args = append(args, prevEqualArgs...)
	}

	whereBuilder.WriteRune(')')

	return whereBuilder.String(), args, orderByBuilder.String()
}
