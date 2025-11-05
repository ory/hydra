// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	//"time"

	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite" //"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/storage"
)

func TestAuthorizeCode_PopulateTokenEndpointResponse(t *testing.T) {
	for k, strategy := range map[string]oauth2.CoreStrategy{
		"hmac": hmacshaStrategy,
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			var h oauth2.AuthorizeExplicitGrantHandler
			for _, c := range []struct {
				areq        *fosite.AccessRequest
				description string
				setup       func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config)
				check       func(t *testing.T, aresp *fosite.AccessResponse)
				expectErr   error
			}{
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"123"},
					},
					description: "should fail because not responsible",
					expectErr:   fosite.ErrUnknownRequest,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					description: "should fail because authcode not found",
					setup: func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config) {
						code, _, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Set("code", code)
					},
					expectErr: fosite.ErrServerError,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{"code": []string{"foo.bar"}},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					description: "should fail because validation failed",
					setup: func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config) {
						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), "bar", areq))
					},
					expectErr: fosite.ErrInvalidRequest,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code", "refresh_token"},
							},
							GrantedScope: fosite.Arguments{"foo", "offline"},
							Session:      &fosite.DefaultSession{},
							RequestedAt:  time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config) {
						code, sig, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Add("code", code)

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), sig, areq))
					},
					description: "should pass with offline scope and refresh token",
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.NotEmpty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Equal(t, "foo offline", aresp.GetExtra("scope"))
					},
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code", "refresh_token"},
							},
							GrantedScope: fosite.Arguments{"foo"},
							Session:      &fosite.DefaultSession{},
							RequestedAt:  time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config) {
						config.RefreshTokenScopes = []string{}
						code, sig, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Add("code", code)

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), sig, areq))
					},
					description: "should pass with refresh token always provided",
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.NotEmpty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Equal(t, "foo", aresp.GetExtra("scope"))
					},
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code"},
							},
							GrantedScope: fosite.Arguments{},
							Session:      &fosite.DefaultSession{},
							RequestedAt:  time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config) {
						config.RefreshTokenScopes = []string{}
						code, sig, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Add("code", code)

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), sig, areq))
					},
					description: "should pass with no refresh token",
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.Empty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Empty(t, aresp.GetExtra("scope"))
					},
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code"},
							},
							GrantedScope: fosite.Arguments{"foo"},
							Session:      &fosite.DefaultSession{},
							RequestedAt:  time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, config *fosite.Config) {
						code, sig, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Add("code", code)

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), sig, areq))
					},
					description: "should not have refresh token",
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.Empty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Equal(t, "foo", aresp.GetExtra("scope"))
					},
				},
			} {
				t.Run("case="+c.description, func(t *testing.T) {
					config := &fosite.Config{
						ScopeStrategy:            fosite.HierarchicScopeStrategy,
						AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
						AccessTokenLifespan:      time.Minute,
						RefreshTokenScopes:       []string{"offline"},
					}
					h = oauth2.AuthorizeExplicitGrantHandler{
						Storage:  store,
						Strategy: strategy,
						Config:   config,
					}

					if c.setup != nil {
						c.setup(t, c.areq, config)
					}

					aresp := fosite.NewAccessResponse()
					err := h.PopulateTokenEndpointResponse(context.Background(), c.areq, aresp)

					if c.expectErr != nil {
						require.EqualError(t, err, c.expectErr.Error(), "%+v", err)
					} else {
						require.NoError(t, err, "%+v", err)
					}

					if c.check != nil {
						c.check(t, aresp)
					}
				})
			}
		})
	}
}

func TestAuthorizeCode_HandleTokenEndpointRequest(t *testing.T) {
	for k, strategy := range map[string]oauth2.CoreStrategy{
		"hmac": hmacshaStrategy,
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			h := oauth2.AuthorizeExplicitGrantHandler{
				Storage:  store,
				Strategy: hmacshaStrategy,
				Config: &fosite.Config{
					ScopeStrategy:            fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
					AuthorizeCodeLifespan:    time.Minute,
				},
			}
			for i, c := range []struct {
				areq        *fosite.AccessRequest
				authreq     *fosite.AuthorizeRequest
				description string
				setup       func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest)
				check       func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest)
				expectErr   error
			}{
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"12345678"},
					},
					description: "should fail because not responsible",
					expectErr:   fosite.ErrUnknownRequest,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client:      &fosite.DefaultClient{ID: "foo", GrantTypes: []string{""}},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					description: "should fail because client is not granted this grant type",
					expectErr:   fosite.ErrUnauthorizedClient,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client:      &fosite.DefaultClient{GrantTypes: []string{"authorization_code"}},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					description: "should fail because authcode could not be retrieved (1)",
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest) {
						token, _, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form = url.Values{"code": {token}}
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form:        url.Values{"code": {"foo.bar"}},
							Client:      &fosite.DefaultClient{GrantTypes: []string{"authorization_code"}},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					description: "should fail because authcode validation failed",
					expectErr:   fosite.ErrInvalidGrant,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client:      &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"authorization_code"}},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.AuthorizeRequest{
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "bar"},
							RequestedScope: fosite.Arguments{"a", "b"},
						},
					},
					description: "should fail because client mismatch",
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest) {
						token, signature, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form = url.Values{"code": {token}}

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), signature, authreq))
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client:      &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"authorization_code"}},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.AuthorizeRequest{
						Request: fosite.Request{
							Client:  &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"authorization_code"}},
							Form:    url.Values{"redirect_uri": []string{"request-redir"}},
							Session: &fosite.DefaultSession{},
						},
					},
					description: "should fail because redirect uri was set during /authorize call, but not in /token call",
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest) {
						token, signature, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form = url.Values{"code": {token}}

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), signature, authreq))
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client:      &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"authorization_code"}},
							Form:        url.Values{"redirect_uri": []string{"request-redir"}},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.AuthorizeRequest{
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"authorization_code"}},
							Session:        &fosite.DefaultSession{},
							RequestedScope: fosite.Arguments{"a", "b"},
							RequestedAt:    time.Now().UTC(),
						},
					},
					description: "should pass",
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest) {
						token, signature, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)

						areq.Form = url.Values{"code": {token}}
						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), signature, authreq))
					},
				},
				{
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"authorization_code"},
							},
							GrantedScope: fosite.Arguments{"foo", "offline"},
							Session:      &fosite.DefaultSession{},
							RequestedAt:  time.Now().UTC(),
						},
					},
					check: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest) {
						assert.Equal(t, time.Now().Add(time.Minute).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.AccessToken))
						assert.Equal(t, time.Now().Add(time.Minute).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.RefreshToken))
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.AuthorizeRequest) {
						code, sig, err := strategy.AuthorizeCodeStrategy().GenerateAuthorizeCode(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Add("code", code)

						require.NoError(t, store.CreateAuthorizeCodeSession(context.Background(), sig, areq))
						require.NoError(t, store.InvalidateAuthorizeCodeSession(context.Background(), sig))
					},
					description: "should fail because code has been used already",
					expectErr:   fosite.ErrInvalidGrant,
				},
			} {
				t.Run(fmt.Sprintf("case=%d/description=%s", i, c.description), func(t *testing.T) {
					if c.setup != nil {
						c.setup(t, c.areq, c.authreq)
					}

					t.Logf("Processing %+v", c.areq.Client)

					err := h.HandleTokenEndpointRequest(context.Background(), c.areq)
					if c.expectErr != nil {
						require.EqualError(t, err, c.expectErr.Error(), "%+v", err)
					} else {
						require.NoError(t, err, "%+v", err)
						if c.check != nil {
							c.check(t, c.areq, c.authreq)
						}
					}
				})
			}
		})
	}
}

func TestAuthorizeCodeTransactional_HandleTokenEndpointRequest(t *testing.T) {
	token, _, err := hmacshaStrategy.GenerateAuthorizeCode(context.Background(), nil)
	require.NoError(t, err)

	request := &fosite.AccessRequest{
		GrantTypes: fosite.Arguments{"authorization_code"},
		Request: fosite.Request{
			Client: &fosite.DefaultClient{
				GrantTypes: fosite.Arguments{"authorization_code", "refresh_token"},
			},
			GrantedScope: fosite.Arguments{"offline"},
			Session:      &fosite.DefaultSession{},
			RequestedAt:  time.Now().UTC(),
		},
	}
	request.Form = url.Values{"code": {token}}
	response := fosite.NewAccessResponse()
	propagatedContext := context.Background()

	for k, c := range []struct {
		description string
		setup       func(
			mockTransactional *internal.MockTransactional,
			tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
			tokenRevocationStorage *internal.MockTokenRevocationStorage,
			authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
			authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
			accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
			accessTokenStorage *internal.MockAccessTokenStorage,
			refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
			refreshTokenStorage *internal.MockRefreshTokenStorage,
			authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
			authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
			accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
			accessTokenStrategy *internal.MockAccessTokenStrategy,
			refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
			refreshTokenStrategy *internal.MockRefreshTokenStrategy,
		)
		expectError error
	}{
		{
			description: "transaction should be committed successfully if no errors occur",
			setup: func(
				mockTransactional *internal.MockTransactional,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
				authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
				authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			) {
				authorizeCodeStrategyProvider.EXPECT().AuthorizeCodeStrategy().Return(authorizeCodeStrategy).Times(2)
				authorizeCodeStrategy.EXPECT().AuthorizeCodeSignature(gomock.Any(), gomock.Any())
				authorizeCodeStrategy.EXPECT().ValidateAuthorizeCode(gomock.Any(), gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the authorize code storage mock
				authorizeCodeStorageProvider.
					EXPECT().
					AuthorizeCodeStorage().
					Return(authorizeCodeStorage).
					Times(2)

				// Set up authorize code storage expectations
				authorizeCodeStorage.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				authorizeCodeStorage.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the access token storage mock
				accessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(accessTokenStorage).
					Times(1)

				// Set up access token storage expectations
				accessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the refresh token storage mock
				refreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(refreshTokenStorage).
					Times(0)

				// Set up refresh token storage expectations
				refreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)

				// Set up transaction expectations
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockTransactional.
					EXPECT().
					Commit(propagatedContext).
					Return(nil).
					Times(1)
			},
		},
		{
			description: "transaction should be rolled back if `InvalidateAuthorizeCodeSession` returns an error",
			setup: func(
				mockTransactional *internal.MockTransactional,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
				authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
				authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			) {
				authorizeCodeStrategyProvider.EXPECT().AuthorizeCodeStrategy().Return(authorizeCodeStrategy).Times(2)
				authorizeCodeStrategy.EXPECT().AuthorizeCodeSignature(gomock.Any(), gomock.Any())
				authorizeCodeStrategy.EXPECT().ValidateAuthorizeCode(gomock.Any(), gomock.Any(), gomock.Any())

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any())

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the authorize code storage mock
				authorizeCodeStorageProvider.
					EXPECT().
					AuthorizeCodeStorage().
					Return(authorizeCodeStorage).
					Times(2)

				// Set up authorize code storage expectations
				authorizeCodeStorage.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				authorizeCodeStorage.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
					Times(1)

				// Set up transaction expectations
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "transaction should be rolled back if `CreateAccessTokenSession` returns an error",
			setup: func(
				mockTransactional *internal.MockTransactional,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
				authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
				authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			) {
				authorizeCodeStrategyProvider.EXPECT().AuthorizeCodeStrategy().Return(authorizeCodeStrategy).Times(2)
				authorizeCodeStrategy.EXPECT().AuthorizeCodeSignature(gomock.Any(), gomock.Any())
				authorizeCodeStrategy.EXPECT().ValidateAuthorizeCode(gomock.Any(), gomock.Any(), gomock.Any())

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any())

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the authorize code storage mock
				authorizeCodeStorageProvider.
					EXPECT().
					AuthorizeCodeStorage().
					Return(authorizeCodeStorage).
					Times(2)

				// Set up authorize code storage expectations
				authorizeCodeStorage.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				authorizeCodeStorage.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				// Set up CoreStorage to return the access token storage mock
				accessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(accessTokenStorage).
					Times(1)

				// Set up access token storage expectations
				accessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
					Times(1)

				// Set up transaction expectations
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should result in a server error if transaction cannot be created",
			setup: func(
				mockTransactional *internal.MockTransactional,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
				authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
				authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			) {
				authorizeCodeStrategyProvider.EXPECT().AuthorizeCodeStrategy().Return(authorizeCodeStrategy).Times(2)
				authorizeCodeStrategy.EXPECT().AuthorizeCodeSignature(gomock.Any(), gomock.Any())
				authorizeCodeStrategy.EXPECT().ValidateAuthorizeCode(gomock.Any(), gomock.Any(), gomock.Any())

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any())

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the authorize code storage mock
				authorizeCodeStorageProvider.
					EXPECT().
					AuthorizeCodeStorage().
					Return(authorizeCodeStorage).
					Times(1)

				// Set up authorize code storage expectations
				authorizeCodeStorage.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)

				// Set up transaction expectations
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(nil, errors.New("Whoops, unable to create transaction!"))
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should result in a server error if transaction cannot be rolled back",
			setup: func(
				mockTransactional *internal.MockTransactional,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
				authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
				authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			) {
				authorizeCodeStrategyProvider.EXPECT().AuthorizeCodeStrategy().Return(authorizeCodeStrategy).Times(2)
				authorizeCodeStrategy.EXPECT().AuthorizeCodeSignature(gomock.Any(), gomock.Any())
				authorizeCodeStrategy.EXPECT().ValidateAuthorizeCode(gomock.Any(), gomock.Any(), gomock.Any())

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any())

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the authorize code storage mock
				authorizeCodeStorageProvider.
					EXPECT().
					AuthorizeCodeStorage().
					Return(authorizeCodeStorage).
					Times(2)

				// Set up authorize code storage expectations
				authorizeCodeStorage.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				authorizeCodeStorage.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
					Times(1)

				// Set up transaction expectations
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(errors.New("Whoops, unable to rollback transaction!")).
					Times(1)
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should result in a server error if transaction cannot be committed",
			setup: func(
				mockTransactional *internal.MockTransactional,
				tokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider,
				tokenRevocationStorage *internal.MockTokenRevocationStorage,
				authorizeCodeStorageProvider *internal.MockAuthorizeCodeStorageProvider,
				authorizeCodeStorage *internal.MockAuthorizeCodeStorage,
				accessTokenStorageProvider *internal.MockAccessTokenStorageProvider,
				accessTokenStorage *internal.MockAccessTokenStorage,
				refreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider,
				refreshTokenStorage *internal.MockRefreshTokenStorage,
				authorizeCodeStrategyProvider *internal.MockAuthorizeCodeStrategyProvider,
				authorizeCodeStrategy *internal.MockAuthorizeCodeStrategy,
				accessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider,
				accessTokenStrategy *internal.MockAccessTokenStrategy,
				refreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider,
				refreshTokenStrategy *internal.MockRefreshTokenStrategy,
			) {
				authorizeCodeStrategyProvider.EXPECT().AuthorizeCodeStrategy().Return(authorizeCodeStrategy).Times(2)
				authorizeCodeStrategy.EXPECT().AuthorizeCodeSignature(gomock.Any(), gomock.Any())
				authorizeCodeStrategy.EXPECT().ValidateAuthorizeCode(gomock.Any(), gomock.Any(), gomock.Any())

				accessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(accessTokenStrategy).Times(1)
				accessTokenStrategy.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any())

				refreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(refreshTokenStrategy).Times(1)
				refreshTokenStrategy.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any())

				// Set up CoreStorage to return the authorize code storage mock
				authorizeCodeStorageProvider.
					EXPECT().
					AuthorizeCodeStorage().
					Return(authorizeCodeStorage).
					Times(2)

				// Set up authorize code storage expectations
				authorizeCodeStorage.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				authorizeCodeStorage.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				// Set up CoreStorage to return the access token storage mock
				accessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(accessTokenStorage).
					Times(1)

				// Set up access token storage expectations
				accessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				// Set up CoreStorage to return the refresh token storage mock
				refreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(refreshTokenStorage).
					Times(0)

				// Set up refresh token storage expectations
				refreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)

				// Set up transaction expectations
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockTransactional.
					EXPECT().
					Commit(propagatedContext).
					Return(errors.New("Whoops, unable to commit transaction!")).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrServerError,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			// Initialize all mocks
			mockTransactional := internal.NewMockTransactional(ctrl)

			tokenRevocationStorageProvider := internal.NewMockTokenRevocationStorageProvider(ctrl)
			tokenRevocationStorage := internal.NewMockTokenRevocationStorage(ctrl)

			authorizeCodeStorageProvider := internal.NewMockAuthorizeCodeStorageProvider(ctrl)
			authorizeCodeStorage := internal.NewMockAuthorizeCodeStorage(ctrl)

			accessTokenStorageProvider := internal.NewMockAccessTokenStorageProvider(ctrl)
			accessTokenStorage := internal.NewMockAccessTokenStorage(ctrl)

			refreshTokenStorageProvider := internal.NewMockRefreshTokenStorageProvider(ctrl)
			refreshTokenStorage := internal.NewMockRefreshTokenStorage(ctrl)

			authorizeCodeStrategyProvider := internal.NewMockAuthorizeCodeStrategyProvider(ctrl)
			authorizeCodeStrategy := internal.NewMockAuthorizeCodeStrategy(ctrl)

			accessTokenStrategyProvider := internal.NewMockAccessTokenStrategyProvider(ctrl)
			accessTokenStrategy := internal.NewMockAccessTokenStrategy(ctrl)

			refreshTokenStrategyProvider := internal.NewMockRefreshTokenStrategyProvider(ctrl)
			refreshTokenStrategy := internal.NewMockRefreshTokenStrategy(ctrl)

			// define concrete types
			mockStorage := struct {
				*internal.MockAuthorizeCodeStorageProvider
				*internal.MockAccessTokenStorageProvider
				*internal.MockRefreshTokenStorageProvider
				*internal.MockTokenRevocationStorageProvider
				*internal.MockTransactional
			}{
				MockAuthorizeCodeStorageProvider:   authorizeCodeStorageProvider,
				MockAccessTokenStorageProvider:     accessTokenStorageProvider,
				MockRefreshTokenStorageProvider:    refreshTokenStorageProvider,
				MockTokenRevocationStorageProvider: tokenRevocationStorageProvider,
				MockTransactional:                  mockTransactional,
			}

			mockStrategy := struct {
				*internal.MockAuthorizeCodeStrategyProvider
				*internal.MockAccessTokenStrategyProvider
				*internal.MockRefreshTokenStrategyProvider
			}{
				MockAuthorizeCodeStrategyProvider: authorizeCodeStrategyProvider,
				MockAccessTokenStrategyProvider:   accessTokenStrategyProvider,
				MockRefreshTokenStrategyProvider:  refreshTokenStrategyProvider,
			}

			handler := oauth2.AuthorizeExplicitGrantHandler{
				Storage:  mockStorage,
				Strategy: mockStrategy,
				Config: &fosite.Config{
					ScopeStrategy:            fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
					AuthorizeCodeLifespan:    time.Minute,
				},
			}

			// set up mock expectations
			c.setup(
				mockTransactional,
				tokenRevocationStorageProvider,
				tokenRevocationStorage,
				authorizeCodeStorageProvider,
				authorizeCodeStorage,
				accessTokenStorageProvider,
				accessTokenStorage,
				refreshTokenStorageProvider,
				refreshTokenStorage,
				authorizeCodeStrategyProvider,
				authorizeCodeStrategy,
				accessTokenStrategyProvider,
				accessTokenStrategy,
				refreshTokenStrategyProvider,
				refreshTokenStrategy,
			)

			// invoke function under test
			if err := handler.PopulateTokenEndpointResponse(propagatedContext, request, response); c.expectError != nil {
				assert.EqualError(t, err, c.expectError.Error())
			}
		})
	}
}
