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

	"github.com/ory/hydra/jwk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
)

func TestHandlerWellKnown(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)
	var testGenerator = &jwk.RS256Generator{}
	conf.MustSet(config.KeyWellKnownKeys, []string{x.OpenIDConnectKeyName, x.OpenIDConnectKeyName})
	router := x.NewRouterPublic()
	h := reg.KeyHandler()
	h.SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
		return h
	})
	testServer := httptest.NewServer(router)
	JWKPath := "/.well-known/jwks.json"

	t.Run("Test_Handler_WellKnown/Run_public_key_With_public_prefix", func(t *testing.T) {
		if conf.HsmEnabled() {
			t.Skip("Skipping test. Not applicable when Hardware Security Module is enabled. Public/private keys on HSM are generated with equal key id's and are not using prefixes")
		}
		IDKS, _ := testGenerator.Generate("test-id-1", "sig")
		require.NoError(t, reg.KeyManager().AddKeySet(context.TODO(), x.OpenIDConnectKeyName, IDKS))
		res, err := http.Get(testServer.URL + JWKPath)
		require.NoError(t, err, "problem in http request")
		defer res.Body.Close()

		var known jose.JSONWebKeySet
		err = json.NewDecoder(res.Body).Decode(&known)
		require.NoError(t, err, "problem in decoding response")

		require.Len(t, known.Keys, 1)

		knownKey := known.Key("public:test-id-1")[0]
		require.NotNil(t, knownKey, "Could not find key public")

		expectedKey, err := jwk.FindPublicKey(IDKS)
		require.NoError(t, err)
		assert.EqualValues(t, canonicalizeThumbprints(*expectedKey), canonicalizeThumbprints(knownKey))
		require.NoError(t, reg.KeyManager().DeleteKeySet(context.TODO(), x.OpenIDConnectKeyName))
	})

	t.Run("Test_Handler_WellKnown/Run_public_key_Without_public_prefix", func(t *testing.T) {
		var IDKS *jose.JSONWebKeySet

		if conf.HsmEnabled() {
			IDKS, _ = reg.KeyManager().GenerateAndPersistKeySet(context.TODO(), x.OpenIDConnectKeyName, "test-id-2", "RS256", "sig")
		} else {
			IDKS, _ = testGenerator.Generate("test-id-2", "sig")
			if strings.ContainsAny(IDKS.Keys[1].KeyID, "public") {
				IDKS.Keys[1].KeyID = "test-id-2"
			} else {
				IDKS.Keys[0].KeyID = "test-id-2"
			}
			require.NoError(t, reg.KeyManager().AddKeySet(context.TODO(), x.OpenIDConnectKeyName, IDKS))
		}

		res, err := http.Get(testServer.URL + JWKPath)
		require.NoError(t, err, "problem in http request")
		defer res.Body.Close()

		var known jose.JSONWebKeySet
		err = json.NewDecoder(res.Body).Decode(&known)
		require.NoError(t, err, "problem in decoding response")
		require.Len(t, known.Keys, 1)

		knownKey := known.Key("test-id-2")[0]
		require.NotNil(t, knownKey, "Could not find key public")

		expectedKey, err := jwk.FindPublicKey(IDKS)
		require.NoError(t, err)
		assert.EqualValues(t, canonicalizeThumbprints(*expectedKey), canonicalizeThumbprints(knownKey))
	})
}

func canonicalizeThumbprints(js jose.JSONWebKey) jose.JSONWebKey {
	if len(js.CertificateThumbprintSHA1) == 0 {
		js.CertificateThumbprintSHA1 = nil
	}
	if len(js.CertificateThumbprintSHA256) == 0 {
		js.CertificateThumbprintSHA256 = nil
	}
	return js
}
