package configuration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ory/hydra/x"
	"github.com/ory/viper"
	"github.com/ory/x/logrusx"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)): "public,pairwise,foobar",
			}),
			e: []string{"public", "pairwise"},
		},
		{
			d: "Load legacy environment variable in legacy format",
			env: setupEnv(map[string]string{
				strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)): "public,pairwise,foobar",
				strings.ToUpper(strings.Replace(ViperKeyAccessTokenStrategy, ".", "_", -1)):   "jwt",
			}),
			e: []string{"public"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			setup, clean := tc.env(t)
			setup()
			p := NewViperProvider(logrus.New(), false, nil)
			viper.Set(ViperKeySubjectIdentifierAlgorithmSalt, "00000000")
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported())
			clean()
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

	viper.Set(ViperKeyAdminDisableHealthAccessLog, "true")

	value = p.AdminDisableHealthAccessLog()
	assert.Equal(t, true, value)
}

func TestViperProvider_PublicDisableHealthAccessLog(t *testing.T) {
	l := logrusx.New()
	l.SetOutput(ioutil.Discard)

	p := NewViperProvider(l, false, nil)

	value := p.PublicDisableHealthAccessLog()
	assert.Equal(t, false, value)

	viper.Set(ViperKeyPublicDisableHealthAccessLog, "true")

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

	viper.Set(ViperKeyCookieSameSiteMode, "none")
	assert.Equal(t, http.SameSiteNoneMode, p.CookieSameSiteMode())
}
