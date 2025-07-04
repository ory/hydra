// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"net/http"
	"net/url"

	"github.com/peterhellberg/link"
)

// PaginationResult represents a parsed result of the link HTTP header.
type PaginationResult struct {
	// NextToken is the next page token. If it's empty, there is no next page.
	NextToken string

	// FirstToken is the first page token.
	FirstToken string
}

// ParseHeader parses the response header's Link.
func ParseHeader(r *http.Response) *PaginationResult {
	links := link.ParseResponse(r)
	return &PaginationResult{
		NextToken:  findRel(links, "next"),
		FirstToken: findRel(links, "first"),
	}
}

func findRel(links link.Group, rel string) string {
	for idx, l := range links {
		if idx == rel {
			parsed, err := url.Parse(l.URI)
			if err != nil {
				continue
			}

			return parsed.Query().Get("page_token")
		}
	}

	return ""
}
