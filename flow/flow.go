// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/oauth2/flowctx"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

// FlowState* constants enumerate the states of a flow. The below graph
// describes possible flow state transitions.
//
// graph TD
//
//	LOGIN_INITIALIZED --> LOGIN_UNUSED
//	LOGIN_UNUSED --> LOGIN_USED
//	LOGIN_UNUSED --> LOGIN_ERROR
//	LOGIN_USED --> CONSENT_INITIALIZED
//	CONSENT_INITIALIZED --> CONSENT_UNUSED
//	CONSENT_UNUSED --> CONSENT_UNUSED
//	CONSENT_UNUSED --> CONSENT_USED
//	CONSENT_UNUSED --> CONSENT_ERROR
const (
	// FlowStateLoginInitialized applies before the login app either
	// accepts or rejects the login request.
	FlowStateLoginInitialized = int16(1)

	// FlowStateLoginUnused indicates that the login has been authenticated, but
	// the User Agent hasn't picked up the result yet.
	FlowStateLoginUnused = int16(2)

	// FlowStateLoginUsed indicates that the User Agent is requesting consent and
	// Hydra has invalidated the login request. This is a short-lived state
	// because the transition to FlowStateConsentInitialized should happen while
	// handling the request that triggered the transition to FlowStateLoginUsed.
	FlowStateLoginUsed = int16(3)

	// FlowStateConsentInitialized applies while Hydra waits for a consent request
	// to be accepted or rejected.
	FlowStateConsentInitialized = int16(4)

	FlowStateConsentUnused = int16(5)
	FlowStateConsentUsed   = int16(6)

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
	FlowStateLoginError   = int16(128)
	FlowStateConsentError = int16(129)
)

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
	// ID is the identifier ("login challenge") of the login request. It is used to
	// identify the session.
	//
	// required: true
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
	LoginSkip bool `db:"login_skip" json:"ls,omitempty"`

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
	IdentityProviderSessionID sqlxx.NullString `db:"identity_provider_session_id" json:"is,omitempty"`

	LoginVerifier string `db:"login_verifier" json:"lv,omitempty"`
	LoginCSRF     string `db:"login_csrf" json:"lc,omitempty"`

	LoginInitializedAt sqlxx.NullTime `db:"login_initialized_at" json:"li,omitempty"`
	RequestedAt        time.Time      `db:"requested_at" json:"ia,omitempty"`

	State int16 `db:"state" json:"q,omitempty"`

	// LoginRemember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store
	// a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she
	// will not be asked to log in again.
	LoginRemember bool `db:"login_remember" json:"lr,omitempty"`

	// LoginRememberFor sets how long the authentication should be remembered for in seconds. If set to `0`, the
	// authorization will be remembered for the duration of the browser session (using a session cookie).
	LoginRememberFor int `db:"login_remember_for" json:"lf,omitempty"`

	// LoginExtendSessionLifespan, if set to true, session cookie expiry time will be updated when session is
	// refreshed (login skip=true).
	LoginExtendSessionLifespan bool `db:"login_extend_session_lifespan" json:"ll,omitempty"`

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
	ForceSubjectIdentifier string `db:"forced_subject_identifier" json:"fs,omitempty"`

	// Context is an optional object which can hold arbitrary data. The data will be made available when fetching the
	// consent request under the "context" field. This is useful in scenarios where login and consent endpoints share
	// data.
	Context sqlxx.JSONRawMessage `db:"context" json:"ct"`

	// LoginWasUsed set to true means that the login request was already handled.
	// This can happen on form double-submit or other errors. If this is set we
	// recommend redirecting the user to `request_url` to re-initiate the flow.
	LoginWasUsed bool `db:"login_was_used" json:"lu,omitempty"`

	LoginError           *RequestDeniedError `db:"login_error" json:"le,omitempty"`
	LoginAuthenticatedAt sqlxx.NullTime      `db:"login_authenticated_at" json:"la,omitempty"`

	// ConsentChallengeID is the identifier ("authorization challenge") of the consent authorization request. It is used to
	// identify the session.
	//
	// required: true
	ConsentChallengeID sqlxx.NullString `db:"consent_challenge_id" json:"cc,omitempty"`

	// ConsentSkip, if true, implies that the client has requested the same scopes from the same user previously.
	// If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the
	// consent request using the usual API call.
	ConsentSkip     bool             `db:"consent_skip" json:"cs,omitempty"`
	ConsentVerifier sqlxx.NullString `db:"consent_verifier" json:"cv,omitempty"`
	ConsentCSRF     sqlxx.NullString `db:"consent_csrf" json:"cr,omitempty"`

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

	// ConsentWasHandled set to true means that the request was already handled.
	// This can happen on form double-submit or other errors. If this is set we
	// recommend redirecting the user to `request_url` to re-initiate the flow.
	ConsentWasHandled  bool                     `db:"consent_was_used" json:"cw,omitempty"`
	ConsentError       *RequestDeniedError      `db:"consent_error" json:"cx"`
	SessionIDToken     sqlxx.MapStringInterface `db:"session_id_token" faker:"-" json:"st"`
	SessionAccessToken sqlxx.MapStringInterface `db:"session_access_token" faker:"-" json:"sa"`
}

func NewFlow(r *LoginRequest) *Flow {
	return &Flow{
		ID:                     r.ID,
		RequestedScope:         r.RequestedScope,
		RequestedAudience:      r.RequestedAudience,
		LoginSkip:              r.Skip,
		Subject:                r.Subject,
		OpenIDConnectContext:   r.OpenIDConnectContext,
		Client:                 r.Client,
		ClientID:               r.ClientID,
		RequestURL:             r.RequestURL,
		SessionID:              r.SessionID,
		LoginWasUsed:           r.WasHandled,
		ForceSubjectIdentifier: r.ForceSubjectIdentifier,
		LoginVerifier:          r.Verifier,
		LoginCSRF:              r.CSRF,
		LoginAuthenticatedAt:   r.AuthenticatedAt,
		RequestedAt:            r.RequestedAt,
		State:                  FlowStateLoginInitialized,
	}
}

func (f *Flow) HandleLoginRequest(h *HandledLoginRequest) error {
	if f.LoginWasUsed {
		return errors.WithStack(x.ErrConflict.WithHint("The login request was already used and can no longer be changed."))
	}

	if f.State != FlowStateLoginInitialized && f.State != FlowStateLoginUnused && f.State != FlowStateLoginError {
		return errors.Errorf("invalid flow state: expected %d/%d/%d, got %d", FlowStateLoginInitialized, FlowStateLoginUnused, FlowStateLoginError, f.State)
	}

	if f.ID != h.ID {
		return errors.Errorf("flow ID %s does not match HandledLoginRequest ID %s", f.ID, h.ID)
	}

	if f.Subject != "" && h.Subject != "" && f.Subject != h.Subject {
		return errors.Errorf("flow Subject %s does not match the HandledLoginRequest Subject %s", f.Subject, h.Subject)
	}

	if f.ForceSubjectIdentifier != "" && h.ForceSubjectIdentifier != "" && f.ForceSubjectIdentifier != h.ForceSubjectIdentifier {
		return errors.Errorf("flow ForceSubjectIdentifier %s does not match the HandledLoginRequest ForceSubjectIdentifier %s", f.ForceSubjectIdentifier, h.ForceSubjectIdentifier)
	}

	if h.Error != nil {
		f.State = FlowStateLoginError
	} else {
		f.State = FlowStateLoginUnused
	}

	if f.Context != nil {
		f.Context = h.Context
	}

	f.ID = h.ID
	f.Subject = h.Subject
	f.ForceSubjectIdentifier = h.ForceSubjectIdentifier
	f.LoginError = h.Error

	f.IdentityProviderSessionID = sqlxx.NullString(h.IdentityProviderSessionID)
	f.LoginRemember = h.Remember
	f.LoginRememberFor = h.RememberFor
	f.LoginExtendSessionLifespan = h.ExtendSessionLifespan
	f.ACR = h.ACR
	f.AMR = h.AMR
	f.LoginWasUsed = h.WasHandled
	f.LoginAuthenticatedAt = h.AuthenticatedAt
	return nil
}

func (f *Flow) GetHandledLoginRequest() HandledLoginRequest {
	return HandledLoginRequest{
		ID:                        f.ID,
		Remember:                  f.LoginRemember,
		RememberFor:               f.LoginRememberFor,
		ExtendSessionLifespan:     f.LoginExtendSessionLifespan,
		ACR:                       f.ACR,
		AMR:                       f.AMR,
		Subject:                   f.Subject,
		IdentityProviderSessionID: f.IdentityProviderSessionID.String(),
		ForceSubjectIdentifier:    f.ForceSubjectIdentifier,
		Context:                   f.Context,
		WasHandled:                f.LoginWasUsed,
		Error:                     f.LoginError,
		LoginRequest:              f.GetLoginRequest(),
		RequestedAt:               f.RequestedAt,
		AuthenticatedAt:           f.LoginAuthenticatedAt,
	}
}

func (f *Flow) GetLoginRequest() *LoginRequest {
	return &LoginRequest{
		ID:                     f.ID,
		RequestedScope:         f.RequestedScope,
		RequestedAudience:      f.RequestedAudience,
		Skip:                   f.LoginSkip,
		Subject:                f.Subject,
		OpenIDConnectContext:   f.OpenIDConnectContext,
		Client:                 f.Client,
		ClientID:               f.ClientID,
		RequestURL:             f.RequestURL,
		SessionID:              f.SessionID,
		WasHandled:             f.LoginWasUsed,
		ForceSubjectIdentifier: f.ForceSubjectIdentifier,
		Verifier:               f.LoginVerifier,
		CSRF:                   f.LoginCSRF,
		AuthenticatedAt:        f.LoginAuthenticatedAt,
		RequestedAt:            f.RequestedAt,
	}
}

// InvalidateLoginRequest shifts the flow state to FlowStateLoginUsed. This
// transition is executed upon login completion.
func (f *Flow) InvalidateLoginRequest() error {
	if f.State != FlowStateLoginUnused && f.State != FlowStateLoginError {
		return errors.Errorf("invalid flow state: expected %d or %d, got %d", FlowStateLoginUnused, FlowStateLoginError, f.State)
	}
	if f.LoginWasUsed {
		return errors.New("login verifier has already been used")
	}
	f.LoginWasUsed = true
	f.State = FlowStateLoginUsed
	return nil
}

func (f *Flow) HandleConsentRequest(r *AcceptOAuth2ConsentRequest) error {
	if time.Time(r.HandledAt).IsZero() {
		return errors.New("refusing to handle a consent request with null HandledAt")
	}

	if f.ConsentWasHandled {
		return x.ErrConflict.WithHint("The consent request was already used and can no longer be changed.")
	}

	if f.State != FlowStateConsentInitialized && f.State != FlowStateConsentUnused && f.State != FlowStateConsentError {
		return errors.Errorf("invalid flow state: expected %d/%d/%d, got %d", FlowStateConsentInitialized, FlowStateConsentUnused, FlowStateConsentError, f.State)
	}

	if f.ConsentChallengeID.String() != r.ID {
		return errors.Errorf("flow.ConsentChallengeID %s doesn't match AcceptOAuth2ConsentRequest.ID %s", f.ConsentChallengeID.String(), r.ID)
	}

	if r.Error != nil {
		f.State = FlowStateConsentError
	} else if r.WasHandled {
		f.State = FlowStateConsentUsed
	} else {
		f.State = FlowStateConsentUnused
	}

	f.GrantedScope = r.GrantedScope
	f.GrantedAudience = r.GrantedAudience
	f.ConsentRemember = r.Remember
	f.ConsentRememberFor = &r.RememberFor
	f.ConsentHandledAt = r.HandledAt
	f.ConsentWasHandled = r.WasHandled
	f.ConsentError = r.Error
	if r.Context != nil {
		f.Context = r.Context
	}

	if r.Session != nil {
		f.SessionIDToken = r.Session.IDToken
		f.SessionAccessToken = r.Session.AccessToken
	}
	return nil
}

func (f *Flow) InvalidateConsentRequest() error {
	if f.ConsentWasHandled {
		return errors.New("consent verifier has already been used")
	}
	if f.State != FlowStateConsentUnused && f.State != FlowStateConsentError {
		return errors.Errorf("unexpected flow state: expected %d or %d, got %d", FlowStateConsentUnused, FlowStateConsentError, f.State)
	}

	f.ConsentWasHandled = true
	f.State = FlowStateConsentUsed
	return nil
}

func (f *Flow) GetConsentRequest() *OAuth2ConsentRequest {
	cs := OAuth2ConsentRequest{
		ID:                     f.ConsentChallengeID.String(),
		RequestedScope:         f.RequestedScope,
		RequestedAudience:      f.RequestedAudience,
		Skip:                   f.ConsentSkip,
		Subject:                f.Subject,
		OpenIDConnectContext:   f.OpenIDConnectContext,
		Client:                 f.Client,
		ClientID:               f.ClientID,
		RequestURL:             f.RequestURL,
		LoginChallenge:         sqlxx.NullString(f.ID),
		LoginSessionID:         f.SessionID,
		ACR:                    f.ACR,
		AMR:                    f.AMR,
		Context:                f.Context,
		WasHandled:             f.ConsentWasHandled,
		ForceSubjectIdentifier: f.ForceSubjectIdentifier,
		Verifier:               f.ConsentVerifier.String(),
		CSRF:                   f.ConsentCSRF.String(),
		AuthenticatedAt:        f.LoginAuthenticatedAt,
		RequestedAt:            f.RequestedAt,
	}
	if cs.AMR == nil {
		cs.AMR = []string{}
	}
	return &cs
}

func (f *Flow) GetHandledConsentRequest() *AcceptOAuth2ConsentRequest {
	crf := 0
	if f.ConsentRememberFor != nil {
		crf = *f.ConsentRememberFor
	}
	return &AcceptOAuth2ConsentRequest{
		ID:                 f.ConsentChallengeID.String(),
		GrantedScope:       f.GrantedScope,
		GrantedAudience:    f.GrantedAudience,
		Session:            &AcceptOAuth2ConsentRequestSession{AccessToken: f.SessionAccessToken, IDToken: f.SessionIDToken},
		Remember:           f.ConsentRemember,
		RememberFor:        crf,
		HandledAt:          f.ConsentHandledAt,
		WasHandled:         f.ConsentWasHandled,
		Context:            f.Context,
		ConsentRequest:     f.GetConsentRequest(),
		Error:              f.ConsentError,
		RequestedAt:        f.RequestedAt,
		AuthenticatedAt:    f.LoginAuthenticatedAt,
		SessionIDToken:     f.SessionIDToken,
		SessionAccessToken: f.SessionAccessToken,
	}
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

// ToLoginChallenge converts the flow into a login challenge.
func (f Flow) ToLoginChallenge(ctx context.Context, cipherProvider CipherProvider) (string, error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return flowctx.Encode(ctx, cipherProvider.FlowCipher(), f, flowctx.AsLoginChallenge)
}

// ToLoginVerifier converts the flow into a login verifier.
func (f Flow) ToLoginVerifier(ctx context.Context, cipherProvider CipherProvider) (string, error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return flowctx.Encode(ctx, cipherProvider.FlowCipher(), f, flowctx.AsLoginVerifier)
}

// ToConsentChallenge converts the flow into a consent challenge.
func (f Flow) ToConsentChallenge(ctx context.Context, cipherProvider CipherProvider) (string, error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return flowctx.Encode(ctx, cipherProvider.FlowCipher(), f, flowctx.AsConsentChallenge)
}

// ToConsentVerifier converts the flow into a consent verifier.
func (f Flow) ToConsentVerifier(ctx context.Context, cipherProvider CipherProvider) (string, error) {
	if f.Client != nil {
		f.ClientID = f.Client.GetID()
	}
	return flowctx.Encode(ctx, cipherProvider.FlowCipher(), f, flowctx.AsConsentVerifier)
}
