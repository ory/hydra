// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"encoding/json"
	stderr "errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/language"

	"github.com/ory/hydra/v2/fosite/i18n"
	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidatedAuthorizeCode is an error indicating that an authorization code has been
	// used previously.
	ErrInvalidatedAuthorizeCode = stderr.New("Authorization code has ben invalidated")
	// ErrInvalidatedDeviceCode is an error indicating that a device code has benn used previously.
	ErrInvalidatedDeviceCode = stderr.New("Device code has been invalidated")
	// ErrExistingUserCodeSignature is an error indicating that a row already exists with the provided user_code signature.
	ErrExistingUserCodeSignature = stderr.New("User code signature already exists in the database")
	// ErrSerializationFailure is an error indicating that the transactional capable storage could not guarantee
	// consistency of Update & Delete operations on the same rows between multiple sessions.
	ErrSerializationFailure = &RFC6749Error{
		ErrorField:       errUnknownErrorName,
		DescriptionField: "The request could not be completed because another request is competing for the same resource.",
		CodeField:        http.StatusConflict,
	}
	ErrUnknownRequest = &RFC6749Error{
		ErrorField:       errUnknownErrorName,
		DescriptionField: "The handler is not responsible for this request.",
		CodeField:        http.StatusBadRequest,
	}
	ErrRequestForbidden = &RFC6749Error{
		ErrorField:       errRequestForbidden,
		DescriptionField: "The request is not allowed.",
		HintField:        "You are not allowed to perform this action.",
		CodeField:        http.StatusForbidden,
	}
	ErrInvalidRequest = &RFC6749Error{
		ErrorField:       errInvalidRequestName,
		DescriptionField: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
		HintField:        "Make sure that the various parameters are correct, be aware of case sensitivity and trim your parameters. Make sure that the client you are using has exactly whitelisted the redirect_uri you specified.",
		CodeField:        http.StatusBadRequest,
	}
	ErrUnauthorizedClient = &RFC6749Error{
		ErrorField:       errUnauthorizedClientName,
		DescriptionField: "The client is not authorized to request a token using this method.",
		HintField:        "Make sure that client id and secret are correctly specified and that the client exists.",
		CodeField:        http.StatusBadRequest,
	}
	ErrAccessDenied = &RFC6749Error{
		ErrorField:       errAccessDeniedName,
		DescriptionField: "The resource owner or authorization server denied the request.",
		HintField:        "Make sure that the request you are making is valid. Maybe the credential or request parameters you are using are limited in scope or otherwise restricted.",
		CodeField:        http.StatusForbidden,
	}
	ErrUnsupportedResponseType = &RFC6749Error{
		ErrorField:       errUnsupportedResponseTypeName,
		DescriptionField: "The authorization server does not support obtaining a token using this method.",
		CodeField:        http.StatusBadRequest,
	}
	ErrUnsupportedResponseMode = &RFC6749Error{
		ErrorField:       errUnsupportedResponseModeName,
		DescriptionField: "The authorization server does not support obtaining a response using this response mode.",
		CodeField:        http.StatusBadRequest,
	}
	ErrInvalidScope = &RFC6749Error{
		ErrorField:       errInvalidScopeName,
		DescriptionField: "The requested scope is invalid, unknown, or malformed.",
		CodeField:        http.StatusBadRequest,
	}
	ErrServerError = &RFC6749Error{
		ErrorField:       errServerErrorName,
		DescriptionField: "The authorization server encountered an unexpected condition that prevented it from fulfilling the request.",
		CodeField:        http.StatusInternalServerError,
	}
	ErrTemporarilyUnavailable = &RFC6749Error{
		ErrorField:       errTemporarilyUnavailableName,
		DescriptionField: "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
		CodeField:        http.StatusServiceUnavailable,
	}
	ErrUnsupportedGrantType = &RFC6749Error{
		ErrorField:       errUnsupportedGrantTypeName,
		DescriptionField: "The authorization grant type is not supported by the authorization server.",
		CodeField:        http.StatusBadRequest,
	}
	ErrInvalidGrant = &RFC6749Error{
		ErrorField:       errInvalidGrantName,
		DescriptionField: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
		CodeField:        http.StatusBadRequest,
	}
	ErrInvalidClient = &RFC6749Error{
		ErrorField:       errInvalidClientName,
		DescriptionField: "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method).",
		CodeField:        http.StatusUnauthorized,
	}
	ErrInvalidState = &RFC6749Error{
		ErrorField:       errInvalidStateName,
		DescriptionField: "The state is missing or does not have enough characters and is therefore considered too weak.",
		CodeField:        http.StatusBadRequest,
	}
	ErrMisconfiguration = &RFC6749Error{
		ErrorField:       errMisconfigurationName,
		DescriptionField: "The request failed because of an internal error that is probably caused by misconfiguration.",
		CodeField:        http.StatusInternalServerError,
	}
	ErrInsufficientEntropy = &RFC6749Error{
		ErrorField:       errInsufficientEntropyName,
		DescriptionField: "The request used a security parameter (e.g., anti-replay, anti-csrf) with insufficient entropy.",
		CodeField:        http.StatusBadRequest,
	}
	ErrNotFound = &RFC6749Error{
		ErrorField:       errNotFoundName,
		DescriptionField: "Could not find the requested resource(s).",
		CodeField:        http.StatusNotFound,
	}
	ErrRequestUnauthorized = &RFC6749Error{
		ErrorField:       errRequestUnauthorizedName,
		DescriptionField: "The request could not be authorized.",
		HintField:        "Check that you provided valid credentials in the right format.",
		CodeField:        http.StatusUnauthorized,
	}
	ErrTokenSignatureMismatch = &RFC6749Error{
		ErrorField:       errTokenSignatureMismatchName,
		DescriptionField: "Token signature mismatch.",
		HintField:        "Check that you provided  a valid token in the right format.",
		CodeField:        http.StatusBadRequest,
	}
	ErrInvalidTokenFormat = &RFC6749Error{
		ErrorField:       errInvalidTokenFormatName,
		DescriptionField: "Invalid token format.",
		HintField:        "Check that you provided a valid token in the right format.",
		CodeField:        http.StatusBadRequest,
	}
	ErrTokenExpired = &RFC6749Error{
		ErrorField:       errTokenExpiredName,
		DescriptionField: "Token expired.",
		HintField:        "The token expired.",
		CodeField:        http.StatusUnauthorized,
	}
	ErrScopeNotGranted = &RFC6749Error{
		ErrorField:       errScopeNotGrantedName,
		DescriptionField: "The token was not granted the requested scope.",
		HintField:        "The resource owner did not grant the requested scope.",
		CodeField:        http.StatusForbidden,
	}
	ErrTokenClaim = &RFC6749Error{
		ErrorField:       errTokenClaimName,
		DescriptionField: "The token failed validation due to a claim mismatch.",
		HintField:        "One or more token claims failed validation.",
		CodeField:        http.StatusUnauthorized,
	}
	ErrInactiveToken = &RFC6749Error{
		ErrorField:       errTokenInactiveName,
		DescriptionField: "Token is inactive because it is malformed, expired or otherwise invalid.",
		HintField:        "Token validation failed.",
		CodeField:        http.StatusUnauthorized,
	}
	ErrLoginRequired = &RFC6749Error{
		ErrorField:       errLoginRequired,
		DescriptionField: "The Authorization Server requires End-User authentication.",
		CodeField:        http.StatusBadRequest,
	}
	ErrInteractionRequired = &RFC6749Error{
		DescriptionField: "The Authorization Server requires End-User interaction of some form to proceed.",
		ErrorField:       errInteractionRequired,
		CodeField:        http.StatusBadRequest,
	}
	ErrConsentRequired = &RFC6749Error{
		DescriptionField: "The Authorization Server requires End-User consent.",
		ErrorField:       errConsentRequired,
		CodeField:        http.StatusBadRequest,
	}
	ErrRequestNotSupported = &RFC6749Error{
		DescriptionField: "The OP does not support use of the request parameter.",
		ErrorField:       errRequestNotSupportedName,
		CodeField:        http.StatusBadRequest,
	}
	ErrRequestURINotSupported = &RFC6749Error{
		DescriptionField: "The OP does not support use of the request_uri parameter.",
		ErrorField:       errRequestURINotSupportedName,
		CodeField:        http.StatusBadRequest,
	}
	ErrRegistrationNotSupported = &RFC6749Error{
		DescriptionField: "The OP does not support use of the registration parameter.",
		ErrorField:       errRegistrationNotSupportedName,
		CodeField:        http.StatusBadRequest,
	}
	ErrInvalidRequestURI = &RFC6749Error{
		DescriptionField: "The request_uri in the Authorization Request returns an error or contains invalid data.",
		ErrorField:       errInvalidRequestURI,
		CodeField:        http.StatusBadRequest,
	}
	ErrInvalidRequestObject = &RFC6749Error{
		DescriptionField: "The request parameter contains an invalid Request Object.",
		ErrorField:       errInvalidRequestObject,
		CodeField:        http.StatusBadRequest,
	}
	ErrJTIKnown = &RFC6749Error{
		DescriptionField: "The jti was already used.",
		ErrorField:       errJTIKnownName,
		CodeField:        http.StatusBadRequest,
	}
	ErrAuthorizationPending = &RFC6749Error{
		DescriptionField: "The authorization request is still pending as the end user hasn't yet completed the user-interaction steps.",
		ErrorField:       errAuthorizationPending,
		CodeField:        http.StatusBadRequest,
	}
	ErrSlowDown = &RFC6749Error{
		DescriptionField: "The authorization request was rate-limited to prevent system overload.",
		HintField:        "Ensure that you don't call the token endpoint sooner than the polling interval",
		ErrorField:       errSlowDown,
		CodeField:        http.StatusBadRequest,
	}
	ErrDeviceExpiredToken = &RFC6749Error{
		DescriptionField: "The device_code has expired, and the device authorization session has concluded.",
		ErrorField:       errDeviceExpiredToken,
		CodeField:        http.StatusBadRequest,
	}
)

const (
	errInvalidRequestURI           = "invalid_request_uri"
	errInvalidRequestObject        = "invalid_request_object"
	errConsentRequired             = "consent_required"
	errInteractionRequired         = "interaction_required"
	errLoginRequired               = "login_required"
	errRequestUnauthorizedName     = "request_unauthorized"
	errRequestForbidden            = "request_forbidden"
	errInvalidRequestName          = "invalid_request"
	errUnauthorizedClientName      = "unauthorized_client"
	errAccessDeniedName            = "access_denied"
	errUnsupportedResponseTypeName = "unsupported_response_type"
	errUnsupportedResponseModeName = "unsupported_response_mode"
	errInvalidScopeName            = "invalid_scope"
	errServerErrorName             = "server_error"
	errTemporarilyUnavailableName  = "temporarily_unavailable"
	errUnsupportedGrantTypeName    = "unsupported_grant_type"
	errInvalidGrantName            = "invalid_grant"
	errInvalidClientName           = "invalid_client"
	errNotFoundName                = "not_found"
	errInvalidStateName            = "invalid_state"
	errMisconfigurationName        = "misconfiguration"
	errInsufficientEntropyName     = "insufficient_entropy"
	errInvalidTokenFormatName      = "invalid_token"
	errTokenSignatureMismatchName  = "token_signature_mismatch"
	errTokenExpiredName            = "invalid_token" // https://tools.ietf.org/html/rfc6750#section-3.1
	errScopeNotGrantedName         = "scope_not_granted"
	errTokenClaimName              = "token_claim"
	errTokenInactiveName           = "token_inactive"
	// errAuthorizationCodeInactiveName = "authorization_code_inactive"
	errUnknownErrorName             = "error"
	errRequestNotSupportedName      = "request_not_supported"
	errRequestURINotSupportedName   = "request_uri_not_supported"
	errRegistrationNotSupportedName = "registration_not_supported"
	errJTIKnownName                 = "jti_known"
	errAuthorizationPending         = "authorization_pending"
	errSlowDown                     = "slow_down"
	errDeviceExpiredToken           = "expired_token"
)

type (
	RFC6749Error struct {
		ErrorField       string
		DescriptionField string
		HintField        string
		CodeField        int
		DebugField       string
		cause            error
		useLegacyFormat  bool
		exposeDebug      bool

		// Fields for globalization
		hintIDField string
		hintArgs    []interface{}
		catalog     i18n.MessageCatalog
		lang        language.Tag
	}
	stackTracer interface {
		StackTrace() errors.StackTrace
	}
)

var (
	_ errorsx.DebugCarrier      = new(RFC6749Error)
	_ errorsx.ReasonCarrier     = new(RFC6749Error)
	_ errorsx.RequestIDCarrier  = new(RFC6749Error)
	_ errorsx.StatusCarrier     = new(RFC6749Error)
	_ errorsx.StatusCodeCarrier = new(RFC6749Error)
	// _ errorsx.DetailsCarrier = new(RFC6749Error)
)

func ErrorToRFC6749Error(err error) *RFC6749Error {
	var e *RFC6749Error
	if errors.As(err, &e) {
		return e
	}
	return &RFC6749Error{
		ErrorField:       errUnknownErrorName,
		DescriptionField: "The error is unrecognizable",
		DebugField:       err.Error(),
		CodeField:        http.StatusInternalServerError,
		cause:            err,
	}
}

// StackTrace returns the error's stack trace.
func (e *RFC6749Error) StackTrace() (trace errors.StackTrace) {
	if e.cause == e || e.cause == nil {
		return
	}

	if st := stackTracer(nil); stderr.As(e.cause, &st) {
		trace = st.StackTrace()
	}

	return
}

func (e RFC6749Error) Unwrap() error {
	return e.cause
}

func (e *RFC6749Error) Wrap(err error) {
	e.cause = err
}

func (e RFC6749Error) WithWrap(cause error) *RFC6749Error {
	e.cause = cause
	return &e
}

func (e RFC6749Error) WithLegacyFormat(useLegacyFormat bool) *RFC6749Error {
	e.useLegacyFormat = useLegacyFormat
	return &e
}

func (e *RFC6749Error) WithTrace(err error) *RFC6749Error {
	if st := stackTracer(nil); !stderr.As(e.cause, &st) {
		e.Wrap(errorsx.WithStack(err))
	} else {
		e.Wrap(err)
	}
	return e
}

func (e RFC6749Error) Is(err error) bool {
	switch te := err.(type) {
	case RFC6749Error:
		return e.ErrorField == te.ErrorField &&
			e.CodeField == te.CodeField
	case *RFC6749Error:
		return e.ErrorField == te.ErrorField &&
			e.CodeField == te.CodeField
	}
	return false
}

func (e *RFC6749Error) Status() string {
	return http.StatusText(e.CodeField)
}

func (e RFC6749Error) Error() string {
	return e.ErrorField
}

func (e *RFC6749Error) RequestID() string {
	return ""
}

func (e *RFC6749Error) Reason() string {
	return e.HintField
}

func (e *RFC6749Error) StatusCode() int {
	return e.CodeField
}

func (e *RFC6749Error) Cause() error {
	return e.cause
}

func (e *RFC6749Error) WithHintf(hint string, args ...interface{}) *RFC6749Error {
	err := *e
	if err.hintIDField == "" {
		err.hintIDField = hint
	}

	err.hintArgs = args
	err.HintField = fmt.Sprintf(hint, args...)
	return &err
}

func (e *RFC6749Error) WithHint(hint string) *RFC6749Error {
	err := *e
	if err.hintIDField == "" {
		err.hintIDField = hint
	}

	err.HintField = hint
	return &err
}

// WithHintIDOrDefaultf accepts the ID of the hint message
func (e *RFC6749Error) WithHintIDOrDefaultf(ID string, def string, args ...interface{}) *RFC6749Error {
	err := *e
	err.hintIDField = ID
	err.hintArgs = args
	err.HintField = fmt.Sprintf(def, args...)
	return &err
}

// WithHintTranslationID accepts the ID of the hint message and should be paired with
// WithHint and WithHintf to add a default message and vaargs.
func (e *RFC6749Error) WithHintTranslationID(ID string) *RFC6749Error {
	err := *e
	err.hintIDField = ID
	return &err
}

func (e *RFC6749Error) Debug() string {
	return e.DebugField
}

func (e *RFC6749Error) WithDebug(debug string) *RFC6749Error {
	err := *e
	err.DebugField = debug
	return &err
}

func (e *RFC6749Error) WithDebugf(debug string, args ...interface{}) *RFC6749Error {
	return e.WithDebug(fmt.Sprintf(debug, args...))
}

func (e *RFC6749Error) WithDescription(description string) *RFC6749Error {
	err := *e
	err.DescriptionField = description
	return &err
}

func (e *RFC6749Error) WithLocalizer(catalog i18n.MessageCatalog, lang language.Tag) *RFC6749Error {
	err := *e
	err.catalog = catalog
	err.lang = lang
	return &err
}

// Sanitize strips the debug field
//
// Deprecated: Use WithExposeDebug instead.
func (e *RFC6749Error) Sanitize() *RFC6749Error {
	err := *e
	err.DebugField = ""
	return &err
}

// WithExposeDebug if set to true exposes debug messages
func (e *RFC6749Error) WithExposeDebug(exposeDebug bool) *RFC6749Error {
	err := *e
	err.exposeDebug = exposeDebug
	return &err
}

// GetDescription returns a more description description, combined with hint and debug (when available).
func (e *RFC6749Error) GetDescription() string {
	description := i18n.GetMessageOrDefault(e.catalog, e.ErrorField, e.lang, e.DescriptionField)
	e.computeHintField()
	if e.HintField != "" {
		description += " " + e.HintField
	}
	if e.DebugField != "" && e.exposeDebug {
		description += " " + e.DebugField
	}
	return strings.ReplaceAll(description, "\"", "'")
}

// RFC6749ErrorJson is a helper struct for JSON encoding/decoding of RFC6749Error.
type RFC6749ErrorJson struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
	Hint        string `json:"error_hint,omitempty"`
	Code        int    `json:"status_code,omitempty"`
	Debug       string `json:"error_debug,omitempty"`
}

func (e *RFC6749Error) UnmarshalJSON(b []byte) error {
	var data RFC6749ErrorJson

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	e.ErrorField = data.Name
	e.CodeField = data.Code
	e.DescriptionField = data.Description

	if len(data.Hint+data.Debug) > 0 {
		e.HintField = data.Hint
		e.DebugField = data.Debug
		e.useLegacyFormat = true
	}

	return nil
}

func (e RFC6749Error) MarshalJSON() ([]byte, error) {
	if !e.useLegacyFormat {
		return json.Marshal(&RFC6749ErrorJson{
			Name:        e.ErrorField,
			Description: e.GetDescription(),
		})
	}

	var debug string
	if e.exposeDebug {
		debug = e.DebugField
	}

	return json.Marshal(&RFC6749ErrorJson{
		Name:        e.ErrorField,
		Description: e.DescriptionField,
		Hint:        e.HintField,
		Code:        e.CodeField,
		Debug:       debug,
	})
}

func (e *RFC6749Error) ToValues() url.Values {
	values := url.Values{}
	values.Set("error", e.ErrorField)
	values.Set("error_description", e.GetDescription())

	if e.useLegacyFormat {
		values.Set("error_description", e.DescriptionField)
		if e.HintField != "" {
			values.Set("error_hint", e.HintField)
		}

		if e.DebugField != "" && e.exposeDebug {
			values.Set("error_debug", e.DebugField)
		}
	}

	return values
}

func (e *RFC6749Error) computeHintField() {
	if e.hintIDField == "" {
		return
	}

	e.HintField = i18n.GetMessageOrDefault(e.catalog, e.hintIDField, e.lang, e.HintField, e.hintArgs...)
}
