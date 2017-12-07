// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	compose2 "github.com/ory/hydra/compose"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntrospectorSDK(t *testing.T) {
	tokens := pkg.Tokens(3)
	memoryStore := storage.NewExampleStore()
	memoryStore.Clients["my-client"].Scopes = []string{"fosite", "openid", "photos", "offline", "foo.*"}

	var localWarden, _ = compose2.NewMockFirewallWithStore("foo", "my-client", fosite.Arguments{"hydra.introspect"}, memoryStore, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"my-client"},
		Resources: []string{"rn:hydra:oauth2:tokens"},
		Actions:   []string{"introspect"},
		Effect:    ladon.AllowAccess,
	})

	l := logrus.New()
	l.Level = logrus.DebugLevel

	router := httprouter.New()
	handler := &oauth2.Handler{
		ScopeStrategy: fosite.WildcardScopeStrategy,
		OAuth2: compose.Compose(
			fc,
			memoryStore,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890")),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(pkg.MustINSECURELOWENTROPYRSAKEYFORTEST()),
			},
			nil,
			compose.OAuth2AuthorizeExplicitFactory,
			compose.OAuth2TokenIntrospectionFactory,
		),
		H:      herodot.NewJSONWriter(l),
		Issuer: "foobariss",
		W:      localWarden,
	}
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	now := time.Now().Round(time.Minute)
	createAccessTokenSession("alice", "my-client", tokens[0][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("siri", "my-client", tokens[1][0], now.Add(-time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("my-client", "my-client", tokens[2][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"hydra.introspect"})

	t.Run("TestIntrospect", func(t *testing.T) {
		for k, c := range []struct {
			token          string
			description    string
			expectInactive bool
			expectCode     int
			scopes         []string
			assert         func(*testing.T, *hydra.OAuth2TokenIntrospection)
			prepare        func(*testing.T) *hydra.OAuth2Api
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
			{
				description:    "should fail because username / password are invalid",
				token:          tokens[0][1],
				expectInactive: true,
				expectCode:     http.StatusForbidden,
				prepare: func(*testing.T) *hydra.OAuth2Api {
					client := hydra.NewOAuth2ApiWithBasePath(server.URL)
					client.Configuration.Username = "foo"
					client.Configuration.Password = "foo"
					return client
				},
			},
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
				prepare: func(*testing.T) *hydra.OAuth2Api {
					client := hydra.NewOAuth2ApiWithBasePath(server.URL)
					client.Configuration.DefaultHeader["Authorization"] = "bearer " + tokens[2][1]
					return client
				},
				token:          tokens[0][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.OAuth2TokenIntrospection) {
					assert.Equal(t, "alice", c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), c.Exp, "expires at")
					assert.Equal(t, now.Unix(), c.Iat, "issued at")
					assert.Equal(t, "foobariss", c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
			{
				description:    "should pass using regular authorization",
				token:          tokens[0][1],
				expectInactive: false,
				scopes:         []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.OAuth2TokenIntrospection) {
					assert.Equal(t, "alice", c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), c.Exp, "expires at")
					assert.Equal(t, now.Unix(), c.Iat, "issued at")
					assert.Equal(t, "foobariss", c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
				var client *hydra.OAuth2Api
				if c.prepare != nil {
					client = c.prepare(t)
				} else {
					client = hydra.NewOAuth2ApiWithBasePath(server.URL)
					client.Configuration.Username = "my-client"
					client.Configuration.Password = "foobar"
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
