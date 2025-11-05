// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	. "github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/storage"
)

func TestIntrospectionResponseTokenUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	validator := internal.NewMockTokenIntrospector(ctrl)
	t.Cleanup(ctrl.Finish)

	ctx := gomock.AssignableToTypeOf(context.WithValue(context.TODO(), ContextKey("test"), nil))

	config := new(Config)
	f := compose.ComposeAllEnabled(config, storage.NewExampleStore(), nil).(*Fosite)
	httpreq := &http.Request{
		Method: "POST",
		Header: http.Header{
			"Authorization": []string{"bearer some-token"},
		},
		PostForm: url.Values{
			"token": []string{"introspect-token"},
		},
	}
	for k, c := range []struct {
		description string
		setup       func()
		expectedTU  TokenUse
		expectedATT string
	}{
		{
			description: "introspecting access token",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				validator.EXPECT().IntrospectToken(ctx, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(AccessToken, nil)
			},
			expectedATT: BearerAccessToken,
			expectedTU:  AccessToken,
		},
		{
			description: "introspecting refresh token",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				validator.EXPECT().IntrospectToken(ctx, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(RefreshToken, nil)
			},
			expectedATT: "",
			expectedTU:  RefreshToken,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			res, err := f.NewIntrospectionRequest(context.TODO(), httpreq, &DefaultSession{})
			require.NoError(t, err)
			assert.Equal(t, c.expectedATT, res.GetAccessTokenType())
			assert.Equal(t, c.expectedTU, res.GetTokenUse())
		})
	}
}

func TestIntrospectionResponse(t *testing.T) {
	r := &fosite.IntrospectionResponse{
		AccessRequester: fosite.NewAccessRequest(nil),
		Active:          true,
	}

	assert.Equal(t, r.AccessRequester, r.GetAccessRequester())
	assert.Equal(t, r.Active, r.IsActive())
}

func TestNewIntrospectionRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	validator := internal.NewMockTokenIntrospector(ctrl)
	t.Cleanup(ctrl.Finish)

	ctx := gomock.AssignableToTypeOf(context.WithValue(context.TODO(), ContextKey("test"), nil))

	config := new(Config)
	f := compose.ComposeAllEnabled(config, storage.NewExampleStore(), nil).(*Fosite)
	httpreq := &http.Request{
		Method: "POST",
		Header: http.Header{},
		Form:   url.Values{},
	}
	newErr := errors.New("asdf")

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		isActive    bool
	}{
		{
			description: "should fail",
			setup: func() {
			},
			expectErr: ErrInvalidRequest,
		},
		{
			description: "should fail",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						"Authorization": []string{"bearer some-token"},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(ctx, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), newErr)
			},
			isActive:  false,
			expectErr: ErrInactiveToken,
		},
		{
			description: "should pass",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						"Authorization": []string{"bearer some-token"},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(ctx, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
			},
			isActive: true,
		},
		{
			description: "should pass with basic auth if username and password encoded",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						// Basic Authorization with username=encoded:client and password=encoded&password
						"Authorization": []string{"Basic ZW5jb2RlZCUzQWNsaWVudDplbmNvZGVkJTI2cGFzc3dvcmQ="},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
			},
			isActive: true,
		},
		{
			description: "should pass with basic auth if username and password not encoded",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						// Basic Authorization with username=my-client and password=foobar
						"Authorization": []string{"Basic bXktY2xpZW50OmZvb2Jhcg=="},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
			},
			isActive: true,
		},
		{
			description: "should pass with basic auth if username and password not encoded",
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						// Basic Authorization with username=my-client and password=foobaz
						"Authorization": []string{"Basic bXktY2xpZW50OmZvb2Jheg=="},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(ctx, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), nil)
			},
			isActive: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			res, err := f.NewIntrospectionRequest(context.TODO(), httpreq, &DefaultSession{})

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.isActive, res.IsActive())
			}
		})
	}
}
