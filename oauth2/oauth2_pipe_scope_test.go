// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/ioutilx"
	"github.com/ory/x/pointerx"

	"github.com/go-jose/go-jose/v3"
)

// TestOAuth2AuthCodeWithPipeCharactersInScopes tests that scopes containing pipe characters
// (e.g., "abc|def") are handled correctly throughout the OAuth2 flow. This is required for
// ONC g(10) certification and FHIR OAuth2 compliance.
// See: https://build.fhir.org/ig/HL7/smart-app-launch/scopes-and-launch-context.html#finer-grained-resource-constraints-using-search-parameters
func TestOAuth2AuthCodeWithPipeCharactersInScopes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for dbName, reg := range testhelpers.ConnectDatabases(t, true, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeyAccessTokenStrategy: "opaque",
		config.KeyRefreshTokenHook:    "",
	}))) {
		t.Run("registry="+dbName, func(t *testing.T) {
			t.Parallel()

			jwk.EnsureAsymmetricKeypairExists(t, reg, string(jose.ES256), x.OpenIDConnectKeyName)
			jwk.EnsureAsymmetricKeypairExists(t, reg, string(jose.ES256), x.OAuth2JWTKeyName)

			publicTS, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)

			adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
			adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

			subject := "test-subject"
			scopeWithPipe := "openid profile patient|read patient|write"
			scopeParts := []string{"openid", "profile", "patient|read", "patient|write"}

			t.Run("case=perform authorize code flow with scopes containing pipe characters", func(t *testing.T) {
				c, conf := newOAuth2Client(
					t,
					reg,
					testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler),
					withScope(scopeWithPipe),
				)

				testhelpers.NewLoginConsentUI(t, reg.Config(),
					func(w http.ResponseWriter, r *http.Request) {
						rr, res, err := adminClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
						require.NoErrorf(t, err, "%s\n%s", res.Request.URL, ioutilx.MustReadAll(res.Body))

						// Verify the login request contains the correct scopes
						assert.ElementsMatch(t, scopeParts, rr.RequestedScope)

						acceptBody := hydra.AcceptOAuth2LoginRequest{
							Subject:  subject,
							Remember: pointerx.Ptr(true),
							Acr:      pointerx.Ptr("1"),
							Amr:      []string{"pwd"},
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
						challenge := r.URL.Query().Get("consent_challenge")
						rr, _, err := adminClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(challenge).Execute()
						require.NoError(t, err)

						// Verify the consent request contains the correct scopes
						assert.ElementsMatch(t, scopeParts, rr.RequestedScope)

						acceptBody := hydra.AcceptOAuth2ConsentRequest{
							GrantScope:               scopeParts,
							GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
							Remember:                 pointerx.Ptr(true),
							RememberFor:              pointerx.Ptr[int64](0),
						}

						v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
							ConsentChallenge(challenge).
							AcceptOAuth2ConsentRequest(acceptBody).
							Execute()
						require.NoError(t, err)
						require.NotEmpty(t, v.RedirectTo)
						http.Redirect(w, r, v.RedirectTo, http.StatusFound)
					},
				)

				code, _ := getAuthorizeCode(t, conf, nil)
				require.NotEmpty(t, code)

				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)
				require.NotEmpty(t, token.AccessToken)

				// Introspect the access token to verify scopes are preserved
				introspect := testhelpers.IntrospectToken(t, token.AccessToken, adminTS)
				assert.True(t, introspect.Get("active").Bool(), "%s", introspect)
				assert.EqualValues(t, conf.ClientID, introspect.Get("client_id").String(), "%s", introspect)
				assert.EqualValues(t, subject, introspect.Get("sub").String(), "%s", introspect)

				// Verify that the scope field contains the correct space-separated scopes
				scopes := introspect.Get("scope").String()
				assert.NotEmpty(t, scopes)

				// Split by space to get individual scopes
				actualScopes := strings.Split(scopes, " ")
				assert.ElementsMatch(t, scopeParts, actualScopes, "Scopes should be preserved as space-separated, with pipe characters intact. Expected: %v, got: %v", scopeParts, actualScopes)

				t.Run("followup=verify refresh token preserves pipe scopes", func(t *testing.T) {
					require.NotEmpty(t, token.RefreshToken)
					token.Expiry = token.Expiry.Add(-time.Hour * 24)

					refreshedToken, err := conf.TokenSource(context.Background(), token).Token()
					require.NoError(t, err)

					// Introspect the refreshed access token
					refreshedIntrospect := testhelpers.IntrospectToken(t, refreshedToken.AccessToken, adminTS)
					assert.True(t, refreshedIntrospect.Get("active").Bool(), "%s", refreshedIntrospect)
					assert.EqualValues(t, conf.ClientID, refreshedIntrospect.Get("client_id").String(), "%s", refreshedIntrospect)
					assert.EqualValues(t, subject, refreshedIntrospect.Get("sub").String(), "%s", refreshedIntrospect)

					// Verify that the scope field still contains the correct scopes with pipe characters
					refreshedScopes := refreshedIntrospect.Get("scope").String()
					assert.NotEmpty(t, refreshedScopes)
					actualRefreshedScopes := strings.Split(refreshedScopes, " ")
					assert.ElementsMatch(t, scopeParts, actualRefreshedScopes, "Refreshed scopes should preserve pipe characters. Expected: %v, got: %v", scopeParts, actualRefreshedScopes)
				})
			})

			t.Run("case=verify JWT access token preserves pipe scopes", func(t *testing.T) {
				reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
				t.Cleanup(func() {
					reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
				})

				c, conf := newOAuth2Client(
					t,
					reg,
					testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler),
					withScope(scopeWithPipe),
				)

				testhelpers.NewLoginConsentUI(t, reg.Config(),
					func(w http.ResponseWriter, r *http.Request) {
						rr, res, err := adminClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
						require.NoErrorf(t, err, "%s\n%s", res.Request.URL, ioutilx.MustReadAll(res.Body))

						acceptBody := hydra.AcceptOAuth2LoginRequest{
							Subject:  subject,
							Remember: pointerx.Ptr(true),
							Acr:      pointerx.Ptr("1"),
							Amr:      []string{"pwd"},
						}

						v, _, err := adminClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
							LoginChallenge(r.URL.Query().Get("login_challenge")).
							AcceptOAuth2LoginRequest(acceptBody).
							Execute()
						require.NoError(t, err)
						http.Redirect(w, r, v.RedirectTo, http.StatusFound)
					},
					func(w http.ResponseWriter, r *http.Request) {
						challenge := r.URL.Query().Get("consent_challenge")
						rr, _, err := adminClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(challenge).Execute()
						require.NoError(t, err)

						acceptBody := hydra.AcceptOAuth2ConsentRequest{
							GrantScope:               scopeParts,
							GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
							Remember:                 pointerx.Ptr(true),
							RememberFor:              pointerx.Ptr[int64](0),
						}

						v, _, err := adminClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
							ConsentChallenge(challenge).
							AcceptOAuth2ConsentRequest(acceptBody).
							Execute()
						require.NoError(t, err)
						http.Redirect(w, r, v.RedirectTo, http.StatusFound)
					},
				)

				code, _ := getAuthorizeCode(t, conf, nil)
				require.NotEmpty(t, code)

				token, err := conf.Exchange(context.Background(), code)
				require.NoError(t, err)
				require.NotEmpty(t, token.AccessToken)

				// Decode JWT access token
				parts := strings.Split(token.AccessToken, ".")
				require.Len(t, parts, 3, "JWT should have 3 parts")

				claims := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(t, token.AccessToken))

				// Verify scopes are preserved in JWT claims
				// The "scp" claim should be an array with pipe characters intact
				scpArray := claims.Get("scp").Array()
				require.NotEmpty(t, scpArray, "JWT should contain 'scp' claim as an array")

				var actualScopes []string
				for _, s := range scpArray {
					actualScopes = append(actualScopes, s.String())
				}

				assert.ElementsMatch(t, scopeParts, actualScopes, "JWT 'scp' claim should preserve pipe characters. Expected: %v, got: %v", scopeParts, actualScopes)

				// Also verify the space-separated "scope" claim
				scopeStr := claims.Get("scope").String()
				assert.NotEmpty(t, scopeStr, "JWT should contain 'scope' claim as a string")
				actualScopesStr := strings.Split(scopeStr, " ")
				assert.ElementsMatch(t, scopeParts, actualScopesStr, "JWT 'scope' claim should be space-separated with pipe characters intact. Expected: %v, got: %v", scopeParts, actualScopesStr)
			})
		})
	}
}
