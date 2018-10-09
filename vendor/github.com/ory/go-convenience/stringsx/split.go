package stringsx

import "strings"

// Splitx is a special case of strings.Split
// which returns an empty slice if the string is empty
func Splitx(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, sep)
}
