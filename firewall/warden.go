// Package firewall defines an API for validating access requests.
package firewall

import (
	"net/http"
	"time"

	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

// Context contains an access token's session data
type Context struct {
	Subject       string                 `json:"sub"`
	GrantedScopes []string               `json:"scopes"`
	Issuer        string                 `json:"iss"`
	Audience      string                 `json:"aud"`
	IssuedAt      time.Time              `json:"iat"`
	ExpiresAt     time.Time              `json:"exp"`
	Extra         map[string]interface{} `json:"ext"`
}

// Firewall offers various validation strategies for access tokens.
type Firewall interface {
	Introspector

	// InspectToken checks if the given token is valid and if the requested scopes are satisfied. Returns
	// a context if the token is valid and an error if not.
	InspectToken(ctx context.Context, token string, scopes ...string) (*Context, error)

	// IsAllowed uses policies to return nil if the access request can be fulfilled or an error if not.
	IsAllowed(ctx context.Context, accessRequest *ladon.Request) error

	// TokenAllowed uses policies and a token to return a context and no error if the access request can be fulfilled or an error if not.
	TokenAllowed(ctx context.Context, token string, accessRequest *ladon.Request, scopes ...string) (*Context, error)

	// TokenFromRequest returns an access token from the HTTP Authorization header.
	TokenFromRequest(r *http.Request) string
}

// Introspection contains an access token's session data as specified by IETF RFC 7662.
type Introspection struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Username  int64  `json:"username,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Issuer    string `json:"iss,omitempty"`
}

// Introspector is capable of introspecting an access token according to IETF RFC 7662.
type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*Introspection, error)
}
