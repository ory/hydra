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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

func TestIntrospectorSDK(t *testing.T) {
	tokens := pkg.Tokens(4)
	memoryStore := storage.NewExampleStore()
	memoryStore.Clients["my-client"].(*fosite.DefaultClient).Scopes = []string{"fosite", "openid", "photos", "offline", "foo.*"}

	l := logrus.New()
	l.Level = logrus.DebugLevel

	jm := &jwk.MemoryManager{Keys: map[string]*jose.JSONWebKeySet{}}
	keys, err := (&jwk.RS256Generator{}).Generate("", "sig")
	require.NoError(t, err)
	require.NoError(t, jm.AddKeySet(context.TODO(), oauth2.OpenIDConnectKeyName, keys))
	jwtStrategy, err := jwk.NewRS256JWTStrategy(jm, oauth2.OpenIDConnectKeyName)

	router := httprouter.New()
	handler := &oauth2.Handler{
		ScopeStrategy: fosite.WildcardScopeStrategy,
		OAuth2: compose.Compose(
			fc,
			memoryStore,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890"), nil),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(fc, pkg.MustINSECURELOWENTROPYRSAKEYFORTEST()),
			},
			nil,
			compose.OAuth2AuthorizeExplicitFactory,
			compose.OAuth2TokenIntrospectionFactory,
		),
		H:                 herodot.NewJSONWriter(l),
		IssuerURL:         "foobariss",
		OpenIDJWTStrategy: jwtStrategy,
	}
	handler.SetRoutes(router, router, func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)

	now := time.Now().UTC().Round(time.Minute)
	createAccessTokenSession("alice", "my-client", tokens[0][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("siri", "my-client", tokens[1][0], now.Add(-time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("my-client", "my-client", tokens[2][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"hydra.introspect"})
	createAccessTokenSessionPairwise("alice", "my-client", tokens[3][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"}, "alice-obfuscated")

	t.Run("TestIntrospect", func(t *testing.T) {
		for k, c := range []struct {
			token          string
			description    string
			expectInactive bool
			expectCode     int
			scopes         []string
			assert         func(*testing.T, *hydra.OAuth2TokenIntrospection)
			prepare        func(*testing.T) *hydra.AdminApi
		}{
			{
				description:    "should fail because invalid token was supplied",
				token:          "invalid",
				expectInactive: true,
			},
			{
				description:    "should fail because token is expired",
				token:          tokens[1][1],
				expectInactive: true,
			},
			//{
			//	description:    "should fail because username / password are invalid",
			//	token:          tokens[0][1],
			//	expectInactive: true,
			//	expectCode:     http.StatusUnauthorized,
			//	prepare: func(*testing.T) *hydra.OAuth2Api {
			//		client := hydra.NewOAuth2ApiWithBasePath(server.URL)
			//		client.Configuration.Username = "foo"
			//		client.Configuration.Password = "foo"
			//		return client
			//	},
			//},
			{
				description:    "should fail because scope `bar` was requested but only `foo` is granted",
				token:          tokens[0][1],
				expectInactive: true,
				scopes:         []string{"bar"},
			},
			{
				description:    "should pass",
				token:          tokens[0][1],
				expectInactive: false,
			},
			{
				description: "should pass using bearer authorization",
				//prepare: func(*testing.T) *hydra.OAuth2Api {
				//	client := hydra.NewOAuth2ApiWithBasePath(server.URL)
				//	client.Configuration.DefaultHeader["Authorization"] = "bearer " + tokens[2][1]
				//	return client
				//},
				token:          tokens[0][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.OAuth2TokenIntrospection) {
					assert.Equal(t, "alice", c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), c.Exp, "expires at")
					assert.Equal(t, now.Unix(), c.Iat, "issued at")
					assert.Equal(t, "foobariss/", c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
			{
				description:    "should pass using regular authorization",
				token:          tokens[0][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.OAuth2TokenIntrospection) {
					assert.Equal(t, "core foo.*", c.Scope)
					assert.Equal(t, "alice", c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), c.Exp, "expires at")
					assert.Equal(t, now.Unix(), c.Iat, "issued at")
					assert.Equal(t, "foobariss/", c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
			{
				description:    "should pass and check for obfuscated subject",
				token:          tokens[3][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.OAuth2TokenIntrospection) {
					assert.Equal(t, "alice", c.Sub)
					assert.Equal(t, "alice-obfuscated", c.ObfuscatedSubject)
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
				var client *hydra.AdminApi
				if c.prepare != nil {
					client = c.prepare(t)
				} else {
					client = hydra.NewAdminApiWithBasePath(server.URL)
					//client.Configuration.Username = "my-client"
					//client.Configuration.Password = "foobar"
				}

				ctx, response, err := client.IntrospectOAuth2Token(c.token, strings.Join(c.scopes, " "))
				require.NoError(t, err)

				if c.expectCode == 0 {
					require.EqualValues(t, http.StatusOK, response.StatusCode)
				} else {
					require.EqualValues(t, c.expectCode, response.StatusCode)
				}

				if c.expectInactive {
					assert.False(t, ctx.Active)
				} else {
					assert.True(t, ctx.Active)
				}

				if !c.expectInactive && c.assert != nil {
					c.assert(t, ctx)
				}
			})
		}
	})
}
