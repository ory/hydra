// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"

	"github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/ioutilx"
	"github.com/ory/x/prometheusx"
	"github.com/ory/x/reqlog"
)

func NewIDToken(t *testing.T, reg *driver.RegistrySQL, subject string) string {
	return NewIDTokenWithExpiry(t, reg, subject, time.Hour)
}

func NewIDTokenWithExpiry(t *testing.T, reg *driver.RegistrySQL, subject string, exp time.Duration) string {
	token, _, err := reg.OpenIDJWTStrategy().Generate(context.Background(), jwt.IDTokenClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(exp),
		IssuedAt:  time.Now(),
	}.ToMapClaims(), jwt.NewHeaders())
	require.NoError(t, err)
	return token
}

func NewIDTokenWithClaims(t *testing.T, reg *driver.RegistrySQL, claims jwt.MapClaims) string {
	token, _, err := reg.OpenIDJWTStrategy().Generate(context.Background(), claims, jwt.NewHeaders())
	require.NoError(t, err)
	return token
}

// NewOAuth2Server
// Deprecated: use NewConfigurableOAuth2Server instead
func NewOAuth2Server(ctx context.Context, t testing.TB, reg *driver.RegistrySQL) (publicTS, adminTS *httptest.Server) {
	reg.Config().MustSet(ctx, config.KeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
	reg.Config().MustSet(ctx, config.KeyAccessTokenLifespan, 10*time.Second)
	reg.Config().MustSet(ctx, config.KeyRefreshTokenLifespan, 20*time.Second)
	reg.Config().MustSet(ctx, config.KeyScopeStrategy, "exact")

	return NewConfigurableOAuth2Server(ctx, t, reg)
}

func NewConfigurableOAuth2Server(ctx context.Context, t testing.TB, reg *driver.RegistrySQL) (publicTS, adminTS *httptest.Server) {
	MustEnsureRegistryKeys(t, reg, x.OpenIDConnectKeyName)
	MustEnsureRegistryKeys(t, reg, x.OAuth2JWTKeyName)

	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	{
		n := negroni.New()
		n.Use(reqlog.NewMiddleware())
		n.UseFunc(httprouterx.TrimTrailingSlashNegroni)
		n.UseFunc(httprouterx.NoCacheNegroni)
		n.UseFunc(httprouterx.AddAdminPrefixIfNotPresentNegroni)

		router := x.NewRouterAdmin(metrics)
		reg.RegisterAdminRoutes(router)
		n.UseHandler(router)

		adminTS = httptest.NewServer(n)
		t.Cleanup(adminTS.Close)
		reg.Config().MustSet(ctx, config.KeyAdminURL, adminTS.URL)
	}
	{
		n := negroni.New()
		n.Use(reqlog.NewMiddleware())
		n.UseFunc(httprouterx.TrimTrailingSlashNegroni)
		n.UseFunc(httprouterx.NoCacheNegroni)

		router := x.NewRouterPublic(metrics)
		reg.RegisterPublicRoutes(ctx, router)
		n.UseHandler(router)

		publicTS = httptest.NewServer(n)
		t.Cleanup(publicTS.Close)
		reg.Config().MustSet(ctx, config.KeyAdminURL, publicTS.URL)
	}

	reg.Config().MustSet(ctx, config.KeyIssuerURL, publicTS.URL)
	return publicTS, adminTS
}

func DecodeIDToken(t *testing.T, token *oauth2.Token) gjson.Result {
	idt, ok := token.Extra("id_token").(string)
	require.True(t, ok)
	assert.NotEmpty(t, idt)

	return gjson.ParseBytes(InsecureDecodeJWT(t, idt))
}

func IntrospectToken(t testing.TB, token string, adminTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token)

	req := httpx.MustNewRequest("POST", adminTS.URL+"/admin/oauth2/introspect",
		strings.NewReader((url.Values{"token": {token}}).Encode()),
		"application/x-www-form-urlencoded")

	res, err := adminTS.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close() //nolint:errcheck
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, res.StatusCode, "Response body: %s", body)
	return gjson.ParseBytes(body)
}

func RevokeToken(t testing.TB, conf *oauth2.Config, token string, publicTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token)

	req := httpx.MustNewRequest("POST", publicTS.URL+"/oauth2/revoke",
		strings.NewReader((url.Values{"token": {token}}).Encode()),
		"application/x-www-form-urlencoded")

	req.SetBasicAuth(conf.ClientID, conf.ClientSecret)
	res, err := publicTS.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close() //nolint:errcheck
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
	defer res.Body.Close() //nolint:errcheck
	require.Equal(t, res.StatusCode, http.StatusOK)
}

func Userinfo(t *testing.T, token *oauth2.Token, publicTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token.AccessToken)

	req := httpx.MustNewRequest("GET", publicTS.URL+"/userinfo", nil, "")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	res, err := publicTS.Client().Do(req)
	require.NoError(t, err)

	defer res.Body.Close() //nolint:errcheck
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

func NewDeviceLoginConsentUI(t testing.TB, c *config.DefaultProvider, device, login, consent http.HandlerFunc) {
	if device == nil {
		device = HTTPServerNotImplementedHandler
	}
	dt := httptest.NewServer(device)
	t.Cleanup(dt.Close)
	c.MustSet(context.Background(), config.KeyDeviceVerificationURL, dt.URL)

	NewLoginConsentUI(t, c, login, consent)
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

// InsecureDecodeJWT decodes a JWT payload without checking the signature.
func InsecureDecodeJWT(t require.TestingT, token string) []byte {
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	dec, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoErrorf(t, err, "failed to decode JWT payload: %s", parts[1])
	return dec
}

var (
	NewEmptyCookieJar = x.NewEmptyCookieJar
	NewEmptyJarClient = x.NewEmptyJarClient
)

func AssertTokenValid(t *testing.T, accessOrIDToken gjson.Result, sub string) {
	assert.Equal(t, sub, accessOrIDToken.Get("sub").Str)
	assert.WithinDurationf(t, time.Now(), time.Unix(accessOrIDToken.Get("iat").Int(), 0), time.Minute, "%s", accessOrIDToken.Raw)
	assert.Truef(t, time.Now().Before(time.Unix(accessOrIDToken.Get("exp").Int(), 0)), "%s", accessOrIDToken.Raw)
}

func AssertAccessToken(t *testing.T, token gjson.Result, sub, clientID string) {
	AssertTokenValid(t, token, sub)
	assert.Equalf(t, clientID, token.Get("client_id").Str, "%s", token.Raw)
	assert.WithinDurationf(t, time.Now(), time.Unix(token.Get("nbf").Int(), 0), time.Minute, "%s", token.Raw)
}

func AssertIDToken(t *testing.T, token gjson.Result, sub, clientID string) {
	AssertTokenValid(t, token, sub)
	assert.Equalf(t, clientID, token.Get("aud.0").Str, "%s", token.Raw)
	assert.WithinDurationf(t, time.Now(), time.Unix(token.Get("rat").Int(), 0), time.Minute, "%s", token.Raw)
}
