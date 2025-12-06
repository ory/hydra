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
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestNewDeviceRequestWithPublicClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)
	deviceClient := &DefaultClient{ID: "client_id"}
	deviceClient.Public = true
	deviceClient.Scopes = []string{"17", "42"}
	deviceClient.Audience = []string{"aud2"}
	deviceClient.GrantTypes = []string{"urn:ietf:params:oauth:grant-type:device_code"}

	authCodeClient := &DefaultClient{ID: "client_id_2"}
	authCodeClient.Public = true
	authCodeClient.Scopes = []string{"17", "42"}
	authCodeClient.GrantTypes = []string{"authorization_code"}

	t.Cleanup(ctrl.Finish)
	config := &Config{ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}
	fosite := &Fosite{Store: store, Config: config}
	for k, c := range []struct {
		header        http.Header
		form          url.Values
		method        string
		expectedError error
		mock          func()
		expect        DeviceRequester
		description   string
	}{{
		description:   "invalid method",
		expectedError: ErrInvalidRequest,
		method:        "GET",
		mock:          func() {},
	}, {
		description:   "empty request",
		expectedError: ErrInvalidRequest,
		method:        "POST",
		mock:          func() {},
	}, {
		description: "invalid client",
		form: url.Values{
			"client_id": {"client_id"},
			"scope":     {"foo bar"},
		},
		expectedError: ErrInvalidClient,
		method:        "POST",
		mock: func() {
			store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
			clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(nil, errors.New(""))
		},
	}, {
		description: "fails because scope not allowed",
		form: url.Values{
			"client_id": {"client_id"},
			"scope":     {"17 42 foo"},
		},
		method: "POST",
		mock: func() {
			store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
			clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(deviceClient, nil)
		},
		expectedError: ErrInvalidScope,
	}, {
		description: "fails because audience not allowed",
		form: url.Values{
			"client_id": {"client_id"},
			"scope":     {"17 42"},
			"audience":  {"random_aud"},
		},
		method: "POST",
		mock: func() {
			store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
			clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(deviceClient, nil)
		},
		expectedError: ErrInvalidRequest,
	}, {
		description: "fails because it doesn't have the proper grant",
		form: url.Values{
			"client_id": {"client_id_2"},
			"scope":     {"17 42"},
		},
		method: "POST",
		mock: func() {
			store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
			clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id_2")).Return(authCodeClient, nil)
		},
		expectedError: ErrInvalidGrant,
	}, {
		description: "success",
		form: url.Values{
			"client_id": {"client_id"},
			"scope":     {"17 42"},
		},
		method: "POST",
		mock: func() {
			store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
			clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(deviceClient, nil)
		},
	}} {
		t.Run(fmt.Sprintf("case=%d description=%s", k, c.description), func(t *testing.T) {
			c.mock()
			r := &http.Request{
				Header:   c.header,
				PostForm: c.form,
				Form:     c.form,
				Method:   c.method,
			}

			ar, err := fosite.NewDeviceRequest(context.Background(), r)
			require.ErrorIs(t, err, c.expectedError)
			if c.expectedError == nil {
				assert.NotNil(t, ar.GetRequestedAt())
			}
		})
	}
}

func TestNewDeviceRequestWithClientAuthn(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockStorage(ctrl)
	clientManager := internal.NewMockClientManager(ctrl)
	hasher := internal.NewMockHasher(ctrl)
	client := &DefaultClient{ID: "client_id"}
	t.Cleanup(ctrl.Finish)
	config := &Config{ClientSecretsHasher: hasher, ScopeStrategy: ExactScopeStrategy, AudienceMatchingStrategy: DefaultAudienceMatchingStrategy}
	fosite := &Fosite{Store: store, Config: config}

	client.Public = false
	client.Secret = []byte("client_secret")
	client.Scopes = []string{"foo", "bar"}
	client.GrantTypes = []string{"urn:ietf:params:oauth:grant-type:device_code"}

	for k, c := range []struct {
		header        http.Header
		form          url.Values
		method        string
		expectedError error
		mock          func()
		expect        DeviceRequester
		description   string
	}{
		{
			form: url.Values{
				"client_id": {"client_id"},
				"scope":     {"foo bar"},
			},
			expectedError: ErrInvalidClient,
			method:        "POST",
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(client, nil)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))
			},
			description: "Should failed becaue no client authn provided.",
		},
		{
			form: url.Values{
				"client_id": {"client_id2"},
				"scope":     {"foo bar"},
			},
			header: http.Header{
				"Authorization": {basicAuth("client_id", "client_secret")},
			},
			expectedError: ErrInvalidRequest,
			method:        "POST",
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(client, nil)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("client_secret")), gomock.Eq([]byte("client_secret"))).Return(nil)
			},
			description: "should fail because different client is used in authn than in form",
		},
		{
			form: url.Values{
				"client_id": {"client_id"},
				"scope":     {"foo bar"},
			},
			header: http.Header{
				"Authorization": {basicAuth("client_id", "client_secret")},
			},
			method: "POST",
			mock: func() {
				store.EXPECT().FositeClientManager().Return(clientManager).Times(1)
				clientManager.EXPECT().GetClient(gomock.Any(), gomock.Eq("client_id")).Return(client, nil)
				hasher.EXPECT().Compare(gomock.Any(), gomock.Eq([]byte("client_secret")), gomock.Eq([]byte("client_secret"))).Return(nil)
			},
			description: "should succeed",
		},
	} {
		t.Run(fmt.Sprintf("case=%d description=%s", k, c.description), func(t *testing.T) {
			c.mock()
			r := &http.Request{
				Header:   c.header,
				PostForm: c.form,
				Form:     c.form,
				Method:   c.method,
			}

			req, err := fosite.NewDeviceRequest(context.Background(), r)
			require.ErrorIs(t, err, c.expectedError)
			if c.expectedError == nil {
				assert.NotZero(t, req.GetRequestedAt())
			}
		})
	}
}
