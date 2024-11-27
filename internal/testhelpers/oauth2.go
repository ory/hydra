// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/ory/fosite/token/jwt"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"

	"github.com/ory/x/httpx"
	"github.com/ory/x/ioutilx"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
)

func NewIDToken(t *testing.T, reg driver.Registry, subject string) string {
	return NewIDTokenWithExpiry(t, reg, subject, time.Hour)
}

func NewIDTokenWithExpiry(t *testing.T, reg driver.Registry, subject string, exp time.Duration) string {
	token, _, err := reg.OpenIDJWTStrategy().Generate(context.Background(), jwt.IDTokenClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(exp),
		IssuedAt:  time.Now(),
	}.ToMapClaims(), jwt.NewHeaders())
	require.NoError(t, err)
	return token
}

func NewIDTokenWithClaims(t *testing.T, reg driver.Registry, claims jwt.MapClaims) string {
	token, _, err := reg.OpenIDJWTStrategy().Generate(context.Background(), claims, jwt.NewHeaders())
	require.NoError(t, err)
	return token
}

func NewOAuth2Server(ctx context.Context, t testing.TB, reg driver.Registry) (publicTS, adminTS *httptest.Server) {
	// Lifespan is two seconds to avoid time synchronization issues with SQL.
	reg.Config().MustSet(ctx, config.KeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
	reg.Config().MustSet(ctx, config.KeyAccessTokenLifespan, time.Second*2)
	reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, time.Second*3)
	reg.Config().MustSet(ctx, config.PublicInterface.Key(config.KeySuffixTLSEnabled), false)
	reg.Config().MustSet(ctx, config.AdminInterface.Key(config.KeySuffixTLSEnabled), false)
	reg.Config().MustSet(ctx, config.KeyScopeStrategy, "exact")

	public, admin := x.NewRouterPublic(), x.NewRouterAdmin(reg.Config().AdminURL)

	MustEnsureRegistryKeys(ctx, reg, x.OpenIDConnectKeyName)
	MustEnsureRegistryKeys(ctx, reg, x.OAuth2JWTKeyName)

	reg.RegisterRoutes(ctx, admin, public)

	publicTS = httptest.NewServer(otelhttp.NewHandler(public, "public", otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
		return r.URL.Path
	})))
	t.Cleanup(publicTS.Close)

	adminTS = httptest.NewServer(otelhttp.NewHandler(admin, "admin", otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
		return r.URL.Path
	})))
	t.Cleanup(adminTS.Close)

	reg.Config().MustSet(ctx, config.KeyIssuerURL, publicTS.URL)
	return publicTS, adminTS
}

func DecodeIDToken(t *testing.T, token *oauth2.Token) gjson.Result {
	idt, ok := token.Extra("id_token").(string)
	require.True(t, ok)
	assert.NotEmpty(t, idt)

	body, err := x.DecodeSegment(strings.Split(idt, ".")[1])
	require.NoError(t, err)

	return gjson.ParseBytes(body)
}

func IntrospectToken(t testing.TB, conf *oauth2.Config, token string, adminTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token)

	req := httpx.MustNewRequest("POST", adminTS.URL+"/admin/oauth2/introspect",
		strings.NewReader((url.Values{"token": {token}}).Encode()),
		"application/x-www-form-urlencoded")

	req.SetBasicAuth(conf.ClientID, conf.ClientSecret)
	res, err := adminTS.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body))
}

func RevokeToken(t testing.TB, conf *oauth2.Config, token string, publicTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token)

	req := httpx.MustNewRequest("POST", publicTS.URL+"/oauth2/revoke",
		strings.NewReader((url.Values{"token": {token}}).Encode()),
		"application/x-www-form-urlencoded")

	req.SetBasicAuth(conf.ClientID, conf.ClientSecret)
	res, err := publicTS.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body))
}

func UpdateClientTokenLifespans(t *testing.T, conf *oauth2.Config, clientID string, lifespans client.Lifespans, adminTS *httptest.Server) {
	b, err := json.Marshal(lifespans)
	require.NoError(t, err)
	req := httpx.MustNewRequest(
		"PUT",
		adminTS.URL+"/admin"+client.ClientsHandlerPath+"/"+clientID+"/lifespans",
		bytes.NewBuffer(b),
		"application/json",
	)
	req.SetBasicAuth(conf.ClientID, conf.ClientSecret)
	res, err := adminTS.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)
}

func Userinfo(t *testing.T, token *oauth2.Token, publicTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token.AccessToken)

	req := httpx.MustNewRequest("GET", publicTS.URL+"/userinfo", nil, "")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	res, err := publicTS.Client().Do(req)
	require.NoError(t, err)

	defer res.Body.Close()
	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body))
}

func HTTPServerNotImplementedHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func HTTPServerNoExpectedCallHandler(t testing.TB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("This should not have been called")
	}
}

func NewLoginConsentUI(t testing.TB, c *config.DefaultProvider, login, consent http.HandlerFunc) {
	if login == nil {
		login = HTTPServerNotImplementedHandler
	}

	if consent == nil {
		login = HTTPServerNotImplementedHandler
	}

	lt := httptest.NewServer(login)
	ct := httptest.NewServer(consent)

	t.Cleanup(lt.Close)
	t.Cleanup(ct.Close)

	c.MustSet(context.Background(), config.KeyLoginURL, lt.URL)
	c.MustSet(context.Background(), config.KeyConsentURL, ct.URL)
}

func NewCallbackURL(t testing.TB, prefix string, h http.HandlerFunc) string {
	if h == nil {
		h = HTTPServerNotImplementedHandler
	}

	r := httprouter.New()
	r.GET("/"+prefix, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h(w, r)
	})
	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	return ts.URL + "/" + prefix
}

func NewEmptyCookieJar(t testing.TB) *cookiejar.Jar {
	c, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)
	return c
}

func NewEmptyJarClient(t testing.TB) *http.Client {
	return &http.Client{
		Jar:       NewEmptyCookieJar(t),
		Transport: &loggingTransport{t},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//t.Logf("Redirect to %s", req.URL.String())

			if len(via) >= 20 {
				for k, v := range via {
					t.Logf("Failed with redirect (%d): %s", k, v.URL.String())
				}
				return errors.New("stopped after 20 redirects")
			}
			return nil
		},
	}
}

type loggingTransport struct{ t testing.TB }

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	//s.t.Logf("%s %s", r.Method, r.URL.String())
	//s.t.Logf("%s %s\nWith Cookies: %v", r.Method, r.URL.String(), r.Cookies())

	return otelhttp.DefaultClient.Transport.RoundTrip(r)
}
