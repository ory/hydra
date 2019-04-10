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

package oauth2_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
)

func TestClientCredentials(t *testing.T) {
	for _, tc := range []struct{ d string }{{d: "opaque"}, {d: "jwt"}} {
		t.Run("tc="+tc.d, func(t *testing.T) {
			conf := internal.NewConfigurationWithDefaults()
			viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Second)
			viper.Set(configuration.ViperKeyAccessTokenStrategy, tc.d)
			reg := internal.NewRegistry(conf)

			router := x.NewRouterPublic()
			ts := httptest.NewServer(router)
			defer ts.Close()
			viper.Set(configuration.ViperKeyIssuerURL, ts.URL)

			handler := NewHandler(reg, conf)
			handler.SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
				return h
			})

			require.NoError(t, reg.ClientManager().CreateClient(context.TODO(), &hc.Client{
				ClientID:      "app-client",
				Secret:        "secret",
				RedirectURIs:  []string{ts.URL + "/callback"},
				ResponseTypes: []string{"token"},
				GrantTypes:    []string{"client_credentials"},
				Scope:         "foobar",
				Audience:      []string{"https://api.ory.sh/"},
			}))

			for k, ccc := range []struct {
				d                 string
				c                 *clientcredentials.Config
				assertAccessToken func(*testing.T, string)
				expectAccessToken bool
				expectError       bool
			}{
				{
					d: "should fail because audience is not allowed",
					c: &clientcredentials.Config{
						ClientID:       "app-client",
						ClientSecret:   "secret",
						TokenURL:       ts.URL + "/oauth2/token",
						Scopes:         []string{"foobar"},
						EndpointParams: url.Values{"audience": {"https://not-api.ory.sh/"}},
					},
					expectAccessToken: false,
					expectError:       true,
				},
				{
					d: "should fail because scope is not allowed",
					c: &clientcredentials.Config{
						ClientID:     "app-client",
						ClientSecret: "secret",
						TokenURL:     ts.URL + "/oauth2/token",
						Scopes:       []string{"not-foobar"},
					},
					expectAccessToken: false,
					expectError:       true,
				},
				{
					d: "should pass with audience",
					c: &clientcredentials.Config{
						ClientID:       "app-client",
						ClientSecret:   "secret",
						TokenURL:       ts.URL + "/oauth2/token",
						Scopes:         []string{"foobar"},
						EndpointParams: url.Values{"audience": {"https://api.ory.sh/"}},
					},
					assertAccessToken: func(t *testing.T, token string) {
						if tc.d != "jwt" {
							return
						}
						body, err := jwt.DecodeSegment(strings.Split(token, ".")[1])
						require.NoError(t, err)

						data := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &data))

						assert.EqualValues(t, "app-client", data["client_id"])
						assert.EqualValues(t, "app-client", data["sub"])
						assert.NotEmpty(t, data["iss"])
						assert.NotEmpty(t, data["jti"])
						assert.NotEmpty(t, data["exp"])
						assert.NotEmpty(t, data["iat"])
						assert.NotEmpty(t, data["nbf"])
						assert.EqualValues(t, data["nbf"], data["iat"])
						assert.EqualValues(t, []interface{}{"foobar"}, data["scp"])
						assert.EqualValues(t, []interface{}{"https://api.ory.sh/"}, data["aud"])
					},
					expectAccessToken: true,
					expectError:       false,
				},
				{
					d: "should pass without audience",
					c: &clientcredentials.Config{
						ClientID:     "app-client",
						ClientSecret: "secret",
						TokenURL:     ts.URL + "/oauth2/token",
						Scopes:       []string{"foobar"},
					},
					assertAccessToken: func(t *testing.T, token string) {
						if tc.d != "jwt" {
							return
						}
						body, err := jwt.DecodeSegment(strings.Split(token, ".")[1])
						require.NoError(t, err)

						data := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &data))

						assert.EqualValues(t, "app-client", data["client_id"])
						assert.EqualValues(t, "app-client", data["sub"])
						assert.NotEmpty(t, data["iss"])
						assert.NotEmpty(t, data["jti"])
						assert.NotEmpty(t, data["exp"])
						assert.NotEmpty(t, data["iat"])
						assert.NotEmpty(t, data["nbf"])
						assert.Empty(t, data["aud"])
						assert.EqualValues(t, data["nbf"], data["iat"])
						assert.EqualValues(t, []interface{}{"foobar"}, data["scp"])
					},
					expectAccessToken: true,
					expectError:       false,
				},
			} {
				t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
					ccc.c.AuthStyle = goauth2.AuthStyleInHeader
					tok, err := ccc.c.Token(context.Background())

					if ccc.expectError {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}

					if ccc.expectAccessToken {
						assert.NotEmpty(t, tok.AccessToken)
						assert.Empty(t, tok.RefreshToken)
						assert.Empty(t, tok.Extra("id_token"))
					} else {
						assert.Nil(t, tok)
					}

					if ccc.assertAccessToken != nil {
						ccc.assertAccessToken(t, tok.AccessToken)
					}
				})
			}
		})
	}
}
