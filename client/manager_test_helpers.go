package client

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/ory/fosite"
)

func TestHelperClientAutoGenerateKey(k string, m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		c := &Client{
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}
		assert.Nil(t, m.CreateClient(c))
		assert.NotEmpty(t, c.ID)
		assert.Nil(t, m.DeleteClient(c.ID))
	}
}

func TestHelperCreateGetDeleteClient(k string, m Storage) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := m.GetClient(nil, "4321")
		assert.NotNil(t, err)

		c := &Client{
			ID:                "1234",
			Name:              "name",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}
		err = m.CreateClient(c)
		assert.Nil(t, err)
		if err == nil {
			compare(t, c, k)
		}

		err = m.CreateClient(&Client{
			ID:                "2-1234",
			Name:              "name",
			Secret:            "secret",
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		})
		assert.Nil(t, err)

		// RethinkDB delay
		time.Sleep(100 * time.Millisecond)

		d, err := m.GetClient(nil, "1234")
		assert.Nil(t, err)
		if err == nil {
			compare(t, d, k)
		}

		ds, err := m.GetClients()
		assert.Nil(t, err)
		assert.Len(t, ds, 2)
		assert.NotEqual(t, ds["1234"].ID, ds["2-1234"].ID)

		err = m.UpdateClient(&Client{
			ID:                "2-1234",
			Name:              "name-new",
			Secret:            "secret-new",
			RedirectURIs:      []string{"http://redirect/new"},
			TermsOfServiceURI: "bar",
		})
		assert.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		nc, err := m.GetConcreteClient("2-1234")
		assert.Nil(t, err)

		if k != "http" {
			// http always returns an empty secret
			assert.NotEqual(t, d.GetHashedSecret(), nc.GetHashedSecret())
		}
		assert.Equal(t, "bar", nc.TermsOfServiceURI)
		assert.Equal(t, "name-new", nc.Name)
		assert.EqualValues(t, []string{"http://redirect/new"}, nc.GetRedirectURIs())
		assert.Zero(t, len(nc.Contacts))

		err = m.DeleteClient("1234")
		assert.Nil(t, err)

		// RethinkDB delay
		time.Sleep(100 * time.Millisecond)

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
