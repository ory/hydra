// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagepagination

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ory/x/pagination"
)

type PagePaginator struct {
	MaxItems     int
	DefaultItems int
}

func (p *PagePaginator) defaults() {
	if p.MaxItems == 0 {
		p.MaxItems = 1000
	}

	if p.DefaultItems == 0 {
		p.DefaultItems = 250
	}
}

// ParsePagination parses limit and page from *http.Request with given limits and defaults.
func (p *PagePaginator) ParsePagination(r *http.Request) (page, itemsPerPage int) {
	p.defaults()

	if offsetParam := r.URL.Query().Get("page"); offsetParam == "" {
		page = 0
	} else {
		if offset, err := strconv.ParseInt(offsetParam, 10, 0); err != nil {
			page = 0
		} else {
			page = int(offset)
		}
	}

	if limitParam := r.URL.Query().Get("per_page"); limitParam == "" {
		itemsPerPage = p.DefaultItems
	} else {
		if limit, err := strconv.ParseInt(limitParam, 10, 0); err != nil {
			itemsPerPage = p.DefaultItems
		} else {
			itemsPerPage = int(limit)
		}
	}

	if itemsPerPage > p.MaxItems {
		itemsPerPage = p.MaxItems
	}

	if itemsPerPage < 1 {
		itemsPerPage = 1
	}

	if page < 0 {
		page = 0
	}

	return
}

func header(u *url.URL, rel string, limit, offset int64) string {
	q := u.Query()
	q.Set("per_page", fmt.Sprintf("%d", limit))
	q.Set("page", fmt.Sprintf("%d", offset/limit))
	u.RawQuery = q.Encode()
	return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), rel)
}

func PaginationHeader(w http.ResponseWriter, u *url.URL, total int64, page, itemsPerPage int) {
	pagination.HeaderWithFormatter(w, u, total, page, itemsPerPage, header)
}
