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

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/herodot"
	. "github.com/ory/hydra/jwk"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

func TestJWKSDK(t *testing.T) {
	manager := new(MemoryManager)

	router := httprouter.New()
	h := Handler{
		Manager: manager,
		H:       herodot.NewJSONWriter(nil),
	}
	h.SetRoutes(router, router, func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)

	client := hydra.NewAdminApiWithBasePath(server.URL)

	t.Run("JSON Web Key", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			// Create a key called set-foo
			resultKeys, response, err := client.CreateJsonWebKeySet("set-foo", hydra.JsonWebKeySetGeneratorRequest{
				Alg: "HS256",
				Kid: "key-bar",
				Use: "sig",
			})
			require.NoError(t, err)
			require.EqualValues(t, http.StatusCreated, response.StatusCode)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
			assert.Equal(t, "sig", resultKeys.Keys[0].Use)
		})

		resultKeys, response, err := client.GetJsonWebKey("key-bar", "set-foo")
		t.Run("GetJwkSetKey after create", func(t *testing.T) {
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, response.StatusCode)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		t.Run("UpdateJwkSetKey", func(t *testing.T) {
			resultKeys.Keys[0].Alg = "RS256"
			resultKey, response, err := client.UpdateJsonWebKey("key-bar", "set-foo", resultKeys.Keys[0])
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, "key-bar", resultKey.Kid)
			assert.Equal(t, "RS256", resultKey.Alg)
		})

		t.Run("DeleteJwkSetKey after delete", func(t *testing.T) {
			response, err := client.DeleteJsonWebKey("key-bar", "set-foo")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, response, err := client.GetJsonWebKey("key-bar", "set-foo")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
		})

	})

	t.Run("JWK Set", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			resultKeys, _, err := client.CreateJsonWebKeySet("set-foo2", hydra.JsonWebKeySetGeneratorRequest{
				Alg: "HS256",
				Kid: "key-bar",
			})
			require.NoError(t, err)

			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		resultKeys, _, err := client.GetJsonWebKeySet("set-foo2")
		t.Run("GetJwkSet after create", func(t *testing.T) {
			require.NoError(t, err)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		t.Run("UpdateJwkSet", func(t *testing.T) {
			resultKeys.Keys[0].Alg = "RS256"
			resultKeys, _, err = client.UpdateJsonWebKeySet("set-foo2", *resultKeys)
			require.NoError(t, err)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
		})

		t.Run("DeleteJwkSet", func(t *testing.T) {
			response, err := client.DeleteJsonWebKeySet("set-foo2")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.Run("GetJwkSet after delete", func(t *testing.T) {
			_, response, err := client.GetJsonWebKeySet("set-foo2")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, response, err := client.GetJsonWebKey("key-bar", "set-foo2")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
		})
	})
}
