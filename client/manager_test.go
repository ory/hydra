package client_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/hash"
	. "github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
)

var clientManagers = map[string]Storage{}

var ts *httptest.Server

func init() {
	clientManagers["memory"] = &MemoryManager{
		Clients: map[string]*fosite.DefaultClient{},
		Hasher:  &hash.BCrypt{},
	}

	localWarden, httpClient := internal.NewFirewall("foo", "alice", fosite.Arguments{Scope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:clients<.*>"},
		Actions:   []string{"create", "get", "delete"},
		Effect:    ladon.AllowAccess,
	})

	s := &Handler{
		Manager: &MemoryManager{
			Clients: map[string]*fosite.DefaultClient{},
			Hasher:  &hash.BCrypt{},
		},
		H: &herodot.JSON{},
		W: localWarden,
	}

	r := httprouter.New()
	s.SetRoutes(r)
	ts = httptest.NewServer(r)

	u, _ := url.Parse(ts.URL + ClientsHandlerPath)
	clientManagers["http"] = &HTTPManager{
		Client:   httpClient,
		Endpoint: u,
	}
}

func TestAuthenticateClient(t *testing.T) {
	var mem = &MemoryManager{
		Clients: map[string]*fosite.DefaultClient{},
		Hasher:  &hash.BCrypt{},
	}
	mem.CreateClient(&fosite.DefaultClient{
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

		c := &fosite.DefaultClient{
			ID:                "1234",
			Secret:            []byte("secret"),
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
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

		ds, err := m.GetClients()
		pkg.AssertError(t, false, err, k)
		assert.Len(t, ds, 1)

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
