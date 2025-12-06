// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/uuidx"

	"github.com/ory/hydra/v2/fosite"
	testhelpersuuid "github.com/ory/hydra/v2/internal/testhelpers/uuid"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/assertx"
	"github.com/ory/x/contextx"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
	"github.com/ory/x/sqlcon"
)

func TestHelperClientAutoGenerateKey(m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()
		c := &Client{
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}
		require.NoError(t, m.CreateClient(ctx, c))
		dbClient, err := m.GetClient(ctx, c.GetID())
		require.NoError(t, err)
		dbClientConcrete, ok := dbClient.(*Client)
		require.True(t, ok)
		testhelpersuuid.AssertUUID(t, dbClientConcrete.ID)
		assert.NoError(t, m.DeleteClient(ctx, c.GetID()))
	}
}

func TestHelperClientAuthenticate(m Manager) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()
		require.NoError(t, m.CreateClient(ctx, &Client{
			ID:           "1234321",
			Secret:       "secret",
			RedirectURIs: []string{"http://redirect"},
		}))

		c, err := m.AuthenticateClient(ctx, "1234321", []byte("secret1"))
		require.Error(t, err)
		require.Nil(t, c)

		c, err = m.AuthenticateClient(ctx, "1234321", []byte("secret"))
		require.NoError(t, err)
		assert.Equal(t, "1234321", c.GetID())
	}
}

func TestHelperUpdateTwoClients(m Manager) func(t *testing.T) {
	return func(t *testing.T) {
		c1, c2 := &Client{Name: "test client 1"}, &Client{Name: "test client 2"}

		require.NoError(t, m.CreateClient(context.Background(), c1))
		require.NoError(t, m.CreateClient(context.Background(), c2))

		c1.Name, c2.Name = "updated client 1", "updated client 2"

		assert.NoError(t, m.UpdateClient(context.Background(), c1))
		assert.NoError(t, m.UpdateClient(context.Background(), c2))
	}
}

func testHelperUpdateClient(t *testing.T, ctx context.Context, store Storage, toUpdate *Client) {
	require.NoError(t, store.UpdateClient(ctx, &Client{
		ID:                toUpdate.ID,
		Name:              "name-new",
		Secret:            "secret-new",
		RedirectURIs:      []string{"http://redirect/new"},
		TermsOfServiceURI: "bar",
		JSONWebKeys:       new(x.JoseJSONWebKeySet),
	}))

	nc, err := store.GetConcreteClient(ctx, toUpdate.ID)
	require.NoError(t, err)

	require.NotZero(t, toUpdate.GetHashedSecret())
	require.NotZero(t, nc.GetHashedSecret())
	assert.NotEqual(t, toUpdate.GetHashedSecret(), nc.GetHashedSecret(), "ensure that the secret hash was updated")

	assert.Equal(t, "bar", nc.TermsOfServiceURI)
	assert.Equal(t, "name-new", nc.Name)
	assert.Equal(t, []string{"http://redirect/new"}, nc.GetRedirectURIs())
	assert.Len(t, nc.Contacts, 0)
}

func TestHelperCreateGetUpdateDeleteClientNext(t *testing.T, m Storage, networks []uuid.UUID) {
	ctx := context.Background()

	resources := map[uuid.UUID][]Client{}
	for k := range networks {
		nid := networks[k]
		resources[nid] = []Client{}

		ctx := contextx.SetNIDContext(ctx, nid)
		t.Run(fmt.Sprintf("nid=%s", nid), func(t *testing.T) {
			var client Client
			require.NoError(t, faker.FakeData(&client))
			client.CreatedAt = time.Now().Truncate(time.Second).UTC()

			t.Run("lifecycle=does not exist", func(t *testing.T) {
				_, err := m.GetClient(ctx, uuidx.NewV4().String())
				require.Error(t, err)
			})

			t.Run("lifecycle=exists", func(t *testing.T) {
				require.NoError(t, m.CreateClient(ctx, &client))
				c, err := m.GetClient(ctx, client.GetID())
				require.NoError(t, err)
				assertx.EqualAsJSONExcept(t, &client, c, []string{
					"registration_access_token",
					"registration_client_uri",
					"updated_at",
				})

				require.ErrorIs(t, m.CreateClient(ctx, &client), sqlcon.ErrUniqueViolation)
			})

			t.Run("lifecycle=update", func(t *testing.T) {
				client.Name = "updated" + nid.String()
				require.NoError(t, m.UpdateClient(ctx, &client))
				c, err := m.GetClient(ctx, client.GetID())
				require.NoError(t, err)
				assertx.EqualAsJSONExcept(t, &client, c, []string{
					"registration_access_token",
					"registration_client_uri",
					"updated_at",
				})
				resources[nid] = append(resources[nid], client)
			})
		})
	}

	for k := range resources {
		original := k
		clients := resources[k]
		for i := range networks {
			check := networks[i]

			t.Run("network="+original.String(), func(t *testing.T) {
				ctx := contextx.SetNIDContext(ctx, check)
				for _, expected := range clients {
					c, err := m.GetClient(ctx, expected.GetID())
					if check != original {
						t.Run(fmt.Sprintf("case=must not find client %s", expected.GetID()), func(t *testing.T) {
							require.ErrorIs(t, err, sqlcon.ErrNoRows)
						})
					} else {
						t.Run("case=updates must not override each other", func(t *testing.T) {
							require.NoError(t, err)
							assert.Equal(t, "updated"+original.String(), c.(*Client).Name)
						})
					}
				}
			})
		}
	}

	for k := range resources {
		clients := resources[k]
		ctx := contextx.SetNIDContext(ctx, k)
		t.Run("network="+k.String(), func(t *testing.T) {
			for _, client := range clients {
				t.Run("lifecycle=cleanup", func(t *testing.T) {
					assert.NoError(t, m.DeleteClient(ctx, client.GetID()))

					_, err := m.GetClient(ctx, client.GetID())
					assert.ErrorIs(t, err, sqlcon.ErrNoRows)

					assert.ErrorIs(t, m.DeleteClient(ctx, client.GetID()), sqlcon.ErrNoRows)
				})
			}
		})
	}
}

func TestHelperCreateGetUpdateDeleteClient(t1, t2 Storage) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := t1.GetClient(ctx, "1234")
		require.Error(t, err)

		t1c1 := &Client{
			ID:                                uuidx.NewV4().String(),
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

		require.NoError(t, t1.CreateClient(ctx, t1c1))
		assert.NotEmpty(t, t1c1.GetHashedSecret())

		t2c1 := *t1c1
		require.NoError(t, t2.CreateClient(ctx, &t2c1), "we should be able to create a client with the same ID in other network")

		t2c3 := *t1c1
		{
			t2c3.ID = uuidx.NewV4().String()
			require.NoError(t, t2.CreateClient(ctx, &t2c3))
			require.ErrorIs(t, t2.CreateClient(ctx, &t2c3), sqlcon.ErrUniqueViolation)
		}

		c2Template := &Client{
			ID:                uuidx.NewV4().String(),
			Name:              "name2",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
			SecretExpiresAt:   1,
		}
		assert.NoError(t, t1.CreateClient(ctx, c2Template))
		assert.NoError(t, t2.CreateClient(ctx, c2Template))

		d, err := t1.GetClient(ctx, t1c1.ID)
		require.NoError(t, err)

		cc := d.(*Client)
		testhelpersuuid.AssertUUID(t, cc.NID)

		compare(t, t1c1, d)

		t.Run("list all", func(t *testing.T) {
			cs, nextPage, err := t1.GetClients(ctx, Filter{})
			require.NoError(t, err)
			require.Len(t, cs, 2)
			assert.ElementsMatch(t, []string{t1c1.ID, c2Template.ID}, []string{cs[0].GetID(), cs[1].GetID()})
			// test if SecretExpiresAt was set properly
			assert.ElementsMatch(t, []int{0, 1}, []int{cs[0].SecretExpiresAt, cs[1].SecretExpiresAt})
			assert.True(t, nextPage.IsLast())
		})

		t.Run("pagination", func(t *testing.T) {
			observedIDs := make([]string, 2)

			cs, nextPage, err := t1.GetClients(ctx, Filter{PageOpts: []keysetpagination.Option{keysetpagination.WithSize(1)}})
			require.NoError(t, err)
			require.Len(t, cs, 1)
			assert.False(t, nextPage.IsLast())
			observedIDs[0] = cs[0].ID

			cs, nextPage, err = t1.GetClients(ctx, Filter{PageOpts: nextPage.ToOptions()})
			require.NoError(t, err)
			require.Len(t, cs, 1)
			assert.True(t, nextPage.IsLast())
			observedIDs[1] = cs[0].ID

			assert.ElementsMatch(t, []string{t1c1.ID, c2Template.ID}, observedIDs)
		})

		t.Run("list by name", func(t *testing.T) {
			cs, nextPage, err := t1.GetClients(ctx, Filter{Name: "name"})
			require.NoError(t, err)
			require.Len(t, cs, 1)
			assert.Equal(t, cs[0].Name, "name")
			assert.True(t, nextPage.IsLast())
		})

		t.Run("list by unknown name", func(t *testing.T) {
			cs, nextPage, err := t1.GetClients(ctx, Filter{Name: "bad name"})
			require.NoError(t, err)
			assert.Len(t, cs, 0)
			assert.True(t, nextPage.IsLast())
		})

		t.Run("list by owner", func(t *testing.T) {
			cs, nextPage, err := t1.GetClients(ctx, Filter{Owner: "aeneas"})
			require.NoError(t, err)
			require.Len(t, cs, 1)
			assert.Equal(t, cs[0].Owner, "aeneas")
			assert.True(t, nextPage.IsLast())
		})

		t.Run("list by ids", func(t *testing.T) {
			cs, nextPage, err := t1.GetClients(ctx, Filter{IDs: []string{t1c1.ID, c2Template.ID}})
			require.NoError(t, err)
			require.Len(t, cs, 2)
			assert.ElementsMatch(t, []string{t1c1.ID, c2Template.ID}, []string{cs[0].GetID(), cs[1].GetID()})
			assert.True(t, nextPage.IsLast())

			cs, nextPage, err = t1.GetClients(ctx, Filter{IDs: []string{t1c1.ID}})
			require.NoError(t, err)
			require.Len(t, cs, 1)
			assert.Equal(t, t1c1.ID, cs[0].GetID())
			assert.True(t, nextPage.IsLast())

			cs, nextPage, err = t1.GetClients(ctx, Filter{IDs: []string{c2Template.ID}})
			require.NoError(t, err)
			require.Len(t, cs, 1)
			assert.Equal(t, c2Template.ID, cs[0].GetID())
			assert.True(t, nextPage.IsLast())

			cs, nextPage, err = t1.GetClients(ctx, Filter{IDs: []string{uuidx.NewV4().String()}})
			require.NoError(t, err)
			require.Len(t, cs, 0)
			assert.True(t, nextPage.IsLast())
		})

		testHelperUpdateClient(t, ctx, t1, t1c1)
		testHelperUpdateClient(t, ctx, t2, &t2c1)

		assert.ErrorIs(t, t1.DeleteClient(ctx, t2c3.ID), sqlcon.ErrNoRows)

		assert.NoError(t, t1.DeleteClient(ctx, t1c1.ID))
		_, err = t1.GetClient(ctx, t1c1.ID)
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	}
}

func compare(t *testing.T, expected *Client, actual fosite.Client) {
	assert.EqualValues(t, expected.GetID(), actual.GetID())
	assert.EqualValues(t, expected.GetHashedSecret(), actual.GetHashedSecret())
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
		// these values are not the same because of https://github.com/ory/pop/issues/591
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
