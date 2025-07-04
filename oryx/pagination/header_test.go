// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	u, err := url.Parse("http://example.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Create previous and first but not next or last if at the end", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 120, 50, 100)

		expect := strings.Join([]string{
			"<http://example.com?limit=50&offset=0>; rel=\"first\"",
			"<http://example.com?limit=50&offset=50>; rel=\"prev\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
		assert.EqualValues(t, "120", r.Result().Header.Get("X-Total-Count"))
	})

	t.Run("Create next and last, but not previous or first if at the beginning", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 120, 50, 0)

		expect := strings.Join([]string{
			"<http://example.com?limit=50&offset=50>; rel=\"next\"",
			"<http://example.com?limit=50&offset=100>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
	})

	t.Run("Create next and last, but not previous or first if on the first page", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 120, 50, 10)

		expect := strings.Join([]string{
			"<http://example.com?limit=50&offset=50>; rel=\"next\"",
			"<http://example.com?limit=50&offset=100>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
	})

	t.Run("Create previous, next, first, and last if in the middle", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 300, 50, 150)

		expect := strings.Join([]string{
			"<http://example.com?limit=50&offset=0>; rel=\"first\"",
			"<http://example.com?limit=50&offset=200>; rel=\"next\"",
			"<http://example.com?limit=50&offset=100>; rel=\"prev\"",
			"<http://example.com?limit=50&offset=250>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
	})

	t.Run("Header should default limit to 1 no limit was provided", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 100, 0, 20)

		expect := strings.Join([]string{
			"<http://example.com?limit=1&offset=0>; rel=\"first\"",
			"<http://example.com?limit=1&offset=21>; rel=\"next\"",
			"<http://example.com?limit=1&offset=19>; rel=\"prev\"",
			"<http://example.com?limit=1&offset=99>; rel=\"last\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
	})

	t.Run("Create previous, next, first, but not last if in the middle and no total was provided", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 0, 50, 150)

		expect := strings.Join([]string{
			"<http://example.com?limit=50&offset=0>; rel=\"first\"",
			"<http://example.com?limit=50&offset=200>; rel=\"next\"",
			"<http://example.com?limit=50&offset=100>; rel=\"prev\"",
		}, ",")

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
	})

	t.Run("Create only first if the limits provided exceeds the number of clients found", func(t *testing.T) {
		r := httptest.NewRecorder()
		Header(r, u, 5, 50, 0)

		expect := "<http://example.com?limit=5&offset=0>; rel=\"first\""

		assert.EqualValues(t, expect, r.Result().Header.Get("Link"))
	})
}
