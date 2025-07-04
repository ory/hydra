// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package urlx

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	assert.EqualValues(t, "http://foo/bar/baz/bar", MustJoin("http://foo", "bar/", "/baz", "bar"))
}

func TestAppendPaths(t *testing.T) {
	u, err := url.Parse("http://localhost/home/")
	require.NoError(t, err)
	assert.Equal(t, "http://localhost/home/", AppendPaths(u).String())

	for k, tc := range []struct {
		give   []string
		expect string
	}{
		{
			give:   []string{"http://localhost/", "/home"},
			expect: "http://localhost/home",
		},
		{
			give:   []string{"http://localhost", "/home"},
			expect: "http://localhost/home",
		},
		{
			give:   []string{"https://localhost/", "/home"},
			expect: "https://localhost/home",
		},
		{
			give:   []string{"http://localhost/", "/home", "home/", "/home/"},
			expect: "http://localhost/home/home/home/",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			u, err := url.Parse(tc.give[0])
			require.NoError(t, err)
			assert.Equal(t, tc.expect, AppendPaths(u, tc.give[1:]...).String())
		})
	}
}

func TestAppendQuery(t *testing.T) {
	u, err := url.Parse("http://localhost/home?foo=bar&baz=bar")
	require.NoError(t, err)

	assert.Equal(t, "http://localhost/home?baz=bar&foo=bar", SetQuery(u, url.Values{}).String())
	assert.Equal(t, "http://localhost/home?bar=baz&baz=bar&foo=bar", SetQuery(u, url.Values{"bar": {"baz"}}).String())
	assert.Equal(t, "http://localhost/home?bar=baz&baz=bar&foo=bar", SetQuery(u, url.Values{"bar": {"baz", "baz"}}).String())
	assert.Equal(t, "http://localhost/home?baz=foo&foo=bar", SetQuery(u, url.Values{"baz": {"foo"}}).String())
}
