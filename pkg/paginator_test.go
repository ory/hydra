package pkg

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	t.Run("case=normal", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo?limit=10&offset=10")
		limit, offset := ParsePagination(&http.Request{URL: u}, 0, 0, 10)
		assert.EqualValues(t, limit, 10)
		assert.EqualValues(t, offset, 10)
	})

	t.Run("case=defaults", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo")
		limit, offset := ParsePagination(&http.Request{URL: u}, 5, 5, 10)
		assert.EqualValues(t, limit, 5)
		assert.EqualValues(t, offset, 5)
	})

	t.Run("case=defaults_and_limits", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo")
		limit, offset := ParsePagination(&http.Request{URL: u}, 5, 5, 2)
		assert.EqualValues(t, limit, 2)
		assert.EqualValues(t, offset, 5)
	})

	t.Run("case=limits", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo?limit=10&offset=10")
		limit, offset := ParsePagination(&http.Request{URL: u}, 0, 0, 5)
		assert.EqualValues(t, limit, 5)
		assert.EqualValues(t, offset, 10)
	})

	t.Run("case=negatives", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo?limit=-1&offset=-1")
		limit, offset := ParsePagination(&http.Request{URL: u}, 0, 0, 5)
		assert.EqualValues(t, limit, 0)
		assert.EqualValues(t, offset, 0)
	})

	t.Run("case=default_negatives", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo")
		limit, offset := ParsePagination(&http.Request{URL: u}, -1, -1, 5)
		assert.EqualValues(t, limit, 0)
		assert.EqualValues(t, offset, 0)
	})

	t.Run("case=invalid_defaults", func(t *testing.T) {
		u, _ := url.Parse("http://localhost/foo?offset=a&limit=b")
		limit, offset := ParsePagination(&http.Request{URL: u}, 10, 10, 15)
		assert.EqualValues(t, limit, 10)
		assert.EqualValues(t, offset, 10)
	})
}
