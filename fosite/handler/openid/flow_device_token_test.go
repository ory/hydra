// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/token/hmac"
	"github.com/ory/hydra/v2/fosite/token/jwt"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/fosite"
)

func TestDeviceToken_HandleTokenEndpointRequest(t *testing.T) {
	h := openid.OpenIDConnectDeviceHandler{
		Config: &fosite.Config{},
	}
	areq := fosite.NewAccessRequest(nil)
	areq.Client = &fosite.DefaultClient{
		ResponseTypes: fosite.Arguments{"code"},
	}

	err := h.HandleTokenEndpointRequest(context.Background(), areq)
	assert.EqualError(t, err, fosite.ErrUnknownRequest.Error())
}

func TestDeviceToken_PopulateTokenEndpointResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	store := internal.NewMockOpenIDConnectRequestStorage(ctrl)
	provider := internal.NewMockOpenIDConnectRequestStorageProvider(ctrl)
	strategyProvider := internal.NewMockDeviceCodeStrategyProvider(ctrl)
	openIDTokenStrategyProvider := internal.NewMockOpenIDConnectTokenStrategyProvider(ctrl)

	config := &fosite.Config{
		MinParameterEntropy:       fosite.MinParameterEntropy,
		DeviceAndUserCodeLifespan: time.Hour * 24,
		IDTokenLifespan:           time.Hour * 24,
	}
	strategy := &rfc8628.DefaultDeviceStrategy{
		Enigma: &hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobar")}},
		Config: config,
	}
	strategyProvider.EXPECT().DeviceCodeStrategy().Return(strategy).AnyTimes()

	signer := &jwt.DefaultSigner{
		GetPrivateKey: func(ctx context.Context) (interface{}, error) {
			return key, nil
		},
	}

	defaultStrategy := &openid.DefaultStrategy{
		Signer: signer,
		Config: config,
	}
	openIDTokenStrategyProvider.EXPECT().OpenIDConnectTokenStrategy().Return(defaultStrategy).AnyTimes()

	h := openid.OpenIDConnectDeviceHandler{
		Storage:  provider,
		Strategy: strategyProvider,
		Config:   config,
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: openIDTokenStrategyProvider,
		},
	}

	session := &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Subject: "foo",
		},
		Headers: &jwt.Headers{},
	}

	client := &fosite.DefaultClient{
		ID:         "foo",
		GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
	}

	testCases := []struct {
		description string
		areq        *fosite.AccessRequest
		aresp       *fosite.AccessResponse
		setup       func(areq *fosite.AccessRequest)
		check       func(t *testing.T, aresp *fosite.AccessResponse)
		expectErr   error
	}{
		{
			description: "should fail because the grant type is invalid",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"authorization_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			aresp:     fosite.NewAccessResponse(),
			expectErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should fail because session not found",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			aresp: fosite.NewAccessResponse(),
			setup: func(areq *fosite.AccessRequest) {
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), gomock.Any(), areq).Return(nil, openid.ErrNoSessionFound)
			},
			expectErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should fail because session lookup fails",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			setup: func(areq *fosite.AccessRequest) {
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), gomock.Any(), areq).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because auth request grant scope is invalid",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			setup: func(areq *fosite.AccessRequest) {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client:       client,
						GrantedScope: fosite.Arguments{"email"},
						Session:      session,
					},
				}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), gomock.Any(), areq).Return(authreq, nil)
			},
			expectErr: fosite.ErrMisconfiguration,
		},
		{
			description: "should fail because auth request is missing session",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			setup: func(areq *fosite.AccessRequest) {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client:       client,
						GrantedScope: fosite.Arguments{"openid", "email"},
					},
				}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), gomock.Any(), areq).Return(authreq, nil)
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because auth request session is missing subject claims",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			setup: func(areq *fosite.AccessRequest) {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client:       client,
						GrantedScope: fosite.Arguments{"openid", "email"},
						Session:      openid.NewDefaultSession(),
					},
				}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), gomock.Any(), areq).Return(authreq, nil)
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			areq: &fosite.AccessRequest{
				GrantTypes: fosite.Arguments{"urn:ietf:params:oauth:grant-type:device_code"},
				Request: fosite.Request{
					Client:  client,
					Form:    url.Values{"device_code": []string{"device_code"}},
					Session: session,
				},
			},
			setup: func(areq *fosite.AccessRequest) {
				authreq := &fosite.DeviceRequest{
					Request: fosite.Request{
						Client:       client,
						GrantedScope: fosite.Arguments{"openid", "email"},
						Session:      session,
					},
				}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(2)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), gomock.Any(), areq).Return(authreq, nil)
				store.EXPECT().DeleteOpenIDConnectSession(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t *testing.T, aresp *fosite.AccessResponse) {
				assert.NotEmpty(t, aresp.GetExtra("id_token"))

				idToken, _ := aresp.GetExtra("id_token").(string)
				decodedIdToken, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
					return key.PublicKey, nil
				})
				require.NoError(t, err)

				claims := decodedIdToken.Claims
				assert.NotEmpty(t, claims["at_hash"])

				idTokenExp := internal.ExtractJwtExpClaim(t, idToken)
				internal.RequireEqualTime(t, time.Now().Add(time.Hour*24), *idTokenExp, time.Minute)
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case=%d/description=%s", i, testCase.description), func(t *testing.T) {
			if testCase.setup != nil {
				testCase.setup(testCase.areq)
			}

			aresp := fosite.NewAccessResponse()
			err := h.PopulateTokenEndpointResponse(context.Background(), testCase.areq, aresp)
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
}
