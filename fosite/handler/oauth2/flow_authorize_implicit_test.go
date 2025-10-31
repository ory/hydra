// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestAuthorizeImplicit_EndpointHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	areq := fosite.NewAuthorizeRequest()
	areq.Session = new(fosite.DefaultSession)
	h, store, chgen, aresp := makeAuthorizeImplicitGrantTypeHandler(ctrl)

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
	}{
		{
			description: "should pass because not responsible for handling the response type",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"a"}
			},
		},
		{
			description: "should fail because access token generation failed",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token"},
				}
				chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return("", "", errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because scope invalid",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token"}
				areq.RequestedScope = fosite.Arguments{"scope"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token"},
				}
			},
			expectErr: fosite.ErrInvalidScope,
		},
		{
			description: "should fail because audience invalid",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token"}
				areq.RequestedScope = fosite.Arguments{"scope"}
				areq.RequestedAudience = fosite.Arguments{"https://www.ory.sh/not-api"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token"},
					Scopes:        []string{"scope"},
					Audience:      []string{"https://www.ory.sh/api"},
				}
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because persistence failed",
			setup: func() {
				areq.RequestedAudience = fosite.Arguments{"https://www.ory.sh/api"}
				chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).AnyTimes().Return("access.ats", "ats", nil)
				store.EXPECT().CreateAccessTokenSession(gomock.Any(), "ats", gomock.Eq(areq.Sanitize([]string{}))).Return(errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			setup: func() {
				areq.State = "state"
				areq.GrantedScope = fosite.Arguments{"scope"}

				store.EXPECT().CreateAccessTokenSession(gomock.Any(), "ats", gomock.Eq(areq.Sanitize([]string{}))).AnyTimes().Return(nil)

				aresp.EXPECT().AddParameter("access_token", "access.ats")
				aresp.EXPECT().AddParameter("expires_in", gomock.Any())
				aresp.EXPECT().AddParameter("token_type", "bearer")
				aresp.EXPECT().AddParameter("state", "state")
				aresp.EXPECT().AddParameter("scope", "scope")
			},
			expectErr: nil,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			err := h.HandleAuthorizeEndpointRequest(context.Background(), areq, aresp)
			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func makeAuthorizeImplicitGrantTypeHandler(ctrl *gomock.Controller) (AuthorizeImplicitGrantTypeHandler,
	*internal.MockAccessTokenStorage, *internal.MockAccessTokenStrategy, *internal.MockAuthorizeResponder) {
	store := internal.NewMockAccessTokenStorage(ctrl)
	chgen := internal.NewMockAccessTokenStrategy(ctrl)
	aresp := internal.NewMockAuthorizeResponder(ctrl)

	h := AuthorizeImplicitGrantTypeHandler{
		AccessTokenStorage:  store,
		AccessTokenStrategy: chgen,
		Config: &fosite.Config{
			AccessTokenLifespan:      time.Hour,
			ScopeStrategy:            fosite.HierarchicScopeStrategy,
			AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
		},
	}

	return h, store, chgen, aresp
}

func TestDefaultResponseMode_AuthorizeImplicit_EndpointHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	areq := fosite.NewAuthorizeRequest()
	areq.Session = new(fosite.DefaultSession)
	h, store, chgen, aresp := makeAuthorizeImplicitGrantTypeHandler(ctrl)

	areq.State = "state"
	areq.GrantedScope = fosite.Arguments{"scope"}
	areq.ResponseTypes = fosite.Arguments{"token"}
	areq.Client = &fosite.DefaultClientWithCustomTokenLifespans{
		DefaultClient: &fosite.DefaultClient{
			GrantTypes:    fosite.Arguments{"implicit"},
			ResponseTypes: fosite.Arguments{"token"},
		},
		TokenLifespans: &internal.TestLifespans,
	}

	store.EXPECT().CreateAccessTokenSession(gomock.Any(), "ats", gomock.Eq(areq.Sanitize([]string{}))).AnyTimes().Return(nil)

	aresp.EXPECT().AddParameter("access_token", "access.ats")
	aresp.EXPECT().AddParameter("expires_in", gomock.Any())
	aresp.EXPECT().AddParameter("token_type", "bearer")
	aresp.EXPECT().AddParameter("state", "state")
	aresp.EXPECT().AddParameter("scope", "scope")
	chgen.EXPECT().GenerateAccessToken(gomock.Any(), areq).AnyTimes().Return("access.ats", "ats", nil)

	err := h.HandleAuthorizeEndpointRequest(context.Background(), areq, aresp)
	assert.NoError(t, err)
	assert.Equal(t, fosite.ResponseModeFragment, areq.GetResponseMode())

	internal.RequireEqualTime(t, time.Now().UTC().Add(*internal.TestLifespans.ImplicitGrantAccessTokenLifespan), areq.Session.GetExpiresAt(fosite.AccessToken), time.Minute)
}
