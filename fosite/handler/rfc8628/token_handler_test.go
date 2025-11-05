// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628_test

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite/internal"

	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/token/hmac"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/storage"
)

var hmacshaStrategyOAuth = oauth2.NewHMACSHAStrategy(
	&hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	&fosite.Config{
		AccessTokenLifespan:   time.Hour * 24,
		AuthorizeCodeLifespan: time.Hour * 24,
	},
)

var RFC8628HMACSHAStrategy = rfc8628.DefaultDeviceStrategy{
	Enigma: &hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	Config: &fosite.Config{
		DeviceAndUserCodeLifespan: time.Minute * 30,
	},
}

type mockDeviceCodeStrategyProvider struct {
	deviceRateLimitStrategy rfc8628.DeviceRateLimitStrategy
	deviceCodeStrategy      rfc8628.DeviceCodeStrategy
	userCodeStrategy        rfc8628.UserCodeStrategy
	coreStrategy            oauth2.CoreStrategy
}

func (t *mockDeviceCodeStrategyProvider) DeviceRateLimitStrategy() rfc8628.DeviceRateLimitStrategy {
	return t.deviceRateLimitStrategy
}

func (t *mockDeviceCodeStrategyProvider) DeviceCodeStrategy() rfc8628.DeviceCodeStrategy {
	return t.deviceCodeStrategy
}

func (t *mockDeviceCodeStrategyProvider) UserCodeStrategy() rfc8628.UserCodeStrategy {
	return t.userCodeStrategy
}

func (t *mockDeviceCodeStrategyProvider) AccessTokenStrategy() oauth2.AccessTokenStrategy {
	return t.coreStrategy.AccessTokenStrategy()
}

func (t *mockDeviceCodeStrategyProvider) RefreshTokenStrategy() oauth2.RefreshTokenStrategy {
	return t.coreStrategy.RefreshTokenStrategy()
}

func TestDeviceUserCode_HandleTokenEndpointRequest(t *testing.T) {
	for k, strategy := range map[string]struct {
		oauth2.CoreStrategy
		rfc8628.DefaultDeviceStrategy
	}{
		"hmac": {hmacshaStrategyOAuth, RFC8628HMACSHAStrategy},
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			h := rfc8628.DeviceCodeTokenEndpointHandler{
				Strategy: &mockDeviceCodeStrategyProvider{
					deviceRateLimitStrategy: &strategy.DefaultDeviceStrategy,
					deviceCodeStrategy:      &strategy.DefaultDeviceStrategy,
					userCodeStrategy:        &strategy.DefaultDeviceStrategy,
					coreStrategy:            strategy.CoreStrategy,
				},
				Storage: store,
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
		rfc8628.DefaultDeviceStrategy
	}{
		"hmac": {hmacshaStrategyOAuth, RFC8628HMACSHAStrategy},
	} {
		t.Run("strategy="+k, func(t *testing.T) {
			store := storage.NewMemoryStore()

			h := rfc8628.DeviceCodeTokenEndpointHandler{
				Strategy: &mockDeviceCodeStrategyProvider{
					deviceRateLimitStrategy: &strategy.DefaultDeviceStrategy,
					deviceCodeStrategy:      &strategy.DefaultDeviceStrategy,
					userCodeStrategy:        &strategy.DefaultDeviceStrategy,
					coreStrategy:            strategy.CoreStrategy,
				},
				Storage: store,
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
		rfc8628.DefaultDeviceStrategy
	}{
		"hmac": {hmacshaStrategyOAuth, RFC8628HMACSHAStrategy},
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

					h := rfc8628.DeviceCodeTokenEndpointHandler{
						Strategy: &mockDeviceCodeStrategyProvider{
							deviceRateLimitStrategy: &strategy.DefaultDeviceStrategy,
							deviceCodeStrategy:      &strategy.DefaultDeviceStrategy,
							userCodeStrategy:        &strategy.DefaultDeviceStrategy,
							coreStrategy:            strategy.CoreStrategy,
						},
						Storage: store,
						Config:  config,
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

	var mockDeviceAuthStorage *internal.MockDeviceAuthStorage
	var mockDeviceAuthStorageProvider *internal.MockDeviceAuthStorageProvider
	var mockAccessTokenStorage *internal.MockAccessTokenStorage
	var mockAccessTokenStorageProvider *internal.MockAccessTokenStorageProvider
	var mockRefreshTokenStorage *internal.MockRefreshTokenStorage
	var mockRefreshTokenStorageProvider *internal.MockRefreshTokenStorageProvider
	var mockTokenRevocationStorageProvider *internal.MockTokenRevocationStorageProvider

	var mockDeviceRateLimitStrategyProvider *internal.MockDeviceRateLimitStrategyProvider
	var mockDeviceCodeStrategy *internal.MockDeviceCodeStrategy
	var mockDeviceCodeStrategyProvider *internal.MockDeviceCodeStrategyProvider
	var mockUserCodeStrategyProvider *internal.MockUserCodeStrategyProvider
	var mockAccessTokenStrategy *internal.MockAccessTokenStrategy
	var mockAccessTokenStrategyProvider *internal.MockAccessTokenStrategyProvider
	var mockRefreshTokenStrategy *internal.MockRefreshTokenStrategy
	var mockRefreshTokenStrategyProvider *internal.MockRefreshTokenStrategyProvider

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

	testCases := []struct {
		description string
		setup       func()
		expectError error
	}{
		{
			description: "transaction should be committed successfully if no errors occur",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy).Times(2)
				mockDeviceAuthStorageProvider.EXPECT().DeviceAuthStorage().Return(mockDeviceAuthStorage).Times(2)
				mockAccessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockAccessTokenStrategy).Times(1)
				mockRefreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(mockRefreshTokenStrategy).Times(1)
				mockAccessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(mockAccessTokenStorage).Times(1)
				mockRefreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(mockRefreshTokenStorage).Times(1)

				mockDeviceCodeStrategy.
					EXPECT().
					DeviceCodeSignature(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), nil)
				mockDeviceAuthStorage.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockDeviceCodeStrategy.
					EXPECT().
					ValidateDeviceCode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockAccessTokenStrategy.
					EXPECT().
					GenerateAccessToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockRefreshTokenStrategy.
					EXPECT().
					GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockDeviceAuthStorage.
					EXPECT().
					InvalidateDeviceCodeSession(propagatedContext, gomock.Any()).
					Return(nil).
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
			description: "transaction should be rolled back if `InvalidateDeviceCodeSession` returns an error",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy).Times(2)
				mockDeviceAuthStorageProvider.EXPECT().DeviceAuthStorage().Return(mockDeviceAuthStorage).Times(2)
				mockAccessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockAccessTokenStrategy).Times(1)
				mockRefreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(mockRefreshTokenStrategy).Times(1)

				mockDeviceCodeStrategy.
					EXPECT().
					DeviceCodeSignature(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), nil)
				mockDeviceAuthStorage.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockDeviceCodeStrategy.
					EXPECT().
					ValidateDeviceCode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockAccessTokenStrategy.
					EXPECT().
					GenerateAccessToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockRefreshTokenStrategy.
					EXPECT().
					GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockDeviceAuthStorage.
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
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy).Times(2)
				mockDeviceAuthStorageProvider.EXPECT().DeviceAuthStorage().Return(mockDeviceAuthStorage).Times(2)
				mockAccessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockAccessTokenStrategy).Times(1)
				mockRefreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(mockRefreshTokenStrategy).Times(1)
				mockAccessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(mockAccessTokenStorage).Times(1)

				mockDeviceCodeStrategy.
					EXPECT().
					DeviceCodeSignature(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), nil)
				mockDeviceAuthStorage.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockDeviceCodeStrategy.
					EXPECT().
					ValidateDeviceCode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockAccessTokenStrategy.
					EXPECT().
					GenerateAccessToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockRefreshTokenStrategy.
					EXPECT().
					GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockDeviceAuthStorage.
					EXPECT().
					InvalidateDeviceCodeSession(propagatedContext, gomock.Any()).
					Return(nil).
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
			description: "should result in a server error if transaction cannot be created",
			setup: func() {
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy).Times(2)
				mockDeviceAuthStorageProvider.EXPECT().DeviceAuthStorage().Return(mockDeviceAuthStorage).Times(1)
				mockAccessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockAccessTokenStrategy).Times(1)
				mockRefreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(mockRefreshTokenStrategy).Times(1)

				mockDeviceCodeStrategy.
					EXPECT().
					DeviceCodeSignature(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), nil)
				mockDeviceAuthStorage.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockDeviceCodeStrategy.
					EXPECT().
					ValidateDeviceCode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockAccessTokenStrategy.
					EXPECT().
					GenerateAccessToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockRefreshTokenStrategy.
					EXPECT().
					GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
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
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy).Times(2)
				mockDeviceAuthStorageProvider.EXPECT().DeviceAuthStorage().Return(mockDeviceAuthStorage).Times(2)
				mockAccessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockAccessTokenStrategy).Times(1)
				mockRefreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(mockRefreshTokenStrategy).Times(1)

				mockDeviceCodeStrategy.
					EXPECT().
					DeviceCodeSignature(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), nil)
				mockDeviceAuthStorage.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockDeviceCodeStrategy.
					EXPECT().
					ValidateDeviceCode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockAccessTokenStrategy.
					EXPECT().
					GenerateAccessToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockRefreshTokenStrategy.
					EXPECT().
					GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockDeviceAuthStorage.
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
				mockDeviceCodeStrategyProvider.EXPECT().DeviceCodeStrategy().Return(mockDeviceCodeStrategy).Times(2)
				mockDeviceAuthStorageProvider.EXPECT().DeviceAuthStorage().Return(mockDeviceAuthStorage).Times(2)
				mockAccessTokenStrategyProvider.EXPECT().AccessTokenStrategy().Return(mockAccessTokenStrategy).Times(1)
				mockRefreshTokenStrategyProvider.EXPECT().RefreshTokenStrategy().Return(mockRefreshTokenStrategy).Times(1)
				mockAccessTokenStorageProvider.EXPECT().AccessTokenStorage().Return(mockAccessTokenStorage).Times(1)
				mockRefreshTokenStorageProvider.EXPECT().RefreshTokenStorage().Return(mockRefreshTokenStorage).Times(1)

				mockDeviceCodeStrategy.
					EXPECT().
					DeviceCodeSignature(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), nil)
				mockDeviceAuthStorage.
					EXPECT().
					GetDeviceCodeSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(authreq, nil).
					Times(1)
				mockDeviceCodeStrategy.
					EXPECT().
					ValidateDeviceCode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockAccessTokenStrategy.
					EXPECT().
					GenerateAccessToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockRefreshTokenStrategy.
					EXPECT().
					GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Return(gomock.Any().String(), gomock.Any().String(), nil)
				mockTransactional.
					EXPECT().
					BeginTX(propagatedContext).
					Return(propagatedContext, nil).
					Times(1)
				mockDeviceAuthStorage.
					EXPECT().
					InvalidateDeviceCodeSession(propagatedContext, gomock.Any()).
					Return(nil).
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
			t.Cleanup(ctrl.Finish)

			mockTransactional = internal.NewMockTransactional(ctrl)

			mockDeviceAuthStorage = internal.NewMockDeviceAuthStorage(ctrl)
			mockDeviceAuthStorageProvider = internal.NewMockDeviceAuthStorageProvider(ctrl)
			mockAccessTokenStorage = internal.NewMockAccessTokenStorage(ctrl)
			mockAccessTokenStorageProvider = internal.NewMockAccessTokenStorageProvider(ctrl)
			mockRefreshTokenStorage = internal.NewMockRefreshTokenStorage(ctrl)
			mockRefreshTokenStorageProvider = internal.NewMockRefreshTokenStorageProvider(ctrl)
			mockTokenRevocationStorageProvider = internal.NewMockTokenRevocationStorageProvider(ctrl)

			mockDeviceRateLimitStrategyProvider = internal.NewMockDeviceRateLimitStrategyProvider(ctrl)
			mockDeviceCodeStrategy = internal.NewMockDeviceCodeStrategy(ctrl)
			mockDeviceCodeStrategyProvider = internal.NewMockDeviceCodeStrategyProvider(ctrl)
			mockUserCodeStrategyProvider = internal.NewMockUserCodeStrategyProvider(ctrl)
			mockAccessTokenStrategy = internal.NewMockAccessTokenStrategy(ctrl)
			mockAccessTokenStrategyProvider = internal.NewMockAccessTokenStrategyProvider(ctrl)
			mockRefreshTokenStrategy = internal.NewMockRefreshTokenStrategy(ctrl)
			mockRefreshTokenStrategyProvider = internal.NewMockRefreshTokenStrategyProvider(ctrl)

			mockStorage := struct {
				*internal.MockDeviceAuthStorageProvider
				*internal.MockAccessTokenStorageProvider
				*internal.MockRefreshTokenStorageProvider
				*internal.MockTokenRevocationStorageProvider
				*internal.MockTransactional
			}{
				MockDeviceAuthStorageProvider:      mockDeviceAuthStorageProvider,
				MockAccessTokenStorageProvider:     mockAccessTokenStorageProvider,
				MockRefreshTokenStorageProvider:    mockRefreshTokenStorageProvider,
				MockTokenRevocationStorageProvider: mockTokenRevocationStorageProvider,
				MockTransactional:                  mockTransactional,
			}

			mockStrategy := struct {
				*internal.MockDeviceRateLimitStrategyProvider
				*internal.MockDeviceCodeStrategyProvider
				*internal.MockUserCodeStrategyProvider
				*internal.MockAccessTokenStrategyProvider
				*internal.MockRefreshTokenStrategyProvider
			}{
				MockDeviceRateLimitStrategyProvider: mockDeviceRateLimitStrategyProvider,
				MockDeviceCodeStrategyProvider:      mockDeviceCodeStrategyProvider,
				MockUserCodeStrategyProvider:        mockUserCodeStrategyProvider,
				MockAccessTokenStrategyProvider:     mockAccessTokenStrategyProvider,
				MockRefreshTokenStrategyProvider:    mockRefreshTokenStrategyProvider,
			}

			testCase.setup()

			h := rfc8628.DeviceCodeTokenEndpointHandler{
				Strategy: mockStrategy,
				Storage:  mockStorage,
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
