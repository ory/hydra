// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	hst "github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestResourceOwnerPasswordCredentialsFlow(t *testing.T) {
	for _, strategy := range []hst.AccessTokenStrategy{
		hmacStrategy,
	} {
		runResourceOwnerPasswordCredentialsGrantTest(t, strategy)
	}
}

func runResourceOwnerPasswordCredentialsGrantTest(t *testing.T, strategy hst.AccessTokenStrategy) {
	f := compose.Compose(new(fosite.Config), fositeStore, strategy, compose.OAuth2ResourceOwnerPasswordCredentialsFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	var username, password string
	oauthClient := newOAuth2Client(ts)
	for k, c := range []struct {
		description string
		setup       func()
		check       func(t *testing.T, token *oauth2.Token)
		err         bool
	}{
		{
			description: "should fail because invalid password",
			setup: func() {
				username = "peter"
				password = "something-wrong"
			},
			err: true,
		},
		{
			description: "should pass",
			setup: func() {
				password = "secret"
			},
		},
		{
			description: "should pass with custom client token lifespans",
			setup: func() {
				oauthClient = newOAuth2Client(ts)
				oauthClient.ClientID = "custom-lifespan-client"
			},
			check: func(t *testing.T, token *oauth2.Token) {
				s, err := fositeStore.GetAccessTokenSession(context.Background(), strings.Split(token.AccessToken, ".")[1], nil)
				require.NoError(t, err)
				atExp := s.GetSession().GetExpiresAt(fosite.AccessToken)
				internal.RequireEqualTime(t, time.Now().UTC().Add(*internal.TestLifespans.PasswordGrantAccessTokenLifespan), atExp, time.Minute)
				atExpIn := time.Duration(token.Extra("expires_in").(float64)) * time.Second
				internal.RequireEqualDuration(t, *internal.TestLifespans.PasswordGrantAccessTokenLifespan, atExpIn, time.Minute)
				rtExp := s.GetSession().GetExpiresAt(fosite.RefreshToken)
				internal.RequireEqualTime(t, time.Now().UTC().Add(*internal.TestLifespans.PasswordGrantRefreshTokenLifespan), rtExp, time.Minute)
			},
		},
	} {
		c.setup()

		token, err := oauthClient.PasswordCredentialsToken(context.Background(), username, password)
		require.Equal(t, c.err, err != nil, "(%d) %s\n%s\n%s", k, c.description, c.err, err)
		if !c.err {
			assert.NotEmpty(t, token.AccessToken, "(%d) %s\n%s", k, c.description, token)

			if c.check != nil {
				c.check(t, token)
			}
		}

		t.Logf("Passed test case %d", k)
	}
}
