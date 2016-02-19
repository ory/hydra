package client

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/ladon/guard/operator"
	"net/http"
)

type AuthorizeRequest struct {
	Resource   string            `json:"resource"`
	Token      string            `json:"token"`
	Permission string            `json:"permission"`
	Context    *operator.Context `json:"context"`
}

type Client interface {
	IsAllowed(ar *AuthorizeRequest) (bool, error)
	IsRequestAllowed(req *http.Request, resource, permission, owner string) (bool, error)
	IsAuthenticated(token string) (bool, error)
}
