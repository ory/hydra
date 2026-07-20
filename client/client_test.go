// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestClientJSONOmitsEmptyOptionalURIFields(t *testing.T) {
	c := &Client{
		ID:           "test-client",
		RedirectURIs: []string{"http://localhost/callback"},
	}

	data, err := json.Marshal(c)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// These optional URI fields should be omitted when empty, not serialized as ""
	_, hasLogoURI := result["logo_uri"]
	assert.False(t, hasLogoURI, "logo_uri should be omitted when empty")

	_, hasTosURI := result["tos_uri"]
	assert.False(t, hasTosURI, "tos_uri should be omitted when empty")

	_, hasPolicyURI := result["policy_uri"]
	assert.False(t, hasPolicyURI, "policy_uri should be omitted when empty")

	_, hasClientURI := result["client_uri"]
	assert.False(t, hasClientURI, "client_uri should be omitted when empty")

	// contacts should be omitted when nil, not serialized as null
	_, hasContacts := result["contacts"]
	assert.False(t, hasContacts, "contacts should be omitted when nil")
}
