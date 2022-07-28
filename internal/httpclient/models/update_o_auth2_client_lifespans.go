// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UpdateOAuth2ClientLifespans UpdateOAuth2ClientLifespans holds default lifespan configuration for the different
// token types that may be issued for the client. This configuration takes
// precedence over fosite's instance-wide default lifespan, but it may be
// overridden by a session's expires_at claim.
//
// The OIDC Hybrid grant type inherits token lifespan configuration from the implicit grant.
//
// swagger:model UpdateOAuth2ClientLifespans
type UpdateOAuth2ClientLifespans struct {

	// authorization code grant access token lifespan
	AuthorizationCodeGrantAccessTokenLifespan NullDuration `json:"authorization_code_grant_access_token_lifespan,omitempty"`

	// authorization code grant id token lifespan
	AuthorizationCodeGrantIDTokenLifespan NullDuration `json:"authorization_code_grant_id_token_lifespan,omitempty"`

	// authorization code grant refresh token lifespan
	AuthorizationCodeGrantRefreshTokenLifespan NullDuration `json:"authorization_code_grant_refresh_token_lifespan,omitempty"`

	// client credentials grant access token lifespan
	ClientCredentialsGrantAccessTokenLifespan NullDuration `json:"client_credentials_grant_access_token_lifespan,omitempty"`

	// implicit grant access token lifespan
	ImplicitGrantAccessTokenLifespan NullDuration `json:"implicit_grant_access_token_lifespan,omitempty"`

	// implicit grant id token lifespan
	ImplicitGrantIDTokenLifespan NullDuration `json:"implicit_grant_id_token_lifespan,omitempty"`

	// jwt bearer grant access token lifespan
	JwtBearerGrantAccessTokenLifespan NullDuration `json:"jwt_bearer_grant_access_token_lifespan,omitempty"`

	// password grant access token lifespan
	PasswordGrantAccessTokenLifespan NullDuration `json:"password_grant_access_token_lifespan,omitempty"`

	// password grant refresh token lifespan
	PasswordGrantRefreshTokenLifespan NullDuration `json:"password_grant_refresh_token_lifespan,omitempty"`

	// refresh token grant access token lifespan
	RefreshTokenGrantAccessTokenLifespan NullDuration `json:"refresh_token_grant_access_token_lifespan,omitempty"`

	// refresh token grant id token lifespan
	RefreshTokenGrantIDTokenLifespan NullDuration `json:"refresh_token_grant_id_token_lifespan,omitempty"`

	// refresh token grant refresh token lifespan
	RefreshTokenGrantRefreshTokenLifespan NullDuration `json:"refresh_token_grant_refresh_token_lifespan,omitempty"`
}

// Validate validates this update o auth2 client lifespans
func (m *UpdateOAuth2ClientLifespans) Validate(formats strfmt.Registry) error {
	var res []error

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

	if err := m.validateRefreshTokenGrantAccessTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRefreshTokenGrantIDTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRefreshTokenGrantRefreshTokenLifespan(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdateOAuth2ClientLifespans) validateAuthorizationCodeGrantAccessTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateAuthorizationCodeGrantIDTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateAuthorizationCodeGrantRefreshTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateClientCredentialsGrantAccessTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateImplicitGrantAccessTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateImplicitGrantIDTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateJwtBearerGrantAccessTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validatePasswordGrantAccessTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validatePasswordGrantRefreshTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateRefreshTokenGrantAccessTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateRefreshTokenGrantIDTokenLifespan(formats strfmt.Registry) error {
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

func (m *UpdateOAuth2ClientLifespans) validateRefreshTokenGrantRefreshTokenLifespan(formats strfmt.Registry) error {
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

// ContextValidate validate this update o auth2 client lifespans based on the context it is used
func (m *UpdateOAuth2ClientLifespans) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

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

	if err := m.contextValidateRefreshTokenGrantAccessTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRefreshTokenGrantIDTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRefreshTokenGrantRefreshTokenLifespan(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateAuthorizationCodeGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AuthorizationCodeGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateAuthorizationCodeGrantIDTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AuthorizationCodeGrantIDTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateAuthorizationCodeGrantRefreshTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.AuthorizationCodeGrantRefreshTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("authorization_code_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateClientCredentialsGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ClientCredentialsGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("client_credentials_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateImplicitGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ImplicitGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("implicit_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateImplicitGrantIDTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ImplicitGrantIDTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("implicit_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateJwtBearerGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.JwtBearerGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("jwt_bearer_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidatePasswordGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.PasswordGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("password_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidatePasswordGrantRefreshTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.PasswordGrantRefreshTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("password_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateRefreshTokenGrantAccessTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RefreshTokenGrantAccessTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_access_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateRefreshTokenGrantIDTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RefreshTokenGrantIDTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_id_token_lifespan")
		}
		return err
	}

	return nil
}

func (m *UpdateOAuth2ClientLifespans) contextValidateRefreshTokenGrantRefreshTokenLifespan(ctx context.Context, formats strfmt.Registry) error {

	if err := m.RefreshTokenGrantRefreshTokenLifespan.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("refresh_token_grant_refresh_token_lifespan")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UpdateOAuth2ClientLifespans) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateOAuth2ClientLifespans) UnmarshalBinary(b []byte) error {
	var res UpdateOAuth2ClientLifespans
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
