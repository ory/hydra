package handler_test

import (
	. "github.com/ory-am/common/handler"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"testing"
)

func TestNewContextAdapter(t *testing.T) {
	called := 0
	NewContextAdapter(
		context.Background(),
		func(next ContextHandler) ContextHandler {
			return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				if called == 0 {
					called = 1
				}
				next.ServeHTTPContext(ctx, rw, req)
			})
		},
		func(next ContextHandler) ContextHandler {
			return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
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
	}).ServeHTTP(nil, nil)
	assert.Equal(t, 3, called)
}
