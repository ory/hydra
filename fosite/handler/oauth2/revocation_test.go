// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestRevokeToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockTokenRevocationStorage(ctrl)
	atStrat := internal.NewMockAccessTokenStrategy(ctrl)
	rtStrat := internal.NewMockRefreshTokenStrategy(ctrl)
	ar := internal.NewMockAccessRequester(ctrl)
	defer ctrl.Finish()

	h := TokenRevocationHandler{
		TokenRevocationStorage: store,
		RefreshTokenStrategy:   rtStrat,
		AccessTokenStrategy:    atStrat,
	}

	var token string
	var tokenType fosite.TokenType

	for k, c := range []struct {
		description string
		mock        func()
		expectErr   error
		client      fosite.Client
	}{
		{
			description: "should fail - token was issued to another client",
			expectErr:   fosite.ErrUnauthorizedClient,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "foo"})
			},
		},
		{
			description: "should pass - refresh token discovery first; refresh token found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)
				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})
				store.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				store.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - access token discovery first; access token found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.AccessToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)
				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})
				store.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				store.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - refresh token discovery first; refresh token not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.AccessToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)
				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})
				store.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				store.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - access token discovery first; access token not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)
				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})
				store.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				store.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - refresh token discovery first; both tokens not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should pass - access token discovery first; both tokens not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.AccessToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{

			description: "should pass - refresh token discovery first; refresh token is inactive",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrInactiveToken)

				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should pass - access token discovery first; refresh token is inactive",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.AccessToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrInactiveToken)
			},
		},
		{
			description: "should fail - store error for access token get",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.AccessToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("random error"))

				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should fail - store error for refresh token get",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("random error"))
			},
		},
		{
			description: "should fail - store error for access token revoke",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.AccessToken
				atStrat.EXPECT().AccessTokenSignature(gomock.Any(), token)
				store.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})
				store.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any()).Return(fosite.ErrNotFound)
				store.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any()).Return(fmt.Errorf("random error"))
			},
		},
		{
			description: "should fail - store error for refresh token revoke",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func() {
				token = "foo"
				tokenType = fosite.RefreshToken
				rtStrat.EXPECT().RefreshTokenSignature(gomock.Any(), token)
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})
				store.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any()).Return(fmt.Errorf("random error"))
				store.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any()).Return(fosite.ErrNotFound)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			c.mock()
			err := h.RevokeToken(context.Background(), token, tokenType, c.client)

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
