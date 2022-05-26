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

	actual, err = ToAccessTokenStrategyType("invalid")
	require.Error(t, err)
}
