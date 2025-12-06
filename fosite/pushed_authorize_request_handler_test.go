// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

// Should pass
//
//   - https://openid.net/specs/oauth-v2-multiple-response-types-1_0.html#Terminology
//     The OAuth 2.0 specification allows for registration of space-separated response_type parameter values.
//     If a Response Type contains one of more space characters (%20), it is compared as a space-delimited list of
//     values in which the order of values does not matter.
func TestNewPushedAuthorizeRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)
	hasher := internal.NewMockHasher(ctrl)
	t.Cleanup(ctrl.Finish)

	config := &Config{
		ScopeStrategy:            ExactScopeStrategy,
		AudienceMatchingStrategy: DefaultAudienceMatchingStrategy,
		ClientSecretsHasher:      hasher,
	}

	fosite := &Fosite{
		Store:  store,
		Config: config,
	}

	redir, _ := url.Parse("https://foo.bar/cb")
	specialCharRedir, _ := url.Parse("web+application://callback")
	for _, c := range []struct {
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
			desc: "empty request fails",
			conf: fosite,
			r: &http.Request{
				Method: "POST",
			},
			expectedError: ErrInvalidClient,
			mock:          func() {},
		},
		/* invalid redirect uri */
		{
			desc:          "invalid redirect uri fails",
			conf:          fosite,
			query:         url.Values{"redirect_uri": []string{"invalid"}},
			expectedError: ErrInvalidClient,
			mock:          func() {},
		},
		/* invalid client */
		{
			desc:          "invalid client fails",
			conf:          fosite,
			query:         url.Values{"redirect_uri": []string{"https://foo.bar/cb"}},
			expectedError: ErrInvalidClient,
			mock:          func() {},
		},
		/* redirect client mismatch */
		{
			desc: "client and request redirects mismatch",
			conf: fosite,
			query: url.Values{
				"client_id":     []string{"1234"},
				"client_secret": []string{"1234"},
			},
			expectedError: ErrInvalidRequest,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"invalid"}, Scopes: []string{}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
		},
		/* redirect client mismatch */
		{
			desc: "client and request redirects mismatch",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  []string{""},
				"client_id":     []string{"1234"},
				"client_secret": []string{"1234"},
			},
			expectedError: ErrInvalidRequest,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"invalid"}, Scopes: []string{}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
		},
		/* redirect client mismatch */
		{
			desc: "client and request redirects mismatch",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  []string{"https://foo.bar/cb"},
				"client_id":     []string{"1234"},
				"client_secret": []string{"1234"},
			},
			expectedError: ErrInvalidRequest,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"invalid"}, Scopes: []string{}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
		},
		/* no state */
		{
			desc: "no state",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  []string{"https://foo.bar/cb"},
				"client_id":     []string{"1234"},
				"client_secret": []string{"1234"},
				"response_type": []string{"code"},
			},
			expectedError: ErrInvalidState,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
		},
		/* short state */
		{
			desc: "short state",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code"},
				"state":         {"short"},
			},
			expectedError: ErrInvalidState,
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
		},
		/* fails because scope not given */
		{
			desc: "should fail because client does not have scope baz",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar baz"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
			expectedError: ErrInvalidScope,
		},
		/* fails because scope not given */
		{
			desc: "should fail because client does not have scope baz",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"},
					Audience: []string{"https://cloud.ory.sh/api"},
					Secret:   []byte("1234"),
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
			expectedError: ErrInvalidRequest,
		},
		/* success case */
		{
			desc: "should pass",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					Secret:        []byte("1234"),
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
						Secret:   []byte("1234"),
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* repeated audience parameter */
		{
			desc: "repeated audience parameter",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					Secret:        []byte("1234"),
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
						Secret:   []byte("1234"),
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* repeated audience parameter with tricky values */
		{
			desc: "repeated audience parameter with tricky values",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"test value", ""},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"test value"},
					Secret:        []byte("1234"),
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
						Secret:   []byte("1234"),
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"test value"},
				},
			},
		},
		/* redirect_uri with special character in protocol*/
		{
			desc: "redirect_uri with special character",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"web+application://callback"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"web+application://callback"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					Secret:        []byte("1234"),
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
						Secret:   []byte("1234"),
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* audience with double spaces between values */
		{
			desc: "audience with double spaces between values",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api  https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{
					ResponseTypes: []string{"code token"},
					RedirectURIs:  []string{"https://foo.bar/cb"},
					Scopes:        []string{"foo", "bar"},
					Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
					Secret:        []byte("1234"),
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
						Secret:   []byte("1234"),
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* fails because unknown response_mode*/
		{
			desc: "should fail because unknown response_mode",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"unknown"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, ResponseTypes: []string{"code token"}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
			expectedError: ErrUnsupportedResponseMode,
		},
		/* fails because response_mode is requested but the OAuth 2.0 client doesn't support response mode */
		{
			desc: "should fail because response_mode is requested but the OAuth 2.0 client doesn't support response mode",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, ResponseTypes: []string{"code token"}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
			expectedError: ErrUnsupportedResponseMode,
		},
		/* fails because requested response mode is not allowed */
		{
			desc: "should fail because requested response mode is not allowed",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code token"},
						Secret:        []byte("1234"),
					},
					ResponseModes: []ResponseModeType{ResponseModeQuery},
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
			expectedError: ErrUnsupportedResponseMode,
		},
		/* success with response mode */
		{
			desc: "success with response mode",
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code token"},
						Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
						Secret:        []byte("1234"),
					},
					ResponseModes: []ResponseModeType{ResponseModeFormPost},
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
							Secret:        []byte("1234"),
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
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code"},
						Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
						Secret:        []byte("1234"),
					},
					ResponseModes: []ResponseModeType{ResponseModeQuery},
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
							Secret:        []byte("1234"),
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
			conf: fosite,
			query: url.Values{
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"audience":      {"https://cloud.ory.sh/api https://www.ory.sh/api"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultResponseModeClient{
					DefaultClient: &DefaultClient{
						RedirectURIs:  []string{"https://foo.bar/cb"},
						Scopes:        []string{"foo", "bar"},
						ResponseTypes: []string{"code token"},
						Audience:      []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
						Secret:        []byte("1234"),
					},
					ResponseModes: []ResponseModeType{ResponseModeFragment},
				}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
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
							Secret:        []byte("1234"),
						},
						ResponseModes: []ResponseModeType{ResponseModeFragment},
					},
					RequestedScope:    []string{"foo", "bar"},
					RequestedAudience: []string{"https://cloud.ory.sh/api", "https://www.ory.sh/api"},
				},
			},
		},
		/* fails because request_uri is included */
		{
			desc: "should fail because request_uri is provided in the request",
			conf: fosite,
			query: url.Values{
				"request_uri":   {"https://foo.bar/ru"},
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"1234"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, ResponseTypes: []string{"code token"}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("1234"))).Return(nil)
			},
			expectedError: ErrInvalidRequest.WithHint("The request must not contain 'request_uri'."),
		},
		/* fails because of invalid client credentials */
		{
			desc: "should fail because of invalid client creds",
			conf: fosite,
			query: url.Values{
				"request_uri":   {"https://foo.bar/ru"},
				"redirect_uri":  {"https://foo.bar/cb"},
				"client_id":     {"1234"},
				"client_secret": []string{"4321"},
				"response_type": {"code token"},
				"state":         {"strong-state"},
				"scope":         {"foo bar"},
				"response_mode": {"form_post"},
			},
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).MaxTimes(2)
				clientManager.EXPECT().GetClient(gomock.Any(), "1234").Return(&DefaultClient{RedirectURIs: []string{"https://foo.bar/cb"}, Scopes: []string{"foo", "bar"}, ResponseTypes: []string{"code token"}, Secret: []byte("1234")}, nil).MaxTimes(2)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("1234")), gomock.Eq([]byte("4321"))).Return(fmt.Errorf("invalid hash"))
			},
			expectedError: ErrInvalidClient,
		},
	} {
		t.Run(fmt.Sprintf("case=%s", c.desc), func(t *testing.T) {
			ctx := NewContext()

			c.mock()
			if c.r == nil {
				c.r = &http.Request{
					Header: http.Header{},
					Method: "POST",
				}
				if c.query != nil {
					c.r.URL = &url.URL{RawQuery: c.query.Encode()}
				}
			}

			ar, err := c.conf.NewPushedAuthorizeRequest(ctx, c.r)
			if c.expectedError != nil {
				assert.EqualError(t, err, c.expectedError.Error(), "Stack: %s", string(debug.Stack()))
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
