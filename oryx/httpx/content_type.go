// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"mime"
	"net/http"
	"slices"
	"strings"
)

// HasContentType determines whether the request `content-type` includes a
// server-acceptable mime-type
//
// Failure should yield an HTTP 415 (`http.StatusUnsupportedMediaType`)
func HasContentType(r *http.Request, mimetypes ...string) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return slices.Contains(mimetypes, "application/octet-stream")
	}

	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err != nil {
		return false
	}
	return slices.Contains(mimetypes, mediaType)
}
