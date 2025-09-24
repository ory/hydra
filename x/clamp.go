// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

// Clamp clamps val to be within the range [min, max] for any integer type.
func Clamp[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](val, min, max T) T {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
