// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"testing"

	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelperClientAutoGenerateKey(k string, m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		c := &Client{
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}
		assert.NoError(t, m.CreateClient(c))
		assert.NotEmpty(t, c.ID)
		assert.NoError(t, m.DeleteClient(c.ID))
	}
}

func TestHelperClientAuthenticate(k string, m Manager) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		m.CreateClient(&Client{
			ID:           "1234321",
			Secret:       "secret",
			RedirectURIs: []string{"http://redirect"},
		})

		c, err := m.Authenticate("1234321", []byte("secret1"))
		require.NotNil(t, err)

		c, err = m.Authenticate("1234321", []byte("secret"))
		require.NoError(t, err)
		assert.Equal(t, "1234321", c.ID)
	}
}

func TestHelperCreateGetDeleteClient(k string, m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		_, err := m.GetClient(nil, "4321")
		assert.NotNil(t, err)

		c := &Client{
			ID:                "1234",
			Name:              "name",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}

		assert.NoError(t, m.CreateClient(c))
		if err == nil {
			compare(t, c, k)
		}

		assert.NoError(t, m.CreateClient(&Client{
			ID:                "2-1234",
			Name:              "name",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}))

		d, err := m.GetClient(nil, "1234")
		assert.NoError(t, err)

		if err == nil {
			compare(t, d, k)
		}

		ds, err := m.GetClients()
		assert.NoError(t, err)
		assert.Len(t, ds, 2)
		assert.NotEqual(t, ds["1234"].ID, ds["2-1234"].ID)

		err = m.UpdateClient(&Client{
			ID:                "2-1234",
			Name:              "name-new",
			Secret:            "secret-new",
			RedirectURIs:      []string{"http://redirect/new"},
			TermsOfServiceURI: "bar",
		})
		assert.NoError(t, err)

		nc, err := m.GetConcreteClient("2-1234")
		assert.NoError(t, err)

		if k != "http" {
			// http always returns an empty secret
			assert.NotEqual(t, d.GetHashedSecret(), nc.GetHashedSecret())
		}
		assert.Equal(t, "bar", nc.TermsOfServiceURI)
		assert.Equal(t, "name-new", nc.Name)
		assert.EqualValues(t, []string{"http://redirect/new"}, nc.GetRedirectURIs())
		assert.Zero(t, len(nc.Contacts))

		err = m.DeleteClient("1234")
		assert.NoError(t, err)

		_, err = m.GetClient(nil, "1234")
		assert.NotNil(t, err)
	}
}

func compare(t *testing.T, c fosite.Client, k string) {
	assert.Equal(t, c.GetID(), "1234")
	if k != "http" {
		assert.NotEmpty(t, c.GetHashedSecret())
	}
	assert.Equal(t, c.GetRedirectURIs(), []string{"http://redirect"})
}
