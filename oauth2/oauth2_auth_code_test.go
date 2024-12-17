// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ory/hydra/v2/jwk"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/fosite"
	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/internal/testhelpers"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/assertx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/ioutilx"
	"github.com/ory/x/josex"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/requirex"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/stringsx"
)

func noopHandler(*testing.T) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type clientCreator interface {
	CreateClient(context.Context, *client.Client) error
}

func getAuthorizeCode(t *testing.T, conf *oauth2.Config, c *http.Client, params ...oauth2.AuthCodeOption) (string, *http.Response) {
	if c == nil {
		c = testhelpers.NewEmptyJarClient(t)
	}

	state := uuid.New()
	resp, err := c.Get(conf.AuthCodeURL(state, params...))
	require.NoError(t, err)
	defer resp.Body.Close()

	q := resp.Request.URL.Query()
	require.EqualValues(t, state, q.Get("state"))
	return q.Get("code"), resp
}

func acceptLoginHandler(t *testing.T, c *client.Client, adminClient *hydra.APIClient, reg driver.Registry, subject string, checkRequestPayload func(request *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rr, _, err := adminClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
		require.NoError(t, err)

		assert.EqualValues(t, c.GetID(), pointerx.Deref(rr.Client.ClientId))
		assert.Empty(t, pointerx.Deref(rr.Client.ClientSecret))
		assert.EqualValues(t, c.GrantTypes, rr.Client.GrantTypes)
		assert.EqualValues(t, c.LogoURI, pointerx.Deref(rr.Client.LogoUri))
		assert.EqualValues(t, c.RedirectURIs, rr.Client.RedirectUris)
		assert.EqualValues(t, r.URL.Query().Get("login_challenge"), rr.Challenge)
		assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
		assert.Contains(t, rr.RequestUrl, reg.Config().OAuth2AuthURL(ctx).String())

		acceptBody := hydra.AcceptOAuth2LoginRequest{
			Subject:  subject,
			Remember: pointerx.Ptr(!rr.Skip),
			Acr:      pointerx.Ptr("1"),
			Amr:      []string{"pwd"},
			Context:  map[string]interface{}{"context": "bar"},
		}
		if checkRequestPayload != nil {
			if b := checkRequestPayload(rr); b != nil {
				acceptBody = *b
			}
		}

		v, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
			LoginChallenge(r.URL.Query().Get("login_challenge")).
			AcceptOAuth2LoginRequest(acceptBody).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)
		http.Redirect(w, r, v.RedirectTo, http.StatusFound)
	}
}

func acceptConsentHandler(t *testing.T, c *client.Client, adminClient *hydra.APIClient, reg driver.Registry, subject string, checkRequestPayload func(*hydra.OAuth2ConsentRequest)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rr, _, err := adminClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(r.URL.Query().Get("consent_challenge")).Execute()
		require.NoError(t, err)

		assert.EqualValues(t, c.GetID(), pointerx.Deref(rr.Client.ClientId))
		assert.Empty(t, pointerx.Deref(rr.Client.ClientSecret))
		assert.EqualValues(t, c.GrantTypes, rr.Client.GrantTypes)
		assert.EqualValues(t, c.LogoURI, pointerx.Deref(rr.Client.LogoUri))
		assert.EqualValues(t, c.RedirectURIs, rr.Client.RedirectUris)
		assert.EqualValues(t, subject, pointerx.Deref(rr.Subject))
		assert.EqualValues(t, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
		assert.EqualValues(t, r.URL.Query().Get("consent_challenge"), rr.Challenge)
		assert.Contains(t, *rr.RequestUrl, reg.Config().OAuth2AuthURL(r.Context()).String())
		if checkRequestPayload != nil {
			checkRequestPayload(rr)
		}

		assert.Equal(t, map[string]interface{}{"context": "bar"}, rr.Context)
		v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
			ConsentChallenge(r.URL.Query().Get("consent_challenge")).
			AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
				GrantScope: []string{"hydra", "offline", "openid"}, Remember: pointerx.Ptr(true), RememberFor: pointerx.Ptr[int64](0),
				GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IdToken:     map[string]interface{}{"bar": "baz", "email": "foo@bar.com"},
				},
			}).
			Execute()
		require.NoError(t, err)
		require.NotEmpty(t, v.RedirectTo)
		http.Redirect(w, r, v.RedirectTo, http.StatusFound)
	}
}

// TestAuthCodeWithDefaultStrategy runs proper integration tests against in-memory and database connectors, specifically
// we test:
//
// - [x] If the flow - in general - works
// - [x] If `authenticatedAt` is properly managed across the lifecycle
//   - [x] The value `authenticatedAt` should be an old time if no user interaction wrt login was required
//   - [x] The value `authenticatedAt` should be a recent time if user interaction wrt login was required
//
// - [x] If `requestedAt` is properly managed across the lifecycle
//   - [x] The value of `requestedAt` must be the initial request time, not some other time (e.g. when accepting login)
//
// - [x] If `id_token_hint` is handled properly
//   - [x] What happens if `id_token_hint` does not match the value from the handled authentication request ("accept login")
func TestAuthCodeWithDefaultStrategy(t *testing.T) {
	setupRegistries(t)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()

	for dbName, reg := range registries {
		t.Run("registry="+dbName, func(t *testing.T) {
			reg := testhelpers.NewRegistrySQLFromURL(t, reg.Config().DSN(), true, &contextx.Default{})

			require.NoError(t, jwk.EnsureAsymmetricKeypairExists(ctx, reg, string(jose.ES256), x.OpenIDConnectKeyName))
			require.NoError(t, jwk.EnsureAsymmetricKeypairExists(ctx, reg, string(jose.ES256), x.OAuth2JWTKeyName))

			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
			reg.Config().MustSet(ctx, config.KeyRefreshTokenHook, "")
			publicTS, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)

			publicClient := hydra.NewAPIClient(hydra.NewConfiguration())
			publicClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: publicTS.URL}}
			adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
			adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

			assertRefreshToken := func(t *testing.T, token *oauth2.Token, c *oauth2.Config, expectedExp time.Time) {
				introspect := testhelpers.IntrospectToken(t, c, token.RefreshToken, adminTS)
				actualExp, err := strconv.ParseInt(introspect.Get("exp").String(), 10, 64)
				require.NoError(t, err, "%s", introspect)
				requirex.EqualTime(t, expectedExp, time.Unix(actualExp, 0), time.Second*3)
			}

			assertIDToken := func(t *testing.T, token *oauth2.Token, c *oauth2.Config, expectedSubject, expectedNonce string, expectedExp time.Time) gjson.Result {
				idt, ok := token.Extra("id_token").(string)
				require.True(t, ok)
				assert.NotEmpty(t, idt)

				body, err := x.DecodeSegment(strings.Split(idt, ".")[1])
				require.NoError(t, err)

				claims := gjson.ParseBytes(body)
				assert.True(t, time.Now().After(time.Unix(claims.Get("iat").Int(), 0)), "%s", claims)
				assert.True(t, time.Now().After(time.Unix(claims.Get("nbf").Int(), 0)), "%s", claims)
				assert.True(t, time.Now().Before(time.Unix(claims.Get("exp").Int(), 0)), "%s", claims)
				if !expectedExp.IsZero() {
					// 1.5s due to rounding
					requirex.EqualTime(t, expectedExp, time.Unix(claims.Get("exp").Int(), 0), 1*time.Second+500*time.Millisecond)
				}
				assert.NotEmpty(t, claims.Get("jti").String(), "%s", claims)
				assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), claims.Get("iss").String(), "%s", claims)
				assert.NotEmpty(t, claims.Get("sid").String(), "%s", claims)
				assert.Equal(t, "1", claims.Get("acr").String(), "%s", claims)
				require.Len(t, claims.Get("amr").Array(), 1, "%s", claims)
				assert.EqualValues(t, "pwd", claims.Get("amr").Array()[0].String(), "%s", claims)

				require.Len(t, claims.Get("aud").Array(), 1, "%s", claims)
				assert.EqualValues(t, c.ClientID, claims.Get("aud").Array()[0].String(), "%s", claims)
				assert.EqualValues(t, expectedSubject, claims.Get("sub").String(), "%s", claims)
				assert.EqualValues(t, expectedNonce, claims.Get("nonce").String(), "%s", claims)
				assert.EqualValues(t, `baz`, claims.Get("bar").String(), "%s", claims)
				assert.EqualValues(t, `foo@bar.com`, claims.Get("email").String(), "%s", claims)
				assert.NotEmpty(t, claims.Get("sid").String(), "%s", claims)

				return claims
			}

			introspectAccessToken := func(t *testing.T, conf *oauth2.Config, token *oauth2.Token, expectedSubject string) gjson.Result {
				require.NotEmpty(t, token.AccessToken)
				i := testhelpers.IntrospectToken(t, conf, token.AccessToken, adminTS)
				assert.True(t, i.Get("active").Bool(), "%s", i)
				assert.EqualValues(t, conf.ClientID, i.Get("client_id").String(), "%s", i)
				assert.EqualValues(t, expectedSubject, i.Get("sub").String(), "%s", i)
				assert.EqualValues(t, `bar`, i.Get("ext.foo").String(), "%s", i)
				return i
			}

			assertJWTAccessToken := func(t *testing.T, strat string, conf *oauth2.Config, token *oauth2.Token, expectedSubject string, expectedExp time.Time, scopes string) gjson.Result {
				require.NotEmpty(t, token.AccessToken)
				parts := strings.Split(token.AccessToken, ".")
				if strat != "jwt" {
					require.Len(t, parts, 2)
					return gjson.Parse("null")
				}
				require.Len(t, parts, 3)

				body, err := x.DecodeSegment(parts[1])
				require.NoError(t, err)

				i := gjson.ParseBytes(body)
				assert.NotEmpty(t, i.Get("jti").String())
				assert.EqualValues(t, conf.ClientID, i.Get("client_id").String(), "%s", i)
				assert.EqualValues(t, expectedSubject, i.Get("sub").String(), "%s", i)
				assert.EqualValues(t, reg.Config().IssuerURL(ctx).String(), i.Get("iss").String(), "%s", i)
				assert.True(t, time.Now().After(time.Unix(i.Get("iat").Int(), 0)), "%s", i)
				assert.True(t, time.Now().After(time.Unix(i.Get("nbf").Int(), 0)), "%s", i)
				assert.True(t, time.Now().Before(time.Unix(i.Get("exp").Int(), 0)), "%s", i)
				requirex.EqualTime(t, expectedExp, time.Unix(i.Get("exp").Int(), 0), time.Second)
				assert.EqualValues(t, `bar`, i.Get("ext.foo").String(), "%s", i)
				assert.EqualValues(t, scopes, i.Get("scp").Raw, "%s", i)
				return i
			}

			waitForRefreshTokenExpiry := func() {
				time.Sleep(reg.Config().GetRefreshTokenLifespan(ctx) + time.Second)
			}

			subject := "aeneas-rekkas"
			nonce := uuid.New()

			t.Run("case=checks if request fails when audience does not match", func(t *testing.T) {
				testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
				_, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				code, _ := getAuthorizeCode(t, conf, nil, oauth2.SetAuthURLParam("audience", "https://not-ory-api/"))
				require.Empty(t, code)
			})

			t.Run("case=perform authorize code flow with ID token and refresh tokens", func(t *testing.T) {
				run := func(t *testing.T, strategy string) {
					c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
					testhelpers.NewLoginConsentUI(t, reg.Config(),
						acceptLoginHandler(t, c, adminClient, reg, subject, nil),
						acceptConsentHandler(t, c, adminClient, reg, subject, nil),
					)

					code, _ := getAuthorizeCode(t, conf, nil, oauth2.SetAuthURLParam("nonce", nonce))
					require.NotEmpty(t, code)
					token, err := conf.Exchange(context.Background(), code)
					iat := time.Now()
					require.NoError(t, err)

					assert.Empty(t, token.Extra("c_nonce_draft_00"), "should not be set if not requested")
					assert.Empty(t, token.Extra("c_nonce_expires_in_draft_00"), "should not be set if not requested")
					introspectAccessToken(t, conf, token, subject)
					assertJWTAccessToken(t, strategy, conf, token, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)
					assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
					assertRefreshToken(t, token, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))

					t.Run("followup=successfully perform refresh token flow", func(t *testing.T) {
						require.NotEmpty(t, token.RefreshToken)
						token.Expiry = token.Expiry.Add(-time.Hour * 24)
						iat = time.Now()
						refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
						require.NoError(t, err)

						require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
						require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)
						require.NotEqual(t, token.Extra("id_token"), refreshedToken.Extra("id_token"))
						introspectAccessToken(t, conf, refreshedToken, subject)

						t.Run("followup=refreshed tokens contain valid tokens", func(t *testing.T) {
							assertJWTAccessToken(t, strategy, conf, refreshedToken, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)
							assertIDToken(t, refreshedToken, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
							assertRefreshToken(t, refreshedToken, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))
						})

						t.Run("followup=original access token is no longer valid", func(t *testing.T) {
							i := testhelpers.IntrospectToken(t, conf, token.AccessToken, adminTS)
							assert.False(t, i.Get("active").Bool(), "%s", i)
						})

						t.Run("followup=original refresh token is no longer valid", func(t *testing.T) {
							_, err := conf.TokenSource(context.Background(), token).Token()
							assert.Error(t, err)
						})

						t.Run("followup=but fail subsequent refresh because expiry was reached", func(t *testing.T) {
							waitForRefreshTokenExpiry()

							// Force golang to refresh token
							refreshedToken.Expiry = refreshedToken.Expiry.Add(-time.Hour * 24)
							_, err := conf.TokenSource(context.Background(), refreshedToken).Token()
							require.Error(t, err)
						})
					})
				}

				t.Run("strategy=jwt", func(t *testing.T) {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
					run(t, "jwt")
				})

				t.Run("strategy=opaque", func(t *testing.T) {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
					run(t, "opaque")
				})
			})

			t.Run("case=perform authorize code flow with verifable credentials", func(t *testing.T) {
				// Make sure we test against all crypto suites that we advertise.
				cfg, _, err := publicClient.OidcAPI.DiscoverOidcConfiguration(ctx).Execute()
				require.NoError(t, err)
				supportedCryptoSuites := cfg.CredentialsSupportedDraft00[0].CryptographicSuitesSupported

				run := func(t *testing.T, strategy string) {
					_, conf := newOAuth2Client(
						t,
						reg,
						testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler),
						withScope("openid userinfo_credential_draft_00"),
					)
					testhelpers.NewLoginConsentUI(t, reg.Config(),
						func(w http.ResponseWriter, r *http.Request) {
							acceptBody := hydra.AcceptOAuth2LoginRequest{
								Subject: subject,
								Acr:     pointerx.Ptr("1"),
								Amr:     []string{"pwd"},
								Context: map[string]interface{}{"context": "bar"},
							}
							v, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
								LoginChallenge(r.URL.Query().Get("login_challenge")).
								AcceptOAuth2LoginRequest(acceptBody).
								Execute()
							require.NoError(t, err)
							require.NotEmpty(t, v.RedirectTo)
							http.Redirect(w, r, v.RedirectTo, http.StatusFound)
						},
						func(w http.ResponseWriter, r *http.Request) {
							rr, _, err := adminClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(r.URL.Query().Get("consent_challenge")).Execute()
							require.NoError(t, err)

							assert.Equal(t, map[string]interface{}{"context": "bar"}, rr.Context)
							v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
								ConsentChallenge(r.URL.Query().Get("consent_challenge")).
								AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
									GrantScope:               []string{"openid", "userinfo_credential_draft_00"},
									GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
									Session: &hydra.AcceptOAuth2ConsentRequestSession{
										AccessToken: map[string]interface{}{"foo": "bar"},
										IdToken:     map[string]interface{}{"email": "foo@bar.com", "bar": "baz"},
									},
								}).
								Execute()
							require.NoError(t, err)
							require.NotEmpty(t, v.RedirectTo)
							http.Redirect(w, r, v.RedirectTo, http.StatusFound)
						},
					)

					code, _ := getAuthorizeCode(t, conf, nil,
						oauth2.SetAuthURLParam("nonce", nonce),
						oauth2.SetAuthURLParam("scope", "openid userinfo_credential_draft_00"),
					)
					require.NotEmpty(t, code)
					token, err := conf.Exchange(context.Background(), code)
					require.NoError(t, err)
					iat := time.Now()

					vcNonce := token.Extra("c_nonce_draft_00").(string)
					assert.NotEmpty(t, vcNonce)
					expiry := token.Extra("c_nonce_expires_in_draft_00")
					assert.NotEmpty(t, expiry)
					assert.NoError(t, reg.Persister().IsNonceValid(ctx, token.AccessToken, vcNonce))

					t.Run("followup=successfully create a verifiable credential", func(t *testing.T) {
						t.Parallel()

						for _, alg := range supportedCryptoSuites {
							alg := alg
							t.Run(fmt.Sprintf("alg=%s", alg), func(t *testing.T) {
								t.Parallel()
								assertCreateVerifiableCredential(t, reg, vcNonce, token, jose.SignatureAlgorithm(alg))
							})
						}
					})

					t.Run("followup=get new nonce from priming request", func(t *testing.T) {
						t.Parallel()
						// Assert that we can fetch a verifiable credential with the nonce.
						res, err := doPrimingRequest(t, reg, token, &hydraoauth2.CreateVerifiableCredentialRequestBody{
							Format: "jwt_vc_json",
							Types:  []string{"VerifiableCredential", "UserInfoCredential"},
						})
						assert.NoError(t, err)

						t.Run("followup=successfully create a verifiable credential from fresh nonce", func(t *testing.T) {
							assertCreateVerifiableCredential(t, reg, res.Nonce, token, jose.ES256)
						})
					})

					t.Run("followup=rejects proof signed by another key", func(t *testing.T) {
						t.Parallel()
						for _, tc := range []struct {
							name      string
							format    string
							proofType string
							proof     func() string
						}{
							{
								name: "proof=mismatching keys",
								proof: func() string {
									// Create mismatching public and private keys.
									pubKey, _, err := josex.NewSigningKey(jose.ES256, 0)
									require.NoError(t, err)
									_, privKey, err := josex.NewSigningKey(jose.ES256, 0)
									require.NoError(t, err)
									pubKeyJWK := &jose.JSONWebKey{Key: pubKey, Algorithm: string(jose.ES256)}
									return createVCProofJWT(t, pubKeyJWK, privKey, vcNonce)
								},
							},
							{
								name:   "proof=invalid format",
								format: "invalid_format",
								proof: func() string {
									// Create mismatching public and private keys.
									pubKey, privKey, err := josex.NewSigningKey(jose.ES256, 0)
									require.NoError(t, err)
									pubKeyJWK := &jose.JSONWebKey{Key: pubKey, Algorithm: string(jose.ES256)}
									return createVCProofJWT(t, pubKeyJWK, privKey, vcNonce)
								},
							},
							{
								name:      "proof=invalid type",
								proofType: "invalid",
								proof: func() string {
									// Create mismatching public and private keys.
									pubKey, privKey, err := josex.NewSigningKey(jose.ES256, 0)
									require.NoError(t, err)
									pubKeyJWK := &jose.JSONWebKey{Key: pubKey, Algorithm: string(jose.ES256)}
									return createVCProofJWT(t, pubKeyJWK, privKey, vcNonce)
								},
							},
							{
								name: "proof=invalid nonce",
								proof: func() string {
									// Create mismatching public and private keys.
									pubKey, privKey, err := josex.NewSigningKey(jose.ES256, 0)
									require.NoError(t, err)
									pubKeyJWK := &jose.JSONWebKey{Key: pubKey, Algorithm: string(jose.ES256)}
									return createVCProofJWT(t, pubKeyJWK, privKey, "invalid nonce")
								},
							},
						} {
							tc := tc
							t.Run(tc.name, func(t *testing.T) {
								t.Parallel()
								_, err := createVerifiableCredential(t, reg, token, &hydraoauth2.CreateVerifiableCredentialRequestBody{
									Format: stringsx.Coalesce(tc.format, "jwt_vc_json"),
									Types:  []string{"VerifiableCredential", "UserInfoCredential"},
									Proof: &hydraoauth2.VerifiableCredentialProof{
										ProofType: stringsx.Coalesce(tc.proofType, "jwt"),
										JWT:       tc.proof(),
									},
								})
								require.Error(t, err)
								assert.Equal(t, "invalid_request", err.Error())
							})
						}

					})

					t.Run("followup=access token and id token are valid", func(t *testing.T) {
						assertJWTAccessToken(t, strategy, conf, token, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["openid","userinfo_credential_draft_00"]`)
						assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
					})
				}

				t.Run("strategy=jwt", func(t *testing.T) {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
					run(t, "jwt")
				})

				t.Run("strategy=opaque", func(t *testing.T) {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
					run(t, "opaque")
				})
			})

			t.Run("suite=invalid query params", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				otherClient, _ := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, nil),
					acceptConsentHandler(t, c, adminClient, reg, subject, nil),
				)

				withWrongClientAfterLogin := &http.Client{
					Jar: testhelpers.NewEmptyCookieJar(t),
					CheckRedirect: func(req *http.Request, _ []*http.Request) error {
						if req.URL.Path != "/oauth2/auth" {
							return nil
						}
						q := req.URL.Query()
						if !q.Has("login_verifier") {
							return nil
						}
						q.Set("client_id", otherClient.GetID())
						req.URL.RawQuery = q.Encode()
						return nil
					},
				}
				withWrongClientAfterConsent := &http.Client{
					Jar: testhelpers.NewEmptyCookieJar(t),
					CheckRedirect: func(req *http.Request, _ []*http.Request) error {
						if req.URL.Path != "/oauth2/auth" {
							return nil
						}
						q := req.URL.Query()
						if !q.Has("consent_verifier") {
							return nil
						}
						q.Set("client_id", otherClient.GetID())
						req.URL.RawQuery = q.Encode()
						return nil
					},
				}

				withWrongScopeAfterLogin := &http.Client{
					Jar: testhelpers.NewEmptyCookieJar(t),
					CheckRedirect: func(req *http.Request, _ []*http.Request) error {
						if req.URL.Path != "/oauth2/auth" {
							return nil
						}
						q := req.URL.Query()
						if !q.Has("login_verifier") {
							return nil
						}
						q.Set("scope", "invalid scope")
						req.URL.RawQuery = q.Encode()
						return nil
					},
				}

				withWrongScopeAfterConsent := &http.Client{
					Jar: testhelpers.NewEmptyCookieJar(t),
					CheckRedirect: func(req *http.Request, _ []*http.Request) error {
						if req.URL.Path != "/oauth2/auth" {
							return nil
						}
						q := req.URL.Query()
						if !q.Has("consent_verifier") {
							return nil
						}
						q.Set("scope", "invalid scope")
						req.URL.RawQuery = q.Encode()
						return nil
					},
				}
				for _, tc := range []struct {
					name             string
					client           *http.Client
					expectedResponse string
				}{{
					name:             "fails with wrong client ID after login",
					client:           withWrongClientAfterLogin,
					expectedResponse: "invalid_client",
				}, {
					name:             "fails with wrong client ID after consent",
					client:           withWrongClientAfterConsent,
					expectedResponse: "invalid_client",
				}, {
					name:             "fails with wrong scopes after login",
					client:           withWrongScopeAfterLogin,
					expectedResponse: "invalid_scope",
				}, {
					name:             "fails with wrong scopes after consent",
					client:           withWrongScopeAfterConsent,
					expectedResponse: "invalid_scope",
				}} {
					t.Run("case="+tc.name, func(t *testing.T) {
						state := uuid.New()
						resp, err := tc.client.Get(conf.AuthCodeURL(state))
						require.NoError(t, err)
						assert.Equal(t, tc.expectedResponse, resp.Request.URL.Query().Get("error"), "%s", resp.Request.URL.RawQuery)
						resp.Body.Close()
					})
				}
			})

			t.Run("case=checks if request fails when subject is empty", func(t *testing.T) {
				testhelpers.NewLoginConsentUI(t, reg.Config(), func(w http.ResponseWriter, r *http.Request) {
					_, res, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(ctx).
						LoginChallenge(r.URL.Query().Get("login_challenge")).
						AcceptOAuth2LoginRequest(hydra.AcceptOAuth2LoginRequest{Subject: "", Remember: pointerx.Ptr(true)}).Execute()
					require.Error(t, err) // expects 400
					body := string(ioutilx.MustReadAll(res.Body))
					assert.Contains(t, body, "Field 'subject' must not be empty", "%s", body)
				}, testhelpers.HTTPServerNoExpectedCallHandler(t))
				_, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))

				_, err := testhelpers.NewEmptyJarClient(t).Get(conf.AuthCodeURL(uuid.New()))
				require.NoError(t, err)
			})

			t.Run("case=perform flow with prompt=registration", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))

				regUI := httptest.NewServer(acceptLoginHandler(t, c, adminClient, reg, subject, nil))
				t.Cleanup(regUI.Close)
				reg.Config().MustSet(ctx, config.KeyRegistrationURL, regUI.URL)

				testhelpers.NewLoginConsentUI(t, reg.Config(),
					nil,
					acceptConsentHandler(t, c, adminClient, reg, subject, nil))

				code, _ := getAuthorizeCode(t, conf, nil,
					oauth2.SetAuthURLParam("prompt", "registration"),
					oauth2.SetAuthURLParam("nonce", nonce))
				require.NotEmpty(t, code)

				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)

				assertIDToken(t, token, conf, subject, nonce, time.Now().Add(reg.Config().GetIDTokenLifespan(ctx)))
			})

			t.Run("case=perform flow with audience", func(t *testing.T) {
				expectAud := "https://api.ory.sh/"
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
						assert.False(t, r.Skip)
						assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
						return nil
					}),
					acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
						assert.False(t, *r.Skip)
						assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
					}))

				code, _ := getAuthorizeCode(t, conf, nil,
					oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"),
					oauth2.SetAuthURLParam("nonce", nonce))
				require.NotEmpty(t, code)

				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)

				claims := introspectAccessToken(t, conf, token, subject)
				aud := claims.Get("aud").Array()
				require.Len(t, aud, 1)
				assert.EqualValues(t, aud[0].String(), expectAud)

				assertIDToken(t, token, conf, subject, nonce, time.Now().Add(reg.Config().GetIDTokenLifespan(ctx)))
			})

			t.Run("case=respects client token lifespan configuration", func(t *testing.T) {
				run := func(t *testing.T, strategy string, c *client.Client, conf *oauth2.Config, expectedLifespans client.Lifespans) {
					testhelpers.NewLoginConsentUI(t, reg.Config(),
						acceptLoginHandler(t, c, adminClient, reg, subject, nil),
						acceptConsentHandler(t, c, adminClient, reg, subject, nil),
					)

					code, _ := getAuthorizeCode(t, conf, nil, oauth2.SetAuthURLParam("nonce", nonce))
					require.NotEmpty(t, code)
					token, err := conf.Exchange(context.Background(), code)
					iat := time.Now()
					require.NoError(t, err)

					body := introspectAccessToken(t, conf, token, subject)
					requirex.EqualTime(t, iat.Add(expectedLifespans.AuthorizationCodeGrantAccessTokenLifespan.Duration), time.Unix(body.Get("exp").Int(), 0), time.Second)

					assertJWTAccessToken(t, strategy, conf, token, subject, iat.Add(expectedLifespans.AuthorizationCodeGrantAccessTokenLifespan.Duration), `["hydra","offline","openid"]`)
					assertIDToken(t, token, conf, subject, nonce, iat.Add(expectedLifespans.AuthorizationCodeGrantIDTokenLifespan.Duration))
					assertRefreshToken(t, token, conf, iat.Add(expectedLifespans.AuthorizationCodeGrantRefreshTokenLifespan.Duration))

					t.Run("followup=successfully perform refresh token flow", func(t *testing.T) {
						require.NotEmpty(t, token.RefreshToken)
						token.Expiry = token.Expiry.Add(-time.Hour * 24)
						refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
						iat = time.Now()
						require.NoError(t, err)
						assertRefreshToken(t, refreshedToken, conf, iat.Add(expectedLifespans.RefreshTokenGrantRefreshTokenLifespan.Duration))
						assertJWTAccessToken(t, strategy, conf, refreshedToken, subject, iat.Add(expectedLifespans.RefreshTokenGrantAccessTokenLifespan.Duration), `["hydra","offline","openid"]`)
						assertIDToken(t, refreshedToken, conf, subject, nonce, iat.Add(expectedLifespans.RefreshTokenGrantIDTokenLifespan.Duration))

						require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
						require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)
						require.NotEqual(t, token.Extra("id_token"), refreshedToken.Extra("id_token"))

						body := introspectAccessToken(t, conf, refreshedToken, subject)
						requirex.EqualTime(t, iat.Add(expectedLifespans.RefreshTokenGrantAccessTokenLifespan.Duration), time.Unix(body.Get("exp").Int(), 0), time.Second)

						t.Run("followup=original access token is no longer valid", func(t *testing.T) {
							i := testhelpers.IntrospectToken(t, conf, token.AccessToken, adminTS)
							assert.False(t, i.Get("active").Bool(), "%s", i)
						})

						t.Run("followup=original refresh token is no longer valid", func(t *testing.T) {
							_, err := conf.TokenSource(context.Background(), token).Token()
							assert.Error(t, err)
						})
					})
				}

				t.Run("case=custom-lifespans-active-jwt", func(t *testing.T) {
					c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
					ls := testhelpers.TestLifespans
					ls.AuthorizationCodeGrantAccessTokenLifespan = x.NullDuration{Valid: true, Duration: 6 * time.Second}
					testhelpers.UpdateClientTokenLifespans(
						t,
						&oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret},
						c.GetID(),
						ls, adminTS,
					)
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
					run(t, "jwt", c, conf, ls)
				})

				t.Run("case=custom-lifespans-active-opaque", func(t *testing.T) {
					c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
					ls := testhelpers.TestLifespans
					ls.AuthorizationCodeGrantAccessTokenLifespan = x.NullDuration{Valid: true, Duration: 6 * time.Second}
					testhelpers.UpdateClientTokenLifespans(
						t,
						&oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret},
						c.GetID(),
						ls, adminTS,
					)
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
					run(t, "opaque", c, conf, ls)
				})

				t.Run("case=custom-lifespans-unset", func(t *testing.T) {
					c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
					testhelpers.UpdateClientTokenLifespans(t, &oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret}, c.GetID(), testhelpers.TestLifespans, adminTS)
					testhelpers.UpdateClientTokenLifespans(t, &oauth2.Config{ClientID: c.GetID(), ClientSecret: conf.ClientSecret}, c.GetID(), client.Lifespans{}, adminTS)
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")

					//goland:noinspection GoDeprecation
					expectedLifespans := client.Lifespans{
						AuthorizationCodeGrantAccessTokenLifespan:  x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
						AuthorizationCodeGrantIDTokenLifespan:      x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
						AuthorizationCodeGrantRefreshTokenLifespan: x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
						ClientCredentialsGrantAccessTokenLifespan:  x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
						ImplicitGrantAccessTokenLifespan:           x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
						ImplicitGrantIDTokenLifespan:               x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
						JwtBearerGrantAccessTokenLifespan:          x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
						PasswordGrantAccessTokenLifespan:           x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
						PasswordGrantRefreshTokenLifespan:          x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
						RefreshTokenGrantIDTokenLifespan:           x.NullDuration{Valid: true, Duration: reg.Config().GetIDTokenLifespan(ctx)},
						RefreshTokenGrantAccessTokenLifespan:       x.NullDuration{Valid: true, Duration: reg.Config().GetAccessTokenLifespan(ctx)},
						RefreshTokenGrantRefreshTokenLifespan:      x.NullDuration{Valid: true, Duration: reg.Config().GetRefreshTokenLifespan(ctx)},
					}
					run(t, "opaque", c, conf, expectedLifespans)
				})
			})

			t.Run("case=use remember feature and prompt=none", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, nil),
					acceptConsentHandler(t, c, adminClient, reg, subject, nil),
				)

				oc := testhelpers.NewEmptyJarClient(t)
				code, _ := getAuthorizeCode(t, conf, oc,
					oauth2.SetAuthURLParam("nonce", nonce),
					oauth2.SetAuthURLParam("prompt", "login consent"),
					oauth2.SetAuthURLParam("max_age", "1"),
				)
				require.NotEmpty(t, code)
				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)
				introspectAccessToken(t, conf, token, subject)

				// Reset UI to check for skip values
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
						require.True(t, r.Skip)
						require.EqualValues(t, subject, r.Subject)
						return nil
					}),
					acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
						require.True(t, *r.Skip)
						require.EqualValues(t, subject, *r.Subject)
					}),
				)

				t.Run("followup=checks if authenticatedAt/requestedAt is properly forwarded across the lifecycle by checking if prompt=none works", func(t *testing.T) {
					// In order to check if authenticatedAt/requestedAt works, we'll sleep first in order to ensure that authenticatedAt is in the past
					// if handled correctly.
					time.Sleep(time.Second + time.Nanosecond)

					code, _ := getAuthorizeCode(t, conf, oc,
						oauth2.SetAuthURLParam("nonce", nonce),
						oauth2.SetAuthURLParam("prompt", "none"),
						oauth2.SetAuthURLParam("max_age", "60"),
					)
					require.NotEmpty(t, code)
					token, err := conf.Exchange(context.Background(), code)
					require.NoError(t, err)
					original := introspectAccessToken(t, conf, token, subject)

					t.Run("followup=run the flow three more times", func(t *testing.T) {
						for i := 0; i < 3; i++ {
							t.Run(fmt.Sprintf("run=%d", i), func(t *testing.T) {
								code, _ := getAuthorizeCode(t, conf, oc,
									oauth2.SetAuthURLParam("nonce", nonce),
									oauth2.SetAuthURLParam("prompt", "none"),
									oauth2.SetAuthURLParam("max_age", "60"),
								)
								require.NotEmpty(t, code)
								token, err := conf.Exchange(context.Background(), code)
								require.NoError(t, err)
								followup := introspectAccessToken(t, conf, token, subject)
								assert.Equal(t, original.Get("auth_time").Int(), followup.Get("auth_time").Int())
							})
						}
					})

					t.Run("followup=fails when max age is reached and prompt is none", func(t *testing.T) {
						code, _ := getAuthorizeCode(t, conf, oc,
							oauth2.SetAuthURLParam("nonce", nonce),
							oauth2.SetAuthURLParam("prompt", "none"),
							oauth2.SetAuthURLParam("max_age", "1"),
						)
						require.Empty(t, code)
					})

					t.Run("followup=passes and resets skip when prompt=login", func(t *testing.T) {
						testhelpers.NewLoginConsentUI(t, reg.Config(),
							acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
								require.False(t, r.Skip)
								require.Empty(t, r.Subject)
								return nil
							}),
							acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
								require.True(t, *r.Skip)
								require.EqualValues(t, subject, *r.Subject)
							}),
						)
						code, _ := getAuthorizeCode(t, conf, oc,
							oauth2.SetAuthURLParam("nonce", nonce),
							oauth2.SetAuthURLParam("prompt", "login"),
							oauth2.SetAuthURLParam("max_age", "1"),
						)
						require.NotEmpty(t, code)
						token, err := conf.Exchange(context.Background(), code)
						require.NoError(t, err)
						introspectAccessToken(t, conf, token, subject)
						assertIDToken(t, token, conf, subject, nonce, time.Now().Add(reg.Config().GetIDTokenLifespan(ctx)))
					})
				})
			})

			t.Run("case=should fail if prompt=none but no auth session given", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, nil),
					acceptConsentHandler(t, c, adminClient, reg, subject, nil),
				)

				oc := testhelpers.NewEmptyJarClient(t)
				code, _ := getAuthorizeCode(t, conf, oc,
					oauth2.SetAuthURLParam("prompt", "none"),
				)
				require.Empty(t, code)
			})

			t.Run("case=requires re-authentication when id_token_hint is set to a user 'patrik-neu' but the session is 'aeneas-rekkas' and then fails because the user id from the log in endpoint is 'aeneas-rekkas'", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
						require.False(t, r.Skip)
						require.Empty(t, r.Subject)
						return nil
					}),
					acceptConsentHandler(t, c, adminClient, reg, subject, nil),
				)

				oc := testhelpers.NewEmptyJarClient(t)

				// Create login session for aeneas-rekkas
				code, _ := getAuthorizeCode(t, conf, oc)
				require.NotEmpty(t, code)

				// Perform authentication for aeneas-rekkas which fails because id_token_hint is patrik-neu
				code, _ = getAuthorizeCode(t, conf, oc,
					oauth2.SetAuthURLParam("id_token_hint", testhelpers.NewIDToken(t, reg, "patrik-neu")),
				)
				require.Empty(t, code)
			})

			t.Run("case=should not cause issues if max_age is very low and consent takes a long time", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
						time.Sleep(time.Second * 2)
						return nil
					}),
					acceptConsentHandler(t, c, adminClient, reg, subject, nil),
				)

				code, _ := getAuthorizeCode(t, conf, nil)
				require.NotEmpty(t, code)
			})

			t.Run("case=ensure consistent claims returned for userinfo", func(t *testing.T) {
				c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					acceptLoginHandler(t, c, adminClient, reg, subject, nil),
					acceptConsentHandler(t, c, adminClient, reg, subject, nil),
				)

				code, _ := getAuthorizeCode(t, conf, nil)
				require.NotEmpty(t, code)

				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)

				idClaims := assertIDToken(t, token, conf, subject, "", time.Now().Add(reg.Config().GetIDTokenLifespan(ctx)))

				uiClaims := testhelpers.Userinfo(t, token, publicTS)

				for _, f := range []string{
					"sub",
					"iss",
					"aud",
					"bar",
					"auth_time",
				} {
					assert.NotEmpty(t, uiClaims.Get(f).Raw, "%s: %s", f, uiClaims)
					assert.EqualValues(t, idClaims.Get(f).Raw, uiClaims.Get(f).Raw, "%s\nuserinfo: %s\nidtoken: %s", f, uiClaims, idClaims)
				}

				for _, f := range []string{
					"at_hash",
					"c_hash",
					"nonce",
					"sid",
					"jti",
				} {
					assert.Empty(t, uiClaims.Get(f).Raw, "%s: %s", f, uiClaims)
				}
			})

			t.Run("case=add ext claims from hook if configured", func(t *testing.T) {
				run := func(strategy string) func(t *testing.T) {
					return func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")
							assert.Equal(t, r.Header.Get("Authorization"), "Bearer secret value")

							var hookReq hydraoauth2.TokenHookRequest
							require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
							require.NotEmpty(t, hookReq.Session)
							require.Equal(t, map[string]interface{}{"foo": "bar"}, hookReq.Session.Extra)
							require.NotEmpty(t, hookReq.Request)
							require.ElementsMatch(t, []string{}, hookReq.Request.GrantedAudience)
							require.Equal(t, map[string][]string{"grant_type": {"authorization_code"}}, hookReq.Request.Payload)

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

						t.Cleanup(func() {
							reg.Config().Delete(ctx, config.KeyTokenHook)
						})

						expectAud := "https://api.ory.sh/"
						c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
						testhelpers.NewLoginConsentUI(t, reg.Config(),
							acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
								assert.False(t, r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
								return nil
							}),
							acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
								assert.False(t, *r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
							}))

						code, _ := getAuthorizeCode(t, conf, nil,
							oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"),
							oauth2.SetAuthURLParam("nonce", nonce))
						require.NotEmpty(t, code)

						token, err := conf.Exchange(context.Background(), code)
						require.NoError(t, err)

						assertJWTAccessToken(t, strategy, conf, token, subject, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)

						// NOTE: using introspect to cover both jwt and opaque strategies
						accessTokenClaims := introspectAccessToken(t, conf, token, subject)
						require.True(t, accessTokenClaims.Get("ext.hooked").Bool())

						idTokenClaims := assertIDToken(t, token, conf, subject, nonce, time.Now().Add(reg.Config().GetIDTokenLifespan(ctx)))
						require.True(t, idTokenClaims.Get("hooked").Bool())
					}
				}

				t.Run("strategy=opaque", run("opaque"))
				t.Run("strategy=jwt", run("jwt"))
			})

			t.Run("case=fail token exchange if hook fails", func(t *testing.T) {
				run := func(strategy string) func(t *testing.T) {
					return func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusInternalServerError)
						}))
						defer hs.Close()

						reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
						reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

						defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

						expectAud := "https://api.ory.sh/"
						c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
						testhelpers.NewLoginConsentUI(t, reg.Config(),
							acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
								assert.False(t, r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
								return nil
							}),
							acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
								assert.False(t, *r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
							}))

						code, _ := getAuthorizeCode(t, conf, nil,
							oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"),
							oauth2.SetAuthURLParam("nonce", nonce))
						require.NotEmpty(t, code)

						_, err := conf.Exchange(context.Background(), code)
						require.Error(t, err)
					}
				}

				t.Run("strategy=opaque", run("opaque"))
				t.Run("strategy=jwt", run("jwt"))
			})

			t.Run("case=fail token exchange if hook denies the request", func(t *testing.T) {
				run := func(strategy string) func(t *testing.T) {
					return func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusForbidden)
						}))
						defer hs.Close()

						reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
						reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

						defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

						expectAud := "https://api.ory.sh/"
						c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
						testhelpers.NewLoginConsentUI(t, reg.Config(),
							acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
								assert.False(t, r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
								return nil
							}),
							acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
								assert.False(t, *r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
							}))

						code, _ := getAuthorizeCode(t, conf, nil,
							oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"),
							oauth2.SetAuthURLParam("nonce", nonce))
						require.NotEmpty(t, code)

						_, err := conf.Exchange(context.Background(), code)
						require.Error(t, err)
					}
				}

				t.Run("strategy=opaque", run("opaque"))
				t.Run("strategy=jwt", run("jwt"))
			})

			t.Run("case=fail token exchange if hook response is malformed", func(t *testing.T) {
				run := func(strategy string) func(t *testing.T) {
					return func(t *testing.T) {
						hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
						}))
						defer hs.Close()

						reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
						reg.Config().MustSet(ctx, config.KeyTokenHook, hs.URL)

						defer reg.Config().MustSet(ctx, config.KeyTokenHook, nil)

						expectAud := "https://api.ory.sh/"
						c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
						testhelpers.NewLoginConsentUI(t, reg.Config(),
							acceptLoginHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest {
								assert.False(t, r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
								return nil
							}),
							acceptConsentHandler(t, c, adminClient, reg, subject, func(r *hydra.OAuth2ConsentRequest) {
								assert.False(t, *r.Skip)
								assert.EqualValues(t, []string{expectAud}, r.RequestedAccessTokenAudience)
							}))

						code, _ := getAuthorizeCode(t, conf, nil,
							oauth2.SetAuthURLParam("audience", "https://api.ory.sh/"),
							oauth2.SetAuthURLParam("nonce", nonce))
						require.NotEmpty(t, code)

						_, err := conf.Exchange(context.Background(), code)
						require.Error(t, err)
					}
				}

				t.Run("strategy=opaque", run("opaque"))
				t.Run("strategy=jwt", run("jwt"))
			})

			t.Run("case=graceful token rotation", func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "2s")
				reg.Config().Delete(ctx, config.KeyTokenHook)
				reg.Config().Delete(ctx, config.KeyRefreshTokenHook)
				reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1m")
				reg.Config().MustSet(ctx, config.KeyAccessTokenLifespan, "1m")
				t.Cleanup(func() {
					reg.Config().Delete(ctx, config.KeyRefreshTokenRotationGracePeriod)
					reg.Config().Delete(ctx, config.KeyRefreshTokenLifespan)
					reg.Config().Delete(ctx, config.KeyAccessTokenLifespan)
				})

				// This is an essential and complex test suite. We need to cover the following cases:
				//
				// * Graceful refresh token rotation invalidates the previous access token.
				// * An expired refresh token cannot be used even if grace period is active.
				// * A used refresh token cannot be re-used once the grace period ends, and it triggers re-use detection.
				// * A test suite with a variety of concurrent refresh token chains.
				run := func(t *testing.T, strategy string) {
					c, conf := newOAuth2Client(t, reg, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
					testhelpers.NewLoginConsentUI(t, reg.Config(),
						acceptLoginHandler(t, c, adminClient, reg, subject, nil),
						acceptConsentHandler(t, c, adminClient, reg, subject, nil),
					)

					issueTokens := func(t *testing.T) *oauth2.Token {
						code, _ := getAuthorizeCode(t, conf, nil, oauth2.SetAuthURLParam("nonce", nonce))
						require.NotEmpty(t, code)
						token, err := conf.Exchange(context.Background(), code)
						iat := time.Now()
						require.NoError(t, err)

						introspectAccessToken(t, conf, token, subject)
						assertJWTAccessToken(t, strategy, conf, token, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)
						assertIDToken(t, token, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
						assertRefreshToken(t, token, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))
						return token
					}

					refreshTokens := func(t *testing.T, token *oauth2.Token) *oauth2.Token {
						require.NotEmpty(t, token.RefreshToken)
						token.Expiry = time.Now().Add(-time.Hour * 24)
						iat := time.Now()
						refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
						require.NoError(t, err)

						require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)
						require.NotEqual(t, token.RefreshToken, refreshedToken.RefreshToken)
						require.NotEqual(t, token.Extra("id_token"), refreshedToken.Extra("id_token"))

						introspectAccessToken(t, conf, refreshedToken, subject)
						assertJWTAccessToken(t, strategy, conf, refreshedToken, subject, iat.Add(reg.Config().GetAccessTokenLifespan(ctx)), `["hydra","offline","openid"]`)
						assertIDToken(t, refreshedToken, conf, subject, nonce, iat.Add(reg.Config().GetIDTokenLifespan(ctx)))
						assertRefreshToken(t, refreshedToken, conf, iat.Add(reg.Config().GetRefreshTokenLifespan(ctx)))
						return refreshedToken
					}

					assertInactive := func(t *testing.T, token string, c *oauth2.Config) {
						t.Helper()
						at := testhelpers.IntrospectToken(t, conf, token, adminTS)
						assert.False(t, at.Get("active").Bool(), "%s", at)
					}

					t.Run("gracefully refreshing a token does invalidate the previous access token", func(t *testing.T) {
						reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "2s")
						reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1m")

						token := issueTokens(t)
						_ = refreshTokens(t, token)

						assertInactive(t, token.AccessToken, conf) // Original access token is invalid

						_ = refreshTokens(t, token)
						assertInactive(t, token.AccessToken, conf) // Original access token is still invalid
					})

					t.Run("an expired refresh token can not be used even if we are in the grace period", func(t *testing.T) {
						reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "5s")
						reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1s")

						token := issueTokens(t)
						time.Sleep(time.Second * 2) // Let token expire - we need 2 seconds to reliably be longer than TTL

						token.Expiry = time.Now().Add(-time.Hour * 24)
						_, err := conf.TokenSource(ctx, token).Token()
						require.Error(t, err, "Rotating an expired token is not possible even when we are in the grace period")

						// The access token is still valid because using an expired refresh token has no effect on the access token.
						assertInactive(t, token.RefreshToken, conf)
					})

					t.Run("a used refresh token can not be re-used once the grace period ends and it triggers re-use detection", func(t *testing.T) {
						reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")
						reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1m")

						token := issueTokens(t)
						refreshed := refreshTokens(t, token)

						time.Sleep(time.Second * 2) // Wait for the grace period to end

						token.Expiry = time.Now().Add(-time.Hour * 24)
						_, err := conf.TokenSource(ctx, token).Token()
						require.Error(t, err, "Rotating a used refresh token is not possible after the grace period")

						assertInactive(t, token.AccessToken, conf)
						assertInactive(t, token.RefreshToken, conf)

						assertInactive(t, refreshed.AccessToken, conf)
						assertInactive(t, refreshed.RefreshToken, conf)
					})

					// This test suite covers complex scenarios where we have multiple generations of tokens and we need to ensure
					// that key security mitigations are in place:
					//
					// - Token re-use detection clears all tokens if a refresh token is re-used after the grace period.
					// - Revoking consent clears all tokens.
					// - Token revokation clears all tokens.
					//
					// The test creates 4 token generations, where each generations has twice as many tokens as the previous generation.
					// The generations are created like this:
					//
					// - In the first scenario, all token generations are created at the same time.
					// - In the second scenario, we create token generations with a delay that is longer than the grace period between them.
					//
					// Tokens for each generation are created in parallel to ensure we have no state leak anywhere.0
					t.Run("token generations", func(t *testing.T) {

						gracePeriod := time.Second
						aboveGracePeriod := time.Second * 2
						reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1m")
						reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, gracePeriod.String())
						reg.Config().Delete(ctx, config.KeyTokenHook)
						reg.Config().Delete(ctx, config.KeyRefreshTokenHook)

						createTokenGenerations := func(t *testing.T, count int, withSleep time.Duration) [][]*oauth2.Token {
							generations := make([][]*oauth2.Token, count)
							generations[0] = []*oauth2.Token{issueTokens(t)}
							// Start from the first generation. For every next generation, we refresh all the tokens of the previous generation twice.
							for i := 1; i < len(generations); i++ {
								generations[i] = make([]*oauth2.Token, 0, len(generations[i-1])*2)

								var wg sync.WaitGroup
								gen := func(i int, token *oauth2.Token) {
									defer wg.Done()
									generations[i] = append(generations[i], refreshTokens(t, token))
								}

								for _, token := range generations[i-1] {
									wg.Add(2)
									if dbName != "cockroach" {
										// We currently only support TX retries on cockroach
										gen(i, token)
										gen(i, token)
									} else {
										go gen(i, token)
										go gen(i, token)
									}
								}

								wg.Wait()
								if withSleep > 0 {
									time.Sleep(withSleep)
								}
							}
							return generations
						}

						t.Run("re-using an old graceful refresh token invalidates all tokens", func(t *testing.T) {
							reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "1s")
							reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1m")
							// This test only works if the refresh token lifespan is longer than the grace period.
							generations := createTokenGenerations(t, 4, time.Second*2)

							generationIndex := rng.Intn(len(generations) - 1) // Exclude the last generation
							tokenIndex := rng.Intn(len(generations[generationIndex]))

							token := generations[generationIndex][tokenIndex]
							token.Expiry = time.Now().Add(-time.Hour * 24)
							_, err := conf.TokenSource(ctx, token).Token()
							require.Error(t, err)

							// Now all tokens are inactive
							for i, generation := range generations {
								t.Run(fmt.Sprintf("generation=%d", i), func(t *testing.T) {
									for j, token := range generation {
										t.Run(fmt.Sprintf("token=%d", j), func(t *testing.T) {
											assertInactive(t, token.AccessToken, conf)
											assertInactive(t, token.RefreshToken, conf)
										})
									}
								})
							}
						})

						for _, withSleep := range []time.Duration{0, aboveGracePeriod} {
							t.Run(fmt.Sprintf("withSleep=%s", withSleep), func(t *testing.T) {
								createTokenGenerations := func(t *testing.T, count int) [][]*oauth2.Token {
									return createTokenGenerations(t, count, withSleep)
								}

								t.Run("only the most recent token generation is valid across the board", func(t *testing.T) {
									generations := createTokenGenerations(t, 4)

									// All generations except the last one are valid.
									for i, generation := range generations[:len(generations)-1] {
										t.Run(fmt.Sprintf("generation=%d", i), func(t *testing.T) {
											for j, token := range generation {
												t.Run(fmt.Sprintf("token=%d", j), func(t *testing.T) {
													assertInactive(t, token.AccessToken, conf)
												})
											}
										})
									}

									// The last generation is valid:
									t.Run(fmt.Sprintf("generation=%d", len(generations)-1), func(t *testing.T) {
										for j, token := range generations[len(generations)-1] {
											t.Run(fmt.Sprintf("token=%d", j), func(t *testing.T) {
												introspectAccessToken(t, conf, token, subject)
												assertIDToken(t, token, conf, subject, nonce, time.Time{})
												assertRefreshToken(t, token, conf, time.Time{})
											})
										}
									})
								})

								t.Run("revoking consent revokes all tokens", func(t *testing.T) {
									generations := createTokenGenerations(t, 4)

									// After revoking consent, all generations are invalid.
									err := reg.ConsentManager().RevokeSubjectConsentSession(context.Background(), subject)
									require.NoError(t, err)

									for i, generation := range generations {
										t.Run(fmt.Sprintf("generation=%d", i), func(t *testing.T) {
											for j, token := range generation {
												t.Run(fmt.Sprintf("token=%d", j), func(t *testing.T) {
													assertInactive(t, token.AccessToken, conf)
													assertInactive(t, token.RefreshToken, conf)
												})
											}
										})
									}
								})

								t.Run("re-using the a recent refresh token after the grace period has ended invalidates all tokens", func(t *testing.T) {
									generations := createTokenGenerations(t, 4)

									token := generations[len(generations)-1][0]

									finalToken := refreshTokens(t, token)
									time.Sleep(aboveGracePeriod) // Wait for the grace period to end

									token.Expiry = time.Now().Add(-time.Hour * 24)
									_, err := conf.TokenSource(ctx, token).Token()
									require.Error(t, err)

									// Now all tokens are inactive
									for i, generation := range append(generations, []*oauth2.Token{finalToken}) {
										t.Run(fmt.Sprintf("generation=%d", i), func(t *testing.T) {
											for j, token := range generation {
												t.Run(fmt.Sprintf("token=%d", j), func(t *testing.T) {
													assertInactive(t, token.AccessToken, conf)
													assertInactive(t, token.RefreshToken, conf)
												})
											}
										})
									}
								})

								t.Run("revoking a refresh token in the chain revokes all tokens", func(t *testing.T) {
									generations := createTokenGenerations(t, 4)

									testhelpers.RevokeToken(t, conf, generations[len(generations)-1][0].RefreshToken, publicTS)

									for i, generation := range generations {
										t.Run(fmt.Sprintf("generation=%d", i), func(t *testing.T) {
											for j, token := range generation {
												token := token
												t.Run(fmt.Sprintf("token=%d", j), func(t *testing.T) {
													assertInactive(t, token.AccessToken, conf)
													assertInactive(t, token.RefreshToken, conf)
												})
											}
										})
									}
								})
							})
						}
					})

					t.Run("it is possible to refresh tokens concurrently", func(t *testing.T) {
						// SQLite can not handle concurrency
						if dbName == "memory" {
							t.Skip("Skipping test because SQLite can not handle concurrency")
						}

						reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, "1m")
						reg.Config().MustSet(ctx, config.KeyRefreshTokenRotationGracePeriod, "5s")

						token := issueTokens(t)

						var wg sync.WaitGroup
						refresh := func(t *testing.T, token *oauth2.Token) *oauth2.Token {
							require.NotEmpty(t, token.RefreshToken)
							token.Expiry = time.Now().Add(-time.Hour * 24)
							tt, err := conf.TokenSource(context.Background(), token).Token()
							require.NoError(t, err)
							return tt
						}

						refreshes := make([]*oauth2.Token, 5)
						for k := range refreshes {
							wg.Add(1)
							go func(k int) {
								defer wg.Done()
								refreshes[k] = refresh(t, token)
							}(k)
						}
						wg.Wait()

						// All tokens are valid.
						for k, actual := range refreshes {
							refresh := actual
							require.NotEmpty(t, refresh.RefreshToken, "token %d:\ntoken:%+v", k, refresh)
							require.NotEmpty(t, refresh.AccessToken, "token %d:\ntoken:%+v", k, refresh)
							require.NotEmpty(t, refresh.Extra("id_token"), "token %d:\ntoken:%+v", k, refresh)

							i := testhelpers.IntrospectToken(t, conf, refresh.AccessToken, adminTS)
							assert.Truef(t, i.Get("active").Bool(), "token %d:\ntoken:%+v\nresult:%s", k, refresh, i)

							i = testhelpers.IntrospectToken(t, conf, refresh.RefreshToken, adminTS)
							assert.Truef(t, i.Get("active").Bool(), "token %d:\ntoken:%+v\nresult:%s", k, refresh, i)
						}
					})
				}

				t.Run("strategy=jwt", func(t *testing.T) {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
					run(t, "jwt")
				})

				t.Run("strategy=opaque", func(t *testing.T) {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
					run(t, "opaque")
				})
			})
		})
	}
}

func assertCreateVerifiableCredential(t *testing.T, reg driver.Registry, nonce string, accessToken *oauth2.Token, alg jose.SignatureAlgorithm) {
	// Build a proof from the nonce.
	pubKey, privKey, err := josex.NewSigningKey(alg, 0)
	require.NoError(t, err)
	pubKeyJWK := &jose.JSONWebKey{Key: pubKey, Algorithm: string(alg)}
	proofJWT := createVCProofJWT(t, pubKeyJWK, privKey, nonce)

	// Assert that we can fetch a verifiable credential with the nonce.
	verifiableCredential, err := createVerifiableCredential(t, reg, accessToken, &hydraoauth2.CreateVerifiableCredentialRequestBody{
		Format: "jwt_vc_json",
		Types:  []string{"VerifiableCredential", "UserInfoCredential"},
		Proof: &hydraoauth2.VerifiableCredentialProof{
			ProofType: "jwt",
			JWT:       proofJWT,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, verifiableCredential)

	_, claims := claimsFromVCResponse(t, reg, verifiableCredential)
	assertClaimsContainPublicKey(t, claims, pubKeyJWK)
}

func claimsFromVCResponse(t *testing.T, reg driver.Registry, vc *hydraoauth2.VerifiableCredentialResponse) (*jwt.Token, *hydraoauth2.VerifableCredentialClaims) {
	ctx := context.Background()
	token, err := jwt.ParseWithClaims(vc.Credential, new(hydraoauth2.VerifableCredentialClaims), func(token *jwt.Token) (interface{}, error) {
		kid, found := token.Header["kid"]
		if !found {
			return nil, errors.New("missing kid header")
		}
		openIDKey, err := reg.OpenIDJWTStrategy().GetPublicKeyID(ctx)
		if err != nil {
			return nil, err
		}
		if kid != openIDKey {
			return nil, errors.New("invalid kid header")
		}

		return x.Must(reg.OpenIDJWTStrategy().GetPublicKey(ctx)).Key, nil
	})
	require.NoError(t, err)

	return token, token.Claims.(*hydraoauth2.VerifableCredentialClaims)
}

func assertClaimsContainPublicKey(t *testing.T, claims *hydraoauth2.VerifableCredentialClaims, pubKeyJWK *jose.JSONWebKey) {
	pubKeyRaw, err := pubKeyJWK.MarshalJSON()
	require.NoError(t, err)
	expectedID := fmt.Sprintf("did:jwk:%s", base64.RawURLEncoding.EncodeToString(pubKeyRaw))
	require.Equal(t, expectedID, claims.VerifiableCredential.Subject["id"])
}

func createVerifiableCredential(
	t *testing.T,
	reg driver.Registry,
	token *oauth2.Token,
	createVerifiableCredentialReq *hydraoauth2.CreateVerifiableCredentialRequestBody,
) (vcRes *hydraoauth2.VerifiableCredentialResponse, vcErr error) {
	var (
		ctx  = context.Background()
		body bytes.Buffer
	)
	require.NoError(t, json.NewEncoder(&body).Encode(createVerifiableCredentialReq))
	req := httpx.MustNewRequest("POST", reg.Config().CredentialsEndpointURL(ctx).String(), &body, "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes fosite.RFC6749Error
		require.NoError(t, json.NewDecoder(res.Body).Decode(&errRes))
		return nil, &errRes
	}
	require.Equal(t, http.StatusOK, res.StatusCode)
	var vc hydraoauth2.VerifiableCredentialResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&vc))

	return &vc, vcErr
}

func doPrimingRequest(
	t *testing.T,
	reg driver.Registry,
	token *oauth2.Token,
	createVerifiableCredentialReq *hydraoauth2.CreateVerifiableCredentialRequestBody,
) (*hydraoauth2.VerifiableCredentialPrimingResponse, error) {
	var (
		ctx  = context.Background()
		body bytes.Buffer
	)
	require.NoError(t, json.NewEncoder(&body).Encode(createVerifiableCredentialReq))
	req := httpx.MustNewRequest("POST", reg.Config().CredentialsEndpointURL(ctx).String(), &body, "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var vc hydraoauth2.VerifiableCredentialPrimingResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&vc))

	return &vc, nil
}

func createVCProofJWT(t *testing.T, pubKey *jose.JSONWebKey, privKey any, nonce string) string {
	proofToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(pubKey.Algorithm)), jwt.MapClaims{"nonce": nonce})
	proofToken.Header["jwk"] = pubKey
	proofJWT, err := proofToken.SignedString(privKey)
	require.NoError(t, err)

	return proofJWT
}

// TestAuthCodeWithMockStrategy runs the authorization_code flow against various ConsentStrategy scenarios.
// For that purpose, the consent strategy is mocked so all scenarios can be applied properly. This test suite checks:
//
// - [x] should pass request if strategy passes
// - [x] should fail because prompt=none and max_age > auth_time
// - [x] should pass because prompt=none and max_age < auth_time
// - [x] should fail because prompt=none but auth_time suggests recent authentication
// - [x] should fail because consent strategy fails
// - [x] should pass with prompt=login when authentication time is recent
// - [x] should fail with prompt=login when authentication time is in the past
func TestAuthCodeWithMockStrategy(t *testing.T) {
	ctx := context.Background()
	for _, strat := range []struct{ d string }{{d: "opaque"}, {d: "jwt"}} {
		t.Run("strategy="+strat.d, func(t *testing.T) {
			conf := testhelpers.NewConfigurationWithDefaults()
			conf.MustSet(ctx, config.KeyAccessTokenLifespan, time.Second*2)
			conf.MustSet(ctx, config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
			conf.MustSet(ctx, config.KeyAccessTokenStrategy, strat.d)
			reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})
			testhelpers.MustEnsureRegistryKeys(ctx, reg, x.OpenIDConnectKeyName)
			testhelpers.MustEnsureRegistryKeys(ctx, reg, x.OAuth2JWTKeyName)

			consentStrategy := &consentMock{}
			router := x.NewRouterPublic()
			ts := httptest.NewServer(router)
			t.Cleanup(ts.Close)

			reg.WithConsentStrategy(consentStrategy)
			handler := reg.OAuth2Handler()
			handler.SetRoutes(httprouterx.NewRouterAdminWithPrefixAndRouter(router.Router, "/admin", conf.AdminURL), router, func(h http.Handler) http.Handler {
				return h
			})

			var callbackHandler *httprouter.Handle
			router.GET("/callback", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				(*callbackHandler)(w, r, ps)
			})
			var mutex sync.Mutex

			require.NoError(t, reg.ClientManager().CreateClient(ctx, &client.Client{
				ID:            "app-client",
				Secret:        "secret",
				RedirectURIs:  []string{ts.URL + "/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scope:         "hydra.* offline openid",
			}))

			oauthConfig := &oauth2.Config{
				ClientID:     "app-client",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:  ts.URL + "/oauth2/auth",
					TokenURL: ts.URL + "/oauth2/token",
				},
				RedirectURL: ts.URL + "/callback",
				Scopes:      []string{"hydra.*", "offline", "openid"},
			}

			var code string
			for k, tc := range []struct {
				cj                        http.CookieJar
				d                         string
				cb                        func(t *testing.T) httprouter.Handle
				authURL                   string
				shouldPassConsentStrategy bool
				expectOAuthAuthError      bool
				expectOAuthTokenError     bool
				checkExpiry               bool
				authTime                  time.Time
				requestTime               time.Time
				assertAccessToken         func(*testing.T, string)
			}{
				{
					d:                         "should pass request if strategy passes",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state"),
					shouldPassConsentStrategy: true,
					checkExpiry:               true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.NotEmpty(t, code)
							_, _ = w.Write([]byte(r.URL.Query().Get("code")))
						}
					},
					assertAccessToken: func(t *testing.T, token string) {
						if strat.d != "jwt" {
							return
						}

						body, err := x.DecodeSegment(strings.Split(token, ".")[1])
						require.NoError(t, err)

						data := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &data))

						assert.EqualValues(t, "app-client", data["client_id"])
						assert.EqualValues(t, "foo", data["sub"])
						assert.NotEmpty(t, data["iss"])
						assert.NotEmpty(t, data["jti"])
						assert.NotEmpty(t, data["exp"])
						assert.NotEmpty(t, data["iat"])
						assert.NotEmpty(t, data["nbf"])
						assert.EqualValues(t, data["nbf"], data["iat"])
						assert.EqualValues(t, []interface{}{"offline", "openid", "hydra.*"}, data["scp"])
					},
				},
				{
					d:                         "should fail because prompt=none and max_age > auth_time",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=none&max_age=1",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC(),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							err := r.URL.Query().Get("error")
							require.Empty(t, code)
							require.EqualValues(t, fosite.ErrLoginRequired.Error(), err)
						}
					},
					expectOAuthAuthError: true,
				},
				{
					d:                         "should pass because prompt=none and max_age is less than auth_time",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=none&max_age=3600",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC(),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.NotEmpty(t, code)
							_, _ = w.Write([]byte(r.URL.Query().Get("code")))
						}
					},
				},
				{
					d:                         "should fail because prompt=none but auth_time suggests recent authentication",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=none",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC().Add(-time.Hour),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							err := r.URL.Query().Get("error")
							require.Empty(t, code)
							require.EqualValues(t, fosite.ErrLoginRequired.Error(), err)
						}
					},
					expectOAuthAuthError: true,
				},
				{
					d:                         "should fail because consent strategy fails",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state"),
					expectOAuthAuthError:      true,
					shouldPassConsentStrategy: false,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							require.Empty(t, r.URL.Query().Get("code"))
							assert.Equal(t, fosite.ErrRequestForbidden.Error(), r.URL.Query().Get("error"))
						}
					},
				},
				{
					d:                         "should pass with prompt=login when authentication time is recent",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=login",
					authTime:                  time.Now().UTC().Add(-time.Second),
					requestTime:               time.Now().UTC().Add(-time.Minute),
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.NotEmpty(t, code)
							_, _ = w.Write([]byte(r.URL.Query().Get("code")))
						}
					},
				},
				{
					d:                         "should fail with prompt=login when authentication time is in the past",
					authURL:                   oauthConfig.AuthCodeURL("some-foo-state") + "&prompt=login",
					authTime:                  time.Now().UTC().Add(-time.Minute),
					requestTime:               time.Now().UTC(),
					expectOAuthAuthError:      true,
					shouldPassConsentStrategy: true,
					cb: func(t *testing.T) httprouter.Handle {
						return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
							code = r.URL.Query().Get("code")
							require.Empty(t, code)
							assert.Equal(t, fosite.ErrLoginRequired.Error(), r.URL.Query().Get("error"))
						}
					},
				},
			} {
				t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
					mutex.Lock()
					defer mutex.Unlock()
					if tc.cb == nil {
						tc.cb = noopHandler
					}

					consentStrategy.deny = !tc.shouldPassConsentStrategy
					consentStrategy.authTime = tc.authTime
					consentStrategy.requestTime = tc.requestTime

					cb := tc.cb(t)
					callbackHandler = &cb

					req, err := http.NewRequest("GET", tc.authURL, nil)
					require.NoError(t, err)

					if tc.cj == nil {
						tc.cj = testhelpers.NewEmptyCookieJar(t)
					}

					resp, err := (&http.Client{Jar: tc.cj}).Do(req)
					require.NoError(t, err, tc.authURL, ts.URL)
					defer resp.Body.Close()

					if tc.expectOAuthAuthError {
						require.Empty(t, code)
						return
					}

					require.NotEmpty(t, code)

					token, err := oauthConfig.Exchange(context.TODO(), code)
					if tc.expectOAuthTokenError {
						require.Error(t, err)
						return
					}

					require.NoError(t, err, code)
					if tc.assertAccessToken != nil {
						tc.assertAccessToken(t, token.AccessToken)
					}

					t.Run("case=userinfo", func(t *testing.T) {
						var makeRequest = func(req *http.Request) *http.Response {
							resp, err = http.DefaultClient.Do(req)
							require.NoError(t, err)
							return resp
						}

						var testSuccess = func(response *http.Response) {
							defer resp.Body.Close()

							require.Equal(t, http.StatusOK, resp.StatusCode)

							var claims map[string]interface{}
							require.NoError(t, json.NewDecoder(resp.Body).Decode(&claims))
							assert.Equal(t, "foo", claims["sub"])
						}

						req, err = http.NewRequest("GET", ts.URL+"/userinfo", nil)
						req.Header.Add("Authorization", "bearer "+token.AccessToken)
						testSuccess(makeRequest(req))

						req, err = http.NewRequest("POST", ts.URL+"/userinfo", nil)
						req.Header.Add("Authorization", "bearer "+token.AccessToken)
						testSuccess(makeRequest(req))

						req, err = http.NewRequest("POST", ts.URL+"/userinfo", bytes.NewBuffer([]byte("access_token="+token.AccessToken)))
						req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						testSuccess(makeRequest(req))

						req, err = http.NewRequest("GET", ts.URL+"/userinfo", nil)
						req.Header.Add("Authorization", "bearer asdfg")
						resp := makeRequest(req)
						require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
					})

					res, err := testRefresh(t, token, ts.URL, tc.checkExpiry)
					require.NoError(t, err)
					assert.Equal(t, http.StatusOK, res.StatusCode)

					body, err := io.ReadAll(res.Body)
					require.NoError(t, err)

					var refreshedToken oauth2.Token
					require.NoError(t, json.Unmarshal(body, &refreshedToken))

					if tc.assertAccessToken != nil {
						tc.assertAccessToken(t, refreshedToken.AccessToken)
					}

					t.Run("the tokens should be different", func(t *testing.T) {
						if strat.d != "jwt" {
							t.Skip()
						}

						body, err := x.DecodeSegment(strings.Split(token.AccessToken, ".")[1])
						require.NoError(t, err)

						origPayload := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &origPayload))

						body, err = x.DecodeSegment(strings.Split(refreshedToken.AccessToken, ".")[1])
						require.NoError(t, err)

						refreshedPayload := map[string]interface{}{}
						require.NoError(t, json.Unmarshal(body, &refreshedPayload))

						if tc.checkExpiry {
							assert.NotEqual(t, refreshedPayload["exp"], origPayload["exp"])
							assert.NotEqual(t, refreshedPayload["iat"], origPayload["iat"])
							assert.NotEqual(t, refreshedPayload["nbf"], origPayload["nbf"])
						}
						assert.NotEqual(t, refreshedPayload["jti"], origPayload["jti"])
						assert.Equal(t, refreshedPayload["client_id"], origPayload["client_id"])
					})

					require.NotEqual(t, token.AccessToken, refreshedToken.AccessToken)

					t.Run("old token should no longer be usable", func(t *testing.T) {
						req, err := http.NewRequest("GET", ts.URL+"/userinfo", nil)
						require.NoError(t, err)
						req.Header.Add("Authorization", "bearer "+token.AccessToken)
						res, err := http.DefaultClient.Do(req)
						require.NoError(t, err)
						assert.EqualValues(t, http.StatusUnauthorized, res.StatusCode)
					})

					t.Run("refreshing new refresh token should work", func(t *testing.T) {
						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusOK, res.StatusCode)

						body, err := io.ReadAll(res.Body)
						require.NoError(t, err)
						require.NoError(t, json.Unmarshal(body, &refreshedToken))
					})

					t.Run("should call refresh token hook if configured", func(t *testing.T) {
						run := func(hookType string) func(t *testing.T) {
							return func(t *testing.T) {
								hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

									expectedGrantedScopes := []string{"openid", "offline", "hydra.*"}
									expectedSubject := "foo"

									exceptKeys := []string{
										"session.kid",
										"session.id_token.expires_at",
										"session.id_token.headers.extra.kid",
										"session.id_token.id_token_claims.iat",
										"session.id_token.id_token_claims.exp",
										"session.id_token.id_token_claims.rat",
										"session.id_token.id_token_claims.auth_time",
									}

									if hookType == "legacy" {
										var hookReq hydraoauth2.RefreshTokenHookRequest
										require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
										require.Equal(t, hookReq.Subject, expectedSubject)
										require.ElementsMatch(t, hookReq.GrantedScopes, expectedGrantedScopes)
										require.ElementsMatch(t, hookReq.GrantedAudience, []string{})
										require.Equal(t, hookReq.ClientID, oauthConfig.ClientID)
										require.NotEmpty(t, hookReq.Session)
										require.Equal(t, hookReq.Session.Subject, expectedSubject)
										require.Equal(t, hookReq.Session.ClientID, oauthConfig.ClientID)
										require.NotEmpty(t, hookReq.Requester)
										require.Equal(t, hookReq.Requester.ClientID, oauthConfig.ClientID)
										require.ElementsMatch(t, hookReq.Requester.GrantedScopes, expectedGrantedScopes)

										snapshotx.SnapshotT(t, hookReq, snapshotx.ExceptPaths(exceptKeys...))
									} else {
										var hookReq hydraoauth2.TokenHookRequest
										require.NoError(t, json.NewDecoder(r.Body).Decode(&hookReq))
										require.NotEmpty(t, hookReq.Session)
										require.Equal(t, hookReq.Session.Subject, expectedSubject)
										require.Equal(t, hookReq.Session.ClientID, oauthConfig.ClientID)
										require.NotEmpty(t, hookReq.Request)
										require.Equal(t, hookReq.Request.ClientID, oauthConfig.ClientID)
										require.ElementsMatch(t, hookReq.Request.GrantedScopes, expectedGrantedScopes)
										require.ElementsMatch(t, hookReq.Request.GrantedAudience, []string{})
										require.Equal(t, hookReq.Request.Payload, map[string][]string{"grant_type": {"refresh_token"}})

										snapshotx.SnapshotT(t, hookReq, snapshotx.ExceptPaths(exceptKeys...))
									}

									claims := map[string]interface{}{
										"hooked": hookType,
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

								if hookType == "legacy" {
									conf.MustSet(ctx, config.KeyRefreshTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyRefreshTokenHook, nil)

								} else {
									conf.MustSet(ctx, config.KeyTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyTokenHook, nil)
								}

								res, err := testRefresh(t, &refreshedToken, ts.URL, false)
								require.NoError(t, err)
								assert.Equal(t, http.StatusOK, res.StatusCode)

								body, err := io.ReadAll(res.Body)
								require.NoError(t, err)
								require.NoError(t, json.Unmarshal(body, &refreshedToken))

								accessTokenClaims := testhelpers.IntrospectToken(t, oauthConfig, refreshedToken.AccessToken, ts)
								require.Equal(t, accessTokenClaims.Get("ext.hooked").String(), hookType)

								idTokenBody, err := x.DecodeSegment(
									strings.Split(
										gjson.GetBytes(body, "id_token").String(),
										".",
									)[1],
								)
								require.NoError(t, err)

								require.Equal(t, gjson.GetBytes(idTokenBody, "hooked").String(), hookType)
							}
						}
						t.Run("hook=legacy", run("legacy"))
						t.Run("hook=new", run("new"))
					})

					t.Run("should not override session data if token refresh hook returns no content", func(t *testing.T) {
						run := func(hookType string) func(t *testing.T) {
							return func(t *testing.T) {
								hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									w.WriteHeader(http.StatusNoContent)
								}))
								defer hs.Close()

								if hookType == "legacy" {
									conf.MustSet(ctx, config.KeyRefreshTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyRefreshTokenHook, nil)
								} else {
									conf.MustSet(ctx, config.KeyTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyTokenHook, nil)
								}

								origAccessTokenClaims := testhelpers.IntrospectToken(t, oauthConfig, refreshedToken.AccessToken, ts)

								res, err := testRefresh(t, &refreshedToken, ts.URL, false)
								require.NoError(t, err)
								assert.Equal(t, http.StatusOK, res.StatusCode)

								body, err = io.ReadAll(res.Body)
								require.NoError(t, err)

								require.NoError(t, json.Unmarshal(body, &refreshedToken))

								refreshedAccessTokenClaims := testhelpers.IntrospectToken(t, oauthConfig, refreshedToken.AccessToken, ts)
								assertx.EqualAsJSONExcept(t, json.RawMessage(origAccessTokenClaims.Raw), json.RawMessage(refreshedAccessTokenClaims.Raw), []string{"exp", "iat", "nbf"})
							}
						}
						t.Run("hook=legacy", run("legacy"))
						t.Run("hook=new", run("new"))
					})

					t.Run("should fail token refresh with `server_error` if refresh hook fails", func(t *testing.T) {
						run := func(hookType string) func(t *testing.T) {
							return func(t *testing.T) {
								hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									w.WriteHeader(http.StatusInternalServerError)
								}))
								defer hs.Close()

								if hookType == "legacy" {
									conf.MustSet(ctx, config.KeyRefreshTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyRefreshTokenHook, nil)
								} else {
									conf.MustSet(ctx, config.KeyTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyTokenHook, nil)
								}

								res, err := testRefresh(t, &refreshedToken, ts.URL, false)
								require.NoError(t, err)
								assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

								var errBody fosite.RFC6749ErrorJson
								require.NoError(t, json.NewDecoder(res.Body).Decode(&errBody))
								require.Equal(t, fosite.ErrServerError.Error(), errBody.Name)
								require.Equal(t, "An error occurred while executing the token hook.", errBody.Description)
							}
						}
						t.Run("hook=legacy", run("legacy"))
						t.Run("hook=new", run("new"))
					})

					t.Run("should fail token refresh with `access_denied` if legacy refresh hook denied the request", func(t *testing.T) {
						run := func(hookType string) func(t *testing.T) {
							return func(t *testing.T) {
								hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									w.WriteHeader(http.StatusForbidden)
								}))
								defer hs.Close()

								if hookType == "legacy" {
									conf.MustSet(ctx, config.KeyRefreshTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyRefreshTokenHook, nil)
								} else {
									conf.MustSet(ctx, config.KeyTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyTokenHook, nil)
								}

								res, err := testRefresh(t, &refreshedToken, ts.URL, false)
								require.NoError(t, err)
								assert.Equal(t, http.StatusForbidden, res.StatusCode)

								var errBody fosite.RFC6749ErrorJson
								require.NoError(t, json.NewDecoder(res.Body).Decode(&errBody))
								require.Equal(t, fosite.ErrAccessDenied.Error(), errBody.Name)
								require.Equal(t, "The token hook target responded with an error. Make sure that the request you are making is valid. Maybe the credential or request parameters you are using are limited in scope or otherwise restricted.", errBody.Description)
							}
						}
						t.Run("hook=legacy", run("legacy"))
						t.Run("hook=new", run("new"))
					})

					t.Run("should fail token refresh with `server_error` if refresh hook response is malformed", func(t *testing.T) {
						run := func(hookType string) func(t *testing.T) {
							return func(t *testing.T) {
								hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
									w.WriteHeader(http.StatusOK)
								}))
								defer hs.Close()

								if hookType == "legacy" {
									conf.MustSet(ctx, config.KeyRefreshTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyRefreshTokenHook, nil)
								} else {
									conf.MustSet(ctx, config.KeyTokenHook, hs.URL)
									defer conf.MustSet(ctx, config.KeyTokenHook, nil)
								}

								res, err := testRefresh(t, &refreshedToken, ts.URL, false)
								require.NoError(t, err)
								assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

								var errBody fosite.RFC6749ErrorJson
								require.NoError(t, json.NewDecoder(res.Body).Decode(&errBody))
								require.Equal(t, fosite.ErrServerError.Error(), errBody.Name)
								require.Equal(t, "The token hook target responded with an error.", errBody.Description)
							}
						}
						t.Run("hook=legacy", run("legacy"))
						t.Run("hook=new", run("new"))
					})

					t.Run("refreshing old token should no longer work", func(t *testing.T) {
						res, err := testRefresh(t, token, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, res.StatusCode)
					})

					t.Run("attempt to refresh old token should revoke new token", func(t *testing.T) {
						res, err := testRefresh(t, &refreshedToken, ts.URL, false)
						require.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, res.StatusCode)
					})

					t.Run("duplicate code exchange fails", func(t *testing.T) {
						token, err := oauthConfig.Exchange(context.TODO(), code)
						require.Error(t, err)
						require.Nil(t, token)
					})

					code = ""
				})
			}
		})
	}
}

func testRefresh(t *testing.T, token *oauth2.Token, u string, sleep bool) (*http.Response, error) {
	if sleep {
		time.Sleep(time.Millisecond * 1001)
	}

	oauthClientConfig := &clientcredentials.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		TokenURL:     u + "/oauth2/token",
		Scopes:       []string{"foobar"},
	}

	req, err := http.NewRequest("POST", oauthClientConfig.TokenURL, strings.NewReader(url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{token.RefreshToken},
	}.Encode()))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(oauthClientConfig.ClientID, oauthClientConfig.ClientSecret)

	return http.DefaultClient.Do(req)
}

func withScope(scope string) func(*client.Client) {
	return func(c *client.Client) {
		c.Scope = scope
	}
}

func newOAuth2Client(
	t *testing.T,
	reg interface {
		config.Provider
		client.Registry
	},
	callbackURL string,
	opts ...func(*client.Client),
) (*client.Client, *oauth2.Config) {
	ctx := context.Background()
	secret := uuid.New()
	c := &client.Client{
		Secret:        secret,
		RedirectURIs:  []string{callbackURL},
		ResponseTypes: []string{"id_token", "code", "token"},
		GrantTypes: []string{
			"implicit",
			"refresh_token",
			"authorization_code",
			"password",
			"client_credentials",
		},
		Scope:    "hydra offline openid",
		Audience: []string{"https://api.ory.sh/"},
	}

	// apply options
	for _, o := range opts {
		o(c)
	}

	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))
	return c, &oauth2.Config{
		ClientID:     c.GetID(),
		ClientSecret: secret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:   reg.Config().OAuth2AuthURL(ctx).String(),
			TokenURL:  reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		Scopes: strings.Split(c.Scope, " "),
	}
}
