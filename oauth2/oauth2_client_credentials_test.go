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
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/hydra/internal/testhelpers"

	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
)

func TestClientCredentials(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	reg.Config().MustSet(config.KeyAccessTokenStrategy, "opaque")
	public, admin := testhelpers.NewOAuth2Server(t, reg)

	var newClient = func(t *testing.T) (*hc.Client, clientcredentials.Config) {
		secret := uuid.New().String()
		c := &hc.Client{
			OutfacingID:   uuid.New().String(),
			Secret:        secret,
			RedirectURIs:  []string{public.URL + "/callback"},
			ResponseTypes: []string{"token"},
			GrantTypes:    []string{"client_credentials"},
			Scope:         "foobar",
			Audience:      []string{"https://api.ory.sh/"},
		}
		require.NoError(t, reg.ClientManager().CreateClient(context.TODO(), c))
		return c, clientcredentials.Config{
			ClientID:       c.OutfacingID,
			ClientSecret:   secret,
			TokenURL:       reg.Config().OAuth2TokenURL().String(),
			Scopes:         strings.Split(c.Scope, " "),
			EndpointParams: url.Values{"audience": c.Audience},
		}
	}

	var getToken = func(t *testing.T, conf clientcredentials.Config) (*goauth2.Token, error) {
		conf.AuthStyle = goauth2.AuthStyleInHeader
		return conf.Token(context.Background())
	}

	var encodeOr = func(t *testing.T, val interface{}, or string) string {
		out, err := json.Marshal(val)
		require.NoError(t, err)
		if string(out) == "null" {
			return or
		}

		return string(out)
	}

	var inspectToken = func(t *testing.T, token *goauth2.Token, cl *hc.Client, conf clientcredentials.Config, strategy string) {
		introspection := testhelpers.IntrospectToken(t, &goauth2.Config{ClientID: cl.OutfacingID, ClientSecret: conf.ClientSecret}, token, admin)

		check := func(res gjson.Result) {
			assert.EqualValues(t, cl.OutfacingID, res.Get("client_id").String(), "%s", res.Raw)
			assert.EqualValues(t, cl.OutfacingID, res.Get("sub").String(), "%s", res.Raw)
			assert.EqualValues(t, reg.Config().IssuerURL().String(), res.Get("iss").String(), "%s", res.Raw)

			assert.EqualValues(t, res.Get("nbf").Int(), res.Get("iat").Int(), "%s", res.Raw)
			assert.True(t, res.Get("exp").Int() >= res.Get("iat").Int()+int64(reg.Config().AccessTokenLifespan().Seconds()), "%s", res.Raw)

			assert.EqualValues(t, encodeOr(t, conf.EndpointParams["audience"], "[]"), res.Get("aud").Raw, "%s", res.Raw)
		}

		check(introspection)
		assert.True(t, introspection.Get("active").Bool())
		assert.EqualValues(t, "access_token", introspection.Get("token_use").String())
		assert.EqualValues(t, "Bearer", introspection.Get("token_type").String())
		assert.EqualValues(t, strings.Join(conf.Scopes, " "), introspection.Get("scope").String(), "%s", introspection.Raw)

		if strategy != "jwt" {
			return
		}

		body, err := x.DecodeSegment(strings.Split(token.AccessToken, ".")[1])
		require.NoError(t, err)

		jwtClaims := gjson.ParseBytes(body)
		assert.NotEmpty(t, jwtClaims.Get("jti").String())
		assert.EqualValues(t, encodeOr(t, conf.Scopes, "[]"), jwtClaims.Get("scp").Raw, "%s", introspection.Raw)
		check(jwtClaims)
	}

	var getAndInspectToken = func(t *testing.T, cl *hc.Client, conf clientcredentials.Config, strategy string) {
		token, err := getToken(t, conf)
		require.NoError(t, err)
		inspectToken(t, token, cl, conf, strategy)
	}

	t.Run("case=should fail because audience is not allowed", func(t *testing.T) {
		_, conf := newClient(t)
		conf.EndpointParams = url.Values{"audience": {"https://not-api.ory.sh/"}}
		_, err := getToken(t, conf)
		require.Error(t, err)
	})

	t.Run("case=should fail because scope is not allowed", func(t *testing.T) {
		_, conf := newClient(t)
		conf.Scopes = []string{"not-allowed-scope"}
		_, err := getToken(t, conf)
		require.Error(t, err)
	})

	t.Run("case=should pass with audience", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				getAndInspectToken(t, cl, conf, strategy)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should pass without audience", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				conf.EndpointParams = url.Values{}
				getAndInspectToken(t, cl, conf, strategy)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should pass without scope", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				conf.Scopes = []string{}
				getAndInspectToken(t, cl, conf, strategy)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should grant default scopes if configured to do ", func(t *testing.T) {
		reg.Config().MustSet(config.KeyGrantAllClientCredentialsScopesPerDefault, true)

		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				defaultScope := conf.Scopes
				conf.Scopes = []string{}

				token, err := getToken(t, conf)
				require.NoError(t, err)

				// We reset this so that introspectToken is going to check for the default scope.
				conf.Scopes = defaultScope
				inspectToken(t, token, cl, conf, strategy)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})
}
