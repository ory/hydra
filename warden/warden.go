package warden

import (
	"github.com/ory-am/ladon/guard/operator"
	"net/http"
	"github.com/ory-am/ladon"
)

type Action struct {
	Resource   string            `json:"resource"`
	Permission string            `json:"permission"`
	Scopes     []string `json:"scopes"`
	Context    *ladon.Context `json:"context"`
}

type Context struct {
	Subject string `json:"subject"`
	Scopes []string `json:"scopes"`
	Issuer string `json:"issuer"`
	Audience string `json:"audience"`
}

type Warden interface {
	ActionAllowed(token string, action *Action) (*Context, error)

	Authorized(token string, scopes ...string) (*Context, error)

	HTTPAuthorized(r *http.Request, scopes ...string) (*Context, error)

	HTTPActionAllowed(r *http.Request, scopes ...string) (*Context, error)
}
