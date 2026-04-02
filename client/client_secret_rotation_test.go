// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRotatedHashes(t *testing.T) {
	t.Run("returns nil when no rotated secrets", func(t *testing.T) {
		c := &Client{
			RotatedSecrets: "",
		}
		hashes := c.GetRotatedHashes()
		assert.Nil(t, hashes)
	})

	t.Run("returns hashes when rotated secrets exist", func(t *testing.T) {
		secrets := []string{"hash1", "hash2", "hash3"}
		secretsJSON, err := json.Marshal(secrets)
		require.NoError(t, err)

		c := &Client{
			RotatedSecrets: string(secretsJSON),
		}
		hashes := c.GetRotatedHashes()
		require.Len(t, hashes, 3)
		assert.Equal(t, []byte("hash1"), hashes[0])
		assert.Equal(t, []byte("hash2"), hashes[1])
		assert.Equal(t, []byte("hash3"), hashes[2])
	})

	t.Run("returns nil on invalid JSON", func(t *testing.T) {
		c := &Client{
			RotatedSecrets: "invalid json",
		}
		hashes := c.GetRotatedHashes()
		assert.Nil(t, hashes)
	})
}
