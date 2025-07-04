// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"
)

func TestBuildWhereAndOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		parts []Column

		expectedWhere   string
		expectedArgs    []any
		expectedOrderBy string
	}{
		{
			name: "single part ascending",
			parts: []Column{
				{Name: "id", Order: OrderAscending, Value: "first"},
			},
			expectedWhere:   "(id > ?)",
			expectedArgs:    []any{"first"},
			expectedOrderBy: "id ASC",
		},
		{
			name: "single part descending",
			parts: []Column{
				{Name: "id", Order: OrderDescending, Value: 1},
			},
			expectedWhere:   "(id < ?)",
			expectedArgs:    []any{1},
			expectedOrderBy: "id DESC",
		},
		{
			name: "two cols",
			parts: []Column{
				{Name: "id", Order: OrderAscending, Value: 1},
				{Name: "name", Order: OrderDescending, Value: "test"},
			},
			expectedWhere:   "(id > ?) OR (id = ? AND name < ?)",
			expectedArgs:    []any{1, 1, "test"},
			expectedOrderBy: "id ASC, name DESC",
		},
		{
			name: "many cols",
			parts: []Column{
				{Name: "id", Order: OrderAscending, Value: 1},
				{Name: "name", Order: OrderAscending, Value: "test"},
				{Name: "created_at", Order: OrderDescending, Value: "2023-01-01"},
				{Name: "owner_id", Order: OrderDescending, Value: "owner123"},
			},
			expectedWhere:   "(id > ?) OR (id = ? AND name > ?) OR (id = ? AND name = ? AND created_at < ?) OR (id = ? AND name = ? AND created_at = ? AND owner_id < ?)",
			expectedArgs:    []any{1, 1, "test", 1, "test", "2023-01-01", 1, "test", "2023-01-01", "owner123"},
			expectedOrderBy: "id ASC, name ASC, created_at DESC, owner_id DESC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			where, args, order := BuildWhereAndOrder(tc.parts, func(s string) string { return s })
			assert.Equal(t, tc.expectedWhere, where)
			assert.Equal(t, tc.expectedArgs, args)
			assert.Equal(t, tc.expectedOrderBy, order)
		})
	}
}

func TestPaginate(t *testing.T) {
	t.Parallel()

	t.Run("paginates correctly", func(t *testing.T) {
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL: "postgres://foo.bar",
		})
		require.NoError(t, err)
		q := pop.Q(c)
		paginator := NewPaginator(WithSize(10), WithToken(NewPageToken(Column{Name: "pk", Value: 666})))
		q = q.Scope(Paginate[testItem](paginator))

		sql, args := q.ToSQL(&pop.Model{Value: new(testItem)})
		assert.Equal(t, `SELECT test_items.created_at, test_items.name, test_items.pk FROM test_items AS test_items WHERE ("test_items"."pk" > $1) ORDER BY "test_items"."pk" ASC LIMIT 11`, sql)
		assert.Equal(t, []interface{}{666}, args)
	})

	t.Run("paginates correctly with negative size", func(t *testing.T) {
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL: "postgres://foo.bar",
		})
		require.NoError(t, err)
		q := pop.Q(c)
		paginator := NewPaginator(WithSize(-1), WithDefaultSize(10), WithToken(NewPageToken(Column{Name: "pk", Value: 123})))
		q = q.Scope(Paginate[testItem](paginator))

		sql, args := q.ToSQL(&pop.Model{Value: new(testItem)})
		assert.Equal(t, `SELECT test_items.created_at, test_items.name, test_items.pk FROM test_items AS test_items WHERE ("test_items"."pk" > $1) ORDER BY "test_items"."pk" ASC LIMIT 11`, sql)
		assert.Equal(t, []interface{}{123}, args)
	})

	t.Run("paginates correctly mysql", func(t *testing.T) {
		c, err := pop.NewConnection(&pop.ConnectionDetails{
			URL: "mysql://user:pass@(host:1337)/database",
		})
		require.NoError(t, err)
		q := pop.Q(c)
		q = q.Scope(Paginate[testItem](NewPaginator(WithSize(10), WithToken(NewPageToken(Column{Name: "pk", Value: 666})))))

		sql, args := q.ToSQL(&pop.Model{Value: new(testItem)})
		assert.Equal(t, "SELECT test_items.created_at, test_items.name, test_items.pk FROM test_items AS test_items WHERE (`test_items`.`pk` > ?) ORDER BY `test_items`.`pk` ASC LIMIT 11", sql)
		assert.Equal(t, []interface{}{666}, args)
	})
}
