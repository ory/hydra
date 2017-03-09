package oauth2

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/internal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeCode_PopulateTokenEndpointResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockAuthorizeCodeGrantStorage(ctrl)
	ach := internal.NewMockAccessTokenStrategy(ctrl)
	rch := internal.NewMockRefreshTokenStrategy(ctrl)
	auch := internal.NewMockAuthorizeCodeStrategy(ctrl)
	aresp := internal.NewMockAccessResponder(ctrl)
	//mockcl := internal.NewMockClient(ctrl)
	defer ctrl.Finish()

	areq := fosite.NewAccessRequest(new(fosite.DefaultSession))
	httpreq := &http.Request{PostForm: url.Values{}}
	authreq := fosite.NewAuthorizeRequest()
	areq.Session = new(fosite.DefaultSession)

	h := AuthorizeExplicitGrantHandler{
		AuthorizeCodeGrantStorage: store,
		AuthorizeCodeStrategy:     auch,
		AccessTokenStrategy:       ach,
		RefreshTokenStrategy:      rch,
		ScopeStrategy:             fosite.HierarchicScopeStrategy,
	}
	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"123"}
			},
		},
		{
			description: "should fail because authcode not found",
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"authorization_code"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"authorization_code"},
				}
				httpreq.PostForm.Add("code", "authcode")
				auch.EXPECT().AuthorizeCodeSignature("authcode").AnyTimes().Return("authsig")
				store.EXPECT().GetAuthorizeCodeSession(nil, "authsig", gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because validation failed",
			setup: func() {
				store.EXPECT().GetAuthorizeCodeSession(nil, "authsig", gomock.Any()).AnyTimes().Return(authreq, nil)
				auch.EXPECT().ValidateAuthorizeCode(nil, areq, "authcode").Return(errors.New(""))
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because access token generation failed",
			setup: func() {
				authreq.GrantedScopes = []string{"offline"}
				auch.EXPECT().ValidateAuthorizeCode(nil, areq, "authcode").AnyTimes().Return(nil)
				ach.EXPECT().GenerateAccessToken(nil, areq).Return("", "", errors.New("error"))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because refresh token generation failed",
			setup: func() {
				ach.EXPECT().GenerateAccessToken(nil, areq).AnyTimes().Return("access.ats", "ats", nil)
				rch.EXPECT().GenerateRefreshToken(nil, areq).Return("", "", errors.New("error"))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because persisting failed",
			setup: func() {
				rch.EXPECT().GenerateRefreshToken(nil, areq).AnyTimes().Return("refresh.rts", "rts", nil)
				store.EXPECT().PersistAuthorizeCodeGrantSession(nil, "authsig", "ats", "rts", areq).Return(errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			setup: func() {
				areq.GrantedScopes = fosite.Arguments{"foo", "offline"}
				store.EXPECT().PersistAuthorizeCodeGrantSession(nil, "authsig", "ats", "rts", areq).Return(nil)

				aresp.EXPECT().SetAccessToken("access.ats")
				aresp.EXPECT().SetTokenType("bearer")
				aresp.EXPECT().SetExtra("refresh_token", "refresh.rts")
				aresp.EXPECT().SetExpiresIn(gomock.Any())
				aresp.EXPECT().SetScopes(areq.GrantedScopes)
			},
		},
	} {
		c.setup()
		err := h.PopulateTokenEndpointResponse(nil, httpreq, areq, aresp)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}

func TestAuthorizeCode_HandleTokenEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockAuthorizeCodeGrantStorage(ctrl)
	ach := internal.NewMockAuthorizeCodeStrategy(ctrl)
	defer ctrl.Finish()

	authreq := fosite.NewAuthorizeRequest()
	areq := fosite.NewAccessRequest(nil)
	httpreq := &http.Request{PostForm: url.Values{}}
	areq.Session = new(fosite.DefaultSession)
	authreq.Session = new(fosite.DefaultSession)

	h := AuthorizeExplicitGrantHandler{
		AuthorizeCodeGrantStorage: store,
		AuthorizeCodeStrategy:     ach,
		ScopeStrategy:             fosite.HierarchicScopeStrategy,
	}
	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"12345678"} // grant_type REQUIRED. Value MUST be set to "authorization_code".
			},
		},
		{
			description: "should fail because authcode could not be retrieved (1)",
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"authorization_code"} // grant_type REQUIRED. Value MUST be set to "authorization_code".
				httpreq.PostForm = url.Values{"code": {"foo.bar"}}
				ach.EXPECT().AuthorizeCodeSignature("foo.bar").AnyTimes().Return("bar")
				store.EXPECT().GetAuthorizeCodeSession(nil, "bar", gomock.Any()).Return(nil, fosite.ErrNotFound)
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because authcode validation failed",
			setup: func() {
				store.EXPECT().GetAuthorizeCodeSession(nil, "bar", gomock.Any()).AnyTimes().Return(authreq, nil)
				ach.EXPECT().ValidateAuthorizeCode(nil, areq, "foo.bar").Return(errors.New(""))
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because client mismatch",
			setup: func() {
				ach.EXPECT().ValidateAuthorizeCode(nil, areq, "foo.bar").AnyTimes().Return(nil)

				areq.Client = &fosite.DefaultClient{ID: "foo"}
				authreq.Scopes = fosite.Arguments{"a", "b"}
				authreq.Client = &fosite.DefaultClient{ID: "bar"}
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because redirect uri not provided",
			setup: func() {
				authreq.Form.Add("redirect_uri", "request-redir")
				authreq.Client = &fosite.DefaultClient{ID: "foo"}
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should pass (2)",
			setup: func() {
				httpreq.PostForm = url.Values{"code": []string{"foo.bar"}}
				authreq.Form.Del("redirect_uri")
				authreq.RequestedAt = time.Now().Add(time.Hour)
			},
		},
	} {
		c.setup()
		err := h.HandleTokenEndpointRequest(nil, httpreq, areq)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}
