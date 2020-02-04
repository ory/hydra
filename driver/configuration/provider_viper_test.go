package configuration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	"github.com/ory/hydra/x"
	"github.com/ory/x/logrusx"
)

func setEnv(key, value string) func(t *testing.T) {
	return func(t *testing.T) {
		require.NoError(t, os.Setenv(key, value))
	}
}

func TestSubjectTypesSupported(t *testing.T) {
	p := NewViperProvider(logrus.New(), false, nil)
	viper.Set(ViperKeySubjectIdentifierAlgorithmSalt, "00000000")
	for k, tc := range []struct {
		d string
		p func(t *testing.T)
		e []string
		c func(t *testing.T)
	}{
		{
			d: "Load legacy environment variable in legacy format",
			p: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), "public,pairwise,foobar"),
			c: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), ""),
			e: []string{"public", "pairwise"},
		},
		{
			d: "Load legacy environment variable in legacy format",
			p: func(t *testing.T) {
				setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), "public,pairwise,foobar")(t)
				setEnv(strings.ToUpper(strings.Replace(ViperKeyAccessTokenStrategy, ".", "_", -1)), "jwt")(t)
			},
			c: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), ""),
			e: []string{"public"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			tc.p(t)
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported())
			tc.c(t)
		})
	}
}

func TestWellKnownKeysUnique(t *testing.T) {
	p := NewViperProvider(logrus.New(), false, nil)
	assert.EqualValues(t, []string{x.OAuth2JWTKeyName, x.OpenIDConnectKeyName}, p.WellKnownKeys(x.OAuth2JWTKeyName, x.OpenIDConnectKeyName, x.OpenIDConnectKeyName))
}

func TestCORSOptions(t *testing.T) {
	p := NewViperProvider(logrus.New(), false, nil)
	viper.Set("serve.public.cors.enabled", true)

	assert.EqualValues(t, cors.Options{
		AllowedOrigins:     []string{},
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:     []string{"Authorization", "Content-Type"},
		ExposedHeaders:     []string{"Content-Type"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		MaxAge:             0,
		Debug:              false,
	}, p.CORSOptions("public"))
}

func TestViperProvider_AdminDisableHealthAccessLog(t *testing.T) {
	l := logrusx.New()
	l.SetOutput(ioutil.Discard)

	p := NewViperProvider(l, false, nil)

	value := p.AdminDisableHealthAccessLog()
	assert.Equal(t, false, value)

	os.Setenv("SERVE_ADMIN_ACCESS_LOG_DISABLE_FOR_HEALTH", "true")

	value = p.AdminDisableHealthAccessLog()
	assert.Equal(t, true, value)
}

func TestViperProvider_PublicDisableHealthAccessLog(t *testing.T) {
	l := logrusx.New()
	l.SetOutput(ioutil.Discard)

	p := NewViperProvider(l, false, nil)

	value := p.PublicDisableHealthAccessLog()
	assert.Equal(t, false, value)

	os.Setenv("SERVE_PUBLIC_ACCESS_LOG_DISABLE_FOR_HEALTH", "true")

	value = p.PublicDisableHealthAccessLog()
	assert.Equal(t, true, value)
}

func TestViperProvider_IssuerURL(t *testing.T) {
	l := logrusx.New()
	l.SetOutput(ioutil.Discard)
	viper.Set(ViperKeyIssuerURL, "http://hydra.localhost")
	p := NewViperProvider(l, false, nil)
	assert.Equal(t, "http://hydra.localhost/", p.IssuerURL().String())

	viper.Set(ViperKeyIssuerURL, "http://hydra.localhost/")
	p2 := NewViperProvider(l, false, nil)
	assert.Equal(t, "http://hydra.localhost/", p2.IssuerURL().String())
}

func TestViperProvider_CookieSameSiteMode(t *testing.T) {
	l := logrusx.New()
	l.SetOutput(ioutil.Discard)

	p := NewViperProvider(l, false, nil)
	assert.Equal(t, http.SameSiteDefaultMode, p.CookieSameSiteMode())

	os.Setenv("COOKIE_SAME_SITE_MODE", "none")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode())
}
