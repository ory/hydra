// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/randx"
	"github.com/ory/x/urlx"
)

func newProvider(t *testing.T, opts ...configx.OptionModifier) *DefaultProvider {
	return MustNew(t, logrusx.New("", ""), opts...)
}

func TestSubjectTypesSupported(t *testing.T) {
	ctx := t.Context()
	for _, tc := range []struct {
		d    string
		vals map[string]any
		e    []string
	}{{
		d:    "no subject types",
		vals: map[string]any{KeySubjectTypesSupported: []string{}},
		e:    []string{"public"},
	}, {
		d:    "public",
		vals: map[string]any{KeySubjectTypesSupported: []string{"public"}},
		e:    []string{"public"},
	}, {
		d: "pairwise",
		vals: map[string]any{
			KeySubjectTypesSupported:          []string{"pairwise"},
			KeySubjectIdentifierAlgorithmSalt: "00000000",
		},
		e: []string{"pairwise"},
	}, {
		d: "public and pairwise",
		vals: map[string]any{
			KeySubjectTypesSupported:          []string{"public", "pairwise"},
			KeySubjectIdentifierAlgorithmSalt: "00000000",
		},
		e: []string{"public", "pairwise"},
	}, {
		d: "pairwise disabled with jwt",
		vals: map[string]any{
			KeySubjectTypesSupported: []string{"public", "pairwise"},
			KeyAccessTokenStrategy:   "jwt",
		},
		e: []string{"public"},
	}, {
		d: "unknown subject type",
		vals: map[string]any{
			KeySubjectTypesSupported:          []string{"public", "pairwise", "unknown"},
			KeySubjectIdentifierAlgorithmSalt: "00000000",
		},
		e: []string{"public", "pairwise"},
	}} {
		t.Run(tc.d, func(t *testing.T) {
			p := newProvider(t, configx.WithValues(tc.vals), configx.SkipValidation())
			assert.Equal(t, tc.e, p.SubjectTypesSupported(ctx))
		})
	}
}

func TestWellKnownKeysUnique(t *testing.T) {
	p := newProvider(t)
	assert.EqualValues(t, []string{x.OpenIDConnectKeyName, x.OAuth2JWTKeyName}, p.WellKnownKeys(t.Context(), x.OAuth2JWTKeyName, x.OpenIDConnectKeyName, x.OpenIDConnectKeyName))
}

func TestCORSOptions(t *testing.T) {
	ctx := context.Background()
	p := newProvider(t, configx.WithValue("serve.public.cors.enabled", true))

	conf, enabled := p.CORSPublic(ctx)
	assert.True(t, enabled)

	assert.EqualValues(t, cors.Options{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"POST", "GET", "PUT", "PATCH", "DELETE", "CONNECT", "HEAD", "OPTIONS", "TRACE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Language", "Content-Language", "Authorization"},
		ExposedHeaders:   []string{"Cache-Control", "Expires", "Last-Modified", "Pragma", "Content-Length", "Content-Language", "Content-Type"},
		AllowCredentials: true,
	}, conf)
}

func TestProviderAdminDisableHealthAccessLog(t *testing.T) {
	p := newProvider(t)
	serve := p.ServeAdmin(t.Context())
	assert.False(t, serve.RequestLog.DisableHealth)

	p = newProvider(t, configx.WithValue("serve.admin.requestlog.disable_health", true))
	serve = p.ServeAdmin(t.Context())
	assert.True(t, serve.RequestLog.DisableHealth)
}

func TestProviderPublicDisableHealthAccessLog(t *testing.T) {
	p := newProvider(t)
	serve := p.ServePublic(t.Context())
	assert.False(t, serve.RequestLog.DisableHealth)

	p = newProvider(t, configx.WithValue("serve.public.requestlog.disable_health", true))
	serve = p.ServePublic(t.Context())
	assert.True(t, serve.RequestLog.DisableHealth)
}

func TestPublicAllowDynamicRegistration(t *testing.T) {
	p := newProvider(t)
	value := p.PublicAllowDynamicRegistration(t.Context())
	assert.False(t, value)

	p = newProvider(t, configx.WithValue(KeyPublicAllowDynamicRegistration, true))
	value = p.PublicAllowDynamicRegistration(t.Context())
	assert.True(t, value)
}

func TestProviderIssuerURL(t *testing.T) {
	p := newProvider(t, configx.WithValue(KeyIssuerURL, "http://hydra.localhost"))
	assert.Equal(t, "http://hydra.localhost", p.IssuerURL(t.Context()).String())
}

func TestProviderIssuerPublicURL(t *testing.T) {
	p := newProvider(t, configx.WithValues(map[string]any{
		KeyIssuerURL: "http://hydra.localhost",
		KeyPublicURL: "http://hydra.example",
	}))

	assert.Equal(t, "http://hydra.localhost", p.IssuerURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/", p.PublicURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p.JWKSURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/consent", p.ConsentURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/login", p.LoginURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/logout", p.LogoutURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/oauth2/token", p.OAuth2TokenURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/oauth2/auth", p.OAuth2AuthURL(t.Context()).String())
	assert.Equal(t, "http://hydra.example/userinfo", p.OIDCDiscoveryUserinfoEndpoint(t.Context()).String())

	p = newProvider(t, configx.WithValue(KeyIssuerURL, "http://hydra.localhost/"))
	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/", p.PublicURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p.JWKSURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/consent", p.ConsentURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/login", p.LoginURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/logout", p.LogoutURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/token", p.OAuth2TokenURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/auth", p.OAuth2AuthURL(t.Context()).String())
	assert.Equal(t, "http://hydra.localhost/userinfo", p.OIDCDiscoveryUserinfoEndpoint(t.Context()).String())
}

func TestProviderCookieSameSiteMode(t *testing.T) {
	for _, tc := range []struct {
		d, mode  string
		others   map[string]any
		expected http.SameSite
	}{{
		d:        "default",
		mode:     "",
		expected: http.SameSiteDefaultMode,
	}, {
		d:        "default dev",
		mode:     "",
		others:   map[string]any{KeyDevelopmentMode: true},
		expected: http.SameSiteLaxMode,
	}, {
		d:        "none with http",
		mode:     "none",
		others:   map[string]any{KeyIssuerURL: "http://example.com"},
		expected: http.SameSiteLaxMode,
	}, {
		d:        "none with https",
		mode:     "none",
		others:   map[string]any{KeyIssuerURL: "https://example.com"},
		expected: http.SameSiteNoneMode,
	}, {
		d:        "lax",
		mode:     "lax",
		expected: http.SameSiteLaxMode,
	}, {
		d:        "strict",
		mode:     "strict",
		expected: http.SameSiteStrictMode,
	}} {
		t.Run(tc.d, func(t *testing.T) {
			p := newProvider(t, configx.WithValue(KeyCookieSameSiteMode, tc.mode), configx.WithValues(tc.others), configx.SkipValidation())
			assert.Equal(t, tc.expected, p.CookieSameSiteMode(t.Context()))
		})
	}
}

func TestProviderValidates(t *testing.T) {
	ctx := t.Context()
	c := newProvider(t, configx.WithConfigFiles("../../internal/.hydra.yaml"))

	// log
	assert.Equal(t, "debug", c.Source(ctx).String(KeyLogLevel))
	assert.Equal(t, "json", c.Source(ctx).String("log.format"))

	// serve
	servePublic, serveAdmin := c.ServePublic(ctx), c.ServeAdmin(ctx)
	assert.Equal(t, "localhost", servePublic.Host)
	assert.Equal(t, 1, servePublic.Port)
	assert.Equal(t, "localhost", serveAdmin.Host)
	assert.Equal(t, 2, serveAdmin.Port)

	expectedPublicPermission := &configx.UnixPermission{
		Owner: "hydra",
		Group: "hydra-public-api",
		Mode:  0775,
	}
	expectedAdminPermission := &configx.UnixPermission{
		Owner: "hydra",
		Group: "hydra-admin-api",
		Mode:  0770,
	}
	assert.Equal(t, expectedPublicPermission, &servePublic.Socket)
	assert.Equal(t, expectedAdminPermission, &serveAdmin.Socket)

	expectedCors := cors.Options{
		AllowedOrigins:   []string{"https://example.com"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Authorization"},
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           1,
		Debug:            false,
	}

	gc, enabled := c.CORSAdmin(ctx)
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	gc, enabled = c.CORSPublic(ctx)
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	assert.Equal(t, []string{"127.0.0.1/32"}, c.Source(ctx).Strings("serve.tls.allow_termination_from"))
	assert.Equal(t, []string{"127.0.0.1/32"}, servePublic.TLS.AllowTerminationFrom)
	assert.Equal(t, []string{"127.0.0.1/32"}, serveAdmin.TLS.AllowTerminationFrom)
	assert.Equal(t, "/path/to/file.pem", c.Source(ctx).String("serve.tls.key.path"))
	assert.Equal(t, "/path/to/file.pem", servePublic.TLS.KeyPath)
	assert.Equal(t, "/path/to/file.pem", serveAdmin.TLS.KeyPath)
	assert.Equal(t, "b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK", c.Source(ctx).String("serve.tls.cert.base64"))
	assert.Equal(t, "b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK", servePublic.TLS.CertBase64)
	assert.Equal(t, "b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK", serveAdmin.TLS.CertBase64)
	assert.Equal(t, http.SameSiteLaxMode, c.CookieSameSiteMode(ctx))
	assert.Equal(t, true, c.CookieSameSiteLegacyWorkaround(ctx))

	// dsn
	assert.Contains(t, c.DSN(), "sqlite://")

	// webfinger
	assert.Equal(t, []string{"hydra.openid.id-token", "hydra.jwt.access-token"}, c.WellKnownKeys(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com"), c.OAuth2ClientRegistrationURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/device_authorization"), c.OAuth2DeviceAuthorisationURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/jwks.json"), c.JWKSURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/auth"), c.OAuth2AuthURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/token"), c.OAuth2TokenURL(ctx))
	assert.Equal(t, []string{"sub", "username"}, c.OIDCDiscoverySupportedClaims(ctx))
	assert.Equal(t, []string{"offline_access", "offline", "openid", "whatever"}, c.OIDCDiscoverySupportedScope(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com"), c.OIDCDiscoveryUserinfoEndpoint(ctx))

	// oidc
	assert.Equal(t, []string{"pairwise"}, c.SubjectTypesSupported(ctx))
	assert.Equal(t, "random_salt", c.SubjectIdentifierAlgorithmSalt(ctx))
	assert.Equal(t, []string{"whatever"}, c.DefaultClientScope(ctx))

	// refresh
	assert.Equal(t, GracefulRefreshTokenRotation{}, c.GracefulRefreshTokenRotation(ctx))
	require.NoError(t, c.Set(ctx, KeyRefreshTokenRotationGracePeriod, "1s"))
	assert.Equal(t, time.Second, c.GracefulRefreshTokenRotation(ctx).Period)
	require.NoError(t, c.Set(ctx, KeyRefreshTokenRotationGracePeriod, "2h"))
	assert.Equal(t, 5*time.Minute, c.GracefulRefreshTokenRotation(ctx).Period)
	require.NoError(t, c.Set(ctx, KeyRefreshTokenRotationGraceReuseCount, "2"))
	assert.Equal(t, GracefulRefreshTokenRotation{Count: 2, Period: 2 * time.Hour}, c.GracefulRefreshTokenRotation(ctx))
	require.NoError(t, c.Set(ctx, KeyRefreshTokenRotationGracePeriod, (time.Hour*24*200).String()))
	assert.Equal(t, GracefulRefreshTokenRotation{Count: 2, Period: time.Hour * 24 * 180}, c.GracefulRefreshTokenRotation(ctx))

	// urls
	assert.Equal(t, urlx.ParseOrPanic("https://issuer"), c.IssuerURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://public/"), c.PublicURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://admin/"), c.AdminURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://login/"), c.LoginURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://consent/"), c.ConsentURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://device/"), c.DeviceVerificationURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://device/callback"), c.DeviceDoneURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://logout/"), c.LogoutURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://error/"), c.ErrorURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://post_logout/"), c.LogoutRedirectURL(ctx))

	// strategies
	assert.True(t, c.GetScopeStrategy(ctx)([]string{"openid"}, "openid"), "should us fosite.ExactScopeStrategy")
	assert.False(t, c.GetScopeStrategy(ctx)([]string{"openid.*"}, "openid.email"), "should us fosite.ExactScopeStrategy")
	assert.Equal(t, AccessTokenDefaultStrategy, c.AccessTokenStrategy(ctx))
	assert.Equal(t, false, c.GrantAllClientCredentialsScopesPerDefault(ctx))
	assert.Equal(t, jwt.JWTScopeFieldList, c.GetJWTScopeField(ctx))

	// ttl
	assert.Equal(t, 2*time.Hour, c.ConsentRequestMaxAge(ctx))
	assert.Equal(t, 2*time.Hour, c.GetAccessTokenLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetRefreshTokenLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetIDTokenLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetAuthorizeCodeLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetDeviceAndUserCodeLifespan(ctx))
	assert.Equal(t, 24*time.Hour, c.GetAuthenticationSessionLifespan(ctx))

	// oauth2
	assert.Equal(t, true, c.GetSendDebugMessagesToClients(ctx))
	assert.Equal(t, 20, c.GetBCryptCost(ctx))
	assert.Equal(t, true, c.GetEnforcePKCE(ctx))
	assert.Equal(t, true, c.GetEnforcePKCEForPublicClients(ctx))
	assert.Equal(t, 2*time.Hour, c.GetDeviceAuthTokenPollingInterval(ctx))
	assert.Equal(t, 8, c.GetUserCodeLength(ctx))
	assert.Equal(t, string(randx.AlphaUpper), string(c.GetUserCodeSymbols(ctx)))

	// secrets
	secret, err := c.GetGlobalSecret(ctx)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x64, 0x40, 0x5f, 0xd4, 0x66, 0xc9, 0x8c, 0x88, 0xa7, 0xf2, 0xcb, 0x95, 0xcd, 0x95, 0xcb, 0xa3, 0x41, 0x49, 0x8b, 0x97, 0xba, 0x9e, 0x92, 0xee, 0x4c, 0xaf, 0xe0, 0x71, 0x23, 0x28, 0xeb, 0xfc}, secret)

	cookieSecret, err := c.GetCookieSecrets(ctx)
	require.NoError(t, err)
	assert.Equal(t, [][]byte{[]byte("some-random-cookie-secret")}, cookieSecret)

	paginationKeys := c.GetPaginationEncryptionKeys(ctx)
	require.Len(t, paginationKeys, 1)
	assert.Equal(t, [32]byte{0x1a, 0x4c, 0x1, 0xbc, 0x1b, 0xd1, 0x4c, 0xdf, 0x23, 0x3, 0xd9, 0x1a, 0x2a, 0x1b, 0x68, 0xdc, 0x69, 0x17, 0xf4, 0x31, 0xd, 0x27, 0x6d, 0x86, 0x70, 0xb0, 0xae, 0x2d, 0x45, 0xe2, 0xf, 0xab}, paginationKeys[0])

	// profiling
	assert.Equal(t, "cpu", c.Source(ctx).String("profiling"))

	// tracing
	assert.EqualValues(t, &otelx.Config{
		ServiceName: "hydra service",
		Provider:    "jaeger",
		Providers: otelx.ProvidersConfig{
			Jaeger: otelx.JaegerConfig{
				LocalAgentAddress: "127.0.0.1:6831",
				Sampling: otelx.JaegerSampling{
					ServerURL:    "http://sampling",
					TraceIdRatio: 1,
				},
			},
			Zipkin: otelx.ZipkinConfig{
				ServerURL: "http://zipkin/api/v2/spans",
			},
			OTLP: otelx.OTLPConfig{
				ServerURL: "localhost:4318",
				Insecure:  true,
				Sampling: otelx.OTLPSampling{
					SamplingRatio: 1.0,
				},
			},
		},
	}, c.Tracing())
}

func TestSetPerm(t *testing.T) {
	f, e := os.CreateTemp("", "test")
	require.NoError(t, e)
	path := f.Name()

	// We cannot test setting owner and group, because we don't know what the
	// tester has access to.
	_ = (&configx.UnixPermission{
		Owner: "",
		Group: "",
		Mode:  0654,
	}).SetPermission(path)

	stat, err := f.Stat()
	require.NoError(t, err)

	assert.Equal(t, os.FileMode(0654), stat.Mode())

	require.NoError(t, f.Close())
	require.NoError(t, os.Remove(path))
}

func TestLoginConsentURL(t *testing.T) {
	p := newProvider(t, configx.WithValues(map[string]any{
		KeyLoginURL:              "http://localhost:8080/oauth/login",
		KeyConsentURL:            "http://localhost:8080/oauth/consent",
		KeyDeviceVerificationURL: "http://localhost:8080/oauth/device",
	}))

	assert.Equal(t, "http://localhost:8080/oauth/login", p.LoginURL(t.Context()).String())
	assert.Equal(t, "http://localhost:8080/oauth/consent", p.ConsentURL(t.Context()).String())
	assert.Equal(t, "http://localhost:8080/oauth/device", p.DeviceVerificationURL(t.Context()).String())

	p = newProvider(t, configx.WithValues(map[string]any{
		KeyLoginURL:              "http://localhost:3000/#/oauth/login",
		KeyConsentURL:            "http://localhost:3000/#/oauth/consent",
		KeyDeviceVerificationURL: "http://localhost:3000/#/oauth/device",
	}))

	assert.Equal(t, "http://localhost:3000/#/oauth/login", p.LoginURL(t.Context()).String())
	assert.Equal(t, "http://localhost:3000/#/oauth/consent", p.ConsentURL(t.Context()).String())
	assert.Equal(t, "http://localhost:3000/#/oauth/device", p.DeviceVerificationURL(t.Context()).String())
}

func TestInfinityRefreshTokenTTL(t *testing.T) {
	c := newProvider(t, configx.WithValue("ttl.refresh_token", -1))

	assert.Equal(t, time.Duration(-1), c.GetRefreshTokenLifespan(t.Context()))
}

func TestLimitAuthSessionLifespan(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	c := MustNew(t, l)
	assert.Equal(t, 30*24*time.Hour, c.GetAuthenticationSessionLifespan(ctx))

	require.NoError(t, c.Set(ctx, KeyAuthenticationSessionLifespan, (time.Hour*24*300).String()))
	assert.Equal(t, 180*24*time.Hour, c.GetAuthenticationSessionLifespan(ctx))
}

func TestCookieSecure(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	c := MustNew(t, l, configx.WithValue(KeyDevelopmentMode, true))

	c.MustSet(ctx, KeyCookieSecure, true)
	assert.True(t, c.CookieSecure(ctx))

	c.MustSet(ctx, KeyCookieSecure, false)
	assert.False(t, c.CookieSecure(ctx))

	c.MustSet(ctx, KeyDevelopmentMode, false)
	assert.True(t, c.CookieSecure(ctx))
}

func TestHookConfigs(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	c := MustNew(t, l, configx.SkipValidation())

	for key, getFunc := range map[string]func(context.Context) *HookConfig{
		KeyRefreshTokenHook: c.TokenRefreshHookConfig,
		KeyTokenHook:        c.TokenHookConfig,
	} {
		assert.Nil(t, getFunc(ctx))
		c.MustSet(ctx, key, "")
		assert.Nil(t, getFunc(ctx))
		c.MustSet(ctx, key, "http://localhost:8080/hook")
		hc := getFunc(ctx)
		require.NotNil(t, hc)
		assert.EqualValues(t, "http://localhost:8080/hook", hc.URL)

		c.MustSet(ctx, key, `
{
	"url": "http://localhost:8080/hook2",
	"auth": {
		"type": "api_key",
		"config": {
			"in": "header",
			"name": "my-header",
			"value": "my-value"
		}
	}
}`)
		hc = getFunc(ctx)
		require.NotNil(t, hc)
		assert.EqualValues(t, "http://localhost:8080/hook2", hc.URL)
		assert.EqualValues(t, "api_key", hc.Auth.Type)
		rawConfig, err := json.Marshal(hc.Auth.Config)
		require.NoError(t, err)
		assert.JSONEq(t, `{"in":"header","name":"my-header","value":"my-value"}`, string(rawConfig))
	}
}

func TestJWTBearer(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	p := MustNew(t, l)

	ctx := context.Background()
	// p.MustSet(ctx, KeyOAuth2GrantJWTClientAuthOptional, false)
	p.MustSet(ctx, KeyOAuth2GrantJWTMaxDuration, "1h")
	p.MustSet(ctx, KeyOAuth2GrantJWTIssuedDateOptional, false)
	p.MustSet(ctx, KeyOAuth2GrantJWTIDOptional, false)

	// assert.Equal(t, false, p.GetGrantTypeJWTBearerCanSkipClientAuth(ctx))
	assert.Equal(t, 1.0, p.GetJWTMaxDuration(ctx).Hours())
	assert.Equal(t, false, p.GetGrantTypeJWTBearerIssuedDateOptional(ctx))
	assert.Equal(t, false, p.GetGrantTypeJWTBearerIDOptional(ctx))

	p2 := MustNew(t, l)

	// p2.MustSet(ctx, KeyOAuth2GrantJWTClientAuthOptional, true)
	p2.MustSet(ctx, KeyOAuth2GrantJWTMaxDuration, "24h")
	p2.MustSet(ctx, KeyOAuth2GrantJWTIssuedDateOptional, true)
	p2.MustSet(ctx, KeyOAuth2GrantJWTIDOptional, true)

	// assert.Equal(t, true, p2.GetGrantTypeJWTBearerCanSkipClientAuth(ctx))
	assert.Equal(t, 24.0, p2.GetJWTMaxDuration(ctx).Hours())
	assert.Equal(t, true, p2.GetGrantTypeJWTBearerIssuedDateOptional(ctx))
	assert.Equal(t, true, p2.GetGrantTypeJWTBearerIDOptional(ctx))
}

func TestJWTScopeClaimStrategy(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	p := MustNew(t, l)

	ctx := context.Background()

	assert.Equal(t, jwt.JWTScopeFieldList, p.GetJWTScopeField(ctx))
	p.MustSet(ctx, KeyJWTScopeClaimStrategy, "list")
	assert.Equal(t, jwt.JWTScopeFieldList, p.GetJWTScopeField(ctx))
	p.MustSet(ctx, KeyJWTScopeClaimStrategy, "string")
	assert.Equal(t, jwt.JWTScopeFieldString, p.GetJWTScopeField(ctx))
	p.MustSet(ctx, KeyJWTScopeClaimStrategy, "both")
	assert.Equal(t, jwt.JWTScopeFieldBoth, p.GetJWTScopeField(ctx))
}

func TestDeviceUserCode(t *testing.T) {
	l := logrusx.New("", "")

	t.Run("preset", func(t *testing.T) {
		p := MustNew(t, l, configx.WithValue(KeyDeviceAuthUserCodeEntropyPreset, "low"))
		assert.Equal(t, 9, p.GetUserCodeLength(t.Context()))
		assert.Equal(t, string(randx.Numeric), string(p.GetUserCodeSymbols(t.Context())))
	})

	t.Run("explicit values", func(t *testing.T) {
		length, charSet := 15, "foobarbaz1234"
		p := MustNew(t, l, configx.WithValues(map[string]any{
			KeyDeviceAuthUserCodeLength:       length,
			KeyDeviceAuthUserCodeCharacterSet: charSet,
		}))
		assert.Equal(t, length, p.GetUserCodeLength(t.Context()))
		assert.Equal(t, charSet, string(p.GetUserCodeSymbols(t.Context())))
	})
}
