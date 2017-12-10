// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import "time"

// ConsentRequest represents a consent request.
type ConsentRequest struct {
	// ID is the id of this consent request.
	ID string `json:"id"`

	// RequestedScopes represents a list of scopes that have been requested by the OAuth2 request initiator.
	RequestedScopes []string `json:"requestedScopes"`

	// ClientID is the client id that initiated the OAuth2 request.
	ClientID string `json:"clientId"`

	// ExpiresAt is the time where the access request will expire.
	ExpiresAt time.Time `json:"expiresAt"`

	// Redirect URL is the URL where the user agent should be redirected to after the consent has been
	// accepted or rejected.
	RedirectURL string `json:"redirectUrl"`

	RequestedACR    []string `json:"requestedAcr"`
	RequestedPrompt string   `json:"requestedPrompt"`
	RequestedMaxAge int64    `json:"requestedMaxAge"`

	CSRF             string                 `json:"-"`
	GrantedScopes    []string               `json:"-"`
	Subject          string                 `json:"-"`
	AccessTokenExtra map[string]interface{} `json:"-"`
	IDTokenExtra     map[string]interface{} `json:"-"`
	Consent          string                 `json:"-"`
	DenyReason       string                 `json:"-"`
	DenyError        string                 `json:"-"`
	AuthTime         int64                  `json:"-"`
	ProvidedACR      string                 `json:"-"`
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

	// ProvidedAuthenticationContextClassReference specifies an Authentication Context Class Reference value that identifies
	// the Authentication Context Class that the authentication performed satisfied. The value "0" indicates the End-User
	// authentication did not meet the requirements of ISO/IEC 29115 [ISO29115] level 1.
	//
	// In summary ISO/IEC 29115 defines four levels, broadly summarized as follows.
	//
	// * acr=0 does not satisfy Level 1 and could be, for example, authentication using a long-lived browser cookie.
	// * Level 1 (acr=1): Minimal confidence in the asserted identity of the entity, but enough confidence that the
	// 		entity is the same over consecutive authentication events. For example presenting a self-registered
	// 		username or password.
	// * Level 2 (acr=2): There is some confidence in the asserted identity of the entity. For example confirming
	// 		authentication using a mobile app ("Something you have").
	// * Level 3 (acr=3): High confidence in an asserted identity of the entity. For example sending a code to a mobile
	// 		phone or using Google Authenticator or a fingerprint scanner ("Something you have and something you know" / "Something you are")
	// * Level 4 (acr=4): Very high confidence in an asserted identity of the entity. Requires in-person identification.
	ProvidedAuthenticationContextClassReference string `json:"providedAcr"`

	// AuthTime is the time when the End-User authentication occurred. Its value is a JSON number representing the
	// number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time.
	AuthTime int64 `json:"authTime"`
}

// RejectConsentRequestPayload represents data that will be used to reject a consent request.
//
// swagger:model consentRequestRejection
type RejectConsentRequestPayload struct {
	// Reason represents the reason why the user rejected the consent request.
	Reason string `json:"reason"`

	// Error can be used to return an OpenID Connect or OAuth 2.0 error to the OAuth 2.0 client, such as login_required,
	// interaction_required, consent_required.
	Error string `json:"error"`
}

type ConsentRequestManager interface {
	PersistConsentRequest(*ConsentRequest) error
	AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error
	RejectConsentRequest(id string, payload *RejectConsentRequestPayload) error
	GetConsentRequest(id string) (*ConsentRequest, error)
}
