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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/viper"

	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
)

func TestHandlerWellKnown(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	viper.Set(configuration.ViperKeyWellKnownKeys, []string{x.OpenIDConnectKeyName, x.OpenIDConnectKeyName})

	router := x.NewRouterPublic()
	IDKS, _ := testGenerator.Generate("test-id", "sig")

	h := reg.KeyHandler()
	require.NoError(t, reg.KeyManager().AddKeySet(context.TODO(), x.OpenIDConnectKeyName, IDKS))

	h.SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
		return h
	})
	testServer := httptest.NewServer(router)

	JWKPath := "/.well-known/jwks.json"
	res, err := http.Get(testServer.URL + JWKPath)
	require.NoError(t, err, "problem in http request")
	defer res.Body.Close()

	var known jose.JSONWebKeySet
	err = json.NewDecoder(res.Body).Decode(&known)
	require.NoError(t, err, "problem in decoding response")

	require.Len(t, known.Keys, 1)

	resp := known.Key("public:test-id")
	require.NotNil(t, resp, "Could not find key public")
	assert.Equal(t, resp, IDKS.Key("public:test-id"))
}

func TestHandlerKeySet(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	viper.Set(configuration.ViperKeyWellKnownKeys, []string{x.OpenIDConnectKeyName, x.OpenIDConnectKeyName})

	router := x.NewRouterAdmin()

	h := reg.KeyHandler()

	h.SetRoutes(router, router.RouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	testServer := httptest.NewServer(router)

	t.Run("CreateJSONWebKeySet", func(t *testing.T) {
		createJWKSetPath := "/keys/test-key"
		createJWKSetReqBody := strings.NewReader(`{"alg": "RS256", "kid": "test-01", "use": "sig"}`)
		createRes, err := http.Post(testServer.URL+createJWKSetPath, "application/json", createJWKSetReqBody)
		require.NoError(t, err, "problem in http request")
		defer createRes.Body.Close()

		var createdJWKSet jose.JSONWebKeySet
		err = json.NewDecoder(createRes.Body).Decode(&createdJWKSet)
		require.NoError(t, err, "problem in decoding response")

		assert.EqualValues(t, createRes.StatusCode, http.StatusCreated)
		assert.Len(t, createdJWKSet.Key("public:test-01"), 1)
		assert.Len(t, createdJWKSet.Key("private:test-01"), 1)
	})

	t.Run("GetJSONWebKeySet", func(t *testing.T) {
		getJWKSetPath := "/keys/test-key"
		getRes, err := http.Get(testServer.URL + getJWKSetPath)
		require.NoError(t, err, "problem in http request")
		defer getRes.Body.Close()

		var getJWKSet jose.JSONWebKeySet
		err = json.NewDecoder(getRes.Body).Decode(&getJWKSet)
		require.NoError(t, err, "problem in decoding response")

		require.Len(t, getJWKSet.Keys, 2)
	})

	t.Run("DeleteJSONWebKeySet", func(t *testing.T) {
		deleteJWKSetPath := "/keys/test-key"
		deleteReq, err := http.NewRequest(http.MethodDelete, testServer.URL+deleteJWKSetPath, nil)
		deleteRes, err := http.DefaultClient.Do(deleteReq)
		deleteReq.Header.Add("Content-Type", "application/json")
		require.NoError(t, err, "problem in http request")
		defer deleteRes.Body.Close()

		assert.Equal(t, http.StatusNoContent, deleteRes.StatusCode)

		getJWKSetPath := "/keys/test-key"
		getRes, err := http.Get(testServer.URL + getJWKSetPath)
		require.NoError(t, err, "problem in http request")
		defer getRes.Body.Close()

		assert.Equal(t, http.StatusNotFound, getRes.StatusCode)
	})
}
