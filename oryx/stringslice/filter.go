// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

import (
	"slices"
	"strings"
	"unicode"
)

// Filter applies the provided filter function and removes all items from the slice for which the filter function returns true.
// Deprecated: use slices.DeleteFunc instead (changes semantics: the original slice is modified)
func Filter(values []string, filter func(string) bool) []string {
	return slices.DeleteFunc(slices.Clone(values), filter)
}

// TrimEmptyFilter applies the strings.TrimFunc function and removes all empty strings
// Deprecated: use slices.DeleteFunc instead (changes semantics: the original slice is modified)
func TrimEmptyFilter(values []string, trim func(rune) bool) (ret []string) {
	return Filter(values, func(value string) bool {
		return strings.TrimFunc(value, trim) == ""
	})
}

// TrimSpaceEmptyFilter applies the strings.TrimSpace function and removes all empty strings
// Deprecated: use slices.DeleteFunc with strings.TrimSpace instead (changes semantics: the original slice is modified)
func TrimSpaceEmptyFilter(values []string) []string {
	return TrimEmptyFilter(values, unicode.IsSpace)
}
