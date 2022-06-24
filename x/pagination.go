package x

import (
	"net/http"
	"net/url"

	"github.com/ory/x/pagination/tokenpagination"
)

// swagger:model pagination
type PaginationParams struct {
	// Items per PageToken
	//
	// This is the number of items per page.
	//
	// required: false
	// in: query
	// default: 250
	// min: 1
	// max: 1000
	PageSize int `json:"page_size"`

	// Pagination PageToken Token
	//
	// The page token.
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
