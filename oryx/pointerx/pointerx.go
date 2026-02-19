// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pointerx

// Deref returns the input values de-referenced value, or zero value if nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
