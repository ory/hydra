// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestNewRevocationRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)
	handler := internal.NewMockRevocationHandler(ctrl)
	hasher := internal.NewMockHasher(ctrl)
	t.Cleanup(ctrl.Finish)

	client := &DefaultClient{}
	config := &Config{ClientSecretsHasher: hasher}
	fosite := &Fosite{Store: store, Config: config}
	for k, c := range []struct {
		header    http.Header
		form      url.Values
		mock      func()
		method    string
		expectErr error
		expect    *AccessRequest
		handlers  RevocationHandlers
	}{
		{
			header:    http.Header{},
			expectErr: ErrInvalidRequest,
			method:    "GET",
			mock:      func() {},
		},
		{
			header:    http.Header{},
			expectErr: ErrInvalidRequest,
			method:    "POST",
			mock:      func() {},
		},
		{
			header: http.Header{},
			method: "POST",
			form: url.Values{
				"token": {"foo"},
			},
			mock:      func() {},
			expectErr: ErrInvalidRequest,
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"token": {"foo"},
			},
			expectErr: ErrInvalidClient,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(nil, errors.New(""))
			},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"token": {"foo"},
			},
			expectErr: ErrInvalidClient,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Secret = []byte("foo")
				client.Public = false
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(errors.New(""))
			},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"token": {"foo"},
			},
			expectErr: nil,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Secret = []byte("foo")
				client.Public = false
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: RevocationHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"token":           {"foo"},
				"token_type_hint": {"access_token"},
			},
			expectErr: nil,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Secret = []byte("foo")
				client.Public = false
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: RevocationHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "")},
			},
			method: "POST",
			form: url.Values{
				"token":           {"foo"},
				"token_type_hint": {"refresh_token"},
			},
			expectErr: nil,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Public = true
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: RevocationHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"token":           {"foo"},
				"token_type_hint": {"refresh_token"},
			},
			expectErr: nil,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Secret = []byte("foo")
				client.Public = false
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: RevocationHandlers{handler},
		},
		{
			header: http.Header{
				"Authorization": {basicAuth("foo", "bar")},
			},
			method: "POST",
			form: url.Values{
				"token":           {"foo"},
				"token_type_hint": {"bar"},
			},
			expectErr: nil,
			mock: func() {
				store.EXPECT().ClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("foo")).Return(client, nil)
				client.Secret = []byte("foo")
				client.Public = false
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: RevocationHandlers{handler},
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
			config.RevocationHandlers = c.handlers
			err := fosite.NewRevocationRequest(ctx, r)

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteRevocationResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	hasher := internal.NewMockHasher(ctrl)
	t.Cleanup(ctrl.Finish)

	config := &Config{ClientSecretsHasher: hasher}
	fosite := &Fosite{Store: store, Config: config}

	type args struct {
		rw  *httptest.ResponseRecorder
		err error
	}
	cases := []struct {
		input      args
		expectCode int
	}{
		{
			input: args{
				rw:  httptest.NewRecorder(),
				err: ErrInvalidRequest,
			},
			expectCode: ErrInvalidRequest.CodeField,
		},
		{
			input: args{
				rw:  httptest.NewRecorder(),
				err: ErrInvalidClient,
			},
			expectCode: ErrInvalidClient.CodeField,
		},
		{
			input: args{
				rw:  httptest.NewRecorder(),
				err: nil,
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		fosite.WriteRevocationResponse(context.Background(), tc.input.rw, tc.input.err)
		assert.Equal(t, tc.expectCode, tc.input.rw.Code)
	}
}
