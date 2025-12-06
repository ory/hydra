// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/fosite"
)

var _ fosite.OpenIDConnectClient = new(Client)
var _ fosite.Client = new(Client)

func TestClient(t *testing.T) {
	c := &Client{
		ID:                      "foo",
		RedirectURIs:            []string{"foo"},
		Scope:                   "foo bar",
		TokenEndpointAuthMethod: "none",
	}

	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())
	assert.EqualValues(t, []byte(c.Secret), c.GetHashedSecret())
	assert.EqualValues(t, fosite.Arguments{"authorization_code"}, c.GetGrantTypes())
	assert.EqualValues(t, fosite.Arguments{"code"}, c.GetResponseTypes())
	assert.EqualValues(t, c.Owner, c.GetOwner())
	assert.True(t, c.IsPublic())
	assert.Len(t, c.GetScopes(), 2)
	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())
}
