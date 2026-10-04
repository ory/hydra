// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/fosite"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
)

func refreshTokenSignature(t testing.TB, reg *driver.RegistrySQL, token string) string {
	t.Helper()
	signature := reg.OAuth2HMACStrategy().RefreshTokenSignature(t.Context(), token)
	require.NotEmpty(t, signature)
	return signature
}

func backdateRefreshTokenFirstUsedAt(t testing.TB, reg *driver.RegistrySQL, token string, elapsed time.Duration) {
	t.Helper()
	backdateRefreshTokenFirstUsedAtBySignature(t, reg, refreshTokenSignature(t, reg, token), elapsed)
}

func backdateRefreshTokenFirstUsedAtBySignature(t testing.TB, reg *driver.RegistrySQL, signature string, elapsed time.Duration) {
	t.Helper()
	updated, err := reg.Persister().Connection(t.Context()).RawQuery(
		"UPDATE hydra_oauth2_refresh SET first_used_at = ? WHERE signature = ?",
		time.Now().UTC().Add(-elapsed),
		signature,
	).ExecWithCount()
	require.NoError(t, err)
	require.Equal(t, 1, updated)
}

func expireRefreshToken(t testing.TB, reg *driver.RegistrySQL, token string) {
	t.Helper()
	ctx := t.Context()
	signature := refreshTokenSignature(t, reg, token)
	expiredAt := time.Now().UTC().Add(-time.Minute)

	var row struct {
		Session []byte `db:"session_data"`
	}
	require.NoError(t, reg.Persister().Connection(ctx).RawQuery(
		"SELECT session_data FROM hydra_oauth2_refresh WHERE signature = ?",
		signature,
	).First(&row))

	sessionData := row.Session
	if reg.Config().EncryptSessionData(ctx) {
		var err error
		sessionData, err = reg.KeyCipher().Decrypt(ctx, string(sessionData), nil)
		require.NoError(t, err)
	}

	var session hydraoauth2.Session
	require.NoError(t, json.Unmarshal(sessionData, &session))
	session.SetExpiresAt(fosite.RefreshToken, expiredAt)
	sessionData, err := json.Marshal(&session)
	require.NoError(t, err)

	if reg.Config().EncryptSessionData(ctx) {
		encrypted, err := reg.KeyCipher().Encrypt(ctx, sessionData, nil)
		require.NoError(t, err)
		sessionData = []byte(encrypted)
	}

	updated, err := reg.Persister().Connection(ctx).RawQuery(
		"UPDATE hydra_oauth2_refresh SET expires_at = ?, session_data = ? WHERE signature = ?",
		expiredAt,
		sessionData,
		signature,
	).ExecWithCount()
	require.NoError(t, err)
	require.Equal(t, 1, updated)
}
