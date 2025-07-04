// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeader(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://www.ory.sh/")
	require.NoError(t, err)
	keys := [][32]byte{{1, 2, 3}}
	defaultToken, nextToken := NewPageToken(Column{Name: "id", Value: "default"}), NewPageToken(Column{Name: "id", Value: "next"})

	t.Run("has next page", func(t *testing.T) {
		p := NewPaginator(WithSize(2), WithDefaultToken(defaultToken), WithToken(nextToken))
		r := httptest.NewRecorder()
		SetLinkHeader(r, keys, u, p)

		first, next, isLast := ParseHeader(&http.Response{Header: r.Header()})
		require.NotEqual(t, first, next, r.Header())
		assert.False(t, isLast)

		parsedFirst, err := ParsePageToken(keys, first)
		require.NoErrorf(t, err, "raw token %q", first)
		assert.Equal(t, defaultToken, parsedFirst, r.Header())

		parsedNext, err := ParsePageToken(keys, next)
		require.NoErrorf(t, err, "raw token %q", next)
		assert.Equal(t, nextToken, parsedNext, r.Header())
	})

	t.Run("is last page", func(t *testing.T) {
		p := NewPaginator(WithSize(2), WithDefaultToken(defaultToken), WithToken(nextToken), withIsLast(true))
		r := httptest.NewRecorder()
		SetLinkHeader(r, keys, u, p)

		first, next, isLast := ParseHeader(&http.Response{Header: r.Header()})
		assert.Empty(t, next, r.Header())
		assert.True(t, isLast)

		parsedFirst, err := ParsePageToken(keys, first)
		require.NoErrorf(t, err, "raw token %q", first)
		assert.Equal(t, defaultToken, parsedFirst, r.Header())
	})
}
