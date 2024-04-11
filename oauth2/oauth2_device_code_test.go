// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"

	"github.com/ory/fosite/token/jwt"

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

	secret := uuid.New()
	c := &client.Client{
		ID:     "device-client",
		Secret: secret,
		GrantTypes: []string{
			string(fosite.GrantTypeDeviceCode),
		},
		Scope:                   "hydra offline openid",
		Audience:                []string{"https://api.ory.sh/"},
		TokenEndpointAuthMethod: "client_secret_post",
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))

	oauthClient := &oauth2.Config{
		ClientID:     c.GetID(),
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:     oauth2.AuthStyleInParams,
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

			resp, err := oauthClient.DeviceAuth(context.Background(), []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("client_secret", secret)}...)

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

	secret := uuid.New()
	c := &client.Client{
		ID:     "device-client",
		Secret: secret,
		GrantTypes: []string{
			string(fosite.GrantTypeDeviceCode),
			string(fosite.GrantTypeRefreshToken),
		},
		Scope:    "hydra offline openid",
		Audience: []string{"https://api.ory.sh/"},
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, c))

	oauthClient := &oauth2.Config{
		ClientID:     c.GetID(),
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: reg.Config().OAuth2DeviceAuthorisationURL(ctx).String(),
			TokenURL:      reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle:     oauth2.AuthStyleInHeader,
		},
		Scopes: strings.Split(c.Scope, " "),
	}

	testCases := []struct {
		description string
		setUp       func(signature string)
		check       func(t *testing.T, token *oauth2.Token, err error)
		cleanUp     func()
	}{
		{
			description: "should pass with refresh token",
			setUp: func(signature string) {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client: &fosite.DefaultClient{
							ID:         c.GetID(),
							GrantTypes: []string{string(fosite.GrantTypeDeviceCode)},
						},
						RequestedScope: []string{"hydra", "offline"},
						GrantedScope:   []string{"hydra", "offline"},
						Session: &hydraoauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "hydra",
								},
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
				assert.NotEmpty(t, token.RefreshToken)
			},
		},
		{
			description: "should pass with ID token",
			setUp: func(signature string) {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client: &fosite.DefaultClient{
							ID:         c.GetID(),
							GrantTypes: []string{string(fosite.GrantTypeDeviceCode)},
						},
						RequestedScope: []string{"hydra", "offline", "openid"},
						GrantedScope:   []string{"hydra", "offline", "openid"},
						Session: &hydraoauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "hydra",
								},
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
				require.NoError(t, reg.OAuth2Storage().CreateOpenIDConnectSession(context.TODO(), signature, authreq))
			},
			check: func(t *testing.T, token *oauth2.Token, err error) {
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)
				assert.NotEmpty(t, token.Extra("id_token"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("case="+testCase.description, func(t *testing.T) {
			code, signature, err := reg.RFC8628HMACStrategy().GenerateDeviceCode(context.TODO())
			require.NoError(t, err)

			if testCase.setUp != nil {
				testCase.setUp(signature)
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
