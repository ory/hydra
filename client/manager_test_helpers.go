/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package client

import (
	"context"
	"crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"

	"github.com/ory/hydra/x"
)

func TestHelperClientAutoGenerateKey(k string, m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()
		c := &Client{
			OutfacingID:       "foo",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}
		assert.NoError(t, m.CreateClient(ctx, c))
		// assert.NotEmpty(t, c.ID)
		assert.NoError(t, m.DeleteClient(ctx, c.GetID()))
	}
}

func TestHelperClientAuthenticate(k string, m Manager) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()
		require.NoError(t, m.CreateClient(ctx, &Client{
			OutfacingID:  "1234321",
			Secret:       "secret",
			RedirectURIs: []string{"http://redirect"},
		}))

		c, err := m.Authenticate(ctx, "1234321", []byte("secret1"))
		require.NotNil(t, err)

		c, err = m.Authenticate(ctx, "1234321", []byte("secret"))
		require.NoError(t, err)
		assert.Equal(t, "1234321", c.GetID())
	}
}

func TestHelperUpdateTwoClients(_ string, m Manager) func(t *testing.T) {
	return func(t *testing.T) {
		c1, c2 := &Client{OutfacingID: "klojdfc", Name: "test client 1"}, &Client{OutfacingID: "jlsdfkj", Name: "test client 2"}

		require.NoError(t, m.CreateClient(context.Background(), c1))
		require.NoError(t, m.CreateClient(context.Background(), c2))

		c1.Name, c2.Name = "updated klojdfc client 1", "updated klojdfc client 2"

		assert.NoError(t, m.UpdateClient(context.Background(), c1))
		assert.NoError(t, m.UpdateClient(context.Background(), c2))
	}
}

func TestHelperCreateGetUpdateDeleteClient(k string, m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetClient(ctx, "1234")
		require.Error(t, err)

		c := &Client{
			OutfacingID:                       "1234",
			Name:                              "name",
			Secret:                            "secret",
			RedirectURIs:                      []string{"http://redirect", "http://redirect1"},
			GrantTypes:                        []string{"implicit", "refresh_token"},
			ResponseTypes:                     []string{"code token", "token id_token", "code"},
			Scope:                             "scope-a scope-b",
			Owner:                             "aeneas",
			PolicyURI:                         "http://policy",
			TermsOfServiceURI:                 "http://tos",
			ClientURI:                         "http://client",
			LogoURI:                           "http://logo",
			Contacts:                          []string{"aeneas1", "aeneas2"},
			SecretExpiresAt:                   0,
			SectorIdentifierURI:               "https://sector",
			JSONWebKeys:                       &x.JoseJSONWebKeySet{JSONWebKeySet: &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{KeyID: "foo", Key: []byte("asdf"), Certificates: []*x509.Certificate{}, CertificateThumbprintSHA1: []uint8{}, CertificateThumbprintSHA256: []uint8{}}}}},
			JSONWebKeysURI:                    "https://...",
			TokenEndpointAuthMethod:           "none",
			TokenEndpointAuthSigningAlgorithm: "RS256",
			RequestURIs:                       []string{"foo", "bar"},
			AllowedCORSOrigins:                []string{"foo", "bar"},
			RequestObjectSigningAlgorithm:     "rs256",
			UserinfoSignedResponseAlg:         "RS256",
			CreatedAt:                         time.Now().Add(-time.Hour).Round(time.Second).UTC(),
			UpdatedAt:                         time.Now().Add(-time.Minute).Round(time.Second).UTC(),
			FrontChannelLogoutURI:             "http://fc-logout",
			FrontChannelLogoutSessionRequired: true,
			PostLogoutRedirectURIs:            []string{"hello", "mister"},
			BackChannelLogoutURI:              "http://bc-logout",
			BackChannelLogoutSessionRequired:  true,
		}

		require.NoError(t, m.CreateClient(ctx, c))
		assert.Equal(t, c.GetID(), "1234")
		if k != "http" {
			assert.NotEmpty(t, c.GetHashedSecret())
		}

		assert.NoError(t, m.CreateClient(ctx, &Client{
			OutfacingID:       "2-1234",
			Name:              "name2",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
			SecretExpiresAt:   1,
		}))

		d, err := m.GetClient(ctx, "1234")
		require.NoError(t, err)

		compare(t, c, d, k)

		ds, err := m.GetClients(ctx, Filter{Limit: 100, Offset: 0})
		assert.NoError(t, err)
		assert.Len(t, ds, 2)
		assert.NotEqual(t, ds[0].OutfacingID, ds[1].OutfacingID)
		assert.NotEqual(t, ds[0].OutfacingID, ds[1].OutfacingID)
		// test if SecretExpiresAt was set properly
		assert.Equal(t, ds[0].SecretExpiresAt, 0)
		assert.Equal(t, ds[1].SecretExpiresAt, 1)

		ds, err = m.GetClients(ctx, Filter{Limit: 1, Offset: 0})
		assert.NoError(t, err)
		assert.Len(t, ds, 1)

		ds, err = m.GetClients(ctx, Filter{Limit: 100, Offset: 100})
		assert.NoError(t, err)

		// get by name
		ds, err = m.GetClients(ctx, Filter{Limit: 100, Offset: 0, Name: "name"})
		assert.NoError(t, err)
		assert.Len(t, ds, 1)
		assert.Equal(t, ds[0].Name, "name")

		// get by name not exist
		ds, err = m.GetClients(ctx, Filter{Limit: 100, Offset: 0, Name: "bad name"})
		assert.NoError(t, err)
		assert.Len(t, ds, 0)

		// get by owner
		ds, err = m.GetClients(ctx, Filter{Limit: 100, Offset: 0, Owner: "aeneas"})
		assert.NoError(t, err)
		assert.Len(t, ds, 1)
		assert.Equal(t, ds[0].Owner, "aeneas")

		err = m.UpdateClient(ctx, &Client{
			OutfacingID:       "2-1234",
			Name:              "name-new",
			Secret:            "secret-new",
			RedirectURIs:      []string{"http://redirect/new"},
			TermsOfServiceURI: "bar",
			JSONWebKeys:       new(x.JoseJSONWebKeySet),
		})
		require.NoError(t, err)

		nc, err := m.GetConcreteClient(ctx, "2-1234")
		require.NoError(t, err)

		if k != "http" {
			// http always returns an empty secret
			assert.NotEqual(t, d.GetHashedSecret(), nc.GetHashedSecret())
		}
		assert.Equal(t, "bar", nc.TermsOfServiceURI)
		assert.Equal(t, "name-new", nc.Name)
		assert.EqualValues(t, []string{"http://redirect/new"}, nc.GetRedirectURIs())
		assert.Zero(t, len(nc.Contacts))

		err = m.DeleteClient(ctx, "1234")
		assert.NoError(t, err)

		_, err = m.GetClient(ctx, "1234")
		assert.NotNil(t, err)

		n, err := m.CountClients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, n)
	}
}

func compare(t *testing.T, expected *Client, actual fosite.Client, k string) {
	assert.EqualValues(t, expected.GetID(), actual.GetID())
	if k != "http" {
		assert.EqualValues(t, expected.GetHashedSecret(), actual.GetHashedSecret())
	}
	assert.EqualValues(t, expected.GetRedirectURIs(), actual.GetRedirectURIs())
	assert.EqualValues(t, expected.GetGrantTypes(), actual.GetGrantTypes())

	assert.EqualValues(t, expected.GetResponseTypes(), actual.GetResponseTypes())
	assert.EqualValues(t, expected.GetScopes(), actual.GetScopes())
	assert.EqualValues(t, expected.IsPublic(), actual.IsPublic())

	if actual, ok := actual.(*Client); ok {
		assert.EqualValues(t, expected.Owner, actual.Owner)
		assert.EqualValues(t, expected.Name, actual.Name)
		assert.EqualValues(t, expected.PolicyURI, actual.PolicyURI)
		assert.EqualValues(t, expected.TermsOfServiceURI, actual.TermsOfServiceURI)
		assert.EqualValues(t, expected.ClientURI, actual.ClientURI)
		assert.EqualValues(t, expected.LogoURI, actual.LogoURI)
		assert.EqualValues(t, expected.Contacts, actual.Contacts)
		assert.EqualValues(t, expected.SecretExpiresAt, actual.SecretExpiresAt)
		assert.EqualValues(t, expected.SectorIdentifierURI, actual.SectorIdentifierURI)
		assert.EqualValues(t, expected.UserinfoSignedResponseAlg, actual.UserinfoSignedResponseAlg)
		assert.EqualValues(t, expected.CreatedAt.UTC().Unix(), actual.CreatedAt.UTC().Unix())
		// these values are not the same because of https://github.com/gobuffalo/pop/issues/591
		//assert.EqualValues(t, expected.UpdatedAt.UTC().Unix(), actual.UpdatedAt.UTC().Unix(), "%s\n%s", expected.UpdatedAt.String(), actual.UpdatedAt.String())
		assert.EqualValues(t, expected.FrontChannelLogoutURI, actual.FrontChannelLogoutURI)
		assert.EqualValues(t, expected.FrontChannelLogoutSessionRequired, actual.FrontChannelLogoutSessionRequired)
		assert.EqualValues(t, expected.PostLogoutRedirectURIs, actual.PostLogoutRedirectURIs)
		assert.EqualValues(t, expected.BackChannelLogoutURI, actual.BackChannelLogoutURI)
		assert.EqualValues(t, expected.BackChannelLogoutSessionRequired, actual.BackChannelLogoutSessionRequired)
	}

	if actual, ok := actual.(fosite.OpenIDConnectClient); ok {
		require.NotNil(t, expected.JSONWebKeys)

		for k, v := range expected.JSONWebKeys.JSONWebKeySet.Keys {
			if v.CertificateThumbprintSHA1 == nil {
				v.CertificateThumbprintSHA1 = make([]byte, 0)
			}
			if v.CertificateThumbprintSHA256 == nil {
				v.CertificateThumbprintSHA256 = make([]byte, 0)
			}
			expected.JSONWebKeys.JSONWebKeySet.Keys[k] = v
		}

		assert.EqualValues(t, expected.JSONWebKeys.JSONWebKeySet, actual.GetJSONWebKeys())
		assert.EqualValues(t, expected.JSONWebKeysURI, actual.GetJSONWebKeysURI())
		assert.EqualValues(t, expected.TokenEndpointAuthMethod, actual.GetTokenEndpointAuthMethod())
		assert.EqualValues(t, expected.RequestURIs, actual.GetRequestURIs())
		assert.EqualValues(t, expected.RequestObjectSigningAlgorithm, actual.GetRequestObjectSigningAlgorithm())
	}
}
