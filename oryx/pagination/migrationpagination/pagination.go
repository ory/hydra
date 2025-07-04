// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package migrationpagination

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ory/x/pagination"
	"github.com/ory/x/pagination/pagepagination"
	"github.com/ory/x/pagination/tokenpagination"
)

type Paginator struct {
	p *pagepagination.PagePaginator
	t *tokenpagination.TokenPaginator
}

func NewPaginator(p *pagepagination.PagePaginator, t *tokenpagination.TokenPaginator) *Paginator {
	return &Paginator{p: p, t: t}
}

func NewDefaultPaginator() *Paginator {
	return &Paginator{p: new(pagepagination.PagePaginator), t: new(tokenpagination.TokenPaginator)}
}

func (p *Paginator) ParsePagination(r *http.Request) (page, itemsPerPage int) {
	if r.URL.Query().Has("page_token") || r.URL.Query().Has("page_size") {
		return p.t.ParsePagination(r)
	}
	return p.p.ParsePagination(r)
}

func header(u *url.URL, rel string, itemsPerPage, offset int64) string {
	q := u.Query()
	q.Set("page_size", fmt.Sprintf("%d", itemsPerPage))
	q.Set("page_token", tokenpagination.Encode(offset))
	q.Set("per_page", fmt.Sprintf("%d", itemsPerPage))
	q.Set("page", fmt.Sprintf("%d", offset/itemsPerPage))
	u.RawQuery = q.Encode()
	return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), rel)
}

func PaginationHeader(w http.ResponseWriter, u *url.URL, total int64, page, itemsPerPage int) {
	pagination.HeaderWithFormatter(w, u, total, page, itemsPerPage, header)
}
