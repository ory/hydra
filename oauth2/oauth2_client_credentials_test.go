// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/json"
	"math"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/hydra/internal/testhelpers"
	"github.com/ory/x/contextx"

	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/x/requirex"
)

func TestClientCredentials(t *testing.T) {
	ctx := context.Background()
	reg := internal.NewMockedRegistry(t, &contextx.Default{})
	reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
	public, admin := testhelpers.NewOAuth2Server(ctx, t, reg)

	var newCustomClient = func(t *testing.T, c *hc.Client) (*hc.Client, clientcredentials.Config) {
		unhashedSecret := c.Secret
		require.NoError(t, reg.ClientManager().CreateClient(ctx, c))
		return c, clientcredentials.Config{
			ClientID:       c.GetID(),
			ClientSecret:   unhashedSecret,
			TokenURL:       reg.Config().OAuth2TokenURL(ctx).String(),
			Scopes:         strings.Split(c.Scope, " "),
			EndpointParams: url.Values{"audience": c.Audience},
		}
	}

	var newClient = func(t *testing.T) (*hc.Client, clientcredentials.Config) {
		cc, config := newCustomClient(t, &hc.Client{
			Secret:        uuid.New().String(),
			RedirectURIs:  []string{public.URL + "/callback"},
			ResponseTypes: []string{"token"},
			GrantTypes:    []string{"client_credentials"},
			Scope:         "foobar",
			Audience:      []string{"https://api.ory.sh/"},
		})
		return cc, config
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

	var inspectToken = func(t *testing.T, token *goauth2.Token, cl *hc.Client, conf clientcredentials.Config, strategy string, expectedExp time.Time) {
		introspection := testhelpers.IntrospectToken(t, &goauth2.Config{ClientID: cl.GetID(), ClientSecret: conf.ClientSecret}, token.AccessToken, admin)

		check := func(res gjson.Result) {
			assert.EqualValues(t, cl.GetID(), res.Get("client_id").String(), "%s", res.Raw)
			assert.EqualValues(t, cl.GetID(), res.Get("sub").String(), "%s", res.Raw)
			assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), res.Get("iss").String(), "%s", res.Raw)

			assert.EqualValues(t, res.Get("nbf").Int(), res.Get("iat").Int(), "%s", res.Raw)
			requirex.EqualTime(t, expectedExp, time.Unix(res.Get("exp").Int(), 0), time.Second)

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

	var getAndInspectToken = func(t *testing.T, cl *hc.Client, conf clientcredentials.Config, strategy string, expectedExp time.Time) {
		token, err := getToken(t, conf)
		require.NoError(t, err)
		inspectToken(t, token, cl, conf, strategy, expectedExp)
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
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)))
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should pass without audience", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				conf.EndpointParams = url.Values{}
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)))
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should pass without scope", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				conf.Scopes = []string{}
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)))
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should grant default scopes if configured to do ", func(t *testing.T) {
		reg.Config().MustSet(ctx, config.KeyGrantAllClientCredentialsScopesPerDefault, true)

		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				cl, conf := newClient(t)
				defaultScope := conf.Scopes
				conf.Scopes = []string{}

				token, err := getToken(t, conf)
				require.NoError(t, err)

				// We reset this so that introspectToken is going to check for the default scope.
				conf.Scopes = defaultScope
				inspectToken(t, token, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)))
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should pass with custom client access token lifespan", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				secret := uuid.New().String()
				cl, conf := newCustomClient(t, &hc.Client{
					Secret:        secret,
					RedirectURIs:  []string{public.URL + "/callback"},
					ResponseTypes: []string{"token"},
					GrantTypes:    []string{"client_credentials"},
					Scope:         "foobar",
					Audience:      []string{"https://api.ory.sh/"},
				})
				testhelpers.UpdateClientTokenLifespans(t, &goauth2.Config{ClientID: cl.GetID(), ClientSecret: conf.ClientSecret}, cl.GetID(), testhelpers.TestLifespans, admin)
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(testhelpers.TestLifespans.ClientCredentialsGrantAccessTokenLifespan.Duration))
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should respect TTL", func(t *testing.T) {
		duration := time.Hour * 24 * 7
		reg.Config().MustSet(ctx, config.KeyAccessTokenLifespan, duration.String())

		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
				cl, conf := newClient(t)
				conf.Scopes = []string{}
				token, err := getToken(t, conf)
				require.NoError(t, err)

				assert.True(t, math.Abs(float64(time.Now().Add(duration).Round(time.Minute).Unix())-float64(token.Expiry.Round(time.Minute).Unix())) < 5)

				introspection := testhelpers.IntrospectToken(t, &goauth2.Config{ClientID: cl.GetID(), ClientSecret: conf.ClientSecret}, token.AccessToken, admin)
				assert.EqualValues(t, time.Now().Add(duration).Round(time.Minute), time.Unix(introspection.Get("exp").Int(), 0).Round(time.Minute))
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})
}
