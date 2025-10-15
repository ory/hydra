// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/flow"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/sqlxx"
)

func TestEncoding(t *testing.T) {
	f := flow.Flow{
		ID:                "test-flow-id",
		NID:               uuid.FromStringOrNil("735c9c15-3d07-4501-9800-4e5e0599e57b"),
		RequestedScope:    []string{"scope1", "scope2"},
		RequestedAudience: []string{"https://api.example.org/v1", "https://api.example.org/v2"},
		LoginSkip:         false,
		Subject:           "some-subject@some-idp-somewhere.com",
		OpenIDConnectContext: &flow.OAuth2ConsentRequestOpenIDConnectContext{
			ACRValues:         []string{"acr1", "acr2"},
			UILocales:         []string{"en-US", "en-GB"},
			Display:           "page",
			IDTokenHintClaims: map[string]interface{}{"claim1": "value1", "claim2": "value2"},
			LoginHint:         "some-login-hint",
		},
		Client: &client.Client{
			ID:     "test-client-id",
			NID:    uuid.FromStringOrNil("735c9c15-3d07-4501-9800-4e5e0599e57b"),
			Name:   "some-client-name",
			Secret: "some-supersafe-secret",
			RedirectURIs: []string{
				"https://redirect1.example.org/callback",
				"https://redirect2.example.org/callback",
			},
			GrantTypes:              []string{"authorization_code", "refresh_token"},
			ResponseTypes:           []string{"code"},
			Scope:                   "scope1 scope2",
			Audience:                sqlxx.StringSliceJSONFormat{"https://api.example.org/v1 https://api.example.org/v2"},
			Owner:                   "some-owner",
			TermsOfServiceURI:       "https://tos.example.org",
			PolicyURI:               "https://policy.example.org",
			ClientURI:               "https://client.example.org",
			LogoURI:                 "https://logo.example.org",
			Contacts:                []string{"contact1", "contact2"},
			SubjectType:             "public",
			JSONWebKeysURI:          "https://jwks.example.org",
			JSONWebKeys:             nil, // TODO?
			TokenEndpointAuthMethod: "client_secret_basic",
			CreatedAt:               time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:               time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			AllowedCORSOrigins:      []string{"https://cors1.example.org", "https://cors2.example.org"},
			Metadata:                sqlxx.JSONRawMessage(`{"client-metadata-key1": "val1"}`),
			AccessTokenStrategy:     "jwt",
			SkipConsent:             true,
		},
		RequestURL:         "https://auth.hydra.local/oauth2/auth?client_id=some-client-id&response_type=code&scope=scope1+scope2&redirect_uri=https%3A%2F%2Fredirect1.example.org%2Fcallback&state=some-state&nonce=some-nonce",
		SessionID:          sqlxx.NullString("some-session-id"),
		LoginCSRF:          "test-login-csrf",
		RequestedAt:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		State:              1,
		LoginRemember:      true,
		LoginRememberFor:   3600,
		Context:            sqlxx.JSONRawMessage(`{"context-key1": "val1"}`),
		GrantedScope:       []string{"scope1", "scope2"},
		GrantedAudience:    []string{"https://api.example.org/v1", "https://api.example.org/v2"},
		ConsentRemember:    true,
		ConsentRememberFor: pointerx.Int(3600),
		ConsentHandledAt:   sqlxx.NullTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		SessionIDToken: sqlxx.MapStringInterface{
			"session-id-token-key1": "val1",
			"session-id-token-key2": "val2",
			"session-id-token-key3": "val3",
			"session-id-token-key4": "val4",
			"session-id-token-key5": "val5",
		},
		SessionAccessToken: sqlxx.MapStringInterface{
			"session-access-token-key1": "val1",
			"session-access-token-key2": "val2",
			"session-access-token-key3": "val3",
			"session-access-token-key4": "val4",
			"session-access-token-key5": "val5",
		},
	}

	ctx := context.Background()
	cp := new(cipherProvider)

	t.Run("encode and decode with snapshots", func(t *testing.T) {
		testCases := []struct {
			name    string
			purpose flow.CodecOption
		}{
			{"login challenge", flow.AsLoginChallenge},
			{"login verifier", flow.AsLoginVerifier},
			{"consent challenge", flow.AsConsentChallenge},
			{"consent verifier", flow.AsConsentVerifier},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				encoded, err := flow.Encode(ctx, cp.FlowCipher(), f, tc.purpose)
				require.NoError(t, err)

				decoded, err := flow.Decode[flow.Flow](ctx, cp.FlowCipher(), encoded, tc.purpose)
				require.NoError(t, err)
				snapshotx.SnapshotT(t, decoded, snapshotx.ExceptPaths("n", "ia"))
			})
		}
	})

	t.Run("purpose validation", func(t *testing.T) {
		testCases := []struct {
			name          string
			encodePurpose flow.CodecOption
			decodePurpose flow.CodecOption
		}{
			{"login challenge decoded as login verifier", flow.AsLoginChallenge, flow.AsLoginVerifier},
			{"login verifier decoded as login challenge", flow.AsLoginVerifier, flow.AsLoginChallenge},
			{"consent challenge decoded as consent verifier", flow.AsConsentChallenge, flow.AsConsentVerifier},
			{"consent verifier decoded as consent challenge", flow.AsConsentVerifier, flow.AsConsentChallenge},
			{"login challenge decoded as consent challenge", flow.AsLoginChallenge, flow.AsConsentChallenge},
			{"consent challenge decoded as login challenge", flow.AsConsentChallenge, flow.AsLoginChallenge},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				encoded, err := flow.Encode(ctx, cp.FlowCipher(), f, tc.encodePurpose)
				require.NoError(t, err)

				_, err = flow.Decode[flow.Flow](ctx, cp.FlowCipher(), encoded, tc.decodePurpose)
				assert.Error(t, err, "decoding with wrong purpose should fail")
			})
		}
	})

	t.Run("with client", func(t *testing.T) {
		j, err := json.Marshal(f)
		require.NoError(t, err)
		t.Logf("Length (JSON): %d", len(j))
		consentVerifier, err := flow.Encode(ctx, cp.FlowCipher(), f, flow.AsConsentVerifier)
		require.NoError(t, err)
		t.Logf("Length (JSON+GZIP+AEAD): %d", len(consentVerifier))
	})

	t.Run("without client", func(t *testing.T) {
		f := f
		f.Client = nil
		j, err := json.Marshal(f)
		require.NoError(t, err)
		t.Logf("Length (JSON): %d", len(j))
		consentVerifier, err := f.ToConsentVerifier(ctx, cp)
		require.NoError(t, err)
		t.Logf("Length (JSON+GZIP+AEAD): %d", len(consentVerifier))
	})
}

type cipherProvider struct{}

func (c *cipherProvider) FlowCipher() *aead.XChaCha20Poly1305 {
	return aead.NewXChaCha20Poly1305(c)
}

func (c *cipherProvider) GetGlobalSecret(context.Context) ([]byte, error) {
	return []byte("supersecret123456789123456789012"), nil
}

func (c *cipherProvider) GetRotatedGlobalSecrets(ctx context.Context) ([][]byte, error) {
	return nil, nil
}
