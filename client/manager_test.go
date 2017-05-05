package client_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	. "github.com/ory/hydra/client"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/integration"
	"github.com/ory/hydra/pkg"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ory-am/dockertest.v3"
)

var clientManagers = map[string]Storage{}

var ts *httptest.Server

func init() {
	clientManagers["memory"] = &MemoryManager{
		Clients: map[string]Client{},
		Hasher:  &fosite.BCrypt{},
	}

	localWarden, httpClient := compose.NewMockFirewall("foo", "alice", fosite.Arguments{Scope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:clients<.*>"},
		Actions:   []string{"create", "get", "delete", "update"},
		Effect:    ladon.AllowAccess,
	})

	s := &Handler{
		Manager: &MemoryManager{
			Clients: map[string]Client{},
			Hasher:  &fosite.BCrypt{},
		},
		H: herodot.NewJSONWriter(nil),
		W: localWarden,
	}

	routing := httprouter.New()
	s.SetRoutes(routing)
	ts = httptest.NewServer(routing)

	u, _ := url.Parse(ts.URL + ClientsHandlerPath)
	clientManagers["http"] = &HTTPManager{
		Client:   httpClient,
		Endpoint: u,
	}
}

var resources []*dockertest.Resource
var pool *dockertest.Pool

func TestMain(m *testing.M) {
	connectToPG()
	connectToMySQL()

	s := m.Run()
	integration.KillAll()
	os.Exit(s)
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func TestClientAutoGenerateKey(t *testing.T) {
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			c := &Client{
				Secret:            "secret",
				RedirectURIs:      []string{"http://redirect"},
				TermsOfServiceURI: "foo",
			}
			assert.Nil(t, m.CreateClient(c))
			assert.NotEmpty(t, c.ID)
			assert.Nil(t, m.DeleteClient(c.ID))
		})
	}
}

func TestAuthenticateClient(t *testing.T) {
	var mem = &MemoryManager{
		Clients: map[string]Client{},
		Hasher:  &fosite.BCrypt{},
	}
	mem.CreateClient(&Client{
		ID:           "1234",
		Secret:       "secret",
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
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetClient("4321")
			assert.NotNil(t, err)

			c := &Client{
				ID:                "1234",
				Name:              "name",
				Secret:            "secret",
				RedirectURIs:      []string{"http://redirect"},
				TermsOfServiceURI: "foo",
			}
			err = m.CreateClient(c)
			assert.Nil(t, err)
			if err == nil {
				compare(t, c, k)
			}

			err = m.CreateClient(&Client{
				ID:                "2-1234",
				Name:              "name",
				Secret:            "secret",
				RedirectURIs:      []string{"http://redirect"},
				TermsOfServiceURI: "foo",
			})
			assert.Nil(t, err)

			// RethinkDB delay
			time.Sleep(100 * time.Millisecond)

			d, err := m.GetClient("1234")
			assert.Nil(t, err)
			if err == nil {
				compare(t, d, k)
			}

			ds, err := m.GetClients()
			assert.Nil(t, err)
			assert.Len(t, ds, 2)
			assert.NotEqual(t, ds["1234"].ID, ds["2-1234"].ID)

			err = m.UpdateClient(&Client{
				ID:                "2-1234",
				Name:              "name-new",
				Secret:            "secret-new",
				TermsOfServiceURI: "bar",
			})
			assert.Nil(t, err)
			time.Sleep(100 * time.Millisecond)

			nc, err := m.GetConcreteClient("2-1234")
			assert.Nil(t, err)

			if k != "http" {
				// http always returns an empty secret
				assert.NotEqual(t, d.GetHashedSecret(), nc.GetHashedSecret(), "%s", k)
			}
			assert.Equal(t, "bar", nc.TermsOfServiceURI, "%s", k)
			assert.Equal(t, "name-new", nc.Name, "%s", k)
			assert.EqualValues(t, []string{"http://redirect"}, nc.GetRedirectURIs(), "%s", k)

			err = m.DeleteClient("1234")
			assert.Nil(t, err)

			// RethinkDB delay
			time.Sleep(100 * time.Millisecond)

			_, err = m.GetClient("1234")
			assert.NotNil(t, err)
		})
	}
}

func compare(t *testing.T, c fosite.Client, k string) {
	assert.Equal(t, c.GetID(), "1234", "%s", k)
	if k != "http" {
		assert.NotEmpty(t, c.GetHashedSecret(), "%s", k)
	}
	assert.Equal(t, c.GetRedirectURIs(), []string{"http://redirect"}, "%s", k)
}
