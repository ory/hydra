// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestGetExpiresIn(t *testing.T) {
	now := time.Now().UTC()
	r := fosite.NewAccessRequest(&fosite.DefaultSession{
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken: now.Add(time.Hour),
		},
	})
	assert.Equal(t, time.Hour, oauth2.CallGetExpiresIn(r, fosite.AccessToken, time.Millisecond, now))
}

func TestIssueAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	areq := &fosite.AccessRequest{}
	aresp := &fosite.AccessResponse{Extra: map[string]interface{}{}}
	accessStrat := internal.NewMockAccessTokenStrategy(ctrl)
	accessStore := internal.NewMockAccessTokenStorage(ctrl)
	provider := internal.NewMockAccessTokenStorageProvider(ctrl)
	t.Cleanup(ctrl.Finish)

	helper := oauth2.HandleHelper{
		Storage:             provider,
		AccessTokenStrategy: accessStrat,
		Config: &fosite.Config{
			AccessTokenLifespan: time.Hour,
		},
	}

	areq.Session = &fosite.DefaultSession{}
	for k, c := range []struct {
		mock func()
		err  error
	}{
		{
			mock: func() {
				accessStrat.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return("", "", errors.New(""))
			},
			err: errors.New(""),
		},
		{
			mock: func() {
				accessStrat.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return("token", "signature", nil)
				provider.EXPECT().AccessTokenStorage().Return(accessStore).Times(1)
				accessStore.EXPECT().CreateAccessTokenSession(gomock.Any(), "signature", gomock.Eq(areq.Sanitize([]string{}))).Return(errors.New(""))
			},
			err: errors.New(""),
		},
		{
			mock: func() {
				accessStrat.EXPECT().GenerateAccessToken(gomock.Any(), areq).Return("token", "signature", nil)
				provider.EXPECT().AccessTokenStorage().Return(accessStore).Times(1)
				accessStore.EXPECT().CreateAccessTokenSession(gomock.Any(), "signature", gomock.Eq(areq.Sanitize([]string{}))).Return(nil)
			},
			err: nil,
		},
	} {
		c.mock()
		signature, err := helper.IssueAccessToken(context.Background(), helper.Config.GetAccessTokenLifespan(context.TODO()), areq, aresp)
		require.Equal(t, err == nil, c.err == nil)
		if c.err != nil {
			assert.EqualError(t, err, c.err.Error(), "Case %d", k)
		} else {
			assert.NotEmpty(t, signature, "Case %d", k)
		}
	}
}
