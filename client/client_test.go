package client

import (
	"github.com/ory-am/fosite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	c := &Client{
		ID:           "foo",
		RedirectURIs: []string{"foo"},
		Scope:        "foo bar",
	}

	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())
	assert.EqualValues(t, []byte(c.Secret), c.GetHashedSecret())
	assert.EqualValues(t, fosite.Arguments{"authorization_code"}, c.GetGrantTypes())
	assert.EqualValues(t, fosite.Arguments{"code"}, c.GetResponseTypes())
	assert.EqualValues(t, (c.Owner), c.GetOwner())
	assert.EqualValues(t, (c.Public), c.IsPublic())
	assert.Len(t, c.GetScopes(), 2)
	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())
}
