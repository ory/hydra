// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hmac

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomBytes(t *testing.T) {
	bytes, err := RandomBytes(128)
	assert.NoError(t, err)
	assert.Len(t, bytes, 128)
}

func TestPseudoRandomness(t *testing.T) {
	runs := 65536
	results := map[string]struct{}{}
	for i := 0; i < runs; i++ {
		bytes, err := RandomBytes(128)
		require.NoError(t, err)
		results[string(bytes)] = struct{}{}
	}
	assert.Len(t, results, runs)
}
