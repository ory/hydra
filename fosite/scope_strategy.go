// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "strings"

// ScopeStrategy is a strategy for matching scopes.
type ScopeStrategy func(haystack []string, needle string) bool

func HierarchicScopeStrategy(haystack []string, needle string) bool {
	for _, this := range haystack {
		// foo == foo -> true
		if this == needle {
			return true
		}

		// picture.read > picture -> false (scope picture includes read, write, ...)
		if len(this) > len(needle) {
			continue
		}

		needles := strings.Split(needle, ".")
		haystack := strings.Split(this, ".")
		haystackLen := len(haystack) - 1
		for k, needle := range needles {
			if haystackLen < k {
				return true
			}

			current := haystack[k]
			if current != needle {
				break
			}
		}
	}

	return false
}

func ExactScopeStrategy(haystack []string, needle string) bool {
	for _, this := range haystack {
		if needle == this {
			return true
		}
	}

	return false
}

func WildcardScopeStrategy(matchers []string, needle string) bool {
	needleParts := strings.Split(needle, ".")
	for _, matcher := range matchers {
		matcherParts := strings.Split(matcher, ".")
		if len(matcherParts) > len(needleParts) {
			continue
		}

		var noteq bool
		for k, c := range matcherParts {
			// this is the last item and the lengths are different
			if k == len(matcherParts)-1 && len(matcherParts) != len(needleParts) {
				if c != "*" {
					noteq = true
					break
				}
			}

			if c == "*" && len(needleParts[k]) > 0 {
				// pass because this satisfies the requirements
				continue
			} else if c != needleParts[k] {
				noteq = true
				break
			}
		}

		if !noteq {
			return true
		}
	}

	return false
}
