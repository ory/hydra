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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/x"

	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/jwk"
)

func TestJWKSDK(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)

	router := x.NewRouterAdmin()
	h := NewHandler(reg, conf)
	h.SetRoutes(router, x.NewRouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)
	sdk := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(server.URL).Host})

	expectedPublicKid := "public:key-bar"
	if conf.HsmEnabled() {
		expectedPublicKid = "key-bar"
	}

	t.Run("JSON Web Key", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			// Create a key called set-foo
			resultKeys, err := sdk.Admin.CreateJSONWebKeySet(admin.NewCreateJSONWebKeySetParams().WithSet("set-foo").WithBody(&models.JSONWebKeySetGeneratorRequest{
				Alg: pointerx.String("RS256"),
				Kid: pointerx.String("key-bar"),
				Use: pointerx.String("sig"),
			}))
			require.NoError(t, err)
			if conf.HsmEnabled() {
				require.Len(t, resultKeys.Payload.Keys, 1)
				assert.Equal(t, expectedPublicKid, *resultKeys.Payload.Keys[0].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[0].Alg)
				assert.Equal(t, "sig", *resultKeys.Payload.Keys[0].Use)
			} else {
				require.Len(t, resultKeys.Payload.Keys, 2)
				assert.Equal(t, "private:key-bar", *resultKeys.Payload.Keys[0].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[0].Alg)
				assert.Equal(t, "sig", *resultKeys.Payload.Keys[0].Use)
				assert.Equal(t, expectedPublicKid, *resultKeys.Payload.Keys[1].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[1].Alg)
				assert.Equal(t, "sig", *resultKeys.Payload.Keys[1].Use)
			}
		})

		var resultKeys *models.JSONWebKeySet
		t.Run("GetJwkSetKey after create", func(t *testing.T) {
			result, err := sdk.Admin.GetJSONWebKey(admin.NewGetJSONWebKeyParams().WithKid(expectedPublicKid).WithSet("set-foo"))
			require.NoError(t, err)
			require.Len(t, result.Payload.Keys, 1)
			require.Equal(t, expectedPublicKid, *result.Payload.Keys[0].Kid)
			require.Equal(t, "RS256", *result.Payload.Keys[0].Alg)

			resultKeys = result.Payload
		})

		t.Run("UpdateJwkSetKey", func(t *testing.T) {
			if conf.HsmEnabled() {
				t.Skip("Skipping test. Keys cannot be updated when Hardware Security Module is enabled")
			}
			require.Len(t, resultKeys.Keys, 1)
			resultKeys.Keys[0].Alg = pointerx.String("ES256")

			resultKey, err := sdk.Admin.UpdateJSONWebKey(admin.NewUpdateJSONWebKeyParams().WithKid(expectedPublicKid).WithSet("set-foo").WithBody(resultKeys.Keys[0]))
			require.NoError(t, err)
			assert.Equal(t, expectedPublicKid, *resultKey.Payload.Kid)
			assert.Equal(t, "ES256", *resultKey.Payload.Alg)
		})

		t.Run("DeleteJwkSetKey after delete", func(t *testing.T) {
			_, err := sdk.Admin.DeleteJSONWebKey(admin.NewDeleteJSONWebKeyParams().WithKid(expectedPublicKid).WithSet("set-foo"))
			require.NoError(t, err)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, err := sdk.Admin.GetJSONWebKey(admin.NewGetJSONWebKeyParams().WithKid(expectedPublicKid).WithSet("set-foo"))
			require.Error(t, err)
		})

	})

	t.Run("JWK Set", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			resultKeys, err := sdk.Admin.CreateJSONWebKeySet(admin.NewCreateJSONWebKeySetParams().WithSet("set-foo2").WithBody(&models.JSONWebKeySetGeneratorRequest{
				Alg: pointerx.String("RS256"),
				Kid: pointerx.String("key-bar"),
			}))
			require.NoError(t, err)
			if conf.HsmEnabled() {
				require.Len(t, resultKeys.Payload.Keys, 1)
				assert.Equal(t, expectedPublicKid, *resultKeys.Payload.Keys[0].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[0].Alg)
			} else {
				require.Len(t, resultKeys.Payload.Keys, 2)
				assert.Equal(t, "private:key-bar", *resultKeys.Payload.Keys[0].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[0].Alg)
				assert.Equal(t, expectedPublicKid, *resultKeys.Payload.Keys[1].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[1].Alg)
			}
		})

		resultKeys, err := sdk.Admin.GetJSONWebKeySet(admin.NewGetJSONWebKeySetParams().WithSet("set-foo2"))
		t.Run("GetJwkSet after create", func(t *testing.T) {
			require.NoError(t, err)
			if conf.HsmEnabled() {
				require.Len(t, resultKeys.Payload.Keys, 1)
				assert.Equal(t, expectedPublicKid, *resultKeys.Payload.Keys[0].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[0].Alg)
			} else {
				require.Len(t, resultKeys.Payload.Keys, 2)
				assert.Equal(t, expectedPublicKid, *resultKeys.Payload.Keys[0].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[0].Alg)
				assert.Equal(t, "private:key-bar", *resultKeys.Payload.Keys[1].Kid)
				assert.Equal(t, "RS256", *resultKeys.Payload.Keys[1].Alg)
			}
		})

		t.Run("UpdateJwkSet", func(t *testing.T) {
			if conf.HsmEnabled() {
				t.Skip("Skipping test. Keys cannot be updated when Hardware Security Module is enabled")
			}
			require.Len(t, resultKeys.Payload.Keys, 2)
			resultKeys.Payload.Keys[0].Alg = pointerx.String("ES256")
			resultKeys.Payload.Keys[1].Alg = pointerx.String("ES256")

			result, err := sdk.Admin.UpdateJSONWebKeySet(admin.NewUpdateJSONWebKeySetParams().WithSet("set-foo2").WithBody(resultKeys.Payload))
			require.NoError(t, err)
			require.Len(t, result.Payload.Keys, 2)
			assert.Equal(t, expectedPublicKid, *result.Payload.Keys[0].Kid)
			assert.Equal(t, "ES256", *result.Payload.Keys[0].Alg)
			assert.Equal(t, "private:key-bar", *result.Payload.Keys[1].Kid)
			assert.Equal(t, "ES256", *result.Payload.Keys[1].Alg)
		})

		t.Run("DeleteJwkSet", func(t *testing.T) {
			_, err := sdk.Admin.DeleteJSONWebKeySet(admin.NewDeleteJSONWebKeySetParams().WithSet("set-foo2"))
			require.NoError(t, err)
		})

		t.Run("GetJwkSet after delete", func(t *testing.T) {
			_, err := sdk.Admin.GetJSONWebKeySet(admin.NewGetJSONWebKeySetParams().WithSet("set-foo2"))
			require.Error(t, err)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, err := sdk.Admin.GetJSONWebKey(admin.NewGetJSONWebKeyParams().WithSet("set-foo2").WithKid(expectedPublicKid))
			require.Error(t, err)
		})
	})
}
