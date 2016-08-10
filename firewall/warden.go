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
	// Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app.
	// This is usually a uuid but you can choose a urn or some other id too.
	Subject       string                 `json:"sub"`

	// GrantedScopes is a list of scopes that the subject authorized when asked for consent.
	GrantedScopes []string               `json:"scopes"`

	// Issuer is the id of the issuer, typically an hydra instance.
	Issuer        string                 `json:"iss"`

	// Audience is who the token was issued for. This is an OAuth2 app usually.
	Audience      string                 `json:"aud"`

	// IssuedAt is the token creation time stamp.
	IssuedAt      time.Time              `json:"iat"`

	// ExpiresAt is the expiry timestamp.
	ExpiresAt     time.Time              `json:"exp"`

	// Extra represents arbitrary session data.
	Extra         map[string]interface{} `json:"ext"`
}

// Firewall offers various validation strategies for access tokens.
type Firewall interface {
	Introspector

	// InspectToken checks if the given token is valid and if the requested scopes are satisfied. Returns
	// a context if the token is valid and an error if not.
	//
	//  ctx, err := firewall.InspectToken(context.Background(), "access-token", "photos", "files")
	//  fmt.Sprintf("%s", ctx.Subject)
	InspectToken(ctx context.Context, token string, scopes ...string) (*Context, error)

	// IsAllowed uses policies to return nil if the access request can be fulfilled or an error if not.
	//
	//  ctx, err := firewall.IsAllowed(context.Background(), &ladon.Request{
	//    Subject:  "alice",
	//    Resource: "matrix",
	//    Action:   "create",
	//    Context:  ladon.Context{},
	//  }, "photos", "files")
	//
	//  fmt.Sprintf("%s", ctx.Subject)
	IsAllowed(ctx context.Context, accessRequest *ladon.Request) error

	// TokenAllowed uses policies and a token to return a context and no error if the access request can be fulfilled or an error if not.
	//
	//  ctx, err := firewall.TokenAllowed(context.Background(), "access-token", &ladon.Request{
	//    Resource: "matrix",
	//    Action:   "create",
	//    Context:  ladon.Context{},
	//  }, "photos", "files")
	//
	//  fmt.Sprintf("%s", ctx.Subject)
	TokenAllowed(ctx context.Context, token string, accessRequest *ladon.Request, scopes ...string) (*Context, error)

	// TokenFromRequest returns an access token from the HTTP Authorization header.
	//
	//  func anyHttpHandler(w http.ResponseWriter, r *http.Request) {
	//    ctx, err := firewall.InspectToken(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
	//    fmt.Sprintf("%s", ctx.Subject)
	//  }
	TokenFromRequest(r *http.Request) string
}

// Introspection contains an access token's session data as specified by IETF RFC 7662.
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
	Active    bool   `json:"active"`

	// Scope is a JSON string containing a space-separated list of
	// scopes associated with this token
	Scope     string `json:"scope,omitempty"`

	// ClientID is aclient identifier for the OAuth 2.0 client that
	// requested this token.
	ClientID  string `json:"client_id,omitempty"`

	// Subject of the token, as defined in JWT [RFC7519].
	// Usually a machine-readable identifier of the resource owner who
	// authorized this token.
	Subject   string `json:"sub,omitempty"`

	// Expires at is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token will expire
	ExpiresAt int64  `json:"exp,omitempty"`

	// Issued at is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token was
	// originally issued
	IssuedAt  int64  `json:"iat,omitempty"`

	// NotBefore is an integer timestamp, measured in the number of seconds
	// since January 1 1970 UTC, indicating when this token is not to be
	// used before
	NotBefore int64  `json:"nbf,omitempty"`

	// Username is a human-readable identifier for the resource owner who
	// authorized this token.
	Username  int64  `json:"username,omitempty"`

	// Audience is a service-specific string identifier or list of string
	// identifiers representing the intended audience for this token
	Audience  string `json:"aud,omitempty"`

	// Issuer is a string representing the issuer of this token
	Issuer    string `json:"iss,omitempty"`
}

// Introspector is capable of introspecting an access token according to IETF RFC 7662, see:
// https://tools.ietf.org/html/rfc7662
type Introspector interface {
	// IntrospectToken performs a token introspection according to IETF RFC 7662, see: https://tools.ietf.org/html/rfc7662
	//
	//  func anyHttpHandler(w http.ResponseWriter, r *http.Request) {
	//    ctx, err := firewall.InspectToken(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
	//    fmt.Sprintf("%s", ctx.Subject)
	//  }
	IntrospectToken(ctx context.Context, token string) (*Introspection, error)
}
