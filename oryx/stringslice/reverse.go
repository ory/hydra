// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

import "slices"

// Reverse reverses the order of a string slice
// Deprecated: use slices.Reverse instead (changes semantics)
func Reverse(s []string) []string {
	c := slices.Clone(s)
	slices.Reverse(c)
	return c
}
