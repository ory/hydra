// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"

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
		{expect: true, rawurl: "https://test.localhost:1234"},
		{expect: true, rawurl: "https://test.localhost"},
	} {
		u, _ := url.Parse(c.rawurl)
		assert.Equal(t, c.expect, fosite.IsLocalhost(u), "case %d", k)
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
		client   fosite.Client
		url      string
		isError  bool
		expected string
	}{
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{""}},
			url:     "https://foo.com/cb",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"wta://auth"}},
			url:      "wta://auth",
			expected: "wta://auth",
			isError:  false,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"wta:///auth"}},
			url:      "wta:///auth",
			expected: "wta:///auth",
			isError:  false,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"wta://foo/auth"}},
			url:      "wta://foo/auth",
			expected: "wta://foo/auth",
			isError:  false,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:     "https://foo.com/cb",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:      "",
			isError:  false,
			expected: "https://bar.com/cb",
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{""}},
			url:     "",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:      "https://bar.com/cb",
			isError:  false,
			expected: "https://bar.com/cb",
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"https://bar.com/cb"}},
			url:     "https://bar.com/cb123",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://[::1]"}},
			url:      "http://[::1]:1024",
			expected: "http://[::1]:1024",
			isError:  false,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://[::1]"}},
			url:     "http://[::1]:1024/cb",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://[::1]/cb"}},
			url:      "http://[::1]:1024/cb",
			expected: "http://[::1]:1024/cb",
			isError:  false,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://[::1]"}},
			url:     "http://foo.bar/bar",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1"}},
			url:      "http://127.0.0.1:1024",
			expected: "http://127.0.0.1:1024",
			isError:  false,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1/cb"}},
			url:      "http://127.0.0.1:64000/cb",
			expected: "http://127.0.0.1:64000/cb",
			isError:  false,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1"}},
			url:     "http://127.0.0.1:64000/cb",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1"}},
			url:      "http://127.0.0.1",
			expected: "http://127.0.0.1",
			isError:  false,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1/Cb"}},
			url:      "http://127.0.0.1:8080/Cb",
			expected: "http://127.0.0.1:8080/Cb",
			isError:  false,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1"}},
			url:     "http://foo.bar/bar",
			isError: true,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1"}},
			url:     ":/invalid.uri)bar",
			isError: true,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1:8080/cb"}},
			url:     "http://127.0.0.1:8080/Cb",
			isError: true,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1:8080/cb"}},
			url:     "http://127.0.0.1:8080/cb?foo=bar",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1:8080/cb?foo=bar"}},
			url:      "http://127.0.0.1:8080/cb?foo=bar",
			expected: "http://127.0.0.1:8080/cb?foo=bar",
			isError:  false,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1:8080/cb?foo=bar"}},
			url:     "http://127.0.0.1:8080/cb?baz=bar&foo=bar",
			isError: true,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1:8080/cb?foo=bar&baz=bar"}},
			url:     "http://127.0.0.1:8080/cb?baz=bar&foo=bar",
			isError: true,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"https://www.ory.sh/cb"}},
			url:     "http://127.0.0.1:8080/cb",
			isError: true,
		},
		{
			client:  &fosite.DefaultClient{RedirectURIs: []string{"http://127.0.0.1:8080/cb"}},
			url:     "https://www.ory.sh/cb",
			isError: true,
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"web+application://callback"}},
			url:      "web+application://callback",
			isError:  false,
			expected: "web+application://callback",
		},
		{
			client:   &fosite.DefaultClient{RedirectURIs: []string{"https://google.com/?foo=bar%20foo+baz"}},
			url:      "https://google.com/?foo=bar%20foo+baz",
			isError:  false,
			expected: "https://google.com/?foo=bar%20foo+baz",
		},
	} {
		redir, err := fosite.MatchRedirectURIWithClientRedirectURIs(c.url, c.client)
		assert.Equal(t, c.isError, err != nil, "%d: %+v", k, c)
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
		{u: "http://test.localhost", err: false},
		{u: "http://127.0.0.1/", err: false},
		{u: "http://[::1]/", err: false},
		{u: "http://127.0.0.1:8080/", err: false},
		{u: "http://[::1]:8080/", err: false},
		{u: "http://testlocalhost", err: true},
		{u: "wta://auth", err: false},
	} {
		uu, err := url.Parse(c.u)
		require.NoError(t, err)
		assert.Equal(t, !c.err, fosite.IsRedirectURISecure(context.Background(), uu), "case %d", d)
	}
}

func TestWriteAuthorizeFormPostResponse(t *testing.T) {
	for d, c := range []struct {
		parameters url.Values
		check      func(code string, state string, customParams url.Values, d int)
	}{
		{
			parameters: url.Values{"code": {"lshr755nsg39fgur"}, "state": {"924659540232"}},
			check: func(code string, state string, customParams url.Values, d int) {
				assert.Equal(t, "lshr755nsg39fgur", code, "case %d", d)
				assert.Equal(t, "924659540232", state, "case %d", d)
			},
		},
		{
			parameters: url.Values{"code": {"lshr75*ns-39f+ur"}, "state": {"9a:* <&)"}},
			check: func(code string, state string, customParams url.Values, d int) {
				assert.Equal(t, "lshr75*ns-39f+ur", code, "case %d", d)
				assert.Equal(t, "9a:* <&)", state, "case %d", d)
			},
		},
		{
			parameters: url.Values{"code": {"1234"}, "custom": {"test2", "test3"}},
			check: func(code string, state string, customParams url.Values, d int) {
				assert.Equal(t, "1234", code, "case %d", d)
				assert.Equal(t, []string{"test2", "test3"}, customParams["custom"], "case %d", d)
			},
		},
		{
			parameters: url.Values{"code": {"1234"}, "custom": {"<b>Bold</b>"}},
			check: func(code string, state string, customParams url.Values, d int) {
				assert.Equal(t, "1234", code, "case %d", d)
				assert.Equal(t, "<b>Bold</b>", customParams.Get("custom"), "case %d", d)
			},
		},
	} {
		var responseBuffer bytes.Buffer
		redirectURL := "https://localhost:8080/cb"
		//parameters :=
		fosite.WriteAuthorizeFormPostResponse(redirectURL, c.parameters, fosite.DefaultFormPostTemplate, &responseBuffer)
		code, state, _, _, customParams, _, err := internal.ParseFormPostResponse(redirectURL, io.NopCloser(bytes.NewReader(responseBuffer.Bytes())))
		assert.NoError(t, err, "case %d", d)
		c.check(code, state, customParams, d)

	}
}

func TestIsRedirectURISecureStrict(t *testing.T) {
	for d, c := range []struct {
		u   string
		err bool
	}{
		{u: "http://google.com", err: true},
		{u: "https://google.com", err: false},
		{u: "http://localhost", err: false},
		{u: "http://test.localhost", err: false},
		{u: "http://127.0.0.1/", err: false},
		{u: "http://[::1]/", err: false},
		{u: "http://127.0.0.1:8080/", err: false},
		{u: "http://[::1]:8080/", err: false},
		{u: "http://testlocalhost", err: true},
		{u: "wta://auth", err: true},
	} {
		uu, err := url.Parse(c.u)
		require.NoError(t, err)
		assert.Equal(t, !c.err, fosite.IsRedirectURISecureStrict(context.Background(), uu), "case %d", d)
	}
}

func TestURLSetFragment(t *testing.T) {
	for d, c := range []struct {
		u string
		a string
		f url.Values
	}{
		{u: "http://google.com", a: "http://google.com#code=567060896", f: url.Values{"code": []string{"567060896"}}},
		{u: "http://google.com", a: "http://google.com#code=567060896&scope=read", f: url.Values{"code": []string{"567060896"}, "scope": []string{"read"}}},
		{u: "http://google.com", a: "http://google.com#code=567060896&scope=read%20mail", f: url.Values{"code": []string{"567060896j"}, "scope": []string{"read mail"}}},
		{u: "http://google.com", a: "http://google.com#code=567060896&scope=read+write", f: url.Values{"code": []string{"567060896"}, "scope": []string{"read+write"}}},
		{u: "http://google.com", a: "http://google.com#code=567060896&scope=api:*", f: url.Values{"code": []string{"567060896"}, "scope": []string{"api:*"}}},
		{u: "https://google.com?foo=bar", a: "https://google.com?foo=bar#code=567060896", f: url.Values{"code": []string{"567060896"}}},
		{u: "http://localhost?foo=bar&baz=foo", a: "http://localhost?foo=bar&baz=foo#code=567060896", f: url.Values{"code": []string{"567060896"}}},
	} {
		uu, err := url.Parse(c.u)
		require.NoError(t, err)
		fosite.URLSetFragment(uu, c.f)
		tURL, err := url.Parse(uu.String())
		require.NoError(t, err)
		r := ParseURLFragment(tURL.Fragment)
		assert.Equal(t, c.f.Get("code"), r.Get("code"), "case %d", d)
		assert.Equal(t, c.f.Get("scope"), r.Get("scope"), "case %d", d)
	}
}
func ParseURLFragment(fragment string) url.Values {
	r := url.Values{}
	kvs := strings.Split(fragment, "&")
	for _, kv := range kvs {
		kva := strings.Split(kv, "=")
		r.Add(kva[0], kva[1])
	}
	return r
}
