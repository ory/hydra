// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"database/sql/driver"
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
		dialect := q.Connection.Dialect.Name()
		where, args, order := BuildWhereAndOrder(p.PageToken().Columns(), quoteAndContextualize, dialect)
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

func BuildWhereAndOrder(columns []Column, quote func(string) string, dialect string) (string, []any, string) {
	var whereBuilder, orderByBuilder, prevEqual strings.Builder

	keysetCols := make([]Column, 0, len(columns))

	// ORDER BY includes all columns (even constrained ones)
	for i, part := range columns {
		column := quote(part.Name)
		_, keyword := part.Order.extract()

		if i > 0 {
			orderByBuilder.WriteString(", ")
		}

		orderByBuilder.WriteString(column + " " + keyword)

		// Postgres orders NULLs differently depending on sort direction;
		// (ASC → NULLS LAST, DESC → NULLS FIRST), which does not match the
		// assumptions of our keyset pagination logic and other supported DBs.
		// We therefore make NULL ordering explicit on Postgres to keep pagination
		// stable and consistent with sqlite/mysql/cockroachdb.
		if dialect == "postgres" && part.Nullable {
			if part.Order == OrderAscending {
				orderByBuilder.WriteString(" NULLS FIRST")
			} else {
				orderByBuilder.WriteString(" NULLS LAST")
			}
		}

		// Build keyset WHERE only from unconstrained columns
		if !part.HasConstraint {
			keysetCols = append(keysetCols, part)
		}
	}

	// If everything is constrained, no keyset predicate is needed.
	if len(keysetCols) == 0 {
		return "", nil, orderByBuilder.String()
	}

	args := make([]any, 0, len(keysetCols)*(len(keysetCols)+1)/2)
	prevEqualArgs := make([]any, 0, len(keysetCols))

	whereBuilder.WriteRune('(')

	for i, part := range keysetCols {
		column := quote(part.Name)
		sign, _ := part.Order.extract()

		if i > 0 {
			whereBuilder.WriteString(") OR (")
		}

		whereBuilder.WriteString(prevEqual.String())
		if prevEqual.Len() > 0 {
			whereBuilder.WriteString(" AND ")
		}

		isNull := part.Nullable && isSQLNull(part.Value)

		if !part.Nullable {
			whereBuilder.WriteString(column + " " + sign + " ?")
		} else if !isNull {
			whereBuilder.WriteString(column + " IS NOT NULL AND " + column + " " + sign + " ?")
		} else {
			whereBuilder.WriteString(column + " IS NOT NULL")
		}

		if i > 0 {
			prevEqual.WriteString(" AND ")
		}

		if !part.Nullable {
			prevEqual.WriteString(column + " = ?")
			prevEqualArgs = append(prevEqualArgs, part.Value)
		} else if !isNull {
			prevEqual.WriteString(column + " = ?")
			prevEqualArgs = append(prevEqualArgs, part.Value)
		} else {
			prevEqual.WriteString(column + " IS NULL")
		}

		args = append(args, prevEqualArgs...)
	}

	whereBuilder.WriteRune(')')

	return whereBuilder.String(), args, orderByBuilder.String()
}

// isSQLNull reports whether v represents a SQL NULL value.
func isSQLNull(v any) bool {
	if v == nil {
		return true
	}

	if valuer, ok := v.(driver.Valuer); ok {
		val, err := valuer.Value()
		return err == nil && val == nil
	}

	return false
}
