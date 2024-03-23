// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/internal/testhelpers"
	hydraoauth2 "github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/contextx"
)

func TestDeviceAuthRequest(t *testing.T) {
	ctx := context.Background()
	reg := internal.NewMockedRegistry(t, &contextx.Default{})
	testhelpers.NewOAuth2Server(ctx, t, reg)

	c := &client.Client{
		ResponseTypes: []string{"id_token", "code", "token"},
		GrantTypes: []string{
			string(fosite.GrantTypeDeviceCode),
		},
		Scope:                   "hydra offline openid",
		Audience:                []string{"https://api.ory.sh/"},
		TokenEndpointAuthMethod: "none",
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))

	oauthClient := &oauth2.Config{
		ClientID: c.GetID(),
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
		},
		Scopes: strings.Split(c.Scope, " "),
	}

	testCases := []struct {
		description string
		setUp       func()
		check       func(t *testing.T, resp *oauth2.DeviceAuthResponse, err error)
		cleanUp     func()
	}{
		{
			description: "should pass",
			check: func(t *testing.T, resp *oauth2.DeviceAuthResponse, _ error) {
				assert.NotEmpty(t, resp.DeviceCode)
				assert.NotEmpty(t, resp.UserCode)
				assert.NotEmpty(t, resp.Interval)
				assert.NotEmpty(t, resp.VerificationURI)
				assert.NotEmpty(t, resp.VerificationURIComplete)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("case="+testCase.description, func(t *testing.T) {
			if testCase.setUp != nil {
				testCase.setUp()
			}

			resp, err := oauthClient.DeviceAuth(context.Background())

			if testCase.check != nil {
				testCase.check(t, resp, err)
			}

			if testCase.cleanUp != nil {
				testCase.cleanUp()
			}
		})
	}
}

func TestDeviceTokenRequest(t *testing.T) {
	ctx := context.Background()
	reg := internal.NewMockedRegistry(t, &contextx.Default{})
	testhelpers.NewOAuth2Server(ctx, t, reg)

	c := &client.Client{
		GrantTypes: []string{
			string(fosite.GrantTypeDeviceCode),
		},
		Scope:                   "hydra offline openid",
		Audience:                []string{"https://api.ory.sh/"},
		TokenEndpointAuthMethod: "none",
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))

	oauthClient := &oauth2.Config{
		ClientID: c.GetID(),
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
		},
		Scopes: strings.Split(c.Scope, " "),
	}

	var code, signature string
	var err error
	code, signature, err = reg.RFC8628HMACStrategy().GenerateDeviceCode(context.TODO())
	require.NoError(t, err)

	testCases := []struct {
		description string
		setUp       func()
		check       func(t *testing.T, token *oauth2.Token, err error)
		cleanUp     func()
	}{
		{
			description: "should pass",
			setUp: func() {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client: &fosite.DefaultClient{ID: c.GetID(), GrantTypes: []string{string(fosite.GrantTypeDeviceCode)}},
						Session: &hydraoauth2.Session{
							DefaultSession: &openid.DefaultSession{
								ExpiresAt: map[fosite.TokenType]time.Time{
									fosite.DeviceCode: time.Now().Add(time.Hour).UTC(),
								},
							},
							BrowserFlowCompleted: true,
						},
						RequestedAt: time.Now(),
					},
				}

				require.NoError(t, reg.OAuth2Storage().CreateDeviceCodeSession(context.TODO(), signature, authreq))
			},
			check: func(t *testing.T, token *oauth2.Token, err error) {
				assert.NotEmpty(t, token.AccessToken)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("case="+testCase.description, func(t *testing.T) {
			if testCase.setUp != nil {
				testCase.setUp()
			}

			var token *oauth2.Token
			token, err = oauthClient.DeviceAccessToken(context.Background(), &oauth2.DeviceAuthResponse{DeviceCode: code})

			if testCase.check != nil {
				testCase.check(t, token, err)
			}

			if testCase.cleanUp != nil {
				testCase.cleanUp()
			}
		})
	}
}
