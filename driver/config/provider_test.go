package config

import (
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
	return MustNew(logrusx.New("", ""))
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

	conf, enabled := p.PublicCORS()
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

	p := MustNew(l)

	value := p.AdminDisableHealthAccessLog()
	assert.Equal(t, false, value)

	p.MustSet(KeyAdminDisableHealthAccessLog, "true")

	value = p.AdminDisableHealthAccessLog()
	assert.Equal(t, true, value)
}

func TestProviderPublicDisableHealthAccessLog(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(l)

	value := p.PublicDisableHealthAccessLog()
	assert.Equal(t, false, value)

	p.MustSet(KeyPublicDisableHealthAccessLog, "true")

	value = p.PublicDisableHealthAccessLog()
	assert.Equal(t, true, value)
}

func TestProviderIssuerURL(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)
	p := MustNew(l)
	p.MustSet(KeyIssuerURL, "http://hydra.localhost")
	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL().String())

	p2 := MustNew(l)
	p2.MustSet(KeyIssuerURL, "http://hydra.localhost/")
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL().String())
}

func TestProviderCookieSameSiteMode(t *testing.T) {
	l := logrusx.New("", "")
	l.Logrus().SetOutput(ioutil.Discard)

	p := MustNew(l, configx.SkipValidation())
	p.MustSet("dangerous-force-http", false)
	p.MustSet(KeyCookieSameSiteMode, "")
	assert.Equal(t, http.SameSiteDefaultMode, p.CookieSameSiteMode())

	p.MustSet(KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode())

	p = MustNew(l, configx.SkipValidation())
	p.MustSet("dangerous-force-http", true)
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode())
	p.MustSet(KeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteLaxMode, p.CookieSameSiteMode())
}

func TestViperProviderValidates(t *testing.T) {
	l := logrusx.New("", "")
	c := MustNew(l, configx.WithConfigFiles("../../internal/.hydra.yaml"))

	// log
	assert.Equal(t, "debug", c.Source().String(KeyLogLevel))
	assert.Equal(t, "json", c.Source().String("log.format"))

	// serve
	assert.Equal(t, "localhost:1", c.PublicListenOn())
	assert.Equal(t, "localhost:2", c.AdminListenOn())

	expectedPublicPermission := &UnixPermission{
		Owner: "hydra",
		Group: "hydra-public-api",
		Mode:  0775,
	}
	expectedAdminPermission := &UnixPermission{
		Owner: "hydra",
		Group: "hydra-admin-api",
		Mode:  0770,
	}
	assert.Equal(t, expectedPublicPermission, c.PublicSocketPermission())
	assert.Equal(t, expectedAdminPermission, c.AdminSocketPermission())

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

	gc, enabled := c.AdminCORS()
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	gc, enabled = c.PublicCORS()
	assert.False(t, enabled)
	assert.Equal(t, expectedCors, gc)

	assert.Equal(t, []string{"127.0.0.1/32"}, c.AllowTLSTerminationFrom())
	assert.Equal(t, "/path/to/file.pem", c.Source().String("serve.tls.key.path"))
	assert.Equal(t, "b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK", c.Source().String("serve.tls.cert.base64"))
	assert.Equal(t, http.SameSiteLaxMode, c.CookieSameSiteMode())
	assert.Equal(t, true, c.CookieSameSiteLegacyWorkaround())

	// dsn
	assert.Equal(t, dbal.InMemoryDSN, c.DSN())

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
		Jaeger: &tracing.JaegerConfig{
			LocalAgentHostPort: "127.0.0.1:6831",
			SamplerType:        "const",
			SamplerValue:       1,
			SamplerServerURL:   "http://sampling",
			Propagation:        "jaeger",
		},
		Zipkin: &tracing.ZipkinConfig{
			ServerURL: "http://zipkin/api/v2/spans",
		},
	}, c.Tracing())
}

func TestSetPerm(t *testing.T) {
	f, e := ioutil.TempFile("", "test")
	require.NoError(t, e)
	path := f.Name()

	// We cannot test setting owner and group, because we don't know what the
	// tester has access to.
	_ = (&UnixPermission{
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
