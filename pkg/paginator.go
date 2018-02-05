package pkg

import (
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request, defaultLimit, defaultOffset, maxLimit int64) (int64, int64) {
	var offset, limit int64
	var err error

	if offsetParam := r.URL.Query().Get("offset"); offsetParam == "" {
		offset = defaultOffset
	} else {
		offset, err = strconv.ParseInt(offsetParam, 10, 64)
		if err != nil {
			offset = defaultOffset
		}
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.ParseInt(limitParam, 10, 64)
		if err != nil {
			limit = defaultLimit
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
