package warden

import (
	"net/http"

	"github.com/ory-am/ladon"
	"time"
)

type Context struct {
	Subject       string   `json:"subject"`
	GrantedScopes []string `json:"scopes"`
	Issuer        string   `json:"issuer"`
	Audience      string   `json:"audience"`
	IssuedAt time.Time
	ExpiresAt time.Time
}

type Warden interface {
	ActionAllowed(token string, *ladon.Request) (*Context, error)

	Authorized(token string, scopes ...string) (*Context, error)

	HTTPAuthorized(r *http.Request, scopes ...string) (*Context, error)

	HTTPActionAllowed(r *http.Request, *ladon.Request, scopes ...string) (*Context, error)
}
