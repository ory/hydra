// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/token/jwt"

	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
)

func TestAuthorizeResponseModes(t *testing.T) {
	session := &defaultSession{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Subject: "peter",
			},
			Headers: &jwt.Headers{},
		},
	}
	f := compose.ComposeAllEnabled(&fosite.Config{
		UseLegacyErrorFormat: true,
		GlobalSecret:         []byte("some-secret-thats-random-some-secret-thats-random-"),
	}, fositeStore, gen.MustRSAKey())
	ts := mockServer(t, f, session)
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	defaultClient := fositeStore.Clients["my-client"].(*fosite.DefaultClient)
	defaultClient.RedirectURIs[0] = ts.URL + "/callback"
	responseModeClient := &fosite.DefaultResponseModeClient{
		DefaultClient: defaultClient,
		ResponseModes: []fosite.ResponseModeType{},
	}
	fositeStore.Clients["response-mode-client"] = responseModeClient
	oauthClient.ClientID = "response-mode-client"

	var state string
	for k, c := range []struct {
		description  string
		setup        func()
		check        func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string)
		responseType string
		responseMode string
	}{
		{
			description:  "Should give err because implicit grant with response mode query",
			responseType: "id_token%20token",
			responseMode: "query",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeQuery}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				assert.NotEmpty(t, err["ErrorField"])
				assert.NotEmpty(t, err["DescriptionField"])
				assert.Equal(t, "Insecure response_mode 'query' for the response_type '[id_token token]'.", err["HintField"])
			},
		},
		{
			description:  "Should pass implicit grant with response mode form_post",
			responseType: "id_token%20token",
			responseMode: "form_post",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeFormPost}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				assert.EqualValues(t, state, stateFromServer)
				assert.NotEmpty(t, token.TokenType)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.Expiry)
				assert.NotEmpty(t, iDToken)
			},
		},
		{
			description:  "Should fail because response mode form_post is not allowed by the client",
			responseType: "id_token%20token",
			responseMode: "form_post",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeQuery}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				assert.NotEmpty(t, err["ErrorField"])
				assert.NotEmpty(t, err["DescriptionField"])
				assert.Equal(t, "The client is not allowed to request response_mode 'form_post'.", err["HintField"])
			},
		},
		{
			description:  "Should fail because response mode form_post is not allowed by the client without legacy format",
			responseType: "id_token%20token",
			responseMode: "form_post",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeQuery}
				f.(*fosite.Fosite).Config.(*fosite.Config).UseLegacyErrorFormat = false
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				f.(*fosite.Fosite).Config.(*fosite.Config).UseLegacyErrorFormat = true // reset
				assert.NotEmpty(t, err["ErrorField"])
				assert.Contains(t, err["DescriptionField"], "The client is not allowed to request response_mode 'form_post'.")
				assert.Empty(t, err["HintField"])
			},
		},
		{
			description:  "Should pass Authorization code grant test with response mode fragment",
			responseType: "code",
			responseMode: "fragment",
			setup: func() {
				state = "12345678901234567890"
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeFragment}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				assert.EqualValues(t, state, stateFromServer)
				assert.NotEmpty(t, code)
			},
		},
		{
			description:  "Should pass Authorization code grant test with response mode form_post",
			responseType: "code",
			responseMode: "form_post",
			setup: func() {
				state = "12345678901234567890"
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeFormPost}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				assert.EqualValues(t, state, stateFromServer)
				assert.NotEmpty(t, code)
			},
		},
		{
			description:  "Should fail Hybrid grant test with query",
			responseType: "token%20code",
			responseMode: "query",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeQuery}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				//assert.EqualValues(t, state, stateFromServer)
				assert.NotEmpty(t, err["ErrorField"])
				assert.NotEmpty(t, err["DescriptionField"])
				assert.Equal(t, "Insecure response_mode 'query' for the response_type '[token code]'.", err["HintField"])
			},
		},
		{
			description:  "Should fail Hybrid grant test with query without legacy fields",
			responseType: "token%20code",
			responseMode: "query",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeQuery}
				f.(*fosite.Fosite).Config.(*fosite.Config).UseLegacyErrorFormat = false
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				f.(*fosite.Fosite).Config.(*fosite.Config).UseLegacyErrorFormat = true // reset

				//assert.EqualValues(t, state, stateFromServer)
				assert.NotEmpty(t, err["ErrorField"])
				assert.Contains(t, err["DescriptionField"], "Insecure response_mode 'query' for the response_type '[token code]'.")
				assert.Empty(t, err["HintField"])
				assert.Empty(t, err["DebugField"])
			},
		},
		{
			description:  "Should pass Hybrid grant test with form_post",
			responseType: "token%20code",
			responseMode: "form_post",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
				responseModeClient.ResponseModes = []fosite.ResponseModeType{fosite.ResponseModeFormPost}
			},
			check: func(t *testing.T, stateFromServer string, code string, token goauth.Token, iDToken string, err map[string]string) {
				assert.EqualValues(t, state, stateFromServer)
				assert.NotEmpty(t, code)
				assert.NotEmpty(t, token.TokenType)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.Expiry)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			c.setup()
			authURL := strings.Replace(oauthClient.AuthCodeURL(state, goauth.SetAuthURLParam("response_mode", c.responseMode), goauth.SetAuthURLParam("nonce", "111111111")), "response_type=code", "response_type="+c.responseType, -1)

			var (
				callbackURL *url.URL
				redirErr    = errors.New("Dont follow redirects")
			)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					callbackURL = req.URL
					return redirErr
				},
			}

			var (
				code, state, iDToken string
				token                goauth.Token
				errResp              map[string]string
			)

			resp, err := client.Get(authURL)
			if fosite.ResponseModeType(c.responseMode) == fosite.ResponseModeFragment {
				// fragment
				require.EqualError(t, errors.Unwrap(err), redirErr.Error())
				fragment, err := url.ParseQuery(callbackURL.Fragment)
				require.NoError(t, err)
				code, state, iDToken, token, errResp = getParameters(t, fragment)
			} else if fosite.ResponseModeType(c.responseMode) == fosite.ResponseModeQuery {
				// query
				require.EqualError(t, errors.Unwrap(err), redirErr.Error())
				query, err := url.ParseQuery(callbackURL.RawQuery)
				require.NoError(t, err)
				code, state, iDToken, token, errResp = getParameters(t, query)
			} else if fosite.ResponseModeType(c.responseMode) == fosite.ResponseModeFormPost {
				// form_post
				require.NoError(t, err)
				code, state, iDToken, token, _, errResp, err = internal.ParseFormPostResponse(fositeStore.Clients["response-mode-client"].GetRedirectURIs()[0], resp.Body)
				require.NoError(t, err)
			} else {
				t.FailNow()
			}

			c.check(t, state, code, token, iDToken, errResp)
		})
	}
}

func getParameters(t *testing.T, param url.Values) (code, state, iDToken string, token goauth.Token, errResp map[string]string) {
	errResp = make(map[string]string)
	if param.Get("error") != "" {
		errResp["ErrorField"] = param.Get("error")
		errResp["DescriptionField"] = param.Get("error_description")
		errResp["HintField"] = param.Get("error_hint")
	} else {
		code = param.Get("code")
		state = param.Get("state")
		iDToken = param.Get("id_token")
		token = goauth.Token{
			AccessToken:  param.Get("access_token"),
			TokenType:    param.Get("token_type"),
			RefreshToken: param.Get("refresh_token"),
		}
		if param.Get("expires_in") != "" {
			expires, err := strconv.Atoi(param.Get("expires_in"))
			require.NoError(t, err)
			token.Expiry = time.Now().UTC().Add(time.Duration(expires) * time.Second)
		}
	}
	return
}
