// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestOpenIDConnectRefreshHandler_HandleTokenEndpointRequest(t *testing.T) {
	h := &openid.OpenIDConnectRefreshHandler{Config: &fosite.Config{}}
	for _, c := range []struct {
		areq        *fosite.AccessRequest
		expectedErr error
		description string
	}{
		{
			description: "should not pass because grant_type is wrong",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"foo"},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should not pass because grant_type is right but scope is missing",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"something"},
				},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should not pass because client may not execute this grant type",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client:       &fosite.DefaultClient{},
				},
			},
			expectedErr: fosite.ErrUnauthorizedClient,
		},
		{
			description: "should pass",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
						// ResponseTypes: []string{"id_token"},
					},
					Session: &openid.DefaultSession{},
				},
			},
		},
	} {
		t.Run("case="+c.description, func(t *testing.T) {
			err := h.HandleTokenEndpointRequest(context.Background(), c.areq)
			if c.expectedErr != nil {
				require.EqualError(t, err, c.expectedErr.Error(), "%v", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOpenIDConnectRefreshHandler_PopulateTokenEndpointResponse(t *testing.T) {
	defaultStrategyProvider := mockOpenIDConnectTokenStrategyProvider{
		strategy: openid.DefaultStrategy{
			Signer: &jwt.DefaultSigner{
				GetPrivateKey: func(ctx context.Context) (interface{}, error) {
					return key, nil
				},
			},
			Config: &fosite.Config{
				MinParameterEntropy: fosite.MinParameterEntropy,
			},
		},
	}

	h := &openid.OpenIDConnectRefreshHandler{
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: defaultStrategyProvider,
		},
		Config: &fosite.Config{},
	}
	for _, c := range []struct {
		areq        *fosite.AccessRequest
		expectedErr error
		check       func(t *testing.T, aresp *fosite.AccessResponse)
		description string
	}{
		{
			description: "should not pass because grant_type is wrong",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"foo"},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should not pass because grant_type is right but scope is missing",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"something"},
				},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		// Disabled because this is already handled at the authorize_request_handler
		//{
		//	description: "should not pass because client may not ask for id_token",
		//	areq: &fosite.AccessRequest{
		//		GrantTypes: []string{"refresh_token"},
		//		Request: fosite.Request{
		//			GrantedScope: []string{"openid"},
		//			Client: &fosite.DefaultClient{
		//				GrantTypes: []string{"refresh_token"},
		//			},
		//		},
		//	},
		//	expectedErr: fosite.ErrUnknownRequest,
		//},
		{
			description: "should pass",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
						// ResponseTypes: []string{"id_token"},
					},
					Session: &openid.DefaultSession{
						Subject: "foo",
						Claims: &jwt.IDTokenClaims{
							Subject: "foo",
						},
					},
				},
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
				require.NotEmpty(t, idTokenExp)
				internal.RequireEqualTime(t, time.Now().Add(time.Hour).UTC(), *idTokenExp, time.Minute)
			},
		},
		{
			description: "should pass",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClientWithCustomTokenLifespans{
						DefaultClient: &fosite.DefaultClient{
							GrantTypes: []string{"refresh_token"},
							// ResponseTypes: []string{"id_token"},
						},
						TokenLifespans: &internal.TestLifespans,
					},
					Session: &openid.DefaultSession{
						Subject: "foo",
						Claims: &jwt.IDTokenClaims{
							Subject: "foo",
						},
					},
				},
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
				require.NotEmpty(t, idTokenExp)
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.RefreshTokenGrantIDTokenLifespan).UTC(), *idTokenExp, time.Minute)
			},
		},
		{
			description: "should fail because missing subject claim",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
						// ResponseTypes: []string{"id_token"},
					},
					Session: &openid.DefaultSession{
						Subject: "foo",
						Claims:  &jwt.IDTokenClaims{},
					},
				},
			},
			expectedErr: fosite.ErrServerError,
		},
		{
			description: "should fail because missing session",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
					},
				},
			},
			expectedErr: fosite.ErrServerError,
		},
	} {
		t.Run("case="+c.description, func(t *testing.T) {
			aresp := fosite.NewAccessResponse()
			err := h.PopulateTokenEndpointResponse(context.Background(), c.areq, aresp)
			if c.expectedErr != nil {
				require.EqualError(t, err, c.expectedErr.Error(), "%v", err)
			} else {
				require.NoError(t, err)
			}

			if c.check != nil {
				c.check(t, aresp)
			}
		})
	}
}
