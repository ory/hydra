// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openapix

// swagger:model tokenPaginationHeaders
type TokenPaginationHeaders struct {
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

// swagger:model tokenPagination
type TokenPaginationParams struct {
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
