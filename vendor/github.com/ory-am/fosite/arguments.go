package fosite

import "strings"

type Arguments []string

func (r Arguments) Matches(items ...string) bool {
	found := make(map[string]bool)
	for _, item := range items {
		if !StringInSlice(item, r) {
			return false
		}
		found[item] = true
	}

	return len(found) == len(r) && len(r) == len(items)
}

func (r Arguments) Has(items ...string) bool {
	for _, item := range items {
		if !StringInSlice(item, r) {
			return false
		}
	}

	return true
}

func (r Arguments) Exact(name string) bool {
	return name == strings.Join(r, " ")
}
