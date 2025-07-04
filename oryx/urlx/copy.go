// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package urlx

import "net/url"

// Copy returns a copy of the input url.
func Copy(src *url.URL) *url.URL {
	var out = new(url.URL)
	*out = *src
	return out
}

// CopyWithQuery returns a copy of the input url with the given query parameters
func CopyWithQuery(src *url.URL, query url.Values) *url.URL {
	out := Copy(src)
	q := out.Query()
	for k := range query {
		q.Set(k, query.Get(k))
	}
	out.RawQuery = q.Encode()
	return out
}
