// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"net/url"
	"testing" //"time"

	"time"

	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite/internal"

	"github.com/ory/hydra/v2/fosite" //"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizeCode_PopulateTokenEndpointResponse(t *testing.T) {
	for k, strategy := range map[string]CoreStrategy{
		"hmac": hmacshaStrategy,
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			var h AuthorizeExplicitGrantHandler
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
						code, _, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						code, sig, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						code, sig, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						code, sig, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						code, sig, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
					h = AuthorizeExplicitGrantHandler{
						CoreStorage:           store,
						AuthorizeCodeStrategy: strategy,
						AccessTokenStrategy:   strategy,
						RefreshTokenStrategy:  strategy,
						Config:                config,
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
	for k, strategy := range map[string]CoreStrategy{
		"hmac": hmacshaStrategy,
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			h := AuthorizeExplicitGrantHandler{
				CoreStorage:            store,
				AuthorizeCodeStrategy:  hmacshaStrategy,
				TokenRevocationStorage: store,
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
						token, _, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						token, signature, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						token, signature, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						token, signature, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
						code, sig, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
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
	var mockTransactional *internal.MockTransactional
	var mockCoreStore *internal.MockCoreStorage
	strategy := hmacshaStrategy
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
	token, _, err := strategy.GenerateAuthorizeCode(context.Background(), nil)
	require.NoError(t, err)
	request.Form = url.Values{"code": {token}}
	response := fosite.NewAccessResponse()
	propagatedContext := context.Background()

	// some storage implementation that has support for transactions, notice the embedded type `storage.Transactional`
	type transactionalStore struct {
		storage.Transactional
		CoreStorage
	}

	for _, testCase := range []struct {
		description string
		setup       func()
		expectError error
	}{
		{
			description: "transaction should be committed successfully if no errors occur",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockCoreStore.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockCoreStore.
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
			description: "transaction should be rolled back if `InvalidateAuthorizeCodeSession` returns an error",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
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
			description: "transaction should be rolled back if `CreateAccessTokenSession` returns an error",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockCoreStore.
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
			description: "should result in a server error if transaction cannot be created",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(nil, errors.New("Whoops, unable to create transaction!"))
			},
			expectError: fosite.ErrServerError,
		},
		{
			description: "should result in a server error if transaction cannot be rolled back",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(errors.New("Whoops, a nasty database error occurred!")).
					Times(1)
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
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetAuthorizeCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(request, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateAuthorizeCodeSession(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockCoreStore.
					EXPECT().
					CreateAccessTokenSession(propagatedContext, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				mockCoreStore.
					EXPECT().
					CreateRefreshTokenSession(propagatedContext, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
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
		t.Run(fmt.Sprintf("scenario=%s", testCase.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransactional = internal.NewMockTransactional(ctrl)
			mockCoreStore = internal.NewMockCoreStorage(ctrl)
			testCase.setup()

			handler := AuthorizeExplicitGrantHandler{
				CoreStorage: transactionalStore{
					mockTransactional,
					mockCoreStore,
				},
				AccessTokenStrategy:   strategy,
				RefreshTokenStrategy:  strategy,
				AuthorizeCodeStrategy: strategy,
				Config: &fosite.Config{
					ScopeStrategy:            fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
					AuthorizeCodeLifespan:    time.Minute,
				},
			}

			if err := handler.PopulateTokenEndpointResponse(propagatedContext, request, response); testCase.expectError != nil {
				assert.EqualError(t, err, testCase.expectError.Error())
			}
		})
	}
}
