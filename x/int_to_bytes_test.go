// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_toBytes_fromBytes(t *testing.T) {
	for _, tc := range []struct {
		name string
		i    int64
	}{
		{
			name: "zero",
			i:    0,
		},
		{
			name: "positive",
			i:    1234567890,
		},
		{
			name: "negative",
			i:    -1234567890,
		},
		{
			name: "now",
			i:    time.Now().Unix(),
		},
		{
			name: "max",
			i:    math.MaxInt64,
		},
		{
			name: "min",
			i:    math.MinInt64,
		},
	} {
		t.Run("case="+tc.name, func(t *testing.T) {
			bytes := IntToBytes(tc.i)
			i, err := BytesToInt(bytes)
			require.NoError(t, err)
			assert.Equal(t, tc.i, i)
		})
	}
}
