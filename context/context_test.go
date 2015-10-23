package context

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

var called = 0

func middleware(h ContextHandler, t *testing.T) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		called++
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func handler(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	called++
}

func TestContextAdapter(t *testing.T) {
	h := &ContextAdapter{
		Ctx:     context.Background(),
		Handler: middleware(ContextHandlerFunc(handler), t),
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/handler", nil)
	require.Nil(t, err)

	m := mux.NewRouter()
	m.Handle("/handler", h).Methods("GET")
	m.ServeHTTP(recorder, req)

	assert.Equal(t, 2, called)
}
