package testhelpers

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/httptest"
)

func FlexibleServer(t *testing.T, h *http.HandlerFunc) string {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(*h)(w, r)
	}))
	t.Cleanup(ts.Close)
	return ts.URL
}
