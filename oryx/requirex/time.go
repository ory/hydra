// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package requirex

import (
	"time"

	"github.com/stretchr/testify/require"
)

// EqualDuration fails if expected and actual are more distant than precision
// Note: The previous implementation incorrectly passed on durations bigger than time.maxDuration (i.e. with zero-time involved) and incorrectly failed on zero durations.
func EqualDuration(t require.TestingT, expected, actual, precision time.Duration) {
	require.Truef(t, expected <= actual+precision && expected >= actual-precision, "expected %s to be within %s of %s", actual, precision, expected)
}

// EqualTime fails if expected and actual are more distant than precision
// Deprecated: use require.WithinDuration instead
// Note: The previous implementation incorrectly passed on durations bigger than time.maxDuration (i.e. with zero-time involved) and incorrectly failed on zero durations.
func EqualTime(t require.TestingT, expected, actual time.Time, precision time.Duration) {
	require.WithinDuration(t, expected, actual, precision)
}
