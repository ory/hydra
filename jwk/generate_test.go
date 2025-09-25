// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWK(t *testing.T) {
	t.Parallel()
	jwks, err := GenerateJWK(jose.RS256, "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, jwks.Keys[0].KeyID)
	assert.EqualValues(t, jose.RS256, jwks.Keys[0].Algorithm)
	assert.EqualValues(t, "sig", jwks.Keys[0].Use)
}
