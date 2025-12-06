// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/storage"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func makeOpenIDConnectImplicitHandler(minParameterEntropy int) openid.OpenIDConnectImplicitHandler {
	config := &fosite.Config{
		MinParameterEntropy: minParameterEntropy,
		AccessTokenLifespan: time.Hour,
		ScopeStrategy:       fosite.HierarchicScopeStrategy,
	}

	defaultStrategyProvider := mockOpenIDConnectTokenStrategyProvider{
		strategy: openid.DefaultStrategy{
			Signer: &jwt.DefaultSigner{
				GetPrivateKey: func(ctx context.Context) (interface{}, error) {
					return gen.MustRSAKey(), nil
				},
			},
			Config: config,
		},
	}

	j := &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(ctx context.Context) (interface{}, error) {
				return key, nil
			},
		},
		Config: config,
	}

	return openid.OpenIDConnectImplicitHandler{
		AuthorizeImplicitGrantTypeHandler: &oauth2.AuthorizeImplicitGrantHandler{
			Config:   config,
			Strategy: hmacStrategy,
			Storage:  storage.NewMemoryStore(),
		},
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: defaultStrategyProvider,
		},
		OpenIDConnectRequestValidator: openid.NewOpenIDConnectRequestValidator(j.Signer, config),
		Config:                        config,
	}
}

func TestImplicit_HandleAuthorizeEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	aresp := fosite.NewAuthorizeResponse()
	areq := fosite.NewAuthorizeRequest()
	areq.Form = url.Values{
		"redirect_uri": {"https://foobar.com"},
	}
	areq.Session = new(fosite.DefaultSession)

	for k, c := range []struct {
		description string
		setup       func() openid.OpenIDConnectImplicitHandler
		expectErr   error
		check       func()
	}{
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.State = "foostate"
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{}
				areq.GrantedScope = fosite.Arguments{"openid"}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
				areq.RequestedScope = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{},
					ResponseTypes: fosite.Arguments{},
					Scopes:        []string{"openid", "fosite"},
				}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			expectErr: fosite.ErrInvalidGrant,
		},
		// Disabled because this is already handled at the authorize_request_handler
		//{
		//	description: "should not do anything because request requirements are not met",
		//	setup: func() OpenIDConnectImplicitHandler {
		//		areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
		//		areq.RequestedScope = fosite.Arguments{"openid"}
		//		areq.Client = &fosite.DefaultClient{
		//			GrantTypes:    fosite.Arguments{"implicit"},
		//			ResponseTypes: fosite.Arguments{},
		//			RequestedScope:        []string{"openid", "fosite"},
		//		}
		//		return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
		//	},
		//	expectErr: fosite.ErrInvalidGrant,
		//},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.RequestedScope = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"implicit"},
					// ResponseTypes: fosite.Arguments{"token", "id_token"},
					Scopes: []string{"openid", "fosite"},
				}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.Form = url.Values{
					"nonce":        {"short"},
					"redirect_uri": {"https://foobar.com"},
				}
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.RequestedScope = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token", "id_token"},
					Scopes:        []string{"openid", "fosite"},
				}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			expectErr: fosite.ErrInsufficientEntropy,
		},
		{
			description: "should fail because session not set",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.Form = url.Values{
					"nonce":        {"long-enough"},
					"redirect_uri": {"https://foobar.com"},
				}
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.RequestedScope = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token", "id_token"},
					Scopes:        []string{"openid", "fosite"},
				}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			expectErr: openid.ErrInvalidSession,
		},
		{
			description: "should pass because nonce set",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.Session = &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{
						Subject: "peter",
					},
					Headers: &jwt.Headers{},
					Subject: "peter",
				}
				areq.Form.Add("nonce", "some-random-foo-nonce-wow")
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should pass",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("state"))
				assert.Empty(t, aresp.GetParameters().Get("access_token"))

				idToken := aresp.GetParameters().Get("id_token")
				assert.NotEmpty(t, idToken)
				idTokenExp := internal.ExtractJwtExpClaim(t, idToken)
				internal.RequireEqualTime(t, time.Now().Add(time.Hour), *idTokenExp, time.Minute)
			},
		},
		{
			description: "should pass with nondefault id token lifespan",
			setup: func() openid.OpenIDConnectImplicitHandler {
				aresp = fosite.NewAuthorizeResponse()
				areq.Session = &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{
						Subject: "peter",
					},
					Headers: &jwt.Headers{},
					Subject: "peter",
				}
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.Client = &fosite.DefaultClientWithCustomTokenLifespans{
					DefaultClient: &fosite.DefaultClient{
						GrantTypes:    fosite.Arguments{"implicit"},
						ResponseTypes: fosite.Arguments{"token", "id_token"},
						Scopes:        []string{"openid", "fosite"},
					},
				}
				areq.Client.(*fosite.DefaultClientWithCustomTokenLifespans).SetTokenLifespans(&internal.TestLifespans)
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				idToken := aresp.GetParameters().Get("id_token")
				assert.NotEmpty(t, idToken)
				assert.NotEmpty(t, aresp.GetParameters().Get("state"))
				assert.Empty(t, aresp.GetParameters().Get("access_token"))
				idTokenExp := internal.ExtractJwtExpClaim(t, idToken)
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.ImplicitGrantIDTokenLifespan), *idTokenExp, time.Minute)
			},
		},
		{
			description: "should pass",
			setup: func() openid.OpenIDConnectImplicitHandler {
				aresp = fosite.NewAuthorizeResponse()
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("state"))

				idToken := aresp.GetParameters().Get("id_token")
				assert.NotEmpty(t, idToken)
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.ImplicitGrantIDTokenLifespan).UTC(), *internal.ExtractJwtExpClaim(t, idToken), time.Minute)

				assert.NotEmpty(t, aresp.GetParameters().Get("access_token"))
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.ImplicitGrantAccessTokenLifespan).UTC(), areq.Session.GetExpiresAt(fosite.AccessToken), time.Minute)
			},
		},
		{
			description: "should pass",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.ResponseTypes = fosite.Arguments{"id_token", "token"}
				areq.RequestedScope = fosite.Arguments{"fosite", "openid"}
				return makeOpenIDConnectImplicitHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("id_token"))
				assert.NotEmpty(t, aresp.GetParameters().Get("state"))
				assert.NotEmpty(t, aresp.GetParameters().Get("access_token"))
				assert.Equal(t, fosite.ResponseModeFragment, areq.GetResponseMode())
			},
		},
		{
			description: "should pass with low min entropy",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.Form.Set("nonce", "short")
				return makeOpenIDConnectImplicitHandler(4)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("id_token"))
				assert.NotEmpty(t, aresp.GetParameters().Get("state"))
				assert.NotEmpty(t, aresp.GetParameters().Get("access_token"))
			},
		},
		{
			description: "should fail without redirect_uri",
			setup: func() openid.OpenIDConnectImplicitHandler {
				areq.Form.Del("redirect_uri")
				return makeOpenIDConnectImplicitHandler(4)
			},
			expectErr: fosite.ErrInvalidRequest,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := c.setup()
			err := h.HandleAuthorizeEndpointRequest(context.Background(), areq, aresp)

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				assert.NoError(t, err)
				if c.check != nil {
					c.check()
				}
			}
		})
	}
}
