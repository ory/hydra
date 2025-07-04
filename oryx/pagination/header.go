// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func header(u *url.URL, rel string, limit, offset int64) string {
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	u.RawQuery = q.Encode()
	return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), rel)
}

type formatter func(location *url.URL, rel string, itemsPerPage int64, offset int64) string

// HeaderWithFormatter adds an HTTP header for pagination which uses a custom formatter for generating the URL links.
func HeaderWithFormatter(w http.ResponseWriter, u *url.URL, total int64, page, itemsPerPage int, f formatter) {
	if itemsPerPage <= 0 {
		itemsPerPage = 1
	}

	itemsPerPage64 := int64(itemsPerPage)
	offset := int64(page) * itemsPerPage64

	// lastOffset will either equal the offset required to contain the remainder,
	// or the limit.
	var lastOffset int64
	if total%itemsPerPage64 == 0 {
		lastOffset = total - itemsPerPage64
	} else {
		lastOffset = (total / itemsPerPage64) * itemsPerPage64
	}

	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))

	// Check for last page
	if offset >= lastOffset {
		if total == 0 {
			w.Header().Set("Link", strings.Join([]string{
				f(u, "first", itemsPerPage64, 0),
				f(u, "next", itemsPerPage64, ((offset/itemsPerPage64)+1)*itemsPerPage64),
				f(u, "prev", itemsPerPage64, ((offset/itemsPerPage64)-1)*itemsPerPage64),
			}, ","))
			return
		}

		if total <= itemsPerPage64 {
			w.Header().Set("link", f(u, "first", total, 0))
			return
		}

		w.Header().Set("Link", strings.Join([]string{
			f(u, "first", itemsPerPage64, 0),
			f(u, "prev", itemsPerPage64, lastOffset-itemsPerPage64),
		}, ","))
		return
	}

	if offset < itemsPerPage64 {
		w.Header().Set("Link", strings.Join([]string{
			f(u, "next", itemsPerPage64, itemsPerPage64),
			f(u, "last", itemsPerPage64, lastOffset),
		}, ","))
		return
	}

	w.Header().Set("Link", strings.Join([]string{
		f(u, "first", itemsPerPage64, 0),
		f(u, "next", itemsPerPage64, ((offset/itemsPerPage64)+1)*itemsPerPage64),
		f(u, "prev", itemsPerPage64, ((offset/itemsPerPage64)-1)*itemsPerPage64),
		f(u, "last", itemsPerPage64, lastOffset),
	}, ","))
}

// Header adds an http header for pagination using a responsewriter where backwards compatibility is required.
// The header will contain links any combination of the first, last, next, or previous (prev) pages in a paginated list (given a limit and an offset, and optionally a total).
// If total is not set, then no "last" page will be calculated.
// If no limit is provided, then it will default to 1.
func Header(w http.ResponseWriter, u *url.URL, total int, limit, offset int) {
	var page int
	if limit == 0 {
		limit = 1
	}

	page = int(math.Floor(float64(offset) / float64(limit)))
	HeaderWithFormatter(w, u, int64(total), page, limit, header)
}
