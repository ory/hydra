// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestOIDCImplicitFlow(t *testing.T) {
	session := &defaultSession{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Subject: "peter",
			},
			Headers: &jwt.Headers{},
		},
	}
	f := compose.ComposeAllEnabled(&fosite.Config{
		GlobalSecret: []byte("some-secret-thats-random-some-secret-thats-random-"),
	}, fositeStore, gen.MustRSAKey())
	ts := mockServer(t, f, session)
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"

	var state = "12345678901234567890"
	for k, c := range []struct {
		responseType string
		description  string
		nonce        string
		setup        func()
		hasToken     bool
		hasIdToken   bool
		hasCode      bool
	}{
		{
			description:  "should pass without id token",
			responseType: "token",
			setup: func() {
				oauthClient.Scopes = []string{"fosite"}
			},
			hasToken: true,
		},
		{

			responseType: "id_token%20token",
			nonce:        "1111111111111111",
			description:  "should pass id token (id_token token)",
			setup: func() {
				oauthClient.Scopes = []string{"fosite", "openid"}
			},
			hasToken:   true,
			hasIdToken: true,
		},
		{

			responseType: "token%20id_token%20code",
			nonce:        "1111111111111111",
			description:  "should pass id token (code id_token token)",
			setup:        func() {},
			hasToken:     true,
			hasCode:      true,
			hasIdToken:   true,
		},
		{

			responseType: "token%20code",
			nonce:        "1111111111111111",
			description:  "should pass id token (code token)",
			setup:        func() {},
			hasToken:     true,
			hasCode:      true,
		},
		{

			responseType: "id_token%20code",
			nonce:        "1111111111111111",
			description:  "should pass id token (id_token code)",
			setup:        func() {},
			hasCode:      true,
			hasIdToken:   true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			c.setup()

			var callbackURL *url.URL
			authURL := strings.Replace(oauthClient.AuthCodeURL(state), "response_type=code", "response_type="+c.responseType, -1) + "&nonce=" + c.nonce
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					callbackURL = req.URL
					return errors.New("Dont follow redirects")
				},
			}
			resp, err := client.Get(authURL)
			require.Error(t, err)

			t.Logf("Response (%d): %s", k, callbackURL.String())
			fragment, err := url.ParseQuery(callbackURL.Fragment)
			require.NoError(t, err)

			if c.hasToken {
				assert.NotEmpty(t, fragment.Get("access_token"))
			} else {
				assert.Empty(t, fragment.Get("access_token"))
			}

			if c.hasCode {
				assert.NotEmpty(t, fragment.Get("code"))
			} else {
				assert.Empty(t, fragment.Get("code"))
			}

			if c.hasIdToken {
				assert.NotEmpty(t, fragment.Get("id_token"))
			} else {
				assert.Empty(t, fragment.Get("id_token"))
			}

			if !c.hasToken {
				return
			}

			expires, err := strconv.Atoi(fragment.Get("expires_in"))
			require.NoError(t, err)

			token := &oauth2.Token{
				AccessToken:  fragment.Get("access_token"),
				TokenType:    fragment.Get("token_type"),
				RefreshToken: fragment.Get("refresh_token"),
				Expiry:       time.Now().UTC().Add(time.Duration(expires) * time.Second),
			}

			httpClient := oauthClient.Client(context.Background(), token)
			resp, err = httpClient.Get(ts.URL + "/info")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			t.Logf("Passed test case (%d) %s", k, c.description)
		})
	}
}
