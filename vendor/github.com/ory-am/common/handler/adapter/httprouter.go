package adapter

import (
	"github.com/julienschmidt/httprouter"
	cc "github.com/ory-am/common/context"
	"github.com/ory-am/common/handler"
	"golang.org/x/net/context"
	"net/http"
)

type httpRouterAdapter struct {
	ctx        context.Context
	final      handler.ContextHandler
	middleware []handler.Middleware
}

func (ca *httpRouterAdapter) Handle(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	final := ca.final
	for i := len(ca.middleware) - 1; i >= 0; i-- {
		final = ca.middleware[i](final)
	}
	final.ServeHTTPContext(cc.NewContextFromRouterParams(ca.ctx, params), rw, req)
}

func (ca *httpRouterAdapter) Then(final handler.ContextHandler) httprouter.Handle {
	ca.final = final
	return ca.Handle
}

func (ca *httpRouterAdapter) ThenFunc(final handler.ContextHandlerFunc) httprouter.Handle {
	ca.final = handler.ContextHandler(final)
	return ca.Handle
}

func NewHttpRouterAdapter(ctx context.Context, middleware ...handler.Middleware) *httpRouterAdapter {
	adapter := &httpRouterAdapter{}
	adapter.ctx = ctx
	adapter.middleware = middleware
	return adapter
}
