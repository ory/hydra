// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	gomock "go.uber.org/mock/gomock"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/storage"
)

func TestRefreshFlow_HandleTokenEndpointRequest(t *testing.T) {
	var areq *fosite.AccessRequest
	sess := &fosite.DefaultSession{Subject: "othersub"}
	expiredSess := &fosite.DefaultSession{
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.RefreshToken: time.Now().UTC().Add(-time.Hour),
		},
	}

	for k, strategy := range map[string]oauth2.RefreshTokenStrategy{
		"hmac": hmacshaStrategy,
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()
			var handler *oauth2.RefreshTokenGrantHandler
			for _, c := range []struct {
				description string
				setup       func(config *fosite.Config)
				expectErr   error
				expect      func(t *testing.T)
			}{
				{
					description: "should fail because not responsible",
					expectErr:   fosite.ErrUnknownRequest,
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"123"}
					},
				},
				{
					description: "should fail because token invalid",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"refresh_token"}}

						areq.Form.Add("refresh_token", "some.refreshtokensig")
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					description: "should fail because token is valid but does not exist",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"refresh_token"}}

						token, _, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)
						areq.Form.Add("refresh_token", token)
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					description: "should fail because client mismatches",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							Client:       &fosite.DefaultClient{ID: ""},
							GrantedScope: []string{"offline"},
							Session:      sess,
						})
						require.NoError(t, err)
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					description: "should fail because token is expired",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
							Scopes:     []string{"foo", "bar", "offline"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo", "offline"},
							RequestedScope: fosite.Arguments{"foo", "bar", "offline"},
							Session:        expiredSess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						})
						require.NoError(t, err)
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					description: "should fail because offline scope has been granted but client no longer allowed to request it",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo", "offline"},
							RequestedScope: fosite.Arguments{"foo", "offline"},
							Session:        sess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						})
						require.NoError(t, err)
					},
					expectErr: fosite.ErrInvalidScope,
				},
				{
					description: "should pass",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
							Scopes:     []string{"foo", "bar", "offline"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)

						orReqID := areq.GetID() + "_OR"
						areq.Form.Add("or_request_id", orReqID)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							ID:             orReqID,
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo", "offline"},
							RequestedScope: fosite.Arguments{"foo", "bar", "offline"},
							Session:        sess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						})
						require.NoError(t, err)
					},
					expect: func(t *testing.T) {
						assert.NotEqual(t, sess, areq.Session)
						assert.NotEqual(t, time.Now().UTC().Add(-time.Hour).Round(time.Hour), areq.RequestedAt)
						assert.Equal(t, fosite.Arguments{"foo", "offline"}, areq.GrantedScope)
						assert.Equal(t, fosite.Arguments{"foo", "bar", "offline"}, areq.RequestedScope)
						assert.NotEqual(t, url.Values{"foo": []string{"bar"}}, areq.Form)
						assert.Equal(t, time.Now().Add(time.Hour).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.AccessToken))
						assert.Equal(t, time.Now().Add(time.Hour).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.RefreshToken))
						assert.EqualValues(t, areq.Form.Get("or_request_id"), areq.GetID(), "Requester ID should be replaced based on the refresh token session")
					},
				},
				{
					description: "should pass with custom client lifespans",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClientWithCustomTokenLifespans{
							DefaultClient: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: fosite.Arguments{"refresh_token"},
								Scopes:     []string{"foo", "bar", "offline"},
							},
						}

						areq.Client.(*fosite.DefaultClientWithCustomTokenLifespans).SetTokenLifespans(&internal.TestLifespans)

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo", "offline"},
							RequestedScope: fosite.Arguments{"foo", "bar", "offline"},
							Session:        sess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						})
						require.NoError(t, err)
					},
					expect: func(t *testing.T) {
						assert.NotEqual(t, sess, areq.Session)
						assert.NotEqual(t, time.Now().UTC().Add(-time.Hour).Round(time.Hour), areq.RequestedAt)
						assert.Equal(t, fosite.Arguments{"foo", "offline"}, areq.GrantedScope)
						assert.Equal(t, fosite.Arguments{"foo", "bar", "offline"}, areq.RequestedScope)
						assert.NotEqual(t, url.Values{"foo": []string{"bar"}}, areq.Form)
						internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.RefreshTokenGrantAccessTokenLifespan).UTC(), areq.GetSession().GetExpiresAt(fosite.AccessToken), time.Minute)
						internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.RefreshTokenGrantRefreshTokenLifespan).UTC(), areq.GetSession().GetExpiresAt(fosite.RefreshToken), time.Minute)
					},
				},
				{
					description: "should fail without offline scope",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
							Scopes:     []string{"foo", "bar"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo"},
							RequestedScope: fosite.Arguments{"foo", "bar"},
							Session:        sess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						})
						require.NoError(t, err)
					},
					expectErr: fosite.ErrScopeNotGranted,
				},
				{
					description: "should pass without offline scope when configured to allow refresh tokens",
					setup: func(config *fosite.Config) {
						config.RefreshTokenScopes = []string{}
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
							Scopes:     []string{"foo", "bar"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", &fosite.Request{
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo"},
							RequestedScope: fosite.Arguments{"foo", "bar"},
							Session:        sess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						})
						require.NoError(t, err)
					},
					expect: func(t *testing.T) {
						assert.NotEqual(t, sess, areq.Session)
						assert.NotEqual(t, time.Now().UTC().Add(-time.Hour).Round(time.Hour), areq.RequestedAt)
						assert.Equal(t, fosite.Arguments{"foo"}, areq.GrantedScope)
						assert.Equal(t, fosite.Arguments{"foo", "bar"}, areq.RequestedScope)
						assert.NotEqual(t, url.Values{"foo": []string{"bar"}}, areq.Form)
						assert.Equal(t, time.Now().Add(time.Hour).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.AccessToken))
						assert.Equal(t, time.Now().Add(time.Hour).UTC().Round(time.Second), areq.GetSession().GetExpiresAt(fosite.RefreshToken))
					},
				},
				{
					description: "should deny access on token reuse",
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.Client = &fosite.DefaultClient{
							ID:         "foo",
							GrantTypes: fosite.Arguments{"refresh_token"},
							Scopes:     []string{"foo", "bar", "offline"},
						}

						token, sig, err := strategy.GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)

						areq.Form.Add("refresh_token", token)
						req := &fosite.Request{
							Client:         areq.Client,
							GrantedScope:   fosite.Arguments{"foo", "offline"},
							RequestedScope: fosite.Arguments{"foo", "bar", "offline"},
							Session:        sess,
							Form:           url.Values{"foo": []string{"bar"}},
							RequestedAt:    time.Now().UTC().Add(-time.Hour).Round(time.Hour),
						}
						err = store.CreateRefreshTokenSession(context.Background(), sig, "", req)
						require.NoError(t, err)

						err = store.RevokeRefreshToken(context.Background(), req.ID)
						require.NoError(t, err)
					},
					expectErr: fosite.ErrInvalidGrant,
				},
			} {
				t.Run("case="+c.description, func(t *testing.T) {
					config := &fosite.Config{
						AccessTokenLifespan:      time.Hour,
						RefreshTokenLifespan:     time.Hour,
						ScopeStrategy:            fosite.HierarchicScopeStrategy,
						AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
						RefreshTokenScopes:       []string{"offline"},
					}
					handler = &oauth2.RefreshTokenGrantHandler{
						Storage: store,
						Strategy: strategy.(interface {
							oauth2.AccessTokenStrategyProvider
							oauth2.RefreshTokenStrategyProvider
						}),
						Config: config,
					}

					areq = fosite.NewAccessRequest(&fosite.DefaultSession{})
					areq.Form = url.Values{}
					c.setup(config)

					err := handler.HandleTokenEndpointRequest(context.Background(), areq)
					if c.expectErr != nil {
						require.EqualError(t, err, c.expectErr.Error())
					} else {
						require.NoError(t, err)
					}

					if c.expect != nil {
						c.expect(t)
					}
				})
			}
		})
	}
}

func TestRefreshFlowTransactional_HandleTokenEndpointRequest(t *testing.T) {
	var (
		mockTransactional                  *internal.MockTransactional
		mockTokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider
		mockTokenRevocationStorage         *internal.MockTokenRevocationStorage
		mockAccessTokenStorageProvider     *internal.MockAccessTokenStorageProvider
		mockRefreshTokenStorageProvider    *internal.MockRefreshTokenStorageProvider
		mockRefreshTokenStorage            *internal.MockRefreshTokenStorage
	)

	request := fosite.NewAccessRequest(&fosite.DefaultSession{})
	propagatedContext := context.Background()

	for _, testCase := range []struct {
		description string
		setup       func()
		expectError error
	}{
		{
			description: "should revoke session on token reuse",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				request.Client = &fosite.DefaultClient{
					ID:         "foo",
					GrantTypes: fosite.Arguments{"refresh_token"},
				}
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(2)
				mockRefreshTokenStorage.
					EXPECT().
					GetRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(request, fosite.ErrInactiveToken).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					DeleteRefreshTokenSession(propagatedContext, gomock.Any()).
					Return(nil).
					Times(1)
				mockTokenRevocationStorageProvider.
					EXPECT().
					TokenRevocationStorage().
					Return(mockTokenRevocationStorage).
					Times(2)
				mockTokenRevocationStorage.
					EXPECT().
					RevokeRefreshToken(propagatedContext, gomock.Any()).
					Return(nil).
					Times(1)
				mockTokenRevocationStorage.
					EXPECT().
					RevokeAccessToken(propagatedContext, gomock.Any()).
					Return(nil).
					Times(1)
				mockTransactional.
					EXPECT().
					Commit(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrInvalidGrant,
		},
	} {
		t.Run(fmt.Sprintf("scenario=%s", testCase.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockTransactional = internal.NewMockTransactional(ctrl)

			mockTokenRevocationStorageProvider = internal.NewMockTokenRevocationStorageProvider(ctrl)
			mockTokenRevocationStorage = internal.NewMockTokenRevocationStorage(ctrl)

			mockAccessTokenStorageProvider = internal.NewMockAccessTokenStorageProvider(ctrl)

			mockRefreshTokenStorageProvider = internal.NewMockRefreshTokenStorageProvider(ctrl)
			mockRefreshTokenStorage = internal.NewMockRefreshTokenStorage(ctrl)

			// define concrete types
			mockStorage := struct {
				*internal.MockAccessTokenStorageProvider
				*internal.MockRefreshTokenStorageProvider
				*internal.MockTokenRevocationStorageProvider
				*internal.MockTransactional
			}{
				MockAccessTokenStorageProvider:     mockAccessTokenStorageProvider,
				MockRefreshTokenStorageProvider:    mockRefreshTokenStorageProvider,
				MockTokenRevocationStorageProvider: mockTokenRevocationStorageProvider,
				MockTransactional:                  mockTransactional,
			}

			handler := oauth2.RefreshTokenGrantHandler{
				Storage:  mockStorage,
				Strategy: hmacshaStrategy,
				Config: &fosite.Config{
					AccessTokenLifespan:      time.Hour,
					ScopeStrategy:            fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
				},
			}

			testCase.setup()

			if err := handler.HandleTokenEndpointRequest(propagatedContext, request); testCase.expectError != nil {
				assert.EqualError(t, err, testCase.expectError.Error())
			}
		})
	}
}

func TestRefreshFlow_PopulateTokenEndpointResponse(t *testing.T) {
	var areq *fosite.AccessRequest
	var aresp *fosite.AccessResponse

	for k, strategy := range map[string]oauth2.CoreStrategy{
		"hmac": hmacshaStrategy,
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			for _, c := range []struct {
				description string
				setup       func(config *fosite.Config)
				check       func(t *testing.T)
				expectErr   error
			}{
				{
					description: "should fail because not responsible",
					expectErr:   fosite.ErrUnknownRequest,
					setup: func(config *fosite.Config) {
						areq.GrantTypes = fosite.Arguments{"313"}
					},
				},
				{
					description: "should pass",
					setup: func(config *fosite.Config) {
						areq.ID = "req-id"
						areq.GrantTypes = fosite.Arguments{"refresh_token"}
						areq.RequestedScope = fosite.Arguments{"foo", "bar"}
						areq.GrantedScope = fosite.Arguments{"foo", "bar"}

						token, signature, err := strategy.RefreshTokenStrategy().GenerateRefreshToken(context.Background(), nil)
						require.NoError(t, err)
						require.NoError(t, store.CreateRefreshTokenSession(context.Background(), signature, "", areq))
						areq.Form.Add("refresh_token", token)
					},
					check: func(t *testing.T) {
						signature := strategy.RefreshTokenStrategy().RefreshTokenSignature(context.Background(), areq.Form.Get("refresh_token"))

						// The old refresh token should be deleted
						_, err := store.GetRefreshTokenSession(context.Background(), signature, nil)
						require.Error(t, err)

						assert.Equal(t, "req-id", areq.ID)
						require.NoError(t, strategy.AccessTokenStrategy().ValidateAccessToken(context.Background(), areq, aresp.GetAccessToken()))
						require.NoError(t, strategy.RefreshTokenStrategy().ValidateRefreshToken(context.Background(), areq, aresp.ToMap()["refresh_token"].(string)))
						assert.Equal(t, "bearer", aresp.GetTokenType())
						assert.NotEmpty(t, aresp.ToMap()["expires_in"])
						assert.Equal(t, "foo bar", aresp.ToMap()["scope"])
					},
				},
			} {
				t.Run("case="+c.description, func(t *testing.T) {
					config := &fosite.Config{
						AccessTokenLifespan:      time.Hour,
						ScopeStrategy:            fosite.HierarchicScopeStrategy,
						AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
					}
					h := oauth2.RefreshTokenGrantHandler{
						Storage: store,
						Strategy: strategy.(interface {
							oauth2.AccessTokenStrategyProvider
							oauth2.RefreshTokenStrategyProvider
						}),
						Config: config,
					}
					areq = fosite.NewAccessRequest(&fosite.DefaultSession{})
					aresp = fosite.NewAccessResponse()
					areq.Client = &fosite.DefaultClient{}
					areq.Form = url.Values{}

					c.setup(config)

					err := h.PopulateTokenEndpointResponse(context.Background(), areq, aresp)
					if c.expectErr != nil {
						assert.EqualError(t, err, c.expectErr.Error())
					} else {
						assert.NoError(t, err)
					}

					if c.check != nil {
						c.check(t)
					}
				})
			}
		})
	}
}

func TestRefreshFlowTransactional_PopulateTokenEndpointResponse(t *testing.T) {
	var (
		mockTransactional                  *internal.MockTransactional
		mockTokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider
		mockAccessTokenStorageProvider     *internal.MockAccessTokenStorageProvider
		mockAccessTokenStorage             *internal.MockAccessTokenStorage
		mockRefreshTokenStorageProvider    *internal.MockRefreshTokenStorageProvider
		mockRefreshTokenStorage            *internal.MockRefreshTokenStorage
	)

	request := fosite.NewAccessRequest(&fosite.DefaultSession{})
	response := fosite.NewAccessResponse()
	propagatedContext := context.Background()

	for _, testCase := range []struct {
		description string
		setup       func()
		expectError error
	}{
		{
			description: "transaction should be committed successfully if no errors occur",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(2)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockTransactional.
					EXPECT().
					Commit(propagatedContext).
					Return(nil).
					Times(1)
			},
		},
		{
			description: "transaction should be rolled back if call to `RevokeAccessToken` results in an error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
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
			description: "should result in a fosite.ErrInvalidRequest if call to `RevokeAccessToken` results in a " +
				"fosite.ErrSerializationFailure error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(fosite.ErrSerializationFailure).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrInvalidRequest,
		},
		{
			description: "transaction should be rolled back if call to `RotateRefreshToken` results in an error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
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
			description: "should result in a fosite.ErrInvalidRequest if call to `RotateRefreshToken` results in a " +
				"fosite.ErrSerializationFailure error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(fosite.ErrSerializationFailure).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrInvalidRequest,
		},
		{
			description: "should result in a fosite.ErrInvalidRequest if call to `CreateAccessTokenSession` results in " +
				"a fosite.ErrSerializationFailure error",
			setup: func() {
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(fosite.ErrSerializationFailure).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrInvalidRequest,
		},
		{
			description: "transaction should be rolled back if call to `CreateAccessTokenSession` results in an error",
			setup: func() {
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
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
			description: "transaction should be rolled back if call to `CreateRefreshTokenSession` results in an error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(2)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
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
			description: "should result in a fosite.ErrInvalidRequest if call to `CreateRefreshTokenSession` results in " +
				"a fosite.ErrSerializationFailure error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(2)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fosite.ErrSerializationFailure).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: fosite.ErrInvalidRequest,
		},
		{
			description: "should result in a server error if transaction cannot be created",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(nil, errors.New("Could not create transaction!")).
					Times(1)
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should result in a server error if transaction cannot be rolled back",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(fosite.ErrNotFound).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(errors.New("Could not rollback transaction!")).
					Times(1)
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should result in a server error if transaction cannot be committed",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(2)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockTransactional.
					EXPECT().
					Commit(propagatedContext).
					Return(errors.New("Could not commit transaction!")).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: nil,
		},
		{
			description: "should result in a `fosite.ErrInvalidRequest` if transaction fails to commit due to a " +
				"`fosite.ErrSerializationFailure` error",
			setup: func() {
				request.GrantTypes = fosite.Arguments{"refresh_token"}
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockRefreshTokenStorageProvider.
					EXPECT().
					RefreshTokenStorage().
					Return(mockRefreshTokenStorage).
					Times(2)
				mockRefreshTokenStorage.
					EXPECT().
					RotateRefreshToken(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockAccessTokenStorageProvider.
					EXPECT().
					AccessTokenStorage().
					Return(mockAccessTokenStorage).
					Times(1)
				mockAccessTokenStorage.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockRefreshTokenStorage.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockTransactional.
					EXPECT().
					Commit(propagatedContext).
					Return(fosite.ErrSerializationFailure).
					Times(1)
				mockTransactional.
					EXPECT().
					Rollback(propagatedContext).
					Return(nil).
					Times(1)
			},
			expectError: nil,
		},
	} {
		t.Run(fmt.Sprintf("scenario=%s", testCase.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockTransactional = internal.NewMockTransactional(ctrl)

			mockTokenRevocationStorageProvider = internal.NewMockTokenRevocationStorageProvider(ctrl)

			mockAccessTokenStorageProvider = internal.NewMockAccessTokenStorageProvider(ctrl)
			mockAccessTokenStorage = internal.NewMockAccessTokenStorage(ctrl)

			mockRefreshTokenStorageProvider = internal.NewMockRefreshTokenStorageProvider(ctrl)
			mockRefreshTokenStorage = internal.NewMockRefreshTokenStorage(ctrl)

			// define concrete types
			mockStorage := struct {
				*internal.MockAccessTokenStorageProvider
				*internal.MockRefreshTokenStorageProvider
				*internal.MockTokenRevocationStorageProvider
				*internal.MockTransactional
			}{
				MockAccessTokenStorageProvider:     mockAccessTokenStorageProvider,
				MockRefreshTokenStorageProvider:    mockRefreshTokenStorageProvider,
				MockTokenRevocationStorageProvider: mockTokenRevocationStorageProvider,
				MockTransactional:                  mockTransactional,
			}

			handler := oauth2.RefreshTokenGrantHandler{
				Storage:  mockStorage,
				Strategy: hmacshaStrategy,
				Config: &fosite.Config{
					AccessTokenLifespan:      time.Hour,
					ScopeStrategy:            fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
				},
			}

			testCase.setup()

			if err := handler.PopulateTokenEndpointResponse(propagatedContext, request, response); testCase.expectError != nil {
				assert.EqualError(t, err, testCase.expectError.Error())
			}
		})
	}
}
