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
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidatedAuthorizeCode is an error indicating that an authorization code has been
	// used previously.
	ErrInvalidatedAuthorizeCode = errors.New("Authorization code has ben invalidated")
	ErrUnknownRequest           = &RFC6749Error{
		Name:        errUnknownErrorName,
		Description: "The handler is not responsible for this request",
		Code:        http.StatusBadRequest,
	}
	ErrRequestForbidden = &RFC6749Error{
		Name:        errRequestForbidden,
		Description: "The request is not allowed",
		Hint:        "You are not allowed to perform this action.",
		Code:        http.StatusForbidden,
	}
	ErrInvalidRequest = &RFC6749Error{
		Name:        errInvalidRequestName,
		Description: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
		Hint:        "Make sure that the various parameters are correct, be aware of case sensitivity and trim your parameters. Make sure that the client you are using has exactly whitelisted the redirect_uri you specified.",
		Code:        http.StatusBadRequest,
	}
	ErrUnauthorizedClient = &RFC6749Error{
		Name:        errUnauthorizedClientName,
		Description: "The client is not authorized to request a token using this method",
		Hint:        "Make sure that client id and secret are correctly specified and that the client exists.",
		Code:        http.StatusBadRequest,
	}
	ErrAccessDenied = &RFC6749Error{
		Name:        errAccessDeniedName,
		Description: "The resource owner or authorization server denied the request",
		Hint:        "Make sure that the request you are making is valid. Maybe the credential or request parameters you are using are limited in scope or otherwise restricted.",
		Code:        http.StatusForbidden,
	}
	ErrUnsupportedResponseType = &RFC6749Error{
		Name:        errUnsupportedResponseTypeName,
		Description: "The authorization server does not support obtaining a token using this method",
		Code:        http.StatusBadRequest,
	}
	ErrInvalidScope = &RFC6749Error{
		Name:        errInvalidScopeName,
		Description: "The requested scope is invalid, unknown, or malformed",
		Code:        http.StatusBadRequest,
	}
	ErrServerError = &RFC6749Error{
		Name:        errServerErrorName,
		Description: "The authorization server encountered an unexpected condition that prevented it from fulfilling the request",
		Code:        http.StatusInternalServerError,
	}
	ErrTemporarilyUnavailable = &RFC6749Error{
		Name:        errTemporarilyUnavailableName,
		Description: "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server",
		Code:        http.StatusServiceUnavailable,
	}
	ErrUnsupportedGrantType = &RFC6749Error{
		Name:        errUnsupportedGrantTypeName,
		Description: "The authorization grant type is not supported by the authorization server",
		Code:        http.StatusBadRequest,
	}
	ErrInvalidGrant = &RFC6749Error{
		Name:        errInvalidGrantName,
		Description: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client",
		Code:        http.StatusBadRequest,
	}
	ErrInvalidClient = &RFC6749Error{
		Name:        errInvalidClientName,
		Description: "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method)",
		Code:        http.StatusUnauthorized,
	}
	ErrInvalidState = &RFC6749Error{
		Name:        errInvalidStateName,
		Description: fmt.Sprintf("The state is missing or has less than %d characters and is therefore considered too weak", MinParameterEntropy),
		Code:        http.StatusBadRequest,
	}
	ErrMisconfiguration = &RFC6749Error{
		Name:        errMisconfigurationName,
		Description: "The request failed because of an internal error that is probably caused by misconfiguration",
		Code:        http.StatusInternalServerError,
	}
	ErrInsufficientEntropy = &RFC6749Error{
		Name:        errInsufficientEntropyName,
		Description: fmt.Sprintf("The request used a security parameter (e.g., anti-replay, anti-csrf) with insufficient entropy (minimum of %d characters)", MinParameterEntropy),
		Code:        http.StatusBadRequest,
	}
	ErrNotFound = &RFC6749Error{
		Name:        errNotFoundName,
		Description: "Could not find the requested resource(s)",
		Code:        http.StatusNotFound,
	}
	ErrRequestUnauthorized = &RFC6749Error{
		Name:        errRequestUnauthorizedName,
		Description: "The request could not be authorized",
		Hint:        "Check that you provided valid credentials in the right format.",
		Code:        http.StatusUnauthorized,
	}
	ErrTokenSignatureMismatch = &RFC6749Error{
		Name:        errTokenSignatureMismatchName,
		Description: "Token signature mismatch",
		Hint:        "Check that you provided  a valid token in the right format.",
		Code:        http.StatusBadRequest,
	}
	ErrInvalidTokenFormat = &RFC6749Error{
		Name:        errInvalidTokenFormatName,
		Description: "Invalid token format",
		Hint:        "Check that you provided a valid token in the right format.",
		Code:        http.StatusBadRequest,
	}
	ErrTokenExpired = &RFC6749Error{
		Name:        errTokenExpiredName,
		Description: "Token expired",
		Hint:        "The token expired.",
		Code:        http.StatusUnauthorized,
	}
	ErrScopeNotGranted = &RFC6749Error{
		Name:        errScopeNotGrantedName,
		Description: "The token was not granted the requested scope",
		Hint:        "The resource owner did not grant the requested scope.",
		Code:        http.StatusForbidden,
	}
	ErrTokenClaim = &RFC6749Error{
		Name:        errTokenClaimName,
		Description: "The token failed validation due to a claim mismatch",
		Hint:        "One or more token claims failed validation.",
		Code:        http.StatusUnauthorized,
	}
	ErrInactiveToken = &RFC6749Error{
		Name:        errTokenInactiveName,
		Description: "Token is inactive because it is malformed, expired or otherwise invalid",
		Hint:        "Token validation failed.",
		Code:        http.StatusUnauthorized,
	}
	ErrRevokationClientMismatch = &RFC6749Error{
		Name:        errRevokationClientMismatchName,
		Description: "Token was not issued to the client making the revokation request",
		Code:        http.StatusBadRequest,
	}
	ErrLoginRequired = &RFC6749Error{
		Name:        errLoginRequired,
		Description: "The Authorization Server requires End-User authentication",
		Code:        http.StatusBadRequest,
	}
	ErrInteractionRequired = &RFC6749Error{
		Description: "The Authorization Server requires End-User interaction of some form to proceed",
		Name:        errInteractionRequired,
		Code:        http.StatusBadRequest,
	}
	ErrConsentRequired = &RFC6749Error{
		Description: "The Authorization Server requires End-User consent",
		Name:        errConsentRequired,
		Code:        http.StatusBadRequest,
	}
	ErrRequestNotSupported = &RFC6749Error{
		Description: "The OP does not support use of the request parameter",
		Name:        errRequestNotSupportedName,
		Code:        http.StatusBadRequest,
	}
	ErrRequestURINotSupported = &RFC6749Error{
		Description: "The OP does not support use of the request_uri parameter",
		Name:        errRequestURINotSupportedName,
		Code:        http.StatusBadRequest,
	}
	ErrRegistrationNotSupported = &RFC6749Error{
		Description: "The OP does not support use of the registration parameter",
		Name:        errRegistrationNotSupportedName,
		Code:        http.StatusBadRequest,
	}
	ErrInvalidRequestURI = &RFC6749Error{
		Description: "The request_uri in the Authorization Request returns an error or contains invalid data. ",
		Name:        errInvalidRequestURI,
		Code:        http.StatusBadRequest,
	}
	ErrInvalidRequestObject = &RFC6749Error{
		Description: "The request parameter contains an invalid Request Object. ",
		Name:        errInvalidRequestObject,
		Code:        http.StatusBadRequest,
	}
)

const (
	errInvalidRequestURI            = "invalid_request_uri"
	errInvalidRequestObject         = "invalid_request_object"
	errConsentRequired              = "consent_required"
	errInteractionRequired          = "interaction_required"
	errLoginRequired                = "login_required"
	errRequestUnauthorizedName      = "request_unauthorized"
	errRequestForbidden             = "request_forbidden"
	errInvalidRequestName           = "invalid_request"
	errUnauthorizedClientName       = "unauthorized_client"
	errAccessDeniedName             = "access_denied"
	errUnsupportedResponseTypeName  = "unsupported_response_type"
	errInvalidScopeName             = "invalid_scope"
	errServerErrorName              = "server_error"
	errTemporarilyUnavailableName   = "temporarily_unavailable"
	errUnsupportedGrantTypeName     = "unsupported_grant_type"
	errInvalidGrantName             = "invalid_grant"
	errInvalidClientName            = "invalid_client"
	errNotFoundName                 = "not_found"
	errInvalidStateName             = "invalid_state"
	errMisconfigurationName         = "misconfiguration"
	errInsufficientEntropyName      = "insufficient_entropy"
	errInvalidTokenFormatName       = "invalid_token"
	errTokenSignatureMismatchName   = "token_signature_mismatch"
	errTokenExpiredName             = "token_expired"
	errScopeNotGrantedName          = "scope_not_granted"
	errTokenClaimName               = "token_claim"
	errTokenInactiveName            = "token_inactive"
	errAuthorizaionCodeInactiveName = "authorization_code_inactive"
	errUnknownErrorName             = "error"
	errRevokationClientMismatchName = "revokation_client_mismatch"
	errRequestNotSupportedName      = "request_not_supported"
	errRequestURINotSupportedName   = "request_uri_not_supported"
	errRegistrationNotSupportedName = "registration_not_supported"
)

func ErrorToRFC6749Error(err error) *RFC6749Error {
	if e, ok := err.(*RFC6749Error); ok {
		return e
	} else if e, ok := errors.Cause(err).(*RFC6749Error); ok {
		return e
	}
	return &RFC6749Error{
		Name:        errUnknownErrorName,
		Description: "The error is unrecognizable.",
		Debug:       err.Error(),
		Code:        http.StatusInternalServerError,
	}
}

type RFC6749Error struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
	Hint        string `json:"error_hint,omitempty"`
	Code        int    `json:"status_code,omitempty"`
	Debug       string `json:"error_debug,omitempty"`
}

func (e *RFC6749Error) Status() string {
	return http.StatusText(e.Code)
}

func (e *RFC6749Error) Error() string {
	return e.Name
}

func (e *RFC6749Error) RequestID() string {
	return ""
}

func (e *RFC6749Error) Reason() string {
	return e.Hint
}

func (e *RFC6749Error) Details() []map[string]interface{} {
	return []map[string]interface{}{}
}

func (e *RFC6749Error) StatusCode() int {
	return e.Code
}

func (e *RFC6749Error) WithHintf(hint string, args ...interface{}) *RFC6749Error {
	return e.WithHint(fmt.Sprintf(hint, args...))
}

func (e *RFC6749Error) WithHint(hint string) *RFC6749Error {
	err := *e
	err.Hint = hint
	return &err
}

func (e *RFC6749Error) WithDebug(debug string) *RFC6749Error {
	err := *e
	err.Debug = debug
	return &err
}

func (e *RFC6749Error) WithDebugf(debug string, args ...interface{}) *RFC6749Error {
	return e.WithDebug(fmt.Sprintf(debug, args...))
}

func (e *RFC6749Error) WithDescription(description string) *RFC6749Error {
	err := *e
	err.Description = description
	return &err
}
