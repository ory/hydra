// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	djwt "github.com/ory/fosite/token/jwt"

	"github.com/ory/fosite/token/jwt"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"

	"github.com/ory/x/httpx"
	"github.com/ory/x/ioutilx"

	"net/http/httptest"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
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

func NewIDTokenWithClaims(t *testing.T, reg driver.Registry, claims djwt.MapClaims) string {
	token, _, err := reg.OpenIDJWTStrategy().Generate(context.Background(), claims, jwt.NewHeaders())
	require.NoError(t, err)
	return token
}

func NewOAuth2Server(ctx context.Context, t *testing.T, reg driver.Registry) (publicTS, adminTS *httptest.Server) {
	// Lifespan is two seconds to avoid time synchronization issues with SQL.
	reg.Config().MustSet(ctx, config.KeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
	reg.Config().MustSet(ctx, config.KeyAccessTokenLifespan, time.Second*2)
	reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, time.Second*3)
	reg.Config().MustSet(ctx, config.PublicInterface.Key(config.KeySuffixTLSEnabled), false)
	reg.Config().MustSet(ctx, config.AdminInterface.Key(config.KeySuffixTLSEnabled), false)
	reg.Config().MustSet(ctx, config.KeyScopeStrategy, "exact")

	public, admin := x.NewRouterPublic(), x.NewRouterAdmin(reg.Config().AdminURL)

	publicTS = httptest.NewServer(public)
	t.Cleanup(publicTS.Close)

	adminTS = httptest.NewServer(admin)
	t.Cleanup(adminTS.Close)

	reg.Config().MustSet(ctx, config.KeyIssuerURL, publicTS.URL)
	// SendDebugMessagesToClients: true,

	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)
	internal.MustEnsureRegistryKeys(reg, x.OAuth2JWTKeyName)

	reg.RegisterRoutes(ctx, admin, public)
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

func IntrospectToken(t *testing.T, conf *oauth2.Config, token string, adminTS *httptest.Server) gjson.Result {
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

func HTTPServerNoExpectedCallHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("This should not have been called")
	}
}

func NewLoginConsentUI(t *testing.T, c *config.DefaultProvider, login, consent http.HandlerFunc) {
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

func NewCallbackURL(t *testing.T, prefix string, h http.HandlerFunc) string {
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

func NewEmptyCookieJar(t *testing.T) *cookiejar.Jar {
	c, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)
	return c
}

func NewEmptyJarClient(t *testing.T) *http.Client {
	return &http.Client{
		Jar: NewEmptyCookieJar(t),
	}
}
