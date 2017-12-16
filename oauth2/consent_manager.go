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

import (
	"time"

	"github.com/ory/hydra/client"
)

// ConsentRequestOpenIDConnectContext represents parsed and sanitized OpenID Conenct request information
//
// swagger:model oAuth2ConsentRequestOIDCContext
type ConsentRequestOpenIDConnectContext struct {
	// Display is an ASCII string value that specifies how the Authorization Server displays the authentication and consent user interface pages to the End-User.
	//
	// For more information head over to http://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	Display string `json:"display,omitempty"`

	// Prompt is a Space delimited, case sensitive list of ASCII string values that specifies whether the Authorization Server prompts the End-User for reauthentication and consent.
	// This value will either be empty or "consent" - no other values will be transmitted.
	//
	// For more information head over to http://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	Prompt string `json:"prompt,omitempty"`

	// UILocales is the End-User's preferred languages and scripts for the user interface, represented as a
	// space-separated list of BCP47 [RFC5646] language tag values, ordered by preference. For instance, the value
	// "fr-CA fr en" represents a preference for French as spoken in Canada, then French (without a region designation), followed by English (without a region designation)
	//
	// For more information head over to http://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	UILocales string `json:"ui_locales,omitempty"`

	// LoginHint hints to the Authorization Server about the login identifier the End-User might use to log in (if necessary).
	// This hint can be used by an RP if it first asks the End-User for their e-mail address (or other identifier)
	// and then wants to pass that value as a hint to the discovered authorization service. It is RECOMMENDED that the
	// hint value match the value used for discovery. This value MAY also be a phone number in the format specified for
	// the phone_number Claim. The use of this parameter is left to the OP's discretion.
	//
	// For more information head over to http://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	LoginHint string `json:"login_hint,omitempty"`

	// ACRValues defines Requested Authentication Context Class Reference values. It is a Space-separated string that specifies
	// the acr values that the Authorization Server is being requested to use for processing this Authentication
	// Request, with the values appearing in order of preference.
	//
	// For more information head over to http://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	ACRValues []string `json:"acr_values,omitempty"`
}

// ConsentRequest represents a consent request
//
// swagger:model oAuth2ConsentRequest
type ConsentRequest struct {
	// ID is the id of this consent request.
	ID string `json:"id"`

	// RequestedScopes represents a list of scopes that have been requested by the OAuth2 request initiator.
	RequestedScopes []string `json:"requestedScopes"`

	// ClientID is the client id that initiated the OAuth2 request.
	ClientID string `json:"clientId"`

	// Client is the client that initiated the OAuth2 request. It contains all fields associated with a client
	// and can be used to display extended information in the consent user interface.
	Client *client.Client `json:"client"`

	// ExpiresAt is the time where the access request will expire.
	ExpiresAt time.Time `json:"expiresAt"`

	// Redirect URL is the URL where the user agent should be redirected to after the consent has been
	// accepted or rejected.
	RedirectURL string `json:"redirectUrl"`

	// OpenIDConnectContext is a payload intended to help decipher OpenID Connect Requests. While it is theoretically possible to
	// extract this information from the `redirectUrl` value yourself, using this context instead is strongly recommended.
	OpenIDConnectContext *ConsentRequestOpenIDConnectContext `json:"oidcContext"`

	RequestedAt      time.Time              `json:"-"`
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

// AcceptConsentRequestPayload represents data that will be used to accept a consent request
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

// RejectConsentRequestPayload represents data that will be used to reject a consent request
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
	GetPreviouslyGrantedConsent(subject string, client string, scopes []string) (*ConsentRequest, error)
}

type byTime []ConsentRequest

func (s byTime) Len() int {
	return len(s)
}

func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byTime) Less(i, j int) bool {
	return s[i].RequestedAt.After(s[j].RequestedAt)
}

func isSubset(subset, set []string) bool {
	values := make(map[string]int)
	for _, value := range set {
		values[value] += 1
	}

	for _, value := range subset {
		if count, found := values[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			values[value] = count - 1
		}
	}

	return true
}
