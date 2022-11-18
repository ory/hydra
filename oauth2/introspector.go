// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

// Introspection contains an access token's session data as specified by
// [IETF RFC 7662](https://tools.ietf.org/html/rfc7662)
//
// swagger:model introspectedOAuth2Token
type Introspection struct {
	// Active is a boolean indicator of whether or not the presented token
	// is currently active.  The specifics of a token's "active" state
	// will vary depending on the implementation of the authorization
	// server and the information it keeps about its tokens, but a "true"
	// value return for the "active" property will generally indicate
	// that a given token has been issued by this authorization server,
	// has not been revoked by the resource owner, and is within its
	// given time window of validity (e.g., after its issuance time and
	// before its expiration time).
	//
	// required: true
	Active bool `json:"active"`

	// Scope is a JSON string containing a space-separated list of
	// scopes associated with this token.
	Scope string `json:"scope,omitempty"`

	// ID is aclient identifier for the OAuth 2.0 client that
	// requested this token.
	ClientID string `json:"client_id"`

	// Subject of the token, as defined in JWT [RFC7519].
	// Usually a machine-readable identifier of the resource owner who
	// authorized this token.
	Subject string `json:"sub"`

	// ObfuscatedSubject is set when the subject identifier algorithm was set to "pairwise" during authorization.
	// It is the `sub` value of the ID Token that was issued.
	ObfuscatedSubject string `json:"obfuscated_subject,omitempty"`

	// Expires at is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token will expire.
	ExpiresAt int64 `json:"exp"`

	// Issued at is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token was
	// originally issued.
	IssuedAt int64 `json:"iat"`

	// NotBefore is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token is not to be
	// used before.
	NotBefore int64 `json:"nbf"`

	// Username is a human-readable identifier for the resource owner who
	// authorized this token.
	Username string `json:"username,omitempty"`

	// Audience contains a list of the token's intended audiences.
	Audience []string `json:"aud"`

	// IssuerURL is a string representing the issuer of this token
	Issuer string `json:"iss"`

	// TokenType is the introspected token's type, typically `Bearer`.
	TokenType string `json:"token_type"`

	// TokenUse is the introspected token's use, for example `access_token` or `refresh_token`.
	TokenUse string `json:"token_use"`

	// Extra is arbitrary data set by the session.
	Extra map[string]interface{} `json:"ext,omitempty"`
}
