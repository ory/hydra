package consent_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/ory/x/pointerx"

	"github.com/tidwall/gjson"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal/testhelpers"

	"github.com/ory/fosite"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
)

func TestStrategyLoginConsentNext(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	reg.Config().MustSet(config.KeyAccessTokenStrategy, "opaque")
	reg.Config().MustSet(config.KeyConsentRequestMaxAge, time.Hour)
	reg.Config().MustSet(config.KeyConsentRequestMaxAge, time.Hour)
	reg.Config().MustSet(config.KeyScopeStrategy, "exact")
	reg.Config().MustSet(config.KeySubjectTypesSupported, []string{"pairwise", "public"})
	reg.Config().MustSet(config.KeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")

	publicTS, adminTS := testhelpers.NewOAuth2Server(t, reg)
	adminClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(adminTS.URL).Host})

	oauth2Config := func(t *testing.T, c *client.Client) *oauth2.Config {
		return &oauth2.Config{
			ClientID:     c.OutfacingID,
			ClientSecret: c.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   publicTS.URL + "/oauth2/auth",
				TokenURL:  publicTS.URL + "/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
			RedirectURL: c.RedirectURIs[0],
		}
	}

	acceptLoginHandler := func(t *testing.T, subject string, payload *models.AcceptLoginRequest) http.HandlerFunc {
		return checkAndAcceptLoginHandler(t, adminClient.Admin, subject, func(*testing.T, *admin.GetLoginRequestOK, error) *models.AcceptLoginRequest {
			if payload == nil {
				return new(models.AcceptLoginRequest)
			}
			return payload
		})
	}

	acceptConsentHandler := func(t *testing.T, payload *models.AcceptConsentRequest) http.HandlerFunc {
		return checkAndAcceptConsentHandler(t, adminClient.Admin, func(*testing.T, *admin.GetConsentRequestOK, error) *models.AcceptConsentRequest {
			if payload == nil {
				return new(models.AcceptConsentRequest)
			}
			return payload
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

		makeRequestAndExpectError(t, nil, c, url.Values{"login_verifier": {"does-not-exist"}}, "The login verifier has already been used, has not been granted, or is invalid.")
	})

	t.Run("case=should fail because a non-existing consent verifier was given", func(t *testing.T) {
		// Covers:
		// - This should fail because consent verifier was set but does not exist
		// - This should fail because a consent verifier was given but no login verifier
		testhelpers.NewLoginConsentUI(t, reg.Config(), testhelpers.HTTPServerNoExpectedCallHandler(t), testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)
		makeRequestAndExpectError(t, nil, c, url.Values{"consent_verifier": {"does-not-exist"}}, "The consent verifier has already been used, has not been granted, or is invalid.")
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
			vr, err := adminClient.Admin.RejectLoginRequest(admin.NewRejectLoginRequestParams().
				WithLoginChallenge(r.URL.Query().Get("login_challenge")).
				WithBody(&models.RejectRequest{
					Error:            fosite.ErrInteractionRequired.ErrorField,
					ErrorDescription: "expect-reject-login",
					StatusCode:       int64(fosite.ErrInteractionRequired.CodeField),
				}))
			require.NoError(t, err)
			assert.NotEmpty(t, *vr.Payload.RedirectTo)
			http.Redirect(w, r, *vr.Payload.RedirectTo, http.StatusFound)
		}, testhelpers.HTTPServerNoExpectedCallHandler(t))
		c := createDefaultClient(t)

		makeRequestAndExpectError(t, nil, c, url.Values{}, "expect-reject-login")
	})

	t.Run("case=should fail because no cookie jar invalid csrf", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", nil),
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		hc := new(http.Client)
		makeRequestAndExpectError(t, hc, c, url.Values{}, "No CSRF value available in the session cookie.")
	})

	t.Run("case=should fail because consent endpoints denies the request after login was granted", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			func(w http.ResponseWriter, r *http.Request) {
				vr, err := adminClient.Admin.RejectConsentRequest(
					admin.NewRejectConsentRequestParams().WithConsentChallenge(r.URL.Query().Get("consent_challenge")).
						WithBody(&models.RejectRequest{
							Error:            fosite.ErrInteractionRequired.ErrorField,
							ErrorDescription: "expect-reject-consent",
							StatusCode:       int64(fosite.ErrInteractionRequired.CodeField)}))
				require.NoError(t, err)
				require.NotEmpty(t, *vr.Payload.RedirectTo)
				http.Redirect(w, r, *vr.Payload.RedirectTo, http.StatusFound)
			})

		makeRequestAndExpectError(t, nil, c, url.Values{}, "expect-reject-consent")
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
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{
				Remember: true,
			}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{
				Remember:   true,
				GrantScope: []string{"openid"},
				Session: &models.ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": "bar"},
					IDToken:     map[string]interface{}{"bar": "baz"},
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

			claims := testhelpers.IntrospectToken(t, conf, token, adminTS)
			assert.Equal(t, "bar", claims.Get("ext.foo").String(), "%s", claims.Raw)

			idClaims := testhelpers.DecodeIDToken(t, token)
			assert.Equal(t, "baz", idClaims.Get("bar").String(), "%s", idClaims.Raw)
			sid = idClaims.Get("sid").String()
			assert.NotNil(t, sid)
		}

		t.Run("perform first flow", run)

		t.Run("perform follow up flows and check if session values are set", func(t *testing.T) {
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				checkAndAcceptLoginHandler(t, adminClient.Admin, subject, func(t *testing.T, res *admin.GetLoginRequestOK, err error) *models.AcceptLoginRequest {
					require.NoError(t, err)
					assert.True(t, *res.Payload.Skip)
					assert.Equal(t, sid, res.Payload.SessionID)
					assert.Equal(t, subject, *res.Payload.Subject)
					assert.Empty(t, res.Payload.Client.ClientSecret)
					return &models.AcceptLoginRequest{
						Subject: &subject,
						Context: map[string]interface{}{"foo": "bar"},
					}
				}),
				checkAndAcceptConsentHandler(t, adminClient.Admin, func(t *testing.T, res *admin.GetConsentRequestOK, err error) *models.AcceptConsentRequest {
					require.NoError(t, err)
					assert.True(t, res.Payload.Skip)
					assert.Equal(t, sid, res.Payload.LoginSessionID)
					assert.Equal(t, subject, res.Payload.Subject)
					assert.Empty(t, res.Payload.Client.ClientSecret)
					return &models.AcceptConsentRequest{
						Remember:   true,
						GrantScope: []string{"openid"},
						Session: &models.ConsentRequestSession{
							AccessToken: map[string]interface{}{"foo": "bar"},
							IDToken:     map[string]interface{}{"bar": "baz"},
						},
					}
				}))

			for k := 0; k < 3; k++ {
				t.Run(fmt.Sprintf("case=%d", k), run)
			}
		})
	})

	t.Run("case=should pass and check if login context is set properly", func(t *testing.T) {
		// This should pass because login was remembered and session id should be set and session context should also work
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{
				Subject: &subject,
				Context: map[string]interface{}{"fooz": "barz"},
			}),
			checkAndAcceptConsentHandler(t, adminClient.Admin, func(t *testing.T, res *admin.GetConsentRequestOK, err error) *models.AcceptConsentRequest {
				require.NoError(t, err)
				assert.Equal(t, map[string]interface{}{"fooz": "barz"}, res.Payload.Context)
				assert.Equal(t, subject, res.Payload.Subject)
				return &models.AcceptConsentRequest{
					Remember:   true,
					GrantScope: []string{"openid"},
					Session: &models.ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IDToken:     map[string]interface{}{"bar": "baz"},
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

		c := &client.Client{OutfacingID: uuid.New(), TokenEndpointAuthMethod: "none",
			RedirectURIs: []string{
				testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler),
				"custom://redirection-scheme/path",
				"custom://localhost/path",
			}}
		require.NoError(t, reg.ClientManager().CreateClient(context.Background(), c))

		subject := "aeneas-rekkas"
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true, RememberFor: 0}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true, RememberFor: 0}))

		hc := testhelpers.NewEmptyJarClient(t)

		t.Run("set up initial session", func(t *testing.T) {
			makeRequestAndExpectCode(t, hc, c, url.Values{"redirect_uri": {c.RedirectURIs[0]}})
		})

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
					_, err := hc.Get(urlx.CopyWithQuery(reg.Config().OAuth2AuthURL(), url.Values{
						"response_type": {"code"},
						"state":         {uuid.New()},
						"redirect_uri":  {redir},
						"client_id":     {c.OutfacingID},
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

	t.Run("case=should fail at login screen because subject in login challenge does not match subject from previous session", func(t *testing.T) {
		// Previously: This should fail at login screen because subject from accept does not match subject from session
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, nil))

		// Init session
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			func(w http.ResponseWriter, r *http.Request) {
				_, err := adminClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
					WithLoginChallenge(r.URL.Query().Get("login_challenge")).
					WithBody(&models.AcceptLoginRequest{
						Subject: pointerx.String("not-aeneas-rekkas"),
					}))
				require.Error(t, err)
				assert.Contains(t, err.Error(), "Field 'subject' does not match subject from previous authentication")
				w.WriteHeader(http.StatusBadRequest)
			},
			testhelpers.HTTPServerNoExpectedCallHandler(t))

		_, res := makeOAuth2Request(t, reg, hc, c, url.Values{})
		assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Request.URL.Query().Get("code"))
	})

	t.Run("case=should require re-authentication when parameters mandate it", func(t *testing.T) {
		// Covers:
		// - should pass and require re-authentication although session is set (because prompt=login)
		// - should pass and require re-authentication although session is set (because max_age=1)
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		resetUI := func(t *testing.T) {
			testhelpers.NewLoginConsentUI(t, reg.Config(),
				checkAndAcceptLoginHandler(t, adminClient.Admin, subject, func(t *testing.T, res *admin.GetLoginRequestOK, err error) *models.AcceptLoginRequest {
					require.NoError(t, err)
					assert.False(t, *res.Payload.Skip) // Skip should always be false here
					return &models.AcceptLoginRequest{
						Remember: true,
					}
				}),
				acceptConsentHandler(t, &models.AcceptConsentRequest{
					Remember: true,
				}))
		}
		resetUI(t)

		// Init session
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		for k, values := range []url.Values{
			{"prompt": {"login"}},
			{"max_age": {"1"}},
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

		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		time.Sleep(time.Second)

		makeRequestAndExpectError(t, hc, c, url.Values{"max_age": {"1"}, "prompt": {"none"}},
			"prompt is set to 'none' and authentication time reached 'max_age'")
	})

	t.Run("case=should fail because prompt is none but no auth session exists", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

		makeRequestAndExpectError(t, nil, c, url.Values{"prompt": {"none"}},
			"Prompt 'none' was requested, but no existing login session was found")
	})

	t.Run("case=should fail because prompt is none and consent is missing a permission which requires re-authorization of the app", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(), acceptLoginHandler(t, "aeneas-rekkas", &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

		// Init cookie
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		// Make request with additional scope and prompt none, which fails
		makeRequestAndExpectError(t, hc, c, url.Values{"prompt": {"none"}, "scope": {"openid"}},
			"Prompt 'none' was requested, but no previous consent was found")
	})

	t.Run("case=pass and properly require authentication as well as authorization because prompt is set to login and consent although previous session exists", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			checkAndAcceptLoginHandler(t, adminClient.Admin, subject, func(t *testing.T, res *admin.GetLoginRequestOK, err error) *models.AcceptLoginRequest {
				require.NoError(t, err)
				assert.False(t, *res.Payload.Skip) // Skip should always be false here because prompt has login
				return &models.AcceptLoginRequest{Remember: true}
			}),
			checkAndAcceptConsentHandler(t, adminClient.Admin, func(t *testing.T, res *admin.GetConsentRequestOK, err error) *models.AcceptConsentRequest {
				require.NoError(t, err)
				assert.False(t, res.Payload.Skip) // Skip should always be false here because prompt has consent
				return &models.AcceptConsentRequest{
					Remember: true,
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
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

		// Init cookie
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		for _, values := range []url.Values{
			{"prompt": {"none"}, "id_token_hint": {testhelpers.NewIDToken(t, reg, notSubject)}},
			{"id_token_hint": {testhelpers.NewIDToken(t, reg, notSubject)}},
		} {
			t.Run(fmt.Sprintf("values=%v", values), func(t *testing.T) {
				testhelpers.NewLoginConsentUI(t, reg.Config(),
					checkAndAcceptLoginHandler(t, adminClient.Admin, subject, func(t *testing.T, res *admin.GetLoginRequestOK, err error) *models.AcceptLoginRequest {
						var b bytes.Buffer
						require.NoError(t, json.NewEncoder(&b).Encode(res.Payload))
						assert.EqualValues(t, notSubject, gjson.GetBytes(b.Bytes(), "oidc_context.id_token_hint_claims.sub"), b.String())
						return &models.AcceptLoginRequest{Remember: true}
					}),
					acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

				makeRequestAndExpectError(t, hc, c, values,
					"Request failed because subject claim from id_token_hint does not match subject from authentication session")
			})
		}
	})

	t.Run("case=should pass and require authentication because id_token_hint does match subject from session", func(t *testing.T) {
		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

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
			sha256.Sum256([]byte(c.SectorIdentifierURI+subject+reg.Config().SubjectIdentifierAlgorithmSalt())))

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true, GrantScope: []string{"openid"}}))

		for _, tc := range []struct {
			d      string
			values url.Values
		}{
			{
				d:      "check all the sub claims",
				values: url.Values{"scope": {"openid"}},
			},
			{
				d:      "works with id_token_hint",
				values: url.Values{"scope": {"openid"}, "id_token_hint": {testhelpers.NewIDToken(t, reg, hash)}},
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
				atClaims := testhelpers.IntrospectToken(t, conf, token, adminTS)
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
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{
				ForceSubjectIdentifier: obfuscated,
			}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{GrantScope: []string{"openid"}}))

		code := makeRequestAndExpectCode(t, nil, c, url.Values{})

		conf := oauth2Config(t, c)
		token, err := conf.Exchange(context.Background(), code)
		require.NoError(t, err)

		// OpenID data must be obfuscated
		idClaims := testhelpers.DecodeIDToken(t, token)
		assert.EqualValues(t, obfuscated, idClaims.Get("sub").String())
		uiClaims := testhelpers.Userinfo(t, token, publicTS)
		assert.EqualValues(t, obfuscated, uiClaims.Get("sub").String())

		// Access token data must not be obfuscated
		atClaims := testhelpers.IntrospectToken(t, conf, token, adminTS)
		assert.EqualValues(t, subject, atClaims.Get("sub").String())
	})

	t.Run("suite=properly clean up session cookies", func(t *testing.T) {
		t.Skip("This test is skipped because we forcibly set remember to true always when skip is also true for a better user experience.")

		subject := "aeneas-rekkas"
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: true}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: true}))

		// Initialize flow
		// Formerly: This should pass as regularly and create a new session and forward data
		hc := testhelpers.NewEmptyJarClient(t)
		makeRequestAndExpectCode(t, hc, c, url.Values{})

		// Re-run flow but do not remember login
		// Formerly: This should pass and also revoke the session cookie
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, subject, &models.AcceptLoginRequest{Remember: false}),
			acceptConsentHandler(t, &models.AcceptConsentRequest{Remember: false}))
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
}
