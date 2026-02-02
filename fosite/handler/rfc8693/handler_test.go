// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8693_test

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/rfc8693"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestHandler_CanHandleTokenEndpointRequest(t *testing.T) {
	h := &rfc8693.Handler{}
	ctx := context.Background()

	req := fosite.NewAccessRequest(&fosite.DefaultSession{})
	req.GrantTypes = fosite.Arguments{string(fosite.GrantTypeTokenExchange)}
	assert.True(t, h.CanHandleTokenEndpointRequest(ctx, req))

	req.GrantTypes = fosite.Arguments{"authorization_code"}
	assert.False(t, h.CanHandleTokenEndpointRequest(ctx, req))
}

func TestHandler_HandleTokenEndpointRequest_MissingSubjectToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := internal.NewMockAccessTokenStorageProvider(ctrl)
	strategy := internal.NewMockAccessTokenStrategyProvider(ctrl)
	cfg := &fosite.Config{
		ScopeStrategy:                        fosite.HierarchicScopeStrategy,
		AudienceMatchingStrategy:             fosite.DefaultAudienceMatchingStrategy,
		AccessTokenLifespan:                  time.Hour,
		GrantTypeTokenExchangeCanSkipClientAuth: false,
	}
	handler := &rfc8693.Handler{Storage: storage, Strategy: strategy, Config: cfg}

	req := fosite.NewAccessRequest(&fosite.DefaultSession{})
	req.GrantTypes = fosite.Arguments{string(fosite.GrantTypeTokenExchange)}
	req.Client = &fosite.DefaultClient{GrantTypes: []string{string(fosite.GrantTypeTokenExchange)}}
	req.Form = url.Values{}
	req.Form.Set("subject_token_type", rfc8693.TokenTypeAccessToken)
	// subject_token missing

	err := handler.HandleTokenEndpointRequest(context.Background(), req)
	require.True(t, errors.Is(err, fosite.ErrInvalidRequest))
	assert.Contains(t, fosite.ErrorToRFC6749Error(err).HintField, "subject_token")
}

func TestHandler_HandleTokenEndpointRequest_MissingSubjectTokenType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := internal.NewMockAccessTokenStorageProvider(ctrl)
	strategy := internal.NewMockAccessTokenStrategyProvider(ctrl)
	cfg := &fosite.Config{
		ScopeStrategy:                        fosite.HierarchicScopeStrategy,
		AudienceMatchingStrategy:             fosite.DefaultAudienceMatchingStrategy,
		AccessTokenLifespan:                  time.Hour,
		GrantTypeTokenExchangeCanSkipClientAuth: false,
	}
	handler := &rfc8693.Handler{Storage: storage, Strategy: strategy, Config: cfg}

	req := fosite.NewAccessRequest(&fosite.DefaultSession{})
	req.GrantTypes = fosite.Arguments{string(fosite.GrantTypeTokenExchange)}
	req.Client = &fosite.DefaultClient{GrantTypes: []string{string(fosite.GrantTypeTokenExchange)}}
	req.Form = url.Values{}
	req.Form.Set("subject_token", "some-token")
	// subject_token_type missing

	err := handler.HandleTokenEndpointRequest(context.Background(), req)
	require.True(t, errors.Is(err, fosite.ErrInvalidRequest))
	assert.Contains(t, fosite.ErrorToRFC6749Error(err).HintField, "subject_token_type")
}

func TestHandler_HandleTokenEndpointRequest_UnsupportedSubjectTokenType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := internal.NewMockAccessTokenStorageProvider(ctrl)
	strategy := internal.NewMockAccessTokenStrategyProvider(ctrl)
	cfg := &fosite.Config{
		ScopeStrategy:                        fosite.HierarchicScopeStrategy,
		AudienceMatchingStrategy:             fosite.DefaultAudienceMatchingStrategy,
		AccessTokenLifespan:                  time.Hour,
		GrantTypeTokenExchangeCanSkipClientAuth: false,
	}
	handler := &rfc8693.Handler{Storage: storage, Strategy: strategy, Config: cfg}

	req := fosite.NewAccessRequest(&fosite.DefaultSession{})
	req.GrantTypes = fosite.Arguments{string(fosite.GrantTypeTokenExchange)}
	req.Client = &fosite.DefaultClient{GrantTypes: []string{string(fosite.GrantTypeTokenExchange)}}
	req.Form = url.Values{}
	req.Form.Set("subject_token", "x")
	req.Form.Set("subject_token_type", "urn:ietf:params:oauth:token-type:saml2")

	err := handler.HandleTokenEndpointRequest(context.Background(), req)
	require.True(t, errors.Is(err, fosite.ErrInvalidRequest))
	assert.Contains(t, fosite.ErrorToRFC6749Error(err).HintField, "Unsupported subject_token_type")
}

func TestHandler_CheckRequest_ClientNotAllowedGrantType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := internal.NewMockAccessTokenStorageProvider(ctrl)
	strategy := internal.NewMockAccessTokenStrategyProvider(ctrl)
	cfg := &fosite.Config{
		ScopeStrategy:                        fosite.HierarchicScopeStrategy,
		AudienceMatchingStrategy:             fosite.DefaultAudienceMatchingStrategy,
		AccessTokenLifespan:                  time.Hour,
		GrantTypeTokenExchangeCanSkipClientAuth: false,
	}
	handler := &rfc8693.Handler{Storage: storage, Strategy: strategy, Config: cfg}

	req := fosite.NewAccessRequest(&fosite.DefaultSession{})
	req.GrantTypes = fosite.Arguments{string(fosite.GrantTypeTokenExchange)}
	req.Client = &fosite.DefaultClient{GrantTypes: []string{"authorization_code"}} // no token-exchange

	err := handler.CheckRequest(context.Background(), req)
	require.True(t, errors.Is(err, fosite.ErrUnauthorizedClient))
}

func TestHandler_PopulateTokenEndpointResponse_SetsIssuedTokenType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockATStorage := internal.NewMockAccessTokenStorage(ctrl)
	mockATStorageProvider := internal.NewMockAccessTokenStorageProvider(ctrl)
	mockATStorageProvider.EXPECT().AccessTokenStorage().Return(mockATStorage).MinTimes(1)

	mockStrategy := internal.NewMockAccessTokenStrategy(ctrl)
	mockStrategyProvider := internal.NewMockAccessTokenStrategyProvider(ctrl)
	mockStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockStrategy).MinTimes(1)

	cfg := &fosite.Config{
		ScopeStrategy:                        fosite.HierarchicScopeStrategy,
		AudienceMatchingStrategy:             fosite.DefaultAudienceMatchingStrategy,
		AccessTokenLifespan:                  time.Hour,
		GrantTypeTokenExchangeCanSkipClientAuth: false,
	}
	handler := &rfc8693.Handler{Storage: mockATStorageProvider, Strategy: mockStrategyProvider, Config: cfg}

	req := fosite.NewAccessRequest(&fosite.DefaultSession{})
	req.GrantTypes = fosite.Arguments{string(fosite.GrantTypeTokenExchange)}
	req.Client = &fosite.DefaultClient{GrantTypes: []string{string(fosite.GrantTypeTokenExchange)}}
	req.GrantScope("openid")
	resp := fosite.NewAccessResponse()

	mockStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return("new-token", "sig", nil)
	mockATStorage.EXPECT().CreateAccessTokenSession(gomock.Any(), "sig", gomock.Any()).Return(nil)

	err := handler.PopulateTokenEndpointResponse(context.Background(), req, resp)
	require.NoError(t, err)
	assert.Equal(t, "new-token", resp.GetAccessToken())
	assert.Equal(t, "bearer", resp.GetTokenType())
	assert.Equal(t, rfc8693.TokenTypeAccessToken, resp.GetExtra("issued_token_type"))
}
