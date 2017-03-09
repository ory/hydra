package fosite

import (
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

func removeEmpty(args []string) (ret []string) {
	for _, v := range args {
		v = strings.TrimSpace(v)
		if v != "" {
			ret = append(ret, v)
		}
	}
	return
}
