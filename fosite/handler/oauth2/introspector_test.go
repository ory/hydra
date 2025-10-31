// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestIntrospectToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockCoreStorage(ctrl)
	chgen := internal.NewMockCoreStrategy(ctrl)
	areq := fosite.NewAccessRequest(nil)
	defer ctrl.Finish()

	config := &fosite.Config{}
	v := &CoreValidator{
		CoreStrategy: chgen,
		CoreStorage:  store,
		Config:       config,
	}
	httpreq := &http.Request{Header: http.Header{}}

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		expectTU    fosite.TokenUse
	}{
		{
			description: "should fail because no bearer token set",
			setup: func() {
				httpreq.Header.Set("Authorization", "bearer")
				chgen.EXPECT().AccessTokenSignature(gomock.Any(), "").Return("")
				store.EXPECT().GetAccessTokenSession(gomock.Any(), "", nil).Return(nil, errors.New(""))
				chgen.EXPECT().RefreshTokenSignature(gomock.Any(), "").Return("")
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), "", nil).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrRequestUnauthorized,
		},
		{
			description: "should fail because retrieval fails",
			setup: func() {
				httpreq.Header.Set("Authorization", "bearer 1234")
				chgen.EXPECT().AccessTokenSignature(gomock.Any(), "1234").AnyTimes().Return("asdf")
				store.EXPECT().GetAccessTokenSession(gomock.Any(), "asdf", nil).Return(nil, errors.New(""))
				chgen.EXPECT().RefreshTokenSignature(gomock.Any(), "1234").Return("asdf")
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), "asdf", nil).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrRequestUnauthorized,
		},
		{
			description: "should fail because validation fails",
			setup: func() {
				store.EXPECT().GetAccessTokenSession(gomock.Any(), "asdf", nil).AnyTimes().Return(areq, nil)
				chgen.EXPECT().ValidateAccessToken(gomock.Any(), areq, "1234").Return(errorsx.WithStack(fosite.ErrTokenExpired))
				chgen.EXPECT().RefreshTokenSignature(gomock.Any(), "1234").Return("asdf")
				store.EXPECT().GetRefreshTokenSession(gomock.Any(), "asdf", nil).Return(nil, errors.New(""))
			},
			expectErr: fosite.ErrTokenExpired,
		},
		{
			description: "should fail because access token invalid",
			setup: func() {
				config.DisableRefreshTokenValidation = true
				chgen.EXPECT().ValidateAccessToken(gomock.Any(), areq, "1234").Return(errorsx.WithStack(fosite.ErrInvalidTokenFormat))
			},
			expectErr: fosite.ErrInvalidTokenFormat,
		},
		{
			description: "should pass",
			setup: func() {
				chgen.EXPECT().ValidateAccessToken(gomock.Any(), areq, "1234").Return(nil)
			},
			expectTU: fosite.AccessToken,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			tu, err := v.IntrospectToken(context.Background(), fosite.AccessTokenFromRequest(httpreq), fosite.AccessToken, areq, []string{})

			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.expectTU, tu)
			}
		})
	}
}
