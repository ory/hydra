// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
)

func TestAuthorizeImplicitFlow(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runTestAuthorizeImplicitGrant(t, strategy)
	}
}

func runTestAuthorizeImplicitGrant(t *testing.T, strategy interface{}) {
	f := compose.Compose(new(fosite.Config), fositeStore, strategy, compose.OAuth2AuthorizeImplicitFactory, compose.OAuth2TokenIntrospectionFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"

	var state string
	for k, c := range []struct {
		description    string
		setup          func()
		check          func(t *testing.T, r *http.Response)
		params         []goauth.AuthCodeOption
		authStatusCode int
	}{
		{
			description: "should fail because of audience",
			params:      []goauth.AuthCodeOption{goauth.SetAuthURLParam("audience", "https://www.ory.sh/not-api")},
			setup: func() {
				state = "12345678901234567890"
			},
			authStatusCode: http.StatusNotAcceptable,
		},
		{
			description: "should fail because of scope",
			params:      []goauth.AuthCodeOption{},
			setup: func() {
				oauthClient.Scopes = []string{"not-exist"}
				state = "12345678901234567890"
			},
			authStatusCode: http.StatusNotAcceptable,
		},
		{
			description: "should pass with proper audience",
			params:      []goauth.AuthCodeOption{goauth.SetAuthURLParam("audience", "https://www.ory.sh/api")},
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"fosite"}
			},
			check: func(t *testing.T, r *http.Response) {
				var b fosite.AccessRequest
				b.Client = new(fosite.DefaultClient)
				b.Session = new(defaultSession)
				require.NoError(t, json.NewDecoder(r.Body).Decode(&b))
				assert.EqualValues(t, fosite.Arguments{"https://www.ory.sh/api"}, b.RequestedAudience)
				assert.EqualValues(t, fosite.Arguments{"https://www.ory.sh/api"}, b.GrantedAudience)
				assert.EqualValues(t, "foo-sub", b.Session.(*defaultSession).Subject)
			},
			authStatusCode: http.StatusOK,
		},
		{
			description: "should pass",
			setup: func() {
				state = "12345678901234567890"
			},
			authStatusCode: http.StatusOK,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			c.setup()

			var callbackURL *url.URL
			authURL := strings.Replace(oauthClient.AuthCodeURL(state, c.params...), "response_type=code", "response_type=token", -1)
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					callbackURL = req.URL
					return errors.New("Dont follow redirects")
				},
			}
			resp, err := client.Get(authURL)
			require.Error(t, err)

			if resp.StatusCode == http.StatusOK {
				fragment, err := url.ParseQuery(callbackURL.Fragment)
				require.NoError(t, err)
				expires, err := strconv.Atoi(fragment.Get("expires_in"))
				require.NoError(t, err)
				token := &goauth.Token{
					AccessToken:  fragment.Get("access_token"),
					TokenType:    fragment.Get("token_type"),
					RefreshToken: fragment.Get("refresh_token"),
					Expiry:       time.Now().UTC().Add(time.Duration(expires) * time.Second),
				}

				httpClient := oauthClient.Client(goauth.NoContext, token)
				resp, err := httpClient.Get(ts.URL + "/info")
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				if c.check != nil {
					c.check(t, resp)
				}
			}
		})
	}
}
