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

	hydra "github.com/ory/hydra-client-go"

	"github.com/ory/x/httprouterx"

	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
)

func TestIntrospectorSDK(t *testing.T) {
	ctx := context.Background()
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(ctx, config.KeyScopeStrategy, "wildcard")
	conf.MustSet(ctx, config.KeyIssuerURL, "https://foobariss")
	reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)
	internal.AddFositeExamples(reg)

	tokens := Tokens(reg.OAuth2ProviderConfig(), 4)

	c, err := reg.ClientManager().GetConcreteClient(context.TODO(), "my-client")
	require.NoError(t, err)
	c.Scope = "fosite,openid,photos,offline,foo.*"
	require.NoError(t, reg.ClientManager().UpdateClient(context.TODO(), c))

	router := x.NewRouterAdmin(conf.AdminURL)
	handler := reg.OAuth2Handler()
	handler.SetRoutes(router, &httprouterx.RouterPublic{Router: router.Router}, func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)
	defer server.Close()

	now := time.Now().UTC().Round(time.Minute)
	createAccessTokenSession("alice", "my-client", tokens[0][0], now.Add(time.Hour), reg.OAuth2Storage(), fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("siri", "my-client", tokens[1][0], now.Add(-time.Hour), reg.OAuth2Storage(), fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("my-client", "my-client", tokens[2][0], now.Add(time.Hour), reg.OAuth2Storage(), fosite.Arguments{"hydra.introspect"})
	createAccessTokenSessionPairwise("alice", "my-client", tokens[3][0], now.Add(time.Hour), reg.OAuth2Storage(), fosite.Arguments{"core", "foo.*"}, "alice-obfuscated")

	t.Run("TestIntrospect", func(t *testing.T) {
		for k, c := range []struct {
			token          string
			description    string
			expectInactive bool
			scopes         []string
			assert         func(*testing.T, *hydra.IntrospectedOAuth2Token)
			prepare        func(*testing.T) *hydra.APIClient
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
			// {
			//	description:    "should fail because username / password are invalid",
			//	token:          tokens[0][1],
			//	expectInactive: true,
			//	expectCode:     http.StatusUnauthorized,
			//	prepare: func(*testing.T) *hydra.OAuth2Api {
			//		client := hydra.NewOAuth2ApiWithBasePath(server.URL)
			//		client.config.Username = "foo"
			//		client.config.Password = "foo"
			//		return client
			//	},
			// },
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
				description:    "should pass using bearer authorization",
				token:          tokens[0][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.IntrospectedOAuth2Token) {
					assert.Equal(t, "alice", *c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), *c.Exp, "expires at")
					assert.Equal(t, now.Unix(), *c.Iat, "issued at")
					assert.Equal(t, "https://foobariss", *c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
			{
				description:    "should pass using regular authorization",
				token:          tokens[0][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.IntrospectedOAuth2Token) {
					assert.Equal(t, "core foo.*", *c.Scope)
					assert.Equal(t, "alice", *c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), *c.Exp, "expires at")
					assert.Equal(t, now.Unix(), *c.Iat, "issued at")
					assert.Equal(t, "https://foobariss", *c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
			{
				description:    "should pass and check for obfuscated subject",
				token:          tokens[3][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.IntrospectedOAuth2Token) {
					assert.Equal(t, "alice", *c.Sub)
					assert.Equal(t, "alice-obfuscated", *c.ObfuscatedSubject)
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
				var client *hydra.APIClient
				if c.prepare != nil {
					client = c.prepare(t)
				} else {
					client = hydra.NewAPIClient(hydra.NewConfiguration())
					client.GetConfig().Servers = hydra.ServerConfigurations{{URL: server.URL}}
				}

				ctx, _, err := client.V0alpha2Api.AdminIntrospectOAuth2Token(context.Background()).
					Token(c.token).
					Scope(strings.Join(c.scopes, " ")).
					Execute()
				require.NoError(t, err)

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
