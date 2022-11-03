// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

// ToPointer returns the pointer to the value.
func ToPointer[T any](val T) *T {
	return &val
}

// FromPointer returns the dereferenced value or if the pointer is nil the zero value.
func FromPointer[T any, TT *T](val *T) (zero T) {
	if val == nil {
		return zero
	}
	return *val
}
