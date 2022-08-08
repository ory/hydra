/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/hydra/driver/config"

	hydra "github.com/ory/hydra-client-go"

	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/jwk"
)

func TestJWKSDK(t *testing.T) {
	ctx := context.Background()
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

	router := x.NewRouterAdmin(conf.AdminURL)
	h := NewHandler(reg)
	h.SetRoutes(router, x.NewRouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)
	conf.MustSet(ctx, config.KeyAdminURL, server.URL)

	sdk := hydra.NewAPIClient(hydra.NewConfiguration())
	sdk.GetConfig().Servers = hydra.ServerConfigurations{{URL: server.URL}}

	expectedKid := "key-bar"
	t.Run("JSON Web Key", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			// Create a key called set-foo
			resultKeys, _, err := sdk.V0alpha2Api.AdminCreateJsonWebKeySet(context.Background(), "set-foo").AdminCreateJsonWebKeySetBody(hydra.AdminCreateJsonWebKeySetBody{
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
			result, _, err := sdk.V0alpha2Api.AdminGetJsonWebKey(ctx, "set-foo", expectedKid).Execute()
			require.NoError(t, err)
			require.Len(t, result.Keys, 1)
			require.Equal(t, expectedKid, result.Keys[0].Kid)
			require.Equal(t, "RS256", result.Keys[0].Alg)

			resultKeys = result
		})

		t.Run("UpdateJwkSetKey", func(t *testing.T) {
			if conf.HSMEnabled() {
				t.Skip("Skipping test. Keys cannot be updated when Hardware Security Module is enabled")
			}
			require.Len(t, resultKeys.Keys, 1)
			resultKeys.Keys[0].Alg = "ES256"

			resultKey, _, err := sdk.V0alpha2Api.AdminUpdateJsonWebKey(ctx, "set-foo", expectedKid).JsonWebKey(resultKeys.Keys[0]).Execute()
			require.NoError(t, err)
			assert.Equal(t, expectedKid, resultKey.Kid)
			assert.Equal(t, "ES256", resultKey.Alg)
		})

		t.Run("DeleteJwkSetKey after delete", func(t *testing.T) {
			_, err := sdk.V0alpha2Api.AdminDeleteJsonWebKey(ctx, "set-foo", expectedKid).Execute()
			require.NoError(t, err)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, res, err := sdk.V0alpha2Api.AdminGetJsonWebKey(ctx, "set-foo", expectedKid).Execute()
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})

	})

	t.Run("JWK Set", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			resultKeys, _, err := sdk.V0alpha2Api.AdminCreateJsonWebKeySet(ctx, "set-foo2").AdminCreateJsonWebKeySetBody(hydra.AdminCreateJsonWebKeySetBody{
				Alg: "RS256",
				Kid: "key-bar",
			}).Execute()
			require.NoError(t, err)
			require.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, expectedKid, resultKeys.Keys[0].Kid)
			assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
		})

		resultKeys, _, err := sdk.V0alpha2Api.AdminGetJsonWebKeySet(ctx, "set-foo2").Execute()
		t.Run("GetJwkSet after create", func(t *testing.T) {
			require.NoError(t, err)
			if conf.HSMEnabled() {
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
			if conf.HSMEnabled() {
				t.Skip("Skipping test. Keys cannot be updated when Hardware Security Module is enabled")
			}
			require.Len(t, resultKeys.Keys, 1)
			resultKeys.Keys[0].Alg = "ES256"

			result, _, err := sdk.V0alpha2Api.AdminUpdateJsonWebKeySet(ctx, "set-foo2").JsonWebKeySet(*resultKeys).Execute()
			require.NoError(t, err)
			require.Len(t, result.Keys, 1)
			assert.Equal(t, expectedKid, result.Keys[0].Kid)
			assert.Equal(t, "ES256", result.Keys[0].Alg)
		})

		t.Run("DeleteJwkSet", func(t *testing.T) {
			_, err := sdk.V0alpha2Api.AdminDeleteJsonWebKeySet(ctx, "set-foo2").Execute()
			require.NoError(t, err)
		})

		t.Run("GetJwkSet after delete", func(t *testing.T) {
			_, res, err := sdk.V0alpha2Api.AdminGetJsonWebKeySet(ctx, "set-foo2").Execute()
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, res, err := sdk.V0alpha2Api.AdminGetJsonWebKey(ctx, "set-foo2", expectedKid).Execute()
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	})
}
