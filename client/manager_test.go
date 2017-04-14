package client_test

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ory-am/fosite"
	. "github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/compose"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/integration"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gopkg.in/ory-am/dockertest.v3"
)

var clientManagers = map[string]Storage{}

var ts *httptest.Server

func init() {
	clientManagers["memory"] = &MemoryManager{
		Clients: map[string]Client{},
		Hasher:  &fosite.BCrypt{},
	}

	localWarden, httpClient := compose.NewFirewall("foo", "alice", fosite.Arguments{Scope}, &ladon.DefaultPolicy{
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
		H: &herodot.JSON{},
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
var rethinkManager *RethinkManager

func TestMain(m *testing.M) {
	connectToPG()
	connectToRethinkDB()
	connectToMySQL()
	connectToRedis()

	s := m.Run()
	integration.KillAll()
	os.Exit(s)
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}
	if err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &SQLManager{DB: db, Hasher: &fosite.BCrypt{WorkFactor: 4}}

	if err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func connectToRethinkDB() {
	var session = integration.ConnectToRethinkDB("hydra", "hydra_clients")
	rethinkManager = &RethinkManager{
		Session: session,
		Table:   r.Table("hydra_clients"),
		Clients: make(map[string]Client),
		Hasher: &fosite.BCrypt{
			// Low workfactor reduces test time
			WorkFactor: 4,
		},
	}

	rethinkManager.Watch(context.Background())
	clientManagers["rethink"] = rethinkManager
}

func connectToRedis() {
	var db = integration.ConnectToRedis()
	clientManagers["redis"] = &RedisManager{
		DB:     db,
		Hasher: &fosite.BCrypt{WorkFactor: 4},
	}
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

func BenchmarkRethinkGet(b *testing.B) {
	b.StopTimer()

	m := rethinkManager
	id := uuid.New()
	c := &Client{
		ID:                id,
		Secret:            "secret",
		RedirectURIs:      []string{"http://redirect"},
		TermsOfServiceURI: "foo",
	}

	var err error
	err = m.CreateClient(c)
	if err != nil {
		b.Fatalf("%s", err)
	}
	time.Sleep(100 * time.Millisecond)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.GetClient(id)
	}
}

func BenchmarkRethinkAuthenticate(b *testing.B) {
	b.StopTimer()

	m := rethinkManager
	id := uuid.New()
	c := &Client{
		ID:                id,
		Secret:            "secret",
		RedirectURIs:      []string{"http://redirect"},
		TermsOfServiceURI: "foo",
	}

	var err error
	err = m.CreateClient(c)
	if err != nil {
		b.Fatalf("%s", err)
	}
	time.Sleep(100 * time.Millisecond)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Authenticate(id, []byte("secret"))
	}
}

func TestColdStartRethinkManager(t *testing.T) {
	assert.Nil(t, rethinkManager.CreateClient(&Client{ID: "foo"}))
	assert.Nil(t, rethinkManager.CreateClient(&Client{ID: "bar"}))

	time.Sleep(time.Second / 2)
	rethinkManager.Clients = make(map[string]Client)
	require.Nil(t, rethinkManager.ColdStart())

	c1, err := rethinkManager.GetClient("foo")
	require.Nil(t, err)
	c2, err := rethinkManager.GetClient("bar")
	require.Nil(t, err)

	assert.NotEqual(t, c1, c2)
	assert.Equal(t, "foo", c1.GetID())

	rethinkManager.Clients = make(map[string]Client)
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
