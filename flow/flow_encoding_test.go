// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow_test

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
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
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/servicelocatorx"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/sqlxx"
)

func createTestFlow(nid uuid.UUID, state flow.State) *flow.Flow {
	return &flow.Flow{
		ID:                "a12bf95e-ccfc-45fc-b10d-1358790772c7",
		NID:               nid,
		RequestedScope:    []string{"openid", "profile"},
		RequestedAudience: []string{"https://api.example.org"},
		LoginSkip:         true,
		Subject:           "test-subject",
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues:         []string{"http://acrvalues.example.org"},
			UILocales:         []string{"en-US", "en-GB"},
			Display:           "page",
			IDTokenHintClaims: map[string]interface{}{"email": "user@example.org"},
			LoginHint:         "login-hint",
		},
		Client: &client.Client{
			ID:  "a12bf95e-ccfc-45fc-b10d-1358790772c7",
			NID: nid,
		},
		ClientID:                   "a12bf95e-ccfc-45fc-b10d-1358790772c7",
		RequestURL:                 "https://example.org/oauth2/auth?client_id=test",
		SessionID:                  "session-123",
		IdentityProviderSessionID:  "session-id",
		LoginCSRF:                  "login-csrf",
		RequestedAt:                time.Now(),
		State:                      state,
		LoginRemember:              true,
		LoginRememberFor:           3000,
		LoginExtendSessionLifespan: true,
		ACR:                        "http://acrvalues.example.org",
		AMR:                        []string{"pwd"},
		ForceSubjectIdentifier:     "forced-subject",
		Context:                    sqlxx.JSONRawMessage(`{"foo":"bar"}`),
		LoginAuthenticatedAt:       sqlxx.NullTime(time.Date(2025, 10, 9, 12, 52, 0, 0, time.UTC)),
		DeviceChallengeID:          "device-challenge",
		DeviceCodeRequestID:        "device-code-request",
		DeviceCSRF:                 "device-csrf",
		DeviceHandledAt:            sqlxx.NullTime{},
		ConsentRequestID:           "consent-request",
		ConsentSkip:                true,
		ConsentCSRF:                "consent-csrf",
		GrantedScope:               []string{"openid"},
		GrantedAudience:            []string{"https://api.example.org"},
		ConsentRemember:            true,
		ConsentRememberFor:         pointerx.Ptr(3000),
		ConsentHandledAt:           sqlxx.NullTime{},
		SessionIDToken:             map[string]interface{}{"sub": "test-subject", "foo": "bar"},
		SessionAccessToken:         map[string]interface{}{"scp": []string{"openid", "profile"}, "aud": []string{"https://api.example.org"}},
	}
}

func TestDecodeFromLoginChallenge(t *testing.T) {
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)
	testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)

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

		t.Run("decodes deterministically", func(t *testing.T) {
			second, err := flow.DecodeFromLoginChallenge(ctx, reg, loginChallenge)
			require.NoError(t, err)
			assert.Equal(t, decoded, second)
		})
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
		flowWithDifferentNID := createTestFlow(uuid.Must(uuid.NewV4()), flow.FlowStateLoginUnused)

		loginChallenge, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsLoginChallenge)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		_, err = flow.DecodeFromLoginChallenge(ctx, reg, loginChallenge)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with expired request", func(t *testing.T) {
		expiredFlow := createTestFlow(nid, flow.FlowStateLoginUnused)
		expiredFlow.RequestedAt = time.Now().Add(-2 * time.Hour)

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
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)
	testFlow := createTestFlow(nid, flow.FlowStateConsentUnused)

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

		t.Run("decodes deterministically", func(t *testing.T) {
			second, err := flow.DecodeFromConsentChallenge(ctx, reg, consentChallenge)
			require.NoError(t, err)
			assert.Equal(t, decoded, second)
		})
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
		flowWithDifferentNID := createTestFlow(uuid.Must(uuid.NewV4()), flow.FlowStateConsentUnused)

		consentChallenge, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsConsentChallenge)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		_, err = flow.DecodeFromConsentChallenge(ctx, reg, consentChallenge)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with expired request", func(t *testing.T) {
		expiredFlow := createTestFlow(nid, flow.FlowStateConsentUnused)
		expiredFlow.RequestedAt = time.Now().Add(-2 * time.Hour)

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
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)

	t.Run("case=successful decode and invalidate with valid login verifier", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginVerifier)

		decoded, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		require.NoError(t, err)

		// Verify that InvalidateLoginRequest was called
		assert.Equal(t, flow.FlowStateLoginUsed, decoded.State, "State should be FlowStateLoginUsed after invalidation")

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
	})

	t.Run("case=fails when flow has already been used", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUsed)

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with invalid flow state", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentUnused)

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
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		differentNID := uuid.Must(uuid.NewV4())
		flowWithDifferentNID := createTestFlow(differentNID, flow.FlowStateLoginUnused)

		loginVerifier, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsLoginVerifier)
		require.NoError(t, err)
		require.NotEmpty(t, loginVerifier)

		_, err = flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with invalid verifier format", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, "invalid-verifier")
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with empty verifier", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, "")
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=works with FlowStateLoginError", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginError)

		loginVerifier, err := testFlow.ToLoginVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginVerifier)

		decoded, err := flow.DecodeAndInvalidateLoginVerifier(ctx, reg, loginVerifier)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.Equal(t, flow.FlowStateLoginError, decoded.State)
	})
}

func TestDecodeFromDeviceChallenge(t *testing.T) {
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)
	testFlow := createTestFlow(nid, flow.DeviceFlowStateUnused)

	t.Run("case=successful decode with valid device challenge", func(t *testing.T) {
		deviceChallenge, err := testFlow.ToDeviceChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, deviceChallenge)

		decoded, err := flow.DecodeFromDeviceChallenge(ctx, reg, deviceChallenge)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.Equal(t, testFlow.ID, decoded.ID)
		assert.Equal(t, testFlow.NID, decoded.NID)
		assert.Equal(t, testFlow.RequestedScope, decoded.RequestedScope)
		assert.Equal(t, testFlow.Subject, decoded.Subject)

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))

		t.Run("decodes deterministically", func(t *testing.T) {
			second, err := flow.DecodeFromDeviceChallenge(ctx, reg, deviceChallenge)
			require.NoError(t, err)
			assert.Equal(t, decoded, second)
		})
	})

	t.Run("case=fails with wrong purpose (login challenge instead of device)", func(t *testing.T) {
		loginChallenge, err := testFlow.ToLoginChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, loginChallenge)

		decoded, err := flow.DecodeFromDeviceChallenge(ctx, reg, loginChallenge)
		assert.Error(t, err)
		assert.Nil(t, decoded)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		flowWithDifferentNID := createTestFlow(uuid.Must(uuid.NewV4()), flow.DeviceFlowStateUnused)

		deviceChallenge, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsDeviceChallenge)
		require.NoError(t, err)
		require.NotEmpty(t, deviceChallenge)

		_, err = flow.DecodeFromDeviceChallenge(ctx, reg, deviceChallenge)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with expired request", func(t *testing.T) {
		expiredFlow := createTestFlow(nid, flow.DeviceFlowStateUnused)
		expiredFlow.RequestedAt = time.Now().Add(-2 * time.Hour)

		deviceChallenge, err := expiredFlow.ToDeviceChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, deviceChallenge)

		_, err = flow.DecodeFromDeviceChallenge(ctx, reg, deviceChallenge)
		assert.ErrorIs(t, err, fosite.ErrRequestUnauthorized)
	})

	t.Run("case=fails with invalid challenge format", func(t *testing.T) {
		_, err := flow.DecodeFromDeviceChallenge(ctx, reg, "invalid-challenge")
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=fails with empty challenge", func(t *testing.T) {
		_, err := flow.DecodeFromDeviceChallenge(ctx, reg, "")
		assert.ErrorIs(t, err, x.ErrNotFound)
	})
}

func TestDecodeAndInvalidateDeviceVerifier(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)

	t.Run("case=successful decode and invalidate with valid device verifier", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.DeviceFlowStateUnused)

		deviceVerifier, err := testFlow.ToDeviceVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, deviceVerifier)

		decoded, err := flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, deviceVerifier)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.Equal(t, flow.DeviceFlowStateUsed, decoded.State, "State should be DeviceFlowStateUsed after invalidation")

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
	})

	t.Run("case=fails when flow has already been used", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.DeviceFlowStateUsed)

		deviceVerifier, err := testFlow.ToDeviceVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, deviceVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with invalid flow state", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)

		deviceVerifier, err := testFlow.ToDeviceVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, deviceVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with wrong purpose (device challenge instead of verifier)", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.DeviceFlowStateUnused)

		deviceChallenge, err := testFlow.ToDeviceChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, deviceChallenge)

		_, err = flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, deviceChallenge)
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		differentNID := uuid.Must(uuid.NewV4())
		flowWithDifferentNID := createTestFlow(differentNID, flow.DeviceFlowStateUnused)

		deviceVerifier, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsDeviceVerifier)
		require.NoError(t, err)
		require.NotEmpty(t, deviceVerifier)

		_, err = flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, deviceVerifier)
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with invalid verifier format", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, "invalid-verifier")
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with empty verifier", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateDeviceVerifier(ctx, reg, "")
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})
}

func TestDecodeAndInvalidateConsentVerifier(t *testing.T) {
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(
		configx.WithValue(config.KeyConsentRequestMaxAge, time.Hour),
	))

	nid := reg.Networker().NetworkID(ctx)

	t.Run("case=successful decode and invalidate with valid consent verifier", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentUnused)

		consentVerifier, err := testFlow.ToConsentVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentVerifier)

		decoded, err := flow.DecodeAndInvalidateConsentVerifier(ctx, reg, consentVerifier)
		require.NoError(t, err)

		// Verify that InvalidateConsentRequest was called
		assert.Equal(t, flow.FlowStateConsentUsed, decoded.State, "State should be FlowStateConsentUsed after invalidation")

		snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
	})

	t.Run("case=fails when flow has already been used", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentUsed)

		consentVerifier, err := testFlow.ToConsentVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateConsentVerifier(ctx, reg, consentVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with invalid flow state", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateLoginUnused)

		consentVerifier, err := testFlow.ToConsentVerifier(ctx, reg)
		require.NoError(t, err)

		_, err = flow.DecodeAndInvalidateConsentVerifier(ctx, reg, consentVerifier)
		assert.ErrorIs(t, err, fosite.ErrInvalidRequest)
	})

	t.Run("case=fails with wrong purpose (consent challenge instead of verifier)", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentUnused)

		consentChallenge, err := testFlow.ToConsentChallenge(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentChallenge)

		_, err = flow.DecodeAndInvalidateConsentVerifier(ctx, reg, consentChallenge)
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with different network ID", func(t *testing.T) {
		differentNID := uuid.Must(uuid.NewV4())
		flowWithDifferentNID := createTestFlow(differentNID, flow.FlowStateConsentUnused)

		consentVerifier, err := flow.Encode(ctx, reg.FlowCipher(), flowWithDifferentNID, flow.AsConsentVerifier)
		require.NoError(t, err)
		require.NotEmpty(t, consentVerifier)

		_, err = flow.DecodeAndInvalidateConsentVerifier(ctx, reg, consentVerifier)
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with invalid verifier format", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateConsentVerifier(ctx, reg, "invalid-verifier")
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=fails with empty verifier", func(t *testing.T) {
		_, err := flow.DecodeAndInvalidateConsentVerifier(ctx, reg, "")
		assert.ErrorIs(t, err, fosite.ErrAccessDenied)
	})

	t.Run("case=works with FlowStateConsentError", func(t *testing.T) {
		testFlow := createTestFlow(nid, flow.FlowStateConsentError)

		consentVerifier, err := testFlow.ToConsentVerifier(ctx, reg)
		require.NoError(t, err)
		require.NotEmpty(t, consentVerifier)

		decoded, err := flow.DecodeAndInvalidateConsentVerifier(ctx, reg, consentVerifier)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		assert.Equal(t, flow.FlowStateConsentError, decoded.State)
	})
}

var (
	//go:embed fixtures/legacy_challenges/*.txt
	LegacyChallenges    embed.FS
	legacyChallengesNID = uuid.Must(uuid.FromString("34b4dd42-f02b-4448-b066-8e4e6655c0bb"))
)

func TestCanUseLegacyChallenges(t *testing.T) {
	reg := testhelpers.NewRegistryMemory(t,
		driver.WithConfigOptions(
			configx.WithValue(config.KeyGetSystemSecret, []string{"well-known-fixture-secret"}),
			configx.WithValue(config.KeyConsentRequestMaxAge, 100*365*24*time.Hour), // 100 years, effectively disabling expiration
		),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: legacyChallengesNID})),
	)

	require.NoError(t, fs.WalkDir(LegacyChallenges, "fixtures/legacy_challenges", func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if d.IsDir() {
			return nil
		}
		t.Run(strings.TrimSuffix(d.Name(), ".txt"), func(t *testing.T) {
			content, err := fs.ReadFile(LegacyChallenges, path)
			require.NoError(t, err)

			var f *flow.Flow
			switch {
			case strings.Contains(d.Name(), "login"):
				f, err = flow.DecodeFromLoginChallenge(t.Context(), reg, string(content))
			case strings.Contains(d.Name(), "consent"):
				f, err = flow.DecodeFromConsentChallenge(t.Context(), reg, string(content))
			case strings.Contains(d.Name(), "device"):
				f, err = flow.DecodeFromDeviceChallenge(t.Context(), reg, string(content))
			default:
				t.Fatalf("unknown challenge type in file name: %s", d.Name())
			}
			require.NoErrorf(t, err, "failed to decode challenge from file: %s\n%+v", d.Name(), errors.Unwrap(errors.Unwrap(err)))

			snapshotx.SnapshotT(t, f)
		})
		return nil
	}))
}

func TestUpdateLegacyChallenges(t *testing.T) {
	t.Skip("this test is used to update the fixtures only, they should not be updated unless we have a breaking change (so probably never)")

	reg := testhelpers.NewRegistryMemory(t,
		driver.WithConfigOptions(configx.WithValue(config.KeyGetSystemSecret, []string{"well-known-fixture-secret"})),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.Static{NID: legacyChallengesNID})),
	)

	for name, flowState := range map[string]flow.State{
		"login_initialized":   flow.FlowStateLoginInitialized,
		"login_unused":        flow.FlowStateLoginUnused,
		"login_used":          flow.FlowStateLoginUsed,
		"login_error":         flow.FlowStateLoginError,
		"consent_initialized": flow.FlowStateConsentInitialized,
		"consent_unused":      flow.FlowStateConsentUnused,
		"consent_used":        flow.FlowStateConsentUsed,
		"consent_error":       flow.FlowStateConsentError,
		"device_initialized":  flow.DeviceFlowStateInitialized,
		"device_unused":       flow.DeviceFlowStateUnused,
		"device_used":         flow.DeviceFlowStateUsed,
	} {
		f := createTestFlow(legacyChallengesNID, flowState)
		var challenge string
		var err error
		switch flowState {
		case flow.FlowStateLoginInitialized, flow.FlowStateLoginUnused, flow.FlowStateLoginUsed, flow.FlowStateLoginError:
			challenge, err = f.ToLoginChallenge(t.Context(), reg)
		case flow.FlowStateConsentInitialized, flow.FlowStateConsentUnused, flow.FlowStateConsentUsed, flow.FlowStateConsentError:
			challenge, err = f.ToConsentChallenge(t.Context(), reg)
		case flow.DeviceFlowStateInitialized, flow.DeviceFlowStateUnused, flow.DeviceFlowStateUsed:
			challenge, err = f.ToDeviceChallenge(t.Context(), reg)
		default:
			t.Fatalf("unknown flow state: %d", flowState)
		}
		require.NoError(t, err)

		require.NoError(t, os.WriteFile(fmt.Sprintf("fixtures/legacy_challenges/%s.txt", name), []byte(challenge), 0644))
	}
}
