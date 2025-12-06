// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"github.com/go-jose/go-jose/v3"
)

// Client represents a client or an app.
type Client interface {
	// GetID returns the client ID.
	GetID() string

	// GetHashedSecret returns the hashed secret as it is stored in the store.
	GetHashedSecret() []byte

	// GetRedirectURIs returns the client's allowed redirect URIs.
	GetRedirectURIs() []string

	// GetGrantTypes returns the client's allowed grant types.
	GetGrantTypes() Arguments

	// GetResponseTypes returns the client's allowed response types.
	// All allowed combinations of response types have to be listed, each combination having
	// response types of the combination separated by a space.
	GetResponseTypes() Arguments

	// GetScopes returns the scopes this client is allowed to request.
	GetScopes() Arguments

	// IsPublic returns true, if this client is marked as public.
	IsPublic() bool

	// GetAudience returns the allowed audience(s) for this client.
	GetAudience() Arguments
}

// ClientWithSecretRotation extends Client interface by a method providing a slice of rotated secrets.
type ClientWithSecretRotation interface {
	Client
	// GetRotatedHashes returns a slice of hashed secrets used for secrets rotation.
	GetRotatedHashes() [][]byte
}

// OpenIDConnectClient represents a client capable of performing OpenID Connect requests.
type OpenIDConnectClient interface {
	// GetRequestURIs is an array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY
	// cache the contents of the files referenced by these URIs and not retrieve them at the time they are used in a request.
	// OPs can require that request_uri values used be pre-registered with the require_request_uri_registration
	// discovery parameter.
	GetRequestURIs() []string

	// GetJSONWebKeys returns the JSON Web Key Set containing the public key used by the client to authenticate.
	GetJSONWebKeys() *jose.JSONWebKeySet

	// GetJSONWebKeys returns the URL for lookup of JSON Web Key Set containing the
	// public key used by the client to authenticate.
	GetJSONWebKeysURI() string

	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP.
	// All Request Objects from this Client MUST be rejected, if not signed with this algorithm.
	GetRequestObjectSigningAlgorithm() string

	// Requested Client Authentication method for the Token Endpoint. The options are client_secret_post,
	// client_secret_basic, private_key_jwt, and none.
	GetTokenEndpointAuthMethod() string

	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing the JWT [JWT] used to authenticate the
	// Client at the Token Endpoint for the private_key_jwt authentication method.
	GetTokenEndpointAuthSigningAlgorithm() string
}

// ResponseModeClient represents a client capable of handling response_mode
type ResponseModeClient interface {
	// GetResponseMode returns the response modes that client is allowed to send
	GetResponseModes() []ResponseModeType
}

// DefaultClient is a simple default implementation of the Client interface.
type DefaultClient struct {
	ID             string   `json:"id"`
	Secret         []byte   `json:"client_secret,omitempty"`
	RotatedSecrets [][]byte `json:"rotated_secrets,omitempty"`
	RedirectURIs   []string `json:"redirect_uris"`
	GrantTypes     []string `json:"grant_types"`
	ResponseTypes  []string `json:"response_types"`
	Scopes         []string `json:"scopes"`
	Audience       []string `json:"audience"`
	Public         bool     `json:"public"`
}

type DefaultOpenIDConnectClient struct {
	*DefaultClient
	JSONWebKeysURI                    string              `json:"jwks_uri"`
	JSONWebKeys                       *jose.JSONWebKeySet `json:"jwks"`
	TokenEndpointAuthMethod           string              `json:"token_endpoint_auth_method"`
	RequestURIs                       []string            `json:"request_uris"`
	RequestObjectSigningAlgorithm     string              `json:"request_object_signing_alg"`
	TokenEndpointAuthSigningAlgorithm string              `json:"token_endpoint_auth_signing_alg"`
}

type DefaultResponseModeClient struct {
	*DefaultClient
	ResponseModes []ResponseModeType `json:"response_modes"`
}

func (c *DefaultClient) GetID() string {
	return c.ID
}

func (c *DefaultClient) IsPublic() bool {
	return c.Public
}

func (c *DefaultClient) GetAudience() Arguments {
	return c.Audience
}

func (c *DefaultClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *DefaultClient) GetHashedSecret() []byte {
	return c.Secret
}

func (c *DefaultClient) GetRotatedHashes() [][]byte {
	return c.RotatedSecrets
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
	// JSON array containing a list of the OAuth 2.0 response_type values that the Client is declaring
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
	if c.TokenEndpointAuthSigningAlgorithm == "" {
		return "RS256"
	} else {
		return c.TokenEndpointAuthSigningAlgorithm
	}
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

func (c *DefaultResponseModeClient) GetResponseModes() []ResponseModeType {
	return c.ResponseModes
}
