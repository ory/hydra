package client

import (
	"github.com/ory-am/ladon/guard/operator"
	"net/http"
)

type Action struct {
	Resource   string            `json:"resource"`
	Permission string            `json:"permission"`
	Scopes     []string `json:"scopes"`
	Context    *operator.Context `json:"context"`
}

type Context struct {
	Subject string
	Scopes []string
	Issuer string
	Audience string
}

type Client interface {
	ActionAllowed(token string, action *Action) (*Context, error)

	Authorized(token string, scopes ...string) (*Context, error)

	HTTPAuthorized(r *http.Request, scopes ...string) (*Context, error)

	HTTPActionAllowed(r *http.Request, scopes ...string) (*Context, error)
}
