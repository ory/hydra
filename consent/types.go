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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/client"
)

const (
	consentRequestDeniedErrorName = "consent request denied"
	loginRequestDeniedErrorName   = "login request denied"
)

// The response payload sent when accepting or rejecting a login or consent request.
//
// swagger:model completedRequest
type RequestHandlerResponse struct {
	// RedirectURL is the URL which you should redirect the user to once the authentication process is completed.
	RedirectTo string `json:"redirect_to"`
}

// swagger:ignore
type LoginSession struct {
	ID              string    `db:"id"`
	AuthenticatedAt time.Time `db:"authenticated_at"`
	Subject         string    `db:"subject"`
	Remember        bool      `db:"remember"`
}

// The request payload used to accept a login or consent request.
//
// swagger:model rejectRequest
type RequestDeniedError struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
	Hint        string `json:"error_hint,omitempty"`
	Code        int    `json:"status_code,omitempty"`
	Debug       string `json:"error_debug,omitempty"`

	valid bool
}

func (e *RequestDeniedError) IsError() bool {
	return e != nil && e.valid
}

func (e *RequestDeniedError) SetDefaults(name string) {
	if e.Name == "" {
		e.Name = name
	}

	if e.Code == 0 {
		e.Code = http.StatusBadRequest
	}
}

func (e *RequestDeniedError) toRFCError() *fosite.RFC6749Error {
	if e.Name == "" {
		e.Name = "request was denied"
	}

	if e.Code == 0 {
		e.Code = fosite.ErrInvalidRequest.Code
	}

	return &fosite.RFC6749Error{
		Name:        e.Name,
		Description: e.Description,
		Hint:        e.Hint,
		Code:        e.Code,
		Debug:       e.Debug,
	}
}

func (e *RequestDeniedError) Scan(value interface{}) error {
	v := fmt.Sprintf("%s", value)
	if len(v) == 0 || v == "{}" {
		return nil
	}
	return errors.WithStack(json.Unmarshal([]byte(v), e))
}

func (e *RequestDeniedError) Value() (driver.Value, error) {
	if e == nil {
		return "{}", nil
	}

	value, err := json.Marshal(e)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return string(value), nil
}

// The request payload used to accept a consent request.
//
// swagger:model acceptConsentRequest
type HandledConsentRequest struct {
	// GrantScope sets the scope the user authorized the client to use. Should be a subset of `requested_scope`.
	GrantedScope sqlxx.StringSlicePipeDelimiter `json:"grant_scope" db:"granted_scope"`

	// GrantedAudience sets the audience the user authorized the client to use. Should be a subset of `requested_access_token_audience`.
	GrantedAudience sqlxx.StringSlicePipeDelimiter `json:"grant_access_token_audience" db:"granted_at_audience"`

	// Session allows you to set (optional) session data for access and ID tokens.
	Session *ConsentRequestSessionData `json:"session" db:"-"`

	// Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same
	// client asks the same user for the same, or a subset of, scope.
	Remember bool `json:"remember" db:"remember"`

	// RememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the
	// authorization will be remembered indefinitely.
	RememberFor int `json:"remember_for" db:"remember_for"`

	// HandledAt contains the timestamp the consent request was handled.
	HandledAt sqlxx.NullTime `json:"handled_at" db:"handled_at"`

	ConsentRequest  *ConsentRequest     `json:"-" db:"-"`
	Error           *RequestDeniedError `json:"-" db:"error"`
	Challenge       string              `json:"-" db:"challenge"`
	RequestedAt     time.Time           `json:"-" db:"requested_at"`
	AuthenticatedAt sqlxx.NullTime      `json:"-" db:"authenticated_at"`
	WasUsed         bool                `json:"-" db:"was_used"`

	SessionIDToken     sqlxx.MapStringInterface `db:"session_id_token" json:"-"`
	SessionAccessToken sqlxx.MapStringInterface `db:"session_access_token" json:"-"`
}

func (r *HandledConsentRequest) HasError() bool {
	return r.Error.IsError()
}

func (r *HandledConsentRequest) prepareSQL() *HandledConsentRequest {
	if r.Session != nil {
		r.SessionAccessToken = r.Session.AccessToken
		r.SessionIDToken = r.Session.IDToken
	}
	return r
}

func (r *HandledConsentRequest) postSQL(cr *ConsentRequest) *HandledConsentRequest {
	r.ConsentRequest = cr
	if r.SessionAccessToken == nil {
		r.SessionAccessToken = make(map[string]interface{})
	}
	if r.SessionIDToken == nil {
		r.SessionIDToken = make(map[string]interface{})
	}
	r.Session = &ConsentRequestSessionData{AccessToken: r.SessionAccessToken, IDToken: r.SessionIDToken}
	return r
}

// The response used to return used consent requests
// same as HandledLoginRequest, just with consent_request exposed as json
type PreviousConsentSession struct {
	// GrantScope sets the scope the user authorized the client to use. Should be a subset of `requested_scope`.
	GrantedScope sqlxx.StringSlicePipeDelimiter `json:"grant_scope" db:"granted_scope"`

	// GrantedAudience sets the audience the user authorized the client to use. Should be a subset of `requested_access_token_audience`.
	GrantedAudience sqlxx.StringSlicePipeDelimiter `json:"grant_access_token_audience" db:"granted_at_audience"`

	// Session allows you to set (optional) session data for access and ID tokens.
	Session *ConsentRequestSessionData `json:"session" db:"-"`

	// Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same
	// client asks the same user for the same, or a subset of, scope.
	Remember bool `json:"remember" db:"remember"`

	// RememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the
	// authorization will be remembered indefinitely.
	RememberFor int `json:"remember_for" db:"remember_for"`

	// HandledAt contains the timestamp the consent request was handled.
	HandledAt sqlxx.NullTime `json:"handled_at" db:"handled_at"`

	ConsentRequest  *ConsentRequest     `json:"consent_request" db:"-"`
	Error           *RequestDeniedError `json:"-" db:"error"`
	Challenge       string              `json:"-" db:"challenge"`
	RequestedAt     time.Time           `json:"-" db:"requested_at"`
	AuthenticatedAt sqlxx.NullTime      `json:"-" db:"authenticated_at"`
	WasUsed         bool                `json:"-" db:"was_used"`

	SessionIDToken     sqlxx.MapStringInterface `db:"session_id_token" json:"-"`
	SessionAccessToken sqlxx.MapStringInterface `db:"session_access_token" json:"-"`
}

// HandledLoginRequest is the request payload used to accept a login request.
//
// swagger:model acceptLoginRequest
type HandledLoginRequest struct {
	// Remember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store
	// a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she
	// will not be asked to log in again.
	Remember bool `json:"remember" db:"remember"`

	// RememberFor sets how long the authentication should be remembered for in seconds. If set to `0`, the
	// authorization will be remembered for the duration of the browser session (using a session cookie).
	RememberFor int `json:"remember_for" db:"remember_for"`

	// ACR sets the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it
	// to express that, for example, a user authenticated using two factor authentication.
	ACR string `json:"acr" db:"acr"`

	// Subject is the user ID of the end-user that authenticated.
	// required: true
	Subject string `json:"subject" db:"subject"`

	// ForceSubjectIdentifier forces the "pairwise" user ID of the end-user that authenticated. The "pairwise" user ID refers to the
	// (Pairwise Identifier Algorithm)[http://openid.net/specs/openid-connect-core-1_0.html#PairwiseAlg] of the OpenID
	// Connect specification. It allows you to set an obfuscated subject ("user") identifier that is unique to the client.
	//
	// Please note that this changes the user ID on endpoint /userinfo and sub claim of the ID Token. It does not change the
	// sub claim in the OAuth 2.0 Introspection.
	//
	// Per default, ORY Hydra handles this value with its own algorithm. In case you want to set this yourself
	// you can use this field. Please note that setting this field has no effect if `pairwise` is not configured in
	// ORY Hydra or the OAuth 2.0 Client does not expect a pairwise identifier (set via `subject_type` key in the client's
	// configuration).
	//
	// Please also be aware that ORY Hydra is unable to properly compute this value during authentication. This implies
	// that you have to compute this value on every authentication process (probably depending on the client ID or some
	// other unique value).
	//
	// If you fail to compute the proper value, then authentication processes which have id_token_hint set might fail.
	ForceSubjectIdentifier string `json:"force_subject_identifier" db:"forced_subject_identifier"`

	// Context is an optional object which can hold arbitrary data. The data will be made available when fetching the
	// consent request under the "context" field. This is useful in scenarios where login and consent endpoints share
	// data.
	Context sqlxx.JSONRawMessage `json:"context" db:"context"`

	LoginRequest    *LoginRequest       `json:"-" db:"-"`
	Error           *RequestDeniedError `json:"-" db:"error"`
	Challenge       string              `json:"-" db:"challenge"`
	RequestedAt     time.Time           `json:"-" db:"requested_at"`
	AuthenticatedAt sqlxx.NullTime      `json:"-" db:"authenticated_at"`
	WasUsed         bool                `json:"-" db:"was_used"`
}

func (r *HandledLoginRequest) HasError() bool {
	return r.Error.IsError()
}

func (r *HandledLoginRequest) postSQL(lr *LoginRequest) *HandledLoginRequest {
	r.LoginRequest = lr
	return r
}

func (r *HandledLoginRequest) prepareSQL() *HandledLoginRequest {
	if string(r.Context) == "" {
		r.Context = sqlxx.JSONRawMessage("{}")
	}
	return r
}

// Contains optional information about the OpenID Connect request.
//
// swagger:model openIDConnectContext
type OpenIDConnectContext struct {
	// ACRValues is the Authentication AuthorizationContext Class Reference requested in the OAuth 2.0 Authorization request.
	// It is a parameter defined by OpenID Connect and expresses which level of authentication (e.g. 2FA) is required.
	//
	// OpenID Connect defines it as follows:
	// > Requested Authentication AuthorizationContext Class Reference values. Space-separated string that specifies the acr values
	// that the Authorization Server is being requested to use for processing this Authentication Request, with the
	// values appearing in order of preference. The Authentication AuthorizationContext Class satisfied by the authentication
	// performed is returned as the acr Claim Value, as specified in Section 2. The acr Claim is requested as a
	// Voluntary Claim by this parameter.
	ACRValues []string `json:"acr_values,omitempty"`

	// UILocales is the End-User'id preferred languages and scripts for the user interface, represented as a
	// space-separated list of BCP47 [RFC5646] language tag values, ordered by preference. For instance, the value
	// "fr-CA fr en" represents a preference for French as spoken in Canada, then French (without a region designation),
	// followed by English (without a region designation). An error SHOULD NOT result if some or all of the requested
	// locales are not supported by the OpenID Provider.
	UILocales []string `json:"ui_locales,omitempty"`

	// Display is a string value that specifies how the Authorization Server displays the authentication and consent user interface pages to the End-User.
	// The defined values are:
	// - page: The Authorization Server SHOULD display the authentication and consent UI consistent with a full User Agent page view. If the display parameter is not specified, this is the default display mode.
	// - popup: The Authorization Server SHOULD display the authentication and consent UI consistent with a popup User Agent window. The popup User Agent window should be of an appropriate size for a login-focused dialog and should not obscure the entire window that it is popping up over.
	// - touch: The Authorization Server SHOULD display the authentication and consent UI consistent with a device that leverages a touch interface.
	// - wap: The Authorization Server SHOULD display the authentication and consent UI consistent with a "feature phone" type display.
	//
	// The Authorization Server MAY also attempt to detect the capabilities of the User Agent and present an appropriate display.
	Display string `json:"display,omitempty"`

	// IDTokenHintClaims are the claims of the ID Token previously issued by the Authorization Server being passed as a hint about the
	// End-User's current or past authenticated session with the Client.
	IDTokenHintClaims map[string]interface{} `json:"id_token_hint_claims,omitempty"`

	// LoginHint hints about the login identifier the End-User might use to log in (if necessary).
	// This hint can be used by an RP if it first asks the End-User for their e-mail address (or other identifier)
	// and then wants to pass that value as a hint to the discovered authorization service. This value MAY also be a
	// phone number in the format specified for the phone_number Claim. The use of this parameter is optional.
	LoginHint string `json:"login_hint,omitempty"`
}

func (n *OpenIDConnectContext) Scan(value interface{}) error {
	v := fmt.Sprintf("%s", value)
	if len(v) == 0 {
		return nil
	}
	return errors.WithStack(json.Unmarshal([]byte(v), n))
}

func (n *OpenIDConnectContext) Value() (driver.Value, error) {
	value, err := json.Marshal(n)
	return value, errors.WithStack(err)
}

// Contains information about an ongoing logout request.
//
// swagger:model logoutRequest
type LogoutRequest struct {
	// Challenge is the identifier ("logout challenge") of the logout authentication request. It is used to
	// identify the session.
	Challenge string `json:"-" db:"challenge"`

	// Subject is the user for whom the logout was request.
	Subject string `json:"subject" db:"subject"`

	// SessionID is the login session ID that was requested to log out.
	SessionID string `json:"sid,omitempty" db:"sid"`

	// RequestURL is the original Logout URL requested.
	RequestURL string `json:"request_url" db:"request_url"`

	// RPInitiated is set to true if the request was initiated by a Relying Party (RP), also known as an OAuth 2.0 Client.
	RPInitiated bool `json:"rp_initiated" db:"rp_initiated"`

	Verifier              string         `json:"-" db:"verifier"`
	PostLogoutRedirectURI string         `json:"-" db:"redir_url"`
	WasUsed               bool           `json:"-" db:"was_used"`
	Accepted              bool           `json:"-" db:"accepted"`
	Rejected              bool           `db:"rejected" json:"-"`
	ClientID              sql.NullString `json:"-" db:"client_id"`
	Client                *client.Client `json:"-" db:"-"`
}

func (r *LogoutRequest) prepareSQL() *LogoutRequest {
	if r.Client != nil {
		r.ClientID = sql.NullString{
			Valid:  true,
			String: r.Client.ClientID,
		}
	}
	return r
}

// Returned when the log out request was used.
//
// swagger:ignore
type LogoutResult struct {
	RedirectTo             string
	FrontChannelLogoutURLs []string
}

// Contains information on an ongoing login request.
//
// swagger:model loginRequest
type LoginRequest struct {
	// Challenge is the identifier ("login challenge") of the login request. It is used to
	// identify the session.
	Challenge string `json:"challenge" db:"challenge"`

	// RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client.
	RequestedScope sqlxx.StringSlicePipeDelimiter `json:"requested_scope" db:"requested_scope"`

	// RequestedScope contains the access token audience as requested by the OAuth 2.0 Client.
	RequestedAudience sqlxx.StringSlicePipeDelimiter `json:"requested_access_token_audience" db:"requested_at_audience"`

	// Skip, if true, implies that the client has requested the same scopes from the same user previously.
	// If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.
	//
	// This feature allows you to update / set session information.
	Skip bool `json:"skip" db:"skip"`

	// Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope
	// requested by the OAuth 2.0 client. If this value is set and `skip` is true, you MUST include this subject type
	// when accepting the login request, or the request will fail.
	Subject string `json:"subject" db:"subject"`

	// OpenIDConnectContext provides context for the (potential) OpenID Connect context. Implementation of these
	// values in your app are optional but can be useful if you want to be fully compliant with the OpenID Connect spec.
	OpenIDConnectContext *OpenIDConnectContext `json:"oidc_context" db:"oidc_context"`

	// Client is the OAuth 2.0 Client that initiated the request.
	Client *client.Client `json:"client"`

	ClientID string `json:"-" db:"client_id"`

	// RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which
	// initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but
	// might come in handy if you want to deal with additional request parameters.
	RequestURL string `json:"request_url" db:"request_url"`

	// SessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag)
	// this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false)
	// this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back-
	// channel logout. It's value can generally be used to associate consecutive login requests by a certain user.
	SessionID sqlxx.NullString `json:"session_id" db:"login_session_id"`

	ForceSubjectIdentifier string `json:"-" db:"-"` // this is here but has no meaning apart from sql_helper working properly.
	Verifier               string `json:"-" db:"verifier"`
	CSRF                   string `json:"-" db:"csrf"`

	AuthenticatedAt sqlxx.NullTime `json:"-" db:"authenticated_at"`
	RequestedAt     time.Time      `json:"-" db:"requested_at"`
	WasHandled      bool           `json:"-" db:"was_handled"`
	Context         string         `json:"-" db:"context"`
}

func (r *LoginRequest) prepareSQL() *LoginRequest {
	if r.Client == nil {
		return r
	}
	r.ClientID = r.Client.ClientID
	return r
}

// Contains information on an ongoing consent request.
//
// swagger:model consentRequest
type ConsentRequest struct {
	// Challenge is the identifier ("authorization challenge") of the consent authorization request. It is used to
	// identify the session.
	Challenge string `json:"challenge" db:"challenge"`

	// RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client.
	RequestedScope sqlxx.StringSlicePipeDelimiter `json:"requested_scope" db:"requested_scope"`

	// RequestedScope contains the access token audience as requested by the OAuth 2.0 Client.
	RequestedAudience sqlxx.StringSlicePipeDelimiter `json:"requested_access_token_audience" db:"requested_at_audience"`

	// Skip, if true, implies that the client has requested the same scopes from the same user previously.
	// If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the
	// consent request using the usual API call.
	Skip bool `json:"skip" db:"skip"`

	// Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope
	// requested by the OAuth 2.0 client.
	Subject string `json:"subject" db:"subject"`

	// OpenIDConnectContext provides context for the (potential) OpenID Connect context. Implementation of these
	// values in your app are optional but can be useful if you want to be fully compliant with the OpenID Connect spec.
	OpenIDConnectContext *OpenIDConnectContext `json:"oidc_context" db:"oidc_context"`

	// Client is the OAuth 2.0 Client that initiated the request.
	Client   *client.Client `json:"client" db:"-"`
	ClientID string         `json:"-" db:"client_id"`

	// RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which
	// initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but
	// might come in handy if you want to deal with additional request parameters.
	RequestURL string `json:"request_url" db:"request_url"`

	// LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate
	// a login and consent request in the login & consent app.
	LoginChallenge sqlxx.NullString `json:"login_challenge" db:"login_challenge"`

	// LoginSessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag)
	// this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false)
	// this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back-
	// channel logout. It's value can generally be used to associate consecutive login requests by a certain user.
	LoginSessionID sqlxx.NullString `json:"login_session_id" db:"login_session_id"`

	// ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it
	// to express that, for example, a user authenticated using two factor authentication.
	ACR string `json:"acr" db:"acr"`

	// Context contains arbitrary information set by the login endpoint or is empty if not set.
	Context sqlxx.JSONRawMessage `json:"context,omitempty" db:"context"`

	// ForceSubjectIdentifier is the value from authentication (if set).
	ForceSubjectIdentifier string         `json:"-" db:"forced_subject_identifier"`
	SubjectIdentifier      string         `json:"-" db:"-"`
	Verifier               string         `json:"-" db:"verifier"`
	CSRF                   string         `json:"-" db:"csrf"`
	AuthenticatedAt        sqlxx.NullTime `json:"-" db:"authenticated_at"`
	RequestedAt            time.Time      `json:"-" db:"requested_at"`
	WasHandled             bool           `json:"-" db:"was_handled"`
}

func (r *ConsentRequest) prepareSQL() *ConsentRequest {
	if r.Client == nil {
		return r
	}
	r.ClientID = r.Client.ClientID
	return r
}

// Used to pass session data to a consent request.
//
// swagger:model consentRequestSession
type ConsentRequestSessionData struct {
	// AccessToken sets session data for the access and refresh token, as well as any future tokens issued by the
	// refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection.
	// If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties
	// can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care!
	AccessToken map[string]interface{} `json:"access_token"`

	// IDToken sets session data for the OpenID Connect ID token. Keep in mind that the session'id payloads are readable
	// by anyone that has access to the ID Challenge. Use with care!
	IDToken map[string]interface{} `json:"id_token"`

	// UserInfo map[string]interface{} `json:"userinfo"`
}

func NewConsentRequestSessionData() *ConsentRequestSessionData {
	return &ConsentRequestSessionData{
		AccessToken: map[string]interface{}{},
		IDToken:     map[string]interface{}{},
	}
}
