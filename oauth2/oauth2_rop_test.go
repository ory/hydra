// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/jwt"
	hydra "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/fositex"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/internal/testhelpers"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/sqlxx"
)

func TestResourceOwnerPasswordGrant(t *testing.T) {
	ctx := context.Background()
	fakeKratos := kratos.NewFake()
	reg := testhelpers.NewMockedRegistry(t, &contextx.Default{})
	reg.WithKratos(fakeKratos)
	reg.WithExtraFositeFactories([]fositex.Factory{compose.OAuth2ResourceOwnerPasswordCredentialsFactory})
	publicTS, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)

	secret := uuid.New().String()
	audience := sqlxx.StringSliceJSONFormat{"https://aud.example.com"}
	client := &hydra.Client{
		Secret:     secret,
		GrantTypes: []string{"password", "refresh_token"},
		Scope:      "offline",
		Audience:   audience,
		Lifespans: hydra.Lifespans{
			PasswordGrantAccessTokenLifespan:  x.NullDuration{Duration: 1 * time.Hour, Valid: true},
			PasswordGrantRefreshTokenLifespan: x.NullDuration{Duration: 1 * time.Hour, Valid: true},
		},
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, client))

	oauth2Config := &oauth2.Config{
		ClientID:     client.GetID(),
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   reg.Config().OAuth2AuthURL(ctx).String(),
			TokenURL:  reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		Scopes: []string{"offline"},
	}

	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")
		assert.Equal(t, r.Header.Get("Authorization"), "Bearer secret value")

		var hookReq hydraoauth2.TokenHookRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
		assert.NotEmpty(t, hookReq.Session)
		assert.NotEmpty(t, hookReq.Request)

		claims := hookReq.Session.Extra
		claims["hooked"] = true
		if hookReq.Request.GrantTypes[0] == "refresh_token" {
			claims["refreshed"] = true
		}

		hookResp := hydraoauth2.TokenHookResponse{
			Session: flow.AcceptOAuth2ConsentRequestSession{
				AccessToken: claims,
				IDToken:     claims,
			},
		}

		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(&hookResp))
	}))
	defer hs.Close()

	reg.Config().MustSet(ctx, config.KeyTokenHook, &config.HookConfig{
		URL: hs.URL,
		Auth: &config.Auth{
			Type: "api_key",
			Config: config.AuthConfig{
				In:    "header",
				Name:  "Authorization",
				Value: "Bearer secret value",
			},
		},
	})
	reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")

	t.Run("case=get ROP grant token with valid username and password", func(t *testing.T) {
		token, err := oauth2Config.PasswordCredentialsToken(ctx, kratos.FakeUsername, kratos.FakePassword)
		require.NoError(t, err)
		require.NotEmpty(t, token.AccessToken)

		// Access token should have hook and identity_id claims
		jwtAT, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
			return reg.AccessTokenJWTStrategy().GetPublicKey(ctx)
		})
		require.NoError(t, err)
		assert.Equal(t, kratos.FakeUsername, jwtAT.Claims["ext"].(map[string]any)["username"])
		assert.Equal(t, kratos.FakeIdentityID, jwtAT.Claims["sub"])
		assert.Equal(t, publicTS.URL, jwtAT.Claims["iss"])
		assert.True(t, jwtAT.Claims["ext"].(map[string]any)["hooked"].(bool))
		assert.ElementsMatch(t, audience, jwtAT.Claims["aud"])

		t.Run("case=introspect token", func(t *testing.T) {
			// Introspected token should have hook and identity_id claims
			i := testhelpers.IntrospectToken(t, oauth2Config, token.AccessToken, adminTS)
			assert.True(t, i.Get("active").Bool(), "%s", i)
			assert.Equal(t, kratos.FakeUsername, i.Get("ext.username").String(), "%s", i)
			assert.Equal(t, kratos.FakeIdentityID, i.Get("sub").String(), "%s", i)
			assert.True(t, i.Get("ext.hooked").Bool(), "%s", i)
			assert.EqualValues(t, oauth2Config.ClientID, i.Get("client_id").String(), "%s", i)
		})

		t.Run("case=refresh token", func(t *testing.T) {
			// Refreshed access token should have hook and identity_id claims
			require.NotEmpty(t, token.RefreshToken)
			token.Expiry = token.Expiry.Add(-time.Hour * 24)
			refreshedToken, err := oauth2Config.TokenSource(context.Background(), token).Token()
			require.NoError(t, err)

			require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
			require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)

			jwtAT, err := jwt.Parse(refreshedToken.AccessToken, func(token *jwt.Token) (interface{}, error) {
				return reg.AccessTokenJWTStrategy().GetPublicKey(ctx)
			})
			require.NoError(t, err)
			assert.Equal(t, kratos.FakeIdentityID, jwtAT.Claims["sub"])
			assert.Equal(t, kratos.FakeUsername, jwtAT.Claims["ext"].(map[string]any)["username"])
			assert.True(t, jwtAT.Claims["ext"].(map[string]any)["hooked"].(bool))
			assert.True(t, jwtAT.Claims["ext"].(map[string]any)["refreshed"].(bool))
		})
	})

	t.Run("case=access denied for invalid password", func(t *testing.T) {
		_, err := oauth2Config.PasswordCredentialsToken(ctx, kratos.FakeUsername, "invalid")
		retrieveError := new(oauth2.RetrieveError)
		require.Error(t, err)
		require.ErrorAs(t, err, &retrieveError)
		assert.Contains(t, retrieveError.ErrorDescription, "Unable to authenticate the provided username and password credentials")
	})
}
