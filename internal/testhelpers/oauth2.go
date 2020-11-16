package testhelpers

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ory/fosite/token/jwt"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"

	"github.com/ory/x/httpx"
	"github.com/ory/x/ioutilx"

	"github.com/gobuffalo/httptest"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/viper"
)

func NewIDToken(t *testing.T, reg driver.Registry, subject string) string {
	token, _, err := reg.OpenIDJWTStrategy().Generate(context.TODO(), jwt.IDTokenClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(time.Hour),
		IssuedAt:  time.Now(),
	}.ToMapClaims(), jwt.NewHeaders())
	require.NoError(t, err)
	return token
}

func NewOAuth2Server(t *testing.T, reg driver.Registry) (publicTS, adminTS *httptest.Server) {
	// Lifespan is two seconds to avoid time synchronization issues with SQL.
	viper.Set(configuration.ViperKeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
	viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Second*2)
	viper.Set(configuration.ViperKeyRefreshTokenLifespan, time.Second*3)
	viper.Set(configuration.ViperKeyScopeStrategy, "exact")

	public, admin := x.NewRouterPublic(), x.NewRouterAdmin()

	publicTS = httptest.NewServer(public)
	t.Cleanup(publicTS.Close)

	adminTS = httptest.NewServer(admin)
	t.Cleanup(adminTS.Close)

	viper.Set(configuration.ViperKeyIssuerURL, publicTS.URL)
	// SendDebugMessagesToClients: true,

	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)
	internal.MustEnsureRegistryKeys(reg, x.OAuth2JWTKeyName)

	reg.RegisterRoutes(admin, public)
	return publicTS, adminTS
}

func IntrospectToken(t *testing.T, conf *oauth2.Config, token *oauth2.Token, adminTS *httptest.Server) gjson.Result {
	require.NotEmpty(t, token.AccessToken)

	req := httpx.MustNewRequest("POST", adminTS.URL+"/oauth2/introspect",
		strings.NewReader((url.Values{"token": {token.AccessToken}}).Encode()),
		"application/x-www-form-urlencoded")

	req.SetBasicAuth(conf.ClientID, conf.ClientSecret)
	res, err := adminTS.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	return gjson.ParseBytes(ioutilx.MustReadAll(res.Body))
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

func NewUI(t *testing.T, login, consent http.HandlerFunc) {
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

	viper.Set(configuration.ViperKeyLoginURL, lt.URL)
	viper.Set(configuration.ViperKeyConsentURL, ct.URL)
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
