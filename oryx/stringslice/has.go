// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

import (
	"slices"
	"strings"
)

// Has returns true if the needle is in the haystack (case-sensitive)
// Deprecated: use slices.Contains instead
func Has(haystack []string, needle string) bool {
	return slices.Contains(haystack, needle)
}

// HasI returns true if the needle is in the haystack (case-insensitive)
func HasI(haystack []string, needle string) bool {
	return slices.ContainsFunc(haystack, func(value string) bool {
		return strings.EqualFold(value, needle)
	})
}
