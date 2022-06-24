// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// OAuth2Client Client represents an OAuth 2.0 Client.
//
// swagger:model oAuth2Client
type OAuth2Client struct {

	// allowed cors origins
	AllowedCorsOrigins StringSliceJSONFormat `json:"allowed_cors_origins,omitempty"`

	// audience
	Audience StringSliceJSONFormat `json:"audience,omitempty"`

	// authorization code grant access token lifespan
	AuthorizationCodeGrantAccessTokenLifespan NullDuration `json:"authorization_code_grant_access_token_lifespan,omitempty"`

	// authorization code grant id token lifespan
	AuthorizationCodeGrantIDTokenLifespan NullDuration `json:"authorization_code_grant_id_token_lifespan,omitempty"`

	// authorization code grant refresh token lifespan
	AuthorizationCodeGrantRefreshTokenLifespan NullDuration `json:"authorization_code_grant_refresh_token_lifespan,omitempty"`

	// Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout
	// Token to identify the RP session with the OP when the backchannel_logout_uri is used.
	// If omitted, the default value is false.
	BackchannelLogoutSessionRequired bool `json:"backchannel_logout_session_required,omitempty"`

	// RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.
	BackchannelLogoutURI string `json:"backchannel_logout_uri,omitempty"`

	// ID is the id for this client.
	//
	// The ID is autogenerated and immutable.
	ClientID string `json:"client_id,omitempty"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	ClientName string `json:"client_name,omitempty"`

	// Secret is the client's secret. The secret will be included in the create request as cleartext, and then
	// never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users
	// that they need to write the secret down as it will not be made available again.
	ClientSecret string `json:"client_secret,omitempty"`

	// SecretExpiresAt is an integer holding the time at which the client
	// secret will expire or 0 if it will not expire. The time is
	// represented as the number of seconds from 1970-01-01T00:00:00Z as
	// measured in UTC until the date/time of expiration.
	//
	// This feature is currently not supported and it's value will always
	// be set to 0.
	ClientSecretExpiresAt int64 `json:"client_secret_expires_at,omitempty"`

	// ClientURI is an URL string of a web page providing information about the client.
	// If present, the server SHOULD display this URL to the end-user in
	// a clickable fashion.
	ClientURI string `json:"client_uri,omitempty"`

	// contacts
	Contacts StringSliceJSONFormat `json:"contacts,omitempty"`

	// CreatedAt returns the timestamp of the client's creation.
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be
	// included to identify the RP session with the OP when the frontchannel_logout_uri is used.
	// If omitted, the default value is false.
	FrontchannelLogoutSessionRequired bool `json:"frontchannel_logout_session_required,omitempty"`

	// RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query
	// parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the
	// request and to determine which of the potentially multiple sessions is to be logged out; if either is
	// included, both MUST be.
	FrontchannelLogoutURI string `json:"frontchannel_logout_uri,omitempty"`

	// grant types
	GrantTypes StringSliceJSONFormat `json:"grant_types,omitempty"`

	// implicit grant access token lifespan
	ImplicitGrantAccessTokenLifespan NullDuration `json:"implicit_grant_access_token_lifespan,omitempty"`

	// implicit grant id token lifespan
	ImplicitGrantIDTokenLifespan NullDuration `json:"implicit_grant_id_token_lifespan,omitempty"`

	// jwks
	Jwks JoseJSONWebKeySet `json:"jwks,omitempty"`

	// URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains
	// the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the
	// Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing
	// and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced
	// JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both
	// signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used
	// to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST
	// match those in the certificate.
	JwksURI string `json:"jwks_uri,omitempty"`

	// jwt bearer grant access token lifespan
	JwtBearerGrantAccessTokenLifespan NullDuration `json:"jwt_bearer_grant_access_token_lifespan,omitempty"`

	// LogoURI is an URL string that references a logo for the client.
	LogoURI string `json:"logo_uri,omitempty"`

	// metadata
	Metadata JSONRawMessage `json:"metadata,omitempty"`

	// Owner is a string identifying the owner of the OAuth 2.0 Client.
	Owner string `json:"owner,omitempty"`

	// password grant access token lifespan
	PasswordGrantAccessTokenLifespan NullDuration `json:"password_grant_access_token_lifespan,omitempty"`

	// password grant refresh token lifespan
	PasswordGrantRefreshTokenLifespan NullDuration `json:"password_grant_refresh_token_lifespan,omitempty"`

	// PolicyURI is a URL string that points to a human-readable privacy policy document
	// that describes how the deployment organization collects, uses,
	// retains, and discloses personal data.
	PolicyURI string `json:"policy_uri,omitempty"`

	// post logout redirect uris
	PostLogoutRedirectUris StringSliceJSONFormat `json:"post_logout_redirect_uris,omitempty"`

	// redirect uris
	RedirectUris StringSliceJSONFormat `json:"redirect_uris,omitempty"`

	// refresh token grant access token lifespan
	RefreshTokenGrantAccessTokenLifespan NullDuration `json:"refresh_token_grant_access_token_lifespan,omitempty"`

	// refresh token grant id token lifespan
	RefreshTokenGrantIDTokenLifespan NullDuration `json:"refresh_token_grant_id_token_lifespan,omitempty"`

	// refresh token grant refresh token lifespan
	RefreshTokenGrantRefreshTokenLifespan NullDuration `json:"refresh_token_grant_refresh_token_lifespan,omitempty"`

	// RegistrationAccessToken can be used to update, get, or delete the OAuth2 Client.
	RegistrationAccessToken string `json:"registration_access_token,omitempty"`

	// RegistrationClientURI is the URL used to update, get, or delete the OAuth2 Client.
	RegistrationClientURI string `json:"registration_client_uri,omitempty"`

	// JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects
	// from this Client MUST be rejected, if not signed with this algorithm.
	RequestObjectSigningAlg string `json:"request_object_signing_alg,omitempty"`

	// request uris
	RequestUris StringSliceJSONFormat `json:"request_uris,omitempty"`

	// response types
	ResponseTypes StringSliceJSONFormat `json:"response_types,omitempty"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	// Example: scope1 scope-2 scope.3 scope:4
	Scope string `json:"scope,omitempty"`

	// URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a
	// file with a single JSON array of redirect_uri values.
	SectorIdentifierURI string `json:"sector_identifier_uri,omitempty"`

	// SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a
	// list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.
	SubjectType string `json:"subject_type,omitempty"`

	// Requested Client Authentication method for the Token Endpoint. The options are client_secret_post,
	// client_secret_basic, private_key_jwt, and none.
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty"`

	// Requested Client Authentication signing algorithm for the Token Endpoint.
	TokenEndpointAuthSigningAlg string `json:"token_endpoint_auth_signing_alg,omitempty"`

	// TermsOfServiceURI is a URL string that points to a human-readable terms of service
	// document for the client that describes a contractual relationship
	// between the end-user and the client that the end-user accepts when
	// authorizing the client.
	TosURI string `json:"tos_uri,omitempty"`

	// UpdatedAt returns the timestamp of the last update.
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updated_at,omitempty"`

	// JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT
	// [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims
	// as a UTF-8 encoded JSON object using the application/json content-type.
	UserinfoSignedResponseAlg string `json:"userinfo_signed_response_alg,omitempty"`
}

// Validate validates this o auth2 client
func (m *OAuth2Client) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAllowedCorsOrigins(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAudience(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAuthorizationCodeGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAuthorizationCodeGrantIDTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAuthorizationCodeGrantRefreshTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateClientCredentialsGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateContacts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGrantTypes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImplicitGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImplicitGrantIDTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateJwtBearerGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePasswordGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePasswordGrantRefreshTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePostLogoutRedirectUris(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRedirectUris(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRefreshTokenGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRefreshTokenGrantIDTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRefreshTokenGrantRefreshTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRequestUris(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResponseTypes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OAuth2Client) validateAllowedCorsOrigins(formats strfmt.Registry) error {
	if swag.IsZero(m.AllowedCorsOrigins) { // not required
		return nil
	}

	if err := m.AllowedCorsOrigins.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("allowed_cors_origins")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateAudience(formats strfmt.Registry) error {
	if swag.IsZero(m.Audience) { // not required
		return nil
	}

	if err := m.Audience.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("audience")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateAuthorizationCodeGrantAccessTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.AuthorizationCodeGrantAccessTokenLifespan) { // not required
		return nil
	}

	if err := m.AuthorizationCodeGrantAccessTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateAuthorizationCodeGrantIDTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.AuthorizationCodeGrantIDTokenLifespan) { // not required
		return nil
	}

	if err := m.AuthorizationCodeGrantIDTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateAuthorizationCodeGrantRefreshTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.AuthorizationCodeGrantRefreshTokenLifespan) { // not required
		return nil
	}

	if err := m.AuthorizationCodeGrantRefreshTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateClientCredentialsGrantAccessTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.ClientCredentialsGrantAccessTokenLifespan) { // not required
		return nil
	}

	if err := m.ClientCredentialsGrantAccessTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("client_credentials_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateContacts(formats strfmt.Registry) error {
	if swag.IsZero(m.Contacts) { // not required
		return nil
	}

	if err := m.Contacts.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("contacts")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *OAuth2Client) validateGrantTypes(formats strfmt.Registry) error {
	if swag.IsZero(m.GrantTypes) { // not required
		return nil
	}

	if err := m.GrantTypes.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("grant_types")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateImplicitGrantAccessTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.ImplicitGrantAccessTokenLifespan) { // not required
		return nil
	}

	if err := m.ImplicitGrantAccessTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("implicit_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateImplicitGrantIDTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.ImplicitGrantIDTokenLifespan) { // not required
		return nil
	}

	if err := m.ImplicitGrantIDTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("implicit_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateJwtBearerGrantAccessTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.JwtBearerGrantAccessTokenLifespan) { // not required
		return nil
	}

	if err := m.JwtBearerGrantAccessTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("jwt_bearer_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validatePasswordGrantAccessTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.PasswordGrantAccessTokenLifespan) { // not required
		return nil
	}

	if err := m.PasswordGrantAccessTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("password_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validatePasswordGrantRefreshTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.PasswordGrantRefreshTokenLifespan) { // not required
		return nil
	}

	if err := m.PasswordGrantRefreshTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("password_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validatePostLogoutRedirectUris(formats strfmt.Registry) error {
	if swag.IsZero(m.PostLogoutRedirectUris) { // not required
		return nil
	}

	if err := m.PostLogoutRedirectUris.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("post_logout_redirect_uris")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateRedirectUris(formats strfmt.Registry) error {
	if swag.IsZero(m.RedirectUris) { // not required
		return nil
	}

	if err := m.RedirectUris.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("redirect_uris")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateRefreshTokenGrantAccessTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.RefreshTokenGrantAccessTokenLifespan) { // not required
		return nil
	}

	if err := m.RefreshTokenGrantAccessTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateRefreshTokenGrantIDTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.RefreshTokenGrantIDTokenLifespan) { // not required
		return nil
	}

	if err := m.RefreshTokenGrantIDTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateRefreshTokenGrantRefreshTokenLifespan(formats strfmt.Registry) error {
	if swag.IsZero(m.RefreshTokenGrantRefreshTokenLifespan) { // not required
		return nil
	}

	if err := m.RefreshTokenGrantRefreshTokenLifespan.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateRequestUris(formats strfmt.Registry) error {
	if swag.IsZero(m.RequestUris) { // not required
		return nil
	}

	if err := m.RequestUris.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("request_uris")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateResponseTypes(formats strfmt.Registry) error {
	if swag.IsZero(m.ResponseTypes) { // not required
		return nil
	}

	if err := m.ResponseTypes.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("response_types")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) validateUpdatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("updated_at", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this o auth2 client based on the context it is used
func (m *OAuth2Client) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAllowedCorsOrigins(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAudience(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAuthorizationCodeGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAuthorizationCodeGrantIDTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAuthorizationCodeGrantRefreshTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateClientCredentialsGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateContacts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateGrantTypes(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateImplicitGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateImplicitGrantIDTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateJwtBearerGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePasswordGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePasswordGrantRefreshTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePostLogoutRedirectUris(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRedirectUris(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRefreshTokenGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRefreshTokenGrantIDTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRefreshTokenGrantRefreshTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRequestUris(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateResponseTypes(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OAuth2Client) contextValidateAllowedCorsOrigins(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AllowedCorsOrigins.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("allowed_cors_origins")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateAudience(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Audience.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("audience")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateAuthorizationCodeGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AuthorizationCodeGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateAuthorizationCodeGrantIDTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AuthorizationCodeGrantIDTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateAuthorizationCodeGrantRefreshTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AuthorizationCodeGrantRefreshTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateClientCredentialsGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ClientCredentialsGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("client_credentials_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateContacts(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Contacts.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("contacts")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateGrantTypes(ctx context.Context, formats strfmt.Registry) error {

	if err := m.GrantTypes.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("grant_types")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateImplicitGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ImplicitGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("implicit_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateImplicitGrantIDTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ImplicitGrantIDTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("implicit_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateJwtBearerGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.JwtBearerGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("jwt_bearer_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidatePasswordGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.PasswordGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("password_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidatePasswordGrantRefreshTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.PasswordGrantRefreshTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("password_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidatePostLogoutRedirectUris(ctx context.Context, formats strfmt.Registry) error {

	if err := m.PostLogoutRedirectUris.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("post_logout_redirect_uris")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateRedirectUris(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RedirectUris.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("redirect_uris")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateRefreshTokenGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RefreshTokenGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateRefreshTokenGrantIDTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RefreshTokenGrantIDTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateRefreshTokenGrantRefreshTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RefreshTokenGrantRefreshTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateRequestUris(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RequestUris.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("request_uris")
		}
		return err
	}

	return nil
}

func (m *OAuth2Client) contextValidateResponseTypes(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ResponseTypes.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("response_types")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OAuth2Client) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OAuth2Client) UnmarshalBinary(b []byte) error {
	var res OAuth2Client
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
