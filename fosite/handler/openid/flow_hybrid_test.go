// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	cristaljwt "github.com/cristalhq/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/storage"
	"github.com/ory/hydra/v2/fosite/token/hmac"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

var hmacStrategy = oauth2.NewHMACSHAStrategy(
	&hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("some-super-cool-secret-that-nobody-knows-nobody-knows")}},
	nil,
)

type mockOpenIDConnectTokenStrategyProvider struct {
	strategy openid.DefaultStrategy
}

func (p mockOpenIDConnectTokenStrategyProvider) OpenIDConnectTokenStrategy() openid.OpenIDConnectTokenStrategy {
	return p.strategy
}

func makeOpenIDConnectHybridHandler(minParameterEntropy int) openid.OpenIDConnectHybridHandler {
	defaultStrategyProvider := mockOpenIDConnectTokenStrategyProvider{
		strategy: openid.DefaultStrategy{
			Signer: &jwt.DefaultSigner{
				GetPrivateKey: func(_ context.Context) (interface{}, error) {
					return gen.MustRSAKey(), nil
				},
			},
			Config: &fosite.Config{
				MinParameterEntropy: minParameterEntropy,
			},
		},
	}

	j := &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return key, nil
			},
		},
		Config: &fosite.Config{
			MinParameterEntropy: minParameterEntropy,
		},
	}

	config := &fosite.Config{
		ScopeStrategy:         fosite.HierarchicScopeStrategy,
		MinParameterEntropy:   minParameterEntropy,
		AccessTokenLifespan:   time.Hour,
		AuthorizeCodeLifespan: time.Hour,
		RefreshTokenLifespan:  time.Hour,
	}
	return openid.OpenIDConnectHybridHandler{
		AuthorizeExplicitGrantHandler: &oauth2.AuthorizeExplicitGrantHandler{
			Strategy: hmacStrategy,
			Storage:  storage.NewMemoryStore(),
			Config:   config,
		},
		AuthorizeImplicitGrantHandler: &oauth2.AuthorizeImplicitGrantHandler{
			Config: &fosite.Config{
				AccessTokenLifespan: time.Hour,
			},
			Strategy: hmacStrategy,
			Storage:  storage.NewMemoryStore(),
		},
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: defaultStrategyProvider,
		},
		Config:                        config,
		OpenIDConnectRequestValidator: openid.NewOpenIDConnectRequestValidator(j.Signer, config),
		OpenIDConnectRequestStorage:   storage.NewMemoryStore(),
	}
}

func TestHybrid_HandleAuthorizeEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	aresp := fosite.NewAuthorizeResponse()
	areq := fosite.NewAuthorizeRequest()
	areq.Form = url.Values{"redirect_uri": {"https://foobar.com"}}

	for k, c := range []struct {
		description string
		setup       func() openid.OpenIDConnectHybridHandler
		check       func()
		expectErr   error
	}{
		{
			description: "should not do anything because not a hybrid request",
			setup: func() openid.OpenIDConnectHybridHandler {
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should not do anything because not a hybrid request",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should fail because nonce set but too short",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form = url.Values{
					"redirect_uri": {"https://foobar.com"},
					"nonce":        {"short"},
				}
				areq.ResponseTypes = fosite.Arguments{"token", "code"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				areq.GrantedScope = fosite.Arguments{"openid"}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			expectErr: fosite.ErrInsufficientEntropy,
		},
		{
			description: "should fail because nonce set but too short for non-default min entropy",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form = url.Values{
					"nonce":        {"some-foobar-nonce-win"},
					"redirect_uri": {"https://foobar.com"},
				}
				areq.ResponseTypes = fosite.Arguments{"token", "code"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				areq.GrantedScope = fosite.Arguments{"openid"}
				return makeOpenIDConnectHybridHandler(42)
			},
			expectErr: fosite.ErrInsufficientEntropy,
		},
		{
			description: "should fail because session not given",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form = url.Values{
					"nonce":        {"long-enough"},
					"redirect_uri": {"https://foobar.com"},
				}
				areq.ResponseTypes = fosite.Arguments{"token", "code"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				areq.GrantedScope = fosite.Arguments{"openid"}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			expectErr: openid.ErrInvalidSession,
		},
		{
			description: "should fail because client missing response types",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.ResponseTypes = fosite.Arguments{"token", "code", "id_token"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				areq.Session = &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{
						Subject: "peter",
					},
					Headers: &jwt.Headers{},
					Subject: "peter",
				}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			expectErr: fosite.ErrInvalidGrant,
		},
		{
			description: "should pass with exact one state parameter in response",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form = url.Values{
					"redirect_uri": {"https://foobar.com"},
					"nonce":        {"long-enough"},
					"state":        {""},
				}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				params := aresp.GetParameters()
				var stateParam []string
				for k, v := range params {
					if k == "state" {
						stateParam = v
						break
					}
				}
				assert.Len(t, stateParam, 1)
			},
		},
		{
			description: "should pass because nonce was set with sufficient entropy",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form.Set("nonce", "some-foobar-nonce-win")
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should pass even if nonce was not set",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
		},
		{
			description: "should pass because nonce was set with low entropy but also with low min entropy",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form.Set("nonce", "short")
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
					ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
					Scopes:        []string{"openid"},
				}
				return makeOpenIDConnectHybridHandler(4)
			},
		},
		{
			description: "should pass because AuthorizeCode's ExpiresAt is set, even if AuthorizeCodeLifespan is zero",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form.Set("nonce", "some-foobar-nonce-win")
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.True(t, !areq.Session.GetExpiresAt(fosite.AuthorizeCode).IsZero())
			},
		},
		{
			description: "should pass",
			setup: func() openid.OpenIDConnectHybridHandler {
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("id_token"))
				assert.NotEmpty(t, aresp.GetParameters().Get("code"))
				assert.NotEmpty(t, aresp.GetParameters().Get("access_token"))
				internal.RequireEqualTime(t, time.Now().Add(time.Hour).UTC(), areq.GetSession().GetExpiresAt(fosite.AuthorizeCode), time.Second)
			},
		},
		{
			description: "should fail if redirect_uri is missing",
			setup: func() openid.OpenIDConnectHybridHandler {
				areq.Form.Del("redirect_uri")
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should pass with custom client lifespans",
			setup: func() openid.OpenIDConnectHybridHandler {
				aresp = fosite.NewAuthorizeResponse()
				areq = fosite.NewAuthorizeRequest()
				areq.Form.Set("nonce", "some-foobar-nonce-win")
				areq.Form.Set("redirect_uri", "https://foobar.com")
				areq.ResponseTypes = fosite.Arguments{"token", "code", "id_token"}
				areq.Client = &fosite.DefaultClientWithCustomTokenLifespans{
					DefaultClient: &fosite.DefaultClient{
						GrantTypes:    fosite.Arguments{"authorization_code", "implicit"},
						ResponseTypes: fosite.Arguments{"token", "code", "id_token"},
						Scopes:        []string{"openid"},
					},
				}
				areq.GrantedScope = fosite.Arguments{"openid"}
				areq.Session = &openid.DefaultSession{
					Claims: &jwt.IDTokenClaims{
						Subject: "peter",
					},
					Headers: &jwt.Headers{},
					Subject: "peter",
				}
				areq.GetClient().(*fosite.DefaultClientWithCustomTokenLifespans).SetTokenLifespans(&internal.TestLifespans)
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("code"))
				internal.RequireEqualTime(t, time.Now().Add(1*time.Hour).UTC(), areq.GetSession().GetExpiresAt(fosite.AuthorizeCode), time.Second)

				idToken := aresp.GetParameters().Get("id_token")
				assert.NotEmpty(t, idToken)
				assert.True(t, areq.GetSession().GetExpiresAt(fosite.IDToken).IsZero())
				jwt, err := cristaljwt.ParseNoVerify([]byte(idToken))
				require.NoError(t, err)
				claims := &cristaljwt.RegisteredClaims{}
				require.NoError(t, json.Unmarshal(jwt.Claims(), claims))
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.ImplicitGrantIDTokenLifespan), claims.ExpiresAt.Time, time.Minute)

				assert.NotEmpty(t, aresp.GetParameters().Get("access_token"))
				internal.RequireEqualTime(t, time.Now().Add(*internal.TestLifespans.ImplicitGrantAccessTokenLifespan).UTC(), areq.GetSession().GetExpiresAt(fosite.AccessToken), time.Second)
			},
		},
		{
			description: "Default responseMode check",
			setup: func() openid.OpenIDConnectHybridHandler {
				return makeOpenIDConnectHybridHandler(fosite.MinParameterEntropy)
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetParameters().Get("id_token"))
				assert.NotEmpty(t, aresp.GetParameters().Get("code"))
				assert.NotEmpty(t, aresp.GetParameters().Get("access_token"))
				assert.Equal(t, fosite.ResponseModeFragment, areq.GetResponseMode())
				assert.WithinDuration(t, time.Now().Add(time.Hour).UTC(), areq.GetSession().GetExpiresAt(fosite.AuthorizeCode), 5*time.Second)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			h := c.setup()
			err := h.HandleAuthorizeEndpointRequest(context.Background(), areq, aresp)

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}

			if c.check != nil {
				c.check()
			}
		})
	}
}
