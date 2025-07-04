// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"net/http"
	"net/url"

	"github.com/peterhellberg/link"
)

// ParseHeader parses the response header's Link and returns the first and next page tokens.
func ParseHeader(r *http.Response) (first, next string, isLast bool) {
	links := link.ParseResponse(r)
	first, _ = findRel(links, "first")
	next, hasNext := findRel(links, "next")
	return first, next, !hasNext
}

func findRel(links link.Group, rel string) (string, bool) {
	for idx, l := range links {
		if idx == rel {
			parsed, err := url.Parse(l.URI)
			if err != nil {
				continue
			}
			q := parsed.Query()

			return q.Get("page_token"), q.Has("page_token")
		}
	}

	return "", false
}
