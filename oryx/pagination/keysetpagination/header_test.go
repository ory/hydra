// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/peterhellberg/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	p := &Paginator{
		defaultToken: StringPageToken("default"),
		token:        StringPageToken("next"),
		size:         2,
	}

	u, err := url.Parse("http://ory.sh/")
	require.NoError(t, err)

	r := httptest.NewRecorder()

	Header(r, u, p)

	assert.Len(t, r.Result().Header.Values("link"), 1, "make sure we send one header with multiple comma-separated values rather than multiple headers")

	links := link.ParseResponse(r.Result())
	assert.Contains(t, links, "first")
	assert.Contains(t, links["first"].URI, "page_token=default")

	assert.Contains(t, links, "next")
	assert.Contains(t, links["next"].URI, "page_token=next")

	p.isLast = true
	r = httptest.NewRecorder()
	Header(r, u, p)
	links = link.ParseResponse(r.Result())

	assert.Contains(t, links, "first")
	assert.Contains(t, links["first"].URI, "page_token=default")

	assert.NotContains(t, links, "next")
}
