// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/pop/v6"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

// FlowState* constants enumerate the states of a flow. The below graph
// describes possible flow state transitions.
//
// stateDiagram-v2
//  [*] --> DEVICE_UNUSED: GET /oauth2/device/verify
//  DEVICE_UNUSED --> DEVICE_USED: submit user code
//  DEVICE_USED --> LOGIN_UNUSED: to verifier
//  [*] --> LOGIN_UNUSED: GET /oauth2/auth
//  LOGIN_UNUSED --> LOGIN_UNUSED: accept login
//  LOGIN_UNUSED --> LOGIN_USED: submit login verifier
//  LOGIN_UNUSED --> LOGIN_ERROR: reject login
//  LOGIN_ERROR --> [*]
//  LOGIN_USED --> CONSENT_UNUSED
//  CONSENT_UNUSED --> CONSENT_UNUSED: accept consent
//  CONSENT_UNUSED --> CONSENT_USED: submit consent verifier
//  CONSENT_UNUSED --> CONSENT_ERROR: reject consent
//  CONSENT_ERROR --> [*]
//  CONSENT_USED --> [*]

type State int16

const (
	// FlowStateLoginInitialized is not used anymore, but is kept for
	// backwards compatibility. New flows start at FlowStateLoginUnused.
	FlowStateLoginInitialized = State(1)

	// FlowStateLoginUnused indicates that the login has been authenticated, but
	// the User Agent hasn't picked up the result yet.
	FlowStateLoginUnused = State(2)

	// FlowStateLoginUsed indicates that the User Agent is requesting consent and
	// Hydra has invalidated the login request. This is a short-lived state
	// because the transition to FlowStateConsentInitialized should happen while
	// handling the request that triggered the transition to FlowStateLoginUsed.
	FlowStateLoginUsed = State(3)

	// FlowStateConsentInitialized is not used anymore, but is kept for
	// backwards compatibility. New flows start at FlowStateConsentUnused.
	FlowStateConsentInitialized = State(4)

	FlowStateConsentUnused = State(5)
	FlowStateConsentUsed   = State(6)

	// DeviceFlowStateInitialized is not used anymore, but is kept for
	// backwards compatibility. New flows start at DeviceFlowStateUnused.
	DeviceFlowStateInitialized = State(7)

	// DeviceFlowStateUnused indicates that the login has been authenticated, but
	// the User Agent hasn't picked up the result yet.
	DeviceFlowStateUnused = State(8)

	// DeviceFlowStateUsed indicates that the User Agent is requesting consent and
	// Hydra has invalidated the login request. This is a short-lived state
	// because the transition to DeviceFlowStateConsentInitialized should happen while
	// handling the request that triggered the transition to DeviceFlowStateUsed.
	DeviceFlowStateUsed = State(9)

	// TODO: Refactor error handling to persist error codes instead of JSON
	// strings. Currently we persist errors as JSON strings in the LoginError
	// and ConsentError fields. This shouldn't be necessary because the different
	// errors are enumerable; most of them have error codes defined in Fosite. It
	// is possible to define a mapping between error codes and the metadata that
	// is currently persisted with each erred Flow. This mapping would be used in
	// GetConsentRequest, HandleConsentRequest, GetHandledLoginRequest, etc. An
	// ErrorContext field can be introduced later if it becomes necessary.
	// If the above is implemented, merge the LoginError and ConsentError fields
	// and use the following FlowStates when converting to/from
	// [Handled]{Login|Consent}Request:
	FlowStateLoginError   = State(128)
	FlowStateConsentError = State(129)
)

func (s State) ConsentWasUsed() bool { return s == FlowStateConsentUsed || s == FlowStateConsentError }
func (s State) LoginWasUsed() bool   { return s == FlowStateLoginUsed || s == FlowStateLoginError }

func (s State) IsAny(expected ...State) error {
	for _, e := range expected {
		if s == e {
			return nil
		}
	}
	return errors.Errorf("invalid flow state: expected one of %v, got %d", expected, s)
}

// Flow is an abstraction used in the persistence layer to unify LoginRequest,
// HandledLoginRequest, ConsentRequest, and AcceptOAuth2ConsentRequest.
//
// TODO: Deprecate the structs that are made obsolete by the Flow concept.
// Context: Before Flow was introduced, the API and the database used the same
// structs, LoginRequest and HandledLoginRequest. These two tables and structs
// were merged into a new concept, Flow, in order to optimize the persistence
// layer. We currently limit the use of Flow to the persistence layer and keep
// using the original structs in the API in order to minimize the impact of the
// database refactoring on the API.
type Flow struct {
	// ID is the identifier of the login request.
	//
	// The struct field is named ID for compatibility with gobuffalo/pop, and is
	// the primary key in the database.
	//
	// The database column should be named `login_challenge_id`, but is not for
	// historical reasons.
	//
	// This is not the same as the login session ID.
	ID  string    `db:"login_challenge" json:"i"`
	NID uuid.UUID `db:"nid" json:"n"`

	// RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client.
	//
	// required: true
	RequestedScope sqlxx.StringSliceJSONFormat `db:"requested_scope" json:"rs,omitempty"`

	// RequestedAudience contains the access token audience as requested by the OAuth 2.0 Client.
	//
	// required: true
	RequestedAudience sqlxx.StringSliceJSONFormat `db:"requested_at_audience" json:"ra,omitempty"`

	// LoginSkip, if true, implies that the client has requested the same scopes from the same user previously.
	// If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.
	//
	// This feature allows you to update / set session information.
	//
	// required: true
	LoginSkip bool `db:"-" json:"ls,omitempty"`

	// Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope
	// requested by the OAuth 2.0 client. If this value is set and `skip` is true, you MUST include this subject type
	// when accepting the login request, or the request will fail.
	//
	// required: true
	Subject string `db:"subject" json:"s,omitempty"`

	// OpenIDConnectContext provides context for the (potential) OpenID Connect context. Implementation of these
	// values in your app are optional but can be useful if you want to be fully compliant with the OpenID Connect spec.
	OpenIDConnectContext *OAuth2ConsentRequestOpenIDConnectContext `db:"oidc_context" json:"oc"`

	// Client is the OAuth 2.0 Client that initiated the request.
	//
	// required: true
	Client   *client.Client `db:"-" json:"c,omitempty"`
	ClientID string         `db:"client_id" json:"ci,omitempty"`

	// RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which
	// initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but
	// might come in handy if you want to deal with additional request parameters.
	//
	// required: true
	RequestURL string `db:"request_url" json:"r,omitempty"`

	// SessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag)
	// this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false)
	// this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back-
	// channel logout. Its value can generally be used to associate consecutive login requests by a certain user.
	SessionID sqlxx.NullString `db:"login_session_id" json:"si,omitempty"`

	// IdentityProviderSessionID is the session ID of the end-user that authenticated.
	// If specified, we will use this value to propagate the logout.
	IdentityProviderSessionID sqlxx.NullString `db:"-" json:"is,omitempty"`

	LoginCSRF string `db:"-" json:"lc,omitempty"`

	RequestedAt time.Time `db:"requested_at" json:"ia,omitempty"`

	State State `db:"-" json:"q,omitempty"`

	// LoginRemember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store
	// a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she
	// will not be asked to log in again.
	LoginRemember bool `db:"-" json:"lr,omitempty"`

	// LoginRememberFor sets how long the authentication should be remembered for in seconds. If set to `0`, the
	// authorization will be remembered for the duration of the browser session (using a session cookie).
	LoginRememberFor int `db:"-" json:"lf,omitempty"`

	// LoginExtendSessionLifespan, if set to true, session cookie expiry time will be updated when session is
	// refreshed (login skip=true).
	LoginExtendSessionLifespan bool `db:"-" json:"ll,omitempty"`

	// ACR sets the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it
	// to express that, for example, a user authenticated using two factor authentication.
	ACR string `db:"acr" json:"a,omitempty"`

	// AMR sets the Authentication Methods References value for this
	// authentication session. You can use it to specify the method a user used to
	// authenticate. For example, if the acr indicates a user used two factor
	// authentication, the amr can express they used a software-secured key.
	AMR sqlxx.StringSliceJSONFormat `db:"amr" json:"am,omitempty"`

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
	ForceSubjectIdentifier string `db:"-" json:"fs,omitempty"`

	// Context is an optional object which can hold arbitrary data. The data will be made available when fetching the
	// consent request under the "context" field. This is useful in scenarios where login and consent endpoints share
	// data.
	Context sqlxx.JSONRawMessage `db:"context" json:"ct"`

	LoginError           *RequestDeniedError `db:"-" json:"le,omitempty"`
	LoginAuthenticatedAt sqlxx.NullTime      `db:"-" json:"la,omitempty"`

	// DeviceChallengeID is the device request's challenge ID
	DeviceChallengeID sqlxx.NullString `db:"device_challenge_id" json:"di,omitempty"`
	// DeviceCodeRequestID is the device request's ID
	DeviceCodeRequestID sqlxx.NullString `db:"device_code_request_id" json:"dr,omitempty"`
	// DeviceCSRF is the device request's CSRF
	DeviceCSRF sqlxx.NullString `db:"-" json:"dc,omitempty"`
	// DeviceHandledAt contains the timestamp the device user_code verification request was handled
	DeviceHandledAt sqlxx.NullTime `db:"-" json:"dh,omitempty"`

	// ConsentRequestID is the identifier of the consent request.
	// The database column should be named `consent_request_id`, but is not for historical reasons.
	ConsentRequestID sqlxx.NullString `db:"consent_challenge_id" json:"cc,omitempty"`
	// ConsentSkip, if true, implies that the client has requested the same scopes from the same user previously.
	// If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the
	// consent request using the usual API call.
	ConsentSkip bool             `db:"consent_skip" json:"cs,omitempty"`
	ConsentCSRF sqlxx.NullString `db:"-" json:"cr,omitempty"`

	// GrantedScope sets the scope the user authorized the client to use. Should be a subset of `requested_scope`.
	GrantedScope sqlxx.StringSliceJSONFormat `db:"granted_scope" json:"gs,omitempty"`

	// GrantedAudience sets the audience the user authorized the client to use. Should be a subset of `requested_access_token_audience`.
	GrantedAudience sqlxx.StringSliceJSONFormat `db:"granted_at_audience" json:"ga,omitempty"`

	// ConsentRemember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same
	// client asks the same user for the same, or a subset of, scope.
	ConsentRemember bool `db:"consent_remember" json:"ce,omitempty"`

	// ConsentRememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the
	// authorization will be remembered indefinitely.
	ConsentRememberFor *int `db:"consent_remember_for" json:"cf"`

	// ConsentHandledAt contains the timestamp the consent request was handled.
	ConsentHandledAt sqlxx.NullTime `db:"consent_handled_at" json:"ch,omitempty"`

	ConsentError       *RequestDeniedError      `db:"-" json:"cx"`
	SessionIDToken     sqlxx.MapStringInterface `db:"session_id_token" faker:"-" json:"st"`
	SessionAccessToken sqlxx.MapStringInterface `db:"session_access_token" faker:"-" json:"sa"`
}

// HandleDeviceUserAuthRequest updates the flows fields from a handled request.
func (f *Flow) HandleDeviceUserAuthRequest(h *HandledDeviceUserAuthRequest) error {
	if err := f.State.IsAny(DeviceFlowStateInitialized, DeviceFlowStateUnused); err != nil {
		return err
	}

	f.State = DeviceFlowStateUnused

	f.Client = h.Client
	f.ClientID = h.Client.GetID()
	f.DeviceCodeRequestID = sqlxx.NullString(h.DeviceCodeRequestID)
	f.DeviceHandledAt = sqlxx.NullTime(time.Now().UTC())
	f.RequestedScope = h.RequestedScope
	f.RequestedAudience = h.RequestedAudience

	return nil
}

// InvalidateDeviceRequest shifts the flow state to DeviceFlowStateUsed. This
// transition is executed upon device completion.
func (f *Flow) InvalidateDeviceRequest() error {
	if err := f.State.IsAny(DeviceFlowStateUnused); err != nil {
		return err
	}
	f.State = DeviceFlowStateUsed
	return nil
}

func (f *Flow) HandleLoginRequest(h *HandledLoginRequest) error {
	if err := f.State.IsAny(FlowStateLoginInitialized, FlowStateLoginUnused, FlowStateLoginError); err != nil {
		return err
	}

	if f.Subject != "" && h.Subject != "" && f.Subject != h.Subject {
		return errors.Errorf("flow Subject %s does not match the HandledLoginRequest Subject %s", f.Subject, h.Subject)
	}

	if f.ForceSubjectIdentifier != "" && h.ForceSubjectIdentifier != "" && f.ForceSubjectIdentifier != h.ForceSubjectIdentifier {
		return errors.Errorf("flow ForceSubjectIdentifier %s does not match the HandledLoginRequest ForceSubjectIdentifier %s", f.ForceSubjectIdentifier, h.ForceSubjectIdentifier)
	}

	f.State = FlowStateLoginUnused

	if f.Context != nil {
		f.Context = h.Context
	}

	f.Subject = h.Subject
	f.ForceSubjectIdentifier = h.ForceSubjectIdentifier

	f.IdentityProviderSessionID = sqlxx.NullString(h.IdentityProviderSessionID)
	f.LoginRemember = h.Remember
	f.LoginRememberFor = h.RememberFor
	f.LoginExtendSessionLifespan = h.ExtendSessionLifespan
	f.ACR = h.ACR
	f.AMR = h.AMR
	return nil
}

func (f *Flow) HandleLoginError(er *RequestDeniedError) error {
	if err := f.State.IsAny(FlowStateLoginInitialized, FlowStateLoginUnused, FlowStateLoginError); err != nil {
		return err
	}

	f.State = FlowStateLoginError

	f.LoginError = er

	// force-reset values
	f.Subject = ""
	f.ForceSubjectIdentifier = ""
	f.LoginAuthenticatedAt = sqlxx.NullTime{}
	f.IdentityProviderSessionID = ""
	f.LoginRemember = false
	f.LoginRememberFor = 0
	f.LoginExtendSessionLifespan = false
	f.ACR = ""
	f.AMR = nil

	return nil
}

func (f *Flow) GetLoginRequest() *LoginRequest {
	return &LoginRequest{
		ID:                   f.ID,
		RequestedScope:       f.RequestedScope,
		RequestedAudience:    f.RequestedAudience,
		Skip:                 f.LoginSkip,
		Subject:              f.Subject,
		OpenIDConnectContext: f.OpenIDConnectContext,
		Client:               f.Client,
		RequestURL:           f.RequestURL,
		SessionID:            f.SessionID,
	}
}

// InvalidateLoginRequest shifts the flow state to FlowStateLoginUsed. This
// transition is executed upon login completion.
func (f *Flow) InvalidateLoginRequest() error {
	if err := f.State.IsAny(FlowStateLoginUnused, FlowStateLoginError); err != nil {
		return err
	}

	if f.State == FlowStateLoginUnused {
		f.State = FlowStateLoginUsed
	} else {
		// FlowStateLoginError is already a terminal state, so we don't need to do anything here.
	}
	return nil
}

func (f *Flow) HandleConsentRequest(r *AcceptOAuth2ConsentRequest) error {
	if err := f.State.IsAny(FlowStateConsentInitialized, FlowStateConsentUnused, FlowStateConsentError); err != nil {
		return err
	}

	f.State = FlowStateConsentUnused

	f.GrantedScope = r.GrantedScope
	f.GrantedAudience = r.GrantedAudience
	f.ConsentRemember = r.Remember
	f.ConsentRememberFor = &r.RememberFor
	f.ConsentHandledAt = sqlxx.NullTime(time.Now().UTC())
	f.ConsentError = nil
	if r.Context != nil {
		f.Context = r.Context
	}

	if r.Session != nil {
		f.SessionIDToken = r.Session.IDToken
		f.SessionAccessToken = r.Session.AccessToken
	}
	return nil
}

func (f *Flow) HandleConsentError(er *RequestDeniedError) error {
	if err := f.State.IsAny(FlowStateConsentInitialized, FlowStateConsentUnused, FlowStateConsentError); err != nil {
		return err
	}

	f.State = FlowStateConsentError

	f.ConsentError = er
	f.ConsentHandledAt = sqlxx.NullTime(time.Now().UTC())

	// force-reset values
	f.GrantedScope = nil
	f.GrantedAudience = nil
	f.ConsentRemember = false
	f.ConsentRememberFor = nil

	return nil
}

func (f *Flow) InvalidateConsentRequest() error {
	if err := f.State.IsAny(FlowStateConsentUnused, FlowStateConsentError); err != nil {
		return err
	}

	if f.State == FlowStateConsentUnused {
		f.State = FlowStateConsentUsed
	} else {
		// FlowStateConsentError is already a terminal state, so we don't need to do anything here.
	}
	return nil
}

func (f *Flow) GetConsentRequest(challenge string) *OAuth2ConsentRequest {
	cs := OAuth2ConsentRequest{
		Challenge:            challenge,
		ConsentRequestID:     f.ConsentRequestID.String(),
		RequestedScope:       f.RequestedScope,
		RequestedAudience:    f.RequestedAudience,
		Skip:                 f.ConsentSkip,
		Subject:              f.Subject,
		OpenIDConnectContext: f.OpenIDConnectContext,
		Client:               f.Client,
		RequestURL:           f.RequestURL,
		LoginChallenge:       sqlxx.NullString(f.ID),
		LoginSessionID:       f.SessionID,
		ACR:                  f.ACR,
		AMR:                  f.AMR,
		Context:              f.Context,
	}
	if cs.AMR == nil {
		cs.AMR = []string{}
	}
	return &cs
}

func (Flow) TableName() string {
	return "hydra_oauth2_flow"
}

func (f *Flow) BeforeSave(_ *pop.Connection) error {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	if f.State == FlowStateLoginUnused && string(f.Context) == "" {
		f.Context = sqlxx.JSONRawMessage("{}")
	}
	return nil
}

func (f *Flow) AfterFind(c *pop.Connection) error {
	// TODO Populate the client field in FindInDB and FindByConsentChallengeID in
	// order to avoid accessing the database twice.
	f.AfterSave(c)
	f.Client = &client.Client{}
	return sqlcon.HandleError(c.Where("id = ? AND nid = ?", f.ClientID, f.NID).First(f.Client))
}

func (f *Flow) AfterSave(_ *pop.Connection) {
	if f.SessionAccessToken == nil {
		f.SessionAccessToken = make(map[string]interface{})
	}
	if f.SessionIDToken == nil {
		f.SessionIDToken = make(map[string]interface{})
	}
}

type CipherProvider interface {
	FlowCipher() *aead.XChaCha20Poly1305
}

// ToDeviceChallenge converts the flow into a device challenge.
func (f *Flow) ToDeviceChallenge(ctx context.Context, cipherProvider CipherProvider) (string, error) {
	return Encode(ctx, cipherProvider.FlowCipher(), f, AsDeviceChallenge)
}

// ToDeviceVerifier converts the flow into a device verifier.
func (f *Flow) ToDeviceVerifier(ctx context.Context, cipherProvider CipherProvider) (string, error) {
	return Encode(ctx, cipherProvider.FlowCipher(), f, AsDeviceVerifier)
}

// ToLoginChallenge converts the flow into a login challenge.
func (f Flow) ToLoginChallenge(ctx context.Context, cipherProvider CipherProvider) (challenge string, err error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return Encode(ctx, cipherProvider.FlowCipher(), f, AsLoginChallenge)
}

// ToLoginVerifier converts the flow into a login verifier.
func (f Flow) ToLoginVerifier(ctx context.Context, cipherProvider CipherProvider) (verifier string, err error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return Encode(ctx, cipherProvider.FlowCipher(), f, AsLoginVerifier)
}

// ToConsentChallenge converts the flow into a consent challenge.
func (f Flow) ToConsentChallenge(ctx context.Context, cipherProvider CipherProvider) (challenge string, err error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return Encode(ctx, cipherProvider.FlowCipher(), f, AsConsentChallenge)
}

// ToConsentVerifier converts the flow into a consent verifier.
func (f Flow) ToConsentVerifier(ctx context.Context, cipherProvider CipherProvider) (verifier string, err error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return Encode(ctx, cipherProvider.FlowCipher(), f, AsConsentVerifier)
}

func (f Flow) ToListConsentSessionResponse() *OAuth2ConsentSession {
	s := &OAuth2ConsentSession{
		ConsentRequestID: f.ConsentRequestID.String(),
		GrantedScope:     f.GrantedScope,
		GrantedAudience:  f.GrantedAudience,
		RememberFor:      pointerx.Deref(f.ConsentRememberFor),
		Session:          &AcceptOAuth2ConsentRequestSession{AccessToken: f.SessionAccessToken, IDToken: f.SessionIDToken},
		Remember:         f.ConsentRemember,
		HandledAt:        f.ConsentHandledAt,
		Context:          f.Context,
		ConsentRequest:   f.GetConsentRequest( /* No longer available and no longer needed: challenge =  */ ""),
	}
	s.ConsentRequest.Client.Secret = "" // do not leak client secret in response
	return s
}
