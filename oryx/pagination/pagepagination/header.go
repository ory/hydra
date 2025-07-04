// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagepagination

// Pagination Request Parameters
//
// The `Link` HTTP header contains multiple links (`first`, `next`, `last`, `previous`) formatted as:
// `<https://{project-slug}.projects.oryapis.com/admin/clients?page_size={limit}&page_token={offset}>; rel="{page}"`
//
// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
//
// swagger:model pagePaginationRequestParameters
type RequestParameters struct {
	// Legacy Items per Page
	//
	// A DEPRECATED alias for `page_size`. Please transition to using `page_size` going forward.
	//
	// required: false
	// in: query
	// default: 250
	// min: 1
	// max: 1000
	PerPage int `json:"per_page"`

	// Legacy Pagination Page
	//
	// A DEPRECATED alias for `page_token`. Please transition to using `page_token` going forward.
	//
	// required: false
	// in: query
	// default: 1
	// min: 1
	Page int `json:"page"`

	// Items per Page
	//
	// This is the number of items per page to return. For details on pagination please head over to the
	// [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// required: false
	// in: query
	// default: 250
	// min: 1
	// max: 500
	PageSize int `json:"page_size"`

	// Next Page Token
	//
	// The next page token. For details on pagination please head over to the
	// [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// required: false
	// in: query
	// default: 1
	// min: 1
	PageToken string `json:"page_token"`
}

// swagger:model pagePaginationResponseHeaders
type ResponseHeaderAnnotation struct {
	// The Link HTTP Header
	//
	// The `Link` header contains a comma-delimited list of links to the following pages:
	//
	// - first: The first page of results.
	// - next: The next page of results.
	// - prev: The previous page of results.
	// - last: The last page of results.
	//
	// Pages are omitted if they do not exist. For example, if there is no next page, the `next` link is omitted.
	//
	// This header will include the `per_page` and `page` parameters for legacy reasons, but these parameters will eventually be removed.
	//
	//	Example: Link: </clients?page_size=5&page_token=0>; rel="first",</clients?page_size=5&page_token=15>; rel="next",</clients?page_size=5&page_token=5>; rel="prev",</clients?page_size=5&page_token=20>; rel="last"
	Link string `json:"link"`

	// The X-Total-Count HTTP Header
	//
	// The `X-Total-Count` header contains the total number of items in the collection.
	//
	// Example: 123
	TotalCount int `json:"x-total-count"`
}
