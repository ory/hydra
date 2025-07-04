// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwtmiddleware_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rakutentech/jwk-go/jwk"
	"github.com/stretchr/testify/assert"

	"github.com/ory/x/jwtmiddleware"

	_ "embed"

	"github.com/tidwall/sjson"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
)

func mustString(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

var key *jwk.KeySpec

//go:embed stub/jwks.json
var rawKey []byte

func init() {
	key = &jwk.KeySpec{}
	if err := json.Unmarshal(rawKey, key); err != nil {
		panic(err)
	}
}

func newKeyServer(t *testing.T) string {
	public, err := key.PublicOnly()
	require.NoError(t, err)
	keys, err := json.Marshal(map[string]interface{}{
		"keys": []interface{}{
			public,
		},
	})
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(keys)
	}))
	t.Cleanup(ts.Close)
	return ts.URL
}

func TestSessionFromRequest(t *testing.T) {
	ks := newKeyServer(t)

	router := httprouter.New()
	router.GET("/anonymous", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("ok"))
	})
	router.GET("/me", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s, err := jwtmiddleware.SessionFromContext(r.Context())
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(s))
	})
	n := negroni.New()
	n.Use(jwtmiddleware.NewMiddleware(ks, jwtmiddleware.MiddlewareExcludePaths("/anonymous")))
	n.UseHandler(router)

	ts := httptest.NewServer(n)
	defer ts.Close()

	for k, tc := range []struct {
		token               string
		expectedStatusCode  int
		expectedErrorReason string
		expectedResponse    string
	}{
		// token without token
		{
			token:               "",
			expectedStatusCode:  401,
			expectedErrorReason: "Authorization header format must be Bearer {token}",
		},
		// token without kid
		{
			token: func() string {
				c := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{})
				delete(c.Header, "kid")
				s, err := c.SignedString(key.Key)
				require.NoError(t, err)
				return s
			}(),
			expectedStatusCode:  401,
			expectedErrorReason: "token is unverifiable: error while executing keyfunc: jwt from authorization HTTP header is missing value for \"kid\" in token header",
		},
		// token with int kid
		{
			token: func() string {
				c := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{})
				c.Header["kid"] = 42
				s, err := c.SignedString(key.Key)
				require.NoError(t, err)
				return s
			}(),
			expectedStatusCode:  401,
			expectedErrorReason: "token is unverifiable: error while executing keyfunc: jwt from authorization HTTP header is expecting string value for \"kid\" in tokenWithoutKid header but got: float64",
		},
		// token with unknown kid
		{
			token: func() string {
				c := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{})
				c.Header["kid"] = "not " + key.KeyID
				s, err := c.SignedString(key.Key)
				require.NoError(t, err)
				return s
			}(),
			expectedStatusCode:  401,
			expectedErrorReason: "token is unverifiable: error while executing keyfunc: unable to find JSON Web Key with ID: not b71ff5bd-a016-4ac0-9f3f-a172552578ea",
		},
		// token with valid kid
		{
			token: func() string {
				c := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
					"identity": map[string]interface{}{"email": "foo@bar.com"},
				})
				c.Header["kid"] = key.KeyID
				s, err := c.SignedString(key.Key)
				require.NoError(t, err)
				return s
			}(),
			expectedStatusCode: 200,
			expectedResponse:   mustString(sjson.SetRaw("{}", "identity.email", `"foo@bar.com"`)),
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/me", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "bearer "+tc.token)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode, string(body))
			assert.Equal(t, tc.expectedErrorReason, gjson.GetBytes(body, "error.reason").String())

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, string(body))
			}
		})
	}

	res, err := http.Get(ts.URL + "/anonymous")
	require.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
