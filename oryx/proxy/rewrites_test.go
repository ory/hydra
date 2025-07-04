// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// This test is a unit test for all the rewrite functions,
// including **all** edge cases. It should not go through the network
// and reverse proxy, but just test all helper functions.

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}

func TestRewrites(t *testing.T) {
	t.Run("suite=HeaderRequest", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com/foo/bar", nil)
		require.NoError(t, err)

		c := &HostConfig{
			CookieDomain:   "example.com",
			originalHost:   "example.com",
			UpstreamHost:   "some-project-1234.oryapis.com",
			UpstreamScheme: "https",
			PathPrefix:     "/foo",
		}

		headerRequestRewrite(req, c)
		assert.Equal(t, c.UpstreamScheme, req.URL.Scheme)
		assert.Equal(t, c.UpstreamHost, req.URL.Host)
		assert.Equal(t, "/bar", req.URL.Path)
	})

	t.Run("suite=HTTPS override", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://example.com/foo/bar", nil)
		require.NoError(t, err)

		c := &HostConfig{}
		c.setScheme(&httputil.ProxyRequest{In: req, Out: &http.Request{}})
		assert.Equal(t, "http", c.originalScheme)

		c.ForceOriginalSchemeHTTPS = true
		c.setScheme(&httputil.ProxyRequest{In: req, Out: &http.Request{}})
		assert.Equal(t, "https", c.originalScheme)
	})

	t.Run("suit=HeaderResponse", func(t *testing.T) {
		newOKResp := func(cookie, location string) *http.Response {
			header := http.Header{}
			if cookie != "" {
				header.Add("Set-Cookie", cookie)
			}
			if location != "" {
				header.Add("Location", location)
			}
			return &http.Response{
				Status:        "ok",
				StatusCode:    200,
				Proto:         "https",
				Header:        header,
				Body:          nil,
				ContentLength: 0,
			}
		}

		t.Run("case=replace location and cookie", func(t *testing.T) {
			upstreamHost := "some-project-1234.oryapis.com"

			c := &HostConfig{
				CookieDomain:   "example.com",
				TargetHost:     upstreamHost,
				UpstreamHost:   upstreamHost,
				PathPrefix:     "/foo",
				UpstreamScheme: "https",
				originalHost:   "example.com",
				originalScheme: "http",
			}
			cookie := http.Cookie{
				Name:   "cookie.example",
				Value:  "1234",
				Domain: upstreamHost,
			}
			location := url.URL{
				Scheme: "https",
				Host:   upstreamHost,
				Path:   "/bar",
			}

			resp := newOKResp(cookie.String(), location.String())

			require.NoError(t, headerResponseRewrite(resp, c))

			loc, err := resp.Location()
			require.NoError(t, err)

			assert.Equal(t, c.originalHost, loc.Host)
			assert.Equal(t, c.originalScheme, loc.Scheme)
			assert.Equal(t, "/foo/bar", loc.Path)

			for _, co := range resp.Cookies() {
				assert.Equal(t, c.CookieDomain, co.Domain)
			}
		})

		t.Run("case=replace location and cookie with different target", func(t *testing.T) {
			c := &HostConfig{
				CookieDomain:   "example.com",
				TargetHost:     "foo.bar",
				UpstreamHost:   "next.hop.com",
				PathPrefix:     "/foo",
				UpstreamScheme: "https",
				originalHost:   "example.com",
				originalScheme: "http",
			}
			cookie := http.Cookie{
				Name:   "cookie.example",
				Value:  "1234",
				Domain: c.TargetHost,
			}
			location := url.URL{
				Scheme: "https",
				Host:   c.TargetHost,
				Path:   "/bar",
			}

			resp := newOKResp(cookie.String(), location.String())

			require.NoError(t, headerResponseRewrite(resp, c))

			loc, err := resp.Location()
			require.NoError(t, err)

			assert.Equal(t, c.originalHost, loc.Host)
			assert.Equal(t, c.originalScheme, loc.Scheme)
			assert.Equal(t, "/foo/bar", loc.Path)

			for _, co := range resp.Cookies() {
				assert.Equal(t, c.CookieDomain, co.Domain)
				assert.Equal(t, false, co.Secure)
				assert.Equal(t, http.SameSiteLaxMode, co.SameSite)
			}
		})

		t.Run("case=replace cookie", func(t *testing.T) {
			upstreamHost := "some-project-1234.oryapis.com"

			c := &HostConfig{
				CookieDomain:   "example.com",
				TargetHost:     upstreamHost,
				UpstreamHost:   upstreamHost,
				PathPrefix:     "/foo",
				UpstreamScheme: "https",
				originalHost:   "example.com",
				originalScheme: "http",
			}

			cookie := http.Cookie{
				Name:   "cookie.example",
				Value:  "1234",
				Domain: upstreamHost,
			}

			resp := newOKResp(cookie.String(), "")

			err := headerResponseRewrite(resp, c)
			require.NoError(t, err)

			_, err = resp.Location()
			require.Error(t, err)

			for _, co := range resp.Cookies() {
				assert.Equal(t, c.CookieDomain, co.Domain)
			}
		})

		t.Run("case=no replaced header fields", func(t *testing.T) {
			upstreamHost := "some-project-1234.oryapis.com"

			c := &HostConfig{
				CookieDomain:   "example.com",
				UpstreamHost:   upstreamHost,
				PathPrefix:     "/foo",
				UpstreamScheme: "https",
				originalHost:   "example.com",
				originalScheme: "http",
			}

			resp := newOKResp("", "")

			require.NoError(t, headerResponseRewrite(resp, c))

			assert.Len(t, resp.Cookies(), 0)
			_, err := resp.Location()
			assert.Error(t, http.ErrNoLocation, err)
		})

	})

	t.Run("suit=BodyResponse", func(t *testing.T) {
		newOKResp := func(body string) *http.Response {
			return &http.Response{
				Status:        "OK",
				StatusCode:    200,
				Proto:         "http",
				Body:          io.NopCloser(strings.NewReader(body)),
				ContentLength: int64(len([]byte(body))),
			}
		}

		t.Run("case=empty body", func(t *testing.T) {
			resp := newOKResp("")
			// we actually want to see if it also handles nil bodies
			resp.Body = nil

			_, _, err := bodyResponseRewrite(resp, &HostConfig{})
			assert.NoError(t, err)
		})

		t.Run("case=json body with path prefix and method rewrite", func(t *testing.T) {
			upstreamHost := "some-project-1234.oryapis.com"

			c := &HostConfig{
				CookieDomain:   "example.com",
				TargetHost:     upstreamHost,
				TargetScheme:   "https",
				UpstreamHost:   upstreamHost,
				UpstreamScheme: "https",
				PathPrefix:     "/foo",
				originalHost:   "auth.example.com",
				originalScheme: "http",
			}

			body, err := sjson.Set("{}", "some_key", "https://"+upstreamHost+"/path")
			require.NoError(t, err)
			body, err = sjson.Set(body, "inner_resp_arr.0.inner_key", "https://"+upstreamHost+"/bar")
			require.NoError(t, err)
			body, err = sjson.Set(body, "inner_resp.inner_key", "https://"+upstreamHost)
			require.NoError(t, err)

			resp := newOKResp(body)

			b, _, err := bodyResponseRewrite(resp, c)
			require.NoError(t, err)

			assert.Equal(t, "http://auth.example.com/foo", gjson.GetBytes(b, "inner_resp.inner_key").Str, "%s", b)
			assert.Equal(t, "http://auth.example.com/foo/path", gjson.GetBytes(b, "some_key").Str, "%s", b)
			assert.Equal(t, "http://auth.example.com/foo/bar", gjson.GetBytes(b, "inner_resp_arr.0.inner_key").Str, "%s", b)
		})

		t.Run("case=string body and no path prefix", func(t *testing.T) {
			c := &HostConfig{
				CookieDomain:   "example.com",
				TargetHost:     "some-project-1234.oryapis.com",
				TargetScheme:   "https",
				UpstreamHost:   "some-project-1234.oryapis.com",
				UpstreamScheme: "https",
				PathPrefix:     "/foo",
				originalHost:   "auth.example.com",
				originalScheme: "https",
			}

			resp := newOKResp(fmt.Sprintf("this is a string body %s://%s", c.TargetScheme, c.TargetHost))

			replaced, _, err := bodyResponseRewrite(resp, c)
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("this is a string body %s://%s", c.originalScheme, c.originalHost+c.PathPrefix), string(replaced))
		})

		t.Run("case=different target and upstream hosts", func(t *testing.T) {
			c := &HostConfig{
				CookieDomain:   "example.com",
				TargetHost:     "actually.host.com",
				TargetScheme:   "https",
				UpstreamHost:   "some-project-1234.oryapis.com",
				UpstreamScheme: "https",
				PathPrefix:     "/foo",
				originalHost:   "auth.example.com",
				originalScheme: "http",
			}

			resp := newOKResp(fmt.Sprintf("I am available at %s://%s", c.TargetScheme, c.TargetHost))

			replaced, _, err := bodyResponseRewrite(resp, c)
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("I am available at %s://%s", c.originalScheme, c.originalHost+c.PathPrefix), string(replaced))
		})
	})
}

func TestHelpers(t *testing.T) {
	t.Run("func=stripPort", func(t *testing.T) {
		for input, output := range map[string]string{
			"example.com":      "example.com",
			"example.com:4321": "example.com",
			"192.168.0.0":      "192.168.0.0",
			"192.168.0.0:8080": "192.168.0.0",
		} {
			assert.Equal(t, output, stripPort(input))
		}
	})

	t.Run("func=readBody", func(t *testing.T) {
		t.Run("case=basic body", func(t *testing.T) {
			rawBody, writer, err := readBody(http.Header{}, io.NopCloser(bytes.NewBufferString("simple body")))
			require.NoError(t, err)
			assert.Equal(t, "simple body", string(rawBody))

			_, err = writer.Write([]byte("not compressed"))
			require.NoError(t, err)
			assert.Equal(t, "not compressed", writer.buf.String())
		})

		t.Run("case=gziped body", func(t *testing.T) {
			header := http.Header{}
			header.Set("Content-Encoding", "gzip")
			body := &bytes.Buffer{}
			w := gzip.NewWriter(body)
			_, err := w.Write([]byte("this is compressed"))
			require.NoError(t, err)
			require.NoError(t, w.Close())

			rawBody, writer, err := readBody(header, io.NopCloser(body))
			require.NoError(t, err)
			assert.Equal(t, "this is compressed", string(rawBody))

			_, err = writer.Write([]byte("should compress"))
			require.NoError(t, err)
			assert.NotEqual(t, "should compress", writer.buf.String())

			r, err := gzip.NewReader(&writer.buf)
			require.NoError(t, err)
			content, err := io.ReadAll(r)
			require.NoError(t, err)
			assert.Equal(t, "should compress", string(content))
		})
	})

	t.Run("func=compressableBody.Read", func(t *testing.T) {
		t.Run("case=empty body", func(t *testing.T) {
			n, err := (*compressableBody)(nil).Read(make([]byte, 10))
			assert.True(t, err == io.EOF)
			assert.Equal(t, 0, n)
		})

		t.Run("case=has content", func(t *testing.T) {
			content := "some test content, who cares"
			b := make([]byte, 128)
			n, err := (&compressableBody{
				buf: *bytes.NewBufferString(content),
			}).Read(b)
			require.NoError(t, err)
			assert.Equal(t, content, string(b[:n]))
		})
	})

	t.Run("func=compressableBody.Write", func(t *testing.T) {
		t.Run("case=empty body", func(t *testing.T) {
			n, err := (*compressableBody)(nil).Write([]byte{0, 1, 2, 3})
			assert.NoError(t, err)
			assert.Equal(t, 0, n)
		})

		t.Run("case=no writer", func(t *testing.T) {
			b := &compressableBody{}
			_, err := b.Write([]byte("foo bar"))
			require.NoError(t, err)
			assert.Equal(t, "foo bar", b.buf.String())
		})

		t.Run("case=wrapped writer", func(t *testing.T) {
			other := &bytes.Buffer{}
			b := &compressableBody{}
			b.w = nopWriteCloser{io.MultiWriter(other, &b.buf)}
			_, err := b.Write([]byte("foo bar"))
			require.NoError(t, err)
			assert.Equal(t, "foo bar", b.buf.String())
			assert.Equal(t, "foo bar", other.String())
		})
	})
}
