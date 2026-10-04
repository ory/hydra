// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cachex

import (
	"net/http"
	"strings"
)

type CacheControl map[string]string

// From Go generated SDK code.
func ParseCacheControl(headers http.Header) CacheControl {
	cc := CacheControl{}
	ccHeader := headers.Get("Cache-Control")
	for _, part := range strings.Split(ccHeader, ",") {
		part = strings.Trim(part, " ")
		if part == "" {
			continue
		}
		if strings.ContainsRune(part, '=') {
			keyval := strings.Split(part, "=")
			cc[strings.Trim(keyval[0], " ")] = strings.Trim(keyval[1], ",")
		} else {
			cc[part] = ""
		}
	}
	return cc
}
