package client_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	. "github.com/ory/hydra/client"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/integration"
	"github.com/ory/ladon"
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
		t.Run(fmt.Sprintf("case=%s", k), TestHelperClientAutoGenerateKey(k, m))
	}
}

func TestAuthenticateClient(t *testing.T) {
	var mem = &MemoryManager{
		Clients: map[string]Client{},
		Hasher:  &fosite.BCrypt{},
	}

	TestHelperClientAuthenticate("", mem)(t)
}

func TestCreateGetDeleteClient(t *testing.T) {
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteClient(k, m))
	}
}
