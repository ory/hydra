// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	. "github.com/ory/hydra/v2/fosite/internal"
)

func TestNewPushedAuthorizeResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	handlers := []*MockPushedAuthorizeEndpointHandler{NewMockPushedAuthorizeEndpointHandler(ctrl)}
	ar := NewMockAuthorizeRequester(ctrl)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	oauth2 := &Fosite{
		Config: &Config{
			PushedAuthorizeEndpointHandlers: PushedAuthorizeEndpointHandlers{handlers[0]},
		},
	}
	ar.EXPECT().SetSession(gomock.Eq(new(DefaultSession))).AnyTimes()
	fooErr := errors.New("foo")
	for k, c := range []struct {
		isErr     bool
		mock      func()
		expectErr error
	}{
		{
			mock: func() {
				handlers[0].EXPECT().HandlePushedAuthorizeEndpointRequest(gomock.Any(), gomock.Eq(ar), gomock.Any()).Return(fooErr)
			},
			isErr:     true,
			expectErr: fooErr,
		},
		{
			mock: func() {
				handlers[0].EXPECT().HandlePushedAuthorizeEndpointRequest(gomock.Any(), gomock.Eq(ar), gomock.Any()).Return(nil)
			},
			isErr: false,
		},
	} {
		c.mock()
		responder, err := oauth2.NewPushedAuthorizeResponse(ctx, ar, new(DefaultSession))
		assert.Equal(t, c.isErr, err != nil, "%d: %s", k, err)
		if err != nil {
			assert.Equal(t, c.expectErr, err, "%d: %s", k, err)
			assert.Nil(t, responder, "%d", k)
		} else {
			assert.NotNil(t, responder, "%d", k)
		}
		t.Logf("Passed test case %d", k)
	}
}
