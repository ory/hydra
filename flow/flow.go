package flow

import (
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

// FlowState* constants enumerate the states of a flow. The below graph
// describes possible flow state transitions.
//
//graph TD
//    LOGIN_INITIALIZED --> LOGIN_UNUSED
//    LOGIN_UNUSED --> CONSENT_INITIALIZED
//    LOGIN_UNUSED --> ERROR
//    CONSENT_INITIALIZED --> CONSENT_UNUSED
//    CONSENT_UNUSED --> CONSENT_USED
//    CONSENT_UNUSED --> ERROR
const (
	FlowStateError              = int16(0)
	FlowStateLoginInitialized   = int16(1)
	FlowStateLoginUnused        = int16(2)
	FlowStateConsentInitialized = int16(3)
	FlowStateConsentUnused      = int16(4)
	FlowStateConsentUsed        = int16(5)
)

// Flow is an abstraction used in the persistence layer to unify LoginRequest
// and HandledLoginRequest.
//
// TODO: Deprecate the structs that are made obsolete by the Flow concept.
// Context: Before Flow was introduced, the API and the database used the same
// structs, LoginRequest and HandledLoginRequest. These two tables and structs
// were merged into a new concept, Flow, in order to optimize the persistence
// layer. We currently limit the use of Flow to the persistence layer and keep
// using the original structs in the API in order to minimize the impact of the
// database refactoring on the API.
//
type Flow struct {
	// ID is the identifier ("login challenge") of the login request. It is used to
	// identify the session.
	//
	// required: true
	ID string `json:"challenge" db:"challenge"`

	// RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client.
	//
	// required: true
	RequestedScope sqlxx.StringSlicePipeDelimiter `json:"requested_scope" db:"requested_scope"`

	// RequestedScope contains the access token audience as requested by the OAuth 2.0 Client.
	//
	// required: true
	RequestedAudience sqlxx.StringSlicePipeDelimiter `json:"requested_access_token_audience" db:"requested_at_audience"`

	// Skip, if true, implies that the client has requested the same scopes from the same user previously.
	// If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.
	//
	// This feature allows you to update / set session information.
	//
	// required: true
	Skip bool `json:"skip" db:"skip"`

	// Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope
	// requested by the OAuth 2.0 client. If this value is set and `skip` is true, you MUST include this subject type
	// when accepting the login request, or the request will fail.
	//
	// required: true
	Subject string `json:"subject" db:"subject"`

	// OpenIDConnectContext provides context for the (potential) OpenID Connect context. Implementation of these
	// values in your app are optional but can be useful if you want to be fully compliant with the OpenID Connect spec.
	OpenIDConnectContext *consent.OpenIDConnectContext `json:"oidc_context" db:"oidc_context"`

	// Client is the OAuth 2.0 Client that initiated the request.
	//
	// required: true
	Client *client.Client `json:"client" db:"-"`

	ClientID string `json:"-" db:"client_id"`

	// RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which
	// initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but
	// might come in handy if you want to deal with additional request parameters.
	//
	// required: true
	RequestURL string `json:"request_url" db:"request_url"`

	// SessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag)
	// this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false)
	// this will be a new random value. This value is used as the "sid" parameter in the ID Token and in OIDC Front-/Back-
	// channel logout. It's value can generally be used to associate consecutive login requests by a certain user.
	SessionID sqlxx.NullString `json:"session_id" db:"login_session_id"`

	Verifier string `json:"-" db:"verifier"`
	CSRF     string `json:"-" db:"csrf"`

	LoginInitializedAt sqlxx.NullTime `json:"-" db:"login_initialized_at"`
	RequestedAt        time.Time      `json:"-" db:"requested_at"`

	State int16 `json:"state" db:"state"`

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

	// AMR sets the Authentication Methods References value for this
	// authentication session. You can use it to specify the method a user used to
	// authenticate. For example, if the acr indicates a user used two factor
	// authentication, the amr can express they used a software-secured key.
	AMR sqlxx.StringSlicePipeDelimiter `json:"amr" db:"amr"`

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

	// If set to true means that the request was already handled. This
	// can happen on form double-submit or other errors. If this is set
	// we recommend redirecting the user to `request_url` to re-initiate
	// the flow.
	WasHandled bool `json:"-" db:"was_used"`

	Error                *consent.RequestDeniedError `json:"-" db:"error"`
	LoginAuthenticatedAt sqlxx.NullTime              `json:"-" db:"login_authenticated_at"`
}

func NewFlow(r *consent.LoginRequest) *Flow {
	return &Flow{
		ID:                     r.ID,
		RequestedScope:         r.RequestedScope,
		RequestedAudience:      r.RequestedAudience,
		Skip:                   r.Skip,
		Subject:                r.Subject,
		OpenIDConnectContext:   r.OpenIDConnectContext,
		Client:                 r.Client,
		ClientID:               r.ClientID,
		RequestURL:             r.RequestURL,
		SessionID:              r.SessionID,
		WasHandled:             r.WasHandled,
		ForceSubjectIdentifier: r.ForceSubjectIdentifier,
		Verifier:               r.Verifier,
		CSRF:                   r.CSRF,
		LoginAuthenticatedAt:   r.AuthenticatedAt,
		RequestedAt:            r.RequestedAt,
		State:                  FlowStateLoginInitialized,
	}
}

func (f *Flow) HandleLoginRequest(h *consent.HandledLoginRequest) error {
	// TODO Validate HandledLoginRequest. This will require updating tests that expect the flow to fail at a later point.
	if f.State != FlowStateLoginInitialized {
		return errors.Errorf("invalid flow state: expected %d, got %d", FlowStateLoginInitialized, f.State)
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

	f.State = FlowStateLoginUnused // TODO FlowStateError if h.Error != nil
	f.ID = h.ID
	f.Subject = h.Subject
	f.ForceSubjectIdentifier = h.ForceSubjectIdentifier
	f.Error = h.Error

	f.Remember = h.Remember
	f.RememberFor = h.RememberFor
	f.ACR = h.ACR
	f.AMR = h.AMR
	f.Context = h.Context
	f.WasHandled = h.WasHandled
	f.LoginAuthenticatedAt = h.AuthenticatedAt
	return nil
}

func (f *Flow) InitializeConsent() error {
	if f.State != FlowStateLoginUnused {
		return errors.Errorf("invalid flow state: expected %d, got %d", FlowStateLoginUnused, f.State)
	}
	if f.WasHandled {
		return errors.New("login verifier has already been used")
	}
	f.WasHandled = true
	f.State = FlowStateConsentInitialized
	return nil
}

func (f *Flow) GetHandledLoginRequest() consent.HandledLoginRequest {
	return consent.HandledLoginRequest{
		ID:                     f.ID,
		Remember:               f.Remember,
		RememberFor:            f.RememberFor,
		ACR:                    f.ACR,
		AMR:                    f.AMR,
		Subject:                f.Subject,
		ForceSubjectIdentifier: f.ForceSubjectIdentifier,
		Context:                f.Context,
		WasHandled:             f.WasHandled,
		Error:                  f.Error,
		LoginRequest:           f.GetLoginRequest(),
		RequestedAt:            f.RequestedAt,
		AuthenticatedAt:        f.LoginAuthenticatedAt,
	}
}

func (f *Flow) GetLoginRequest() *consent.LoginRequest {
	return &consent.LoginRequest{
		ID:                     f.ID,
		RequestedScope:         f.RequestedScope,
		RequestedAudience:      f.RequestedAudience,
		Skip:                   f.Skip,
		Subject:                f.Subject,
		OpenIDConnectContext:   f.OpenIDConnectContext,
		Client:                 f.Client,
		ClientID:               f.ClientID,
		RequestURL:             f.RequestURL,
		SessionID:              f.SessionID,
		WasHandled:             f.WasHandled,
		ForceSubjectIdentifier: f.ForceSubjectIdentifier,
		Verifier:               f.Verifier,
		CSRF:                   f.CSRF,
		AuthenticatedAt:        f.LoginAuthenticatedAt,
		RequestedAt:            f.RequestedAt,
	}
}

func (_ Flow) TableName() string {
	return "hydra_oauth2_flow"
}

func (f *Flow) FindInDB(c *pop.Connection, id string) error {
	return c.Find(f, id)
}

func (f *Flow) BeforeSave(_ *pop.Connection) error {
	if f.Client != nil {
		f.ClientID = f.Client.OutfacingID
	}
	if f.State == FlowStateLoginUnused && string(f.Context) == "" {
		f.Context = sqlxx.JSONRawMessage("{}")
	}
	return nil
}

// TODO Populate the client field in FindInDB in order to avoid accessing the database twice.
func (f *Flow) AfterFind(c *pop.Connection) error {
	f.Client = &client.Client{}
	return sqlcon.HandleError(c.Where("id = ?", f.ClientID).First(f.Client))
}
