package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerWellKnown(t *testing.T) {
	h := &Handler{
		H:      herodot.NewJSONWriter(nil),
		Issuer: "http://hydra.localhost",
	}

	AuthPathT := "/oauth2/auth"
	TokenPathT := "/oauth2/token"
	JWKPathT := "/.well-known/jwks.json"

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + "/.well-known/openid-configuration")

	defer res.Body.Close()

	trueConfig := WellKnown{
		Issuer:        h.Issuer,
		AuthURL:       h.Issuer + AuthPathT,
		TokenURL:      h.Issuer + TokenPathT,
		JWKsURI:       h.Issuer + JWKPathT,
		SubjectTypes:  []string{"pairwise", "public"},
		SigningAlgs:   []string{"RS256"},
		ResponseTypes: []string{"code", "code id_token", "id_token", "token id_token", "token"},
	}
	var wellKnownResp WellKnown
	err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
	require.NoError(t, err, "problem decoding wellknown json response: %+v", err)
	assert.Equal(t, trueConfig, wellKnownResp)
}

type FakeConsentStrategy struct {
	RedirectURL string
}

func (s *FakeConsentStrategy) ValidateResponse(authorizeRequest fosite.AuthorizeRequester, token string, session *sessions.Session) (claims *Session, err error) {
	return nil, nil
}

func (s *FakeConsentStrategy) IssueChallenge(authorizeRequest fosite.AuthorizeRequester, redirectURL string, session *sessions.Session) (token string, err error) {
	s.RedirectURL = redirectURL
	return "token", nil
}

func TestIssuerRedirect(t *testing.T) {
	storage := storage.NewExampleStore()
	secret := []byte("my super secret password")
	config := compose.Config{}
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	consentUrl, _ := url.Parse("http://consent.localhost")

	cs := &FakeConsentStrategy{}

	h := &Handler{
		H:           herodot.NewJSONWriter(nil),
		Issuer:      "http://127.0.0.1/some/proxied/path",
		OAuth2:      compose.ComposeAllEnabled(&config, storage, secret, privateKey),
		ConsentURL:  *consentUrl,
		CookieStore: sessions.NewCookieStore([]byte("my super secret password")),
		Consent:     cs,
		L:           logrus.New(),
	}

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	authUrl, _ := url.Parse(ts.URL)
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", "my-client")
	v.Set("redirect_uri", "http://localhost:3846/callback")
	v.Set("scope", "openid")
	v.Set("state", "my super secret state")
	authUrl.Path = "/oauth2/auth"
	authUrl.RawQuery = v.Encode()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, _ := client.Get(authUrl.String())

	authRedirect, _ := url.Parse(cs.RedirectURL)
	assert.Equal(t, "/some/proxied/path/oauth2/auth", authRedirect.Path, "The redirect URL sent in the challenge includes the full issuer path")
	assert.Equal(t, authUrl.Query(), authRedirect.Query(), "The auth redirect should have the same parameters with the addition of challenge")

	defer res.Body.Close()
}
