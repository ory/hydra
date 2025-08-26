// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	. "github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/prometheusx"
)

func TestJWKSDK(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t)

	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	router := x.NewRouterAdmin(metrics)
	h := NewHandler(reg)
	h.SetAdminRoutes(router)
	server := httptest.NewServer(router)
	reg.Config().MustSet(ctx, config.KeyAdminURL, server.URL)

	sdk := hydra.NewAPIClient(hydra.NewConfiguration())
	sdk.GetConfig().Servers = hydra.ServerConfigurations{{URL: server.URL}}

	expectedKid := "key-bar"
	t.Run("JSON Web Key", func(t *testing.T) {
		t.Parallel()
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			// Create a key called set-foo
			resultKeys, _, err := sdk.JwkAPI.CreateJsonWebKeySet(context.Background(), "set-foo").CreateJsonWebKeySet(hydra.CreateJsonWebKeySet{
				Alg: "RS256",
				Kid: "key-bar",
				Use: "sig",
			}).Execute()
			require.NoError(t, err)
			require.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
			assert.Equal(t, "sig", resultKeys.Keys[0].Use)
		})

		var resultKeys *hydra.JsonWebKeySet
		t.Run("GetJwkSetKey after create", func(t *testing.T) {
			result, _, err := sdk.JwkAPI.GetJsonWebKey(ctx, "set-foo", expectedKid).Execute()
			require.NoError(t, err)
			require.Len(t, result.Keys, 1)
			require.Equal(t, expectedKid, result.Keys[0].Kid)
			require.Equal(t, "RS256", result.Keys[0].Alg)

			resultKeys = result
		})

		t.Run("UpdateJwkSetKey", func(t *testing.T) {
			if reg.Config().HSMEnabled() {
				t.Skip("Skipping test. Keys cannot be updated when Hardware Security Module is enabled")
			}
			require.Len(t, resultKeys.Keys, 1)
			resultKeys.Keys[0].Alg = "ES256"

			resultKey, _, err := sdk.JwkAPI.SetJsonWebKey(ctx, "set-foo", expectedKid).JsonWebKey(resultKeys.Keys[0]).Execute()
			require.NoError(t, err)
			assert.Equal(t, expectedKid, resultKey.Kid)
			assert.Equal(t, "ES256", resultKey.Alg)
		})

		t.Run("DeleteJwkSetKey after delete", func(t *testing.T) {
			_, err := sdk.JwkAPI.DeleteJsonWebKey(ctx, "set-foo", expectedKid).Execute()
			require.NoError(t, err)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, res, err := sdk.JwkAPI.GetJsonWebKey(ctx, "set-foo", expectedKid).Execute()
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})

	})

	t.Run("JWK Set", func(t *testing.T) {
		t.Parallel()
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			resultKeys, _, err := sdk.JwkAPI.CreateJsonWebKeySet(ctx, "set-foo2").CreateJsonWebKeySet(hydra.CreateJsonWebKeySet{
				Alg: "RS256",
				Kid: "key-bar",
				Use: "sig",
			}).Execute()
			require.NoError(t, err)
			require.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, expectedKid, resultKeys.Keys[0].Kid)
			assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
		})

		resultKeys, _, err := sdk.JwkAPI.GetJsonWebKeySet(ctx, "set-foo2").Execute()
		t.Run("GetJwkSet after create", func(t *testing.T) {
			require.NoError(t, err)
			if reg.Config().HSMEnabled() {
				require.Len(t, resultKeys.Keys, 1)
				assert.Equal(t, expectedKid, resultKeys.Keys[0].Kid)
				assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
			} else {
				require.Len(t, resultKeys.Keys, 1)
				assert.Equal(t, expectedKid, resultKeys.Keys[0].Kid)
				assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
			}
		})

		t.Run("UpdateJwkSet", func(t *testing.T) {
			if reg.Config().HSMEnabled() {
				t.Skip("Skipping test. Keys cannot be updated when Hardware Security Module is enabled")
			}
			require.Len(t, resultKeys.Keys, 1)
			resultKeys.Keys[0].Alg = "ES256"

			result, _, err := sdk.JwkAPI.SetJsonWebKeySet(ctx, "set-foo2").JsonWebKeySet(*resultKeys).Execute()
			require.NoError(t, err)
			require.Len(t, result.Keys, 1)
			assert.Equal(t, expectedKid, result.Keys[0].Kid)
			assert.Equal(t, "ES256", result.Keys[0].Alg)
		})

		t.Run("DeleteJwkSet", func(t *testing.T) {
			_, err := sdk.JwkAPI.DeleteJsonWebKeySet(ctx, "set-foo2").Execute()
			require.NoError(t, err)
		})

		t.Run("GetJwkSet after delete", func(t *testing.T) {
			_, res, err := sdk.JwkAPI.GetJsonWebKeySet(ctx, "set-foo2").Execute()
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, res, err := sdk.JwkAPI.GetJsonWebKey(ctx, "set-foo2", expectedKid).Execute()
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	})
}
