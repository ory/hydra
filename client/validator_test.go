// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/x/httpx"

	"github.com/gofrs/uuid"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()
	c := internal.NewConfigurationWithDefaults()
	c.MustSet(ctx, config.KeySubjectTypesSupported, []string{"pairwise", "public"})
	c.MustSet(ctx, config.KeyDefaultClientScope, []string{"openid"})
	reg := internal.NewRegistryMemory(t, c, &contextx.Static{C: c.Source(ctx)})
	v := NewValidator(reg)

	testCtx := context.TODO()

	for k, tc := range []struct {
		in        *Client
		check     func(t *testing.T, c *Client)
		expectErr bool
		v         func(t *testing.T) *Validator
	}{
		{
			in: new(Client),
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, uuid.Nil.String(), c.GetID())
				assert.EqualValues(t, c.GetID(), c.ID.String())
				assert.Empty(t, c.LegacyClientID)
			},
		},
		{
			in: &Client{LegacyClientID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.EqualValues(t, c.GetID(), c.LegacyClientID)
			},
		},
		{
			in: &Client{LegacyClientID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.EqualValues(t, c.GetID(), c.LegacyClientID)
			},
		},
		{
			in:        &Client{LegacyClientID: "foo", UserinfoSignedResponseAlg: "foo"},
			expectErr: true,
		},
		{
			in:        &Client{LegacyClientID: "foo", TokenEndpointAuthMethod: "private_key_jwt"},
			expectErr: true,
		},
		{
			in:        &Client{LegacyClientID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: new(jose.JSONWebKeySet)}, JSONWebKeysURI: "asdf", TokenEndpointAuthMethod: "private_key_jwt"},
			expectErr: true,
		},
		{
			in:        &Client{LegacyClientID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: new(jose.JSONWebKeySet)}, TokenEndpointAuthMethod: "private_key_jwt", TokenEndpointAuthSigningAlgorithm: "HS256"},
			expectErr: true,
		},
		{
			in:        &Client{LegacyClientID: "foo", PostLogoutRedirectURIs: []string{"https://bar/"}, RedirectURIs: []string{"https://foo/"}},
			expectErr: true,
		},
		{
			in:        &Client{LegacyClientID: "foo", PostLogoutRedirectURIs: []string{"http://foo/"}, RedirectURIs: []string{"https://foo/"}},
			expectErr: true,
		},
		{
			in:        &Client{LegacyClientID: "foo", PostLogoutRedirectURIs: []string{"https://foo:1234/"}, RedirectURIs: []string{"https://foo/"}},
			expectErr: true,
		},
		{
			in: &Client{LegacyClientID: "foo", PostLogoutRedirectURIs: []string{"https://foo/"}, RedirectURIs: []string{"https://foo/"}},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, []string{"https://foo/"}, []string(c.PostLogoutRedirectURIs))
			},
		},
		{
			in: &Client{LegacyClientID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "public", c.SubjectType)
			},
		},
		{
			v: func(t *testing.T) *Validator {
				c.MustSet(ctx, config.KeySubjectTypesSupported, []string{"pairwise"})
				return NewValidator(reg)
			},
			in: &Client{LegacyClientID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "pairwise", c.SubjectType)
			},
		},
		{
			in: &Client{LegacyClientID: "foo", SubjectType: "pairwise"},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "pairwise", c.SubjectType)
			},
		},
		{
			in:        &Client{LegacyClientID: "foo", SubjectType: "foo"},
			expectErr: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			if tc.v == nil {
				tc.v = func(t *testing.T) *Validator {
					return v
				}
			}
			err := tc.v(t).Validate(testCtx, tc.in)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tc.check(t, tc.in)
			}
		})
	}
}

type fakeHTTP struct {
	driver.Registry
	c *http.Client
}

func (f *fakeHTTP) HTTPClient(ctx context.Context, opts ...httpx.ResilientOptions) *retryablehttp.Client {
	return httpx.NewResilientClient(httpx.ResilientClientWithClient(f.c))
}

func TestValidateSectorIdentifierURL(t *testing.T) {
	reg := internal.NewMockedRegistry(t, &contextx.Default{})
	var payload string

	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(payload))
	}
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	v := NewValidator(&fakeHTTP{Registry: reg, c: ts.Client()})
	for k, tc := range []struct {
		p         string
		r         []string
		u         string
		expectErr bool
	}{
		{
			u:         "",
			expectErr: true,
		},
		{
			u:         "http://foo/bar",
			expectErr: true,
		},
		{
			u:         ts.URL,
			expectErr: true,
		},
		{
			p:         `["http://foo"]`,
			u:         ts.URL,
			expectErr: false,
			r:         []string{"http://foo"},
		},
		{
			p:         `["http://foo"]`,
			u:         ts.URL,
			expectErr: true,
			r:         []string{"http://foo", "http://not-foo"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			payload = tc.p
			err := v.ValidateSectorIdentifierURL(context.Background(), tc.u, tc.r)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateIPRanges(t *testing.T) {
	ctx := context.Background()
	c := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, c, &contextx.Static{C: c.Source(ctx)})

	v := NewValidator(reg)
	c.MustSet(ctx, config.ViperKeyClientHTTPNoPrivateIPRanges, true)
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{}))
	require.ErrorContains(t, v.ValidateDynamicRegistration(ctx, &Client{JSONWebKeysURI: "https://localhost:1234"}), "invalid_client_metadata")
	require.ErrorContains(t, v.ValidateDynamicRegistration(ctx, &Client{BackChannelLogoutURI: "https://localhost:1234"}), "invalid_client_metadata")
	require.ErrorContains(t, v.ValidateDynamicRegistration(ctx, &Client{RequestURIs: []string{"https://google", "https://localhost:1234"}}), "invalid_client_metadata")

	c.MustSet(ctx, config.ViperKeyClientHTTPNoPrivateIPRanges, false)
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{}))
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{JSONWebKeysURI: "https://localhost:1234"}))
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{BackChannelLogoutURI: "https://localhost:1234"}))
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{RequestURIs: []string{"https://google", "https://localhost:1234"}}))
}

func TestValidateDynamicRegistration(t *testing.T) {
	ctx := context.Background()
	c := internal.NewConfigurationWithDefaults()
	c.MustSet(ctx, config.KeySubjectTypesSupported, []string{"pairwise", "public"})
	c.MustSet(ctx, config.KeyDefaultClientScope, []string{"openid"})
	reg := internal.NewRegistryMemory(t, c, &contextx.Static{C: c.Source(ctx)})

	testCtx := context.TODO()
	v := NewValidator(reg)
	for k, tc := range []struct {
		in        *Client
		check     func(t *testing.T, c *Client)
		expectErr bool
		v         func(t *testing.T) *Validator
	}{
		{
			in: &Client{
				LegacyClientID:         "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
				Metadata:               []byte("{\"access_token_ttl\":10}"),
			},
			expectErr: true,
		},
		{
			in: &Client{
				LegacyClientID:         "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
				Metadata:               []byte("{\"id_token_ttl\":10}"),
			},
			expectErr: true,
		},
		{
			in: &Client{
				LegacyClientID:         "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
				Metadata:               []byte("{\"anything\":10}"),
			},
			expectErr: true,
		},
		{
			in: &Client{
				LegacyClientID:         "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
			},
			check: func(t *testing.T, c *Client) {
				assert.EqualValues(t, "foo", c.LegacyClientID)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			if tc.v == nil {
				tc.v = func(t *testing.T) *Validator {
					return v
				}
			}
			err := tc.v(t).ValidateDynamicRegistration(testCtx, tc.in)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tc.check(t, tc.in)
			}
		})
	}
}
