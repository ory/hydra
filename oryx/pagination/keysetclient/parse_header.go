// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package keysetclient parses keyset pagination metadata from HTTP responses.
// It is the client-side counterpart to keysetpagination and deliberately
// carries no SQL dependencies, so API clients such as the Ory CLI do not link
// the database layer.
package keysetclient

import (
	"net/http"
	"net/url"

	"github.com/peterhellberg/link"
)

// Result represents a parsed result of the link HTTP header.
type Result struct {
	// NextToken is the next page token. If it's empty, there is no next page.
	NextToken string

	// FirstToken is the first page token.
	FirstToken string
}

// ParseHeader parses the response header's Link.
func ParseHeader(r *http.Response) *Result {
	links := link.ParseResponse(r)
	return &Result{
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
