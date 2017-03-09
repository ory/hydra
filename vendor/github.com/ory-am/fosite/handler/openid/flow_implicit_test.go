package openid

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/fosite/storage"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestImplicit_HandleAuthorizeEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aresp := fosite.NewAuthorizeResponse()
	areq := fosite.NewAuthorizeRequest()
	httpreq := &http.Request{Form: url.Values{}}
	areq.Session = new(fosite.DefaultSession)

	h := OpenIDConnectImplicitHandler{
		AuthorizeImplicitGrantTypeHandler: &oauth2.AuthorizeImplicitGrantTypeHandler{
			AccessTokenLifespan: time.Hour,
			AccessTokenStrategy: hmacStrategy,
			AccessTokenStorage:  storage.NewMemoryStore(),
		},
		IDTokenHandleHelper: &IDTokenHandleHelper{
			IDTokenStrategy: idStrategy,
		},
		ScopeStrategy: fosite.HierarchicScopeStrategy,
	}
	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		check       func()
	}{
		{
			description: "should not do anything because request requirements are not met",
			setup:       func() {},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.State = "foostate"
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{}
				areq.GrantedScopes = fosite.Arguments{"openid"}
			},
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
				areq.Scopes = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{},
					ResponseTypes: fosite.Arguments{},
					Scopes:        []string{"openid", "fosite"},
				}
			},
			expectErr: fosite.ErrInvalidGrant,
		},
		{
			description: "should not do anything because request requirements are not met",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
				areq.Scopes = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{},
					Scopes:        []string{"openid", "fosite"},
				}
			},
			expectErr: fosite.ErrInvalidGrant,
		},
		{
			description: "should fail because session not set",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"id_token"}
				areq.Scopes = fosite.Arguments{"openid"}
				areq.Client = &fosite.DefaultClient{
					GrantTypes:    fosite.Arguments{"implicit"},
					ResponseTypes: fosite.Arguments{"token", "id_token"},
					Scopes:        []string{"openid", "fosite"},
				}
			},
			expectErr: ErrInvalidSession,
		},
		{
			description: "should fail because nonce not set",
			setup: func() {
				areq.Session = &DefaultSession{
					Claims: &jwt.IDTokenClaims{
						Subject: "peter",
					},
					Headers: &jwt.Headers{},
				}
				areq.Form.Add("nonce", "some-random-foo-nonce-wow")
			},
		},
		{
			description: "should pass",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"id_token"}
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetFragment().Get("id_token"))
				assert.NotEmpty(t, aresp.GetFragment().Get("state"))
				assert.Empty(t, aresp.GetFragment().Get("access_token"))
			},
		},
		{
			description: "should pass",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"token", "id_token"}
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetFragment().Get("id_token"))
				assert.NotEmpty(t, aresp.GetFragment().Get("state"))
				assert.NotEmpty(t, aresp.GetFragment().Get("access_token"))
			},
		},
		{
			description: "should pass",
			setup: func() {
				areq.ResponseTypes = fosite.Arguments{"id_token", "token"}
				areq.Scopes = fosite.Arguments{"fosite", "openid"}
			},
			check: func() {
				assert.NotEmpty(t, aresp.GetFragment().Get("id_token"))
				assert.NotEmpty(t, aresp.GetFragment().Get("state"))
				assert.NotEmpty(t, aresp.GetFragment().Get("access_token"))
			},
		},
	} {
		c.setup()
		err := h.HandleAuthorizeEndpointRequest(nil, httpreq, areq, aresp)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
		if c.check != nil {
			c.check()
		}
	}
}
