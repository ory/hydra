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
	"github.com/ory/x/dbal"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/tracing"
	"github.com/ory/x/urlx"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/x"
)

func newProvider() *Provider {
	return MustNew(context.Background(), logrusx.New("", ""))
}

func setupEnv(env map[string]string) func(t *testing.T) (func(), func()) {
	return func(t *testing.T) (setup func(), clean func()) {
		setup = func() {
			for k, v := range env {
				require.NoError(t, os.Setenv(k, v))
			}
		}

		clean = func() {
			for k := range env {
				require.NoError(t, os.Unsetenv(k))
			}
		}
		return
	}
}

func TestSubjectTypesSupported(t *testing.T) {
	for k, tc := range []struct {
		d   string
		env func(t *testing.T) (func(), func())
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
			d: "Load legacy environment variable in legacy format",
			env: setupEnv(map[string]string{
				strings.ToUpper(strings.Replace(KeySubjectTypesSupported, ".", "_", -1)):                 "public,pairwise",
				strings.ToUpper(strings.Replace("oidc.subject_identifiers.pairwise.salt", ".", "_", -1)): "some-salt",
				strings.ToUpper(strings.Replace(KeyAccessTokenStrategy, ".", "_", -1)):                   "jwt",
			}),
			e: []string{"public"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			setup, clean := tc.env(t)
			setup()
			p := newProvider()
			p.MustSet(KeySubjectIdentifierAlgorithmSalt, "00000000")
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported())
			clean()
		})
	}
}

func TestWellKnownKeysUnique(t *testing.T) {
	p := newProvider()
	assert.EqualValues(t, []string{x.OpenIDConnectKeyName, x.OAuth2JWTKeyName}, p.WellKnownKeys(x.OAuth2JWTKeyName, x.OpenIDConnectKeyName, x.OpenIDConnectKeyName))
}

func TestCORSOptions(t *testing.T) {
	p := newProvider()
	p.MustSet("serve.public.cors.enabled", true)

	conf, enabled := p.CORS(PublicInterface)
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
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l)

	value := p.DisableHealthAccessLog(AdminInterface)
	assert.Equal(t, false, value)

	p.MustSet(AdminInterface.Key(KeySuffixDisableHealthAccessLog), "true")

	value = p.DisableHealthAccessLog(AdminInterface)
	assert.Equal(t, true, value)
}

func TestProviderPublicDisableHealthAccessLog(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l)

	value := p.DisableHealthAccessLog(PublicInterface)
	assert.Equal(t, false, value)

	p.MustSet(PublicInterface.Key(KeySuffixDisableHealthAccessLog), "true")

	value = p.DisableHealthAccessLog(PublicInterface)
	assert.Equal(t, true, value)
}

func TestPublicAllowDynamicRegistration(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l)

	value := p.PublicAllowDynamicRegistration()
	assert.Equal(t, false, value)

	p.MustSet(KeyPublicAllowDynamicRegistration, "true")

	value = p.PublicAllowDynamicRegistration()
	assert.Equal(t, true, value)
}

func TestProviderIssuerURL(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(KeyIssuerURL, "http://hydra.localhost")
	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL().String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(KeyIssuerURL, "http://hydra.localhost/")
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL().String())
}

func TestProviderIssuerPublicURL(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(KeyIssuerURL, "http://hydra.localhost")
	p.MustSet(KeyPublicURL, "http://hydra.example")

	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL().String())
	assert.Equal(t, "http://hydra.example/", p.PublicURL().String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p.JWKSURL().String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/consent", p.ConsentURL().String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/login", p.LoginURL().String())
	assert.Equal(t, "http://hydra.example/oauth2/fallbacks/logout", p.LogoutURL().String())
	assert.Equal(t, "http://hydra.example/oauth2/token", p.OAuth2TokenURL().String())
	assert.Equal(t, "http://hydra.example/oauth2/auth", p.OAuth2AuthURL().String())
	assert.Equal(t, "http://hydra.example/userinfo", p.OIDCDiscoveryUserinfoEndpoint().String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(KeyIssuerURL, "http://hydra.localhost")
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL().String())
	assert.Equal(t, "http://hydra.localhost/", p2.PublicURL().String())
	assert.Equal(t, "http://hydra.localhost/.well-known/jwks.json", p2.JWKSURL().String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/consent", p2.ConsentURL().String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/login", p2.LoginURL().String())
	assert.Equal(t, "http://hydra.localhost/oauth2/fallbacks/logout", p2.LogoutURL().String())
	assert.Equal(t, "http://hydra.localhost/oauth2/token", p2.OAuth2TokenURL().String())
	assert.Equal(t, "http://hydra.localhost/oauth2/auth", p2.OAuth2AuthURL().String())
	assert.Equal(t, "http://hydra.localhost/userinfo", p2.OIDCDiscoveryUserinfoEndpoint().String())
}

func TestProviderCookieSameSiteMode(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(context.Background(), l, configx.SkipValidation())
	p.MustSet(KeyCookieSameSiteMode, "")
	assert.Equal(t, http.SameSiteDefaultMode, p.CookieSameSiteMode())

	p.MustSet(KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode())

	p = MustNew(context.Background(), l, configx.SkipValidation())
	p.MustSet("dangerous-force-http", true)
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode())
	p.MustSet(KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode())
}

func TestViperProviderValidates(t *testing.T) {
	l := logrusx.New("", "")
	c := MustNew(context.Background(), l, configx.WithConfigFiles("../../internal/.hydra.yaml"))

	// log
	assert.Equal(t, "debug", c.Source().String(KeyLogLevel))
	assert.Equal(t, "json", c.Source().String("log.format"))

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
		AllowedOrigins:     []string{"https://example.com"},
		AllowedMethods:     []string{"GET"},
		AllowedHeaders:     []string{"Authorization"},
		ExposedHeaders:     []string{"Content-Type"},
		AllowCredentials:   true,
		MaxAge:             1,
		Debug:              false,
		OptionsPassthrough: true,
	}

	gc, enabled := c.CORS(AdminInterface)
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	gc, enabled = c.CORS(PublicInterface)
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	assert.Equal(t, []string{"127.0.0.1/32"}, c.TLS(PublicInterface).AllowTerminationFrom())
	assert.Equal(t, "/path/to/file.pem", c.Source().String("serve.tls.key.path"))
	assert.Equal(t, "b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK", c.Source().String("serve.tls.cert.base64"))
	assert.Equal(t, http.SameSiteLaxMode, c.CookieSameSiteMode())
	assert.Equal(t, true, c.CookieSameSiteLegacyWorkaround())

	// dsn
	assert.Equal(t, dbal.SQLiteInMemory, c.DSN())

	// webfinger
	assert.Equal(t, []string{"hydra.openid.id-token"}, c.WellKnownKeys())
	assert.Equal(t, urlx.ParseOrPanic("https://example.com"), c.OAuth2ClientRegistrationURL())
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/jwks.json"), c.JWKSURL())
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/auth"), c.OAuth2AuthURL())
	assert.Equal(t, urlx.ParseOrPanic("https://example.com/token"), c.OAuth2TokenURL())
	assert.Equal(t, []string{"sub", "username"}, c.OIDCDiscoverySupportedClaims())
	assert.Equal(t, []string{"offline_access", "offline", "openid", "whatever"}, c.OIDCDiscoverySupportedScope())
	assert.Equal(t, urlx.ParseOrPanic("https://example.com"), c.OIDCDiscoveryUserinfoEndpoint())

	// oidc
	assert.Equal(t, []string{"pairwise"}, c.SubjectTypesSupported())
	assert.Equal(t, "random_salt", c.SubjectIdentifierAlgorithmSalt())
	assert.Equal(t, []string{"whatever"}, c.DefaultClientScope())

	// urls
	assert.Equal(t, urlx.ParseOrPanic("https://issuer/"), c.IssuerURL())
	assert.Equal(t, urlx.ParseOrPanic("https://public/"), c.PublicURL())
	assert.Equal(t, urlx.ParseOrPanic("https://login/"), c.LoginURL())
	assert.Equal(t, urlx.ParseOrPanic("https://consent/"), c.ConsentURL())
	assert.Equal(t, urlx.ParseOrPanic("https://logout/"), c.LogoutURL())
	assert.Equal(t, urlx.ParseOrPanic("https://error/"), c.ErrorURL())
	assert.Equal(t, urlx.ParseOrPanic("https://post_logout/"), c.LogoutRedirectURL())

	// strategies
	assert.Equal(t, "exact", c.ScopeStrategy())
	assert.Equal(t, "opaque", c.AccessTokenStrategy())
	assert.Equal(t, false, c.GrantAllClientCredentialsScopesPerDefault())

	// ttl
	assert.Equal(t, 2*time.Hour, c.ConsentRequestMaxAge())
	assert.Equal(t, 2*time.Hour, c.AccessTokenLifespan())
	assert.Equal(t, 2*time.Hour, c.RefreshTokenLifespan())
	assert.Equal(t, 2*time.Hour, c.IDTokenLifespan())
	assert.Equal(t, 2*time.Hour, c.AuthCodeLifespan())

	// oauth2
	assert.Equal(t, true, c.ShareOAuth2Debug())
	assert.Equal(t, true, c.OAuth2LegacyErrors())
	assert.Equal(t, 20, c.BCryptCost())
	assert.Equal(t, true, c.PKCEEnforced())
	assert.Equal(t, true, c.EnforcePKCEForPublicClients())

	// secrets
	assert.Equal(t, []byte{0x64, 0x40, 0x5f, 0xd4, 0x66, 0xc9, 0x8c, 0x88, 0xa7, 0xf2, 0xcb, 0x95, 0xcd, 0x95, 0xcb, 0xa3, 0x41, 0x49, 0x8b, 0x97, 0xba, 0x9e, 0x92, 0xee, 0x4c, 0xaf, 0xe0, 0x71, 0x23, 0x28, 0xeb, 0xfc}, c.GetSystemSecret())
	assert.Equal(t, [][]uint8{[]byte("some-random-cookie-secret")}, c.GetCookieSecrets())

	// profiling
	assert.Equal(t, "cpu", c.Source().String("profiling"))

	// tracing
	assert.EqualValues(t, &tracing.Config{
		ServiceName: "hydra service",
		Provider:    "jaeger",
		Providers: &tracing.ProvidersConfig{
			Jaeger: &tracing.JaegerConfig{
				LocalAgentAddress: "127.0.0.1:6831",
				Sampling: &tracing.JaegerSampling{
					Type:      "const",
					Value:     1,
					ServerURL: "http://sampling",
				},
				Propagation:       "jaeger",
				MaxTagValueLength: 1024,
			},
			Zipkin: &tracing.ZipkinConfig{
				ServerURL: "http://zipkin/api/v2/spans",
			},
		},
	}, c.Tracing())
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
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)
	p := MustNew(context.Background(), l)
	p.MustSet(KeyLoginURL, "http://localhost:8080/oauth/login")
	p.MustSet(KeyConsentURL, "http://localhost:8080/oauth/consent")

	assert.Equal(t, "http://localhost:8080/oauth/login", p.LoginURL().String())
	assert.Equal(t, "http://localhost:8080/oauth/consent", p.ConsentURL().String())

	p2 := MustNew(context.Background(), l)
	p2.MustSet(KeyLoginURL, "http://localhost:3000/#/oauth/login")
	p2.MustSet(KeyConsentURL, "http://localhost:3000/#/oauth/consent")

	assert.Equal(t, "http://localhost:3000/#/oauth/login", p2.LoginURL().String())
	assert.Equal(t, "http://localhost:3000/#/oauth/consent", p2.ConsentURL().String())
}

func TestInfinitRefreshTokenTTL(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)
	c := MustNew(context.Background(), l, configx.WithValue("ttl.refresh_token", -1))

	assert.Equal(t, -1*time.Nanosecond, c.RefreshTokenLifespan())
}
