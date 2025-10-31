// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
)

func TestPushedAuthorizeCodeFlow(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runPushedAuthorizeCodeGrantTest(t, strategy)
	}
}

func runPushedAuthorizeCodeGrantTest(t *testing.T, strategy interface{}) {
	f := compose.Compose(new(fosite.Config), fositeStore, strategy, compose.OAuth2AuthorizeExplicitFactory, compose.OAuth2TokenIntrospectionFactory, compose.PushedAuthorizeHandlerFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{Subject: "foo-sub"})
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"

	var state string
	for k, c := range []struct {
		description    string
		setup          func()
		check          func(t *testing.T, r *http.Response)
		params         map[string]string
		authStatusCode int
		parStatusCode  int
	}{
		{
			description: "should fail because of audience",
			params:      map[string]string{"audience": "https://www.ory.sh/not-api"},
			setup: func() {
				oauthClient = newOAuth2Client(ts)
				state = "12345678901234567890"
			},
			parStatusCode:  http.StatusBadRequest,
			authStatusCode: http.StatusNotAcceptable,
		},
		{
			description: "should fail because of scope",
			params:      nil,
			setup: func() {
				oauthClient = newOAuth2Client(ts)
				oauthClient.Scopes = []string{"not-exist"}
				state = "12345678901234567890"
			},
			parStatusCode:  http.StatusBadRequest,
			authStatusCode: http.StatusNotAcceptable,
		},
		{
			description: "should pass with proper audience",
			params:      map[string]string{"audience": "https://www.ory.sh/api"},
			setup: func() {
				oauthClient = newOAuth2Client(ts)
				state = "12345678901234567890"
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
			parStatusCode:  http.StatusCreated,
			authStatusCode: http.StatusOK,
		},
		{
			description: "should pass",
			setup: func() {
				oauthClient = newOAuth2Client(ts)
				state = "12345678901234567890"
			},
			parStatusCode:  http.StatusCreated,
			authStatusCode: http.StatusOK,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			c.setup()

			// build request from the OAuth client
			data := url.Values{}
			data.Set("client_id", oauthClient.ClientID)
			data.Set("client_secret", oauthClient.ClientSecret)
			data.Set("response_type", "code")
			data.Set("state", state)
			data.Set("scope", strings.Join(oauthClient.Scopes, " "))
			data.Set("redirect_uri", oauthClient.RedirectURL)
			for k, v := range c.params {
				data.Set(k, v)
			}

			req, err := http.NewRequest("POST", ts.URL+"/par", strings.NewReader(data.Encode()))
			require.NoError(t, err)

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			resp, err := http.DefaultClient.Do(req)

			require.NoError(t, err)

			body, err := checkStatusAndGetBody(t, resp, c.parStatusCode)
			require.NoError(t, err, "Unable to get body after PAR. Err=%v", err)

			if resp.StatusCode != http.StatusCreated {
				return
			}

			m := map[string]interface{}{}
			err = json.Unmarshal(body, &m)

			assert.NoError(t, err, "Error occurred when unamrshaling the body: %v", err)

			// validate request_uri
			requestURI, _ := m["request_uri"].(string)
			assert.NotEmpty(t, requestURI, "request_uri is empty")
			assert.Condition(t, func() bool {
				return strings.HasPrefix(requestURI, "urn:ietf:params:oauth:request_uri:")
			}, "PAR Prefix is incorrect: %s", requestURI)

			// validate expires_in
			assert.EqualValues(t, 300, int(m["expires_in"].(float64)), "Invalid expires_in value=%v", m["expires_in"])

			// call authorize
			data = url.Values{}
			data.Set("client_id", oauthClient.ClientID)
			data.Set("request_uri", m["request_uri"].(string))
			req, err = http.NewRequest("POST", ts.URL+"/auth", strings.NewReader(data.Encode()))
			require.NoError(t, err)

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, c.authStatusCode, resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				return
			}

			require.NotEmpty(t, resp.Request.URL.Query().Get("code"), "Auth code is empty")

			token, err := oauthClient.Exchange(goauth.NoContext, resp.Request.URL.Query().Get("code"))
			require.NoError(t, err)
			require.NotEmpty(t, token.AccessToken)

			httpClient := oauthClient.Client(goauth.NoContext, token)
			resp, err = httpClient.Get(ts.URL + "/info")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			if c.check != nil {
				c.check(t, resp)
			}
		})
	}
}

func checkStatusAndGetBody(t *testing.T, resp *http.Response, expectedStatusCode int) ([]byte, error) {
	defer resp.Body.Close()

	require.Equal(t, expectedStatusCode, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	if err == nil {
		fmt.Printf("PAR response: body=%s\n", string(b))
	}
	if expectedStatusCode != resp.StatusCode {
		return nil, fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}

	return b, err
}
