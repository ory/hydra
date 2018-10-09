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
 *
 */

package fosite

import (
	"gopkg.in/square/go-jose.v2"
)

// Client represents a client or an app.
type Client interface {
	// GetID returns the client ID.
	GetID() string

	// GetHashedSecret returns the hashed secret as it is stored in the store.
	GetHashedSecret() []byte

	// Returns the client's allowed redirect URIs.
	GetRedirectURIs() []string

	// Returns the client's allowed grant types.
	GetGrantTypes() Arguments

	// Returns the client's allowed response types.
	GetResponseTypes() Arguments

	// Returns the scopes this client is allowed to request.
	GetScopes() Arguments

	// IsPublic returns true, if this client is marked as public.
	IsPublic() bool
}

type OpenIDConnectClient interface {
	// Array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY cache the
	// contents of the files referenced by these URIs and not retrieve them at the time they are used in a request.
	// OPs can require that request_uri values used be pre-registered with the require_request_uri_registration
	// discovery parameter.
	GetRequestURIs() []string

	// GetJSONWebKeys returns the JSON Web Key Set containing the public keys used by the client to authenticate.
	GetJSONWebKeys() *jose.JSONWebKeySet

	// GetJSONWebKeys returns the URL for lookup of JSON Web Key Set containing the
	// public keys used by the client to authenticate.
	GetJSONWebKeysURI() string

	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP.
	// All Request Objects from this Client MUST be rejected, if not signed with this algorithm.
	GetRequestObjectSigningAlgorithm() string

	// Requested Client Authentication method for the Token Endpoint. The options are client_secret_post,
	// client_secret_basic, client_secret_jwt, private_key_jwt, and none.
	GetTokenEndpointAuthMethod() string

	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing the JWT [JWT] used to authenticate the
	// Client at the Token Endpoint for the private_key_jwt and client_secret_jwt authentication methods.
	GetTokenEndpointAuthSigningAlgorithm() string
}

// DefaultClient is a simple default implementation of the Client interface.
type DefaultClient struct {
	ID            string   `json:"id"`
	Secret        []byte   `json:"client_secret,omitempty"`
	RedirectURIs  []string `json:"redirect_uris"`
	GrantTypes    []string `json:"grant_types"`
	ResponseTypes []string `json:"response_types"`
	Scopes        []string `json:"scopes"`
	Public        bool     `json:"public"`
}

type DefaultOpenIDConnectClient struct {
	*DefaultClient
	JSONWebKeysURI                string              `json:"jwks_uri"`
	JSONWebKeys                   *jose.JSONWebKeySet `json:"jwks"`
	TokenEndpointAuthMethod       string              `json:"token_endpoint_auth_method"`
	RequestURIs                   []string            `json:"request_uris"`
	RequestObjectSigningAlgorithm string              `json:"request_object_signing_alg"`
}

func (c *DefaultClient) GetID() string {
	return c.ID
}

func (c *DefaultClient) IsPublic() bool {
	return c.Public
}

func (c *DefaultClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *DefaultClient) GetHashedSecret() []byte {
	return c.Secret
}

func (c *DefaultClient) GetScopes() Arguments {
	return c.Scopes
}

func (c *DefaultClient) GetGrantTypes() Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 Grant Types that the Client is declaring
	// that it will restrict itself to using.
	// If omitted, the default is that the Client will use only the authorization_code Grant Type.
	if len(c.GrantTypes) == 0 {
		return Arguments{"authorization_code"}
	}
	return Arguments(c.GrantTypes)
}

func (c *DefaultClient) GetResponseTypes() Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// <JSON array containing a list of the OAuth 2.0 response_type values that the Client is declaring
	// that it will restrict itself to using. If omitted, the default is that the Client will use
	// only the code Response Type.
	if len(c.ResponseTypes) == 0 {
		return Arguments{"code"}
	}
	return Arguments(c.ResponseTypes)
}

func (c *DefaultOpenIDConnectClient) GetJSONWebKeysURI() string {
	return c.JSONWebKeysURI
}

func (c *DefaultOpenIDConnectClient) GetJSONWebKeys() *jose.JSONWebKeySet {
	return c.JSONWebKeys
}

func (c *DefaultOpenIDConnectClient) GetTokenEndpointAuthSigningAlgorithm() string {
	return "RS256"
}

func (c *DefaultOpenIDConnectClient) GetRequestObjectSigningAlgorithm() string {
	return c.RequestObjectSigningAlgorithm
}

func (c *DefaultOpenIDConnectClient) GetTokenEndpointAuthMethod() string {
	return c.TokenEndpointAuthMethod
}

func (c *DefaultOpenIDConnectClient) GetRequestURIs() []string {
	return c.RequestURIs
}
