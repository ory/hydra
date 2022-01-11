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
 */

package client

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"

	jose "gopkg.in/square/go-jose.v2" // Naming the dependency jose is important for go-swagger to work, see https://github.com/go-swagger/go-swagger/issues/1587

	"github.com/ory/fosite"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlxx"
)

// Client represents an OAuth 2.0 Client.
//
// swagger:model oAuth2Client
type Client struct {
	ID int64 `json:"-" db:"pk"`

	// ID  is the id for this client.
	OutfacingID string `json:"client_id" db:"id"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `json:"client_name" db:"client_name"`

	// Secret is the client's secret. The secret will be included in the create request as cleartext, and then
	// never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users
	// that they need to write the secret down as it will not be made available again.
	Secret string `json:"client_secret,omitempty" db:"client_secret"`

	// RedirectURIs is an array of allowed redirect urls for the client, for example http://mydomain/oauth/callback .
	RedirectURIs sqlxx.StringSlicePipeDelimiter `json:"redirect_uris" db:"redirect_uris"`

	// GrantTypes is an array of grant types the client is allowed to use.
	//
	// Pattern: client_credentials|authorization_code|implicit|refresh_token
	GrantTypes sqlxx.StringSlicePipeDelimiter `json:"grant_types" db:"grant_types"`

	// ResponseTypes is an array of the OAuth 2.0 response type strings that the client can
	// use at the authorization endpoint.
	//
	// Pattern: id_token|code|token
	ResponseTypes sqlxx.StringSlicePipeDelimiter `json:"response_types" db:"response_types"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	//
	// Pattern: ([a-zA-Z0-9\.\*]+\s?)+
	Scope string `json:"scope" db:"scope"`

	// Audience is a whitelist defining the audiences this client is allowed to request tokens for. An audience limits
	// the applicability of an OAuth 2.0 Access Token to, for example, certain API endpoints. The value is a list
	// of URLs. URLs MUST NOT contain whitespaces.
	Audience sqlxx.StringSlicePipeDelimiter `json:"audience" db:"audience"`

	// Owner is a string identifying the owner of the OAuth 2.0 Client.
	Owner string `json:"owner" db:"owner"`

	// PolicyURI is a URL string that points to a human-readable privacy policy document
	// that describes how the deployment organization collects, uses,
	// retains, and discloses personal data.
	PolicyURI string `json:"policy_uri" db:"policy_uri"`

	// AllowedCORSOrigins are one or more URLs (scheme://host[:port]) which are allowed to make CORS requests
	// to the /oauth/token endpoint. If this array is empty, the sever's CORS origin configuration (`CORS_ALLOWED_ORIGINS`)
	// will be used instead. If this array is set, the allowed origins are appended to the server's CORS origin configuration.
	// Be aware that environment variable `CORS_ENABLED` MUST be set to `true` for this to work.
	AllowedCORSOrigins sqlxx.StringSlicePipeDelimiter `json:"allowed_cors_origins" db:"allowed_cors_origins"`

	// TermsOfServiceURI is a URL string that points to a human-readable terms of service
	// document for the client that describes a contractual relationship
	// between the end-user and the client that the end-user accepts when
	// authorizing the client.
	TermsOfServiceURI string `json:"tos_uri" db:"tos_uri"`

	// ClientURI is an URL string of a web page providing information about the client.
	// If present, the server SHOULD display this URL to the end-user in
	// a clickable fashion.
	ClientURI string `json:"client_uri" db:"client_uri"`

	// LogoURI is an URL string that references a logo for the client.
	LogoURI string `json:"logo_uri" db:"logo_uri"`

	// Contacts is a array of strings representing ways to contact people responsible
	// for this client, typically email addresses.
	Contacts sqlxx.StringSlicePipeDelimiter `json:"contacts" db:"contacts"`

	// SecretExpiresAt is an integer holding the time at which the client
	// secret will expire or 0 if it will not expire. The time is
	// represented as the number of seconds from 1970-01-01T00:00:00Z as
	// measured in UTC until the date/time of expiration.
	//
	// This feature is currently not supported and it's value will always
	// be set to 0.
	SecretExpiresAt int `json:"client_secret_expires_at" db:"client_secret_expires_at"`

	// SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a
	// list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.
	SubjectType string `json:"subject_type" db:"subject_type"`

	// URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a
	// file with a single JSON array of redirect_uri values.
	SectorIdentifierURI string `json:"sector_identifier_uri,omitempty" db:"sector_identifier_uri"`

	// URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains
	// the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the
	// Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing
	// and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced
	// JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both
	// signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used
	// to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST
	// match those in the certificate.
	JSONWebKeysURI string `json:"jwks_uri,omitempty" db:"jwks_uri"`

	// Client's JSON Web Key Set [JWK] document, passed by value. The semantics of the jwks parameter are the same as
	// the jwks_uri parameter, other than that the JWK Set is passed by value, rather than by reference. This parameter
	// is intended only to be used by Clients that, for some reason, are unable to use the jwks_uri parameter, for
	// instance, by native applications that might not have a location to host the contents of the JWK Set. If a Client
	// can use jwks_uri, it MUST NOT use jwks. One significant downside of jwks is that it does not enable key rotation
	// (which jwks_uri does, as described in Section 10 of OpenID Connect Core 1.0 [OpenID.Core]). The jwks_uri and jwks
	// parameters MUST NOT be used together.
	JSONWebKeys *x.JoseJSONWebKeySet `json:"jwks,omitempty" db:"jwks"`

	// Requested Client Authentication method for the Token Endpoint. The options are client_secret_post,
	// client_secret_basic, private_key_jwt, and none.
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty" db:"token_endpoint_auth_method"`

	// Requested Client Authentication signing algorithm for the Token Endpoint.
	TokenEndpointAuthSigningAlgorithm string `json:"token_endpoint_auth_signing_alg,omitempty" db:"token_endpoint_auth_signing_alg"`

	// Array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY cache the
	// contents of the files referenced by these URIs and not retrieve them at the time they are used in a request.
	// OPs can require that request_uri values used be pre-registered with the require_request_uri_registration
	// discovery parameter.
	RequestURIs sqlxx.StringSlicePipeDelimiter `json:"request_uris,omitempty" db:"request_uris"`

	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects
	// from this Client MUST be rejected, if not signed with this algorithm.
	RequestObjectSigningAlgorithm string `json:"request_object_signing_alg,omitempty" db:"request_object_signing_alg"`

	// JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT
	// [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims
	// as a UTF-8 encoded JSON object using the application/json content-type.
	UserinfoSignedResponseAlg string `json:"userinfo_signed_response_alg,omitempty" db:"userinfo_signed_response_alg"`

	// CreatedAt returns the timestamp of the client's creation.
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`

	// UpdatedAt returns the timestamp of the last update.
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query
	// parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the
	// request and to determine which of the potentially multiple sessions is to be logged out; if either is
	// included, both MUST be.
	FrontChannelLogoutURI string `json:"frontchannel_logout_uri,omitempty" db:"frontchannel_logout_uri"`

	// Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be
	// included to identify the RP session with the OP when the frontchannel_logout_uri is used.
	// If omitted, the default value is false.
	FrontChannelLogoutSessionRequired bool `json:"frontchannel_logout_session_required,omitempty" db:"frontchannel_logout_session_required"`

	// Array of URLs supplied by the RP to which it MAY request that the End-User's User Agent be redirected using the
	// post_logout_redirect_uri parameter after a logout has been performed.
	PostLogoutRedirectURIs sqlxx.StringSlicePipeDelimiter `json:"post_logout_redirect_uris,omitempty" db:"post_logout_redirect_uris"`

	// RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.
	BackChannelLogoutURI string `json:"backchannel_logout_uri,omitempty" db:"backchannel_logout_uri"`

	// Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout
	// Token to identify the RP session with the OP when the backchannel_logout_uri is used.
	// If omitted, the default value is false.
	BackChannelLogoutSessionRequired bool `json:"backchannel_logout_session_required,omitempty" db:"backchannel_logout_session_required"`

	// Metadata is arbitrary data.
	Metadata sqlxx.JSONRawMessage `json:"metadata,omitempty" db:"metadata"`

	// RegistrationAccessTokenSignature is contains the signature of the registration token for managing the OAuth2 Client.
	RegistrationAccessTokenSignature string `json:"-" db:"registration_access_token_signature"`

	// RegistrationAccessToken can be used to update, get, or delete the OAuth2 Client.
	RegistrationAccessToken string `json:"registration_access_token,omitempty" db:"-"`

	// RegistrationClientURI is the URL used to update, get, or delete the OAuth2 Client.
	RegistrationClientURI string `json:"registration_client_uri,omitempty" db:"-"`
}

func (Client) TableName() string {
	return "hydra_client"
}

func (c *Client) BeforeSave(_ *pop.Connection) error {
	if c.JSONWebKeys == nil {
		c.JSONWebKeys = new(x.JoseJSONWebKeySet)
	}

	if c.Metadata == nil {
		c.Metadata = []byte("{}")
	}

	if c.Audience == nil {
		c.Audience = sqlxx.StringSlicePipeDelimiter{}
	}

	if c.AllowedCORSOrigins == nil {
		c.AllowedCORSOrigins = sqlxx.StringSlicePipeDelimiter{}
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.CreatedAt = c.CreatedAt.UTC()

	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}
	c.UpdatedAt = c.UpdatedAt.UTC()

	return nil
}

func (c *Client) GetID() string {
	return c.OutfacingID
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

func (c *Client) GetAudience() fosite.Arguments {
	return fosite.Arguments(c.Audience)
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

func (c *Client) GetResponseModes() []fosite.ResponseModeType {
	return []fosite.ResponseModeType{
		fosite.ResponseModeDefault,
		fosite.ResponseModeFormPost,
		fosite.ResponseModeQuery,
		fosite.ResponseModeFragment,
	}
}

func (c *Client) GetOwner() string {
	return c.Owner
}

func (c *Client) IsPublic() bool {
	return c.TokenEndpointAuthMethod == "none"
}

func (c *Client) GetJSONWebKeysURI() string {
	return c.JSONWebKeysURI
}

func (c *Client) GetJSONWebKeys() *jose.JSONWebKeySet {
	if c.JSONWebKeys == nil {
		return nil
	}
	return c.JSONWebKeys.JSONWebKeySet
}

func (c *Client) GetTokenEndpointAuthSigningAlgorithm() string {
	if c.TokenEndpointAuthSigningAlgorithm == "" {
		return "RS256"
	}
	return c.TokenEndpointAuthSigningAlgorithm
}

func (c *Client) GetRequestObjectSigningAlgorithm() string {
	return c.RequestObjectSigningAlgorithm
}

func (c *Client) GetTokenEndpointAuthMethod() string {
	if c.TokenEndpointAuthMethod == "" {
		return "client_secret_basic"
	}
	return c.TokenEndpointAuthMethod
}

func (c *Client) GetRequestURIs() []string {
	return c.RequestURIs
}
