// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2cors_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/hydra/v2/x/oauth2cors"
	"github.com/ory/x/configx"
	"github.com/ory/x/dbal"
)

func TestOAuth2AwareCORSMiddleware(t *testing.T) {
	ctx := context.Background()
	dsn := dbal.NewSQLiteInMemoryDatabase(t.Name())
	r := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue("dsn", dsn)))
	token, signature, _ := r.OAuth2HMACStrategy().GenerateAccessToken(ctx, nil)

	for k, tc := range []struct {
		prep         func(*testing.T, *driver.RegistrySQL)
		configs      map[string]any
		d            string
		mw           func(http.Handler) http.Handler
		code         int
		header       http.Header
		expectHeader http.Header
		method       string
		body         io.Reader
	}{
		{
			d:            "should ignore when disabled",
			prep:         func(t *testing.T, r *driver.RegistrySQL) {},
			code:         http.StatusNotImplemented,
			header:       http.Header{},
			expectHeader: http.Header{},
		},
		{
			d: "should reject when basic auth but client does not exist and cors enabled",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo", "bar"))}},
			expectHeader: http.Header{"Vary": {"Origin"}},
		},
		{
			d: "should reject when post auth client exists but origin not allowed",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-2", Secret: "bar", AllowedCORSOrigins: []string{"http://not-foobar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Content-Type": {"application/x-www-form-urlencoded"}},
			expectHeader: http.Header{"Vary": {"Origin"}},
			method:       http.MethodPost,
			body:         bytes.NewBufferString(url.Values{"client_id": []string{"foo-2"}}.Encode()),
		},
		{
			d: "should accept when post auth client exists and origin allowed",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-3", Secret: "bar", AllowedCORSOrigins: []string{"http://foobar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Content-Type": {"application/x-www-form-urlencoded"}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
			method:       http.MethodPost,
			body:         bytes.NewBufferString(url.Values{"client_id": {"foo-3"}}.Encode()),
		},
		{
			d: "should reject when basic auth client exists but origin not allowed",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-2", Secret: "bar", AllowedCORSOrigins: []string{"http://not-foobar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-2", "bar"))}},
			expectHeader: http.Header{"Vary": {"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and origin allowed",
			configs: map[string]any{
				"serve.public.cors.enabled": true,
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-3", Secret: "bar", AllowedCORSOrigins: []string{"http://foobar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-3", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and origin allowed",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-3", Secret: "bar", AllowedCORSOrigins: []string{"http://foobar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-3", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and origin (with partial wildcard) is allowed per client",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-4", Secret: "bar", AllowedCORSOrigins: []string{"http://*.foobar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foo.foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-4", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foo.foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and wildcard origin is allowed per client",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-4", Secret: "bar", AllowedCORSOrigins: []string{"http://*"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foo.foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-4", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foo.foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and origin (with full wildcard) is allowed globally",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"*"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-5", Secret: "bar", AllowedCORSOrigins: []string{"http://barbar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"*"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-5", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"*"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and origin (with partial wildcard) is allowed globally",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://*.foobar.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-6", Secret: "bar", AllowedCORSOrigins: []string{"http://barbar.com"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foo.foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-6", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foo.foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept when basic auth client exists and origin (with full wildcard) allowed per client",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-7", Secret: "bar", AllowedCORSOrigins: []string{"*"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-7", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should succeed on pre-flight request when token introspection fails",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {"Bearer 1234"}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
			method:       "OPTIONS",
		},
		{
			d: "should fail when token introspection fails",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {"Bearer 1234"}},
			expectHeader: http.Header{"Vary": {"Origin"}},
		},
		{
			d: "should work when token introspection returns a session",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://not-test-domain.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				sess := oauth2.NewTestSession(t, "foo-9")
				sess.SetExpiresAt(fosite.AccessToken, time.Now().Add(time.Hour))
				ar := fosite.NewAccessRequest(sess)
				cl := &client.Client{ID: "foo-9", Secret: "bar", AllowedCORSOrigins: []string{"http://foobar.com"}}
				ar.Client = cl

				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, cl)
				_ = r.OAuth2Storage().CreateAccessTokenSession(ctx, signature, ar)
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foobar.com"}, "Authorization": {"Bearer " + token}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept any allowed specified origin protocol",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://*", "https://*"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-11", Secret: "bar", AllowedCORSOrigins: []string{"*"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://foo.foobar.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-11", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://foo.foobar.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept client origin when basic auth client exists and origin is set at the client as well as the server",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://**.example.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-12", Secret: "bar", AllowedCORSOrigins: []string{"http://myapp.example.biz"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://myapp.example.biz"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-12", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://myapp.example.biz"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
		{
			d: "should accept server origin when basic auth client exists and origin is set at the client as well as the server",
			configs: map[string]any{
				"serve.public.cors.enabled":         true,
				"serve.public.cors.allowed_origins": []string{"http://**.example.com"},
			},
			prep: func(t *testing.T, r *driver.RegistrySQL) {
				// Ignore unique violations
				_ = r.ClientManager().CreateClient(ctx, &client.Client{ID: "foo-13", Secret: "bar", AllowedCORSOrigins: []string{"http://myapp.example.biz"}})
			},
			code:         http.StatusNotImplemented,
			header:       http.Header{"Origin": {"http://client-app.example.com"}, "Authorization": {fmt.Sprintf("Basic %s", x.BasicAuth("foo-13", "bar"))}},
			expectHeader: http.Header{"Access-Control-Allow-Credentials": []string{"true"}, "Access-Control-Allow-Origin": []string{"http://client-app.example.com"}, "Access-Control-Expose-Headers": []string{"Cache-Control, Expires, Last-Modified, Pragma, Content-Length, Content-Language, Content-Type"}, "Vary": []string{"Origin"}},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			r := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue("dsn", dsn), configx.WithValues(tc.configs)))

			if tc.prep != nil {
				tc.prep(t, r)
			}

			method := "GET"
			if tc.method != "" {
				method = tc.method
			}
			req, err := http.NewRequest(method, "http://foobar.com/", tc.body)
			require.NoError(t, err)
			for k := range tc.header {
				req.Header.Set(k, tc.header.Get(k))
			}

			res := httptest.NewRecorder()
			oauth2cors.Middleware(r)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
			})).ServeHTTP(res, req)
			require.NoError(t, err)
			assert.EqualValues(t, tc.code, res.Code)
			assert.EqualValues(t, tc.expectHeader, res.Header())
		})
	}
}
