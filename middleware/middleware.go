package middleware

import (
	chd "github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/handler"
)

type Middleware interface {
	IsAuthorized(resource, permission string, environment *Env) func(chd.ContextHandler) chd.ContextHandler
	IsAuthenticated(next chd.ContextHandler) chd.ContextHandler
}
