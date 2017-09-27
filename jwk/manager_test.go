package jwk_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/jwk"
)

var managers = map[string]Manager{
	"memory": new(MemoryManager),
}

var testGenerator = &RS256Generator{}

var encryptionKey, _ = RandomBytes(32)

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

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	managers["postgres"] = s
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	managers["mysql"] = s
}

func TestManagerKey(t *testing.T) {
	ks, _ := testGenerator.Generate("")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), func(t *testing.T) {
			TestHelperManagerKey(m, ks)(t)
		})
	}
}

func TestManagerKeySet(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	ks.Key("private")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), func(t *testing.T) {
			TestHelperManagerKeySet(m, ks)(t)
		})
	}
}
