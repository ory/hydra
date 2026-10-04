// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/configx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/ioutilx"
)

func TestSecretRotationSDK(t *testing.T) {
	ctx := t.Context()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeySubjectTypesSupported: []string{"public"},
		config.KeyDefaultClientScope:    []string{"openid"},
		"log.level":                     "fatal",
	})))

	routerAdmin := httprouterx.NewRouterAdminWithPrefix()
	routerPublic := httprouterx.NewRouterPublic()
	clHandler := client.NewHandler(reg)
	clHandler.SetPublicRoutes(routerPublic)
	clHandler.SetAdminRoutes(routerAdmin)
	o2Handler := oauth2.NewHandler(reg)
	o2Handler.SetPublicRoutes(routerPublic, func(h http.Handler) http.Handler { return h })
	o2Handler.SetAdminRoutes(routerAdmin)

	adminServer := httptest.NewServer(routerAdmin)
	t.Cleanup(adminServer.Close)
	publicServer := httptest.NewServer(routerPublic)
	t.Cleanup(publicServer.Close)

	reg.Config().MustSet(ctx, config.KeyAdminURL, adminServer.URL)
	reg.Config().MustSet(ctx, config.KeyOAuth2TokenURL, publicServer.URL+"/oauth2/token")

	sdk := hydra.NewAPIClient(hydra.NewConfiguration())
	sdk.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminServer.URL}}

	tokenURL := reg.Config().OAuth2TokenURL(ctx).String()

	newClientCredentialsConfig := func(clientID, secret string) *clientcredentials.Config {
		return &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: secret,
			TokenURL:     tokenURL,
			AuthStyle:    goauth2.AuthStyleInHeader,
		}
	}

	assertNoRotatedSecretsInResponse := func(t *testing.T, resp *http.Response) {
		t.Helper()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var raw map[string]any
		require.NoError(t, json.Unmarshal(body, &raw))
		assert.NotContains(t, raw, "rotated", "Response should not contain rotated secrets")
	}

	t.Run("case=secret rotation max old secrets", func(t *testing.T) {
		created, _, err := sdk.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(hydra.OAuth2Client{
			ClientSecret:            new("unused"),
			GrantTypes:              []string{"client_credentials"},
			Scope:                   new("openid"),
			TokenEndpointAuthMethod: new("client_secret_basic"),
		}).Execute()
		require.NoError(t, err)
		clientID := *created.ClientId
		require.NotNil(t, created.ClientSecret)
		assert.Equal(t, "unused", *created.ClientSecret)

		secrets := make([]string, 0, 8)
		for range cap(secrets) {
			rotated, _, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, clientID).Execute()
			require.NoError(t, err)
			require.NotNil(t, rotated.ClientSecret)
			newSecret := *rotated.ClientSecret
			assert.NotEmpty(t, newSecret)
			secrets = append(secrets, newSecret)
		}
		require.Len(t, secrets, 8)

		// First 2 secrets should no longer work
		for i, secret := range secrets[:2] {
			_, err = newClientCredentialsConfig(clientID, secret).Token(ctx)
			assert.Error(t, err, "Old secret %d should not work", i)
		}
		// Most recent 6 (one active, five rotated) should work
		for i, secret := range secrets[2:] {
			token, err := newClientCredentialsConfig(clientID, secret).Token(ctx)
			require.NoError(t, err, "Secret %d should work", i+2)
			assert.NotEmpty(t, token.AccessToken)
		}
	})

	t.Run("case=secret rotation lifecycle", func(t *testing.T) {
		// Create a client with a known secret
		secret1 := "original-test-secret"
		created, _, err := sdk.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(hydra.OAuth2Client{
			ClientSecret:            new(secret1),
			GrantTypes:              []string{"client_credentials"},
			Scope:                   new("openid"),
			TokenEndpointAuthMethod: new("client_secret_basic"),
		}).Execute()
		require.NoError(t, err)
		clientID := *created.ClientId
		require.NotNil(t, created.ClientSecret)
		assert.Equal(t, secret1, *created.ClientSecret)

		// Verify original secret works
		token, err := newClientCredentialsConfig(clientID, secret1).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// First rotation
		rotated1, resp1, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, clientID).Execute()
		require.NoError(t, err)
		require.NotNil(t, rotated1.ClientSecret)
		secret2 := *rotated1.ClientSecret
		assert.NotEmpty(t, secret2)
		assert.NotEqual(t, secret1, secret2)
		assertNoRotatedSecretsInResponse(t, resp1)

		// Both secrets work after first rotation
		for _, secret := range []string{secret1, secret2} {
			token, err := newClientCredentialsConfig(clientID, secret).Token(ctx)
			require.NoError(t, err, "Secret %q should work after first rotation", secret)
			assert.NotEmpty(t, token.AccessToken)
		}

		// Second rotation
		rotated2, resp2, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, clientID).Execute()
		require.NoError(t, err)
		require.NotNil(t, rotated2.ClientSecret)
		secret3 := *rotated2.ClientSecret
		assert.NotEmpty(t, secret3)
		assert.NotEqual(t, secret2, secret3)
		assertNoRotatedSecretsInResponse(t, resp2)

		// All three secrets work after second rotation
		for _, secret := range []string{secret1, secret2, secret3} {
			token, err := newClientCredentialsConfig(clientID, secret).Token(ctx)
			require.NoError(t, err, "Secret %q should work after second rotation", secret)
			assert.NotEmpty(t, token.AccessToken)
		}

		// GetOAuth2Client should not return rotated_secrets
		fetched, resp, err := sdk.OAuth2API.GetOAuth2Client(ctx, clientID).Execute()
		require.NoError(t, err)
		assertNoRotatedSecretsInResponse(t, resp)
		assert.Nil(t, fetched.ClientSecret)

		// SetOAuth2Client should not return rotated_secrets
		fetched.ClientName = new("updated-name")
		set, resp, err := sdk.OAuth2API.SetOAuth2Client(ctx, clientID).OAuth2Client(*fetched).Execute()
		require.NoError(t, err)
		assert.Nil(t, set.ClientSecret)
		assertNoRotatedSecretsInResponse(t, resp)

		// PatchOAuth2Client should not return rotated_secrets
		patched, resp, err := sdk.OAuth2API.PatchOAuth2Client(ctx, clientID).JsonPatch([]hydra.JsonPatch{
			{Op: "replace", Path: "/client_name", Value: "patched-name"},
		}).Execute()
		require.NoError(t, err)
		assert.Nil(t, patched.ClientSecret)
		assertNoRotatedSecretsInResponse(t, resp)

		// ListOAuth2Clients should not return rotated_secrets
		list, resp, err := sdk.OAuth2API.ListOAuth2Clients(ctx).ClientName("patched-name").Execute()
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Nil(t, list[0].ClientSecret)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var raw []map[string]any
		require.NoError(t, json.Unmarshal(body, &raw))
		for _, item := range raw {
			assert.NotContains(t, item, "rotated")
		}

		// Rotated secrets still work after all API operations
		for _, secret := range []string{secret1, secret2, secret3} {
			_, err = newClientCredentialsConfig(clientID, secret).Token(ctx)
			require.NoError(t, err, "Secret %q should still work after API operations", secret)
		}

		// Delete rotated secrets
		fetched, resp, err = sdk.OAuth2API.DeleteRotatedOAuth2ClientSecrets(ctx, clientID).Execute()
		require.NoError(t, err)
		assertNoRotatedSecretsInResponse(t, resp)
		assert.Nil(t, fetched.ClientSecret)

		// Old secrets no longer work
		for _, secret := range []string{secret1, secret2} {
			_, err = newClientCredentialsConfig(clientID, secret).Token(ctx)
			require.Error(t, err, "Secret %q should fail after deleting rotated secrets", secret)
		}

		// Most recent secret still works after deleting rotated secrets
		token, err = newClientCredentialsConfig(clientID, secret3).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// Rotate once more
		rotated3, _, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, clientID).Execute()
		require.NoError(t, err)
		secret4 := *rotated3.ClientSecret
		assert.NotEmpty(t, secret4)
		assert.NotEqual(t, secret3, secret4)

		// New secret works
		token, err = newClientCredentialsConfig(clientID, secret4).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// Explicitly setting a new secret revokes all previous secrets, including rotated secrets
		secret5 := "second-explicit-test-secret"
		patched, _, err = sdk.OAuth2API.PatchOAuth2Client(ctx, clientID).JsonPatch([]hydra.JsonPatch{
			{Op: "replace", Path: "/client_secret", Value: secret5},
		}).Execute()
		require.NoError(t, err)
		require.NotNil(t, patched.ClientSecret)
		assert.Equal(t, secret5, *patched.ClientSecret)

		// New secret works
		token, err = newClientCredentialsConfig(clientID, secret5).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// Previous secrets no longer work
		for i, secret := range []string{secret1, secret2, secret3, secret4} {
			_, err = newClientCredentialsConfig(clientID, secret).Token(ctx)
			assert.Error(t, err, "secret%d should fail after setting new secret", i+1)
		}

		// Rotate one last time
		rotated6, _, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, clientID).Execute()
		require.NoError(t, err)
		secret6 := *rotated6.ClientSecret
		assert.NotEmpty(t, secret6)
		assert.NotEqual(t, secret5, secret6)

		// New secret works
		token, err = newClientCredentialsConfig(clientID, secret6).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// Previous secret still works after rotation
		token, err = newClientCredentialsConfig(clientID, secret5).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// Explicitly setting a new secret again revokes all previous secrets, including rotated secrets
		secret7 := "third-explicit-test-secret"
		client7 := *rotated6
		client7.ClientSecret = &secret7
		set, _, err = sdk.OAuth2API.SetOAuth2Client(ctx, clientID).OAuth2Client(client7).Execute()
		require.NoError(t, err)
		require.NotNil(t, set.ClientSecret)
		assert.Equal(t, secret7, *set.ClientSecret)

		// New secret works
		token, err = newClientCredentialsConfig(clientID, secret7).Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)

		// Previous secrets no longer work
		for _, secret := range []string{secret1, secret2, secret3, secret4, secret5, secret6} {
			_, err = newClientCredentialsConfig(clientID, secret).Token(ctx)
			assert.Error(t, err, "Secret %q should fail after setting new secret", secret)
		}
	})

	t.Run("case=rotate secret for nonexistent client fails", func(t *testing.T) {
		_, _, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, "nonexistent-client-id").Execute()
		require.Error(t, err)
	})

	t.Run("case=delete rotated secrets for nonexistent client fails", func(t *testing.T) {
		_, _, err := sdk.OAuth2API.DeleteRotatedOAuth2ClientSecrets(ctx, "nonexistent-client-id").Execute()
		require.Error(t, err)
	})

	t.Run("case=public clients cannot rotate secrets", func(t *testing.T) {
		created, _, err := sdk.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(hydra.OAuth2Client{
			GrantTypes:              []string{"authorization_code"},
			Scope:                   new("offline"),
			TokenEndpointAuthMethod: new("none"),
		}).Execute()
		require.NoError(t, err)
		require.NotNil(t, created.ClientId)
		require.Nil(t, created.ClientSecret)
		_, res, err := sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, *created.ClientId).Execute()
		require.Error(t, err)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "does not support client secret rotation")
		_, res, err = sdk.OAuth2API.DeleteRotatedOAuth2ClientSecrets(ctx, *created.ClientId).Execute()
		require.Error(t, err)
		body, err = io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "does not support client secret rotation")
	})

	t.Run("case=clients with private_key_jwt cannot rotate secrets", func(t *testing.T) {
		created, res, err := sdk.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(hydra.OAuth2Client{
			GrantTypes:              []string{"authorization_code"},
			Scope:                   new("offline"),
			TokenEndpointAuthMethod: new("private_key_jwt"),
			Jwks: &hydra.JsonWebKeySet{
				Keys: []hydra.JsonWebKey{{
					Kid: "test-key-id",
					Kty: "RSA",
					Use: "sig",
					Alg: "RS256",
					N:   new("l80jJJqcc1PpefIGVIjuPvA1D7NscnuF9aQqLa7I9rDUK4IaSOO3kL_EF13k-jTzcA5q4OZn5dR0kmqIMZT2gQ"),
					E:   new("AQAB"),
				}},
			},
		}).Execute()
		require.NoError(t, err, "Failed to create client: %s", ioutilx.MustReadAll(res.Body))
		require.NotNil(t, created.ClientId)
		_, res, err = sdk.OAuth2API.RotateOAuth2ClientSecret(ctx, *created.ClientId).Execute()
		require.Error(t, err)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "does not support client secret rotation")
		_, res, err = sdk.OAuth2API.DeleteRotatedOAuth2ClientSecrets(ctx, *created.ClientId).Execute()
		require.Error(t, err)
		body, err = io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "does not support client secret rotation")
	})
}
