// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package verifiable

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

type mockNonceManager struct{ t *testing.T }

func (m *mockNonceManager) NewNonce(ctx context.Context, accessToken string, expiresAt time.Time) (string, error) {
	assert.Equal(m.t, "fake access token", accessToken)
	assert.WithinDuration(m.t, time.Now().Add(time.Hour), expiresAt, 5*time.Second)
	return "mocked nonce", nil
}

func (m *mockNonceManager) IsNonceValid(context.Context, string, string) error {
	return nil
}

func TestHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("case=correct scopes", func(t *testing.T) {
		t.Parallel()
		handler := newHandler(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		req := internal.NewMockAccessRequester(ctrl)
		req.EXPECT().GetGrantedScopes().Return(fosite.Arguments{"openid", draftScope}).AnyTimes()

		resp := internal.NewMockAccessResponder(ctrl)
		resp.EXPECT().GetAccessToken().Return("fake access token")
		resp.EXPECT().SetExtra(gomock.Eq(draftNonceField), gomock.Eq("mocked nonce"))
		resp.EXPECT().SetExtra(gomock.Eq(draftNonceExpField), gomock.Any())

		assert.NoError(t, handler.HandleTokenEndpointRequest(ctx, req))
		assert.NoError(t, handler.PopulateTokenEndpointResponse(ctx, req, resp))
	})

	t.Run("case=incorrect scopes", func(t *testing.T) {
		t.Parallel()
		handler := newHandler(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		req := internal.NewMockAccessRequester(ctrl)
		req.EXPECT().GetGrantedScopes().Return(fosite.Arguments{"openid"}).AnyTimes()

		resp := internal.NewMockAccessResponder(ctrl)

		assert.ErrorIs(t, handler.HandleTokenEndpointRequest(ctx, req), fosite.ErrUnknownRequest)
		assert.ErrorIs(t, handler.PopulateTokenEndpointResponse(ctx, req, resp), fosite.ErrUnknownRequest)
	})
}

func newHandler(t *testing.T) *Handler {
	return &Handler{
		Config:       new(fosite.Config),
		NonceManager: &mockNonceManager{t: t},
	}
}
