package client

import (
	"strings"

	"github.com/ory-am/fosite"
)

type Client struct {
	ID                string   `json:"id" gorethink:"id"`
	Name              string   `json:"client_name" gorethink:"client_name"`
	Secret            string   `json:"client_secret,omitempty" gorethink:"client_secret"`
	RedirectURIs      []string `json:"redirect_uris" gorethink:"redirect_uris"`
	GrantTypes        []string `json:"grant_types" gorethink:"grant_types"`
	ResponseTypes     []string `json:"response_types" gorethink:"response_types"`
	Scope             string   `json:"scope" gorethink:"scope"`
	Owner             string   `json:"owner" gorethink:"owner"`
	PolicyURI         string   `json:"policy_uri" gorethink:"policy_uri"`
	TermsOfServiceURI string   `json:"tos_uri" gorethink:"tos_uri"`
	ClientURI         string   `json:"client_uri" gorethink:"client_uri"`
	LogoURI           string   `json:"logo_uri" gorethink:"logo_uri"`
	Contacts          []string `json:"contacts" gorethink:"contacts"`
	Public            bool     `json:"public" gorethink:"public"`
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
	return fosite.Arguments(strings.Split(c.Scope, " "))
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
