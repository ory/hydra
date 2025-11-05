// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestNewAccessRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)
	handler := internal.NewMockTokenEndpointHandler(ctrl)
	handler.EXPECT().CanHandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	handler.EXPECT().CanSkipClientAuth(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	hasher := internal.NewMockHasher(ctrl)
	t.Cleanup(ctrl.Finish)

	client := &DefaultClient{}
	config := &Config{ClientSecretsHasher: hasher, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}
	fosite := &Fosite{Store: store, Config: config}
	for k, c := range []struct {
		header    http.Header
		form      url.Values
		mock      func()
		method    string
		expectErr error
		expect    *AccessRequest
		handlers  TokenEndpointHandlers
	}{
		{
			header:    http.Header{},
			expectErr: ErrInvalidRequest,
			form:      url.Values{},
			method:    "POST",
			mock:      func() {},
		},
		{
			header: http.Header{},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock:      func() {},
			expectErr: ErrInvalidRequest,
		},
		{
			header: http.Header{},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
				"client_id":  {""},
			},
			expectErr: ErrInvalidRequest,
			mock:      func() {},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			expectErr: ErrInvalidClient,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(nil, errors.New(""))
			},
			handlers: TokenEndpointHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "GET",
			form: url.Values{
				"grant_type": {"foo"},
			},
			expectErr: ErrInvalidRequest,
			mock:      func() {},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			expectErr: ErrInvalidClient,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(nil, errors.New(""))
			},
			handlers: TokenEndpointHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			expectErr: ErrInvalidClient,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = false
				client.Secret = []byte("foo")
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(errors.New(""))
			},
			handlers: TokenEndpointHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			expectErr: ErrServerError,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = false
				client.Secret = []byte("foo")
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(ErrServerError)
			},
			handlers: TokenEndpointHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = false
				client.Secret = []byte("foo")
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: TokenEndpointHandlers{handler},
			expect: &AccessRequest{
				GrantTypes: Arguments{"foo"},
				Request: Request{
					Client: client,
				},
			},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = true
				handler.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: TokenEndpointHandlers{handler},
			expect: &AccessRequest{
				GrantTypes: Arguments{"foo"},
				Request: Request{
					Client: client,
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			r := &http.Request{
				Header:   c.header,
				PostForm: c.form,
				Form:     c.form,
				Method:   c.method,
			}
			c.mock()
			ctx := NewContext()
			config.TokenEndpointHandlers = c.handlers
			ar, err := fosite.NewAccessRequest(ctx, r, new(DefaultSession))

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				AssertObjectKeysEqual(t, c.expect, ar, "GrantTypes", "Client")
				assert.NotNil(t, ar.GetRequestedAt())
			}
		})
	}
}

func TestNewAccessRequestWithoutClientAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)
	handler := internal.NewMockTokenEndpointHandler(ctrl)
	handler.EXPECT().CanHandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	handler.EXPECT().CanSkipClientAuth(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	hasher := internal.NewMockHasher(ctrl)
	t.Cleanup(ctrl.Finish)

	client := &DefaultClient{}
	anotherClient := &DefaultClient{ID: "another"}
	config := &Config{ClientSecretsHasher: hasher, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}
	fosite := &Fosite{Store: store, Config: config}
	for k, c := range []struct {
		header    http.Header
		form      url.Values
		mock      func()
		method    string
		expectErr error
		expect    *AccessRequest
		handlers  TokenEndpointHandlers
	}{
		// No grant type -> error
		{
			form: url.Values{},
			mock: func() {
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Any()).Times(0)
			},
			method:    "POST",
			expectErr: ErrInvalidRequest,
		},
		// No registered handlers -> error
		{
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Any()).Times(0)
			},
			method:    "POST",
			expectErr: ErrInvalidRequest,
			handlers:  TokenEndpointHandlers{},
		},
		// Handler can skip client auth and ignores missing client.
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				// despite error from storage, we should success, because client auth is not required
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "foo").Return(nil, errors.New("no client")).Times(1)
				handler.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			method: "POST",
			expect: &AccessRequest{
				GrantTypes: Arguments{"foo"},
				Request: Request{
					Client: client,
				},
			},
			handlers: TokenEndpointHandlers{handler},
		},
		// Should pass if no auth is set in the header and can skip!
		{
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				handler.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			method: "POST",
			expect: &AccessRequest{
				GrantTypes: Arguments{"foo"},
				Request: Request{
					Client: client,
				},
			},
			handlers: TokenEndpointHandlers{handler},
		},
		// Should also pass if client auth is set!
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), "foo").Return(anotherClient, nil).Times(1)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				handler.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			method: "POST",
			expect: &AccessRequest{
				GrantTypes: Arguments{"foo"},
				Request: Request{
					Client: anotherClient,
				},
			},
			handlers: TokenEndpointHandlers{handler},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			r := &http.Request{
				Header:   c.header,
				PostForm: c.form,
				Form:     c.form,
				Method:   c.method,
			}
			c.mock()
			ctx := NewContext()
			config.TokenEndpointHandlers = c.handlers
			ar, err := fosite.NewAccessRequest(ctx, r, new(DefaultSession))

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				AssertObjectKeysEqual(t, c.expect, ar, "GrantTypes", "Client")
				assert.NotNil(t, ar.GetRequestedAt())
			}
		})
	}
}

// In this test case one handler requires client auth and another handler not.
func TestNewAccessRequestWithMixedClientAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)

	handlerWithClientAuth := internal.NewMockTokenEndpointHandler(ctrl)
	handlerWithClientAuth.EXPECT().CanHandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	handlerWithClientAuth.EXPECT().CanSkipClientAuth(gomock.Any(), gomock.Any()).Return(false).AnyTimes()

	handlerWithoutClientAuth := internal.NewMockTokenEndpointHandler(ctrl)
	handlerWithoutClientAuth.EXPECT().CanHandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	handlerWithoutClientAuth.EXPECT().CanSkipClientAuth(gomock.Any(), gomock.Any()).Return(true).AnyTimes()

	hasher := internal.NewMockHasher(ctrl)
	t.Cleanup(ctrl.Finish)

	client := &DefaultClient{}
	config := &Config{ClientSecretsHasher: hasher, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}
	fosite := &Fosite{Store: store, Config: config}
	for k, c := range []struct {
		header    http.Header
		form      url.Values
		mock      func()
		method    string
		expectErr error
		expect    *AccessRequest
		handlers  TokenEndpointHandlers
	}{
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = false
				client.Secret = []byte("foo")
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(errors.New("hash err"))
				handlerWithoutClientAuth.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			method:    "POST",
			expectErr: ErrInvalidClient,
			handlers:  TokenEndpointHandlers{handlerWithoutClientAuth, handlerWithClientAuth},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = false
				client.Secret = []byte("foo")
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handlerWithoutClientAuth.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
				handlerWithClientAuth.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			method: "POST",
			expect: &AccessRequest{
				GrantTypes: Arguments{"foo"},
				Request: Request{
					Client: client,
				},
			},
			handlers: TokenEndpointHandlers{handlerWithoutClientAuth, handlerWithClientAuth},
		},
		{
			header: http.Header{},
			form: url.Values{
				"grant_type": {"foo"},
			},
			mock: func() {
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Any()).Times(0)
				handlerWithoutClientAuth.EXPECT().HandleTokenEndpointRequest(gomock.Any(), gomock.Any()).Return(nil)
			},
			method:    "POST",
			expectErr: ErrInvalidRequest,
			handlers:  TokenEndpointHandlers{handlerWithoutClientAuth, handlerWithClientAuth},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			r := &http.Request{
				Header:   c.header,
				PostForm: c.form,
				Form:     c.form,
				Method:   c.method,
			}
			c.mock()
			ctx := NewContext()
			config.TokenEndpointHandlers = c.handlers
			ar, err := fosite.NewAccessRequest(ctx, r, new(DefaultSession))

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				AssertObjectKeysEqual(t, c.expect, ar, "GrantTypes", "Client")
				assert.NotNil(t, ar.GetRequestedAt())
			}
		})
	}
}

func basicAuth(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
}
