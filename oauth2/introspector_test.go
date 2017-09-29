package oauth2_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"strings"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntrospectorSDK(t *testing.T) {
	tokens := pkg.Tokens(3)
	memoryStore := storage.NewExampleStore()
	memoryStore.Clients["my-client"].Scopes = []string{"fosite", "openid", "photos", "offline", "foo.*"}

	router := httprouter.New()
	handler := &oauth2.Handler{
		ScopeStrategy: fosite.WildcardScopeStrategy,
		OAuth2: compose.Compose(
			fc,
			memoryStore,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890")),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(pkg.MustRSAKey()),
			},
			nil,
			compose.OAuth2AuthorizeExplicitFactory,
			compose.OAuth2TokenIntrospectionFactory,
		),
		H:      herodot.NewJSONWriter(nil),
		Issuer: "foobariss",
	}
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	now := time.Now().Round(time.Minute)
	createAccessTokenSession("alice", "siri", tokens[0][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"})
	createAccessTokenSession("siri", "siri", tokens[1][0], now.Add(time.Hour), memoryStore, fosite.Arguments{"core", "foo"})
	createAccessTokenSession("siri", "doesnt-exist", tokens[2][0], now.Add(-time.Hour), memoryStore, fosite.Arguments{"core", "foo.*"})

	client := hydra.NewOAuth2ApiWithBasePath(server.URL)
	client.Configuration.Username = "my-client"
	client.Configuration.Password = "foobar"

	t.Run("TestIntrospect", func(t *testing.T) {
		for k, c := range []struct {
			token       string
			description string
			expectErr   bool
			scopes      []string
			assert      func(*testing.T, *hydra.OAuth2TokenIntrospection)
		}{
			{
				description: "should fail because invalid token was supplied",
				token:       "invalid",
				expectErr:   true,
			},
			{
				description: "should fail because token is expired",
				token:       tokens[2][1],
				expectErr:   true,
			},
			{
				description: "should pass",
				token:       tokens[1][1],
				expectErr:   false,
			},
			{
				description: "should fail because scope `foo.bar` was requested but only `foo` is granted",
				token:       tokens[1][1],
				expectErr:   true,
				scopes:      []string{"foo.bar"},
			},
			{
				description: "should pass",
				token:       tokens[0][1],
				expectErr:   false,
			},
			{
				description: "should pass",
				token:       tokens[0][1],
				expectErr:   false,
				scopes:      []string{"foo.bar"},
				assert: func(t *testing.T, c *hydra.OAuth2TokenIntrospection) {
					assert.Equal(t, "alice", c.Sub)
					assert.Equal(t, now.Add(time.Hour).Unix(), c.Exp, "expires at")
					assert.Equal(t, now.Unix(), c.Iat, "issued at")
					assert.Equal(t, "foobariss", c.Iss, "issuer")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Ext)
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
				ctx, response, err := client.IntrospectOAuth2Token(c.token, strings.Join(c.scopes, " "))
				require.NoError(t, err)
				require.EqualValues(t, http.StatusOK, response.StatusCode)
				t.Logf("Got %s", response.Payload)

				if c.expectErr {
					assert.False(t, ctx.Active)
				} else {
					assert.True(t, ctx.Active)
				}

				if !c.expectErr && c.assert != nil {
					c.assert(t, ctx)
				}
			})
		}
	})
}
