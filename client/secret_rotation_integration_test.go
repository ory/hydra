// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/internal/testhelpers"
)

func TestSecretRotation_E2E(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t)
	public, _ := testhelpers.NewOAuth2Server(ctx, t, reg)

	t.Run("case=complete_secret_rotation_workflow", func(t *testing.T) {
		// Step 1: Create a client with a known secret
		originalSecret := uuid.Must(uuid.NewV4()).String()
		c := &client.Client{
			Secret:        originalSecret,
			RedirectURIs:  []string{public.URL + "/callback"},
			ResponseTypes: []string{"token"},
			GrantTypes:    []string{"client_credentials"},
			Scope:         "foobar",
		}
		require.NoError(t, reg.ClientManager().CreateClient(ctx, c))
		clientID := c.GetID()

		conf := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: originalSecret,
			TokenURL:     reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:    goauth2.AuthStyleInHeader,
		}

		// Step 2: Verify original secret works for authentication
		t.Run("step=authenticate_with_original_secret", func(t *testing.T) {
			token, err := conf.Token(ctx)
			require.NoError(t, err)
			assert.NotEmpty(t, token.AccessToken)
		})

		// Step 3: Rotate the secret
		newSecret := uuid.Must(uuid.NewV4()).String()
		t.Run("step=rotate_secret", func(t *testing.T) {
			c, err := reg.ClientManager().GetConcreteClient(ctx, clientID)
			require.NoError(t, err)

			c.Secret = newSecret
			err = reg.ClientManager().UpdateClient(ctx, c)
			require.NoError(t, err)
		})

		// Step 4: Verify NEW secret works
		t.Run("step=authenticate_with_new_secret", func(t *testing.T) {
			newConf := clientcredentials.Config{
				ClientID:     clientID,
				ClientSecret: newSecret,
				TokenURL:     reg.Config().OAuth2TokenURL(ctx).String(),
				AuthStyle:    goauth2.AuthStyleInHeader,
			}
			token, err := newConf.Token(ctx)
			require.NoError(t, err)
			assert.NotEmpty(t, token.AccessToken)
		})

		// Step 5: Verify OLD secret STILL works (key test for rotation)
		t.Run("step=authenticate_with_old_secret_after_rotation", func(t *testing.T) {
			token, err := conf.Token(ctx)
			require.NoError(t, err, "Old secret should still work after rotation")
			assert.NotEmpty(t, token.AccessToken)
		})

		// Step 6: Verify rotated secrets are stored
		t.Run("step=verify_rotated_secrets_stored", func(t *testing.T) {
			c, err := reg.ClientManager().GetConcreteClient(ctx, clientID)
			require.NoError(t, err)

			assert.NotEmpty(t, c.RotatedSecrets, "Rotated secrets should be stored")

			var rotated []string
			err = json.Unmarshal([]byte(c.RotatedSecrets), &rotated)
			require.NoError(t, err)
			assert.Len(t, rotated, 1, "Should have one rotated secret")
		})

		// Step 7: Clear rotated secrets
		t.Run("step=clear_rotated_secrets", func(t *testing.T) {
			c, err := reg.ClientManager().GetConcreteClient(ctx, clientID)
			require.NoError(t, err)

			c.RotatedSecrets = "[]"
			c.Secret = ""
			err = reg.ClientManager().UpdateClient(ctx, c)
			require.NoError(t, err)

			c, err = reg.ClientManager().GetConcreteClient(ctx, clientID)
			require.NoError(t, err)
			assert.Equal(t, "[]", c.RotatedSecrets)
		})

		// Step 8: Verify OLD secret no longer works after cleanup
		t.Run("step=old_secret_fails_after_cleanup", func(t *testing.T) {
			_, err := conf.Token(ctx)
			require.Error(t, err, "Old secret should not work after cleanup")
		})

		// Step 9: Verify NEW secret still works after cleanup
		t.Run("step=new_secret_works_after_cleanup", func(t *testing.T) {
			newConf := clientcredentials.Config{
				ClientID:     clientID,
				ClientSecret: newSecret,
				TokenURL:     reg.Config().OAuth2TokenURL(ctx).String(),
				AuthStyle:    goauth2.AuthStyleInHeader,
			}
			token, err := newConf.Token(ctx)
			require.NoError(t, err)
			assert.NotEmpty(t, token.AccessToken)
		})
	})

	t.Run("case=multiple_rotations", func(t *testing.T) {
		// Create client with first secret
		secret1 := uuid.Must(uuid.NewV4()).String()
		c := &client.Client{
			Secret:        secret1,
			RedirectURIs:  []string{public.URL + "/callback"},
			ResponseTypes: []string{"token"},
			GrantTypes:    []string{"client_credentials"},
			Scope:         "foobar",
		}
		require.NoError(t, reg.ClientManager().CreateClient(ctx, c))
		clientID := c.GetID()

		// Rotate to secret2
		secret2 := uuid.Must(uuid.NewV4()).String()
		c, err := reg.ClientManager().GetConcreteClient(ctx, clientID)
		require.NoError(t, err)
		c.Secret = secret2
		err = reg.ClientManager().UpdateClient(ctx, c)
		require.NoError(t, err)

		// Rotate to secret3
		secret3 := uuid.Must(uuid.NewV4()).String()
		c, err = reg.ClientManager().GetConcreteClient(ctx, clientID)
		require.NoError(t, err)
		c.Secret = secret3
		err = reg.ClientManager().UpdateClient(ctx, c)
		require.NoError(t, err)

		// Verify rotated secrets array has 2 entries
		c, err = reg.ClientManager().GetConcreteClient(ctx, clientID)
		require.NoError(t, err)

		var rotated []string
		err = json.Unmarshal([]byte(c.RotatedSecrets), &rotated)
		require.NoError(t, err)
		assert.Len(t, rotated, 2, "Should have two rotated secrets")

		// Test all three secrets work
		conf1 := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: secret1,
			TokenURL:     reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:    goauth2.AuthStyleInHeader,
		}
		token1, err := conf1.Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token1.AccessToken)

		conf2 := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: secret2,
			TokenURL:     reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:    goauth2.AuthStyleInHeader,
		}
		token2, err := conf2.Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token2.AccessToken)

		conf3 := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: secret3,
			TokenURL:     reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:    goauth2.AuthStyleInHeader,
		}
		token3, err := conf3.Token(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, token3.AccessToken)
	})
}
