package middleware

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	for k, c := range []struct {
		expectedIP    string
		expectedOwner string
		expectedUA    string
		request       *http.Request
		owner         string
		headers       map[string]string
	}{
		{"1234", "peter", "firefox", &http.Request{RemoteAddr:"1234"}, "peter", map[string]string{"user-Agent": "firefox"}},
		{"1234", "peter", "", &http.Request{RemoteAddr:"4321"}, "peter", map[string]string{"x-FoRwarded-FoR": "1234"}},
		{"1234", "peter", "firefox", &http.Request{RemoteAddr:"4321"}, "peter", map[string]string{"x-FoRwarded-FoR": "1234", "user-Agent": "firefox"}},
	} {
		h := http.Header{}
		for y, x := range c.headers {
			h.Set(y, x)
		}
		r := c.request
		r.Header = h
		t.Logf("%v", r)
		e := NewEnv(r)
		e.Owner(c.owner)
		assert.Equal(t, c.expectedIP, e.Ctx().ClientIP, "Case %d", k)
		assert.Equal(t, c.expectedOwner, e.Ctx().Owner, "Case %d", k)
		assert.Equal(t, c.expectedUA, e.Ctx().UserAgent, "Case %d", k)
	}
}
