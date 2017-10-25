package group_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/warden/group"
)

var clientManagers = map[string]Manager{
	"memory": &MemoryManager{
		Groups: map[string]Group{},
	},
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		connectToPG()
		connectToMySQL()
	}

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
	t.Parallel()

	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperManagers(m))
	}
}
