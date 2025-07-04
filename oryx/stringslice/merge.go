// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

import "slices"

// Merge merges several string slices into one.
// Deprecated: use slices.Concat instead
func Merge(parts ...[]string) []string {
	return slices.Concat(parts...)
}
