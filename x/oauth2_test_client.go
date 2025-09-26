// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/x/ioutilx"
	"github.com/ory/x/uuidx"
)

func NewEmptyCookieJar(t testing.TB) *cookiejar.Jar {
	c, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)
	return c
}

func NewEmptyJarClient(t testing.TB) *http.Client {
	return &http.Client{
		Jar: NewEmptyCookieJar(t),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//t.Logf("Redirect to %s", req.URL.String())

			if len(via) >= 20 {
				for k, v := range via {
					t.Logf("Failed with redirect (%d): %s", k, v.URL.String())
				}
				return errors.New("stopped after 20 redirects")
			}
			return nil
		},
	}
}

func GetExpectRedirect(t *testing.T, cl *http.Client, uri string) *url.URL {
	resp, err := cl.Get(uri)
	require.NoError(t, err)
	require.Equalf(t, 3, resp.StatusCode/100, "status: %d\nresponse: %s", resp.StatusCode, ioutilx.MustReadAll(resp.Body))
	loc, err := resp.Location()
	require.NoError(t, err)
	return loc
}

const (
	ClientCallbackURL = "https://client.ory/callback"
	LoginURL          = "https://ui.ory/login"
	ConsentURL        = "https://ui.ory/consent"
)

func PerformAuthCodeFlow(ctx context.Context, t *testing.T, baseClient *http.Client, cfg *oauth2.Config, admin *hydra.APIClient, lr func(*testing.T, *hydra.OAuth2LoginRequest) hydra.AcceptOAuth2LoginRequest, cr func(*testing.T, *hydra.OAuth2ConsentRequest) hydra.AcceptOAuth2ConsentRequest, authCodeOpts ...oauth2.AuthCodeOption) *oauth2.Token {
	var cl http.Client
	if baseClient != nil {
		cl = *baseClient
	}
	if cl.Jar == nil {
		cl.Jar = NewEmptyCookieJar(t)
	}
	cl.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

	// start the auth code flow
	state := uuidx.NewV4().String()
	loc := GetExpectRedirect(t, &cl, cfg.AuthCodeURL(state, authCodeOpts...))
	require.Equal(t, LoginURL, fmt.Sprintf("%s://%s%s", loc.Scheme, loc.Host, loc.Path))

	// get & submit the login request
	lReq, _, err := admin.OAuth2API.GetOAuth2LoginRequest(ctx).LoginChallenge(loc.Query().Get("login_challenge")).Execute()
	require.NoError(t, err)

	v, _, err := admin.OAuth2API.AcceptOAuth2LoginRequest(ctx).
		LoginChallenge(lReq.Challenge).
		AcceptOAuth2LoginRequest(lr(t, lReq)).
		Execute()
	require.NoError(t, err)

	loc = GetExpectRedirect(t, &cl, v.RedirectTo)
	require.Equal(t, ConsentURL, fmt.Sprintf("%s://%s%s", loc.Scheme, loc.Host, loc.Path))

	// get & submit the consent request
	cReq, _, err := admin.OAuth2API.GetOAuth2ConsentRequest(ctx).ConsentChallenge(loc.Query().Get("consent_challenge")).Execute()
	require.NoError(t, err)

	v, _, err = admin.OAuth2API.AcceptOAuth2ConsentRequest(ctx).
		ConsentChallenge(cReq.Challenge).
		AcceptOAuth2ConsentRequest(cr(t, cReq)).
		Execute()
	require.NoError(t, err)
	loc = GetExpectRedirect(t, &cl, v.RedirectTo)
	// ensure we got redirected to the client callback URL
	require.Equal(t, ClientCallbackURL, fmt.Sprintf("%s://%s%s", loc.Scheme, loc.Host, loc.Path))
	require.Equal(t, state, loc.Query().Get("state"))

	// exchange the code for a token
	code := loc.Query().Get("code")
	require.NotEmpty(t, code)
	token, err := cfg.Exchange(ctx, code)
	require.NoError(t, err)

	return token
}
