// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	hc "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/internal/testhelpers"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/configx"
)

func TestClientCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue(config.KeyAccessTokenStrategy, "opaque")))
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
		return newCustomClient(t, &hc.Client{
			Secret:        uuid.Must(uuid.NewV4()).String(),
			RedirectURIs:  []string{public.URL + "/callback"},
			ResponseTypes: []string{"token"},
			GrantTypes:    []string{"client_credentials"},
			Scope:         "foobar",
			Audience:      []string{"https://api.ory.sh/"},
		})
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

	var inspectToken = func(t *testing.T, token *goauth2.Token, cl *hc.Client, conf clientcredentials.Config, strategy string, expectedExp time.Time, checkExtraClaims bool) {
		introspection := testhelpers.IntrospectToken(t, token.AccessToken, admin)

		check := func(res gjson.Result) {
			assert.EqualValues(t, cl.GetID(), res.Get("client_id").String(), "%s", res.Raw)
			assert.EqualValues(t, cl.GetID(), res.Get("sub").String(), "%s", res.Raw)
			assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), res.Get("iss").String(), "%s", res.Raw)

			assert.EqualValues(t, res.Get("nbf").Int(), res.Get("iat").Int(), "%s", res.Raw)
			assert.WithinDuration(t, expectedExp, time.Unix(res.Get("exp").Int(), 0), 2*time.Second)

			assert.EqualValues(t, encodeOr(t, conf.EndpointParams["audience"], "[]"), res.Get("aud").Raw, "%s", res.Raw)

			if checkExtraClaims {
				require.True(t, res.Get("ext.hooked").Bool())
			}
		}

		check(introspection)
		assert.True(t, introspection.Get("active").Bool())
		assert.EqualValues(t, "access_token", introspection.Get("token_use").String())
		assert.EqualValues(t, "Bearer", introspection.Get("token_type").String())
		assert.EqualValues(t, strings.Join(conf.Scopes, " "), introspection.Get("scope").String(), "%s", introspection.Raw)

		if strategy != "jwt" {
			return
		}

		jwtClaims := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, token.AccessToken))
		assert.NotEmpty(t, jwtClaims.Get("jti").String())
		assert.EqualValues(t, encodeOr(t, conf.Scopes, "[]"), jwtClaims.Get("scp").Raw, "%s", introspection.Raw)
		check(jwtClaims)
	}

	var getAndInspectToken = func(t *testing.T, cl *hc.Client, conf clientcredentials.Config, strategy string, expectedExp time.Time, checkExtraClaims bool) {
		token, err := getToken(t, conf)
		require.NoError(t, err)
		inspectToken(t, token, cl, conf, strategy, expectedExp, checkExtraClaims)
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
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), false)
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
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), false)
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
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), false)
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
				inspectToken(t, token, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), false)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=should pass with custom client access token lifespan", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				cl, conf := newCustomClient(t, &hc.Client{
					Secret:        uuid.Must(uuid.NewV4()).String(),
					RedirectURIs:  []string{public.URL + "/callback"},
					ResponseTypes: []string{"token"},
					GrantTypes:    []string{"client_credentials"},
					Scope:         "foobar",
					Audience:      []string{"https://api.ory.sh/"},
				})
				testhelpers.UpdateClientTokenLifespans(t, &goauth2.Config{ClientID: cl.GetID(), ClientSecret: conf.ClientSecret}, cl.GetID(), testhelpers.TestLifespans, admin)
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(testhelpers.TestLifespans.ClientCredentialsGrantAccessTokenLifespan.Duration), false)
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
				_, conf := newClient(t)
				conf.Scopes = []string{}
				token, err := getToken(t, conf)
				require.NoError(t, err)
				expected := time.Now().Add(duration)
				assert.WithinDuration(t, expected, token.Expiry, 5*time.Second)
				introspection := testhelpers.IntrospectToken(t, token.AccessToken, admin)
				assert.WithinDuration(t, expected, time.Unix(introspection.Get("exp").Int(), 0), 5*time.Second)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("should call token hook if configured", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				scope := "foobar"
				audience := []string{"https://api.ory.sh/"}

				hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")
					assert.Equal(t, r.Header.Get("Authorization"), "Bearer secret value")

					expectedGrantedScopes := []string{"foobar"}
					expectedGrantedAudience := []string{"https://api.ory.sh/"}

					var hookReq hydraoauth2.TokenHookRequest
					require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
					require.NotEmpty(t, hookReq.Session)
					require.Equal(t, hookReq.Session.Extra, map[string]interface{}{})
					require.NotEmpty(t, hookReq.Request)
					require.ElementsMatch(t, hookReq.Request.GrantedScopes, expectedGrantedScopes)
					require.ElementsMatch(t, hookReq.Request.GrantedAudience, expectedGrantedAudience)
					require.Equal(t, hookReq.Request.Payload, map[string][]string{
						"grant_type": {"client_credentials"},
						"scope":      {"foobar"},
					})

					claims := map[string]interface{}{
						"hooked": true,
					}

					hookResp := hydraoauth2.TokenHookResponse{
						Session: flow.AcceptOAuth2ConsentRequestSession{
							AccessToken: claims,
							IDToken:     claims,
						},
					}

					w.WriteHeader(http.StatusOK)
					require.NoError(t, json.NewEncoder(w).Encode(&hookResp))
				}))
				defer hs.Close()

				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
				reg.Config().MustSet(ctx, config.KeyTokenHook, &config.HookConfig{
					URL: hs.URL,
					Auth: &config.Auth{
						Type: "api_key",
						Config: config.AuthConfig{
							In:    "header",
							Name:  "Authorization",
							Value: "Bearer secret value",
						},
					},
				})

				defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

				cl, conf := newCustomClient(t, &hc.Client{
					Secret:        uuid.Must(uuid.NewV4()).String(),
					RedirectURIs:  []string{public.URL + "/callback"},
					ResponseTypes: []string{"token"},
					GrantTypes:    []string{"client_credentials"},
					Scope:         scope,
					Audience:      audience,
				})
				getAndInspectToken(t, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), true)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("should fail token if hook fails", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				defer hs.Close()

				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
				reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

				defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

				_, conf := newClient(t)

				_, err := getToken(t, conf)
				require.Error(t, err)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("should fail token if hook denied the request", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				}))
				defer hs.Close()

				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
				reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

				defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

				_, conf := newClient(t)

				_, err := getToken(t, conf)
				require.Error(t, err)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("should fail token if hook response is malformed", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				defer hs.Close()

				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
				reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

				defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

				_, conf := newClient(t)

				_, err := getToken(t, conf)
				require.Error(t, err)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})
}
