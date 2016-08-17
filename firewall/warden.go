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

// Firewall offers various validation strategies for access tokens.
type Firewall interface {
	// TokenValid checks if the given token is valid and if the requested scopes are satisfied. Returns
	// a context if the token is valid and an error if not.
	//
	//  ctx, err := firewall.TokenValid(context.Background(), "access-token", "photos", "files")
	//  fmt.Sprintf("%s", ctx.Subject)
	TokenValid(ctx context.Context, token string, scopes ...string) (*Context, error)

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
	//    ctx, err := firewall.TokenValid(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
	//    fmt.Sprintf("%s", ctx.Subject)
	//  }
	TokenFromRequest(r *http.Request) string
}
