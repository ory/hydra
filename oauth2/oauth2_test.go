package oauth2

import (
	"testing"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/fosite/hash"
	"net/url"
	"net/http/httptest"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"net/http"
	"github.com/ory-am/hydra/key"
)

var store = pkg.FositeStore()

var handler = &Handler{
	OAuth2: &fosite.Fosite{
		Store: store,
		MandatoryScope       :       "hydra",
		AuthorizeEndpointHandlers  : fosite.AuthorizeEndpointHandlers{},
		TokenEndpointHandlers      : fosite.TokenEndpointHandlers{},
		AuthorizedRequestValidators: fosite.AuthorizedRequestValidators{},
		Hasher                 :     &hash.BCrypt{},
	},
	Consent  :  &DefaultConsentStrategy{
		Issuer: "https://hydra.localhost",
		KeyManager: &key.MemoryManager{
			AsymmetricKeys: map[string]*key.AsymmetricKey{},
			SymmetricKeys:  map[string]*key.SymmetricKey{},
			Strategy: &key.DefaultKeyStrategy{
				AsymmetricKeyStrategy: &key.RSAPEMStrategy{},
				SymmetricKeyStrategy:  &key.SHAStrategy{},
			},
		},
		PrivateKey: []byte{},
		PublicKey: []byte{},
	},
}

var r = httprouter.New()

var ts *httptest.Server

func init() {
	ts = httptest.NewServer(r)
	handler.SetRoutes(r)

	store.Clients["app"] = &fosite.DefaultClient{
		ID: "app",
		Secret: []byte("secret"),
		RedirectURIs: []string{ts.URL + "/callback"},
	}

	s, _ := url.Parse(ts.URL)
	handler.SelfURL = s

	c, _ := url.Parse(ts.URL + "/consent")
	handler.ConsentURL = c
}

func TestAuthCode(t *testing.T) {
	c := oauth2.Config{
		ClientID: "",
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			AuthURL: ts.URL + "/oauth2/auth",
			TokenURL:ts.URL + "/oauth2/token",
		},
		RedirectURL:ts.URL +  "/callback",
		Scopes :[]string{},
	}

	var token string
	r.GET("/consent", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	})
	r.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	})

	var consent Session
	err := pkg.NewSuperAgent(c.AuthCodeURL("some-foo-state")).GET(&consent)
	pkg.RequireError(t, false, err)

	resp, err := http.Get(c.AuthCodeURL("some-foo-state") + "&consent="+token)
	pkg.RequireError(t, false, err)
	defer resp.Body.Close()
}