// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToAccessTokenStrategyType(t *testing.T) {
	actual, err := ToAccessTokenStrategyType("opaque")
	require.NoError(t, err)
	assert.Equal(t, AccessTokenDefaultStrategy, actual)

	actual, err = ToAccessTokenStrategyType("jwt")
	require.NoError(t, err)
	assert.Equal(t, AccessTokenJWTStrategy, actual)

	_, err = ToAccessTokenStrategyType("invalid")
	require.Error(t, err)
}
