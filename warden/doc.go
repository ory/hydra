// Package warden implements endpoints capable of making access control decisions based on Access Control Policies
package warden

import (
	"github.com/ory/hydra/firewall"
)

// The warden access request response
// swagger:response wardenAccessRequestResponse
type swaggerWardenAccessRequestResponseParameters struct {
	// in: body
	Body swaggerWardenAccessRequestResponse
}

// The warden access request response
// swagger:model wardenAccessRequestResponse
type swaggerWardenAccessRequestResponse struct {
	// Allowed is true if the request is allowed and false otherwise.
	Allowed bool `json:"allowed"`
}

// swagger:parameters doesWardenAllowAccessRequest
type swaggerDoesWardenAllowAccessRequestParameters struct {
	// in: body
	Body firewall.AccessRequest
}

// swagger:parameters doesWardenAllowTokenAccessRequest
type swaggerDoesWardenAllowTokenAccessRequestParameters struct {
	// in: body
	Body swaggerWardenTokenAccessRequest
}

// swagger:model wardenTokenAccessRequest
type swaggerWardenTokenAccessRequest struct {
	// Scopes is an array of scopes that are requried.
	Scopes []string `json:"scopes"`

	// Token is the token to introspect.
	Token string `json:"token"`

	// Resource is the resource that access is requested to.
	Resource string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action string `json:"action"`

	// Context is the request's environmental context.
	Context map[string]interface{} `json:"context"`
}

// The warden access request (with token) response
// swagger:response wardenTokenAccessRequestResponse
type swaggerWardenTokenAccessRequestResponse struct {
	// in: body
	Body swaggerWardenTokenAccessRequestResponsePayload
}

// The warden access request (with token) response
// swagger:model wardenTokenAccessRequestResponse
type swaggerWardenTokenAccessRequestResponsePayload struct {
	// Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app.
	// This is usually a uuid but you can choose a urn or some other id too.
	Subject string `json:"subject"`

	// GrantedScopes is a list of scopes that the subject authorized when asked for consent.
	GrantedScopes []string `json:"grantedScopes"`

	// Issuer is the id of the issuer, typically an hydra instance.
	Issuer string `json:"issuer"`

	// ClientID is the id of the OAuth2 client that requested the token.
	ClientID string `json:"clientId"`

	// IssuedAt is the token creation time stamp.
	IssuedAt string `json:"issuedAt"`

	// ExpiresAt is the expiry timestamp.
	ExpiresAt string `json:"expiresAt"`

	// Extra represents arbitrary session data.
	Extra map[string]interface{} `json:"accessTokenExtra"`

	// Allowed is true if the request is allowed and false otherwise.
	Allowed bool `json:"allowed"`
}
