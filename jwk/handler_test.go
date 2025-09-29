// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/urlx"
)

func TestHandlerWellKnown(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue(config.KeyWellKnownKeys, []string{x.OpenIDConnectKeyName, x.OpenIDConnectKeyName})))
	router := x.NewRouterPublic(prometheusx.NewMetricsManager("", "", "", ""))
	h := jwk.NewHandler(reg)
	h.SetPublicRoutes(router, func(h http.Handler) http.Handler {
		return h
	})
	testServer := httptest.NewServer(router)
	JWKPath := "/.well-known/jwks.json"

	t.Run("Test_Handler_WellKnown/Run_public_key_With_public_prefix", func(t *testing.T) {
		t.Parallel()
		if reg.Config().HSMEnabled() {
			t.Skip("Skipping test. Not applicable when Hardware Security Module is enabled. Public/private keys on HSM are generated with equal key id's and are not using prefixes")
		}
		IDKS, _ := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
		require.NoError(t, reg.KeyManager().AddKeySet(context.TODO(), x.OpenIDConnectKeyName, IDKS))
		res, err := http.Get(urlx.MustJoin(testServer.URL, JWKPath))
		require.NoError(t, err, "problem in http request")
		defer res.Body.Close() //nolint:errcheck

		var known jose.JSONWebKeySet
		err = json.NewDecoder(res.Body).Decode(&known)
		require.NoError(t, err, "problem in decoding response")

		require.GreaterOrEqual(t, len(known.Keys), 1)

		knownKey := known.Key("test-id-1")[0].Public()
		require.NotNil(t, knownKey, "Could not find key public")

		expectedKey, err := jwk.FindPublicKey(IDKS)
		require.NoError(t, err)
		assert.EqualValues(t, canonicalizeThumbprints(*expectedKey), canonicalizeThumbprints(knownKey))
		require.NoError(t, reg.KeyManager().DeleteKeySet(context.TODO(), x.OpenIDConnectKeyName))
	})

	t.Run("Test_Handler_WellKnown/Run_public_key_Without_public_prefix", func(t *testing.T) {
		t.Parallel()
		var IDKS *jose.JSONWebKeySet

		if reg.Config().HSMEnabled() {
			var err error
			IDKS, err = reg.KeyManager().GenerateAndPersistKeySet(context.TODO(), x.OpenIDConnectKeyName, "test-id-2", "RS256", "sig")
			require.NoError(t, err, "problem in generating keys")
		} else {
			var err error
			IDKS, err = jwk.GenerateJWK(jose.RS256, "test-id-2", "sig")
			require.NoError(t, err, "problem in generating keys")
			IDKS.Keys[0].KeyID = "test-id-2"
			require.NoError(t, reg.KeyManager().AddKeySet(context.TODO(), x.OpenIDConnectKeyName, IDKS))
		}

		res, err := http.Get(urlx.MustJoin(testServer.URL, JWKPath))
		require.NoError(t, err, "problem in http request")
		defer res.Body.Close() //nolint:errcheck

		var known jose.JSONWebKeySet
		err = json.NewDecoder(res.Body).Decode(&known)
		require.NoError(t, err, "problem in decoding response")
		if reg.Config().HSMEnabled() {
			require.GreaterOrEqual(t, len(known.Keys), 2)
		} else {
			require.GreaterOrEqual(t, len(known.Keys), 1)
		}

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
