// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package safecast

import "math"

// Clamp if needed.
func Uint64ToInt64(in uint64) int64 {
	if in > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(in)
}
