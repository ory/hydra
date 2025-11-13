// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/x/uuidx"
)

func Test_pairwiseObfuscate(t *testing.T) {
	salt, sub, clientURI := uuidx.NewV4().String(), uuidx.NewV4().String(), uuidx.NewV4().String()

	t.Run("same result with same parameters", func(t *testing.T) {
		baseLine, err := pairwiseObfuscate(salt, sub, &client.Client{SectorIdentifierURI: clientURI})
		require.NoError(t, err)

		other, err := pairwiseObfuscate(salt, sub, &client.Client{SectorIdentifierURI: clientURI})
		require.NoError(t, err)
		assert.Equal(t, baseLine, other)

		other, err = pairwiseObfuscate(salt, sub, &client.Client{RedirectURIs: []string{"https://" + clientURI}})
		require.NoError(t, err)
		assert.Equal(t, baseLine, other)
	})

	t.Run("different result with different parameters", func(t *testing.T) {
		baseLine, err := pairwiseObfuscate(salt, sub, &client.Client{SectorIdentifierURI: clientURI})
		require.NoError(t, err)

		other, err := pairwiseObfuscate(uuidx.NewV4().String(), sub, &client.Client{SectorIdentifierURI: clientURI})
		require.NoError(t, err)
		assert.NotEqual(t, baseLine, other)

		other, err = pairwiseObfuscate(salt, uuidx.NewV4().String(), &client.Client{SectorIdentifierURI: clientURI})
		require.NoError(t, err)
		assert.NotEqual(t, baseLine, other)

		other, err = pairwiseObfuscate(salt, sub, &client.Client{SectorIdentifierURI: uuidx.NewV4().String()})
		require.NoError(t, err)
		assert.NotEqual(t, baseLine, other)

		other, err = pairwiseObfuscate(salt, sub, &client.Client{RedirectURIs: []string{"https://" + uuidx.NewV4().String()}})
		require.NoError(t, err)
		assert.NotEqual(t, baseLine, other)
	})

	t.Run("errors with invalid client setup", func(t *testing.T) {
		_, err := pairwiseObfuscate(salt, sub, &client.Client{})
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
		_, err = pairwiseObfuscate(salt, sub, &client.Client{RedirectURIs: []string{"https://" + uuidx.NewV4().String(), "https://" + uuidx.NewV4().String()}})
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})
}
