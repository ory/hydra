// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
)

func TestIntrospectToken(t *testing.T) {
	for _, c := range []struct {
		description string
		strategy    oauth2.AccessTokenStrategy
		factory     compose.Factory
	}{
		{
			description: "HMAC strategy with OAuth2TokenIntrospectionFactory",
			strategy:    hmacStrategy,
			factory:     compose.OAuth2TokenIntrospectionFactory,
		},
		{
			description: "JWT strategy with OAuth2TokenIntrospectionFactory",
			strategy:    jwtStrategy,
			factory:     compose.OAuth2TokenIntrospectionFactory,
		},
		{
			description: "JWT strategy with OAuth2StatelessJWTIntrospectionFactory",
			strategy:    jwtStrategy,
			factory:     compose.OAuth2StatelessJWTIntrospectionFactory,
		},
	} {
		t.Logf("testing %v", c.description)
		runIntrospectTokenTest(t, c.strategy, c.factory)
	}
}

func runIntrospectTokenTest(t *testing.T, strategy oauth2.AccessTokenStrategy, introspectionFactory compose.Factory) {
	f := compose.Compose(new(fosite.Config), fositeStore, strategy, compose.OAuth2ClientCredentialsGrantFactory, introspectionFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2AppClient(ts)
	a, err := oauthClient.Token(goauth.NoContext)
	require.NoError(t, err)
	b, err := oauthClient.Token(goauth.NoContext)
	require.NoError(t, err)

	for k, c := range []struct {
		prepare  func(*gorequest.SuperAgent) *gorequest.SuperAgent
		isActive bool
		scopes   string
	}{
		{
			prepare: func(s *gorequest.SuperAgent) *gorequest.SuperAgent {
				return s.SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret)
			},
			isActive: true,
			scopes:   "",
		},
		{
			prepare: func(s *gorequest.SuperAgent) *gorequest.SuperAgent {
				return s.Set("Authorization", "bearer "+a.AccessToken)
			},
			isActive: true,
			scopes:   "fosite",
		},
		{
			prepare: func(s *gorequest.SuperAgent) *gorequest.SuperAgent {
				return s.Set("Authorization", "bearer "+a.AccessToken)
			},
			isActive: true,
			scopes:   "",
		},
		{
			prepare: func(s *gorequest.SuperAgent) *gorequest.SuperAgent {
				return s.Set("Authorization", "bearer "+a.AccessToken)
			},
			isActive: false,
			scopes:   "foo",
		},
		{
			prepare: func(s *gorequest.SuperAgent) *gorequest.SuperAgent {
				return s.Set("Authorization", "bearer "+b.AccessToken)
			},
			isActive: false,
			scopes:   "",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			res := struct {
				Active    bool    `json:"active"`
				ClientId  string  `json:"client_id"`
				Scope     string  `json:"scope"`
				ExpiresAt float64 `json:"exp"`
				IssuedAt  float64 `json:"iat"`
			}{}
			s := gorequest.New()
			s = s.Post(ts.URL + "/introspect").
				Type("form").
				SendStruct(map[string]string{"token": b.AccessToken, "scope": c.scopes})
			_, bytes, errs := c.prepare(s).End()

			assert.Nil(t, json.Unmarshal([]byte(bytes), &res))
			t.Logf("Got answer: %s", bytes)

			assert.Len(t, errs, 0)
			assert.Equal(t, c.isActive, res.Active)
			if c.isActive {
				assert.Equal(t, "fosite", res.Scope)
				assert.True(t, res.ExpiresAt > 0)
				assert.True(t, res.IssuedAt > 0)
				assert.True(t, res.IssuedAt < res.ExpiresAt)
			}
		})
	}
}
