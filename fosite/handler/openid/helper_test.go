// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

var stratProvider = mockOpenIDConnectTokenStrategyProvider{
	strategy: openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return gen.MustRSAKey(), nil
			},
		},
		Config: &fosite.Config{
			MinParameterEntropy: fosite.MinParameterEntropy,
		},
	},
}

var fooErr = errors.New("foo")

func TestGenerateIDToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	chgen := internal.NewMockOpenIDConnectTokenStrategy(ctrl)
	chgenp := internal.NewMockOpenIDConnectTokenStrategyProvider(ctrl)
	t.Cleanup(ctrl.Finish)

	ar := fosite.NewAccessRequest(nil)
	sess := &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Subject: "peter",
		},
		Headers: &jwt.Headers{},
	}
	h := &openid.IDTokenHandleHelper{IDTokenStrategy: chgenp}

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
	}{
		{
			description: "should fail because generator failed",
			setup: func() {
				ar.Form.Set("nonce", "11111111111111111111111111111111111")
				ar.SetSession(sess)
				chgenp.EXPECT().OpenIDConnectTokenStrategy().Return(chgen).Times(1)
				chgen.EXPECT().GenerateIDToken(gomock.Any(), time.Duration(0), ar).Return("", fooErr)
			},
			expectErr: fooErr,
		},
		{
			description: "should pass",
			setup: func() {
				chgenp.EXPECT().OpenIDConnectTokenStrategy().Return(chgen).Times(1)
				chgen.EXPECT().GenerateIDToken(gomock.Any(), time.Duration(0), ar).AnyTimes().Return("asdf", nil)
			},
		},
	} {
		c.setup()
		token, err := openid.CallGenerateIDToken(context.Background(), time.Duration(0), ar, h)
		assert.True(t, err == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		if err == nil {
			assert.NotEmpty(t, token, "(%d) %s", k, c.description)
		}
		t.Logf("Passed test case %d", k)
	}
}

func TestIssueExplicitToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	resp := internal.NewMockAccessResponder(ctrl)
	t.Cleanup(ctrl.Finish)

	ar := fosite.NewAuthorizeRequest()
	ar.Form = url.Values{"nonce": {"111111111111"}}
	ar.SetSession(&openid.DefaultSession{Claims: &jwt.IDTokenClaims{
		Subject: "peter",
	}, Headers: &jwt.Headers{}})

	resp.EXPECT().SetExtra("id_token", gomock.Any())
	h := &openid.IDTokenHandleHelper{IDTokenStrategy: stratProvider}
	err := h.IssueExplicitIDToken(context.Background(), time.Duration(0), ar, resp)
	assert.NoError(t, err)
}

func TestIssueImplicitToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	resp := internal.NewMockAuthorizeResponder(ctrl)
	t.Cleanup(ctrl.Finish)

	ar := fosite.NewAuthorizeRequest()
	ar.Form = url.Values{"nonce": {"111111111111"}}
	ar.SetSession(&openid.DefaultSession{Claims: &jwt.IDTokenClaims{
		Subject: "peter",
	}, Headers: &jwt.Headers{}})

	resp.EXPECT().AddParameter("id_token", gomock.Any())
	h := &openid.IDTokenHandleHelper{IDTokenStrategy: stratProvider}
	err := h.IssueImplicitIDToken(context.Background(), time.Duration(0), ar, resp)
	assert.NoError(t, err)
}

func TestGetAccessTokenHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	req := internal.NewMockAccessRequester(ctrl)
	resp := internal.NewMockAccessResponder(ctrl)

	t.Cleanup(ctrl.Finish)

	req.EXPECT().GetSession().Return(nil)
	resp.EXPECT().GetAccessToken().Return("7a35f818-9164-48cb-8c8f-e1217f44228431c41102-d410-4ed5-9276-07ba53dfdcd8")

	h := &openid.IDTokenHandleHelper{IDTokenStrategy: stratProvider}

	hash := h.GetAccessTokenHash(context.Background(), req, resp)
	assert.Equal(t, "Zfn_XBitThuDJiETU3OALQ", hash)
}

func TestGetAccessTokenHashWithDifferentKeyLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	req := internal.NewMockAccessRequester(ctrl)
	resp := internal.NewMockAccessResponder(ctrl)

	t.Cleanup(ctrl.Finish)

	headers := &jwt.Headers{
		Extra: map[string]interface{}{
			"alg": "RS384",
		},
	}
	req.EXPECT().GetSession().Return(&openid.DefaultSession{Headers: headers})
	resp.EXPECT().GetAccessToken().Return("7a35f818-9164-48cb-8c8f-e1217f44228431c41102-d410-4ed5-9276-07ba53dfdcd8")

	h := &openid.IDTokenHandleHelper{IDTokenStrategy: stratProvider}

	hash := h.GetAccessTokenHash(context.Background(), req, resp)
	assert.Equal(t, "VNX38yiOyeqBPheW5jDsWQKa6IjJzK66", hash)
}

func TestGetAccessTokenHashWithBadAlg(t *testing.T) {
	ctrl := gomock.NewController(t)
	req := internal.NewMockAccessRequester(ctrl)
	resp := internal.NewMockAccessResponder(ctrl)

	t.Cleanup(ctrl.Finish)

	headers := &jwt.Headers{
		Extra: map[string]interface{}{
			"alg": "R",
		},
	}
	req.EXPECT().GetSession().Return(&openid.DefaultSession{Headers: headers})
	resp.EXPECT().GetAccessToken().Return("7a35f818-9164-48cb-8c8f-e1217f44228431c41102-d410-4ed5-9276-07ba53dfdcd8")

	h := &openid.IDTokenHandleHelper{IDTokenStrategy: stratProvider}

	hash := h.GetAccessTokenHash(context.Background(), req, resp)
	assert.Equal(t, "Zfn_XBitThuDJiETU3OALQ", hash)
}

func TestGetAccessTokenHashWithMissingKeyLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	req := internal.NewMockAccessRequester(ctrl)
	resp := internal.NewMockAccessResponder(ctrl)

	t.Cleanup(ctrl.Finish)

	headers := &jwt.Headers{
		Extra: map[string]interface{}{
			"alg": "RS",
		},
	}
	req.EXPECT().GetSession().Return(&openid.DefaultSession{Headers: headers})
	resp.EXPECT().GetAccessToken().Return("7a35f818-9164-48cb-8c8f-e1217f44228431c41102-d410-4ed5-9276-07ba53dfdcd8")

	h := &openid.IDTokenHandleHelper{IDTokenStrategy: stratProvider}

	hash := h.GetAccessTokenHash(context.Background(), req, resp)
	assert.Equal(t, "Zfn_XBitThuDJiETU3OALQ", hash)
}
