// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestIntrospectToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	accessTokenStorageProvider := internal.NewMockAccessTokenStorageProvider(ctrl)
	accessTokenStorage := internal.NewMockAccessTokenStorage(ctrl)
	refreshTokenStorageProvider := internal.NewMockRefreshTokenStorageProvider(ctrl)
	refreshTokenStorage := internal.NewMockRefreshTokenStorage(ctrl)
	accessTokenStrategyProvider := internal.NewMockAccessTokenStrategyProvider(ctrl)
	accessTokenStrategy := internal.NewMockAccessTokenStrategy(ctrl)
	refreshTokenStrategyProvider := internal.NewMockRefreshTokenStrategyProvider(ctrl)
	refreshTokenStrategy := internal.NewMockRefreshTokenStrategy(ctrl)
	areq := fosite.NewAccessRequest(nil)
	t.Cleanup(ctrl.Finish)

	mockStorage := struct {
		*internal.MockAccessTokenStorageProvider
		*internal.MockRefreshTokenStorageProvider
	}{
		MockAccessTokenStorageProvider:  accessTokenStorageProvider,
		MockRefreshTokenStorageProvider: refreshTokenStorageProvider,
	}
	mockStrategy := struct {
		*internal.MockAccessTokenStrategyProvider
		*internal.MockRefreshTokenStrategyProvider
	}{
		MockAccessTokenStrategyProvider:  accessTokenStrategyProvider,
		MockRefreshTokenStrategyProvider: refreshTokenStrategyProvider,
	}

	config := &fosite.Config{}
	v := &oauth2.CoreValidator{
		Strategy: mockStrategy,
		Storage:  mockStorage,
		Config:   config,
	}
	httpreq := &http.Request{Header: http.Header{}}

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		expectTU    fosite.TokenUse
	}{
		{
			description: "should fail because no bearer token set",
			setup: func() {
				httpreq.Header.Set("Authorization", "bearer")
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), "").Return("")
				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), "", nil).Return(nil, errors.New(""))
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), "").Return("")
				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), "", nil).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrRequestUnauthorized,
		},
		{
			description: "should fail because retrieval fails",
			setup: func() {
				httpreq.Header.Set("Authorization", "bearer 1234")
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().AccessTokenSignature(gomock.Any(), "1234").AnyTimes().Return("asdf")
				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).Times(1)
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), "asdf", nil).Return(nil, errors.New(""))
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), "1234").Return("asdf")
				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), "asdf", nil).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrRequestUnauthorized,
		},
		{
			description: "should fail because validation fails",
			setup: func() {
				accessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(accessTokenStorage).AnyTimes()
				accessTokenStorage.EXPECT().GetAccessTokenSession(gomock.Any(), "asdf", nil).AnyTimes().Return(areq, nil)
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(2)
				accessTokenStrategy.EXPECT().ValidateAccessToken(gomock.Any(), areq, "1234").Return(errorsx.WithStack(fosite.ErrTokenExpired))
				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().RefreshTokenSignature(gomock.Any(), "1234").Return("asdf")
				refreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(refreshTokenStorage).Times(1)
				refreshTokenStorage.EXPECT().GetRefreshTokenSession(gomock.Any(), "asdf", nil).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrTokenExpired,
		},
		{
			description: "should fail because access token invalid",
			setup: func() {
				config.DisableRefreshTokenValidation = true
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(2)
				accessTokenStrategy.EXPECT().ValidateAccessToken(gomock.Any(), areq, "1234").Return(errorsx.WithStack(fosite.ErrInvalidTokenFormat))
			},
			expectErr: fosite.ErrInvalidTokenFormat,
		},
		{
			description: "should pass",
			setup: func() {
				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(2)
				accessTokenStrategy.EXPECT().ValidateAccessToken(gomock.Any(), areq, "1234").Return(nil)
			},
			expectTU: fosite.AccessToken,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			tu, err := v.IntrospectToken(context.Background(), fosite.AccessTokenFromRequest(httpreq), fosite.AccessToken, areq, []string{})

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.expectTU, tu)
			}
		})
	}
}
