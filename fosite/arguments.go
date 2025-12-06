// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "strings"

type Arguments []string

// Matches performs an case-insensitive, out-of-order check that the items
// provided exist and equal all of the args in arguments.
// Note:
//   - Providing a list that includes duplicate string-case items will return not
//     matched.
func (r Arguments) Matches(items ...string) bool {
	if len(r) != len(items) {
		return false
	}

	found := make(map[string]bool)
	for _, item := range items {
		if !StringInSlice(item, r) {
			return false
		}
		found[item] = true
	}

	return len(found) == len(r)
}

// Has checks, in a case-insensitive manner, that all of the items
// provided exists in arguments.
func (r Arguments) Has(items ...string) bool {
	for _, item := range items {
		if !StringInSlice(item, r) {
			return false
		}
	}

	return true
}

// HasOneOf checks, in a case-insensitive manner, that one of the items
// provided exists in arguments.
func (r Arguments) HasOneOf(items ...string) bool {
	for _, item := range items {
		if StringInSlice(item, r) {
			return true
		}
	}

	return false
}

// Deprecated: Use ExactOne, Matches or MatchesExact
func (r Arguments) Exact(name string) bool {
	return name == strings.Join(r, " ")
}

// ExactOne checks, by string case, that a single argument equals the provided
// string.
func (r Arguments) ExactOne(name string) bool {
	return len(r) == 1 && r[0] == name
}

// MatchesExact checks, by order and string case, that the items provided equal
// those in arguments.
func (r Arguments) MatchesExact(items ...string) bool {
	if len(r) != len(items) {
		return false
	}

	for i, item := range items {
		if item != r[i] {
			return false
		}
	}

	return true
}
