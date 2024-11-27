// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/fosite/token/jwt"
	"github.com/ory/x/configx"
	"github.com/ory/x/otelx"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/urlx"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/v2/x"
)

func newProvider() *DefaultProvider {
	return MustNew(context.Background(), logrusx.New("", ""))
}

func setupEnv(env map[string]string) func(t *testing.T) func() {
	return func(t *testing.T) (setup func()) {
		setup = func() {
			for k, v := range env {
				t.Setenv(k, v)
			}
		}
		return
	}
}

func TestSubjectTypesSupported(t *testing.T) {
	ctx := context.Background()
	for k, tc := range []struct {
		d   string
		env func(t *testing.T) func()
		e   []string
	}{
		{
			d: "Load legacy environment variable in legacy format",
			env: setupEnv(map[string]string{
				strings.ToUpper(strings.Replace(KeySubjectTypesSupported, ".", "_", -1)):                 "public,pairwise",
				strings.ToUpper(strings.Replace("oidc.subject_identifiers.pairwise.salt", ".", "_", -1)): "some-salt",
			}),
			e: []string{"public", "pairwise"},
		},
		{
			d: "Load legacy environment variable in legacy format with JWT enabled",
			env: setupEnv(map[string]string{
				strings.ToUpper(strings.Replace(KeySubjectTypesSupported, ".", "_", -1)):                 "public,pairwise",
				strings.ToUpper(strings.Replace("oidc.subject_identifiers.pairwise.salt", ".", "_", -1)): "some-salt",
				strings.ToUpper(strings.Replace(KeyAccessTokenStrategy, ".", "_", -1)):                   "jwt",
			}),
			e: []string{"public"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			setup := tc.env(t)
			setup()
			p := newProvider()
			p.MustSet(ctx, KeySubjectIdentifierAlgorithmSalt, "00000000")
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported(ctx))
		})
	}
}

func TestWellKnownKeysUnique(t *testing.T) {
	p := newProvider()
	assert.EqualValues(t, []string{x.OpenIDConnectKeyName, x.OAuth2JWTKeyName}, p.WellKnownKeys(context.Background(), x.OAuth2JWTKeyName, x.OpenIDConnectKeyName, x.OpenIDConnectKeyName))
}

func TestCORSOptions(t *testing.T) {
	ctx := context.Background()
	p := newProvider()
	p.MustSet(ctx, "serve.public.cors.enabled", true)

	conf, enabled := p.CORS(ctx, PublicInterface)
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
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)

	p := MustNew(context.Background(), l)

	value := p.DisableHealthAccessLog(AdminInterface)
	assert.Equal(t, false, value)

	p.MustSet(ctx, AdminInterface.Key(KeySuffixDisableHealthAccessLog), "true")

	value = p.DisableHealthAccessLog(AdminInterface)
	assert.Equal(t, true, value)
}

func TestProviderPublicDisableHealthAccessLog(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)

	p := MustNew(context.Background(), l)

	value := p.DisableHealthAccessLog(PublicInterface)
	assert.Equal(t, false, value)

	p.MustSet(ctx, PublicInterface.Key(KeySuffixDisableHealthAccessLog), "true")

	value = p.DisableHealthAccessLog(PublicInterface)
	assert.Equal(t, true, value)
}

func TestPublicAllowDynamicRegistration(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)

	p := MustNew(context.Background(), l)

	value := p.PublicAllowDynamicRegistration(ctx)
	assert.Equal(t, false, value)

	p.MustSet(ctx, KeyPublicAllowDynamicRegistration, "true")

	value = p.PublicAllowDynamicRegistration(ctx)
	assert.Equal(t, true, value)
}

func TestProviderIssuerURL(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost")
	assert.Equal(t, "http://hydra.localhost", p.IssuerURL(ctx).String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost/")
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL(ctx).String())
}

func TestProviderIssuerPublicURL(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost")
	p.MustSet(ctx, KeyPublicURL, "http://hydra.example")

	assert.Equal(t, "http://hydra.localhost", p.IssuerURL(ctx).String())
	assert.Equal(t, "http://hydra.example/", p.PublicURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p.JWKSURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/consent", p.ConsentURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/login", p.LoginURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/logout", p.LogoutURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/token", p.OAuth2TokenURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/auth", p.OAuth2AuthURL(ctx).String())
	assert.Equal(t, "http://hydra.example/userinfo", p.OIDCDiscoveryUserinfoEndpoint(ctx).String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost/")
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/", p2.PublicURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p2.JWKSURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/consent", p2.ConsentURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/login", p2.LoginURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/logout", p2.LogoutURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/token", p2.OAuth2TokenURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/oauth2/auth", p2.OAuth2AuthURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/userinfo", p2.OIDCDiscoveryUserinfoEndpoint(ctx).String())
}

func TestProviderCookieSameSiteMode(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)

	p := MustNew(context.Background(), l, configx.SkipValidation())
	p.MustSet(ctx, KeyTLSEnabled, true)

	p.MustSet(ctx, KeyCookieSameSiteMode, "")
	assert.Equal(t, http.SameSiteDefaultMode, p.CookieSameSiteMode(ctx))

	p.MustSet(ctx, KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode(ctx))

	p.MustSet(ctx, KeyCookieSameSiteMode, "lax")
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode(ctx))

	p.MustSet(ctx, KeyCookieSameSiteMode, "strict")
	assert.Equal(t, http.SameSiteStrictMode, p.CookieSameSiteMode(ctx))

	p = MustNew(context.Background(), l, configx.SkipValidation())
	p.MustSet(ctx, "dev", true)
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode(ctx))
	p.MustSet(ctx, KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode(ctx))

	p.MustSet(ctx, KeyIssuerURL, "https://example.com")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode(ctx))
}

func TestViperProviderValidates(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	c := MustNew(context.Background(), l, configx.WithConfigFiles("../../internal/.hydra.yaml"))

	// log
	assert.Equal(t, "debug", c.Source(ctx).String(KeyLogLevel))
	assert.Equal(t, "json", c.Source(ctx).String("log.format"))

	// serve
	assert.Equal(t, "localhost:1", c.ListenOn(PublicInterface))
	assert.Equal(t, "localhost:2", c.ListenOn(AdminInterface))

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
	assert.Equal(t, expectedPublicPermission, c.SocketPermission(PublicInterface))
	assert.Equal(t, expectedAdminPermission, c.SocketPermission(AdminInterface))

	expectedCors := cors.Options{
		AllowedOrigins:   []string{"https://example.com"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Authorization"},
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           1,
		Debug:            false,
	}

	gc, enabled := c.CORS(ctx, AdminInterface)
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	gc, enabled = c.CORS(ctx, PublicInterface)
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	assert.Equal(t, []string{"127.0.0.1/32"}, c.TLS(ctx, PublicInterface).AllowTerminationFrom())
	assert.Equal(t, "/path/to/file.pem", c.Source(ctx).String("serve.tls.key.path"))
	assert.Equal(t, "b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK", c.Source(ctx).String("serve.tls.cert.base64"))
	assert.Equal(t, http.SameSiteLaxMode, c.CookieSameSiteMode(ctx))
	assert.Equal(t, true, c.CookieSameSiteLegacyWorkaround(ctx))

	// dsn
	assert.Contains(t, c.DSN(), "sqlite://")

	// webfinger
	assert.Equal(t, []string{"hydra.openid.id-token", "hydra.jwt.access-token"}, c.WellKnownKeys(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://example.com"), c.OAuth2ClientRegistrationURL(ctx))
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
	assert.Equal(t, time.Duration(0), c.RefreshTokenRotationGracePeriod(ctx))
	require.NoError(t, c.Set(ctx, KeyRefreshTokenRotationGracePeriod, "1s"))
	assert.Equal(t, time.Second, c.RefreshTokenRotationGracePeriod(ctx))
	require.NoError(t, c.Set(ctx, KeyRefreshTokenRotationGracePeriod, "2h"))
	assert.Equal(t, time.Minute*5, c.RefreshTokenRotationGracePeriod(ctx))

	// urls
	assert.Equal(t, urlx.ParseOrPanic("https://issuer"), c.IssuerURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://public/"), c.PublicURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://admin/"), c.AdminURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://login/"), c.LoginURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://consent/"), c.ConsentURL(ctx))
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

	// oauth2
	assert.Equal(t, true, c.GetSendDebugMessagesToClients(ctx))
	assert.Equal(t, 20, c.GetBCryptCost(ctx))
	assert.Equal(t, true, c.GetEnforcePKCE(ctx))
	assert.Equal(t, true, c.GetEnforcePKCEForPublicClients(ctx))

	// secrets
	secret, err := c.GetGlobalSecret(ctx)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x64, 0x40, 0x5f, 0xd4, 0x66, 0xc9, 0x8c, 0x88, 0xa7, 0xf2, 0xcb, 0x95, 0xcd, 0x95, 0xcb, 0xa3, 0x41, 0x49, 0x8b, 0x97, 0xba, 0x9e, 0x92, 0xee, 0x4c, 0xaf, 0xe0, 0x71, 0x23, 0x28, 0xeb, 0xfc}, secret)

	cookieSecret, err := c.GetCookieSecrets(ctx)
	require.NoError(t, err)
	assert.Equal(t, [][]uint8{[]byte("some-random-cookie-secret")}, cookieSecret)

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
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(ctx, KeyLoginURL, "http://localhost:8080/oauth/login")
	p.MustSet(ctx, KeyConsentURL, "http://localhost:8080/oauth/consent")

	assert.Equal(t, "http://localhost:8080/oauth/login", p.LoginURL(ctx).String())
	assert.Equal(t, "http://localhost:8080/oauth/consent", p.ConsentURL(ctx).String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(ctx, KeyLoginURL, "http://localhost:3000/#/oauth/login")
	p2.MustSet(ctx, KeyConsentURL, "http://localhost:3000/#/oauth/consent")

	assert.Equal(t, "http://localhost:3000/#/oauth/login", p2.LoginURL(ctx).String())
	assert.Equal(t, "http://localhost:3000/#/oauth/consent", p2.ConsentURL(ctx).String())
}

func TestInfinitRefreshTokenTTL(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	c := MustNew(context.Background(), l, configx.WithValue("ttl.refresh_token", -1))

	assert.Equal(t, -1*time.Nanosecond, c.GetRefreshTokenLifespan(ctx))
}

func TestCookieSecure(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(io.Discard)
	c := MustNew(context.Background(), l, configx.WithValue(KeyDevelopmentMode, true))

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
	c := MustNew(context.Background(), l, configx.SkipValidation())

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
	p := MustNew(context.Background(), l)

	ctx := context.Background()
	// p.MustSet(ctx, KeyOAuth2GrantJWTClientAuthOptional, false)
	p.MustSet(ctx, KeyOAuth2GrantJWTMaxDuration, "1h")
	p.MustSet(ctx, KeyOAuth2GrantJWTIssuedDateOptional, false)
	p.MustSet(ctx, KeyOAuth2GrantJWTIDOptional, false)

	// assert.Equal(t, false, p.GetGrantTypeJWTBearerCanSkipClientAuth(ctx))
	assert.Equal(t, 1.0, p.GetJWTMaxDuration(ctx).Hours())
	assert.Equal(t, false, p.GetGrantTypeJWTBearerIssuedDateOptional(ctx))
	assert.Equal(t, false, p.GetGrantTypeJWTBearerIDOptional(ctx))

	p2 := MustNew(context.Background(), l)

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
	p := MustNew(context.Background(), l)

	ctx := context.Background()

	assert.Equal(t, jwt.JWTScopeFieldList, p.GetJWTScopeField(ctx))
	p.MustSet(ctx, KeyJWTScopeClaimStrategy, "list")
	assert.Equal(t, jwt.JWTScopeFieldList, p.GetJWTScopeField(ctx))
	p.MustSet(ctx, KeyJWTScopeClaimStrategy, "string")
	assert.Equal(t, jwt.JWTScopeFieldString, p.GetJWTScopeField(ctx))
	p.MustSet(ctx, KeyJWTScopeClaimStrategy, "both")
	assert.Equal(t, jwt.JWTScopeFieldBoth, p.GetJWTScopeField(ctx))
}
