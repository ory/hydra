package warden

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

type Warden interface {
	Authorized(ctx context.Context, token string, scopes ...string) (*Context, error)
	HTTPAuthorized(ctx context.Context, r *http.Request, scopes ...string) (*Context, error)

	ActionAllowed(ctx context.Context, token string, accessRequest *ladon.Request, scopes ...string) (*Context, error)
	HTTPActionAllowed(ctx context.Context, r *http.Request, accessRequest *ladon.Request, scopes ...string) (*Context, error)
}
