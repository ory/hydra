// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/configx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"
	"github.com/ory/x/uuidx"
)

func TestStrategyLoginConsentNext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyAccessTokenStrategy:            "opaque",
		config.KeyConsentRequestMaxAge:           time.Hour,
		config.KeyScopeStrategy:                  "exact",
		config.KeySubjectTypesSupported:          []string{"pairwise", "public"},
		config.KeySubjectIdentifierAlgorithmSalt: "76d5d2bf-747f-4592-9fbd-d2b895a54b3a",
	})))

	publicTS, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)
	adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
	adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

	oauth2Config := func(t *testing.T, c *client.Client) *oauth2.Config {
		return &oauth2.Config{
			ClientID:     c.GetID(),
			ClientSecret: c.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   publicTS.URL + "/oauth2/auth",
				TokenURL:  publicTS.URL + "/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
			RedirectURL: c.RedirectURIs[0],
		}
	}

	acceptLoginHandler := func(t *testing.T, subject string, payload *hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
		return checkAndAcceptLoginHandler(t, adminClient, subject, func(*testing.T, *hydra.OAuth2LoginRequest, error) hydra.AcceptOAuth2LoginRequest {
			if payload == nil {
				return hydra.AcceptOAuth2LoginRequest{}
			}
			return *payload
		})
	}

	acceptConsentHandler := func(t *testing.T, payload *hydra.AcceptOAuth2ConsentRequest) http.HandlerFunc {
		return checkAndAcceptConsentHandler(t, adminClient, func(*testing.T, *hydra.OAuth2ConsentRequest, error) hydra.AcceptOAuth2ConsentRequest {
			if payload == nil {
				return hydra.AcceptOAuth2ConsentRequest{}
			}
			return *payload
		})
	}

	createClientWithRedir := func(t *testing.T, redir string) *client.Client {
		c := &client.Client{RedirectURIs: []string{redir}}
		return createClient(t, reg, c)
	}

	createDefaultClient := func(t *testing.T) *client.Client {
		return createClientWithRedir(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
	}

	makeRequestAndExpectCode := func(t *testing.T, hc *http.Client, c *client.Client, values url.Values) string {
		_, res := makeOAuth2Request(t, reg, hc, c, values)
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		code := res.Request.URL.Query().Get("code")
		assert.NotEmpty(t, code)
		return code
	}

	makeRequestAndExpectError := func(t *testing.T, hc *http.Client, c *client.Client, values url.Values, errContains string) {
		_, res := makeOAuth2Request(t, reg, hc, c, values)
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Empty(t, res.Request.URL.Query().Get("code"))
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), errContains, "%v", res.Request.URL.Query())
	}

	t.Run("case=should fail because a login verifier was given that doesn't exist in the store", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)

		makeRequestAndExpectError(
			t, hc, c, url.Values{"login_verifier": {"does-not-exist"}},
			"The resource owner or authorization server denied the request. The login verifier has already been used, has not been granted, or is invalid.",
		)
	})

	t.Run("case=should fail because a non-existing consent verifier was given", func(t *testing.T) {
		// Covers:
		// - This should fail because consent verifier was set but does not exist
		// - This should fail because a consent verifier was given but no login verifier
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)

		makeRequestAndExpectError(
			t, hc, c, url.Values{"consent_verifier": {"does-not-exist"}},
			"The consent verifier has already been used, has not been granted, or is invalid.",
		)
	})

	t.Run("case=should fail because the request was redirected but the login endpoint doesn't do anything (like redirecting back)", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNotImplementedHandler, testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createClientWithRedir(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNoExpectedCallHandler(t)))

		_, res := makeOAuth2Request(t, reg, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("login_challenge"), "%s", res.Request.URL)
	})

	t.Run("case=should fail because the request was redirected but consent endpoint doesn't do anything (like redirecting back)", func(t *testing.T) {
		// "This should fail because consent endpoints idles after login was granted - but consent endpoint should be called because cookie jar exists"
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", nil), testhelpers.HTTPServerNotImplementedHandler)
		c := createClientWithRedir(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNoExpectedCallHandler(t)))

		_, res := makeOAuth2Request(t, reg, nil, c, url.Values{})
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.NotEmpty(t, res.Request.URL.Query().Get("consent_challenge"), "%s", res.Request.URL)
	})

	t.Run("case=should fail because the request was redirected but the login endpoint rejected the request", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), func(w http.ResponseWriter, r *http.Request) {
			vr, _, err := adminClient.OAuth2API.RejectOAuth2LoginRequest(context.Background()).
				LoginChallenge(r.URL.Query().Get("login_challenge")).
				RejectOAuth2Request(hydra.RejectOAuth2Request{
					Error:            pointerx.String(fosite.ErrInteractionRequired.ErrorField),
					ErrorDescription: pointerx.String("expect-reject-login"),
					StatusCode:       pointerx.Int64(int64(fosite.ErrInteractionRequired.CodeField)),
				}).Execute()
			require.NoError(t, err)
			assert.NotEmpty(t, vr.RedirectTo)
			http.Redirect(w, r, vr.RedirectTo, http.StatusFound)
		}, testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)

		makeRequestAndExpectError(t, nil, c, url.Values{}, "expect-reject-login")
	})

	t.Run("case=should fail because no cookie jar invalid csrf", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", nil),
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		hc := new(http.Client)
		hc.Jar = DropCookieJar(regexp.MustCompile("ory_hydra_.*_csrf_.*"))
		makeRequestAndExpectError(t, hc, c, url.Values{}, "No CSRF value available in the session cookie.")
	})

	t.Run("case=should fail because consent endpoints denies the request after login was granted", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			func(w http.ResponseWriter, r *http.Request) {
				vr, _, err := adminClient.OAuth2API.RejectOAuth2ConsentRequest(context.Background()).
					ConsentChallenge(r.URL.Query().Get("consent_challenge")).
					RejectOAuth2Request(hydra.RejectOAuth2Request{
						Error:            pointerx.String(fosite.ErrInteractionRequired.ErrorField),
						ErrorDescription: pointerx.String("expect-reject-consent"),
						StatusCode:       pointerx.Int64(int64(fosite.ErrInteractionRequired.CodeField))}).Execute()
				require.NoError(t, err)
				require.NotEmpty(t, vr.RedirectTo)
				http.Redirect(w, r, vr.RedirectTo, http.StatusFound)
			})

		makeRequestAndExpectError(t, nil, c, url.Values{}, "expect-reject-consent")
	})

	t.Run("suite=double-submit", func(t *testing.T) {
		ctx := context.Background()
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)
		var loginChallenge, consentChallenge string

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			func(w http.ResponseWriter, r *http.Request) {
				loginChallenge = r.URL.Query().Get("login_challenge")
				res, _, err := adminClient.OAuth2API.GetOAuth2LoginRequest(ctx).
					LoginChallenge(loginChallenge).
					Execute()
				require.NoError(t, err)
				require.Equal(t, loginChallenge, res.Challenge)

				v, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(ctx).
					LoginChallenge(loginChallenge).
					AcceptOAuth2LoginRequest(hydra.AcceptOAuth2LoginRequest{Subject: "aeneas-rekkas"}).
					Execute()
				require.NoError(t, err)
				require.NotEmpty(t, v.RedirectTo)
				http.Redirect(w, r, v.RedirectTo, http.StatusFound)
			},
			func(w http.ResponseWriter, r *http.Request) {
				consentChallenge = r.URL.Query().Get("consent_challenge")
				res, _, err := adminClient.OAuth2API.GetOAuth2ConsentRequest(ctx).
					ConsentChallenge(consentChallenge).
					Execute()
				require.NoError(t, err)
				require.Equal(t, consentChallenge, res.Challenge)

				v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(ctx).
					ConsentChallenge(consentChallenge).
					AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{}).
					Execute()
				require.NoError(t, err)
				require.NotEmpty(t, v.RedirectTo)
				http.Redirect(w, r, v.RedirectTo, http.StatusFound)
			})

		makeRequestAndExpectCode(t, hc, c, url.Values{})

		t.Run("case=double-submit login verifier", func(t *testing.T) {
			v, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(ctx).
				LoginChallenge(loginChallenge).
				AcceptOAuth2LoginRequest(hydra.AcceptOAuth2LoginRequest{Subject: "aeneas-rekkas"}).
				Execute()
			require.NoError(t, err)
			res, err := hc.Get(v.RedirectTo)
			require.NoError(t, err)
			q := res.Request.URL.Query()
			assert.Equal(t,
				"The resource owner or authorization server denied the request. The consent verifier has already been used.",
				q.Get("error_description"), q)
		})

		t.Run("case=double-submit consent verifier", func(t *testing.T) {
			v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(ctx).
				ConsentChallenge(consentChallenge).
				AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{}).
				Execute()
			require.NoError(t, err)
			res, err := hc.Get(v.RedirectTo)
			require.NoError(t, err)
			q := res.Request.URL.Query()
			assert.Equal(t,
				"The resource owner or authorization server denied the request. The consent verifier has already been used.",
				q.Get("error_description"), q)
		})

	})

	t.Run("case=should pass and set acr values properly", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			acceptConsentHandler(t, nil))

		makeRequestAndExpectCode(t, nil, c, url.Values{})
	})

	t.Run("case=should pass if both login and consent are granted and check remember flows as well as various payloads", func(t *testing.T) {
		// Covers old test cases:
		// - This should pass because login and consent have been granted, this time we remember the decision
		// - This should pass because login and consent have been granted, this time we remember the decision#2
		// - This should pass because login and consent have been granted, this time we remember the decision#3
		// - This should pass because login was remembered and session id should be set and session context should also work
		// - This should pass and confirm previous authentication and consent because it is a authorization_code

		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		now := 1723546027 // Unix timestamps must round-trip through Hydra without converting to floats or similar
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Remember: pointerx.Bool(true),
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
				Remember:   pointerx.Bool(true),
				GrantScope: []string{"openid"},
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{
						"foo": "bar",
						"ts1": now,
					},
					IdToken: map[string]interface{}{
						"bar": "baz",
						"ts2": now,
					},
				},
			}))

		hc := testhelpers.NewEmptyJarClient(t)
		conf := oauth2Config(t, c)

		var sid string
		var run = func(t *testing.T) {
			code := makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]},
				"scope": {"openid"}})

			token, err := conf.Exchange(context.Background(), code)
			require.NoError(t, err)

			claims := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
			assert.Equalf(t, `"bar"`, claims.Get("ext.foo").Raw, "%s", claims.Raw)      // Raw rather than .Int() or .Value() to verify the exact JSON payload
			assert.Equalf(t, "1723546027", claims.Get("ext.ts1").Raw, "%s", claims.Raw) // must round-trip as integer

			idClaims := testhelpers.DecodeIDToken(t, token)
			assert.Equalf(t, `"baz"`, idClaims.Get("bar").Raw, "%s", idClaims.Raw)      // Raw rather than .Int() or .Value() to verify the exact JSON payload
			assert.Equalf(t, "1723546027", idClaims.Get("ts2").Raw, "%s", idClaims.Raw) // must round-trip as integer
			sid = idClaims.Get("sid").String()
			assert.NotEmpty(t, sid)
		}

		t.Run("perform first flow", run)

		t.Run("perform follow up flows and check if session values are set", func(t *testing.T) {
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
					require.NoError(t, err)
					assert.True(t, res.Skip)
					assert.Equal(t, sid, *res.SessionId)
					assert.Equal(t, subject, res.Subject)
					assert.Empty(t, pointerx.StringR(res.Client.ClientSecret))
					return hydra.AcceptOAuth2LoginRequest{
						Subject: subject,
						Context: map[string]interface{}{"xyz": "abc"},
					}
				}),
				checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, req *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
					require.NoError(t, err)
					assert.True(t, *req.Skip)
					assert.Equal(t, sid, *req.LoginSessionId)
					assert.Equal(t, subject, *req.Subject)
					assert.Empty(t, pointerx.StringR(req.Client.ClientSecret))
					assert.Equal(t, map[string]interface{}{"xyz": "abc"}, req.Context)
					return hydra.AcceptOAuth2ConsentRequest{
						Remember:   pointerx.Bool(true),
						GrantScope: []string{"openid"},
						Session: &hydra.AcceptOAuth2ConsentRequestSession{
							AccessToken: map[string]interface{}{
								"foo": "bar",
								"ts1": now,
							},
							IdToken: map[string]interface{}{
								"bar": "baz",
								"ts2": now,
							},
						},
					}
				}))

			for k := 0; k < 3; k++ {
				t.Run(fmt.Sprintf("case=%d", k), run)
			}
		})
	})

	t.Run("case=should set client specific csrf cookie names", func(t *testing.T) {
		subject := "subject-1"
		consentRequestMaxAge := reg.Config().ConsentRequestMaxAge(ctx).Seconds()
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Remember: pointerx.Bool(true),
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
				Remember:   pointerx.Bool(true),
				GrantScope: []string{"openid"},
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IdToken:     map[string]interface{}{"bar": "baz"},
				},
			}))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
				require.NoError(t, err)
				assert.Empty(t, res.Subject)
				assert.Empty(t, pointerx.StringR(res.Client.ClientSecret))
				return hydra.AcceptOAuth2LoginRequest{
					Subject: subject,
					Context: map[string]interface{}{"foo": "bar"},
				}
			}),
			checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, subject, *res.Subject)
				assert.Empty(t, pointerx.StringR(res.Client.ClientSecret))
				return hydra.AcceptOAuth2ConsentRequest{
					Remember:   pointerx.Bool(true),
					GrantScope: []string{"openid"},
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}
			}))
		hc := &http.Client{
			Jar:       testhelpers.NewEmptyCookieJar(t),
			Transport: &http.Transport{},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		_, oauthRes := makeOAuth2Request(t, reg, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}, "scope": {"openid"}})
		assert.EqualValues(t, http.StatusFound, oauthRes.StatusCode)
		loginChallengeRedirect, err := oauthRes.Location()
		require.NoError(t, err)
		defer oauthRes.Body.Close() //nolint:errcheck

		foundLoginCookie := slices.ContainsFunc(oauthRes.Header.Values("set-cookie"), func(sc string) bool {
			ok, err := regexp.MatchString(fmt.Sprintf("ory_hydra_login_csrf_dev_%s=.*Max-Age=%.0f;.*", c.CookieSuffix(), consentRequestMaxAge), sc)
			require.NoError(t, err)
			return ok
		})
		require.True(t, foundLoginCookie, "client-specific login cookie with max age set")

		loginChallengeRes, err := hc.Get(loginChallengeRedirect.String())
		require.NoError(t, err)
		defer loginChallengeRes.Body.Close() //nolint:errcheck

		loginVerifierRedirect, err := loginChallengeRes.Location()
		require.NoError(t, err)
		loginVerifierRes, err := hc.Get(loginVerifierRedirect.String())
		require.NoError(t, err)
		defer loginVerifierRes.Body.Close() //nolint:errcheck

		foundConsentCookie := slices.ContainsFunc(loginVerifierRes.Header.Values("set-cookie"), func(sc string) bool {
			ok, err := regexp.MatchString(fmt.Sprintf("ory_hydra_consent_csrf_dev_%s=.*Max-Age=%.0f;.*", c.CookieSuffix(), consentRequestMaxAge), sc)
			require.NoError(t, err)
			return ok
		})
		require.True(t, foundConsentCookie, "client-specific consent cookie with max age set")
	})

	t.Run("case=should pass if both login and consent are granted and check remember flows with refresh session cookie", func(t *testing.T) {

		subject := "subject-1"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Remember: pointerx.Bool(true),
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
				Remember:   pointerx.Bool(true),
				GrantScope: []string{"openid"},
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IdToken:     map[string]interface{}{"bar": "baz"},
				},
			}))

		hc := testhelpers.NewEmptyJarClient(t)

		followUpHandler := func(extendSessionLifespan bool) {
			rememberFor := int64(12345)
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
					require.NoError(t, err)
					assert.True(t, res.Skip)
					assert.Equal(t, subject, res.Subject)
					assert.Empty(t, res.Client.ClientSecret)
					return hydra.AcceptOAuth2LoginRequest{
						Subject:               subject,
						Remember:              pointerx.Bool(true),
						RememberFor:           pointerx.Int64(rememberFor),
						ExtendSessionLifespan: pointerx.Bool(extendSessionLifespan),
						Context:               map[string]interface{}{"foo": "bar"},
					}
				}),
				checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
					require.NoError(t, err)
					assert.True(t, *res.Skip)
					assert.Equal(t, subject, res.Subject)
					assert.Empty(t, res.Client.ClientSecret)
					return hydra.AcceptOAuth2ConsentRequest{
						Remember:   pointerx.Bool(true),
						GrantScope: []string{"openid"},
						Session: &hydra.AcceptOAuth2ConsentRequestSession{
							AccessToken: map[string]interface{}{"foo": "bar"},
							IdToken:     map[string]interface{}{"bar": "baz"},
						},
					}
				}))

			hc := &http.Client{
				Jar:       hc.Jar,
				Transport: &http.Transport{},
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			_, oauthRes := makeOAuth2Request(t, reg, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}, "scope": {"openid"}})
			assert.EqualValues(t, http.StatusFound, oauthRes.StatusCode)
			loginChallengeRedirect, err := oauthRes.Location()
			require.NoError(t, err)
			defer oauthRes.Body.Close() //nolint:errcheck

			loginChallengeRes, err := hc.Get(loginChallengeRedirect.String())
			require.NoError(t, err)
			defer loginChallengeRes.Body.Close() //nolint:errcheck
			loginVerifierRedirect, err := loginChallengeRes.Location()
			require.NoError(t, err)

			loginVerifierRes, err := hc.Get(loginVerifierRedirect.String())
			require.NoError(t, err)
			defer loginVerifierRes.Body.Close() //nolint:errcheck

			setCookieHeader := loginVerifierRes.Header.Get("set-cookie")
			assert.NotNil(t, setCookieHeader)
			if extendSessionLifespan {
				assert.Regexp(t, fmt.Sprintf("ory_hydra_session_dev=.*; Path=/; Expires=.*Max-Age=%d; HttpOnly; SameSite=Lax", rememberFor), setCookieHeader)
			} else {
				assert.NotContains(t, setCookieHeader, "ory_hydra_session_dev")
			}
		}

		t.Run("perform first flow", func(t *testing.T) {
			makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]},
				"scope": {"openid"}})
		})

		t.Run("perform follow up flow with extend_session_lifespan=false", func(t *testing.T) {
			followUpHandler(false)
		})

		t.Run("perform follow up flow with extend_session_lifespan=true", func(t *testing.T) {
			followUpHandler(true)
		})
	})

	t.Run("case=should set session cookie with correct configuration", func(t *testing.T) {
		cookiePath := "/foo"
		reg.Config().MustSet(ctx, config.KeyCookieSessionPath, cookiePath)
		defer reg.Config().MustSet(ctx, config.KeyCookieSessionPath, "/")

		subject := "subject-1"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Remember: pointerx.Bool(true),
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
				Remember:   pointerx.Bool(true),
				GrantScope: []string{"openid"},
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IdToken:     map[string]interface{}{"bar": "baz"},
				},
			}))
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
				require.NoError(t, err)
				assert.Empty(t, res.Subject)
				assert.Empty(t, pointerx.StringR(res.Client.ClientSecret))
				return hydra.AcceptOAuth2LoginRequest{
					Subject: subject,
					Context: map[string]interface{}{"foo": "bar"},
				}
			}),
			checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, subject, *res.Subject)
				assert.Empty(t, pointerx.StringR(res.Client.ClientSecret))
				return hydra.AcceptOAuth2ConsentRequest{
					Remember:   pointerx.Bool(true),
					GrantScope: []string{"openid"},
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}
			}))
		hc := &http.Client{
			Jar:       testhelpers.NewEmptyCookieJar(t),
			Transport: &http.Transport{},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		_, oauthRes := makeOAuth2Request(t, reg, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}, "scope": {"openid"}})
		assert.EqualValues(t, http.StatusFound, oauthRes.StatusCode)
		loginChallengeRedirect, err := oauthRes.Location()
		require.NoError(t, err)
		defer oauthRes.Body.Close() //nolint:errcheck

		loginChallengeRes, err := hc.Get(loginChallengeRedirect.String())
		require.NoError(t, err)
		defer loginChallengeRes.Body.Close() //nolint:errcheck

		loginVerifierRedirect, err := loginChallengeRes.Location()
		require.NoError(t, err)
		loginVerifierRes, err := hc.Get(loginVerifierRedirect.String())
		require.NoError(t, err)
		defer loginVerifierRes.Body.Close() //nolint:errcheck

		setCookieHeader := loginVerifierRes.Header.Get("set-cookie")
		assert.NotNil(t, setCookieHeader)

		assert.Regexp(t, fmt.Sprintf("ory_hydra_session_dev=.*; Path=%s; Expires=.*; Max-Age=0; HttpOnly; SameSite=Lax", cookiePath), setCookieHeader)
	})

	t.Run("case=should pass and check if login context is set properly", func(t *testing.T) {
		// This should pass because login was remembered and session id should be set and session context should also work
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Subject: subject,
				Context: map[string]interface{}{"fooz": "barz"},
			}),
			checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, map[string]interface{}{"fooz": "barz"}, res.Context)
				assert.Equal(t, subject, *res.Subject)
				return hydra.AcceptOAuth2ConsentRequest{
					Remember:   pointerx.Bool(true),
					GrantScope: []string{"openid"},
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}
			}))

		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})
	})

	t.Run("case=perform flows with a public client", func(t *testing.T) {
		// This test covers old cases:
		// - This should fail because prompt=none, client is public, and redirection scheme is not HTTPS but a custom scheme and a custom domain
		// - This should fail because prompt=none, client is public, and redirection scheme is not HTTPS but a custom scheme
		// - This should pass because prompt=none, client is public, redirection scheme is HTTP and host is localhost

		c := &client.Client{ID: uuidx.NewV4().String(), TokenEndpointAuthMethod: "none",
			RedirectURIs: []string{
				testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler),
				"custom://redirection-scheme/path",
				"custom://localhost/path",
			}}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))

		subject := "aeneas-rekkas"
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true), RememberFor: pointerx.Int64(0)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true), RememberFor: pointerx.Int64(0)}))

		hc := testhelpers.NewEmptyJarClient(t)

		// set up initial session
		makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})

		// By not waiting here we ensure that there are no race conditions when it comes to authenticated_at and
		// requested_at time comparisons:
		//
		//	time.Sleep(time.Second)

		t.Run("followup=should pass when prompt=none, redirection scheme is HTTP and host is localhost", func(t *testing.T) {
			makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}, "prompt": {"none"}})
		})

		t.Run("followup=should pass when prompt=none, redirection scheme is HTTP and host is a custom scheme", func(t *testing.T) {
			for _, redir := range c.RedirectURIs[1:] {
				t.Run("redir=should pass because prompt=none, client is public, and redirection is "+redir, func(t *testing.T) {
					_, err := hc.Get(urlx.CopyWithQuery(reg.Config().OAuth2AuthURL(ctx), url.Values{
						"response_type": {"code"},
						"state":         {uuid.New()},
						"redirect_uri":  {redir},
						"client_id":     {c.GetID()},
						"prompt":        {"none"},
					}).String())

					require.Error(t, err)
					assert.Contains(t, err.Error(), redir)

					// https://tools.ietf.org/html/rfc6749
					//
					// As stated in Section 10.2 of OAuth 2.0 [RFC6749], the authorization
					// server SHOULD NOT process authorization requests automatically
					// without user consent or interaction, except when the identity of the
					// client can be assured.  This includes the case where the user has
					// previously approved an authorization request for a given client id --
					// unless the identity of the client can be proven, the request SHOULD
					// be processed as if no previous request had been approved.
					//
					// Measures such as claimed "https" scheme redirects MAY be accepted by
					// authorization servers as identity proof.  Some operating systems may
					// offer alternative platform-specific identity features that MAY be
					assert.Contains(t, err.Error(), "error=consent_required")
				})
			}
		})
	})

	t.Run("case=should retry the authorization with prompt=login if subject in login challenge does not match subject from previous session", func(t *testing.T) {
		// Previously: This should fail at login screen because subject from accept does not match subject from session
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, nil))

		// Init session
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			func(w http.ResponseWriter, r *http.Request) {
				res, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
					LoginChallenge(r.URL.Query().Get("login_challenge")).
					AcceptOAuth2LoginRequest(hydra.AcceptOAuth2LoginRequest{
						Subject: "not-aeneas-rekkas",
					}).Execute()
				require.NoError(t, err)
				redirectURL, err := url.Parse(res.RedirectTo)
				require.NoError(t, err)
				assert.Equal(t, "login", redirectURL.Query().Get("prompt"))
				w.WriteHeader(http.StatusBadRequest)
			},
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		_, res := makeOAuth2Request(t, reg, hc, c, url.Values{})
		assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Request.URL.Query().Get("code"))
	})

	t.Run("case=should forward the identity schema in the login URL", func(t *testing.T) {
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "custom-id-schema", r.URL.Query().Get("identity_schema"))
				w.WriteHeader(http.StatusBadRequest) // We do not want to continue the flow here, we just want to check the query parameter
			},
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		_, res := makeOAuth2Request(t, reg, hc, c, url.Values{"identity_schema": {"custom-id-schema"}})
		assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("case=should require re-authentication when parameters mandate it", func(t *testing.T) {
		// Covers:
		// - should pass and require re-authentication although session is set (because prompt=login)
		// - should pass and require re-authentication although session is set (because max_age=1)
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		resetUI := func(t *testing.T) {
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
					require.NoError(t, err)
					assert.False(t, res.Skip) // Skip should always be false here
					return hydra.AcceptOAuth2LoginRequest{
						Remember: pointerx.Bool(true),
					}
				}),
				acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
					Remember: pointerx.Bool(true),
				}))
		}
		resetUI(t)

		// Init session
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		for k, values := range []url.Values{
			{"prompt": {"login"}},
			{"max_age": {"1"}},
			{"max_age": {"0"}},
		} {
			t.Run("values="+values.Encode(), func(t *testing.T) {
				if k == 1 {
					// If this is the max_age case we need to wait for max age to pass.
					time.Sleep(time.Second)
				}

				resetUI(t)
				makeRequestAndExpectCode(t, hc, c, values)
			})
		}
	})

	t.Run("case=should fail because max_age=1 but prompt=none", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)

		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		time.Sleep(time.Second)

		makeRequestAndExpectError(t, hc, c, url.Values{"max_age": {"1"}, "prompt": {"none"}},
			"prompt is set to 'none' and authentication time reached 'max_age'")
	})

	t.Run("case=should fail because prompt is none but no auth session exists", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

		makeRequestAndExpectError(t, nil, c, url.Values{"prompt": {"none"}},
			"Prompt 'none' was requested, but no existing login session was found")
	})

	t.Run("case=should fail because prompt is none and consent is missing a permission which requires re-authorization of the app", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

		// Init cookie
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		// Make request with additional scope and prompt none, which fails
		makeRequestAndExpectError(t, hc, c, url.Values{"prompt": {"none"}, "scope": {"openid"}, "redirect_uri": {c.RedirectURIs[0]}},
			"Prompt 'none' was requested, but no previous consent was found")
	})

	t.Run("case=pass and properly require authentication as well as authorization because prompt is set to login and consent although previous session exists", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
				require.NoError(t, err)
				assert.False(t, res.Skip) // Skip should always be false here because prompt has login
				return hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}
			}),
			checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.False(t, *res.Skip) // Skip should always be false here because prompt has consent
				return hydra.AcceptOAuth2ConsentRequest{
					Remember: pointerx.Bool(true),
				}
			}))

		// Init cookie
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		// Rerun with login and consent set
		makeRequestAndExpectCode(t, hc, c, url.Values{"prompt": {"login consent"}})
	})

	t.Run("case=should fail because id_token_hint does not match value from accepted login request", func(t *testing.T) {
		// Covers former tests:
		// - This should pass and require authentication because id_token_hint does not match subject from session
		// - This should fail because id_token_hint does not match authentication session and prompt is none
		// - This should fail because the user from the ID token does not match the user from the accept login request

		subject := "aeneas-rekkas"
		notSubject := "not-aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

		// Init cookie
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		for _, values := range []url.Values{
			{"prompt": {"none"}, "id_token_hint": {testhelpers.NewIDToken(t, reg, notSubject)}},
			{"id_token_hint": {testhelpers.NewIDToken(t, reg, notSubject)}},
		} {
			t.Run(fmt.Sprintf("values=%v", values), func(t *testing.T) {
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
						var b bytes.Buffer
						require.NoError(t, json.NewEncoder(&b).Encode(res))
						assert.EqualValues(t, notSubject, gjson.GetBytes(b.Bytes(), "oidc_context.id_token_hint_claims.sub"), b.String())
						return hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}
					}),
					acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

				makeRequestAndExpectError(t, hc, c, values,
					"Request failed because subject claim from id_token_hint does not match subject from authentication session")
			})
		}
	})

	t.Run("case=should pass and require authentication because id_token_hint does match subject from session", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

		makeRequestAndExpectCode(t, nil, c, url.Values{"id_token_hint": {testhelpers.NewIDToken(t, reg, subject)}})

		t.Run("case=should pass even though id_token_hint is expired", func(t *testing.T) {
			// Formerly: should pass as regularly even though id_token_hint is expired
			makeRequestAndExpectCode(t, nil, c, url.Values{
				"id_token_hint": {testhelpers.NewIDTokenWithExpiry(t, reg, subject, -time.Hour)}})
		})
	})

	t.Run("suite=pairwise auth", func(t *testing.T) {
		// Covers former tests:
		// - This should pass as regularly and create a new session with pairwise subject set by hydra
		// - This should pass as regularly and create a new session with pairwise subject and also with the ID token set

		c := createClient(t, reg, &client.Client{
			SubjectType:         "pairwise",
			SectorIdentifierURI: "foo",
			RedirectURIs:        []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler)},
		})

		subject := "auth-user"
		hash := fmt.Sprintf("%x",
			sha256.Sum256([]byte(c.SectorIdentifierURI+subject+reg.Config().SubjectIdentifierAlgorithmSalt(ctx))))

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true), GrantScope: []string{"openid"}}))

		for _, tc := range []struct {
			d      string
			values url.Values
		}{
			{
				d:      "check all the sub claims",
				values: url.Values{"scope": {"openid"}, "redirect_uri": {c.RedirectURIs[0]}},
			},
			{
				d:      "works with id_token_hint",
				values: url.Values{"scope": {"openid"}, "redirect_uri": {c.RedirectURIs[0]}, "id_token_hint": {testhelpers.NewIDToken(t, reg, hash)}},
			},
		} {
			t.Run("case="+tc.d, func(t *testing.T) {
				code := makeRequestAndExpectCode(t, nil, c, tc.values)

				conf := oauth2Config(t, c)
				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)

				// OpenID data must be obfuscated
				idClaims := testhelpers.DecodeIDToken(t, token)
				assert.EqualValues(t, hash, idClaims.Get("sub").String())
				uiClaims := testhelpers.Userinfo(t, token, publicTS)
				assert.EqualValues(t, hash, uiClaims.Get("sub").String())

				// Access token data must not be obfuscated
				atClaims := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
				assert.EqualValues(t, subject, atClaims.Get("sub").String())
			})
		}
	})

	t.Run("suite=pairwise auth with forced identifier", func(t *testing.T) {
		// Covers:
		// - This should pass as regularly and create a new session with pairwise subject set login request
		// - This should pass as regularly and create a new session with pairwise subject set on login request and also with the ID token set
		c := createClient(t, reg, &client.Client{
			SubjectType:         "pairwise",
			SectorIdentifierURI: "foo",
			RedirectURIs:        []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler)},
		})
		subject := "aeneas-rekkas"
		obfuscated := "obfuscated-friedrich-kaiser"
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				ForceSubjectIdentifier: &obfuscated,
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{GrantScope: []string{"openid"}}))

		code := makeRequestAndExpectCode(t, nil, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})

		conf := oauth2Config(t, c)
		token, err := conf.Exchange(context.Background(), code)
		require.NoError(t, err)

		// OpenID data must be obfuscated
		idClaims := testhelpers.DecodeIDToken(t, token)
		assert.EqualValues(t, obfuscated, idClaims.Get("sub").String())
		uiClaims := testhelpers.Userinfo(t, token, publicTS)
		assert.EqualValues(t, obfuscated, uiClaims.Get("sub").String())

		// Access token data must not be obfuscated
		atClaims := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
		assert.EqualValues(t, subject, atClaims.Get("sub").String())
	})

	t.Run("suite=properly clean up session cookies", func(t *testing.T) {
		t.Skip("This test is skipped because we forcibly set remember to true always when skip is also true for a better user experience.")

		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(true)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(true)}))

		// Initialize flow
		// Formerly: This should pass as regularly and create a new session and forward data
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		// Re-run flow but do not remember login
		// Formerly: This should pass and also revoke the session cookie
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{Remember: pointerx.Bool(false)}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{Remember: pointerx.Bool(false)}))
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		// Formerly: This should require re-authentication because the session was revoked in the previous test
		makeRequestAndExpectError(t, hc, c, url.Values{"prompt": {"none"}}, "...")
	})

	t.Run("case=should require re-authentication because the session does not exist in the store", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, subject, nil), acceptConsentHandler(t, nil))

		hc := &http.Client{Jar: newAuthCookieJar(t, reg, publicTS.URL, "i-do-not-exist")}
		makeRequestAndExpectError(t, hc, c, url.Values{"prompt": {"none"}}, "The Authorization Server requires End-User authentication.")
	})

	t.Run("case=should be able to retry accept consent request", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Subject: subject,
				Context: map[string]interface{}{"fooz": "barz"},
			}),
			checkAndDuplicateAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, map[string]interface{}{"fooz": "barz"}, res.Context)
				assert.Equal(t, subject, *res.Subject)
				return hydra.AcceptOAuth2ConsentRequest{
					Remember:   pointerx.Bool(true),
					GrantScope: []string{"openid"},
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}
			}))

		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})

	})

	t.Run("case=should be able to retry accept login request", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndDuplicateAcceptLoginHandler(t, adminClient, subject, func(*testing.T, *hydra.OAuth2LoginRequest, error) hydra.AcceptOAuth2LoginRequest {
				return hydra.AcceptOAuth2LoginRequest{
					Subject: subject,
					Context: map[string]interface{}{"fooz": "barz"},
				}
			}),
			checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, map[string]interface{}{"fooz": "barz"}, res.Context)
				assert.Equal(t, subject, *res.Subject)
				return hydra.AcceptOAuth2ConsentRequest{
					Remember:   pointerx.Bool(true),
					GrantScope: []string{"openid"},
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}
			}))

		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})
	})

	t.Run("case=should be able to retry both accept login and consent requests", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndDuplicateAcceptLoginHandler(t, adminClient, subject, func(*testing.T, *hydra.OAuth2LoginRequest, error) hydra.AcceptOAuth2LoginRequest {
				return hydra.AcceptOAuth2LoginRequest{
					Subject: subject,
					Context: map[string]interface{}{"fooz": "barz"},
				}
			}),
			checkAndDuplicateAcceptConsentHandler(t, adminClient, func(t *testing.T, res *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, map[string]interface{}{"fooz": "barz"}, res.Context)
				assert.Equal(t, subject, *res.Subject)
				return hydra.AcceptOAuth2ConsentRequest{
					Remember:   pointerx.Bool(true),
					GrantScope: []string{"openid"},
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}
			}))

		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})
	})
}

func TestStrategyDeviceLoginConsent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t)
	reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
	reg.Config().MustSet(ctx, config.KeyConsentRequestMaxAge, time.Hour)
	reg.Config().MustSet(ctx, config.KeyConsentRequestMaxAge, time.Hour)
	reg.Config().MustSet(ctx, config.KeyScopeStrategy, "exact")
	reg.Config().MustSet(ctx, config.KeySubjectTypesSupported, []string{"pairwise", "public"})
	reg.Config().MustSet(ctx, config.KeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")

	publicTS, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)
	adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
	adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

	oauth2Config := func(t *testing.T, c *client.Client) *oauth2.Config {
		return &oauth2.Config{
			ClientID:     c.GetID(),
			ClientSecret: c.Secret,
			Endpoint: oauth2.Endpoint{
				DeviceAuthURL: publicTS.URL + "/oauth2/device/auth",
				TokenURL:      publicTS.URL + "/oauth2/token",
				AuthStyle:     oauth2.AuthStyleInHeader,
			},
		}
	}

	now := 1723546027 // Unix timestamps must round-trip through Hydra without converting to floats or similar
	acceptDeviceHandler := func(t *testing.T) http.HandlerFunc {
		return checkAndAcceptDeviceHandler(t, adminClient)
	}

	acceptLoginHandler := func(t *testing.T, subject string, payload *hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
		return checkAndAcceptLoginHandler(t, adminClient, subject, func(*testing.T, *hydra.OAuth2LoginRequest, error) hydra.AcceptOAuth2LoginRequest {
			if payload == nil {
				return hydra.AcceptOAuth2LoginRequest{}
			}
			return *payload
		})
	}

	acceptConsentHandler := func(t *testing.T, payload *hydra.AcceptOAuth2ConsentRequest) http.HandlerFunc {
		return checkAndAcceptConsentHandler(t, adminClient, func(*testing.T, *hydra.OAuth2ConsentRequest, error) hydra.AcceptOAuth2ConsentRequest {
			if payload == nil {
				return hydra.AcceptOAuth2ConsentRequest{}
			}
			return *payload
		})
	}

	createDefaultClient := func(t *testing.T) *client.Client {
		c := &client.Client{GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}}
		return createClient(t, reg, c)
	}
	t.Run("case=should pass if both login and consent are granted and check remember flows as well as various payloads", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			acceptDeviceHandler(t),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Remember: pointerx.Bool(true),
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
				Remember:   pointerx.Bool(true),
				GrantScope: []string{"openid"},
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{
						"foo": "bar",
						"ts1": now,
					},
					IdToken: map[string]interface{}{
						"bar": "baz",
						"ts1": now,
					},
				},
			}))

		hc := testhelpers.NewEmptyJarClient(t)

		var sid string
		var run = func(t *testing.T) {
			res, resp := makeOAuth2DeviceAuthRequest(t, reg, hc, c, "openid")
			assert.EqualValues(t, http.StatusOK, resp.StatusCode)

			devResp := new(oauth2.DeviceAuthResponse)
			require.NoError(t, json.Unmarshal([]byte(res.Raw), devResp))

			resp, err := hc.Get(devResp.VerificationURIComplete)
			require.NoError(t, err)
			require.Contains(t, reg.Config().DeviceDoneURL(ctx).String(), resp.Request.URL.Path, "did not end up in post device URL")
			require.Equal(t, resp.Request.URL.Query().Get("client_id"), c.ID)

			conf := oauth2Config(t, c)
			token, err := conf.DeviceAccessToken(ctx, devResp)
			require.NoError(t, err)

			claims := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
			assert.Equal(t, "bar", claims.Get("ext.foo").String(), "%s", claims.Raw)

			idClaims := testhelpers.DecodeIDToken(t, token)
			assert.Equal(t, "baz", idClaims.Get("bar").String(), "%s", idClaims.Raw)
			sid = idClaims.Get("sid").String()
			assert.NotNil(t, sid)
		}

		t.Run("perform first flow", run)

		t.Run("perform follow up flows and check if session values are set", func(t *testing.T) {
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				checkAndAcceptLoginHandler(t, adminClient, subject, func(t *testing.T, res *hydra.OAuth2LoginRequest, err error) hydra.AcceptOAuth2LoginRequest {
					require.NoError(t, err)
					assert.True(t, res.Skip)
					assert.Equal(t, sid, *res.SessionId)
					assert.Equal(t, subject, res.Subject)
					assert.Empty(t, pointerx.StringR(res.Client.ClientSecret))
					return hydra.AcceptOAuth2LoginRequest{
						Subject: subject,
						Context: map[string]interface{}{"xyz": "abc"},
					}
				}),
				checkAndAcceptConsentHandler(t, adminClient, func(t *testing.T, req *hydra.OAuth2ConsentRequest, err error) hydra.AcceptOAuth2ConsentRequest {
					require.NoError(t, err)
					assert.True(t, *req.Skip)
					assert.Equal(t, sid, *req.LoginSessionId)
					assert.Equal(t, subject, *req.Subject)
					assert.Empty(t, pointerx.StringR(req.Client.ClientSecret))
					assert.Equal(t, map[string]interface{}{"xyz": "abc"}, req.Context)
					return hydra.AcceptOAuth2ConsentRequest{
						Remember:   pointerx.Bool(true),
						GrantScope: []string{"openid"},
						Session: &hydra.AcceptOAuth2ConsentRequestSession{
							AccessToken: map[string]interface{}{
								"foo": "bar",
								"ts1": now,
							},
							IdToken: map[string]interface{}{
								"bar": "baz",
								"ts2": now,
							},
						},
					}
				}))

			for k := 0; k < 3; k++ {
				t.Run(fmt.Sprintf("case=%d", k), run)
			}
		})
	})
	t.Run("case=should fail because we are reusing the same verifier", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(),
			acceptDeviceHandler(t),
			acceptLoginHandler(t, subject, &hydra.AcceptOAuth2LoginRequest{
				Remember: pointerx.Bool(true),
			}),
			acceptConsentHandler(t, &hydra.AcceptOAuth2ConsentRequest{
				Remember:   pointerx.Bool(true),
				GrantScope: []string{"openid"},
				Session: &hydra.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IdToken:     map[string]interface{}{"bar": "baz"},
				},
			}))

		hc := testhelpers.NewEmptyJarClient(t)

		res, resp := makeOAuth2DeviceAuthRequest(t, reg, hc, c, "openid")
		assert.EqualValues(t, http.StatusOK, resp.StatusCode)

		devResp := new(oauth2.DeviceAuthResponse)
		require.NoError(t, json.Unmarshal([]byte(res.Raw), devResp))

		resp, err := hc.Get(devResp.VerificationURIComplete)
		require.NoError(t, err)
		require.Contains(t, reg.Config().DeviceDoneURL(ctx).String(), resp.Request.URL.Path, "did not end up in post device URL")
		require.Equal(t, resp.Request.URL.Query().Get("client_id"), c.ID)

		conf := oauth2Config(t, c)
		token, err := conf.DeviceAccessToken(ctx, devResp)
		require.NoError(t, err)

		claims := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
		assert.Equal(t, "bar", claims.Get("ext.foo").String(), "%s", claims.Raw)

		idClaims := testhelpers.DecodeIDToken(t, token)
		assert.Equal(t, "baz", idClaims.Get("bar").String(), "%s", idClaims.Raw)
		sid := idClaims.Get("sid").String()
		assert.NotNil(t, sid)

	})
	t.Run("case=should fail because a device verifier was given that doesn't exist in the store", func(t *testing.T) {
		testhelpers.NewDeviceLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)

		_, res := makeOAuth2DeviceVerificationRequest(t, reg, hc, c, url.Values{"device_verifier": {"does-not-exist"}})
		assert.EqualValues(t, http.StatusForbidden, res.StatusCode)
	})

	t.Run("case=should fail because a login verifier was given that doesn't exist in the store", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)

		_, res := makeOAuth2DeviceVerificationRequest(t, reg, hc, c, url.Values{"login_verifier": {"does-not-exist"}})
		assert.EqualValues(t, http.StatusForbidden, res.StatusCode)
	})

	t.Run("case=should fail because a consent verifier was given that doesn't exist in the store", func(t *testing.T) {
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)
		hc := testhelpers.NewEmptyJarClient(t)

		_, res := makeOAuth2DeviceVerificationRequest(t, reg, hc, c, url.Values{"consent_verifier": {"does-not-exist"}})
		assert.EqualValues(t, http.StatusForbidden, res.StatusCode)
	})
}

func DropCookieJar(drop *regexp.Regexp) http.CookieJar {
	jar, _ := cookiejar.New(nil)
	return &dropCSRFCookieJar{
		jar:  jar,
		drop: drop,
	}
}

type dropCSRFCookieJar struct {
	jar  *cookiejar.Jar
	drop *regexp.Regexp
}

var _ http.CookieJar = (*dropCSRFCookieJar)(nil)

func (d *dropCSRFCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		if d.drop.MatchString(c.Name) {
			continue
		}
		d.jar.SetCookies(u, []*http.Cookie{c})
	}
}

func (d *dropCSRFCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return d.jar.Cookies(u)
}
