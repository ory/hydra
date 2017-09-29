package oauth2

import "time"

// ConsentRequest represents a consent request.
//
// swagger:model oAuth2consentRequest
type ConsentRequest struct {
	// ID is the id of this consent request.
	ID string `json:"id"`

	// RequestedScopes represents a list of scopes that have been requested by the OAuth2 request initiator.
	RequestedScopes []string `json:"requestedScopes,omitempty"`

	// Audience is the client id that initiated the OAuth2 request.
	Audience string `json:"audience"`

	// ExpiresAt is the time where the access request will expire.
	ExpiresAt time.Time `json:"expiresAt"`

	// Redirect URL is the URL where the user agent should be redirected to after the consent has been
	// accepted or rejected.
	RedirectURL string `json:"redirectUrl"`

	CSRF             string                 `json:"-"`
	GrantedScopes    []string               `json:"-"`
	Subject          string                 `json:"-"`
	AccessTokenExtra map[string]interface{} `json:"-"`
	IDTokenExtra     map[string]interface{} `json:"-"`
	Consent          string                 `json:"-"`
	DenyReason       string                 `json:"-"`
}

func (c *ConsentRequest) IsConsentGranted() bool {
	return c.Consent == ConsentRequestAccepted
}

// AcceptConsentRequestPayload represents data that will be used to accept a consent request.
//
// swagger:model consentRequestAcceptance
type AcceptConsentRequestPayload struct {
	// AccessTokenExtra represents arbitrary data that will be added to the access token and that will be returned
	// on introspection and warden requests.
	AccessTokenExtra map[string]interface{} `json:"accessTokenExtra"`

	// IDTokenExtra represents arbitrary data that will be added to the ID token. The ID token will only be issued
	// if the user agrees to it and if the client requested an ID token.
	IDTokenExtra map[string]interface{} `json:"idTokenExtra"`

	// Subject represents a unique identifier of the user (or service, or legal entity, ...) that accepted the
	// OAuth2 request.
	Subject string `json:"subject"`

	// A list of scopes that the user agreed to grant. It should be a subset of requestedScopes from the consent request.
	GrantScopes []string `json:"grantScopes"`
}

// RejectConsentRequestPayload represents data that will be used to reject a consent request.
//
// swagger:model consentRequestRejection
type RejectConsentRequestPayload struct {
	// Reason represents the reason why the user rejected the consent request.
	Reason string `json:"reason"`
}

type ConsentRequestManager interface {
	PersistConsentRequest(*ConsentRequest) error
	AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error
	RejectConsentRequest(id string, payload *RejectConsentRequestPayload) error
	GetConsentRequest(id string) (*ConsentRequest, error)
}
