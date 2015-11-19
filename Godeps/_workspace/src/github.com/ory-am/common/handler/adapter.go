package handler

import (
	"golang.org/x/net/context"
	"net/http"
)

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
