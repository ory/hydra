// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestHandleTokenEndpointRequest(t *testing.T) {
	h := &openid.ExplicitHandler{Config: &fosite.Config{}}
	areq := fosite.NewAccessRequest(nil)
	areq.Client = &fosite.DefaultClient{
		// ResponseTypes: fosite.Arguments{"id_token"},
	}
	assert.EqualError(t, h.HandleTokenEndpointRequest(context.Background(), areq), fosite.ErrUnknownRequest.Error())
}

func TestExplicit_PopulateTokenEndpointResponse(t *testing.T) {
	for k, c := range []struct {
		description string
		setup       func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest)
		expectErr   error
		check       func(t *testing.T, aresp *fosite.AccessResponse)
	}{
		{
			description: "should fail because current request has invalid grant type",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.GrantTypes = fosite.Arguments{"some_other_grant_type"}
			},
			expectErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should fail because storage lookup returns not found",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(nil, openid.ErrNoSessionFound)
			},
			expectErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should fail because storage lookup fails",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because stored request is missing openid scope",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(fosite.NewAuthorizeRequest(), nil)
			},
			expectErr: fosite.ErrMisconfiguration,
		},
		{
			description: "should fail because current request's client does not have authorization_code grant type",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.Client = &fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"some_other_grant_type"},
				}
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				storedReq := fosite.NewAuthorizeRequest()
				storedReq.GrantedScope = fosite.Arguments{"openid"}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(storedReq, nil)
			},
			expectErr: fosite.ErrUnauthorizedClient,
		},
		{
			description: "should pass with custom client lifespans",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.Client = &fosite.DefaultClientWithCustomTokenLifespans{
					DefaultClient: &fosite.DefaultClient{
						GrantTypes: fosite.Arguments{"authorization_code"},
					},
					TokenLifespans: &internal.TestLifespans,
				}
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				storedSession := &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{Subject: "peter"},
				}
				storedReq := fosite.NewAuthorizeRequest()
				storedReq.Session = storedSession
				storedReq.GrantedScope = fosite.Arguments{"openid"}
				storedReq.Form.Set("nonce", "1111111111111111")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(2)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(storedReq, nil)
				store.EXPECT().DeleteOpenIDConnectSession(gomock.Any(), "foobar").Return(nil)
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
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.AuthorizationCodeGrantIDTokenLifespan).UTC(), *idTokenExp, time.Minute)
			},
		},
		{
			description: "should pass",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.Client = &fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"authorization_code"},
				}
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				storedSession := &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{Subject: "peter"},
				}
				storedReq := fosite.NewAuthorizeRequest()
				storedReq.Session = storedSession
				storedReq.GrantedScope = fosite.Arguments{"openid"}
				storedReq.Form.Set("nonce", "1111111111111111")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(2)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(storedReq, nil)
				store.EXPECT().DeleteOpenIDConnectSession(gomock.Any(), "foobar").Return(nil)
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
				internal.RequireEqualTime(t, time.Now().Add(time.Hour), *idTokenExp, time.Minute)
			},
		},
		{
			description: "should fail because stored request's session is missing subject claim",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				storedSession := &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{Subject: ""},
				}
				storedReq := fosite.NewAuthorizeRequest()
				storedReq.Session = storedSession
				storedReq.GrantedScope = fosite.Arguments{"openid"}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(storedReq, nil)
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because stored request is missing session",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				storedReq := fosite.NewAuthorizeRequest()
				storedReq.Session = nil
				storedReq.GrantScope("openid")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(storedReq, nil)
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because storage returns error when deleting openid session",
			setup: func(provider *internal.MockOpenIDConnectRequestStorageProvider, store *internal.MockOpenIDConnectRequestStorage, req *fosite.AccessRequest) {
				req.Client = &fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"authorization_code"},
				}
				req.GrantTypes = fosite.Arguments{"authorization_code"}
				req.Form.Set("code", "foobar")
				storedSession := &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{Subject: "peter"},
				}
				storedReq := fosite.NewAuthorizeRequest()
				storedReq.Session = storedSession
				storedReq.GrantedScope = fosite.Arguments{"openid"}
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(2)
				store.EXPECT().GetOpenIDConnectSession(gomock.Any(), "foobar", req).Return(storedReq, nil)
				store.EXPECT().DeleteOpenIDConnectSession(gomock.Any(), "foobar").Return(errors.New("delete openid session err"))
			},
			expectErr: fosite.ErrServerError,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, c.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := internal.NewMockOpenIDConnectRequestStorage(ctrl)
			provider := internal.NewMockOpenIDConnectRequestStorageProvider(ctrl)
			openIDTokenStrategyProvider := internal.NewMockOpenIDConnectTokenStrategyProvider(ctrl)
			t.Cleanup(ctrl.Finish)

			session := &openid.DefaultSession{
				Claims: &jwt.IDTokenClaims{
					Subject: "peter",
				},
				Headers: &jwt.Headers{},
			}
			aresp := fosite.NewAccessResponse()
			areq := fosite.NewAccessRequest(session)

			j := &openid.DefaultStrategy{
				Signer: &jwt.DefaultSigner{
					GetPrivateKey: func(ctx context.Context) (interface{}, error) {
						return key, nil
					},
				},
				Config: &fosite.Config{
					MinParameterEntropy: fosite.MinParameterEntropy,
				},
			}
			openIDTokenStrategyProvider.EXPECT().OpenIDConnectTokenStrategy().Return(j).AnyTimes()

			h := &openid.ExplicitHandler{
				Storage: provider,
				IDTokenHandleHelper: &openid.IDTokenHandleHelper{
					IDTokenStrategy: openIDTokenStrategyProvider,
				},
				Config: &fosite.Config{},
			}

			c.setup(provider, store, areq)
			err := h.PopulateTokenEndpointResponse(context.Background(), areq, aresp)

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
			if c.check != nil {
				c.check(t, aresp)
			}
		})
	}
}
