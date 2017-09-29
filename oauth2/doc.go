package oauth2

// swagger:parameters revokeOAuth2Token
type swaggerRevokeOAuth2TokenParameters struct {
	// in: formData
	// required: true
	Token string `json:"token"`
}

// swagger:parameters rejectOAuth2ConsentRequest
type swaggerRejectConsentRequest struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body RejectConsentRequestPayload
}

// swagger:parameters acceptOAuth2ConsentRequest
type swaggerAcceptConsentRequest struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body AcceptConsentRequestPayload
}

// The consent request response
// swagger:response oAuth2ConsentRequest
type swaggerOAuthConsentRequest struct {
	// in: body
	Body ConsentRequest
}

// The token response
// swagger:response oauthTokenResponse
type swaggerOAuthTokenResponse struct {
	// in: body
	Body struct {
		// The lifetime in seconds of the access token.  For
		//  example, the value "3600" denotes that the access token will
		// expire in one hour from the time the response was generated.
		ExpiresIn int `json:"expires_in"`

		// The scope of the access token
		Scope int `json:"scope"`

		// To retrieve a refresh token request the id_token scope.
		IDToken int `json:"id_token"`

		// The access token issued by the authorization server.
		AccessToken string `json:"access_token"`

		// The refresh token, which can be used to obtain new
		// access tokens. To retrieve it add the scope "offline" to your access token request.
		RefreshToken string `json:"refresh_token"`

		// The type of the token issued
		TokenType string `json:"token_type"`
	}
}

// The token introspection response
// swagger:response introspectOAuth2TokenResponse
type swaggerOAuthIntrospectionResponse struct {
	// in: body
	Body swaggerOAuthIntrospectionResponsePayload
}

// swagger:model oAuth2TokenIntrospection
type swaggerOAuthIntrospectionResponsePayload struct {
	// Active is a boolean indicator of whether or not the presented token
	// is currently active.  The specifics of a token's "active" state
	// will vary depending on the implementation of the authorization
	// server and the information it keeps about its tokens, but a "true"
	// value return for the "active" property will generally indicate
	// that a given token has been issued by this authorization server,
	// has not been revoked by the resource owner, and is within its
	// given time window of validity (e.g., after its issuance time and
	// before its expiration time).
	Active bool `json:"active"`

	// Scope is a JSON string containing a space-separated list of
	// scopes associated with this token.
	Scope string `json:"scope,omitempty"`

	// ClientID is aclient identifier for the OAuth 2.0 client that
	// requested this token.
	ClientID string `json:"client_id,omitempty"`

	// Subject of the token, as defined in JWT [RFC7519].
	// Usually a machine-readable identifier of the resource owner who
	// authorized this token.
	Subject string `json:"sub,omitempty"`

	// Expires at is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token will expire.
	ExpiresAt int64 `json:"exp,omitempty"`

	// Issued at is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token was
	// originally issued.
	IssuedAt int64 `json:"iat,omitempty"`

	// NotBefore is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token is not to be
	// used before.
	NotBefore int64 `json:"nbf,omitempty"`

	// Username is a human-readable identifier for the resource owner who
	// authorized this token.
	Username string `json:"username,omitempty"`

	// Audience is a service-specific string identifier or list of string
	// identifiers representing the intended audience for this token.
	Audience string `json:"aud,omitempty"`

	// Issuer is a string representing the issuer of this token
	Issuer string `json:"iss,omitempty"`

	// Extra is arbitrary data set by the session.
	Extra map[string]interface{} `json:"ext,omitempty"`
}

// swagger:parameters introspectOAuth2Token
type swaggerOAuthIntrospectionRequest struct {
	// The string value of the token. For access tokens, this
	// is the "access_token" value returned from the token endpoint
	// defined in OAuth 2.0 [RFC6749], Section 5.1.
	// This endpoint DOES NOT accept refresh tokens for validation.
	//
	// required: true
	// in: formData
	Token string `json:"token"`

	// An optional, space separated list of required scopes. If the access token was not granted one of the
	// scopes, the result of active will be false.
	//
	// in: formData
	Scope string `json:"scope"`
}

// swagger:parameters getOAuth2ConsentRequest acceptConsentRequest rejectConsentRequest
type swaggerOAuthConsentRequestPayload struct {
	// The id of the OAuth 2.0 Consent Request.
	//
	// unique: true
	// required: true
	// in: path
	ID string `json:"id"`
}
