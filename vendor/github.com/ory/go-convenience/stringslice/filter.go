package stringslice

import (
	"strings"
	"unicode"
)

// Filter applies the provided filter function and removes all items from the slice for which the filter function returns true.
// This function uses append and might cause
func Filter(values []string, filter func(string) bool) (ret []string) {
	for _, value := range values {
		if !filter(value) {
			ret = append(ret, value)
		}
	}

	if ret == nil {
		return []string{}
	}

	return
}

// TrimEmptyFilter applies the strings.TrimFunc function and removes all empty strings
func TrimEmptyFilter(values []string, trim func(rune) bool) (ret []string) {
	return Filter(values, func(value string) bool {
		return strings.TrimFunc(value, trim) == ""
	})
}

// TrimSpaceEmptyFilter applies the strings.TrimSpace function and removes all empty strings
func TrimSpaceEmptyFilter(values []string) []string {
	return TrimEmptyFilter(values, unicode.IsSpace)
}
