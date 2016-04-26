package oauth2_test

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite/client"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/server"
	wcl "github.com/ory-am/hydra/warden/client"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"github.com/ory-am/fosite"
	. "github.com/ory-am/hydra/oauth2"
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

var warden = &wcl.LocalWarden{
	Warden: ladonWarden,
	TokenValidator: &core.CoreValidator{
		AccessTokenStrategy: pkg.HMACStrategy,
		AccessTokenStorage:  fositeStore,
	},
	Issuer: "tests",
}

var ts *httptest.Server

var tokens = pkg.Tokens(1)

func init() {
	ar := fosite.NewAccessRequest(&Session{Subject: "alice"})
	ar.GrantedScopes = fosite.Arguments{"hydra.oauth2.clients"}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	clientManagers["memory"] = &MemoryClientManager{Clients: map[string]*OAuth2Client{}}

	s := server.OAuth2Client{
		Manager: &MemoryClientManager{Clients: map[string]*OAuth2Client{}},
		H:       &herodot.JSON{},
		W:       warden,
	}
	r := httprouter.New()
	s.SetRoutes(r)
	ts = httptest.NewServer(r)
	conf := &oauth2.Config{Scopes: []string{}, Endpoint: oauth2.Endpoint{}}

	u, _ := url.Parse(ts.URL + "/clients")
	clientManagers["http"] = &HTTPClientManager{
		Client: conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		}),
		Endpoint: u,
	}
}

func TestCreateGetDeleteClient(t *testing.T) {
	for k, m := range clientManagers {
		_, err := m.GetClient("4321")
		pkg.AssertError(t, true, err, k)

		c := &OAuth2Client{
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

func compare(t *testing.T, c client.Client, k string) {
	assert.Equal(t, c.GetID(), "1234", "%s", k)
	assert.NotEmpty(t, c.GetHashedSecret(), "%s", k)
	assert.Equal(t, c.GetRedirectURIs(), []string{"http://redirect"}, "%s", k)
}
