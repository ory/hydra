// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/sqlcon"
)

func createTestFlow(nid uuid.UUID, state int16) *flow.Flow {
	return &flow.Flow{
		ID:                "a12bf95e-ccfc-45fc-b10d-1358790772c7",
		NID:               nid,
		RequestedScope:    []string{"openid", "profile"},
		RequestedAudience: []string{"https://api.example.org"},
		Subject:           "test-subject",
		Client: &client.Client{
			ID:  "a12bf95e-ccfc-45fc-b10d-1358790772c7",
			NID: nid,
		},
		RequestURL:  "https://example.org/oauth2/auth?client_id=test",
		SessionID:   "session-123",
		RequestedAt: time.Now(),
		State:       state,
	}
}

func TestDecodeFromLoginChallenge(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)
	testFlow := createTestFlow(nid, flow.FlowStateLoginInitialized)

	t.Run("case=successful decode with valid login challenge", func(t *testing.T) {
		loginChallenge, err := testFlow.ToLoginChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		decoded, err := flow.DecodeFromLoginChallenge(ctx, reg, loginChallenge)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.Equal(t, testFlow.ID, decoded.ID)
		assert.Equal(t, testFlow.NID, decoded.NID)
		assert.Equal(t, testFlow.RequestedScope, decoded.RequestedScope)
		assert.Equal(t, testFlow.Subject, decoded.Subject)

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
	})

	t.Run("case=fails with wrong purpose (consent challenge instead of login)", func(t *testing.T) {
		consentChallenge, err := testFlow.ToConsentChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		decoded, err := flow.DecodeFromLoginChallenge(ctx, reg, consentChallenge)
		assert.Error(t, err)
		assert.Nil(t, decoded)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		differentNID := uuid.Must(uuid.NewV4())
		flowWithDifferentNID := &flow.Flow{
			ID:                uuid.Must(uuid.NewV4()).String(),
			NID:               differentNID,
			RequestedScope:    []string{"openid"},
			RequestedAudience: []string{"https://api.example.org"},
			Subject:           "test-subject",
			Client: &client.Client{
				ID:  uuid.Must(uuid.NewV4()).String(),
				NID: differentNID,
			},
			RequestURL:  "https://example.org/oauth2/auth",
			RequestedAt: time.Now(),
			State:       flow.FlowStateLoginInitialized,
		}

		loginChallenge, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsLoginChallenge)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		_, err = flow.DecodeFromLoginChallenge(ctx, reg, loginChallenge)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with expired request", func(t *testing.T) {
		expiredFlow := &flow.Flow{
			ID:                uuid.Must(uuid.NewV4()).String(),
			NID:               nid,
			RequestedScope:    []string{"openid"},
			RequestedAudience: []string{"https://api.example.org"},
			Subject:           "test-subject",
			Client: &client.Client{
				ID:  uuid.Must(uuid.NewV4()).String(),
				NID: nid,
			},
			RequestURL:  "https://example.org/oauth2/auth",
			RequestedAt: time.Now().Add(-2 * time.Hour), // 2 hours ago, beyond the 1 hour max age
			State:       flow.FlowStateLoginInitialized,
		}

		loginChallenge, err := expiredFlow.ToLoginChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		_, err = flow.DecodeFromLoginChallenge(ctx, reg, loginChallenge)
		assert.ErrorIs(t, err, fosite.ErrRequestUnauthorized)
	})

	t.Run("case=fails with invalid challenge format", func(t *testing.T) {
		_, err := flow.DecodeFromLoginChallenge(ctx, reg, "invalid-challenge")
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with empty challenge", func(t *testing.T) {
		_, err := flow.DecodeFromLoginChallenge(ctx, reg, "")
		assert.ErrorIs(t, err, x.ErrNotFound)
	})
}

func TestDecodeFromConsentChallenge(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)
	testFlow := createTestFlow(nid, flow.FlowStateConsentInitialized)

	t.Run("case=successful decode with valid consent challenge", func(t *testing.T) {
		consentChallenge, err := testFlow.ToConsentChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		decoded, err := flow.DecodeFromConsentChallenge(ctx, reg, consentChallenge)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.Equal(t, testFlow.ID, decoded.ID)
		assert.Equal(t, testFlow.NID, decoded.NID)
		assert.Equal(t, testFlow.RequestedScope, decoded.RequestedScope)
		assert.Equal(t, testFlow.Subject, decoded.Subject)

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
	})

	t.Run("case=fails with wrong purpose (login challenge instead of consent)", func(t *testing.T) {
		loginChallenge, err := testFlow.ToLoginChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		decoded, err := flow.DecodeFromConsentChallenge(ctx, reg, loginChallenge)
		assert.Error(t, err)
		assert.Nil(t, decoded)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		differentNID := uuid.Must(uuid.NewV4())
		flowWithDifferentNID := &flow.Flow{
			ID:                uuid.Must(uuid.NewV4()).String(),
			NID:               differentNID,
			RequestedScope:    []string{"openid"},
			RequestedAudience: []string{"https://api.example.org"},
			Subject:           "test-subject",
			Client: &client.Client{
				ID:  uuid.Must(uuid.NewV4()).String(),
				NID: differentNID,
			},
			RequestURL:  "https://example.org/oauth2/auth",
			RequestedAt: time.Now(),
			State:       flow.FlowStateConsentInitialized,
		}

		consentChallenge, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsConsentChallenge)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		_, err = flow.DecodeFromConsentChallenge(ctx, reg, consentChallenge)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with expired request", func(t *testing.T) {
		expiredFlow := &flow.Flow{
			ID:                uuid.Must(uuid.NewV4()).String(),
			NID:               nid,
			RequestedScope:    []string{"openid"},
			RequestedAudience: []string{"https://api.example.org"},
			Subject:           "test-subject",
			Client: &client.Client{
				ID:  uuid.Must(uuid.NewV4()).String(),
				NID: nid,
			},
			RequestURL:  "https://example.org/oauth2/auth",
			RequestedAt: time.Now().Add(-2 * time.Hour), // 2 hours ago, beyond the 1 hour max age
			State:       flow.FlowStateConsentInitialized,
		}

		consentChallenge, err := expiredFlow.ToConsentChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		_, err = flow.DecodeFromConsentChallenge(ctx, reg, consentChallenge)
		assert.ErrorIs(t, err, fosite.ErrRequestUnauthorized)
	})

	t.Run("case=fails with invalid challenge format", func(t *testing.T) {
		_, err := flow.DecodeFromConsentChallenge(ctx, reg, "invalid-challenge")
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with empty challenge", func(t *testing.T) {
		_, err := flow.DecodeFromConsentChallenge(ctx, reg, "")
		assert.ErrorIs(t, err, x.ErrNotFound)
	})
}

func TestDecodeAndInvalidateLoginVerifier(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)

	t.Run("case=successful decode and invalidate with valid login verifier", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)
		testFlow.LoginWasUsed = false

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginVerifier)

		decoded, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		require.NoError(t, err)

		// Verify that InvalidateLoginRequest was called
		assert.True(t, decoded.LoginWasUsed, "LoginWasUsed should be true after invalidation")
		assert.Equal(t, flow.FlowStateLoginUsed, decoded.State, "State should be FlowStateLoginUsed after invalidation")

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
	})

	t.Run("case=fails when flow has already been used", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)
		testFlow.LoginWasUsed = true

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with invalid flow state", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentInitialized)

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with wrong purpose (login challenge instead of verifier)", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)

		loginChallenge, err := testFlow.ToLoginChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginChallenge)
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})

	t.Run("case=fails with wrong purpose (consent challenge instead of verifier)", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentInitialized)

		consentChallenge, err := testFlow.ToConsentChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, consentChallenge)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		differentNID := uuid.Must(uuid.NewV4())
		flowWithDifferentNID := createTestFlow(differentNID, flow.FlowStateLoginUnused)

		loginVerifier, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsLoginVerifier)
		require.NoError(t, err)
		require.NotEmpty(t, loginVerifier)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})

	t.Run("case=fails with invalid verifier format", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, "invalid-verifier")
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})

	t.Run("case=fails with empty verifier", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, "")
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	})


	t.Run("case=works with FlowStateLoginError", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginError)

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginVerifier)

		decoded, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.True(t, decoded.LoginWasUsed)
		assert.Equal(t, flow.FlowStateLoginUsed, decoded.State)
	})
}
