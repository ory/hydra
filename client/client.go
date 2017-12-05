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
	"strings"

	"github.com/ory/fosite"
)

// Client represents an OAuth 2.0 Client.
//
// swagger:model oAuth2Client
type Client struct {
	// ID is the id for this client.
	ID string `json:"id" gorethink:"id"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `json:"client_name" gorethink:"client_name"`

	// Secret is the client's secret. The secret will be included in the create request as cleartext, and then
	// never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users
	// that they need to write the secret down as it will not be made available again.
	Secret string `json:"client_secret,omitempty" gorethink:"client_secret"`

	// RedirectURIs is an array of allowed redirect urls for the client, for example http://mydomain/oauth/callback .
	RedirectURIs []string `json:"redirect_uris" gorethink:"redirect_uris"`

	// GrantTypes is an array of grant types the client is allowed to use.
	//
	// Pattern: client_credentials|authorize_code|implicit|refresh_token
	GrantTypes []string `json:"grant_types" gorethink:"grant_types"`

	// ResponseTypes is an array of the OAuth 2.0 response type strings that the client can
	// use at the authorization endpoint.
	//
	// Pattern: id_token|code|token
	ResponseTypes []string `json:"response_types" gorethink:"response_types"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	//
	// Pattern: ([a-zA-Z0-9\.\*]+\s?)+
	Scope string `json:"scope" gorethink:"scope"`

	// Owner is a string identifying the owner of the OAuth 2.0 Client.
	Owner string `json:"owner" gorethink:"owner"`

	// PolicyURI is a URL string that points to a human-readable privacy policy document
	// that describes how the deployment organization collects, uses,
	// retains, and discloses personal data.
	PolicyURI string `json:"policy_uri" gorethink:"policy_uri"`

	// TermsOfServiceURI is a URL string that points to a human-readable terms of service
	// document for the client that describes a contractual relationship
	// between the end-user and the client that the end-user accepts when
	// authorizing the client.
	TermsOfServiceURI string `json:"tos_uri" gorethink:"tos_uri"`

	// ClientURI is an URL string of a web page providing information about the client.
	// If present, the server SHOULD display this URL to the end-user in
	// a clickable fashion.
	ClientURI string `json:"client_uri" gorethink:"client_uri"`

	// LogoURI is an URL string that references a logo for the client.
	LogoURI string `json:"logo_uri" gorethink:"logo_uri"`

	// Contacts is a array of strings representing ways to contact people responsible
	// for this client, typically email addresses.
	Contacts []string `json:"contacts" gorethink:"contacts"`

	// Public is a boolean that identifies this client as public, meaning that it
	// does not have a secret. It will disable the client_credentials grant type for this client if set.
	Public bool `json:"public" gorethink:"public"`
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

func (c *Client) GetScopes() fosite.Arguments {
	return fosite.Arguments(strings.Fields(c.Scope))
}

func (c *Client) GetGrantTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 Grant Types that the Client is declaring
	// that it will restrict itself to using.
	// If omitted, the default is that the Client will use only the authorization_code Grant Type.
	if len(c.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}
	return fosite.Arguments(c.GrantTypes)
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// <JSON array containing a list of the OAuth 2.0 response_type values that the Client is declaring
	// that it will restrict itself to using. If omitted, the default is that the Client will use
	// only the code Response Type.
	if len(c.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}
	return fosite.Arguments(c.ResponseTypes)
}

func (c *Client) GetOwner() string {
	return c.Owner
}

func (c *Client) IsPublic() bool {
	return c.Public
}
