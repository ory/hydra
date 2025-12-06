// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestRevokeToken(t *testing.T) {
	for k, c := range []struct {
		description string
		mock        func(
			ar *internal.MockAccessRequester,
			tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
			tokenRevocationStorage *internal.MockTokenRevocationStorage,
			accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
			accessTokenStorage *internal.MockAccessTokenStorage,
			refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
			refreshTokenStorage *internal.MockRefreshTokenStorage,
			accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
			accessTokenStrategy *internal.MockAccessTokenStrategy,
			refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
			refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			token *string,
			tokenType *fosite.TokenType,
		)
		expectErr error
		client    fosite.Client
	}{
		{
			description: "should fail - token was issued to another client",
			expectErr:   fosite.ErrUnauthorizedClient,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "foo"})
			},
		},
		{
			description: "should pass - refresh token discovery first; refresh token found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})

				tokenRevocationStorageProvider.EXPECT().TokenRevocationStorage().Return(tokenRevocationStorage).Times(2)
				tokenRevocationStorage.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				tokenRevocationStorage.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - access token discovery first; access token found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.AccessToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})

				tokenRevocationStorageProvider.EXPECT().TokenRevocationStorage().Return(tokenRevocationStorage).Times(2)
				tokenRevocationStorage.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				tokenRevocationStorage.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - refresh token discovery first; refresh token not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.AccessToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})

				tokenRevocationStorageProvider.EXPECT().TokenRevocationStorage().Return(tokenRevocationStorage).Times(2)
				tokenRevocationStorage.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				tokenRevocationStorage.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - access token discovery first; access token not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})

				tokenRevocationStorageProvider.EXPECT().TokenRevocationStorage().Return(tokenRevocationStorage).Times(2)
				tokenRevocationStorage.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any())
				tokenRevocationStorage.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any())
			},
		},
		{
			description: "should pass - refresh token discovery first; both tokens not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should pass - access token discovery first; both tokens not found",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.AccessToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should pass - refresh token discovery first; refresh token is inactive",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrInactiveToken)

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should pass - access token discovery first; refresh token is inactive",
			expectErr:   nil,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.AccessToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrInactiveToken)
			},
		},
		{
			description: "should fail - store error for access token get",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.AccessToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("random error"))

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
		},
		{
			description: "should fail - store error for refresh token get",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fosite.ErrNotFound)

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("random error"))
			},
		},
		{
			description: "should fail - store error for access token revoke",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.AccessToken
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), *token)

				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})

				tokenRevocationStorageProvider.EXPECT().TokenRevocationStorage().Return(tokenRevocationStorage).Times(2)
				tokenRevocationStorage.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any()).Return(fosite.ErrNotFound)
				tokenRevocationStorage.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any()).Return(fmt.Errorf("random error"))
			},
		},
		{
			description: "should fail - store error for refresh token revoke",
			expectErr:   fosite.ErrTemporarilyUnavailable,
			client:      &fosite.DefaultClient{ID: "bar"},
			mock: func(
				ar *internal.MockAccessRequester,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
				token *string,
				tokenType *fosite.TokenType,
			) {
				*token = "foo"
				*tokenType = fosite.RefreshToken
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), *token)

				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(ar, nil)

				ar.EXPECT().GetID()
				ar.EXPECT().GetClient().Return(&fosite.DefaultClient{ID: "bar"})

				tokenRevocationStorageProvider.EXPECT().TokenRevocationStorage().Return(tokenRevocationStorage).Times(2)
				tokenRevocationStorage.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any()).Return(fmt.Errorf("random error"))
				tokenRevocationStorage.EXPECT().RevokeAccessToken(gomock.Any(), gomock.Any()).Return(fosite.ErrNotFound)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			// define mocks
			ar := internal.NewMockAccessRequester(ctrl)

			tokenRevocationStorageProvider := internal.NewMockTokenRevocationStorageProvider(ctrl)
			tokenRevocationStorage := internal.NewMockTokenRevocationStorage(ctrl)

			accessTokenStorageProvider := internal.NewMockAccessTokenStorageProvider(ctrl)
			accessTokenStorage := internal.NewMockAccessTokenStorage(ctrl)

			refreshTokenStorageProvider := internal.NewMockRefreshTokenStorageProvider(ctrl)
			refreshTokenStorage := internal.NewMockRefreshTokenStorage(ctrl)

			accessTokenStrategyProvider := internal.NewMockAccessTokenStrategyProvider(ctrl)
			accessTokenStrategy := internal.NewMockAccessTokenStrategy(ctrl)

			refreshTokenStrategyProvider := internal.NewMockRefreshTokenStrategyProvider(ctrl)
			refreshTokenStrategy := internal.NewMockRefreshTokenStrategy(ctrl)

			// define concrete types
			var token string
			var tokenType fosite.TokenType

			mockStorage := struct {
				*internal.MockTokenRevocationStorageProvider
				*internal.MockAccessTokenStorageProvider
				*internal.MockRefreshTokenStorageProvider
			}{
				MockTokenRevocationStorageProvider: tokenRevocationStorageProvider,
				MockAccessTokenStorageProvider:     accessTokenStorageProvider,
				MockRefreshTokenStorageProvider:    refreshTokenStorageProvider,
			}

			mockStrategy := struct {
				*internal.MockAccessTokenStrategyProvider
				*internal.MockRefreshTokenStrategyProvider
			}{
				MockAccessTokenStrategyProvider:  accessTokenStrategyProvider,
				MockRefreshTokenStrategyProvider: refreshTokenStrategyProvider,
			}

			h := oauth2.TokenRevocationHandler{
				Storage:  mockStorage,
				Strategy: mockStrategy,
			}

			// set up mock expectations
			c.mock(
				ar,
				tokenRevocationStorageProvider,
				tokenRevocationStorage,
				accessTokenStorageProvider,
				accessTokenStorage,
				refreshTokenStorageProvider,
				refreshTokenStorage,
				accessTokenStrategyProvider,
				accessTokenStrategy,
				refreshTokenStrategyProvider,
				refreshTokenStrategy,
				&token,
				&tokenType,
			)

			// invoke function under test
			err := h.RevokeToken(context.Background(), token, tokenType, c.client)
			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
