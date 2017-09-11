package pkg

import "strings"

// SplitNonEmpty is a special case of strings.Split
// which returns an empty slice if string is empty
func SplitNonEmpty(s, sep string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, sep)
}
