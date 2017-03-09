package integration_test

import (
	"net/http"
	"testing"

	"fmt"

	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/handler/openid"
	"github.com/ory-am/fosite/internal"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestOpenIDConnectExplicitFlow(t *testing.T) {
	session := &defaultSession{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Subject: "peter",
			},
			Headers: &jwt.Headers{},
		},
	}
	f := compose.ComposeAllEnabled(new(compose.Config), fositeStore, []byte("some-secret-thats-random"), internal.MustRSAKey())
	ts := mockServer(t, f, session)

	defer ts.Close()
	oauthClient := newOAuth2Client(ts)
	fositeStore.Clients["my-client"].RedirectURIs[0] = ts.URL + "/callback"

	var state string
	for k, c := range []struct {
		description    string
		setup          func()
		authStatusCode int
	}{
		{
			description: "should pass",
			setup: func() {
				state = "12345678901234567890"
				oauthClient.Scopes = []string{"openid"}
			},
			authStatusCode: http.StatusOK,
		},
	} {
		c.setup()

		resp, err := http.Get(oauthClient.AuthCodeURL(state) + "&nonce=1234567890")
		require.Nil(t, err)
		require.Equal(t, c.authStatusCode, resp.StatusCode, "(%d) %s", k, c.description)

		if resp.StatusCode == http.StatusOK {
			token, err := oauthClient.Exchange(oauth2.NoContext, resp.Request.URL.Query().Get("code"))
			fmt.Printf("after exchange: %s\n\n", fositeStore.AuthorizeCodes)
			require.Nil(t, err, "(%d) %s", k, c.description)
			require.NotEmpty(t, token.AccessToken, "(%d) %s", k, c.description)
			require.NotEmpty(t, token.Extra("id_token"), "(%d) %s", k, c.description)
		}
		t.Logf("Passed test case (%d) %s", k, c.description)
	}
}
