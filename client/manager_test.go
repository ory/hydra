package client_test

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	. "github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/herodot"
	ioa2 "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var clientManagers = map[string]ClientStorage{}

var fositeStore = pkg.FositeStore()

var ladonWarden = pkg.LadonWarden(map[string]ladon.Policy{
	"1": &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:clients<.*>"},
		Actions:   []string{"create", "get", "delete"},
		Effect:    ladon.AllowAccess,
	},
})

var localWarden = &warden.LocalWarden{
	Warden: ladonWarden,
	TokenValidator: &core.CoreValidator{
		AccessTokenStrategy: pkg.HMACStrategy,
		AccessTokenStorage:  fositeStore,
	},
	Issuer: "tests",
}

var ts *httptest.Server

var tokens = pkg.Tokens(1)

var httpClientManager *HTTPClientManager

func init() {
	ar := fosite.NewAccessRequest(&ioa2.Session{Subject: "alice"})
	ar.GrantedScopes = fosite.Arguments{Scope}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	clientManagers["memory"] = &MemoryClientManager{Clients: map[string]*Client{}}

	s := &ClientHandler{
		Manager: &MemoryClientManager{Clients: map[string]*Client{}},
		H:       &herodot.JSON{},
		W:       localWarden,
	}
	r := httprouter.New()
	s.SetRoutes(r)
	ts = httptest.NewServer(r)
	conf := &oauth2.Config{Scopes: []string{}, Endpoint: oauth2.Endpoint{}}

	u, _ := url.Parse(ts.URL + ClientsHandlerPath)
	clientManagers["http"] = &HTTPClientManager{
		Client: conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		}),
		Endpoint: u,
	}
}

func TestAuthenticateClient(t *testing.T) {
	var mem = &MemoryClientManager{Clients: map[string]*Client{}}
	mem.CreateClient(&Client{
		ID:           "1234",
		Secret:       []byte("secret"),
		RedirectURIs: []string{"http://redirect"},
	})

	c, err := mem.Authenticate("1234", []byte("secret1"))
	pkg.AssertError(t, true, err)

	c, err = mem.Authenticate("1234", []byte("secret"))
	pkg.AssertError(t, false, err)
	assert.Equal(t, "1234", c.ID)
}

func TestCreateGetDeleteClient(t *testing.T) {
	for k, m := range clientManagers {
		_, err := m.GetClient("4321")
		pkg.AssertError(t, true, err, k)

		c := &Client{
			ID:           "1234",
			Secret:       []byte("secret"),
			RedirectURIs: []string{"http://redirect"},
		}
		err = m.CreateClient(c)
		pkg.AssertError(t, false, err, k)
		if err == nil {
			compare(t, c, k)
		}

		d, err := m.GetClient("1234")
		pkg.AssertError(t, false, err, k)
		if err == nil {
			compare(t, d, k)
		}

		err = m.DeleteClient("1234")
		pkg.AssertError(t, false, err, k)

		_, err = m.GetClient("1234")
		pkg.AssertError(t, true, err, k)
	}
}

func compare(t *testing.T, c fosite.Client, k string) {
	assert.Equal(t, c.GetID(), "1234", "%s", k)
	assert.NotEmpty(t, c.GetHashedSecret(), "%s", k)
	assert.Equal(t, c.GetRedirectURIs(), []string{"http://redirect"}, "%s", k)
}
