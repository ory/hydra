// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/prometheusx"
)

func TestIntrospectorSDK(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyScopeStrategy: "wildcard",
		config.KeyIssuerURL:     "https://foobariss",
	})))

	testhelpers.MustEnsureRegistryKeys(t, reg, x.OpenIDConnectKeyName)
	internal.AddFositeExamples(t, reg)

	tokens := Tokens(reg.OAuth2ProviderConfig(), 4)

	c, err := reg.ClientManager().GetConcreteClient(context.TODO(), "my-client")
	require.NoError(t, err)
	c.Scope = "fosite,openid,photos,offline,foo.*"
	require.NoError(t, reg.ClientManager().UpdateClient(context.TODO(), c))

	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	router := x.NewRouterAdmin(metrics)
	handler := oauth2.NewHandler(reg)
	handler.SetAdminRoutes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	now := time.Now().UTC().Round(time.Minute)
	createAccessTokenSession(t, "alice", "my-client", tokens[0].sig, now.Add(time.Hour), reg.OAuth2Storage(), fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession(t, "siri", "my-client", tokens[1].sig, now.Add(-time.Hour), reg.OAuth2Storage(), fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession(t, "my-client", "my-client", tokens[2].sig, now.Add(time.Hour), reg.OAuth2Storage(), fosite.Arguments{"hydra.introspect"})
	createAccessTokenSessionPairwise(t, "alice", "my-client", tokens[3].sig, now.Add(time.Hour), reg.OAuth2Storage(), fosite.Arguments{"core", "foo.*"}, "alice-obfuscated")

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
				token:          tokens[1].tok,
				expectInactive: true,
			},
			// {
			//	description:    "should fail because username / password are invalid",
			//	token:          tokens[0][1],
			//	expectInactive: true,
			//	expectCode:     http.StatusUnauthorized,
			//	prepare: func(*testing.T) *hydra.OAuth2API.{
			//		client := hydra.Ne.OAuth2API.ithBasePath(server.URL)
			//		client.config.Username = "foo"
			//		client.config.Password = "foo"
			//		return client
			//	},
			// },
			{
				description:    "should fail because scope `bar` was requested but only `foo` is granted",
				token:          tokens[0].tok,
				expectInactive: true,
				scopes:         []string{"bar"},
			},
			{
				description:    "should pass",
				token:          tokens[0].tok,
				expectInactive: false,
			},
			{
				description:    "should pass using bearer authorization",
				token:          tokens[0].tok,
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
				token:          tokens[0].tok,
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
				token:          tokens[3].tok,
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

				ctx, _, err := client.OAuth2API.IntrospectOAuth2Token(context.Background()).
					Token(c.token).Scope(strings.Join(c.scopes, " ")).Execute()
				require.NoError(t, err)

				assert.Equal(t, c.expectInactive, !ctx.Active)

				if !c.expectInactive && c.assert != nil {
					c.assert(t, ctx)
				}
			})
		}
	})
}
