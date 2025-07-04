// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"net/http"
	"strconv"
)

// Parse parses limit and offset from *http.Request with given limits and defaults.
func Parse(r *http.Request, defaultLimit, defaultOffset, maxLimit int) (int, int) {
	var offset, limit int

	if offsetParam := r.URL.Query().Get("offset"); offsetParam == "" {
		offset = defaultOffset
	} else {
		if offset64, err := strconv.ParseInt(offsetParam, 10, 64); err != nil {
			offset = defaultOffset
		} else {
			offset = int(offset64)
		}
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam == "" {
		limit = defaultLimit
	} else {
		if limit64, err := strconv.ParseInt(limitParam, 10, 64); err != nil {
			limit = defaultLimit
		} else {
			limit = int(limit64)
		}
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
}
