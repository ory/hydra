// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package tokenpagination

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/ory/x/pagination"

	"github.com/ory/herodot"
)

func Encode(offset int64) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"offset":"%d","v":2}`, offset)))
}

func decode(s string) (int, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, errors.WithStack(herodot.ErrBadRequest.WithWrap(err).WithReasonf("Unable to parse pagination token: %s", err))
	}

	return int(gjson.Get(string(b), "offset").Int()), nil
}

type TokenPaginator struct {
	MaxItems     int
	DefaultItems int
}

func (p *TokenPaginator) defaults() {
	if p.MaxItems == 0 {
		p.MaxItems = 1000
	}

	if p.DefaultItems == 0 {
		p.DefaultItems = 250
	}
}

// ParsePagination parses limit and page from *http.Request with given limits and defaults.
func (p *TokenPaginator) ParsePagination(r *http.Request) (page, itemsPerPage int) {
	p.defaults()

	var offset int
	if offsetParam := r.URL.Query().Get("page_token"); len(offsetParam) > 0 {
		offset, _ = decode(offsetParam)
	}

	if gotLimit, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 0); err == nil {
		itemsPerPage = int(gotLimit)
	} else {
		itemsPerPage = p.DefaultItems
	}

	if itemsPerPage > p.MaxItems {
		itemsPerPage = p.MaxItems
	}

	if itemsPerPage < 1 {
		itemsPerPage = 1
	}

	if offset > 0 {
		page = offset / itemsPerPage
	}

	if page < 0 {
		page = 0
	}

	return
}

func header(u *url.URL, rel string, itemsPerPage, offset int64) string {
	q := u.Query()
	q.Set("page_size", fmt.Sprintf("%d", itemsPerPage))
	q.Set("page_token", Encode(offset))
	u.RawQuery = q.Encode()
	return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), rel)
}

func PaginationHeader(w http.ResponseWriter, u *url.URL, total int64, page, itemsPerPage int) {
	pagination.HeaderWithFormatter(w, u, total, page, itemsPerPage, header)
}
