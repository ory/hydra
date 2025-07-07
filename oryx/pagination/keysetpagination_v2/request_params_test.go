// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/peterhellberg/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetLinkHeader(t *testing.T) {
	t.Parallel()

	keys := [][32]byte{{1, 2, 3}}
	defaultToken, nextToken := NewPageToken(Column{Name: "id", Value: "default"}), NewPageToken(Column{Name: "id", Value: "next"})
	opts := []Option{WithSize(2), WithDefaultToken(defaultToken), WithToken(nextToken)}

	u, err := url.Parse("https://ory.sh/")
	require.NoError(t, err)

	getParsedToken := func(t *testing.T, uri string) PageToken {
		u, err := url.Parse(uri)
		require.NoError(t, err)
		assert.Equal(t, "https", u.Scheme)
		assert.Equal(t, "ory.sh", u.Host)
		raw := u.Query().Get("page_token")
		token, err := ParsePageToken(keys, raw)
		require.NoError(t, err)
		return token
	}

	t.Run("case=not last page", func(t *testing.T) {
		r := httptest.NewRecorder()
		p := NewPaginator(opts...)

		SetLinkHeader(r, keys, u, p)

		assert.Len(t, r.Result().Header.Values("link"), 1, "make sure we send one header with multiple comma-separated values rather than multiple headers")
		links := link.ParseResponse(r.Result())

		require.Contains(t, links, "first")
		assert.Equal(t, defaultToken, getParsedToken(t, links["first"].URI))

		require.Contains(t, links, "next")
		assert.Equal(t, nextToken, getParsedToken(t, links["next"].URI))
	})

	t.Run("case=last page", func(t *testing.T) {
		r := httptest.NewRecorder()
		p := NewPaginator(append(opts, withIsLast(true))...)

		SetLinkHeader(r, keys, u, p)

		assert.Len(t, r.Result().Header.Values("link"), 1, "make sure we send one header with multiple comma-separated values rather than multiple headers")
		links := link.ParseResponse(r.Result())

		require.Contains(t, links, "first")
		assert.Equal(t, defaultToken, getParsedToken(t, links["first"].URI))

		assert.NotContains(t, links, "next")
	})
}

func TestParsePageToken(t *testing.T) {
	t.Parallel()

	keys := [][32]byte{{1, 2, 3}, {4, 5, 6}}

	expectedToken := NewPageToken(Column{Name: "id", Value: "token"}, Column{Name: "name", Order: OrderDescending, Value: "test"})
	encryptedToken := expectedToken.Encrypt(keys)

	t.Run("with valid key", func(t *testing.T) {
		token, err := ParsePageToken(keys, encryptedToken)
		require.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("with rotated key", func(t *testing.T) {
		encryptedToken := expectedToken.Encrypt(keys[1:])
		token, err := ParsePageToken(keys, encryptedToken)
		require.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("with invalid key", func(t *testing.T) {
		token, err := ParsePageToken([][32]byte{{7, 8, 9}}, encryptedToken)
		require.ErrorContains(t, err, "decrypt token")
		assert.Zero(t, token)
	})

	t.Run("uses fallback key", func(t *testing.T) {
		fallbackEncryptedToken := expectedToken.Encrypt(nil)
		for _, noKeys := range [][][32]byte{nil, {}} {
			token, err := ParsePageToken(noKeys, fallbackEncryptedToken)
			require.NoError(t, err)
			assert.Equal(t, expectedToken, token)
		}
	})
}

func TestParse(t *testing.T) {
	t.Parallel()

	keys := [][32]byte{{1, 2, 3}}
	token := NewPageToken(Column{Name: "id", Value: "token"}, Column{Name: "name", Order: OrderDescending, Value: "test"})
	defaultToken := NewPageToken(Column{Name: "id", Value: "default"}, Column{Name: "name", Order: OrderDescending, Value: "default name"})
	encryptedToken := token.Encrypt(keys)

	for _, tc := range []struct {
		name          string
		q             url.Values
		expectedSize  int
		expectedToken PageToken
	}{
		{
			name:          "no query parameters",
			q:             url.Values{},
			expectedSize:  DefaultSize,
			expectedToken: defaultToken,
		},
		{
			name:          "with page token",
			q:             url.Values{"page_token": {encryptedToken}},
			expectedSize:  DefaultSize,
			expectedToken: token,
		},
		{
			name:          "with page size",
			q:             url.Values{"page_size": {"123"}},
			expectedSize:  123,
			expectedToken: defaultToken,
		},
		{
			name:          "with page size and page token",
			q:             url.Values{"page_size": {"123"}, "page_token": {encryptedToken}},
			expectedSize:  123,
			expectedToken: token,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := ParseQueryParams(keys, tc.q)
			require.NoError(t, err)
			paginator := NewPaginator(append(opts, WithDefaultToken(defaultToken))...)
			assert.Equal(t, tc.expectedSize, paginator.Size())
			assert.Equal(t, tc.expectedToken, paginator.PageToken())
		})
	}

	t.Run("invalid page size leads to err", func(t *testing.T) {
		_, err := ParseQueryParams(keys, url.Values{"page_size": {"invalid-int"}})
		require.ErrorIs(t, err, strconv.ErrSyntax)
	})

	t.Run("empty tokens and page sizes work as if unset, empty values are skipped", func(t *testing.T) {
		opts, err := ParseQueryParams(keys, url.Values{})
		require.NoError(t, err)
		paginator := NewPaginator(append(opts, WithDefaultToken(defaultToken))...)
		assert.Equal(t, defaultToken, paginator.PageToken())
		assert.Equal(t, DefaultSize, paginator.Size())

		opts, err = ParseQueryParams(keys, url.Values{"page_token": {""}, "page_size": {""}})
		require.NoError(t, err)
		paginator = NewPaginator(append(opts, WithDefaultToken(defaultToken))...)
		assert.Equal(t, defaultToken, paginator.PageToken())
		assert.Equal(t, DefaultSize, paginator.Size())

		opts, err = ParseQueryParams(keys, url.Values{"page_token": {"", encryptedToken, ""}, "page_size": {"", "123", ""}})
		require.NoError(t, err)
		paginator = NewPaginator(append(opts, WithDefaultToken(defaultToken))...)
		assert.Equal(t, token, paginator.PageToken())
		assert.Equal(t, 123, paginator.Size())
	})
}
