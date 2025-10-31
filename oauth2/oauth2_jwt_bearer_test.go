// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
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
	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/x/configx"
)

func TestJWTBearer(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue(config.KeyAccessTokenStrategy, "opaque")))
	_, admin := testhelpers.NewOAuth2Server(ctx, t, reg)

	secret := uuid.Must(uuid.NewV4()).String()
	client := &hc.Client{
		Secret:     secret,
		GrantTypes: []string{"client_credentials", "urn:ietf:params:oauth:grant-type:jwt-bearer"},
		Scope:      "offline_access",
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, client))

	newConf := func(client *hc.Client) *clientcredentials.Config {
		return &clientcredentials.Config{
			ClientID:       client.GetID(),
			ClientSecret:   secret,
			TokenURL:       reg.Config().OAuth2TokenURL(ctx).String(),
			Scopes:         strings.Split(client.Scope, " "),
			EndpointParams: url.Values{"audience": client.Audience},
		}
	}

	var getToken = func(t *testing.T, conf *clientcredentials.Config) (*goauth2.Token, error) {
		if conf.AuthStyle == goauth2.AuthStyleAutoDetect {
			conf.AuthStyle = goauth2.AuthStyleInHeader
		}
		return conf.Token(context.Background())
	}

	var inspectToken = func(t *testing.T, token *goauth2.Token, cl *hc.Client, strategy string, grant trust.Grant, checkExtraClaims bool) {
		introspection := testhelpers.IntrospectToken(t, token.AccessToken, admin)

		check := func(res gjson.Result) {
			assert.EqualValues(t, cl.GetID(), res.Get("client_id").String(), "%s", res.Raw)
			assert.EqualValues(t, grant.Subject, res.Get("sub").String(), "%s", res.Raw)
			assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), res.Get("iss").String(), "%s", res.Raw)

			assert.EqualValues(t, res.Get("nbf").Int(), res.Get("iat").Int(), "%s", res.Raw)
			assert.True(t, res.Get("exp").Int() >= res.Get("iat").Int()+int64(reg.Config().GetAccessTokenLifespan(ctx).Seconds()), "%s", res.Raw)

			assert.EqualValues(t, fmt.Sprintf(`["%s"]`, reg.Config().OAuth2TokenURL(ctx).String()), res.Get("aud").Raw, "%s", res.Raw)

			if checkExtraClaims {
				require.True(t, res.Get("ext.hooked").Bool())
			}
		}

		check(introspection)
		assert.True(t, introspection.Get("active").Bool())
		assert.EqualValues(t, "access_token", introspection.Get("token_use").String())
		assert.EqualValues(t, "Bearer", introspection.Get("token_type").String())
		assert.EqualValues(t, "offline_access", introspection.Get("scope").String(), "%s", introspection.Raw)

		if strategy != "jwt" {
			return
		}

		jwtClaims := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, token.AccessToken))
		assert.NotEmpty(t, jwtClaims.Get("jti").String())
		assert.NotEmpty(t, jwtClaims.Get("iss").String())
		assert.NotEmpty(t, jwtClaims.Get("client_id").String())
		assert.EqualValues(t, "offline_access", introspection.Get("scope").String(), "%s", introspection.Raw)

		header, err := base64.RawURLEncoding.DecodeString(strings.Split(token.AccessToken, ".")[0])
		require.NoError(t, err)
		jwtHeader := gjson.ParseBytes(header)
		assert.NotEmpty(t, jwtHeader.Get("kid").String())
		assert.EqualValues(t, "offline_access", introspection.Get("scope").String(), "%s", introspection.Raw)

		check(jwtClaims)
	}

	t.Run("case=unable to exchange invalid jwt", func(t *testing.T) {
		conf := newConf(client)
		conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {"not-a-jwt"}}
		_, err := getToken(t, conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to parse JSON Web Token")
	})

	t.Run("case=unable to request grant if not set", func(t *testing.T) {
		client := &hc.Client{
			Secret:     secret,
			GrantTypes: []string{"client_credentials"},
			Scope:      "offline_access",
		}
		require.NoError(t, reg.ClientManager().CreateClient(ctx, client))

		conf := newConf(client)
		conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {"not-a-jwt"}}
		_, err := getToken(t, conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "urn:ietf:params:oauth:grant-type:jwt-bearer")
	})

	set, kid := uuid.Must(uuid.NewV4()).String(), uuid.Must(uuid.NewV4()).String()
	keys, err := jwk.GenerateJWK(jose.RS256, kid, "sig")
	require.NoError(t, err)
	trustGrant := trust.Grant{
		ID:              uuid.Must(uuid.NewV4()),
		Issuer:          set,
		Subject:         uuid.Must(uuid.NewV4()).String(),
		AllowAnySubject: false,
		Scope:           []string{"offline_access"},
		ExpiresAt:       time.Now().Add(time.Hour),
		PublicKey:       trust.PublicKey{Set: set, KeyID: kid},
	}
	require.NoError(t, reg.GrantManager().CreateGrant(ctx, trustGrant, keys.Keys[0].Public()))
	signer := jwk.NewDefaultJWTSigner(reg, set)
	signer.GetPrivateKey = func(ctx context.Context) (interface{}, error) {
		return keys.Keys[0], nil
	}

	t.Run("case=unable to exchange token with a non-allowed subject", func(t *testing.T) {
		token, _, err := signer.Generate(ctx, jwt.MapClaims{
			"jti": uuid.Must(uuid.NewV4()).String(),
			"iss": trustGrant.Issuer,
			"sub": uuid.Must(uuid.NewV4()).String(),
			"aud": reg.Config().OAuth2TokenURL(ctx).String(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Add(-time.Minute).Unix(),
		}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
		require.NoError(t, err)

		conf := newConf(client)
		conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}
		_, err = getToken(t, conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "public key is required to check signature of JWT")
	})

	t.Run("case=unable to exchange token with non-allowed scope", func(t *testing.T) {
		token, _, err := signer.Generate(ctx, jwt.MapClaims{
			"jti": uuid.Must(uuid.NewV4()).String(),
			"iss": trustGrant.Issuer,
			"sub": trustGrant.Subject,
			"aud": reg.Config().OAuth2TokenURL(ctx).String(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Add(-time.Minute).Unix(),
		}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
		require.NoError(t, err)

		conf := newConf(client)
		conf.Scopes = []string{"i_am_not_allowed"}
		conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}
		_, err = getToken(t, conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "i_am_not_allowed")
	})

	t.Run("case=unable to exchange token with an unknown kid", func(t *testing.T) {
		token, _, err := signer.Generate(ctx, jwt.MapClaims{
			"jti": uuid.Must(uuid.NewV4()).String(),
			"iss": trustGrant.Issuer,
			"sub": trustGrant.Subject,
			"aud": reg.Config().OAuth2TokenURL(ctx).String(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Add(-time.Minute).Unix(),
		}, &jwt.Headers{Extra: map[string]interface{}{"kid": uuid.Must(uuid.NewV4()).String()}})
		require.NoError(t, err)

		conf := newConf(client)
		conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}
		_, err = getToken(t, conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "public key is required to check signature of JWT")
	})

	t.Run("case=unable to exchange token with an invalid key", func(t *testing.T) {
		keys, err := jwk.GenerateJWK(jose.RS256, kid, "sig")
		require.NoError(t, err)
		signer := jwk.NewDefaultJWTSigner(reg, set)
		signer.GetPrivateKey = func(ctx context.Context) (interface{}, error) {
			return keys.Keys[0], nil
		}

		token, _, err := signer.Generate(ctx, jwt.MapClaims{
			"jti": uuid.Must(uuid.NewV4()).String(),
			"iss": trustGrant.Issuer,
			"sub": trustGrant.Subject,
			"aud": reg.Config().OAuth2TokenURL(ctx).String(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Add(-time.Minute).Unix(),
		}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
		require.NoError(t, err)

		conf := newConf(client)
		conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}
		_, err = getToken(t, conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to verify the integrity")
	})

	t.Run("case=should exchange for an access token", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": reg.Config().OAuth2TokenURL(ctx).String(),
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				conf := newConf(client)
				conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}

				result, err := getToken(t, conf)
				require.NoError(t, err)

				inspectToken(t, result, client, strategy, trustGrant, false)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("case=exchange for an access token without client", func(t *testing.T) {
		t.Skip("This currently does not work because the client is a required foreign key and also required throughout the code base.")

		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
				reg.Config().MustSet(ctx, "config.KeyOAuth2GrantJWTClientAuthOptional", true)

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": reg.Config().OAuth2TokenURL(ctx).String(),
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				res, err := http.DefaultClient.PostForm(reg.Config().OAuth2TokenURL(ctx).String(), url.Values{
					"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
					"assertion":  {token},
				})
				require.NoError(t, err)
				defer res.Body.Close() //nolint:errcheck
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				require.EqualValues(t, http.StatusOK, res.StatusCode, "%s", body)

				var result goauth2.Token
				require.NoError(t, json.Unmarshal(body, &result))
				assert.NotEmpty(t, result.AccessToken, "%s", body)

				inspectToken(t, &result, client, strategy, trustGrant, false)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("should call token hook if configured", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				audience := reg.Config().OAuth2TokenURL(ctx).String()
				grantType := "urn:ietf:params:oauth:grant-type:jwt-bearer"

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": audience,
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

					expectedGrantedScopes := []string{client.Scope}
					expectedGrantedAudience := []string{audience}
					expectedPayload := map[string][]string{
						"assertion":  {token},
						"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
						"scope":      {"offline_access"},
					}

					var hookReq hydraoauth2.TokenHookRequest
					require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
					require.NotEmpty(t, hookReq.Session)
					require.Equal(t, hookReq.Session.Extra, map[string]interface{}{})
					require.NotEmpty(t, hookReq.Request)
					require.ElementsMatch(t, hookReq.Request.GrantedScopes, expectedGrantedScopes)
					require.ElementsMatch(t, hookReq.Request.GrantedAudience, expectedGrantedAudience)
					require.Equal(t, expectedPayload, hookReq.Request.Payload)

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
				reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

				defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

				conf := newConf(client)
				conf.EndpointParams = url.Values{"grant_type": {grantType}, "assertion": {token}}

				result, err := getToken(t, conf)
				require.NoError(t, err)

				inspectToken(t, result, client, strategy, trustGrant, true)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})

	t.Run("should call token hook if configured and omit client_secret from payload", func(t *testing.T) {
		run := func(strategy string) func(t *testing.T) {
			return func(t *testing.T) {
				audience := reg.Config().OAuth2TokenURL(ctx).String()
				grantType := "urn:ietf:params:oauth:grant-type:jwt-bearer"

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": audience,
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				client := &hc.Client{
					Secret:                  secret,
					GrantTypes:              []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
					Scope:                   "offline_access",
					TokenEndpointAuthMethod: "client_secret_post",
				}
				require.NoError(t, reg.ClientManager().CreateClient(ctx, client))

				hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

					expectedGrantedScopes := []string{client.Scope}
					expectedGrantedAudience := []string{audience}
					expectedPayload := map[string][]string{
						"assertion":  {token},
						"client_id":  {client.GetID()},
						"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
						"scope":      {"offline_access"},
					}

					var hookReq hydraoauth2.TokenHookRequest
					require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
					require.NotEmpty(t, hookReq.Session)
					require.Equal(t, hookReq.Session.Extra, map[string]interface{}{})
					require.NotEmpty(t, hookReq.Request)
					require.ElementsMatch(t, hookReq.Request.GrantedScopes, expectedGrantedScopes)
					require.ElementsMatch(t, hookReq.Request.GrantedAudience, expectedGrantedAudience)
					require.Equal(t, hookReq.Request.Payload, expectedPayload)

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
				reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

				defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

				conf := newConf(client)
				conf.AuthStyle = goauth2.AuthStyleInParams
				conf.EndpointParams = url.Values{"grant_type": {grantType}, "assertion": {token}}

				result, err := getToken(t, conf)
				require.NoError(t, err)

				inspectToken(t, result, client, strategy, trustGrant, true)
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

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": reg.Config().OAuth2TokenURL(ctx).String(),
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				conf := newConf(client)
				conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}

				_, tokenError := getToken(t, conf)
				require.Error(t, tokenError)
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

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": reg.Config().OAuth2TokenURL(ctx).String(),
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				conf := newConf(client)
				conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}

				_, tokenError := getToken(t, conf)
				require.Error(t, tokenError)
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

				token, _, err := signer.Generate(ctx, jwt.MapClaims{
					"jti": uuid.Must(uuid.NewV4()).String(),
					"iss": trustGrant.Issuer,
					"sub": trustGrant.Subject,
					"aud": reg.Config().OAuth2TokenURL(ctx).String(),
					"exp": time.Now().Add(time.Hour).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
				}, &jwt.Headers{Extra: map[string]interface{}{"kid": kid}})
				require.NoError(t, err)

				conf := newConf(client)
				conf.EndpointParams = url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {token}}

				_, tokenError := getToken(t, conf)
				require.Error(t, tokenError)
			}
		}

		t.Run("strategy=opaque", run("opaque"))
		t.Run("strategy=jwt", run("jwt"))
	})
}
