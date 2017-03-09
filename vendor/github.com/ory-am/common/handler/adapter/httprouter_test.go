package adapter

import (
	"github.com/julienschmidt/httprouter"
	cc "github.com/ory-am/common/context"
	"github.com/ory-am/common/handler"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"testing"
)

func TestNewHttpRouterAdapter(t *testing.T) {
	called := 0
	params := httprouter.Params{{"foo", "bar"}}
	NewHttpRouterAdapter(
		context.Background(),
		func(next handler.ContextHandler) handler.ContextHandler {
			return handler.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				if called == 0 {
					called = 1
				}
				next.ServeHTTPContext(ctx, rw, req)
			})
		},
		func(next handler.ContextHandler) handler.ContextHandler {
			return handler.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				if called == 1 {
					called = 2
				}
				next.ServeHTTPContext(ctx, rw, req)
			})
		},
	).ThenFunc(func(ctx context.Context, r http.ResponseWriter, w *http.Request) {
		if called == 2 {
			called = 3
		}
		vars, err := cc.FetchRouterParamsFromContext(ctx, "foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", vars["foo"])
	})(nil, nil, params)
	assert.Equal(t, 3, called)
}
