// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"net/http"
	"net/url"

	"github.com/ory/x/pagination/tokenpagination"
)

// swagger:model paginationHeaders
type PaginationHeaders struct {
	// The link header contains pagination links.
	//
	// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// in: header
	Link string `json:"link"`

	// The total number of clients.
	//
	// in: header
	XTotalCount string `json:"x-total-count"`
}

// swagger:model pagination
type PaginationParams struct {
	// Items per page
	//
	// This is the number of items per page to return.
	// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// required: false
	// in: query
	// default: 250
	// min: 1
	// max: 1000
	PageSize int `json:"page_size"`

	// Next Page Token
	//
	// The next page token.
	// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// required: false
	// in: query
	// default: 1
	// min: 1
	PageToken string `json:"page_token"`
}

const paginationMaxItems = 1000
const paginationDefaultItems = 250

var paginator = &tokenpagination.TokenPaginator{
	MaxItems:     paginationMaxItems,
	DefaultItems: paginationDefaultItems,
}

// ParsePagination parses limit and page from *http.Request with given limits and defaults.
func ParsePagination(r *http.Request) (page, itemsPerPage int) {
	return paginator.ParsePagination(r)
}

func PaginationHeader(w http.ResponseWriter, u *url.URL, total int64, page, itemsPerPage int) {
	tokenpagination.PaginationHeader(w, u, total, page, itemsPerPage)
}
