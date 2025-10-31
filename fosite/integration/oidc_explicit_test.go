// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func newIDSession(j *jwt.IDTokenClaims) *defaultSession {
	return &defaultSession{
		DefaultSession: &openid.DefaultSession{
			Claims:  j,
			Headers: &jwt.Headers{},
			Subject: j.Subject,
		},
	}
}

func TestOpenIDConnectExplicitFlow(t *testing.T) {
	f := compose.ComposeAllEnabled(&fosite.Config{
		GlobalSecret: []byte("some-secret-thats-random-some-secret-thats-random-")}, fositeStore, gen.MustRSAKey())

	for k, c := range []struct {
		description    string
		setup          func(oauthClient *oauth2.Config) string
		authStatusCode int
		authCodeURL    string
		session        *defaultSession
		expectAuthErr  string
		expectTokenErr string
	}{
		{
			session:     newIDSession(&jwt.IDTokenClaims{Subject: "peter"}),
			description: "should pass",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=11234123"
			},
			authStatusCode: http.StatusOK,
		},
		{
			session:     newIDSession(&jwt.IDTokenClaims{Subject: "peter"}),
			description: "should fail registered single redirect uri but no redirect uri in request",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				oauthClient.RedirectURL = ""

				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=11234123"
			},
			authStatusCode: http.StatusBadRequest,
			expectAuthErr:  `{"error":"invalid_request","error_description":"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. The 'redirect_uri' parameter is required when using OpenID Connect 1.0."}`,
		},
		{
			session:     newIDSession(&jwt.IDTokenClaims{Subject: "peter"}),
			description: "should fail registered single redirect uri but no redirect uri in request",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				oauthClient.RedirectURL = ""

				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=11234123"
			},
			authStatusCode: http.StatusBadRequest,
			expectAuthErr:  `{"error":"invalid_request","error_description":"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. The 'redirect_uri' parameter is required when using OpenID Connect 1.0."}`,
		},
		{
			session:     newIDSession(&jwt.IDTokenClaims{Subject: "peter"}),
			description: "should fail because nonce is not long enough",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=1"
			},
			authStatusCode: http.StatusOK,
			expectTokenErr: "insufficient_entropy",
		},
		{
			session: newIDSession(&jwt.IDTokenClaims{
				Subject:     "peter",
				RequestedAt: time.Now().UTC(),
				AuthTime:    time.Now().Add(time.Second).UTC(),
			}),
			description: "should not pass missing redirect uri",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.RedirectURL = ""
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=1234567890&prompt=login"
			},
			expectAuthErr:  `{"error":"invalid_request","error_description":"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. The 'redirect_uri' parameter is required when using OpenID Connect 1.0."}`,
			authStatusCode: http.StatusBadRequest,
		},
		{
			session:     newIDSession(&jwt.IDTokenClaims{Subject: "peter"}),
			description: "should fail because state is not long enough",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("123") + "&nonce=1234567890"
			},
			expectAuthErr:  "invalid_state",
			authStatusCode: http.StatusNotAcceptable, // code from internal test callback handler when error occurs
		},
		{
			session: newIDSession(&jwt.IDTokenClaims{
				Subject:     "peter",
				RequestedAt: time.Now().UTC(),
				AuthTime:    time.Now().Add(time.Second).UTC(),
			}),
			description: "should pass",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=1234567890&prompt=login"
			},
			authStatusCode: http.StatusOK,
		},
		{
			session: newIDSession(&jwt.IDTokenClaims{
				Subject:     "peter",
				RequestedAt: time.Now().UTC(),
				AuthTime:    time.Now().Add(time.Second).UTC(),
			}),
			description: "should not pass missing redirect uri",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.RedirectURL = ""
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=1234567890&prompt=login"
			},
			expectAuthErr:  `{"error":"invalid_request","error_description":"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. The 'redirect_uri' parameter is required when using OpenID Connect 1.0."}`,
			authStatusCode: http.StatusBadRequest,
		},
		{
			session: newIDSession(&jwt.IDTokenClaims{
				Subject:     "peter",
				RequestedAt: time.Now().UTC(),
				AuthTime:    time.Now().Add(-time.Minute).UTC(),
			}),
			description: "should fail because authentication was in the past",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=1234567890&prompt=login"
			},
			authStatusCode: http.StatusNotAcceptable, // code from internal test callback handler when error occurs
			expectAuthErr:  "login_required",
		},
		{
			session: newIDSession(&jwt.IDTokenClaims{
				Subject:     "peter",
				RequestedAt: time.Now().UTC(),
				AuthTime:    time.Now().Add(-time.Minute).UTC(),
			}),
			description: "should pass because authorization was in the past and no login was required",
			setup: func(oauthClient *oauth2.Config) string {
				oauthClient.Scopes = []string{"openid"}
				return oauthClient.AuthCodeURL("12345678901234567890") + "&nonce=1234567890&prompt=none"
			},
			authStatusCode: http.StatusOK,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			ts := mockServer(t, f, c.session)
			defer ts.Close()

			oauthClient := newOAuth2Client(ts)

			fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs = []string{ts.URL + "/callback"}

			resp, err := http.Get(c.setup(oauthClient))
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			require.Equal(t, c.authStatusCode, resp.StatusCode, "Got response: %s", body)
			if resp.StatusCode >= 400 {
				assert.Equal(t, c.expectAuthErr, strings.Replace(string(body), "error: ", "", 1))
			}

			if c.expectAuthErr != "" {
				assert.Empty(t, resp.Request.URL.Query().Get("code"))
			}

			if resp.StatusCode == http.StatusOK {
				time.Sleep(time.Second)

				token, err := oauthClient.Exchange(context.Background(), resp.Request.URL.Query().Get("code"))
				if c.expectTokenErr != "" {
					require.Error(t, err)
					assert.True(t, strings.Contains(err.Error(), c.expectTokenErr), err.Error())
				} else {
					require.NoError(t, err)
					assert.NotEmpty(t, token.AccessToken)
					assert.NotEmpty(t, token.Extra("id_token"))
				}
			}
		})
	}
}
