package oauth2_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"net/http/cookiejar"

	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	ejwt "github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/jwk"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestAuthCode(t *testing.T) {
	var code string
	var validConsent bool
	router.GET("/consent", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tok, err := jwt.Parse(r.URL.Query().Get("challenge"), func(tt *jwt.Token) (interface{}, error) {
			if _, ok := tt.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.Errorf("Unexpected signing method: %v", tt.Header["alg"])
			}

			pk, err := keyManager.GetKey(ConsentChallengeKey, "public")
			pkg.RequireError(t, false, err)
			return jwk.MustRSAPublic(jwk.First(pk.Keys)), nil
		})
		pkg.RequireError(t, false, err)
		require.True(t, tok.Valid)

		jwtClaims, ok := tok.Claims.(jwt.MapClaims)
		require.True(t, ok)
		require.NotEmpty(t, jwtClaims)

		expl := map[string]interface{}{"foo": "bar", "baz": map[string]interface{}{"foo": "baz"}}
		consent, err := signConsentToken(map[string]interface{}{
			"jti":    jwtClaims["jti"],
			"exp":    time.Now().Add(time.Hour).Unix(),
			"iat":    time.Now().Unix(),
			"sub":    "foo",
			"aud":    "app-client",
			"scp":    []string{"hydra", "offline"},
			"at_ext": expl,
			"id_ext": expl,
		})
		pkg.RequireError(t, false, err)

		http.Redirect(w, r, ejwt.ToString(jwtClaims["redir"])+"&consent="+consent, http.StatusFound)
		validConsent = true
	})

	router.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		code = r.URL.Query().Get("code")
		w.Write([]byte(r.URL.Query().Get("code")))
	})

	cookieJar, _ := cookiejar.New(nil)
	req, err := http.NewRequest("GET", oauthConfig.AuthCodeURL("some-foo-state"), nil)
	pkg.RequireError(t, false, err)

	resp, err := (&http.Client{Jar: cookieJar}).Do(req)
	pkg.RequireError(t, false, err)
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	pkg.RequireError(t, false, err)

	require.True(t, validConsent)
	require.NotEmpty(t, code)

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	pkg.RequireError(t, false, err, code)

	time.Sleep(time.Second * 5)

	res, err := testRefresh(t, token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func testRefresh(t *testing.T, token *oauth2.Token) (*http.Response, error) {
	req, err := http.NewRequest("POST", oauthClientConfig.TokenURL, strings.NewReader(url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{token.RefreshToken},
	}.Encode()))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(oauthClientConfig.ClientID, oauthClientConfig.ClientSecret)

	return http.DefaultClient.Do(req)
}
