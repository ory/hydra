// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"fmt"
	"strings"
)

// StringInSlice returns true if needle exists in haystack
func StringInSlice(needle string, haystack []string) bool {
	for _, b := range haystack {
		if strings.ToLower(b) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

func RemoveEmpty(args []string) (ret []string) {
	for _, v := range args {
		v = strings.TrimSpace(v)
		if v != "" {
			ret = append(ret, v)
		}
	}
	return
}

// EscapeJSONString does a poor man's JSON encoding. Useful when we do not want to use full JSON encoding
// because we just had an error doing the JSON encoding. The characters that MUST be escaped: quotation mark,
// reverse solidus, and the control characters (U+0000 through U+001F).
// See: https://tools.ietf.org/html/std90#section-7
func EscapeJSONString(str string) string {
	// Escape reverse solidus.
	str = strings.ReplaceAll(str, `\`, `\\`)
	// Escape control characters.
	for r := rune(0); r < ' '; r++ {
		str = strings.ReplaceAll(str, string(r), fmt.Sprintf(`\u%04x`, r))
	}
	// Escape quotation mark.
	str = strings.ReplaceAll(str, `"`, `\"`)
	return str
}
