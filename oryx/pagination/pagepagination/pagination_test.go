// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagepagination

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/urlx"
)

func TestPaginationHeader(t *testing.T) {
	u := urlx.ParseOrPanic("http://example.com")

	t.Run("Create previous and first but not next or last if at the end", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 120, 2, 50)

		expect := strings.Join([]string{
			"<http://example.com?page=0&per_page=50>; rel=\"first\"",
			"<http://example.com?page=1&per_page=50>; rel=\"prev\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "120", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create next and last, but not previous or first if at the beginning", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 120, 0, 50)

		expect := strings.Join([]string{
			"<http://example.com?page=1&per_page=50>; rel=\"next\"",
			"<http://example.com?page=2&per_page=50>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "120", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create previous, next, first, and last if in the middle", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 300, 3, 50)

		expect := strings.Join([]string{
			"<http://example.com?page=0&per_page=50>; rel=\"first\"",
			"<http://example.com?page=4&per_page=50>; rel=\"next\"",
			"<http://example.com?page=2&per_page=50>; rel=\"prev\"",
			"<http://example.com?page=5&per_page=50>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "300", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Header should default limit to 1 no limit was provided", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 100, 20, 0)

		expect := strings.Join([]string{
			"<http://example.com?page=0&per_page=1>; rel=\"first\"",
			"<http://example.com?page=21&per_page=1>; rel=\"next\"",
			"<http://example.com?page=19&per_page=1>; rel=\"prev\"",
			"<http://example.com?page=99&per_page=1>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "100", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create previous, next, first, but not last if in the middle and no total was provided", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 0, 3, 50)

		expect := strings.Join([]string{
			"<http://example.com?page=0&per_page=50>; rel=\"first\"",
			"<http://example.com?page=4&per_page=50>; rel=\"next\"",
			"<http://example.com?page=2&per_page=50>; rel=\"prev\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "0", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create only first if the limits provided exceeds the number of clients found", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 5, 0, 50)

		expect := "<http://example.com?page=0&per_page=5>; rel=\"first\""

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "5", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create only first if the limits provided equals the number of clients found", func(t *testing.T) {
		r := httptest.NewRecorder()
		PaginationHeader(r, u, 50, 0, 50)

		expect := "<http://example.com?page=0&per_page=50>; rel=\"first\""

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "50", r.Result().Header.Get("X-Total-Count"))
	})
}

func TestParsePagination(t *testing.T) {
	for _, tc := range []struct {
		d                    string
		url                  string
		expectedItemsPerPage int
		expectedPage         int
	}{
		{"normal", "http://localhost/foo?per_page=10&page=10", 10, 10},
		{"defaults", "http://localhost/foo", 250, 0},
		{"limits", "http://localhost/foo?per_page=2000", 1000, 0},
		{"negatives", "http://localhost/foo?per_page=-1&page=-1", 1, 0},
		{"invalid_params", "http://localhost/foo?per_page=a&page=b", 250, 0},
	} {
		t.Run(fmt.Sprintf("case=%s", tc.d), func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			page, perPage := new(PagePaginator).ParsePagination(&http.Request{URL: u})
			assert.EqualValues(t, perPage, tc.expectedItemsPerPage, "per_page")
			assert.EqualValues(t, page, tc.expectedPage, "page")
		})
	}
}
