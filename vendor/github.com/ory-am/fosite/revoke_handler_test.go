package fosite_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/ory-am/fosite"
	"github.com/ory-am/fosite/internal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewRevocationRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	client := internal.NewMockClient(ctrl)
	handler := internal.NewMockRevocationHandler(ctrl)
	hasher := internal.NewMockHasher(ctrl)
	defer ctrl.Finish()

	fosite := &Fosite{Store: store, Hasher: hasher}
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(nil, errors.New(""))
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(client, nil)
				client.EXPECT().GetHashedSecret().Return([]byte("foo"))
				client.EXPECT().IsPublic().Return(false)
				hasher.EXPECT().Compare(gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(errors.New(""))
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(client, nil)
				client.EXPECT().GetHashedSecret().Return([]byte("foo"))
				client.EXPECT().IsPublic().Return(false)
				hasher.EXPECT().Compare(gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(client, nil)
				client.EXPECT().GetHashedSecret().Return([]byte("foo"))
				client.EXPECT().IsPublic().Return(false)
				hasher.EXPECT().Compare(gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(client, nil)
				client.EXPECT().IsPublic().Return(true)
				hasher.EXPECT().Compare(gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(client, nil)
				client.EXPECT().GetHashedSecret().Return([]byte("foo"))
				client.EXPECT().IsPublic().Return(false)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				store.EXPECT().GetClient(gomock.Eq("foo")).Return(client, nil)
				client.EXPECT().GetHashedSecret().Return([]byte("foo"))
				client.EXPECT().IsPublic().Return(false)
				hasher.EXPECT().Compare(gomock.Eq([]byte("foo")), gomock.Eq([]byte("bar"))).Return(nil)
				handler.EXPECT().RevokeToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			handlers: RevocationHandlers{handler},
		},
	} {
		r := &http.Request{
			Header:   c.header,
			PostForm: c.form,
			Form:     c.form,
			Method:   c.method,
		}
		c.mock()
		ctx := NewContext()
		fosite.RevocationHandlers = c.handlers
		err := fosite.NewRevocationRequest(ctx, r)
		assert.True(t, errors.Cause(err) == c.expectErr, "%d\nwant: %s \ngot: %s", k, c.expectErr, err)
		if err != nil {
			t.Logf("Error occured: %v", err)
		}
		t.Logf("Passed test case %d", k)
	}
}
