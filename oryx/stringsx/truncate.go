// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import "unicode/utf8"

// TruncateByteLen returns string truncated at the end with the length specified
func TruncateByteLen(s string, length int) string {
	if length <= 0 || len(s) <= length {
		return s
	}

	res := s[:length]

	// in case we cut in the middle of an utf8 rune, we have to remove the last byte as well until it fits
	for !utf8.ValidString(res) {
		res = res[:len(res)-1]
	}
	return res
}
