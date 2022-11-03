// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
)

func TestGenerateJWK(t *testing.T) {
	jwks, err := GenerateJWK(context.Background(), jose.RS256, "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, jwks.Keys[0].KeyID)
	assert.EqualValues(t, jose.RS256, jwks.Keys[0].Algorithm)
	assert.EqualValues(t, "sig", jwks.Keys[0].Use)
}
