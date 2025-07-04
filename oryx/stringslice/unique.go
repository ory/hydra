// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

// Unique returns the given string slice with unique values, preserving order.
// Consider using slices.Compact with slices.Sort instead when you don't care about order.
func Unique(i []string) []string {
	u := make([]string, 0, len(i))
	m := make(map[string]struct{}, len(i))

	for _, val := range i {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}

	return u
}
