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

func TestRefreshFlow_HandleTokenEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockRefreshTokenGrantStorage(ctrl)
	chgen := internal.NewMockRefreshTokenStrategy(ctrl)
	defer ctrl.Finish()

	areq := fosite.NewAccessRequest(nil)
	sess := &fosite.DefaultSession{Subject: "othersub"}
	httpreq := &http.Request{PostForm: url.Values{}}

	h := RefreshTokenGrantHandler{
		RefreshTokenGrantStorage: store,
		RefreshTokenStrategy:     chgen,
		AccessTokenLifespan:      time.Hour,
	}
	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		expect      func()
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"123"}
			},
		},
		{
			description: "should fail because lookup failed",
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"refresh_token"}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"refresh_token"}}
				httpreq.PostForm.Add("refresh_token", "some.refreshtokensig")
				chgen.EXPECT().RefreshTokenSignature("some.refreshtokensig").AnyTimes().Return("refreshtokensig")
				store.EXPECT().GetRefreshTokenSession(nil, "refreshtokensig", nil).Return(nil, fosite.ErrNotFound)
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because validation failed",
			setup: func() {
				store.EXPECT().GetRefreshTokenSession(nil, "refreshtokensig", nil).Return(&fosite.Request{
					GrantedScopes: []string{"offline"},
				}, nil)
				chgen.EXPECT().ValidateRefreshToken(nil, areq, "some.refreshtokensig").Return(errors.New(""))
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should fail because client mismatches",
			setup: func() {
				areq.Client = &fosite.DefaultClient{
					ID:         "foo",
					GrantTypes: fosite.Arguments{"refresh_token"},
				}
				store.EXPECT().GetRefreshTokenSession(nil, "refreshtokensig", nil).Return(&fosite.Request{
					Client:        &fosite.DefaultClient{ID: ""},
					GrantedScopes: []string{"offline"},
				}, nil)
				chgen.EXPECT().ValidateRefreshToken(nil, areq, "some.refreshtokensig").AnyTimes().Return(nil)
			},
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			description: "should pass",
			setup: func() {
				store.EXPECT().GetRefreshTokenSession(nil, "refreshtokensig", nil).Return(&fosite.Request{
					Client:        &fosite.DefaultClient{ID: "foo"},
					GrantedScopes: fosite.Arguments{"foo", "offline"},
					Scopes:        fosite.Arguments{"foo", "bar"},
					Session:       sess,
					Form:          url.Values{"foo": []string{"bar"}},
					RequestedAt:   time.Now().Add(-time.Hour).Round(time.Hour),
				}, nil)
			},
			expect: func() {
				assert.NotEqual(t, sess, areq.Session)
				assert.NotEqual(t, time.Now().Add(-time.Hour).Round(time.Hour), areq.RequestedAt)
				assert.Equal(t, fosite.Arguments{"foo", "offline"}, areq.GrantedScopes)
				assert.Equal(t, fosite.Arguments{"foo", "bar"}, areq.Scopes)
				assert.NotEqual(t, url.Values{"foo": []string{"bar"}}, areq.Form)
			},
		},
	} {
		c.setup()
		err := h.HandleTokenEndpointRequest(nil, httpreq, areq)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		if c.expect != nil {
			c.expect()
		}
		t.Logf("Passed test case %d", k)
	}
}

func TestRefreshFlow_PopulateTokenEndpointResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockRefreshTokenGrantStorage(ctrl)
	rcts := internal.NewMockRefreshTokenStrategy(ctrl)
	acts := internal.NewMockAccessTokenStrategy(ctrl)
	areq := fosite.NewAccessRequest(nil)
	aresp := internal.NewMockAccessResponder(ctrl)
	httpreq := &http.Request{PostForm: url.Values{}}
	defer ctrl.Finish()

	areq.Client = &fosite.DefaultClient{}
	h := RefreshTokenGrantHandler{
		RefreshTokenGrantStorage: store,
		RefreshTokenStrategy:     rcts,
		AccessTokenStrategy:      acts,
		AccessTokenLifespan:      time.Hour,
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
				areq.GrantTypes = fosite.Arguments{"313"}
			},
		},
		{
			description: "should fail because access token generation fails",
			setup: func() {
				areq.GrantTypes = fosite.Arguments{"refresh_token"}
				httpreq.PostForm.Add("refresh_token", "foo.reftokensig")
				rcts.EXPECT().RefreshTokenSignature("foo.reftokensig").AnyTimes().Return("reftokensig")
				acts.EXPECT().GenerateAccessToken(nil, areq).Return("", "", errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because access token generation fails",
			setup: func() {
				acts.EXPECT().GenerateAccessToken(nil, areq).AnyTimes().Return("access.atsig", "atsig", nil)
				rcts.EXPECT().GenerateRefreshToken(nil, areq).Return("", "", errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should fail because persisting fails",
			setup: func() {
				rcts.EXPECT().GenerateRefreshToken(nil, areq).AnyTimes().Return("refresh.resig", "resig", nil)
				store.EXPECT().PersistRefreshTokenGrantSession(nil, "reftokensig", "atsig", "resig", areq).Return(errors.New(""))
			},
			expectErr: fosite.ErrServerError,
		},
		{
			description: "should pass",
			setup: func() {
				areq.Session = &fosite.DefaultSession{}
				store.EXPECT().PersistRefreshTokenGrantSession(nil, "reftokensig", "atsig", "resig", areq).AnyTimes().Return(nil)

				aresp.EXPECT().SetAccessToken("access.atsig")
				aresp.EXPECT().SetTokenType("bearer")
				aresp.EXPECT().SetExpiresIn(gomock.Any())
				aresp.EXPECT().SetScopes(gomock.Any())
				aresp.EXPECT().SetExtra("refresh_token", "refresh.resig")
			},
		},
	} {
		c.setup()
		err := h.PopulateTokenEndpointResponse(nil, httpreq, areq, aresp)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}
