package context

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/net/context"
	"net/http"
)

type key int

type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	h(ctx, rw, req)
}

type Middleware func(next ContextHandler) ContextHandler

// Implements http.Handler
type contextAdapter struct {
	ctx        context.Context
	final      ContextHandler
	middleware []Middleware
}

func NewContextAdapter(ctx context.Context, middleware ...Middleware) *contextAdapter {
	return &contextAdapter{
		ctx:        ctx,
		middleware: append([]Middleware{}, middleware...),
	}
}

func (ca *contextAdapter) Then(final ContextHandler) *contextAdapter {
	ca.final = final
	return ca
}

func (ca *contextAdapter) ThenFunc(final ContextHandlerFunc) *contextAdapter {
	ca.final = ContextHandler(final)
	return ca
}

func (ca *contextAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	final := ca.final
	for i := len(ca.middleware) - 1; i >= 0; i-- {
		final = ca.middleware[i](final)
	}
	final.ServeHTTPContext(ca.ctx, rw, req)
}
