package oauth2

import "time"

type ConsentRequest struct {
	ID               string                 `json:"id"`
	CSRF             string                 `json:"-"`
	GrantedScopes    []string               `json:"-"`
	RequestedScope   []string               `json:"requested_scope,omitempty"`
	Audience         string                 `json:"audience"`
	Subject          string                 `json:"-"`
	ExpiresAt        time.Time              `json:"expires_at"`
	RedirectURL      string                 `json:"redirect_url"`
	AccessTokenExtra map[string]interface{} `json:"-"`
	IDTokenExtra     map[string]interface{} `json:"-"`
	Consent          string                 `json:"-"`
}

func (c *ConsentRequest) IsConsentGranted() bool {
	return c.Consent == ConsentRequestAccepted
}

type AcceptConsentRequestPayload struct {
	AccessTokenExtra map[string]interface{} `json:"access_token_extra"`
	IDTokenExtra     map[string]interface{} `json:"id_token_extra"`
	Subject          string                 `json:"subject"`
	GrantScopes      []string               `json:"grant_scopes"`
}

type ConsentRequestClient interface {
	AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error
	RejectConsentRequest(id string) error
	GetConsentRequest(id string) (*ConsentRequest, error)
}

type ConsentRequestManager interface {
	PersistConsentRequest(*ConsentRequest) error
	ConsentRequestClient
}
