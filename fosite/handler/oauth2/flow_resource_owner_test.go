// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestResourceOwnerFlow_HandleTokenEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockResourceOwnerPasswordCredentialsGrantStorage(ctrl)
	defer ctrl.Finish()

	areq := fosite.NewAccessRequest(new(fosite.DefaultSession))
	areq.Form = url.Values{}
	for k, c := range []struct {
		description string
		setup       func(config *fosite.Config)
		expectErr   error
		check       func(areq *fosite.AccessRequest)
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			setup: func(config *fosite.Config) {
				areq.GrantTypes = fosite.Arguments{"123"}
			},
		},
		{
			description: "should fail because scope missing",
			setup: func(config *fosite.Config) {
				areq.GrantTypes = fosite.Arguments{"password"}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"password"}, Scopes: []string{}}
				areq.RequestedScope = []string{"foo-scope"}
			},
			expectErr: fosite.ErrInvalidScope,
		},
		{
			description: "should fail because audience missing",
			setup: func(config *fosite.Config) {
				areq.RequestedAudience = fosite.Arguments{"https://www.ory.sh/api"}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"password"}, Scopes: []string{"foo-scope"}}
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because invalid grant_type specified",
			setup: func(config *fosite.Config) {
				areq.GrantTypes = fosite.Arguments{"password"}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"authorization_code"}, Scopes: []string{"foo-scope"}}
			},
			expectErr: fosite.ErrUnauthorizedClient,
		},
		{
			description: "should fail because invalid credentials",
			setup: func(config *fosite.Config) {
				areq.Form.Set("username", "peter")
				areq.Form.Set("password", "pan")
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"password"}, Scopes: []string{"foo-scope"}, Audience: []string{"https://www.ory.sh/api"}}

				store.EXPECT().Authenticate(gomock.Any(), "peter", "pan").Return("", fosite.ErrNotFound)
			},
			expectErr: fosite.ErrInvalidGrant,
		},
		{
			description: "should fail because error on lookup",
			setup: func(config *fosite.Config) {
				store.EXPECT().Authenticate(gomock.Any(), "peter", "pan").Return("", errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			setup: func(config *fosite.Config) {
				store.EXPECT().Authenticate(gomock.Any(), "peter", "pan").Return("", nil)
			},
			check: func(areq *fosite.AccessRequest) {
				//assert.NotEmpty(t, areq.GetSession().GetExpiresAt(fosite.AccessToken))
				assert.Equal(t, time.Now().Add(time.Hour).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.AccessToken))
				assert.Equal(t, time.Now().Add(time.Hour).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.RefreshToken))
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			config := &fosite.Config{
				AccessTokenLifespan:      time.Hour,
				RefreshTokenLifespan:     time.Hour,
				ScopeStrategy:            fosite.HierarchicScopeStrategy,
				AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
			}
			h := ResourceOwnerPasswordCredentialsGrantHandler{
				ResourceOwnerPasswordCredentialsGrantStorage: store,
				HandleHelper: &HandleHelper{
					AccessTokenStorage: store,
					Config:             config,
				},
				Config: config,
			}
			c.setup(config)
			err := h.HandleTokenEndpointRequest(context.Background(), areq)

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				if c.check != nil {
					c.check(areq)
				}
			}
		})
	}
}

func TestResourceOwnerFlow_PopulateTokenEndpointResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockResourceOwnerPasswordCredentialsGrantStorage(ctrl)
	chgen := internal.NewMockAccessTokenStrategy(ctrl)
	rtstr := internal.NewMockRefreshTokenStrategy(ctrl)
	mockAT := "accesstoken.foo.bar"
	mockRT := "refreshtoken.bar.foo"
	defer ctrl.Finish()

	var areq *fosite.AccessRequest
	var aresp *fosite.AccessResponse
	config := &fosite.Config{}
	var h ResourceOwnerPasswordCredentialsGrantHandler
	h.Config = config

	for k, c := range []struct {
		description string
		setup       func(*fosite.Config)
		expectErr   error
		expect      func()
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			setup: func(config *fosite.Config) {
				areq.GrantTypes = fosite.Arguments{""}
			},
		},
		{
			description: "should pass",
			setup: func(config *fosite.Config) {
				areq.GrantTypes = fosite.Arguments{"password"}
				chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return(mockAT, "bar", nil)
				store.EXPECT().CreateAccessTokenSession(gomock.Any(), "bar", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
			},
			expect: func() {
				assert.Nil(t, aresp.GetExtra("refresh_token"), "unexpected refresh token")
			},
		},
		{
			description: "should pass - offline scope",
			setup: func(config *fosite.Config) {
				areq.GrantTypes = fosite.Arguments{"password"}
				areq.GrantScope("offline")
				rtstr.EXPECT().GenerateRefreshToken(gomock.Any(), areq).Return(mockRT, "bar", nil)
				store.EXPECT().CreateRefreshTokenSession(gomock.Any(), "bar", "bar", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
				chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return(mockAT, "bar", nil)
				store.EXPECT().CreateAccessTokenSession(gomock.Any(), "bar", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
			},
			expect: func() {
				assert.NotNil(t, aresp.GetExtra("refresh_token"), "expected refresh token")
			},
		},
		{
			description: "should pass - refresh token without offline scope",
			setup: func(config *fosite.Config) {
				config.RefreshTokenScopes = []string{}
				areq.GrantTypes = fosite.Arguments{"password"}
				chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return(mockAT, "bar", nil)
				store.EXPECT().CreateAccessTokenSession(gomock.Any(), "bar", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
				rtstr.EXPECT().GenerateRefreshToken(gomock.Any(), areq).Return(mockRT, "bar", nil)
				store.EXPECT().CreateRefreshTokenSession(gomock.Any(), "bar", "bar", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
			},
			expect: func() {
				assert.NotNil(t, aresp.GetExtra("refresh_token"), "expected refresh token")
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			areq = fosite.NewAccessRequest(nil)
			aresp = fosite.NewAccessResponse()
			areq.Session = &fosite.DefaultSession{}
			config := &fosite.Config{
				RefreshTokenScopes:  []string{"offline"},
				AccessTokenLifespan: time.Hour,
			}
			h = ResourceOwnerPasswordCredentialsGrantHandler{
				ResourceOwnerPasswordCredentialsGrantStorage: store,
				HandleHelper: &HandleHelper{
					AccessTokenStorage:  store,
					AccessTokenStrategy: chgen, Config: config,
				},
				RefreshTokenStrategy: rtstr, Config: config,
			}
			c.setup(config)
			err := h.PopulateTokenEndpointResponse(context.Background(), areq, aresp)
			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				if c.expect != nil {
					c.expect()
				}
			}
		})
	}
}
