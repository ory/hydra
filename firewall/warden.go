// Package firewall defines an API for validating access requests.
package firewall

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// Context contains an access token's session data
type Context struct {
	// Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app.
	// This is usually a uuid but you can choose a urn or some other id too.
	Subject string `json:"sub"`

	// GrantedScopes is a list of scopes that the subject authorized when asked for consent.
	GrantedScopes []string `json:"scopes"`

	// Issuer is the id of the issuer, typically an hydra instance.
	Issuer string `json:"iss"`

	// Audience is who the token was issued for. This is an OAuth2 app usually.
	Audience string `json:"aud"`

	// IssuedAt is the token creation time stamp.
	IssuedAt time.Time `json:"iat"`

	// ExpiresAt is the expiry timestamp.
	ExpiresAt time.Time `json:"exp"`

	// Extra represents arbitrary session data.
	Extra map[string]interface{} `json:"ext"`
}

// AccessRequest is the warden's request object.
type AccessRequest struct {
	// Resource is the resource that access is requested to.
	Resource string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action string `json:"action"`

	// Subejct is the subject that is requesting access.
	Subject string `json:"subject"`

	// Context is the request's environmental context.
	Context map[string]interface{} `json:"context"`
}

type TokenAccessRequest struct {
	// Resource is the resource that access is requested to.
	Resource string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action string `json:"action"`

	// Context is the request's environmental context.
	Context map[string]interface{} `json:"context"`
}

// Firewall offers various validation strategies for access tokens.
type Firewall interface {
	// IsAllowed uses policies to return nil if the access request can be fulfilled or an error if not.
	//
	//  ctx, err := firewall.IsAllowed(context.Background(), &AccessRequest{
	//    Subject:  "alice",
	//    Resource: "matrix",
	//    Action:   "create",
	//    Context:  ladon.Context{},
	//  }, "photos", "files")
	//
	//  fmt.Sprintf("%s", ctx.Subject)
	IsAllowed(ctx context.Context, accessRequest *AccessRequest) error

	// TokenAllowed uses policies and a token to return a context and no error if the access request can be fulfilled or an error if not.
	//
	//  ctx, err := firewall.TokenAllowed(context.Background(), "access-token", &TokenAccessRequest{
	//    Resource: "matrix",
	//    Action:   "create",
	//    Context:  ladon.Context{},
	//  }, "photos", "files")
	//
	//  fmt.Sprintf("%s", ctx.Subject)
	TokenAllowed(ctx context.Context, token string, accessRequest *TokenAccessRequest, scopes ...string) (*Context, error)

	// TokenFromRequest returns an access token from the HTTP Authorization header.
	//
	//  func anyHttpHandler(w http.ResponseWriter, r *http.Request) {
	//    ctx, err := firewall.TokenValid(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
	//    fmt.Sprintf("%s", ctx.Subject)
	//  }
	TokenFromRequest(r *http.Request) string
}
