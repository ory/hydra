// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestClientCredentials_HandleTokenEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockClientCredentialsGrantStorage(ctrl)
	chgen := internal.NewMockAccessTokenStrategy(ctrl)
	areq := internal.NewMockAccessRequester(ctrl)
	defer ctrl.Finish()

	h := ClientCredentialsGrantHandler{
		HandleHelper: &HandleHelper{
			AccessTokenStorage:  store,
			AccessTokenStrategy: chgen,
			Config: &fosite.Config{
				AccessTokenLifespan: time.Hour,
			},
		},
		Config: &fosite.Config{
			ScopeStrategy:            fosite.HierarchicScopeStrategy,
			AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
		},
	}
	for k, c := range []struct {
		description string
		mock        func()
		req         *http.Request
		expectErr   error
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			mock: func() {
				areq.EXPECT().GetGrantTypes().Return(fosite.Arguments{""})
			},
		},
		{
			description: "should fail because audience not valid",
			expectErr:   fosite.ErrInvalidRequest,
			mock: func() {
				areq.EXPECT().GetGrantTypes().Return(fosite.Arguments{"client_credentials"})
				areq.EXPECT().GetRequestedScopes().Return([]string{})
				areq.EXPECT().GetRequestedAudience().Return([]string{"https://www.ory.sh/not-api"})
				areq.EXPECT().GetClient().Return(&fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"client_credentials"},
					Audience:   []string{"https://www.ory.sh/api"},
				})
			},
		},
		{
			description: "should fail because scope not valid",
			expectErr:   fosite.ErrInvalidScope,
			mock: func() {
				areq.EXPECT().GetGrantTypes().Return(fosite.Arguments{"client_credentials"})
				areq.EXPECT().GetRequestedScopes().Return([]string{"foo", "bar", "baz.bar"})
				areq.EXPECT().GetClient().Return(&fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"client_credentials"},
					Scopes:     []string{"foo"},
				})
			},
		},
		{
			description: "should pass",
			mock: func() {
				areq.EXPECT().GetSession().Return(new(fosite.DefaultSession))
				areq.EXPECT().GetGrantTypes().Return(fosite.Arguments{"client_credentials"})
				areq.EXPECT().GetRequestedScopes().Return([]string{"foo", "bar", "baz.bar"})
				areq.EXPECT().GetRequestedAudience().Return([]string{})
				areq.EXPECT().GetClient().Return(&fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"client_credentials"},
					Scopes:     []string{"foo", "bar", "baz"},
				})
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.mock()
			err := h.HandleTokenEndpointRequest(context.Background(), areq)
			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClientCredentials_PopulateTokenEndpointResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockClientCredentialsGrantStorage(ctrl)
	chgen := internal.NewMockAccessTokenStrategy(ctrl)
	areq := fosite.NewAccessRequest(new(fosite.DefaultSession))
	aresp := fosite.NewAccessResponse()
	defer ctrl.Finish()

	h := ClientCredentialsGrantHandler{
		HandleHelper: &HandleHelper{
			AccessTokenStorage:  store,
			AccessTokenStrategy: chgen,
			Config: &fosite.Config{
				AccessTokenLifespan: time.Hour,
			},
		},
		Config: &fosite.Config{
			ScopeStrategy: fosite.HierarchicScopeStrategy,
		},
	}
	for k, c := range []struct {
		description string
		mock        func()
		req         *http.Request
		expectErr   error
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			mock: func() {
				areq.GrantTypes = fosite.Arguments{""}
			},
		},
		{
			description: "should fail because grant_type not allowed",
			expectErr:   fosite.ErrUnauthorizedClient,
			mock: func() {
				areq.GrantTypes = fosite.Arguments{"client_credentials"}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"authorization_code"}}
			},
		},
		{
			description: "should pass",
			mock: func() {
				areq.GrantTypes = fosite.Arguments{"client_credentials"}
				areq.Session = &fosite.DefaultSession{}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"client_credentials"}}
				chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return("tokenfoo.bar", "bar", nil)
				store.EXPECT().CreateAccessTokenSession(gomock.Any(), "bar", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.mock()
			err := h.PopulateTokenEndpointResponse(context.Background(), areq, aresp)
			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
