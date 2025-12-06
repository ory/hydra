// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	. "github.com/ory/hydra/v2/fosite/internal"
)

// Should pass
//
//   - https://openid.net/specs/oauth-v2-multiple-response-types-1_0.html#Terminology
//     The OAuth 2.0 specification allows for registration of space-separated response_type parameter values.
//     If a Response Type contains one of more space characters (%20), it is compared as a space-delimited list of
//     values in which the order of values does not matter.
func TestNewAuthorizeRequest(t *testing.T) {
	var store *MockStorage
	var clientManager *MockClientManager

	redir, _ := url.Parse("https://foo.bar/cb")
	specialCharRedir, _ := url.Parse("web+application://callback")
	for k, c := range []struct {
		desc          string
		conf          *Fosite
		r             *http.Request
		query         url.Values
		expectedError error
		mock          func()
		expect        *AuthorizeRequest
	}{
		/* empty request */
		{
			desc:          "empty request fails",
			conf:          &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			r:             &http.Request{},
			expectedError: ErrInvalidClient,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo"))
			},
		},
		/* invalid redirect uri */
		{
			desc:          "invalid redirect uri fails",
			conf:          &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query:         url.Values{"redirect_uri": []string{"invalid"}},
			expectedError: ErrInvalidClient,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo"))
			},
		},
		/* invalid client */
		{
			desc:          "invalid client fails",
			conf:          &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query:         url.Values{"redirect_uri": []string{"https://foo.bar/cb"}},
			expectedError: ErrInvalidClient,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo"))
			},
		},
		/* redirect client mismatch */
		{
			desc: "client and request redirects mismatch",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"client_id": []string{"1234"},
			},
			expectedError: ErrInvalidRequest,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"invalid"}, Scopes: []string{}}, nil)
			},
		},
		/* redirect client mismatch */
		{
			desc: "client and request redirects mismatch",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri": []string{""},
				"client_id":    []string{"1234"},
			},
			expectedError: ErrInvalidRequest,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"invalid"}, Scopes: []string{}}, nil)
			},
		},
		/* redirect client mismatch */
		{
			desc: "client and request redirects mismatch",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri": []string{"https://foo.bar/cb"},
				"client_id":    []string{"1234"},
			},
			expectedError: ErrInvalidRequest,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"invalid"}, Scopes: []string{}}, nil)
			},
		},
		/* no state */
		{
			desc: "no state",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  []string{"https://foo.bar/cb"},
				"client_id":     []string{"1234"},
				"response_type": []string{"code"},
			},
			expectedError: ErrInvalidState,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{}}, nil)
			},
		},
		/* short state */
		{
			desc: "short state",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code"},
				"state":         {"short"},
			},
			expectedError: ErrInvalidState,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{}}, nil)
			},
		},
		/* fails because scope not given */
		{
			desc: "should fail because client does not have scope baz",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar baz"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}}, nil)
			},
			expectedError: ErrInvalidScope,
		},
		/* fails because scope not given */
		{
			desc: "should fail because client does not have scope baz",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"},
					Audience: []string{"https://cloud.ory.sh/api"},
				}, nil)
			},
			expectedError: ErrInvalidRequest,
		},
		/* success case */
		{
			desc: "should pass",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultClient{
						ResponseTypes: []string{"code token"}, RedirectURIs: []string{"https://foo.bar/cb"},
						Scopes:   []string{"foo", "bar"},
						Audience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* repeated audience parameter */
		{
			desc: "repeated audience parameter",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultClient{
						ResponseTypes: []string{"code token"}, RedirectURIs: []string{"https://foo.bar/cb"},
						Scopes:   []string{"foo", "bar"},
						Audience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* repeated audience parameter with tricky values */
		{
			desc: "repeated audience parameter with tricky values",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: ExactAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"test value", ""},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"test value"},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultClient{
						ResponseTypes: []string{"code token"}, RedirectURIs: []string{"https://foo.bar/cb"},
						Scopes:   []string{"foo", "bar"},
						Audience: []string{"test value"},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"test value"},
				},
			},
		},
		/* redirect_uri with special character in protocol*/
		{
			desc: "redirect_uri with special character",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"web+application://callback"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"web+application://callback"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   specialCharRedir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultClient{
						ResponseTypes: []string{"code token"}, RedirectURIs: []string{"web+application://callback"},
						Scopes:   []string{"foo", "bar"},
						Audience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* audience with double spaces between values */
		{
			desc: "audience with double spaces between values",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api  https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultClient{
						ResponseTypes: []string{"code token"}, RedirectURIs: []string{"https://foo.bar/cb"},
						Scopes:   []string{"foo", "bar"},
						Audience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* fails because unknown response_mode*/
		{
			desc: "should fail because unknown response_mode",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"unknown"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, ResponseTypes: []string{"code token"}}, nil)
			},
			expectedError: ErrUnsupportedResponseMode,
		},
		/* fails because response_mode is requested but the OAuth 2.0 client doesn't support response mode */
		{
			desc: "should fail because response_mode is requested but the OAuth 2.0 client doesn't support response mode",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, ResponseTypes: []string{"code token"}}, nil)
			},
			expectedError: ErrUnsupportedResponseMode,
		},
		/* fails because requested response mode is not allowed */
		{
			desc: "should fail because requested response mode is not allowed",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code token"},
					},
					ResponseModes: []ResponseModeType{ResponseModeQuery},
				}, nil)
			},
			expectedError: ErrUnsupportedResponseMode,
		},
		/* success with response mode */
		{
			desc: "success with response mode",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code token"},
						Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					ResponseModes: []ResponseModeType{ResponseModeFormPost},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultResponseModeClient{
						DefaultClient: &DefaultClient{
							RedirectURIs:  []string{"https://foo.bar/cb"},
							Scopes:        []string{"foo", "bar"},
							ResponseTypes: []string{"code token"},
							Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
						},
						ResponseModes: []ResponseModeType{ResponseModeFormPost},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* determine correct response mode if default */
		{
			desc: "success with response mode",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code"},
						Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					ResponseModes: []ResponseModeType{ResponseModeQuery},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultResponseModeClient{
						DefaultClient: &DefaultClient{
							RedirectURIs:  []string{"https://foo.bar/cb"},
							Scopes:        []string{"foo", "bar"},
							ResponseTypes: []string{"code"},
							Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
						},
						ResponseModes: []ResponseModeType{ResponseModeQuery},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* determine correct response mode if default */
		{
			desc: "success with response mode",
			conf: &Fosite{Store: store, Config: &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}},
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code token"},
						Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					},
					ResponseModes: []ResponseModeType{ResponseModeFragment},
				}, nil)
			},
			expect: &AuthorizeRequest{
				RedirectURI:   redir,
				ResponseTypes: []string{"code", "token"},
				State:         "strong-state",
				Request: Request{
					Client: &DefaultResponseModeClient{
						DefaultClient: &DefaultClient{
							RedirectURIs:  []string{"https://foo.bar/cb"},
							Scopes:        []string{"foo", "bar"},
							ResponseTypes: []string{"code token"},
							Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
						},
						ResponseModes: []ResponseModeType{ResponseModeFragment},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store = NewMockStorage(ctrl)
			clientManager = NewMockClientManager(ctrl)
			t.Cleanup(ctrl.Finish)

			c.mock()
			if c.r == nil {
				c.r = &http.Request{Header: http.Header{}}
				if c.query != nil {
					c.r.URL = &url.URL{RawQuery: c.query.Encode()}
				}
			}

			c.conf.Store = store
			ar, err := c.conf.NewAuthorizeRequest(context.Background(), c.r)
			if c.expectedError != nil {
				assert.EqualError(t, err, c.expectedError.Error())
				// https://github.com/ory/hydra/issues/1642
				AssertObjectKeysEqual(t, &AuthorizeRequest{State: c.query.Get("state")}, ar, "State")
			} else {
				require.NoError(t, err)
				AssertObjectKeysEqual(t, c.expect, ar, "ResponseTypes", "RequestedAudience", "RequestedScope", "Client", "RedirectURI", "State")
				assert.NotNil(t, ar.GetRequestedAt())
			}
		})
	}
}
