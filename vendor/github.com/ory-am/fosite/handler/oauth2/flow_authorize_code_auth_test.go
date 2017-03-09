package oauth2

import (
	"net/http"
	"net/url"
	"testing"

	"strings"

	"github.com/golang/mock/gomock"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/internal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeCode_HandleAuthorizeEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockAuthorizeCodeGrantStorage(ctrl)
	chgen := internal.NewMockAuthorizeCodeStrategy(ctrl)
	aresp := internal.NewMockAuthorizeResponder(ctrl)
	defer ctrl.Finish()

	areq := fosite.NewAuthorizeRequest()
	httpreq := &http.Request{Form: url.Values{}}

	areq.Session = new(fosite.DefaultSession)
	h := AuthorizeExplicitGrantHandler{
		AuthorizeCodeGrantStorage: store,
		AuthorizeCodeStrategy:     chgen,
		ScopeStrategy:             fosite.HierarchicScopeStrategy,
	}
	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
	}{
		{
			description: "should pass because not responsible for handling an empty response type",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{""}
			},
		},
		{
			description: "should pass because not responsible for handling an invalid response type",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"foo"}
			},
		},
		{
			description: "should fail because redirect uri is not https",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"code"}
				areq.Client = &fosite.DefaultClient{ResponseTypes: fosite.Arguments{"code"}}
				areq.RedirectURI, _ = url.Parse("http://asdf.de/cb")
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because authorize code generation failed",
			setup: func() {
				areq.RedirectURI, _ = url.Parse("https://foobar.com/cb")
				chgen.EXPECT().GenerateAuthorizeCode(nil, areq).Return("", "", errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because could not presist authorize code session",
			setup: func() {
				chgen.EXPECT().GenerateAuthorizeCode(nil, areq).AnyTimes().Return("someauthcode.authsig", "authsig", nil)
				store.EXPECT().CreateAuthorizeCodeSession(nil, "authsig", areq).Return(errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			setup: func() {
				areq.GrantedScopes = fosite.Arguments{"a", "b"}
				areq.State = "superstate"
				store.EXPECT().CreateAuthorizeCodeSession(nil, "authsig", areq).Return(nil)
				aresp.EXPECT().AddQuery("code", "someauthcode.authsig")
				aresp.EXPECT().AddQuery("scope", strings.Join(areq.GrantedScopes, " "))
				aresp.EXPECT().AddQuery("state", areq.State)
			},
		},
	} {
		c.setup()
		err := h.HandleAuthorizeEndpointRequest(nil, httpreq, areq, aresp)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}
