// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pkg/errors"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func mustGenerateAssertion(t *testing.T, claims jwt.MapClaims, key *rsa.PrivateKey, kid string) string {
	token := jwt.NewWithClaims(jose.RS256, claims)
	if kid != "" {
		token.Header["kid"] = kid
	}
	tokenString, err := token.SignedString(key)
	require.NoError(t, err)
	return tokenString
}

func mustGenerateHSAssertion(t *testing.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jose.HS256, claims)
	tokenString, err := token.SignedString([]byte("aaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbcccccccccccccccccccccddddddddddddddddddddddd"))
	require.NoError(t, err)
	return tokenString
}

func mustGenerateNoneAssertion(t *testing.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	return tokenString
}

func TestAuthorizeRequestParametersFromOpenIDConnectRequest(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	jwks := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				KeyID: "kid-foo",
				Use:   "sig",
				Key:   &key.PublicKey,
			},
		},
	}

	validRequestObject := mustGenerateAssertion(t, jwt.MapClaims{"scope": "foo", "foo": "bar", "baz": "baz", "response_type": "token", "response_mode": "post_form"}, key, "kid-foo")
	validRequestObjectWithoutKid := mustGenerateAssertion(t, jwt.MapClaims{"scope": "foo", "foo": "bar", "baz": "baz"}, key, "")
	validNoneRequestObject := mustGenerateNoneAssertion(t, jwt.MapClaims{"scope": "foo", "foo": "bar", "baz": "baz", "state": "some-state"})

	var reqH http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(validRequestObject))
	}
	reqTS := httptest.NewServer(reqH)
	defer reqTS.Close()

	var hJWK http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewEncoder(rw).Encode(jwks))
	}
	reqJWK := httptest.NewServer(hJWK)
	defer reqJWK.Close()

	f := &Fosite{Config: &Config{JWKSFetcherStrategy: NewDefaultJWKSFetcherStrategy()}}
	for k, tc := range []struct {
		client Client
		form   url.Values
		d      string

		expectErr       error
		expectErrReason string
		expectForm      url.Values
	}{
		{
			d:          "should pass because no request context given and not openid",
			form:       url.Values{},
			expectErr:  nil,
			expectForm: url.Values{},
		},
		{
			d:          "should pass because no request context given",
			form:       url.Values{"scope": {"openid"}},
			expectErr:  nil,
			expectForm: url.Values{"scope": {"openid"}},
		},
		{
			d:          "should pass because request context given but not openid",
			form:       url.Values{"request": {"foo"}},
			expectErr:  nil,
			expectForm: url.Values{"request": {"foo"}},
		},
		{
			d:          "should fail because not an OpenIDConnect compliant client",
			form:       url.Values{"scope": {"openid"}, "request": {"foo"}},
			expectErr:  ErrRequestNotSupported,
			expectForm: url.Values{"scope": {"openid"}},
		},
		{
			d:          "should fail because not an OpenIDConnect compliant client",
			form:       url.Values{"scope": {"openid"}, "request_uri": {"foo"}},
			expectErr:  ErrRequestURINotSupported,
			expectForm: url.Values{"scope": {"openid"}},
		},
		{
			d:          "should fail because token invalid an no key set",
			form:       url.Values{"scope": {"openid"}, "request_uri": {"foo"}},
			client:     &DefaultOpenIDConnectClient{RequestObjectSigningAlgorithm: "RS256"},
			expectErr:  ErrInvalidRequest,
			expectForm: url.Values{"scope": {"openid"}},
		},
		{
			d:          "should fail because token invalid",
			form:       url.Values{"scope": {"openid"}, "request": {"foo"}},
			client:     &DefaultOpenIDConnectClient{JSONWebKeys: jwks, RequestObjectSigningAlgorithm: "RS256"},
			expectErr:  ErrInvalidRequestObject,
			expectForm: url.Values{"scope": {"openid"}},
		},
		{
			d:               "should fail because kid does not exist",
			form:            url.Values{"scope": {"openid"}, "request": {mustGenerateAssertion(t, jwt.MapClaims{}, key, "does-not-exists")}},
			client:          &DefaultOpenIDConnectClient{JSONWebKeys: jwks, RequestObjectSigningAlgorithm: "RS256"},
			expectErr:       ErrInvalidRequestObject,
			expectErrReason: "Unable to retrieve RSA signing key from OAuth 2.0 Client. The JSON Web Token uses signing key with kid 'does-not-exists', which could not be found.",
			expectForm:      url.Values{"scope": {"openid"}},
		},
		{
			d:               "should fail because not RS256 token",
			form:            url.Values{"scope": {"openid"}, "request": {mustGenerateHSAssertion(t, jwt.MapClaims{})}},
			client:          &DefaultOpenIDConnectClient{JSONWebKeys: jwks, RequestObjectSigningAlgorithm: "RS256"},
			expectErr:       ErrInvalidRequestObject,
			expectErrReason: "The request object uses signing algorithm 'HS256', but the requested OAuth 2.0 Client enforces signing algorithm 'RS256'.",
			expectForm:      url.Values{"scope": {"openid"}},
		},
		{
			d:      "should pass and set request parameters properly",
			form:   url.Values{"scope": {"openid"}, "response_type": {"code"}, "response_mode": {"none"}, "request": {validRequestObject}},
			client: &DefaultOpenIDConnectClient{JSONWebKeys: jwks, RequestObjectSigningAlgorithm: "RS256"},
			// The values from form are overwritten by the request object.
			expectForm: url.Values{"response_type": {"token"}, "response_mode": {"post_form"}, "scope": {"foo openid"}, "request": {validRequestObject}, "foo": {"bar"}, "baz": {"baz"}},
		},
		{
			d:          "should pass even if kid is unset",
			form:       url.Values{"scope": {"openid"}, "request": {validRequestObjectWithoutKid}},
			client:     &DefaultOpenIDConnectClient{JSONWebKeys: jwks, RequestObjectSigningAlgorithm: "RS256"},
			expectForm: url.Values{"scope": {"foo openid"}, "request": {validRequestObjectWithoutKid}, "foo": {"bar"}, "baz": {"baz"}},
		},
		{
			d:          "should fail because request uri is not whitelisted",
			form:       url.Values{"scope": {"openid"}, "request_uri": {reqTS.URL}},
			client:     &DefaultOpenIDConnectClient{JSONWebKeysURI: reqJWK.URL, RequestObjectSigningAlgorithm: "RS256"},
			expectForm: url.Values{"scope": {"foo openid"}, "request_uri": {reqTS.URL}, "foo": {"bar"}, "baz": {"baz"}},
			expectErr:  ErrInvalidRequestURI,
		},
		{
			d:          "should pass and set request_uri parameters properly and also fetch jwk from remote",
			form:       url.Values{"scope": {"openid"}, "request_uri": {reqTS.URL}},
			client:     &DefaultOpenIDConnectClient{JSONWebKeysURI: reqJWK.URL, RequestObjectSigningAlgorithm: "RS256", RequestURIs: []string{reqTS.URL}},
			expectForm: url.Values{"response_type": {"token"}, "response_mode": {"post_form"}, "scope": {"foo openid"}, "request_uri": {reqTS.URL}, "foo": {"bar"}, "baz": {"baz"}},
		},
		{
			d:          "should pass when request object uses algorithm none",
			form:       url.Values{"scope": {"openid"}, "request": {validNoneRequestObject}},
			client:     &DefaultOpenIDConnectClient{JSONWebKeysURI: reqJWK.URL, RequestObjectSigningAlgorithm: "none"},
			expectForm: url.Values{"state": {"some-state"}, "scope": {"foo openid"}, "request": {validNoneRequestObject}, "foo": {"bar"}, "baz": {"baz"}},
		},
		{
			d:          "should pass when request object uses algorithm none and the client did not explicitly allow any algorithm",
			form:       url.Values{"scope": {"openid"}, "request": {validNoneRequestObject}},
			client:     &DefaultOpenIDConnectClient{JSONWebKeysURI: reqJWK.URL},
			expectForm: url.Values{"state": {"some-state"}, "scope": {"foo openid"}, "request": {validNoneRequestObject}, "foo": {"bar"}, "baz": {"baz"}},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			req := &AuthorizeRequest{
				Request: Request{
					Client: tc.client,
					Form:   tc.form,
				},
			}

			err := f.authorizeRequestParametersFromOpenIDConnectRequest(context.Background(), req, false)
			if tc.expectErr != nil {
				require.EqualError(t, err, tc.expectErr.Error(), "%+v", err)
				if tc.expectErrReason != "" {
					real := new(RFC6749Error)
					require.True(t, errors.As(err, &real))
					assert.EqualValues(t, tc.expectErrReason, real.Reason())
				}
			} else {
				if err != nil {
					real := new(RFC6749Error)
					errors.As(err, &real)
					require.NoErrorf(t, err, "Hint: %v\nDebug:%v", real.HintField, real.DebugField)
				}
				require.NoErrorf(t, err, "%+v", err)
				require.Equal(t, len(tc.expectForm), len(req.Form))
				for k, v := range tc.expectForm {
					assert.EqualValues(t, v, req.Form[k])
				}
			}
		})
	}
}
