// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowctx_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/oauth2/flowctx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/sqlxx"
)

func TestEncoding(t *testing.T) {
	f := flow.Flow{
		ID:                uuid.Must(uuid.NewV4()).String(),
		NID:               uuid.Must(uuid.NewV4()),
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
			ID:     uuid.Must(uuid.NewV4()).String(),
			NID:    uuid.Must(uuid.NewV4()),
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
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
			AllowedCORSOrigins:      []string{"https://cors1.example.org", "https://cors2.example.org"},
			Metadata:                sqlxx.JSONRawMessage(`{"client-metadata-key1": "val1"}`),
			AccessTokenStrategy:     "jwt",
			SkipConsent:             true,
		},
		RequestURL:         "https://auth.hydra.local/oauth2/auth?client_id=some-client-id&response_type=code&scope=scope1+scope2&redirect_uri=https%3A%2F%2Fredirect1.example.org%2Fcallback&state=some-state&nonce=some-nonce",
		SessionID:          sqlxx.NullString("some-session-id"),
		LoginCSRF:          uuid.Must(uuid.NewV4()).String(),
		LoginInitializedAt: sqlxx.NullTime(time.Now()),
		RequestedAt:        time.Now(),
		State:              1,
		LoginRemember:      true,
		LoginRememberFor:   3600,
		Context:            sqlxx.JSONRawMessage(`{"context-key1": "val1"}`),
		GrantedScope:       []string{"scope1", "scope2"},
		GrantedAudience:    []string{"https://api.example.org/v1", "https://api.example.org/v2"},
		ConsentRemember:    true,
		ConsentRememberFor: pointerx.Int(3600),
		ConsentHandledAt:   sqlxx.NullTime(time.Now()),
		SessionIDToken: sqlxx.MapStringInterface{
			"session-id-token-key1":          "val1",
			"session-id-token-key2":          "val2",
			uuid.Must(uuid.NewV4()).String(): "val3",
			uuid.Must(uuid.NewV4()).String(): "val4",
			uuid.Must(uuid.NewV4()).String(): "val5",
		},
		SessionAccessToken: sqlxx.MapStringInterface{
			"session-access-token-key1":      "val1",
			"session-access-token-key2":      "val2",
			uuid.Must(uuid.NewV4()).String(): "val3",
			uuid.Must(uuid.NewV4()).String(): "val4",
			uuid.Must(uuid.NewV4()).String(): "val5",
		},
	}

	ctx := context.Background()

	t.Run("with client", func(t *testing.T) {
		j, err := json.Marshal(f)
		require.NoError(t, err)
		t.Logf("Length (JSON): %d", len(j))
		cp := new(cipherProvider)
		consentVerifier, err := flowctx.Encode(ctx, cp.FlowCipher(), f, flowctx.AsConsentVerifier)
		require.NoError(t, err)
		t.Logf("Length (JSON+GZIP+AEAD): %d", len(consentVerifier))
	})
	t.Run("without client", func(t *testing.T) {
		f := f
		f.Client = nil
		j, err := json.Marshal(f)
		require.NoError(t, err)
		t.Logf("Length (JSON): %d", len(j))
		cp := new(cipherProvider)
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
