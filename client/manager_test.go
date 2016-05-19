package client_test

import (
	r "github.com/dancannon/gorethink"
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
	"time"
	"os"
	"gopkg.in/ory-am/dockertest.v2"
	"log"
	"golang.org/x/net/context"
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

func TestMain(m *testing.M) {
	var session *r.Session
	var rethinkManager *RethinkManager
	var err error

	c, err := dockertest.ConnectToRethinkDB(20, time.Second, func(url string) bool {
		if session, err = r.Connect(r.ConnectOpts{Address:  url, Database: "hydra"}); err != nil {
			return false
		} else if _, err = r.DBCreate("hydra").RunWrite(session); err != nil {
			log.Printf("Database exists: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_clients").RunWrite(session); err != nil {
			log.Printf("Could not create table: %s", err)
			return false
		}

		rethinkManager = &RethinkManager{
			Session: session,
			Table: r.Table("hydra_clients"),
			Clients:make(map[string]*fosite.DefaultClient),
			Hasher: &hash.BCrypt{},
		}
		err := rethinkManager.Watch(context.Background())
		if err != nil {
			log.Printf("Could not watch: %s", err)
			return false
		}
		return true
	})
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	clientManagers["rethink"] = rethinkManager

	retCode := m.Run()
	c.KillRemove()
	os.Exit(retCode)
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
		pkg.AssertError(t, true, err, "%s", k)

		c := &fosite.DefaultClient{
			ID:                "1234",
			Secret:            []byte("secret"),
			RedirectURIs:      []string{"http://redirect"},
			TermsOfServiceURI: "foo",
		}
		err = m.CreateClient(c)
		pkg.AssertError(t, false, err, "%s", k)
		if err == nil {
			compare(t, c, k)
		}

		// RethinkDB delay
		time.Sleep(time.Microsecond * 500)

		d, err := m.GetClient("1234")
		pkg.AssertError(t, false, err, "%s", k)
		if err == nil {
			compare(t, d, k)
		}

		ds, err := m.GetClients()
		pkg.AssertError(t, false, err, "%s", k)
		assert.Len(t, ds, 1)

		err = m.DeleteClient("1234")
		pkg.AssertError(t, false, err, "%s", k)

		// RethinkDB delay
		time.Sleep(time.Microsecond * 500)

		_, err = m.GetClient("1234")
		pkg.AssertError(t, true, err, "%s", k)
	}
}

func compare(t *testing.T, c fosite.Client, k string) {
	assert.Equal(t, c.GetID(), "1234", "%s", k)
	assert.NotEmpty(t, c.GetHashedSecret(), "%s", k)
	assert.Equal(t, c.GetRedirectURIs(), []string{"http://redirect"}, "%s", k)
}
