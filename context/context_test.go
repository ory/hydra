package context

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/stretchr/testify/require"
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

var called = 0

func middleware(h ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		called++
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func handler(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	called++
}

func TestContextAdapterThenExports(t *testing.T) {
	h := NewContextAdapter(context.Background(), middleware).Then(ContextHandler(ContextHandlerFunc(handler)))

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/handler", nil)
	require.Nil(t, err)

	m := mux.NewRouter()
	m.Handle("/handler", h).Methods("GET")
	m.ServeHTTP(recorder, req)

	assert.Equal(t, 2, called)
	called = 0
}

func TestContextAdapterThenFuncExports(t *testing.T) {
	h := NewContextAdapter(context.Background(), middleware).ThenFunc(handler)

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/handler", nil)
	require.Nil(t, err)

	m := mux.NewRouter()
	m.Handle("/handler", h).Methods("GET")
	m.ServeHTTP(recorder, req)

	assert.Equal(t, 2, called)
	called = 0
}

func TestContextAdapter(t *testing.T) {
	h := &contextAdapter{
		ctx:   context.Background(),
		final: middleware(ContextHandlerFunc(handler)),
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/handler", nil)
	require.Nil(t, err)

	m := mux.NewRouter()
	m.Handle("/handler", h).Methods("GET")
	m.ServeHTTP(recorder, req)

	assert.Equal(t, 2, called)
	called = 0
}
