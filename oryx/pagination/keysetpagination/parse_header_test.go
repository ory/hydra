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
	u, err := url.Parse("https://www.ory.sh/")
	require.NoError(t, err)

	t.Run("has next page", func(t *testing.T) {
		p := &Paginator{
			defaultToken: StringPageToken("default"),
			token:        StringPageToken("next"),
			size:         2,
		}

		r := httptest.NewRecorder()
		Header(r, u, p)

		result := ParseHeader(&http.Response{Header: r.Header()})
		assert.Equal(t, "next", result.NextToken, r.Header())
		assert.Equal(t, "default", result.FirstToken, r.Header())
	})

	t.Run("is last page", func(t *testing.T) {
		p := &Paginator{
			defaultToken: StringPageToken("default"),
			size:         1,
			isLast:       true,
		}

		r := httptest.NewRecorder()
		Header(r, u, p)

		result := ParseHeader(&http.Response{Header: r.Header()})
		assert.Equal(t, "", result.NextToken, r.Header())
		assert.Equal(t, "default", result.FirstToken, r.Header())
	})
}
