package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/x/configx"
	"github.com/ory/x/otelx"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/urlx"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/x"
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
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:     []string{"Authorization", "Content-Type"},
		ExposedHeaders:     []string{"Content-Type"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		MaxAge:             0,
		Debug:              false,
	}, conf)
}

func TestProviderAdminDisableHealthAccessLog(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l)

	value := p.DisableHealthAccessLog(ctx, AdminInterface)
	assert.Equal(t, false, value)

	p.MustSet(ctx, AdminInterface.Key(KeySuffixDisableHealthAccessLog), "true")

	value = p.DisableHealthAccessLog(ctx, AdminInterface)
	assert.Equal(t, true, value)
}

func TestProviderPublicDisableHealthAccessLog(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l)

	value := p.DisableHealthAccessLog(ctx, PublicInterface)
	assert.Equal(t, false, value)

	p.MustSet(ctx, PublicInterface.Key(KeySuffixDisableHealthAccessLog), "true")

	value = p.DisableHealthAccessLog(ctx, PublicInterface)
	assert.Equal(t, true, value)
}

func TestPublicAllowDynamicRegistration(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

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
	l.Logrus().SetOutput(ioutil.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost")
	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL(ctx).String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost/")
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL(ctx).String())
}

func TestProviderIssuerPublicURL(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost")
	p.MustSet(ctx, KeyPublicURL, "http://hydra.example")

	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL(ctx).String())
	assert.Equal(t, "http://hydra.example/", p.PublicURL(ctx).String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p.JWKSURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/consent", p.ConsentURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/login", p.LoginURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/logout", p.LogoutURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/token", p.OAuth2TokenURL(ctx).String())
	assert.Equal(t, "http://hydra.example/oauth2/auth", p.OAuth2AuthURL(ctx).String())
	assert.Equal(t, "http://hydra.example/userinfo", p.OIDCDiscoveryUserinfoEndpoint(ctx).String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(ctx, KeyIssuerURL, "http://hydra.localhost")
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
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l, configx.SkipValidation())
	p.MustSet(ctx, KeyCookieSameSiteMode, "")
	assert.Equal(t, http.SameSiteDefaultMode, p.CookieSameSiteMode(ctx))

	p.MustSet(ctx, KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode(ctx))

	p = MustNew(context.Background(), l, configx.SkipValidation())
	p.MustSet(ctx, "dangerous-force-http", true)
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode(ctx))
	p.MustSet(ctx, KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode(ctx))
}

func TestViperProviderValidates(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	c := MustNew(context.Background(), l, configx.WithConfigFiles("../../internal/.hydra.yaml"))

	// log
	assert.Equal(t, "debug", c.Source(ctx).String(KeyLogLevel))
	assert.Equal(t, "json", c.Source(ctx).String("log.format"))

	// serve
	assert.Equal(t, "localhost:1", c.ListenOn(ctx, PublicInterface))
	assert.Equal(t, "localhost:2", c.ListenOn(ctx, AdminInterface))

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
	assert.Equal(t, expectedPublicPermission, c.SocketPermission(ctx, PublicInterface))
	assert.Equal(t, expectedAdminPermission, c.SocketPermission(ctx, AdminInterface))

	expectedCors := cors.Options{
		AllowedOrigins:     []string{"https://example.com"},
		AllowedMethods:     []string{"GET"},
		AllowedHeaders:     []string{"Authorization"},
		ExposedHeaders:     []string{"Content-Type"},
		AllowCredentials:   true,
		MaxAge:             1,
		Debug:              false,
		OptionsPassthrough: true,
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
	assert.Contains(t, c.DSN(ctx), "sqlite://")

	// webfinger
	assert.Equal(t, []string{"hydra.openid.id-token"}, c.WellKnownKeys(ctx))
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

	// urls
	assert.Equal(t, urlx.ParseOrPanic("https://issuer/"), c.IssuerURL(ctx))
	assert.Equal(t, urlx.ParseOrPanic("https://public/"), c.PublicURL(ctx))
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

	// ttl
	assert.Equal(t, 2*time.Hour, c.ConsentRequestMaxAge(ctx))
	assert.Equal(t, 2*time.Hour, c.GetAccessTokenLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetRefreshTokenLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetIDTokenLifespan(ctx))
	assert.Equal(t, 2*time.Hour, c.GetAuthorizeCodeLifespan(ctx))

	// oauth2
	assert.Equal(t, true, c.GetSendDebugMessagesToClients(ctx))
	assert.Equal(t, true, c.GetUseLegacyErrorFormat(ctx))
	assert.Equal(t, 20, c.GetBCryptCost(ctx))
	assert.Equal(t, true, c.GetEnforcePKCE(ctx))
	assert.Equal(t, true, c.GetEnforcePKCEForPublicClients(ctx))

	// secrets
	assert.Equal(t, []byte{0x64, 0x40, 0x5f, 0xd4, 0x66, 0xc9, 0x8c, 0x88, 0xa7, 0xf2, 0xcb, 0x95, 0xcd, 0x95, 0xcb, 0xa3, 0x41, 0x49, 0x8b, 0x97, 0xba, 0x9e, 0x92, 0xee, 0x4c, 0xaf, 0xe0, 0x71, 0x23, 0x28, 0xeb, 0xfc}, c.GetGlobalSecret(ctx))
	assert.Equal(t, [][]uint8{[]byte("some-random-cookie-secret")}, c.GetCookieSecrets(ctx))

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
					ServerURL: "http://sampling",
				},
			},
		},
	}, c.Tracing(ctx))
}

func TestSetPerm(t *testing.T) {
	f, e := ioutil.TempFile("", "test")
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
	l.Logrus().SetOutput(ioutil.Discard)
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
	l.Logrus().SetOutput(ioutil.Discard)
	c := MustNew(context.Background(), l, configx.WithValue("ttl.refresh_token", -1))

	assert.Equal(t, -1*time.Nanosecond, c.GetRefreshTokenLifespan(ctx))
}
