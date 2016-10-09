package oauth2_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	ejwt "github.com/ory-am/fosite/token/jwt"
	"github.com/ory-am/hydra/jwk"
	. "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
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

		consent, err := signConsentToken(map[string]interface{}{
			"jti": uuid.New(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Unix(),
			"aud": "app-client",
			"scp": []string{"hydra"},
		})
		pkg.RequireError(t, false, err)

		http.Redirect(w, r, ejwt.ToString(jwtClaims["redir"])+"&consent="+consent, http.StatusFound)
		validConsent = true
	})

	router.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		code = r.URL.Query().Get("code")
		w.Write([]byte(r.URL.Query().Get("code")))
	})

	resp, err := http.Get(oauthConfig.AuthCodeURL("some-foo-state"))
	pkg.RequireError(t, false, err)
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	pkg.RequireError(t, false, err)

	require.True(t, validConsent)
	require.NotEmpty(t, code)

	_, err = oauthConfig.Exchange(oauth2.NoContext, code)
	pkg.RequireError(t, false, err, code)
}
