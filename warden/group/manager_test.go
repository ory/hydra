package group_test

import (
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"fmt"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
)

var clientManagers = map[string]Manager{}
var ts *httptest.Server

func init() {
	clientManagers["memory"] = &MemoryManager{
		Groups: map[string]Group{},
	}

	localWarden, httpClient := compose.NewMockFirewall("foo", "alice", fosite.Arguments{Scope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:warden<.*>"},
		Actions:   []string{"create", "get", "delete", "update", "members.add", "members.remove"},
		Effect:    ladon.AllowAccess,
	})

	s := &Handler{
		Manager: &MemoryManager{
			Groups: map[string]Group{},
		},
		H: herodot.NewJSONWriter(nil),
		W: localWarden,
	}

	routing := httprouter.New()
	s.SetRoutes(routing)
	ts = httptest.NewServer(routing)

	u, _ := url.Parse(ts.URL + GroupsHandlerPath)
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
	s := &SQLManager{DB: db}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &SQLManager{DB: db}

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func TestManagers(t *testing.T) {
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperManagers(m))
	}
}
