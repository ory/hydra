// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/storage"
)

func TestAccessTokenFromRequestNoToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)

	assert.Equal(t, AccessTokenFromRequest(req), "", "No token should produce an empty string")
}

func TestAccessTokenFromRequestHeader(t *testing.T) {
	token := "TokenFromHeader"

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	assert.Equal(t, AccessTokenFromRequest(req), token, "Token should be obtainable from header")
}

func TestAccessTokenFromRequestQuery(t *testing.T) {
	token := "TokenFromQueryParam"

	req, _ := http.NewRequest("GET", "http://example.com/test?access_token="+token, nil)

	assert.Equal(t, AccessTokenFromRequest(req), token, "Token should be obtainable from access_token query parameter")
}

func TestIntrospect(t *testing.T) {
	ctrl := gomock.NewController(t)
	validator := internal.NewMockTokenIntrospector(ctrl)
	t.Cleanup(ctrl.Finish)

	config := new(Config)
	f := compose.ComposeAllEnabled(config, storage.NewMemoryStore(), nil).(*Fosite)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Add("Authorization", "bearer some-token")

	for k, c := range []struct {
		description string
		scopes      []string
		setup       func()
		expectErr   error
	}{
		{
			description: "should fail",
			scopes:      []string{},
			setup: func() {
			},
			expectErr: ErrRequestUnauthorized,
		},
		{
			description: "should fail",
			scopes:      []string{"foo"},
			setup: func() {
				config.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				validator.EXPECT().IntrospectToken(gomock.Any(), "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), ErrUnknownRequest)
			},
			expectErr: ErrRequestUnauthorized,
		},
		{
			description: "should fail",
			scopes:      []string{"foo"},
			setup: func() {
				validator.EXPECT().IntrospectToken(gomock.Any(), "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenUse(""), ErrInvalidClient)
			},
			expectErr: ErrInvalidClient,
		},
		{
			description: "should pass",
			setup: func() {
				validator.EXPECT().IntrospectToken(gomock.Any(), "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, _ string, _ TokenUse, accessRequest AccessRequester, _ []string) {
					accessRequest.(*AccessRequest).GrantedScope = []string{"bar"}
				}).Return(TokenUse(""), nil)
			},
		},
		{
			description: "should pass",
			scopes:      []string{"bar"},
			setup: func() {
				validator.EXPECT().IntrospectToken(gomock.Any(), "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, _ string, _ TokenType, accessRequest AccessRequester, _ []string) {
					accessRequest.(*AccessRequest).GrantedScope = []string{"bar"}
				}).Return(TokenUse(""), nil)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			_, _, err := f.IntrospectToken(context.Background(), AccessTokenFromRequest(req), AccessToken, nil, c.scopes...)
			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
