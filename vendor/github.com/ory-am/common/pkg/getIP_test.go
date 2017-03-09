package pkg_test

import (
	"github.com/ory-am/common/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetIP(t *testing.T) {
	for k, c := range []struct {
		expectedIP string
		request    *http.Request
		headers    map[string]string
	}{
		{"1234", &http.Request{RemoteAddr: "1234"}, map[string]string{"user-Agent": "firefox"}},
		{"1234", &http.Request{RemoteAddr: "4321"}, map[string]string{"x-FoRwarded-FoR": "1234"}},
		{"1234", &http.Request{RemoteAddr: "4321"}, map[string]string{"x-FoRwarded-FoR": "1234", "user-Agent": "firefox"}},
	} {
		h := http.Header{}
		for y, x := range c.headers {
			h.Set(y, x)
		}
		r := c.request
		r.Header = h
		assert.Equal(t, c.expectedIP, pkg.GetIP(r), "Case %d", k)
	}
}
