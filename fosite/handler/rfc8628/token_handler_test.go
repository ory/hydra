// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite/internal"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/token/hmac"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hmacshaStrategy = oauth2.NewHMACSHAStrategy(
	&hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	&fosite.Config{
		AccessTokenLifespan:   time.Hour * 24,
		AuthorizeCodeLifespan: time.Hour * 24,
	},
)

var RFC8628HMACSHAStrategy = DefaultDeviceStrategy{
	Enigma: &hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	Config: &fosite.Config{
		DeviceAndUserCodeLifespan: time.Minute * 30,
	},
}

func TestDeviceUserCode_HandleTokenEndpointRequest(t *testing.T) {
	for k, strategy := range map[string]struct {
		oauth2.CoreStrategy
		RFC8628CodeStrategy
	}{
		"hmac": {hmacshaStrategy, &RFC8628HMACSHAStrategy},
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			h := DeviceCodeTokenEndpointHandler{
				DeviceRateLimitStrategy: strategy,
				DeviceCodeStrategy:      strategy,
				UserCodeStrategy:        strategy,
				CoreStorage:             store,
				AccessTokenStrategy:     strategy.CoreStrategy,
				RefreshTokenStrategy:    strategy.CoreStrategy,
				Config: &fosite.Config{
					ScopeStrategy:             fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy:  fosite.DefaultAudienceMatchingStrategy,
					DeviceAndUserCodeLifespan: time.Minute,
				},
			}

			testCases := []struct {
				description string
				areq        *fosite.AccessRequest
				authreq     *fosite.DeviceRequest
				setup       func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest)
				check       func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest)
				expectErr   error
			}{
				{
					description: "should fail because not responsible for handling the request",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					expectErr: fosite.ErrUnknownRequest,
				},
				{
					description: "should fail because client is not granted the correct grant type",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{""},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					expectErr: fosite.ErrUnauthorizedClient,
				},
				{
					description: "should fail because device code could not be retrieved",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, _ *fosite.DeviceRequest) {
						deviceCode, _, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						areq.Form = url.Values{"device_code": {deviceCode}}
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					description: "should fail because user has not completed the browser flow",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeUnused,
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							RequestedScope: fosite.Arguments{"foo"},
							GrantedScope:   fosite.Arguments{"foo"},
							Session: &fosite.DefaultSession{
								ExpiresAt: map[fosite.TokenType]time.Time{
									fosite.DeviceCode: time.Now().Add(-time.Hour).UTC(),
								},
							},
							RequestedAt: time.Now().Add(-2 * time.Hour).UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest) {
						code, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)
						areq.Form.Add("device_code", code)

						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
					expectErr: fosite.ErrAuthorizationPending,
				},
				{
					description: "should fail because device code has expired",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							GrantedScope: fosite.Arguments{"foo", "offline"},
							Session:      &fosite.DefaultSession{},
							RequestedAt:  time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeAccepted,
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}},
							RequestedScope: fosite.Arguments{"foo"},
							GrantedScope:   fosite.Arguments{"foo"},
							Session: &fosite.DefaultSession{
								ExpiresAt: map[fosite.TokenType]time.Time{
									fosite.DeviceCode: time.Now().Add(-time.Hour).UTC(),
								},
							},
							RequestedAt: time.Now().Add(-2 * time.Hour).UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest) {
						code, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)
						areq.Form.Add("device_code", code)

						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
					expectErr: fosite.ErrDeviceExpiredToken,
				},
				{
					description: "should fail because client mismatch",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeAccepted,
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "bar"},
							RequestedScope: fosite.Arguments{"foo"},
							GrantedScope:   fosite.Arguments{"foo"},
							Session: &fosite.DefaultSession{
								ExpiresAt: map[fosite.TokenType]time.Time{
									fosite.DeviceCode: time.Now().Add(time.Hour).UTC(),
								},
							},
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest) {
						token, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)
						areq.Form = url.Values{"device_code": {token}}

						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
					expectErr: fosite.ErrInvalidGrant,
				},
				{
					description: "should pass",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeAccepted,
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								ID:         "foo",
								GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							RequestedScope: fosite.Arguments{"foo"},
							GrantedScope:   fosite.Arguments{"foo"},
							Session:        &fosite.DefaultSession{},
							RequestedAt:    time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest) {
						token, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)

						areq.Form = url.Values{"device_code": {token}}
						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
				},
			}

			for i, testCase := range testCases {
				t.Run(fmt.Sprintf("case=%d/description=%s", i, testCase.description), func(t *testing.T) {
					if testCase.setup != nil {
						testCase.setup(t, testCase.areq, testCase.authreq)
					}

					t.Logf("Processing %+v", testCase.areq.Client)

					err := h.HandleTokenEndpointRequest(context.Background(), testCase.areq)
					if testCase.expectErr != nil {
						require.EqualError(t, err, testCase.expectErr.Error(), "%+v", err)
					} else {
						require.NoError(t, err, "%+v", err)
						if testCase.check != nil {
							testCase.check(t, testCase.areq, testCase.authreq)
						}
					}
				})
			}
		})
	}
}

func TestDeviceUserCode_HandleTokenEndpointRequest_RateLimiting(t *testing.T) {
	for k, strategy := range map[string]struct {
		oauth2.CoreStrategy
		RFC8628CodeStrategy
	}{
		"hmac": {hmacshaStrategy, &RFC8628HMACSHAStrategy},
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			h := DeviceCodeTokenEndpointHandler{
				DeviceRateLimitStrategy: strategy,
				DeviceCodeStrategy:      strategy,
				UserCodeStrategy:        strategy,
				CoreStorage:             store,
				AccessTokenStrategy:     strategy.CoreStrategy,
				RefreshTokenStrategy:    strategy.CoreStrategy,
				Config: &fosite.Config{
					ScopeStrategy:             fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy:  fosite.DefaultAudienceMatchingStrategy,
					DeviceAndUserCodeLifespan: time.Minute,
				},
			}
			areq := &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Form: url.Values{},
					Client: &fosite.DefaultClient{
						ID:         "foo",
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
					},
					Session:     &fosite.DefaultSession{},
					RequestedAt: time.Now().UTC(),
				},
			}
			authreq := &fosite.DeviceRequest{
				UserCodeState: fosite.UserCodeAccepted,
				Request: fosite.Request{
					Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}},
					RequestedScope: fosite.Arguments{"foo"},
					GrantedScope:   fosite.Arguments{"foo"},
					Session:        &fosite.DefaultSession{},
					RequestedAt:    time.Now().UTC(),
				},
			}

			token, signature, err := strategy.GenerateDeviceCode(context.TODO())
			require.NoError(t, err)
			_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
			require.NoError(t, err)

			areq.Form = url.Values{"device_code": {token}}
			require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
			err = h.HandleTokenEndpointRequest(context.Background(), areq)
			require.NoError(t, err, "%+v", err)
			err = h.HandleTokenEndpointRequest(context.Background(), areq)
			require.Error(t, fosite.ErrSlowDown, err)
			time.Sleep(10 * time.Second)
			err = h.HandleTokenEndpointRequest(context.Background(), areq)
			require.NoError(t, err, "%+v", err)
		})
	}
}

func TestDeviceUserCode_PopulateTokenEndpointResponse(t *testing.T) {
	for k, strategy := range map[string]struct {
		oauth2.CoreStrategy
		RFC8628CodeStrategy
	}{
		"hmac": {hmacshaStrategy, &RFC8628HMACSHAStrategy},
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			testCases := []struct {
				description string
				areq        *fosite.AccessRequest
				authreq     *fosite.DeviceRequest
				setup       func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest, config *fosite.Config)
				check       func(t *testing.T, aresp *fosite.AccessResponse)
				expectErr   error
			}{
				{
					description: "should fail because not responsible for handling the request",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"authorization_code"},
						Request: fosite.Request{
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					expectErr: fosite.ErrUnknownRequest,
				},
				{
					description: "should fail because device code cannot be retrieved",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, _ *fosite.DeviceRequest, _ *fosite.Config) {
						code, _, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						areq.Form.Set("device_code", code)
					},
					expectErr: fosite.ErrServerError,
				},
				{
					description: "should pass with offline scope and refresh token grant type",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code", "refresh_token"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeAccepted,
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}},
							RequestedScope: fosite.Arguments{"foo", "bar", "offline"},
							GrantedScope:   fosite.Arguments{"foo", "offline"},
							Session:        &fosite.DefaultSession{},
							RequestedAt:    time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest, _ *fosite.Config) {
						code, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)
						areq.Form.Add("device_code", code)

						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.NotEmpty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Equal(t, "foo offline", aresp.GetExtra("scope"))
					},
				},
				{
					description: "should pass with refresh token grant type",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code", "refresh_token"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeAccepted,
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}},
							RequestedScope: fosite.Arguments{"foo", "bar"},
							GrantedScope:   fosite.Arguments{"foo"},
							Session:        &fosite.DefaultSession{},
							RequestedAt:    time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest, config *fosite.Config) {
						config.RefreshTokenScopes = []string{}
						code, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)
						areq.Form.Add("device_code", code)

						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.NotEmpty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Equal(t, "foo", aresp.GetExtra("scope"))
					},
				},
				{
					description: "pass and response should not have refresh token",
					areq: &fosite.AccessRequest{
						GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
						Request: fosite.Request{
							Form: url.Values{},
							Client: &fosite.DefaultClient{
								GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
							},
							Session:     &fosite.DefaultSession{},
							RequestedAt: time.Now().UTC(),
						},
					},
					authreq: &fosite.DeviceRequest{
						UserCodeState: fosite.UserCodeAccepted,
						Request: fosite.Request{
							Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}},
							RequestedScope: fosite.Arguments{"foo", "bar"},
							GrantedScope:   fosite.Arguments{"foo"},
							Session:        &fosite.DefaultSession{},
							RequestedAt:    time.Now().UTC(),
						},
					},
					setup: func(t *testing.T, areq *fosite.AccessRequest, authreq *fosite.DeviceRequest, config *fosite.Config) {
						code, signature, err := strategy.GenerateDeviceCode(context.TODO())
						require.NoError(t, err)
						_, userCodeSignature, err := strategy.GenerateUserCode(context.TODO())
						require.NoError(t, err)
						areq.Form.Add("device_code", code)

						require.NoError(t, store.CreateDeviceAuthSession(context.TODO(), signature, userCodeSignature, authreq))
					},
					check: func(t *testing.T, aresp *fosite.AccessResponse) {
						assert.NotEmpty(t, aresp.AccessToken)
						assert.Equal(t, "bearer", aresp.TokenType)
						assert.Empty(t, aresp.GetExtra("refresh_token"))
						assert.NotEmpty(t, aresp.GetExtra("expires_in"))
						assert.Equal(t, "foo", aresp.GetExtra("scope"))
					},
				},
			}

			for _, testCase := range testCases {
				t.Run("case="+testCase.description, func(t *testing.T) {
					config := &fosite.Config{
						ScopeStrategy:            fosite.HierarchicScopeStrategy,
						AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
						AccessTokenLifespan:      time.Minute,
						RefreshTokenScopes:       []string{"offline"},
					}
					h := DeviceCodeTokenEndpointHandler{
						DeviceRateLimitStrategy: strategy,
						DeviceCodeStrategy:      strategy,
						UserCodeStrategy:        strategy,
						AccessTokenStrategy:     strategy.CoreStrategy,
						RefreshTokenStrategy:    strategy.CoreStrategy,
						Config:                  config,
						CoreStorage:             store,
						TokenRevocationStorage:  store,
					}

					if testCase.setup != nil {
						testCase.setup(t, testCase.areq, testCase.authreq, config)
					}

					aresp := fosite.NewAccessResponse()
					err := h.PopulateTokenEndpointResponse(context.TODO(), testCase.areq, aresp)

					if testCase.expectErr != nil {
						require.EqualError(t, err, testCase.expectErr.Error(), "%+v", err)
					} else {
						require.NoError(t, err, "%+v", err)
					}

					if testCase.check != nil {
						testCase.check(t, aresp)
					}
				})
			}
		})
	}
}

func TestDeviceUserCodeTransactional_HandleTokenEndpointRequest(t *testing.T) {
	var mockTransactional *internal.MockTransactional
	var mockCoreStore *internal.MockRFC8628CoreStorage
	var mockDeviceRateLimitStrategy *internal.MockDeviceRateLimitStrategy
	strategy := hmacshaStrategy
	deviceStrategy := RFC8628HMACSHAStrategy

	authreq := &fosite.DeviceRequest{
		UserCodeState: fosite.UserCodeAccepted,
		Request: fosite.Request{
			Client:         &fosite.DefaultClient{ID: "foo", GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code"}},
			RequestedScope: fosite.Arguments{"foo", "offline"},
			GrantedScope:   fosite.Arguments{"foo", "offline"},
			Session:        &fosite.DefaultSession{},
			RequestedAt:    time.Now().UTC(),
		},
	}

	areq := &fosite.AccessRequest{
		GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
		Request: fosite.Request{
			Client: &fosite.DefaultClient{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code", "refresh_token"},
			},
			Session:     &fosite.DefaultSession{},
			RequestedAt: time.Now().UTC(),
		},
	}
	aresp := fosite.NewAccessResponse()
	propagatedContext := context.Background()

	code, _, err := deviceStrategy.GenerateDeviceCode(context.Background())
	require.NoError(t, err)
	areq.Form = url.Values{"device_code": {code}}

	// some storage implementation that has support for transactions, notice the embedded type `storage.Transactional`

	type deviceTransactionalStore struct {
		storage.Transactional
		RFC8628CoreStorage
	}

	testCases := []struct {
		description string
		setup       func()
		expectError error
	}{
		{
			description: "transaction should be committed successfully if no errors occur",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateDeviceCodeSession(gomock.Any(), gomock.Any()).
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
			description: "transaction should be rolled back if `InvalidateDeviceCodeSession` returns an error",
			setup: func() {
				mockCoreStore.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateDeviceCodeSession(gomock.Any(), gomock.Any()).
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
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateDeviceCodeSession(gomock.Any(), gomock.Any()).
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
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
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
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateDeviceCodeSession(gomock.Any(), gomock.Any()).
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
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil)
				mockCoreStore.
					EXPECT().
					InvalidateDeviceCodeSession(gomock.Any(), gomock.Any()).
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
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("scenario=%s", testCase.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransactional = internal.NewMockTransactional(ctrl)
			mockCoreStore = internal.NewMockRFC8628CoreStorage(ctrl)
			mockDeviceRateLimitStrategy = internal.NewMockDeviceRateLimitStrategy(ctrl)
			testCase.setup()

			h := DeviceCodeTokenEndpointHandler{
				DeviceCodeStrategy:      &deviceStrategy,
				UserCodeStrategy:        &deviceStrategy,
				DeviceRateLimitStrategy: mockDeviceRateLimitStrategy,
				CoreStorage: deviceTransactionalStore{
					mockTransactional,
					mockCoreStore,
				},
				AccessTokenStrategy:  strategy,
				RefreshTokenStrategy: strategy,
				Config: &fosite.Config{
					ScopeStrategy:             fosite.HierarchicScopeStrategy,
					AudienceMatchingStrategy:  fosite.DefaultAudienceMatchingStrategy,
					DeviceAndUserCodeLifespan: time.Minute,
				},
			}

			if err = h.PopulateTokenEndpointResponse(propagatedContext, areq, aresp); testCase.expectError != nil {
				assert.EqualError(t, err, testCase.expectError.Error())
			}
		})
	}
}
