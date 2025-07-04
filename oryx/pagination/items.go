// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

// MaxItemsPerPage is used to prevent DoS attacks against large lists by limiting the items per page to 500.
func MaxItemsPerPage(max, is int) int {
	if is > max {
		return max
	}
	return is
}
