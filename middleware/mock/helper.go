package mock

import (
	chd "github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/handler"
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/net/context"
	"net/http"
)

var MockFailAuthenticationHandler = chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusUnauthorized)
})

var MockPassAuthenticationHandler = func(next chd.ContextHandler) chd.ContextHandlerFunc {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTPContext(ctx, rw, req)
	}
}
