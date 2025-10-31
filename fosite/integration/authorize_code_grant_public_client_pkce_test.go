// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/oauth2" // "github.com/stretchr/testify/assert"
)

func TestAuthorizeCodeFlowWithPublicClientAndPKCE(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runAuthorizeCodeGrantWithPublicClientAndPKCETest(t, strategy)
	}
}

func runAuthorizeCodeGrantWithPublicClientAndPKCETest(t *testing.T, strategy interface{}) {
	c := new(fosite.Config)
	c.EnforcePKCE = true
	c.EnablePKCEPlainChallengeMethod = true
	f := compose.Compose(c, fositeStore, strategy, compose.OAuth2AuthorizeExplicitFactory, compose.OAuth2PKCEFactory, compose.OAuth2TokenIntrospectionFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	oauthClient.ClientSecret = ""
	oauthClient.ClientID = "public-client"
	fositeStore.Clients["public-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"

	var authCodeUrl string
	var verifier string
	for k, c := range []struct {
		description     string
		setup           func()
		authStatusCode  int
		tokenStatusCode int
	}{
		{
			description: "should fail because no challenge was given",
			setup: func() {
				authCodeUrl = oauthClient.AuthCodeURL("12345678901234567890")
			},
			authStatusCode: http.StatusNotAcceptable,
		},
		{
			description: "should pass",
			setup: func() {
				verifier = "somechallengesomechallengesomechallengesomechallengesomechallengesomechallenge"
				authCodeUrl = oauthClient.AuthCodeURL("12345678901234567890") + "&code_challenge=somechallengesomechallengesomechallengesomechallengesomechallengesomechallenge"
			},
			authStatusCode: http.StatusOK,
		},
		{
			description: "should fail because the verifier is mismatching",
			setup: func() {
				verifier = "failchallengefailchallengefailchallengefailchallengefailchallengefailchallengefailchallengefailchallenge"
				authCodeUrl = oauthClient.AuthCodeURL("12345678901234567890") + "&code_challenge=somechallengesomechallengesomechallengesomechallengesomechallengesomechallengesomechallengesomechallenge"
			},
			authStatusCode:  http.StatusOK,
			tokenStatusCode: http.StatusBadRequest,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			c.setup()

			t.Logf("Got url: %s", authCodeUrl)

			resp, err := http.Get(authCodeUrl)
			require.NoError(t, err)
			require.Equal(t, resp.StatusCode, c.authStatusCode)

			if resp.StatusCode == http.StatusOK {
				// This should fail because no verifier was given
				// _, err := oauthClient.Exchange(goauth.NoContext, resp.Request.URL.Query().Get("code"))
				// require.Error(t, err)
				// require.Empty(t, token.AccessToken)
				t.Logf("Got redirect url: %s", resp.Request.URL)

				resp, err := http.PostForm(ts.URL+"/token", url.Values{
					"code":          {resp.Request.URL.Query().Get("code")},
					"grant_type":    {"authorization_code"},
					"client_id":     {"public-client"},
					"redirect_uri":  {ts.URL + "/callback"},
					"code_verifier": {verifier},
				})
				require.NoError(t, err)
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				if c.tokenStatusCode != 0 {
					require.Equal(t, c.tokenStatusCode, resp.StatusCode)
					token := goauth.Token{}
					require.NoError(t, json.Unmarshal(body, &token))
					require.Empty(t, token.AccessToken)
					return
				}

				assert.Equal(t, resp.StatusCode, http.StatusOK)
				token := goauth.Token{}
				require.NoError(t, json.Unmarshal(body, &token))

				require.NotEmpty(t, token.AccessToken, "Got body: %s", string(body))

				httpClient := oauthClient.Client(goauth.NoContext, &token)
				resp, err = httpClient.Get(ts.URL + "/info")
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
		})
	}
}
