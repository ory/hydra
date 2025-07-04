// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	ID        int    `db:"pk"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
}

func nTestItems(n int) []testItem {
	items := make([]testItem, n)
	for i := range items {
		items[i] = testItem{
			ID:        i + 1,
			Name:      "item" + strconv.Itoa(i+1),
			CreatedAt: "2023-01-01T00:00:00Z",
		}
	}
	return items
}

func TestResult(t *testing.T) {
	t.Parallel()

	defaultToken := NewPageToken(Column{Name: "pk", Value: 0}, Column{Name: "name", Order: OrderDescending, Value: ""})
	paginator := NewPaginator(WithSize(10), WithDefaultToken(defaultToken))

	t.Run("not last page", func(t *testing.T) {
		items := nTestItems(11)
		croppedItems, nextPage := Result(items, paginator)
		assert.Len(t, croppedItems, 10)
		assert.Equal(t, 10, nextPage.Size())
		assert.False(t, nextPage.IsLast())
		assert.Equal(t, NewPageToken(
			Column{Name: "pk", Value: 10},
			Column{Name: "name", Order: OrderDescending, Value: items[9].Name},
		), nextPage.PageToken())
		assert.NotContains(t, croppedItems, items[10], "last item should not be included in the result")
		assert.Equal(t, croppedItems, items[:10], "cropped items should match the first 10 items")
	})

	t.Run("last page is full", func(t *testing.T) {
		items := nTestItems(10)
		croppedItems, nextPage := Result(items, paginator)
		assert.Len(t, croppedItems, 10)
		assert.Equal(t, 10, nextPage.Size())
		assert.True(t, nextPage.IsLast())
		assert.Equal(t, defaultToken, nextPage.PageToken())
		assert.Equal(t, croppedItems, items)
	})

	t.Run("last page not full", func(t *testing.T) {
		items := nTestItems(2)
		croppedItems, nextPage := Result(items, paginator)
		assert.Len(t, croppedItems, 2)
		assert.Equal(t, 10, nextPage.Size())
		assert.True(t, nextPage.IsLast())
		assert.Equal(t, defaultToken, nextPage.PageToken())
		assert.Equal(t, croppedItems, items)
	})
}

func TestPaginator_Size(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		opts     []Option
		expected int
	}{
		{
			name:     "default",
			opts:     nil,
			expected: DefaultSize,
		},
		{
			name:     "enforced default max size",
			opts:     []Option{WithSize(2 * DefaultMaxSize)},
			expected: DefaultMaxSize,
		},
		{
			name:     "with size",
			opts:     []Option{WithSize(10)},
			expected: 10,
		},
		{
			name:     "with custom default",
			opts:     []Option{WithDefaultSize(10)},
			expected: 10,
		},
		{
			name:     "with custom default and size",
			opts:     []Option{WithDefaultSize(10), WithSize(20)},
			expected: 20,
		},
		{
			name:     "with size and default bigger than max",
			opts:     []Option{WithSize(10), WithDefaultSize(20), WithMaxSize(5)},
			expected: 5,
		},
		{
			name:     "with negative size",
			opts:     []Option{WithSize(-1), WithDefaultSize(20), WithMaxSize(100)},
			expected: 20,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, NewPaginator(tc.opts...).Size())
		})
	}
}

func TestPaginator_Token(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		opts     []Option
		expected PageToken
	}{
		{
			name:     "no options",
			opts:     nil,
			expected: PageToken{},
		},
		{
			name:     "with token",
			opts:     []Option{WithToken(NewPageToken(Column{Name: "id", Value: "token"}))},
			expected: NewPageToken(Column{Name: "id", Value: "token"}),
		},
		{
			name:     "with default token",
			opts:     []Option{WithDefaultToken(NewPageToken(Column{Name: "id", Value: "default"}))},
			expected: NewPageToken(Column{Name: "id", Value: "default"}),
		},
		{
			name:     "with both tokens",
			opts:     []Option{WithToken(NewPageToken(Column{Name: "id", Value: "token"})), WithDefaultToken(NewPageToken(Column{Name: "id", Value: "default"}))},
			expected: NewPageToken(Column{Name: "id", Value: "token"}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			paginator := NewPaginator(tc.opts...)
			assert.Equal(t, tc.expected, paginator.PageToken())
		})
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		opts []Option
	}{
		{
			name: "no options",
			opts: nil,
		},
		{
			name: "with token",
			opts: []Option{WithToken(NewPageToken(Column{Name: "id", Value: "token"}))},
		},
		{
			name: "with size",
			opts: []Option{WithSize(10)},
		},
		{
			name: "with all options",
			opts: []Option{
				WithSize(20),
				WithDefaultSize(30),
				WithMaxSize(50),
				WithToken(NewPageToken(Column{Name: "id", Value: 123})),
				WithDefaultToken(NewPageToken(Column{Name: "id", Value: 456})),
				withIsLast(true),
			},
		},
		{
			name: "with explicit defaults",
			opts: []Option{WithMaxSize(DefaultMaxSize), WithDefaultSize(DefaultSize)},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			paginator := NewPaginator(tc.opts...)
			assert.Equal(t, paginator, NewPaginator(paginator.ToOptions()...))
		})
	}
}
