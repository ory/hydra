// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRotatedHashes(t *testing.T) {
	t.Run("returns nil when no rotated secrets", func(t *testing.T) {
		c := &Client{
			RotatedSecrets: []string{},
		}
		hashes := c.GetRotatedHashes()
		assert.Nil(t, hashes)
		c = &Client{
			RotatedSecrets: nil,
		}
		hashes = c.GetRotatedHashes()
		assert.Nil(t, hashes)
	})

	t.Run("returns hashes when rotated secrets exist", func(t *testing.T) {
		secrets := []string{"hash1", "hash2", "hash3"}
		c := &Client{
			RotatedSecrets: secrets,
		}
		hashes := c.GetRotatedHashes()
		require.Len(t, hashes, 3)
		assert.Equal(t, []byte("hash1"), hashes[0])
		assert.Equal(t, []byte("hash2"), hashes[1])
		assert.Equal(t, []byte("hash3"), hashes[2])
	})
}
