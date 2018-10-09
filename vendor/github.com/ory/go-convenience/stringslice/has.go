package stringslice

import "strings"

// Has returns true if the needle is in the haystack (case-sensitive)
func Has(haystack []string, needle string) bool {
	for _, current := range haystack {
		if current == needle {
			return true
		}
	}
	return false
}

// HasI returns true if the needle is in the haystack (case-insensitive)
func HasI(haystack []string, needle string) bool {
	for _, current := range haystack {
		if strings.ToLower(current) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}
