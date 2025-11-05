// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/storage"
)

type mockStrategyProvider struct {
	strategy *rfc8628.DefaultDeviceStrategy
}

func (p mockStrategyProvider) DeviceRateLimitStrategy() rfc8628.DeviceRateLimitStrategy {
	return p.strategy
}

func (p mockStrategyProvider) DeviceCodeStrategy() rfc8628.DeviceCodeStrategy {
	return p.strategy
}

func (p mockStrategyProvider) UserCodeStrategy() rfc8628.UserCodeStrategy {
	return p.strategy
}

func Test_HandleDeviceEndpointRequest(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := rfc8628.DeviceAuthHandler{
		Storage:  store,
		Strategy: mockStrategyProvider{strategy: &hmacshaStrategyDefault},
		Config: &fosite.Config{
			DeviceAndUserCodeLifespan:      time.Minute * 10,
			DeviceAuthTokenPollingInterval: time.Second * 5,
			DeviceVerificationURL:          "www.test.com",
			AccessTokenLifespan:            time.Hour,
			RefreshTokenLifespan:           time.Hour,
			ScopeStrategy:                  fosite.HierarchicScopeStrategy,
			AudienceMatchingStrategy:       fosite.DefaultAudienceMatchingStrategy,
			RefreshTokenScopes:             []string{"offline"},
		},
	}

	req := &fosite.DeviceRequest{
		Request: fosite.Request{
			Client: &fosite.DefaultClient{
				Audience: []string{"https://www.ory.sh/api"},
			},
			Session: &fosite.DefaultSession{},
		},
	}
	resp := fosite.NewDeviceResponse()
	err := handler.HandleDeviceEndpointRequest(context.Background(), req, resp)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetDeviceCode())
	assert.NotEmpty(t, resp.GetUserCode())
	assert.Equal(t, len(resp.GetUserCode()), 8)
	assert.Contains(t, resp.GetDeviceCode(), "ory_dc_")
	assert.Contains(t, resp.GetDeviceCode(), ".")
	assert.Equal(t, resp.GetVerificationURI(), "www.test.com")
}

func Test_HandleDeviceEndpointRequestWithRetry(t *testing.T) {
	var mockDeviceAuthStorage *internal.MockDeviceAuthStorage
	var mockDeviceAuthStorageProvider *internal.MockDeviceAuthStorageProvider
	var mockAccessTokenStorageProvider *internal.MockAccessTokenStorageProvider
	var mockRefreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider
	var mockDeviceRateLimitStrategyProvider *internal.MockDeviceRateLimitStrategyProvider
	var mockDeviceCodeStrategy *internal.MockDeviceCodeStrategy
	var mockDeviceCodeStrategyProvider *internal.MockDeviceCodeStrategyProvider
	var mockUserCodeStrategy *internal.MockUserCodeStrategy
	var mockUserCodeStrategyProvider *internal.MockUserCodeStrategyProvider

	ctx := context.Background()
	req := &fosite.DeviceRequest{
		Request: fosite.Request{
			Client: &fosite.DefaultClient{
				Audience: []string{"https://www.ory.sh/api"},
			},
			Session: &fosite.DefaultSession{},
		},
	}

	testCases := []struct {
		description string
		setup       func()
		check       func(t *testing.T, resp *fosite.DeviceResponse)
		expectError error
	}{
		{
			description: "should pass when generating a unique user code at the first attempt",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy)
				mockDeviceCodeStrategy.
					EXPECT().
					GenerateDeviceCode(ctx).
					Return("deviceCode", "signature", nil)
				mockUserCodeStrategyProvider.EXPECT().UserCodeStrategy().Return(mockUserCodeStrategy)
				mockUserCodeStrategy.
					EXPECT().
					GenerateUserCode(ctx).
					Return("userCode", "signature2", nil).
					Times(1)
				mockDeviceAuthStorageProvider.
					EXPECT().
					DeviceAuthStorage().
					Return(mockDeviceAuthStorage).
					Times(1)
				mockDeviceAuthStorage.
					EXPECT().
					CreateDeviceAuthSession(ctx, "signature", "signature2", gomock.Any()).
					Return(nil)
			},
			check: func(t *testing.T, resp *fosite.DeviceResponse) {
				assert.Equal(t, "userCode", resp.GetUserCode())
			},
		},
		{
			description: "should pass when generating a unique user code within allowed attempts",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy)
				mockDeviceCodeStrategy.
					EXPECT().
					GenerateDeviceCode(ctx).
					Return("deviceCode", "signature", nil)
				gomock.InOrder(
					mockUserCodeStrategyProvider.EXPECT().UserCodeStrategy().Return(mockUserCodeStrategy),
					mockUserCodeStrategy.
						EXPECT().
						GenerateUserCode(ctx).
						Return("duplicatedUserCode", "duplicatedSignature", nil),
					mockDeviceAuthStorageProvider.
						EXPECT().
						DeviceAuthStorage().
						Return(mockDeviceAuthStorage).
						Times(1),
					mockDeviceAuthStorage.
						EXPECT().
						CreateDeviceAuthSession(ctx, "signature", "duplicatedSignature", gomock.Any()).
						Return(fosite.ErrExistingUserCodeSignature),
					mockUserCodeStrategyProvider.EXPECT().UserCodeStrategy().Return(mockUserCodeStrategy),
					mockUserCodeStrategy.
						EXPECT().
						GenerateUserCode(ctx).
						Return("uniqueUserCode", "uniqueSignature", nil),
					mockDeviceAuthStorageProvider.
						EXPECT().
						DeviceAuthStorage().
						Return(mockDeviceAuthStorage).
						Times(1),
					mockDeviceAuthStorage.
						EXPECT().
						CreateDeviceAuthSession(ctx, "signature", "uniqueSignature", gomock.Any()).
						Return(nil),
				)
			},
			check: func(t *testing.T, resp *fosite.DeviceResponse) {
				assert.Equal(t, "uniqueUserCode", resp.GetUserCode())
			},
		},
		{
			description: "should fail after maximum retries to generate a unique user code",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy)
				mockDeviceCodeStrategy.
					EXPECT().
					GenerateDeviceCode(ctx).
					Return("deviceCode", "signature", nil)
				mockUserCodeStrategyProvider.EXPECT().UserCodeStrategy().Return(mockUserCodeStrategy).Times(rfc8628.MaxAttempts)
				mockUserCodeStrategy.
					EXPECT().
					GenerateUserCode(ctx).
					Return("duplicatedUserCode", "duplicatedSignature", nil).
					Times(rfc8628.MaxAttempts)
				mockDeviceAuthStorageProvider.
					EXPECT().
					DeviceAuthStorage().
					Return(mockDeviceAuthStorage).
					Times(rfc8628.MaxAttempts)
				mockDeviceAuthStorage.
					EXPECT().
					CreateDeviceAuthSession(ctx, "signature", "duplicatedSignature", gomock.Any()).
					Return(fosite.ErrExistingUserCodeSignature).
					Times(rfc8628.MaxAttempts)
			},
			check: func(t *testing.T, resp *fosite.DeviceResponse) {
				assert.Empty(t, resp.GetUserCode())
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should fail if another error is returned",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy)
				mockDeviceCodeStrategy.
					EXPECT().
					GenerateDeviceCode(ctx).
					Return("deviceCode", "signature", nil)
				mockUserCodeStrategyProvider.EXPECT().UserCodeStrategy().Return(mockUserCodeStrategy)
				mockUserCodeStrategy.
					EXPECT().
					GenerateUserCode(ctx).
					Return("userCode", "userCodeSignature", nil)
				mockDeviceAuthStorageProvider.
					EXPECT().
					DeviceAuthStorage().
					Return(mockDeviceAuthStorage).
					Times(1)
				mockDeviceAuthStorage.
					EXPECT().
					CreateDeviceAuthSession(ctx, "signature", "userCodeSignature", gomock.Any()).
					Return(errors.New("some error"))
			},
			check: func(t *testing.T, resp *fosite.DeviceResponse) {
				assert.Empty(t, resp.GetUserCode())
			},
			expectError: fosite.ErrServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("scenario=%s", testCase.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockDeviceAuthStorage = internal.NewMockDeviceAuthStorage(ctrl)
			mockDeviceAuthStorageProvider = internal.NewMockDeviceAuthStorageProvider(ctrl)
			mockAccessTokenStorageProvider = internal.NewMockAccessTokenStorageProvider(ctrl)
			mockRefreshTokenStorageProvider = internal.NewMockRefreshTokenStorageProvider(ctrl)
			mockDeviceRateLimitStrategyProvider = internal.NewMockDeviceRateLimitStrategyProvider(ctrl)
			mockDeviceCodeStrategy = internal.NewMockDeviceCodeStrategy(ctrl)
			mockDeviceCodeStrategyProvider = internal.NewMockDeviceCodeStrategyProvider(ctrl)
			mockUserCodeStrategy = internal.NewMockUserCodeStrategy(ctrl)
			mockUserCodeStrategyProvider = internal.NewMockUserCodeStrategyProvider(ctrl)

			mockStorage := struct {
				*internal.MockDeviceAuthStorageProvider
				*internal.MockAccessTokenStorageProvider
				*internal.MockRefreshTokenStorageProvider
			}{
				MockDeviceAuthStorageProvider:   mockDeviceAuthStorageProvider,
				MockAccessTokenStorageProvider:  mockAccessTokenStorageProvider,
				MockRefreshTokenStorageProvider: mockRefreshTokenStorageProvider,
			}

			mockStrategy := struct {
				*internal.MockDeviceRateLimitStrategyProvider
				*internal.MockDeviceCodeStrategyProvider
				*internal.MockUserCodeStrategyProvider
			}{
				MockDeviceRateLimitStrategyProvider: mockDeviceRateLimitStrategyProvider,
				MockDeviceCodeStrategyProvider:      mockDeviceCodeStrategyProvider,
				MockUserCodeStrategyProvider:        mockUserCodeStrategyProvider,
			}

			h := rfc8628.DeviceAuthHandler{
				Storage:  mockStorage,
				Strategy: mockStrategy,
				Config: &fosite.Config{
					DeviceAndUserCodeLifespan:      time.Minute * 10,
					DeviceAuthTokenPollingInterval: time.Second * 5,
					DeviceVerificationURL:          "www.test.com",
					AccessTokenLifespan:            time.Hour,
					RefreshTokenLifespan:           time.Hour,
					ScopeStrategy:                  fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy:       fosite.DefaultAudienceMatchingStrategy,
					RefreshTokenScopes:             []string{"offline"},
				},
			}

			if testCase.setup != nil {
				testCase.setup()
			}

			resp := fosite.NewDeviceResponse()
			err := h.HandleDeviceEndpointRequest(ctx, req, resp)

			if testCase.expectError != nil {
				require.EqualError(t, err, testCase.expectError.Error(), "%+v", err)
			} else {
				require.NoError(t, err, "%+v", err)
			}

			if testCase.check != nil {
				testCase.check(t, resp)
			}
		})
	}
}
