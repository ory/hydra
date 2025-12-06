// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

// expose key to verify id_token
var key = gen.MustRSAKey()

var oidcParameters = []string{
	"grant_type",
	"max_age",
	"prompt",
	"acr_values",
	"id_token_hint",
	"nonce",
}

func makeOpenIDConnectExplicitHandler(ctrl *gomock.Controller, minParameterEntropy int) (openid.ExplicitHandler, *internal.MockOpenIDConnectRequestStorage, *internal.MockOpenIDConnectRequestStorageProvider) {
	store := internal.NewMockOpenIDConnectRequestStorage(ctrl)
	provider := internal.NewMockOpenIDConnectRequestStorageProvider(ctrl)
	openIDTokenStrategyProvider := internal.NewMockOpenIDConnectTokenStrategyProvider(ctrl)
	config := &fosite.Config{MinParameterEntropy: minParameterEntropy}

	defaultStrategy := &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(ctx context.Context) (interface{}, error) {
				return key, nil
			},
		},
		Config: config,
	}
	openIDTokenStrategyProvider.EXPECT().OpenIDConnectTokenStrategy().Return(defaultStrategy).AnyTimes()

	return openid.ExplicitHandler{
		Storage: provider,
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: openIDTokenStrategyProvider,
		},
		OpenIDConnectRequestValidator: openid.NewOpenIDConnectRequestValidator(defaultStrategy.Signer, config),
		Config:                        config,
	}, store, provider
}

func TestExplicit_HandleAuthorizeEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	aresp := internal.NewMockAuthorizeResponder(ctrl)
	t.Cleanup(ctrl.Finish)

	areq := fosite.NewAuthorizeRequest()

	session := openid.NewDefaultSession()
	session.Claims.Subject = "foo"
	areq.Session = session
	areq.Form = url.Values{
		"redirect_uri": {"https://foobar.com"},
	}

	for k, c := range []struct {
		description string
		setup       func() openid.ExplicitHandler
		expectErr   error
	}{
		{
			description: "should pass because not responsible for handling an empty response type",
			setup: func() openid.ExplicitHandler {
				h, _, _ := makeOpenIDConnectExplicitHandler(ctrl, fosite.MinParameterEntropy)
				areq.ResponseTypes = fosite.Arguments{""}
				return h
			},
		},
		{
			description: "should pass because scope openid is not set",
			setup: func() openid.ExplicitHandler {
				h, _, _ := makeOpenIDConnectExplicitHandler(ctrl, fosite.MinParameterEntropy)
				areq.ResponseTypes = fosite.Arguments{"code"}
				areq.Client = &fosite.DefaultClient{
					ResponseTypes: fosite.Arguments{"code"},
				}
				areq.RequestedScope = fosite.Arguments{""}
				return h
			},
		},
		{
			description: "should fail because no code set",
			setup: func() openid.ExplicitHandler {
				h, _, _ := makeOpenIDConnectExplicitHandler(ctrl, fosite.MinParameterEntropy)
				areq.GrantedScope = fosite.Arguments{"openid"}
				areq.Form.Set("nonce", "11111111111111111111111111111")
				aresp.EXPECT().GetCode().Return("")
				return h
			},
			expectErr: fosite.ErrMisconfiguration,
		},
		{
			description: "should fail because lookup fails",
			setup: func() openid.ExplicitHandler {
				h, store, provider := makeOpenIDConnectExplicitHandler(ctrl, fosite.MinParameterEntropy)
				aresp.EXPECT().GetCode().AnyTimes().Return("codeexample")
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().CreateOpenIDConnectSession(gomock.Any(), "codeexample", gomock.Eq(areq.Sanitize(oidcParameters))).Return(errors.New(""))
				return h
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			setup: func() openid.ExplicitHandler {
				h, store, provider := makeOpenIDConnectExplicitHandler(ctrl, fosite.MinParameterEntropy)
				provider.EXPECT().OpenIDConnectRequestStorage().Return(store).Times(1)
				store.EXPECT().CreateOpenIDConnectSession(gomock.Any(), "codeexample", gomock.Eq(areq.Sanitize(oidcParameters))).AnyTimes().Return(nil)
				return h
			},
		},
		{
			description: "should fail because redirect url is missing",
			setup: func() openid.ExplicitHandler {
				areq.Form.Del("redirect_uri")
				h, store, _ := makeOpenIDConnectExplicitHandler(ctrl, fosite.MinParameterEntropy)
				store.EXPECT().CreateOpenIDConnectSession(gomock.Any(), "codeexample", gomock.Eq(areq.Sanitize(oidcParameters))).AnyTimes().Return(nil)
				return h
			},
			expectErr: fosite.ErrInvalidRequest,
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
		})
	}
}
