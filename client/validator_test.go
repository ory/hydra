// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/httpx"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeySubjectTypesSupported: []string{"pairwise", "public"},
		config.KeyDefaultClientScope:    []string{"openid"},
	})))
	v := NewValidator(reg)

	dec := json.NewDecoder(strings.NewReader(validJWKS))
	dec.DisallowUnknownFields()
	var goodJWKS jose.JSONWebKeySet
	require.NoError(t, dec.Decode(&goodJWKS))

	for k, tc := range []struct {
		in        *Client
		check     func(*testing.T, *Client)
		assertErr func(t assert.TestingT, err error, msg ...interface{}) bool
		v         func(*testing.T) *Validator
	}{
		{
			in: new(Client),
			check: func(t *testing.T, c *Client) {
				assert.Zero(t, c.GetID())
				assert.EqualValues(t, c.GetID(), c.ID)
			},
		},
		{
			in: &Client{ID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.EqualValues(t, c.GetID(), c.ID)
			},
		},
		{
			in: &Client{ID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.EqualValues(t, c.GetID(), c.ID)
			},
		},
		{
			in:        &Client{ID: "foo", UserinfoSignedResponseAlg: "foo"},
			assertErr: assert.Error,
		},
		{
			in:        &Client{ID: "foo", TokenEndpointAuthMethod: "private_key_jwt"},
			assertErr: assert.Error,
		},
		{
			in:        &Client{ID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: new(jose.JSONWebKeySet)}, JSONWebKeysURI: "asdf", TokenEndpointAuthMethod: "private_key_jwt"},
			assertErr: assert.Error,
		},
		{
			in:        &Client{ID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: new(jose.JSONWebKeySet)}, TokenEndpointAuthMethod: "private_key_jwt", TokenEndpointAuthSigningAlgorithm: "HS256"},
			assertErr: assert.Error,
		},
		{
			in: &Client{ID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: new(jose.JSONWebKeySet)}, JSONWebKeysURI: "https://example.org/jwks.json"},
			assertErr: func(t assert.TestingT, err error, msg ...interface{}) bool {
				e := new(fosite.RFC6749Error)
				assert.ErrorAs(t, err, &e)
				assert.Contains(t, e.HintField, "jwks and jwks_uri can not both be set")
				return true
			},
		},
		{
			in: &Client{ID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: &goodJWKS}},
			check: func(t *testing.T, c *Client) {
				assert.Len(t, c.JSONWebKeys.Keys, 2)
				assert.Equal(t, c.JSONWebKeys.Keys[0].KeyID, "1")
				assert.Equal(t, c.JSONWebKeys.Keys[1].KeyID, "2011-04-29")
			},
		},
		{
			in: &Client{ID: "foo", JSONWebKeys: &x.JoseJSONWebKeySet{JSONWebKeySet: &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{}}}}},
			assertErr: func(t assert.TestingT, err error, msg ...interface{}) bool {
				e := new(fosite.RFC6749Error)
				assert.ErrorAs(t, err, &e)
				assert.Contains(t, e.HintField, "Invalid JSON web key in set")
				return true
			},
		},
		{
			in: &Client{ID: "foo", JSONWebKeys: new(x.JoseJSONWebKeySet), JSONWebKeysURI: "https://example.org/jwks.json"},
			check: func(t *testing.T, c *Client) {
				assert.Nil(t, c.GetJSONWebKeys())
			},
		},
		{
			in:        &Client{ID: "foo", PostLogoutRedirectURIs: []string{"https://bar/"}, RedirectURIs: []string{"https://foo/"}},
			assertErr: assert.Error,
		},
		{
			in:        &Client{ID: "foo", PostLogoutRedirectURIs: []string{"http://foo/"}, RedirectURIs: []string{"https://foo/"}},
			assertErr: assert.Error,
		},
		{
			in:        &Client{ID: "foo", PostLogoutRedirectURIs: []string{"https://foo:1234/"}, RedirectURIs: []string{"https://foo/"}},
			assertErr: assert.Error,
		},
		{
			in: &Client{ID: "foo", PostLogoutRedirectURIs: []string{"https://foo/"}, RedirectURIs: []string{"https://foo/"}},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, []string{"https://foo/"}, []string(c.PostLogoutRedirectURIs))
			},
		},
		{
			in:        &Client{ID: "foo", TermsOfServiceURI: "https://example.org"},
			assertErr: assert.NoError,
		},
		{
			in:        &Client{ID: "foo", TermsOfServiceURI: "javascript:alert('XSS')"},
			assertErr: assert.Error,
		},
		{
			in: &Client{ID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "public", c.SubjectType)
			},
		},
		{
			v: func(t *testing.T) *Validator {
				reg.Config().MustSet(ctx, config.KeySubjectTypesSupported, []string{"pairwise"})
				return NewValidator(reg)
			},
			in: &Client{ID: "foo"},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "pairwise", c.SubjectType)
			},
		},
		{
			in: &Client{ID: "foo", SubjectType: "pairwise"},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "pairwise", c.SubjectType)
			},
		},
		{
			in:        &Client{ID: "foo", SubjectType: "foo"},
			assertErr: assert.Error,
		},
	} {
		tc := tc
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			if tc.v == nil {
				tc.v = func(t *testing.T) *Validator {
					return v
				}
			}
			err := tc.v(t).Validate(ctx, tc.in)
			if tc.assertErr != nil {
				tc.assertErr(t, err)
			} else {
				require.NoError(t, err)
				tc.check(t, tc.in)
			}
		})
	}
}

type fakeHTTP struct {
	*driver.RegistrySQL
	c *http.Client
}

func (f *fakeHTTP) HTTPClient(_ context.Context, opts ...httpx.ResilientOptions) *retryablehttp.Client {
	c := httpx.NewResilientClient(opts...)
	c.HTTPClient = f.c
	return c
}

func TestValidateSectorIdentifierURL(t *testing.T) {
	reg := testhelpers.NewRegistryMemory(t)
	var payload string

	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(payload))
	}
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	v := NewValidator(&fakeHTTP{RegistrySQL: reg, c: ts.Client()})
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

// from https://datatracker.ietf.org/doc/html/rfc7517#appendix-A.2
const validJWKS = `
{"keys":
[
  {"kty":"EC",
   "crv":"P-256",
   "x":"MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
   "y":"4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
   "d":"870MB6gfuTJ4HtUnUvYMyJpr5eUZNP4Bk43bVdj3eAE",
   "use":"enc",
   "kid":"1"},

  {"kty":"RSA",
   "n":"0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw",
   "e":"AQAB",
   "d":"X4cTteJY_gn4FYPsXB8rdXix5vwsg1FLN5E3EaG6RJoVH-HLLKD9M7dx5oo7GURknchnrRweUkC7hT5fJLM0WbFAKNLWY2vv7B6NqXSzUvxT0_YSfqijwp3RTzlBaCxWp4doFk5N2o8Gy_nHNKroADIkJ46pRUohsXywbReAdYaMwFs9tv8d_cPVY3i07a3t8MN6TNwm0dSawm9v47UiCl3Sk5ZiG7xojPLu4sbg1U2jx4IBTNBznbJSzFHK66jT8bgkuqsk0GjskDJk19Z4qwjwbsnn4j2WBii3RL-Us2lGVkY8fkFzme1z0HbIkfz0Y6mqnOYtqc0X4jfcKoAC8Q",
   "p":"83i-7IvMGXoMXCskv73TKr8637FiO7Z27zv8oj6pbWUQyLPQBQxtPVnwD20R-60eTDmD2ujnMt5PoqMrm8RfmNhVWDtjjMmCMjOpSXicFHj7XOuVIYQyqVWlWEh6dN36GVZYk93N8Bc9vY41xy8B9RzzOGVQzXvNEvn7O0nVbfs",
   "q":"3dfOR9cuYq-0S-mkFLzgItgMEfFzB2q3hWehMuG0oCuqnb3vobLyumqjVZQO1dIrdwgTnCdpYzBcOfW5r370AFXjiWft_NGEiovonizhKpo9VVS78TzFgxkIdrecRezsZ-1kYd_s1qDbxtkDEgfAITAG9LUnADun4vIcb6yelxk",
   "dp":"G4sPXkc6Ya9y8oJW9_ILj4xuppu0lzi_H7VTkS8xj5SdX3coE0oimYwxIi2emTAue0UOa5dpgFGyBJ4c8tQ2VF402XRugKDTP8akYhFo5tAA77Qe_NmtuYZc3C3m3I24G2GvR5sSDxUyAN2zq8Lfn9EUms6rY3Ob8YeiKkTiBj0",
   "dq":"s9lAH9fggBsoFR8Oac2R_E2gw282rT2kGOAhvIllETE1efrA6huUUvMfBcMpn8lqeW6vzznYY5SSQF7pMdC_agI3nG8Ibp1BUb0JUiraRNqUfLhcQb_d9GF4Dh7e74WbRsobRonujTYN1xCaP6TO61jvWrX-L18txXw494Q_cgk",
   "qi":"GyM_p6JrXySiz1toFgKbWV-JdI3jQ4ypu9rbMWx3rQJBfmt0FoYzgUIZEVFEcOqwemRN81zoDAaa-Bk0KWNGDjJHZDdDmFhW3AN7lI-puxk_mHZGJ11rxyR8O55XLSe3SPmRfKwZI6yU24ZxvQKFYItdldUKGzO6Ia6zTKhAVRU",
   "alg":"RS256",
   "kid":"2011-04-29"}
]
}
`

func TestValidateIPRanges(t *testing.T) {
	ctx := t.Context()
	reg := testhelpers.NewRegistryMemory(t)

	v := NewValidator(reg)
	reg.Config().MustSet(t.Context(), config.KeyClientHTTPNoPrivateIPRanges, true)
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{}))
	require.ErrorContains(t, v.ValidateDynamicRegistration(ctx, &Client{JSONWebKeysURI: "https://localhost:1234"}), "invalid_client_metadata")
	require.ErrorContains(t, v.ValidateDynamicRegistration(ctx, &Client{BackChannelLogoutURI: "https://localhost:1234"}), "invalid_client_metadata")
	require.ErrorContains(t, v.ValidateDynamicRegistration(ctx, &Client{RequestURIs: []string{"https://google", "https://localhost:1234"}}), "invalid_client_metadata")

	reg.Config().MustSet(t.Context(), config.KeyClientHTTPNoPrivateIPRanges, false)
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{}))
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{JSONWebKeysURI: "https://localhost:1234"}))
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{BackChannelLogoutURI: "https://localhost:1234"}))
	require.NoError(t, v.ValidateDynamicRegistration(ctx, &Client{RequestURIs: []string{"https://google", "https://localhost:1234"}}))
}

func TestValidateDynamicRegistration(t *testing.T) {
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(map[string]any{
		config.KeySubjectTypesSupported: []string{"pairwise", "public"},
		config.KeyDefaultClientScope:    []string{"openid"},
	})))

	v := NewValidator(reg)
	for k, tc := range []struct {
		in        *Client
		check     func(t *testing.T, c *Client)
		expectErr bool
		v         func(t *testing.T) *Validator
	}{
		{
			in: &Client{
				ID:                     "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
				Metadata:               []byte(`{"access_token_ttl":10}`),
			},
			expectErr: true,
		},
		{
			in: &Client{
				ID:                     "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
				Metadata:               []byte(`{"id_token_ttl":10}`),
			},
			expectErr: true,
		},
		{
			in: &Client{
				ID:                     "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
				Metadata:               []byte(`{"anything":10}`),
			},
			expectErr: true,
		},
		{
			in: &Client{
				ID:                     "foo",
				PostLogoutRedirectURIs: []string{"https://foo/"},
				RedirectURIs:           []string{"https://foo/"},
			},
			check: func(t *testing.T, c *Client) {
				assert.EqualValues(t, "foo", c.ID)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			if tc.v == nil {
				tc.v = func(t *testing.T) *Validator {
					return v
				}
			}
			err := tc.v(t).ValidateDynamicRegistration(t.Context(), tc.in)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tc.check(t, tc.in)
			}
		})
	}
}
