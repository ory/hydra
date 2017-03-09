package fosite

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLocalhost(t *testing.T) {
	for k, c := range []struct {
		expect bool
		rawurl string
	}{
		{expect: false, rawurl: "https://foo.bar"},
		{expect: true, rawurl: "https://localhost"},
		{expect: true, rawurl: "https://localhost:1234"},
		{expect: true, rawurl: "https://127.0.0.1:1234"},
		{expect: true, rawurl: "https://127.0.0.1"},
	} {
		u, _ := url.Parse(c.rawurl)
		assert.Equal(t, c.expect, isLocalhost(u), "case %d", k)
	}
}

// Test for
// * https://tools.ietf.org/html/rfc6749#section-3.1.2
//   The endpoint URI MAY include an
//   "application/x-www-form-urlencoded" formatted (per Appendix B) query
//   component ([RFC3986] Section 3.4), which MUST be retained when adding
//   additional query parameters.
func TestGetRedirectURI(t *testing.T) {
	for k, c := range []struct {
		in       string
		isError  bool
		expected string
	}{
		{in: "", isError: false, expected: ""},
		{in: "https://google.com/", isError: false, expected: "https://google.com/"},
		{in: "https://google.com/?foo=bar%20foo+baz", isError: false, expected: "https://google.com/?foo=bar foo baz"},
	} {
		values := url.Values{}
		values.Set("redirect_uri", c.in)
		res, err := GetRedirectURIFromRequestValues(values)
		assert.Equal(t, c.isError, err != nil, "%s", err)
		if err == nil {
			assert.Equal(t, c.expected, res)
		}
		t.Logf("Passed test case %d", k)
	}
}

// rfc6749 10.6.
// Authorization Code Redirection URI Manipulation
// The authorization server	MUST require public clients and SHOULD require confidential clients
// to register their redirection URIs.  If a redirection URI is provided
// in the request, the authorization server MUST validate it against the
// registered value.
//
// rfc6819 4.4.1.7.
// Threat: Authorization "code" Leakage through Counterfeit Client
// The authorization server may also enforce the usage and validation
// of pre-registered redirect URIs (see Section 5.2.3.5).
func TestDoesClientWhiteListRedirect(t *testing.T) {
	for k, c := range []struct {
		client   Client
		url      string
		isError  bool
		expected string
	}{
		{
			client:  &DefaultClient{RedirectURIs: []string{""}},
			url:     "https://foo.com/cb",
			isError: true,
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"wta://auth"}},
			url:      "wta://auth",
			expected: "wta://auth",
			isError:  false,
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"wta:///auth"}},
			url:      "wta:///auth",
			expected: "wta:///auth",
			isError:  false,
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"wta://foo/auth"}},
			url:      "wta://foo/auth",
			expected: "wta://foo/auth",
			isError:  false,
		},
		{
			client:  &DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:     "https://foo.com/cb",
			isError: true,
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:      "",
			isError:  false,
			expected: "https://bar.com/cb",
		},
		{
			client:  &DefaultClient{RedirectURIs: []string{""}},
			url:     "",
			isError: true,
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:      "https://bar.com/cb",
			isError:  false,
			expected: "https://bar.com/cb",
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"https://bar.Com/cb"}},
			url:      "https://bar.com/cb",
			isError:  false,
			expected: "https://bar.com/cb",
		},
		{
			client:   &DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:      "https://bar.Com/cb",
			isError:  false,
			expected: "https://bar.Com/cb",
		},
		{
			client:  &DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:     "https://bar.com/cb123",
			isError: true,
		},
	} {
		redir, err := MatchRedirectURIWithClientRedirectURIs(c.url, c.client)
		assert.Equal(t, c.isError, err != nil, "%d: %s", k, err)
		if err == nil {
			require.NotNil(t, redir, "%d", k)
			assert.Equal(t, c.expected, redir.String(), "%d", k)
		}
		t.Logf("Passed test case %d", k)
	}
}

func TestIsRedirectURISecure(t *testing.T) {
	for d, c := range []struct {
		u   string
		err bool
	}{
		{u: "http://google.com", err: true},
		{u: "https://google.com", err: false},
		{u: "http://localhost", err: false},
		{u: "wta://auth", err: false},
	} {
		uu, err := url.Parse(c.u)
		require.Nil(t, err)
		assert.Equal(t, !c.err, IsRedirectURISecure(uu), "case %d", d)
	}
}
