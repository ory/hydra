/*
 * Copyright © 2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */
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
