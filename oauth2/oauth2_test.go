package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/enigma/hmac"
	ejwt "github.com/ory-am/fosite/enigma/jwt"
	"github.com/ory-am/fosite/handler/core/explicit"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/hydra/key"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

var store = pkg.FositeStore()

var keyStrategy = &key.DefaultKeyStrategy{
	AsymmetricKeyStrategy: &key.RSAPEMStrategy{},
	SymmetricKeyStrategy:  &key.SHAStrategy{},
}

var keyManager = &key.MemoryManager{
	AsymmetricKeys: map[string]*key.AsymmetricKey{},
	SymmetricKeys:  map[string]*key.SymmetricKey{},
	Strategy:       keyStrategy,
}

var hmacStrategy = &strategy.HMACSHAStrategy{
	Enigma: &hmac.Enigma{
		GlobalSecret: []byte("some-super-cool-secret-that-nobody-knows"),
	},
}

var authCodeHandler = &explicit.AuthorizeExplicitGrantTypeHandler{
	AccessTokenStrategy:       hmacStrategy,
	RefreshTokenStrategy:      hmacStrategy,
	AuthorizeCodeStrategy:     hmacStrategy,
	AuthorizeCodeGrantStorage: store,
	AuthCodeLifespan:          time.Hour,
	AccessTokenLifespan:       time.Hour,
}

var hasher = &hash.BCrypt{}

var handler = &Handler{
	OAuth2: &fosite.Fosite{
		Store:          store,
		MandatoryScope: "hydra",
		AuthorizeEndpointHandlers: fosite.AuthorizeEndpointHandlers{
			authCodeHandler,
		},
		TokenEndpointHandlers: fosite.TokenEndpointHandlers{
			authCodeHandler,
		},
		AuthorizedRequestValidators: fosite.AuthorizedRequestValidators{},
		Hasher: hasher,
	},
	Consent: &DefaultConsentStrategy{
		Issuer:     "https://hydra.localhost",
		KeyManager: keyManager,
	},
}

var r = httprouter.New()

var ts *httptest.Server

func init() {
	keyManager.CreateAsymmetricKey(ConsentChallengeKey)
	keyManager.CreateAsymmetricKey(ConsentEndpointKey)
	ts = httptest.NewServer(r)

	handler.SetRoutes(r)
	store.Clients["app"] = &fosite.DefaultClient{
		ID:           "app",
		Secret:       []byte("secret"),
		RedirectURIs: []string{ts.URL + "/callback"},
	}

	s, _ := url.Parse(ts.URL)
	handler.SelfURL = *s

	c, _ := url.Parse(ts.URL + "/consent")
	handler.ConsentURL = *c

	h, _ := hasher.Hash( []byte("secret"))

	store.Clients["app-client"] = &fosite.DefaultClient{
		ID:           "app-client",
		Secret:       h,
		RedirectURIs: []string{ts.URL + "/callback"},
	}
}

func TestAuthCode(t *testing.T) {
	t.Log(ts.URL)
	c := oauth2.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/oauth2/auth",
			TokenURL: ts.URL + "/oauth2/token",
		},
		RedirectURL: ts.URL + "/callback",
		Scopes:      []string{"hydra"},
	}

	var code string
	var validConsent bool
	r.GET("/consent", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t.Logf("Consent request at %s", r.URL)
		tok, err := jwt.Parse(r.URL.Query().Get("challenge"), func(tt *jwt.Token) (interface{}, error) {
			if _, ok := tt.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.Errorf("Unexpected signing method: %v", tt.Header["alg"])
			}

			pk, err := keyManager.GetAsymmetricKey(ConsentChallengeKey)
			pkg.RequireError(t, false, err)
			return jwt.ParseRSAPublicKeyFromPEM(pk.Public)
		})
		pkg.RequireError(t, false, err)
		require.True(t, tok.Valid)
		t.Logf("Consent request at %v", tok.Claims)

		validConsent = true
		consent, err := SignToken(map[string]interface{}{
			"jti": uuid.New(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Unix(),
			"aud": "app-client",
		})
		pkg.RequireError(t, false, err)
		http.Redirect(w, r, ejwt.ToString(tok.Claims["redir"])+"&consent="+consent, http.StatusFound)
	})

	r.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t.Logf("Callback request at %s", r.URL)

		code = r.URL.Query().Get("code")
		w.Write([]byte(r.URL.Query().Get("code")))
	})

	resp, err := http.Get(c.AuthCodeURL("some-foo-state"))
	pkg.RequireError(t, false, err)
	defer resp.Body.Close()

	out, _ := ioutil.ReadAll(resp.Body)
	t.Logf("Got response: %s", out)

	require.True(t, validConsent)
	require.NotEmpty(t, code)

	ot, err := c.Exchange(oauth2.NoContext, code)
	pkg.RequireError(t, false, err)

	t.Logf("OAuth2 Token: %v", ot)
}

func SignToken(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = claims

	key, err := keyManager.GetAsymmetricKey(ConsentEndpointKey)
	if err != nil {
		return "", errors.New(err)
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(key.Private)
	if err != nil {
		return "", errors.New(err)
	}

	var signature, encoded string
	if encoded, err = token.SigningString(); err != nil {
		return "", errors.New(err)
	} else if signature, err = token.Method.Sign(encoded, rsaKey); err != nil {
		return "", errors.New(err)
	}

	return fmt.Sprintf("%s.%s", encoded, signature), nil
}
