package oauth2

// swagger:parameters revokeOAuthToken
type swaggerCreateClientPayload struct {
	// in: body
	// required: true
	Body struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
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
// swagger:response introspectOAuthTokenResponse
type swaggerOAuthIntrospectionResponse struct {
	// in: body
	Body struct {
		// Boolean indicator of whether or not the presented token
		// is currently active.  The specifics of a token's "active" state
		// will vary depending on the implementation of the authorization
		// server and the information it keeps about its tokens, but a "true"
		// value return for the "active" property will generally indicate
		// that a given token has been issued by this authorization server,
		// has not been revoked by the resource owner, and is within its
		// given time window of validity (e.g., after its issuance time and
		// before its expiration time).
		Active bool `json:"active"`

		// Client identifier for the OAuth 2.0 client that
		// requested this token.
		ClientID string `json:"client_id,omitempty"`

		// A JSON string containing a space-separated list of
		// scopes associated with this token
		Scope string `json:"scope,omitempty"`

		// Integer timestamp, measured in the number of seconds
		// since January 1 1970 UTC, indicating when this token will expire
		ExpiresAt int64 `json:"exp,omitempty"`
		// Integer timestamp, measured in the number of seconds
		// since January 1 1970 UTC, indicating when this token was
		//originally issued
		IssuedAt int64 `json:"iat,omitempty"`

		// Subject of the token, as defined in JWT [RFC7519].
		// Usually a machine-readable identifier of the resource owner who
		// authorized this token.
		Subject string `json:"sub,omitempty"`

		// Human-readable identifier for the resource owner who
		// authorized this token. Currently not supported by Hydra.
		Username string `json:"username,omitempty"`

		// Extra session information set using the at_ext key in the consent response.
		Session Session `json:"sess,omitempty"`
	}
}
