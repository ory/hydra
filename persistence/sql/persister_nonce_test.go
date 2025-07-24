// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/randx"
)

func TestPersister_Nonce(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := testhelpers.NewRegistryMemory(t).Persister()

	accessToken := randx.MustString(100, randx.AlphaNum)
	anotherToken := randx.MustString(100, randx.AlphaNum)
	validNonce, err := p.NewNonce(ctx, accessToken, time.Now().Add(1*time.Hour))
	require.NoError(t, err)

	expiredNonce, err := p.NewNonce(ctx, accessToken, time.Now().Add(-1*time.Hour))
	require.NoError(t, err)

	nonceForAnotherAccessToken, err := p.NewNonce(ctx, anotherToken, time.Now().Add(-1*time.Hour))
	require.NoError(t, err)

	for _, tc := range []struct {
		name      string
		nonce     string
		assertErr assert.ErrorAssertionFunc
	}{{
		name:      "valid nonce",
		nonce:     validNonce,
		assertErr: assert.NoError,
	}, {
		name:      "expired nonce",
		nonce:     expiredNonce,
		assertErr: assertInvalidRequest,
	}, {
		name:      "nonce for another access token",
		nonce:     nonceForAnotherAccessToken,
		assertErr: assertInvalidRequest,
	},
	} {
		t.Run("case="+tc.name, func(t *testing.T) {
			err := p.IsNonceValid(ctx, accessToken, tc.nonce)
			tc.assertErr(t, err)
		})
	}
}

func assertInvalidRequest(t assert.TestingT, err error, i ...interface{}) bool {
	return assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
}
