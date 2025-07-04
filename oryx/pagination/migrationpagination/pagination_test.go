// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package migrationpagination

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ory/x/pagination/pagepagination"
	"github.com/ory/x/pagination/tokenpagination"

	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/urlx"
)

func TestPaginationHeader(t *testing.T) {
	u := urlx.ParseOrPanic("http://example.com")

	matches := func(t *testing.T, r *httptest.ResponseRecorder) {
		snapshotx.SnapshotT(t, strings.Split(r.Result().Header.Get("Link"), "; "))
	}

	t.Run("Create previous and first but not next or last if at the end", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 120, 2, 50)

		matches(t, r)
		assert.EqualValues(t, "120", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create next and last, but not previous or first if at the beginning", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 120, 0, 50)

		matches(t, r)
		assert.EqualValues(t, "120", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create previous, next, first, and last if in the middle", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 300, 3, 50)

		matches(t, r)
		assert.EqualValues(t, "300", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Header should default limit to 1 no limit was provided", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 100, 20, 0)

		matches(t, r)
		assert.EqualValues(t, "100", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create previous, next, first, but not last if in the middle and no total was provided", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 0, 3, 50)

		matches(t, r)
		assert.EqualValues(t, "0", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create only first if the limits provided exceeds the number of clients found", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 5, 0, 50)

		matches(t, r)
		assert.EqualValues(t, "5", r.Result().Header.Get("X-Total-Count"))
	})
}

func TestParsePagination(t *testing.T) {
	for _, tc := range []struct {
		d                    string
		url                  string
		expectedItemsPerPage int
		expectedPage         int
	}{
		{"normal", "http://localhost/foo?page_size=10&page_token=eyJvZmZzZXQiOjEwfQ", 10, 1},
		{"normal-encoded", fmt.Sprintf("http://localhost/foo?page_size=10&page_token=%s", tokenpagination.Encode(10)), 10, 1},
		{"defaults", "http://localhost/foo", 250, 0},
		{"limits", "http://localhost/foo?page_size=2000", 1000, 0},
		{"negatives", "http://localhost/foo?page_size=-1&page=eyJvZmZzZXQiOi0xfQ", 1, 0},
		{"negatives-encoded", fmt.Sprintf("http://localhost/foo?page_size=-1&page=%s", tokenpagination.Encode(-1)), 1, 0},
		{"invalid_params", "http://localhost/foo?page_size=a&page=b", 250, 0},
		{"legacy-normal", "http://localhost/foo?per_page=10&page=10", 10, 10},
		{"legacy-defaults", "http://localhost/foo", 250, 0},
		{"legacy-limits", "http://localhost/foo?per_page=2000", 1000, 0},
		{"legacy-negatives", "http://localhost/foo?per_page=-1&page=-1", 1, 0},
		{"legacy-invalid_params", "http://localhost/foo?per_page=a&page=b", 250, 0},
	} {
		t.Run(fmt.Sprintf("case=%s", tc.d), func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			page, perPage := NewPaginator(&pagepagination.PagePaginator{}, &tokenpagination.TokenPaginator{}).
				ParsePagination(&http.Request{URL: u})
			assert.EqualValues(t, tc.expectedItemsPerPage, perPage, "page_size")
			assert.EqualValues(t, tc.expectedPage, page, "page_token")
			assert.EqualValues(t, tc.expectedItemsPerPage, perPage, "per_page")
			assert.EqualValues(t, tc.expectedPage, page, "page")
		})
	}
}
