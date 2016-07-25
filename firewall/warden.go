package firewall

import (
	"net/http"
	"time"

	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

type Context struct {
	Subject       string    `json:"sub"`
	GrantedScopes []string  `json:"scopes"`
	Issuer        string    `json:"iss"`
	Audience      string    `json:"aud"`
	IssuedAt      time.Time `json:"iat"`
	ExpiresAt     time.Time `json:"exp"`
}

type Firewall interface {
	// Authorized checks if the given token is valid and if the requested scopes are satisfied and returns
	// a context if yes and an error if no.
	Authorized(ctx context.Context, token string, scopes ...string) (*Context, error)

	// HTTPAuthorized checks if the given HTTP request is authorized and if the requested scopes are satisfied and returns
	// a context if yes and an error if no.
	HTTPAuthorized(ctx context.Context, r *http.Request, scopes ...string) (*Context, error)

	// ActionAllowed, apart from doing the same thing as Authorized, checks if the token's subject is allowed to perform
	// the given action on the given resource. Returns an error if any of the parameters are not fulfilled.
	ActionAllowed(ctx context.Context, token string, accessRequest *ladon.Request, scopes ...string) (*Context, error)

	// HTTPActionAllowed, apart from doing the same thing as Authorized, checks if the token's subject (extracted from the HTTP Request) is allowed to perform
	// the given action on the given resource. Returns an error if any of the parameters are not fulfilled.
	HTTPActionAllowed(ctx context.Context, r *http.Request, accessRequest *ladon.Request, scopes ...string) (*Context, error)
}
