// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/twmb/murmur3"

	"github.com/ory/pop/v6"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/sqlxx"
)

var (
	_ fosite.OpenIDConnectClient = (*Client)(nil)
	_ fosite.Client              = (*Client)(nil)
)

// OAuth 2.0 Client
//
// OAuth 2.0 Clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
// generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.
//
// swagger:model oAuth2Client
type Client struct {
	NID uuid.UUID `db:"nid" faker:"-" json:"-"`

	// OAuth 2.0 Client ID
	//
	// The ID is immutable. If no ID is provided, a UUID4 will be generated.
	ID string `json:"client_id" db:"id"`

	// OAuth 2.0 Client Name
	//
	// The human-readable name of the client to be presented to the
	// end-user during authorization.
	Name string `json:"client_name" db:"client_name"`

	// OAuth 2.0 Client Secret
	//
	// The secret will be included in the create request as cleartext, and then
	// never again. The secret is kept in hashed format and is not recoverable once lost.
	Secret string `json:"client_secret,omitempty" db:"client_secret"`

	// OAuth 2.0 Client Redirect URIs
	//
	// RedirectURIs is an array of allowed redirect urls for the client.
	//
	// Example: http://mydomain/oauth/callback
	RedirectURIs sqlxx.StringSliceJSONFormat `json:"redirect_uris" db:"redirect_uris"`

	// OAuth 2.0 Client Grant Types
	//
	// An array of OAuth 2.0 grant types the client is allowed to use. Can be one
	// of:
	//
	// - Client Credentials Grant: `client_credentials`
	// - Authorization Code Grant: `authorization_code`
	// - OpenID Connect Implicit Grant (deprecated!): `implicit`
	// - Refresh Token Grant: `refresh_token`
	// - OAuth 2.0 Token Exchange: `urn:ietf:params:oauth:grant-type:jwt-bearer`
	// - OAuth 2.0 Device Code Grant: `urn:ietf:params:oauth:grant-type:device_code`
	GrantTypes sqlxx.StringSliceJSONFormat `json:"grant_types" db:"grant_types"`

	// OAuth 2.0 Client Response Types
	//
	// An array of the OAuth 2.0 response type strings that the client can
	// use at the authorization endpoint. Can be one of:
	//
	// - Needed for OpenID Connect Implicit Grant:
	//   - Returns ID Token to redirect URI: `id_token`
	//   - Returns Access token redirect URI: `token`
	// - Needed for Authorization Code Grant: `code`
	ResponseTypes sqlxx.StringSliceJSONFormat `json:"response_types" db:"response_types"`

	// OAuth 2.0 Client Scope
	//
	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	//
	// Example: scope1 scope-2 scope.3 scope:4
	Scope string `json:"scope" db:"scope"`

	// OAuth 2.0 Client Audience
	//
	// An allow-list defining the audiences this client is allowed to request tokens for. An audience limits
	// the applicability of an OAuth 2.0 Access Token to, for example, certain API endpoints. The value is a list
	// of URLs. URLs MUST NOT contain whitespaces.
	//
	// Example: https://mydomain.com/api/users, https://mydomain.com/api/posts
	Audience sqlxx.StringSliceJSONFormat `json:"audience" db:"audience"`

	// OAuth 2.0 Client Owner
	//
	// Owner is a string identifying the owner of the OAuth 2.0 Client.
	Owner string `json:"owner" db:"owner"`

	// OAuth 2.0 Client Policy URI
	//
	// PolicyURI is a URL string that points to a human-readable privacy policy document
	// that describes how the deployment organization collects, uses,
	// retains, and discloses personal data.
	PolicyURI string `json:"policy_uri" db:"policy_uri"`

	// OAuth 2.0 Client Allowed CORS Origins
	//
	// One or more URLs (scheme://host[:port]) which are allowed to make CORS requests
	// to the /oauth/token endpoint. If this array is empty, the sever's CORS origin configuration (`CORS_ALLOWED_ORIGINS`)
	// will be used instead. If this array is set, the allowed origins are appended to the server's CORS origin configuration.
	// Be aware that environment variable `CORS_ENABLED` MUST be set to `true` for this to work.
	AllowedCORSOrigins sqlxx.StringSliceJSONFormat `json:"allowed_cors_origins" db:"allowed_cors_origins"`

	// OAuth 2.0 Client Terms of Service URI
	//
	// A URL string pointing to a human-readable terms of service
	// document for the client that describes a contractual relationship
	// between the end-user and the client that the end-user accepts when
	// authorizing the client.
	TermsOfServiceURI string `json:"tos_uri" db:"tos_uri"`

	// OAuth 2.0 Client URI
	//
	// ClientURI is a URL string of a web page providing information about the client.
	// If present, the server SHOULD display this URL to the end-user in
	// a clickable fashion.
	ClientURI string `json:"client_uri" db:"client_uri"`

	// OAuth 2.0 Client Logo URI
	//
	// A URL string referencing the client's logo.
	LogoURI string `json:"logo_uri" db:"logo_uri"`

	// OAuth 2.0 Client Contact
	//
	// An array of strings representing ways to contact people responsible
	// for this client, typically email addresses.
	//
	// Example: help@example.org
	Contacts sqlxx.StringSliceJSONFormat `json:"contacts" db:"contacts"`

	// OAuth 2.0 Client Secret Expires At
	//
	// The field is currently not supported and its value is always 0.
	SecretExpiresAt int `json:"client_secret_expires_at" db:"client_secret_expires_at"`

	// OpenID Connect Subject Type
	//
	// The `subject_types_supported` Discovery parameter contains a
	// list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.
	SubjectType string `json:"subject_type" db:"subject_type" faker:"len=15"`

	// OpenID Connect Sector Identifier URI
	//
	// URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a
	// file with a single JSON array of redirect_uri values.
	SectorIdentifierURI string `json:"sector_identifier_uri,omitempty" db:"sector_identifier_uri"`

	// OAuth 2.0 Client JSON Web Key Set URL
	//
	// URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains
	// the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the
	// Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing
	// and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced
	// JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both
	// signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used
	// to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST
	// match those in the certificate.
	JSONWebKeysURI string `json:"jwks_uri,omitempty" db:"jwks_uri"`

	// OAuth 2.0 Client JSON Web Key Set
	//
	// Client's JSON Web Key Set [JWK] document, passed by value. The semantics of the jwks parameter are the same as
	// the jwks_uri parameter, other than that the JWK Set is passed by value, rather than by reference. This parameter
	// is intended only to be used by Clients that, for some reason, are unable to use the jwks_uri parameter, for
	// instance, by native applications that might not have a location to host the contents of the JWK Set. If a Client
	// can use jwks_uri, it MUST NOT use jwks. One significant downside of jwks is that it does not enable key rotation
	// (which jwks_uri does, as described in Section 10 of OpenID Connect Core 1.0 [OpenID.Core]). The jwks_uri and jwks
	// parameters MUST NOT be used together.
	JSONWebKeys *x.JoseJSONWebKeySet `json:"jwks,omitempty" db:"jwks" faker:"-"`

	// OAuth 2.0 Token Endpoint Authentication Method
	//
	// Requested Client Authentication method for the Token Endpoint. The options are:
	//
	// - `client_secret_basic`: (default) Send `client_id` and `client_secret` as `application/x-www-form-urlencoded` encoded in the HTTP Authorization header.
	// - `client_secret_post`: Send `client_id` and `client_secret` as `application/x-www-form-urlencoded` in the HTTP body.
	// - `private_key_jwt`: Use JSON Web Tokens to authenticate the client.
	// - `none`: Used for public clients (native apps, mobile apps) which can not have secrets.
	//
	// default: client_secret_basic
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty" db:"token_endpoint_auth_method" faker:"len=25"`

	// OAuth 2.0 Token Endpoint Signing Algorithm
	//
	// Requested Client Authentication signing algorithm for the Token Endpoint.
	TokenEndpointAuthSigningAlgorithm string `json:"token_endpoint_auth_signing_alg,omitempty" db:"token_endpoint_auth_signing_alg" faker:"len=10"`

	// OpenID Connect Request URIs
	//
	// Array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY cache the
	// contents of the files referenced by these URIs and not retrieve them at the time they are used in a request.
	// OPs can require that request_uri values used be pre-registered with the require_request_uri_registration
	// discovery parameter.
	RequestURIs sqlxx.StringSliceJSONFormat `json:"request_uris,omitempty" db:"request_uris"`

	// OpenID Connect Request Object Signing Algorithm
	//
	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects
	// from this Client MUST be rejected, if not signed with this algorithm.
	RequestObjectSigningAlgorithm string `json:"request_object_signing_alg,omitempty" db:"request_object_signing_alg" faker:"len=10"`

	// OpenID Connect Request Userinfo Signed Response Algorithm
	//
	// JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT
	// [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims
	// as a UTF-8 encoded JSON object using the application/json content-type.
	UserinfoSignedResponseAlg string `json:"userinfo_signed_response_alg,omitempty" db:"userinfo_signed_response_alg" faker:"len=10"`

	// OAuth 2.0 Client Creation Date
	//
	// CreatedAt returns the timestamp of the client's creation.
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`

	// OAuth 2.0 Client Last Update Date
	//
	// UpdatedAt returns the timestamp of the last update.
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// OpenID Connect Front-Channel Logout URI
	//
	// RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query
	// parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the
	// request and to determine which of the potentially multiple sessions is to be logged out; if either is
	// included, both MUST be.
	FrontChannelLogoutURI string `json:"frontchannel_logout_uri,omitempty" db:"frontchannel_logout_uri"`

	// OpenID Connect Front-Channel Logout Session Required
	//
	// Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be
	// included to identify the RP session with the OP when the frontchannel_logout_uri is used.
	// If omitted, the default value is false.
	FrontChannelLogoutSessionRequired bool `json:"frontchannel_logout_session_required,omitempty" db:"frontchannel_logout_session_required"`

	// Allowed Post-Redirect Logout URIs
	//
	// Array of URLs supplied by the RP to which it MAY request that the End-User's User Agent be redirected using the
	// post_logout_redirect_uri parameter after a logout has been performed.
	PostLogoutRedirectURIs sqlxx.StringSliceJSONFormat `json:"post_logout_redirect_uris,omitempty" db:"post_logout_redirect_uris"`

	// OpenID Connect Back-Channel Logout URI
	//
	// RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.
	BackChannelLogoutURI string `json:"backchannel_logout_uri,omitempty" db:"backchannel_logout_uri"`

	// OpenID Connect Back-Channel Logout Session Required
	//
	// Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout
	// Token to identify the RP session with the OP when the backchannel_logout_uri is used.
	// If omitted, the default value is false.
	BackChannelLogoutSessionRequired bool `json:"backchannel_logout_session_required,omitempty" db:"backchannel_logout_session_required"`

	// OAuth 2.0 Client Metadata
	//
	// Use this field to story arbitrary data about the OAuth 2.0 Client. Can not be modified using OpenID Connect Dynamic Client Registration protocol.
	Metadata sqlxx.JSONRawMessage `json:"metadata,omitempty" db:"metadata" faker:"-"`

	// OpenID Connect Dynamic Client Registration Access Token
	//
	// RegistrationAccessTokenSignature is contains the signature of the registration token for managing the OAuth2 Client.
	RegistrationAccessTokenSignature string `json:"-" db:"registration_access_token_signature"`

	// OpenID Connect Dynamic Client Registration Access Token
	//
	// RegistrationAccessToken can be used to update, get, or delete the OAuth2 Client. It is sent when creating a client
	// using Dynamic Client Registration.
	RegistrationAccessToken string `json:"registration_access_token,omitempty" db:"-"`

	// OpenID Connect Dynamic Client Registration URL
	//
	// RegistrationClientURI is the URL used to update, get, or delete the OAuth2 Client.
	RegistrationClientURI string `json:"registration_client_uri,omitempty" db:"-"`

	// OAuth 2.0 Access Token Strategy
	//
	// AccessTokenStrategy is the strategy used to generate access tokens.
	// Valid options are `jwt` and `opaque`. `jwt` is a bad idea, see https://www.ory.sh/docs/oauth2-oidc/jwt-access-token
	// Setting the strategy here overrides the global setting in `strategies.access_token`.
	AccessTokenStrategy string `json:"access_token_strategy,omitempty" db:"access_token_strategy" faker:"-"`

	// SkipConsent skips the consent screen for this client. This field can only
	// be set from the admin API.
	SkipConsent bool `json:"skip_consent" db:"skip_consent" faker:"-"`

	// SkipLogoutConsent skips the logout consent screen for this client. This field can only
	// be set from the admin API.
	SkipLogoutConsent sqlxx.NullBool `json:"skip_logout_consent" db:"skip_logout_consent" faker:"-"`

	Lifespans
}

// OAuth 2.0 Client Token Lifespans
//
// Lifespans of different token types issued for this OAuth 2.0 Client.
//
// swagger:model oAuth2ClientTokenLifespans
type Lifespans struct {
	// OAuth2 Authorization Code Grant Access Token Lifespan
	//
	// The lifespan of an access token issued by the OAuth 2.0 Authorization Code Grant for this OAuth 2.0 Client.
	AuthorizationCodeGrantAccessTokenLifespan x.NullDuration `json:"authorization_code_grant_access_token_lifespan,omitempty" db:"authorization_code_grant_access_token_lifespan"`

	// OAuth2 Authorization Code Grant Access ID Lifespan
	//
	// The lifespan of an ID token issued by the OAuth 2.0 Authorization Code Grant for this OAuth 2.0 Client.
	AuthorizationCodeGrantIDTokenLifespan x.NullDuration `json:"authorization_code_grant_id_token_lifespan,omitempty" db:"authorization_code_grant_id_token_lifespan"`

	// OAuth2 Authorization Code Grant Access Refresh Lifespan
	//
	// The lifespan of a refresh token issued by the OAuth 2.0 Authorization Code Grant for this OAuth 2.0 Client.
	AuthorizationCodeGrantRefreshTokenLifespan x.NullDuration `json:"authorization_code_grant_refresh_token_lifespan,omitempty" db:"authorization_code_grant_refresh_token_lifespan"`

	// OAuth2 Client Credentials Grant Access Token Lifespan
	//
	// The lifespan of an access token issued by the OAuth 2.0 Client Credentials Grant for this OAuth 2.0 Client.
	ClientCredentialsGrantAccessTokenLifespan x.NullDuration `json:"client_credentials_grant_access_token_lifespan,omitempty" db:"client_credentials_grant_access_token_lifespan"`

	// OpenID Connect Implicit Grant Access Token Lifespan
	//
	// The lifespan of an access token issued by the OpenID Connect Implicit Grant for this OAuth 2.0 Client.
	ImplicitGrantAccessTokenLifespan x.NullDuration `json:"implicit_grant_access_token_lifespan,omitempty" db:"implicit_grant_access_token_lifespan"`

	// OpenID Connect Implicit Grant ID Token Lifespan
	//
	// The lifespan of an ID token issued by the OpenID Connect Implicit Grant for this OAuth 2.0 Client.
	ImplicitGrantIDTokenLifespan x.NullDuration `json:"implicit_grant_id_token_lifespan,omitempty" db:"implicit_grant_id_token_lifespan"`

	// OpenID Connect Implicit Grant Access Token Lifespan
	//
	// The lifespan of an access token issued by the OpenID Connect Implicit Grant for this OAuth 2.0 Client.
	JwtBearerGrantAccessTokenLifespan x.NullDuration `json:"jwt_bearer_grant_access_token_lifespan,omitempty" db:"jwt_bearer_grant_access_token_lifespan"`

	// DEPRECATED: This field has no effect.
	PasswordGrantAccessTokenLifespan x.NullDuration `json:"-" db:"password_grant_access_token_lifespan"`

	// DEPRECATED: This field has no effect.
	PasswordGrantRefreshTokenLifespan x.NullDuration `json:"-" db:"password_grant_refresh_token_lifespan"`

	// OAuth2 2.0 Refresh Token Grant ID Token Lifespan
	//
	// The lifespan of an ID token issued by the OAuth2 2.0 Refresh Token Grant for this OAuth 2.0 Client.
	RefreshTokenGrantIDTokenLifespan x.NullDuration `json:"refresh_token_grant_id_token_lifespan,omitempty" db:"refresh_token_grant_id_token_lifespan"`

	// OAuth2 2.0 Refresh Token Grant Access Token Lifespan
	//
	// The lifespan of an access token issued by the OAuth2 2.0 Refresh Token Grant for this OAuth 2.0 Client.
	RefreshTokenGrantAccessTokenLifespan x.NullDuration `json:"refresh_token_grant_access_token_lifespan,omitempty" db:"refresh_token_grant_access_token_lifespan"`

	// OAuth2 2.0 Refresh Token Grant Refresh Token Lifespan
	//
	// The lifespan of a refresh token issued by the OAuth2 2.0 Refresh Token Grant for this OAuth 2.0 Client.
	RefreshTokenGrantRefreshTokenLifespan x.NullDuration `json:"refresh_token_grant_refresh_token_lifespan,omitempty" db:"refresh_token_grant_refresh_token_lifespan"`

	// OAuth2 2.0 Device Authorization Grant ID Token Lifespan
	//
	// The lifespan of an ID token issued by the OAuth2 2.0 Device Authorization Grant for this OAuth 2.0 Client.
	DeviceAuthorizationGrantIDTokenLifespan x.NullDuration `json:"device_authorization_grant_id_token_lifespan,omitempty" db:"device_authorization_grant_id_token_lifespan"`

	// OAuth2 2.0 Device Authorization Grant Access Token Lifespan
	//
	// The lifespan of an access token issued by the OAuth2 2.0 Device Authorization Grant for this OAuth 2.0 Client.
	DeviceAuthorizationGrantAccessTokenLifespan x.NullDuration `json:"device_authorization_grant_access_token_lifespan,omitempty" db:"device_authorization_grant_access_token_lifespan"`

	// OAuth2 2.0 Device Authorization Grant Device Authorization Lifespan
	//
	// The lifespan of a Device Authorization issued by the OAuth2 2.0 Device Authorization Grant for this OAuth 2.0 Client.
	DeviceAuthorizationGrantRefreshTokenLifespan x.NullDuration `json:"device_authorization_grant_refresh_token_lifespan,omitempty" db:"device_authorization_grant_refresh_token_lifespan"`
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
		c.Audience = sqlxx.StringSliceJSONFormat{}
	}

	if c.AllowedCORSOrigins == nil {
		c.AllowedCORSOrigins = sqlxx.StringSliceJSONFormat{}
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
	return c.ID
}

func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

func (c *Client) GetScopes() fosite.Arguments {
	return strings.Fields(c.Scope)
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

var _ fosite.ClientWithCustomTokenLifespans = &Client{}

func (c *Client) GetEffectiveLifespan(gt fosite.GrantType, tt fosite.TokenType, fallback time.Duration) time.Duration {
	var cl *time.Duration
	if gt == fosite.GrantTypeAuthorizationCode {
		if tt == fosite.AccessToken && c.AuthorizationCodeGrantAccessTokenLifespan.Valid {
			cl = &c.AuthorizationCodeGrantAccessTokenLifespan.Duration
		} else if tt == fosite.IDToken && c.AuthorizationCodeGrantIDTokenLifespan.Valid {
			cl = &c.AuthorizationCodeGrantIDTokenLifespan.Duration
		} else if tt == fosite.RefreshToken && c.AuthorizationCodeGrantRefreshTokenLifespan.Valid {
			cl = &c.AuthorizationCodeGrantRefreshTokenLifespan.Duration
		}
	} else if gt == fosite.GrantTypeClientCredentials {
		if tt == fosite.AccessToken && c.ClientCredentialsGrantAccessTokenLifespan.Valid {
			cl = &c.ClientCredentialsGrantAccessTokenLifespan.Duration
		}
	} else if gt == fosite.GrantTypeImplicit {
		if tt == fosite.AccessToken && c.ImplicitGrantAccessTokenLifespan.Valid {
			cl = &c.ImplicitGrantAccessTokenLifespan.Duration
		} else if tt == fosite.IDToken && c.ImplicitGrantIDTokenLifespan.Valid {
			cl = &c.ImplicitGrantIDTokenLifespan.Duration
		}
	} else if gt == fosite.GrantTypeJWTBearer {
		if tt == fosite.AccessToken && c.JwtBearerGrantAccessTokenLifespan.Valid {
			cl = &c.JwtBearerGrantAccessTokenLifespan.Duration
		}
	} else if gt == fosite.GrantTypePassword {
		if tt == fosite.AccessToken && c.PasswordGrantAccessTokenLifespan.Valid {
			cl = &c.PasswordGrantAccessTokenLifespan.Duration
		} else if tt == fosite.RefreshToken && c.PasswordGrantRefreshTokenLifespan.Valid {
			cl = &c.PasswordGrantRefreshTokenLifespan.Duration
		}
	} else if gt == fosite.GrantTypeRefreshToken {
		if tt == fosite.AccessToken && c.RefreshTokenGrantAccessTokenLifespan.Valid {
			cl = &c.RefreshTokenGrantAccessTokenLifespan.Duration
		} else if tt == fosite.IDToken && c.RefreshTokenGrantIDTokenLifespan.Valid {
			cl = &c.RefreshTokenGrantIDTokenLifespan.Duration
		} else if tt == fosite.RefreshToken && c.RefreshTokenGrantRefreshTokenLifespan.Valid {
			cl = &c.RefreshTokenGrantRefreshTokenLifespan.Duration
		}
	} else if gt == fosite.GrantTypeDeviceCode {
		if tt == fosite.AccessToken && c.DeviceAuthorizationGrantAccessTokenLifespan.Valid {
			cl = &c.DeviceAuthorizationGrantAccessTokenLifespan.Duration
		} else if tt == fosite.IDToken && c.DeviceAuthorizationGrantIDTokenLifespan.Valid {
			cl = &c.DeviceAuthorizationGrantIDTokenLifespan.Duration
		} else if tt == fosite.RefreshToken && c.DeviceAuthorizationGrantRefreshTokenLifespan.Valid {
			cl = &c.DeviceAuthorizationGrantRefreshTokenLifespan.Duration
		}
	}

	if cl == nil {
		return fallback
	}
	return *cl
}

func (c *Client) GetAccessTokenStrategy() config.AccessTokenStrategyType {
	// We ignore the error here, because the empty string will default to
	// the global access token strategy.
	s, _ := config.ToAccessTokenStrategyType(c.AccessTokenStrategy)
	return s
}

func AccessTokenStrategySource(client fosite.Client) config.AccessTokenStrategySource {
	if source, ok := client.(config.AccessTokenStrategySource); ok {
		return source
	}
	return nil
}

func (c *Client) CookieSuffix() string {
	return CookieSuffix(c)
}

type IDer interface{ GetID() string }

func CookieSuffix(client IDer) string {
	return strconv.Itoa(int(murmur3.Sum32([]byte(client.GetID()))))
}
