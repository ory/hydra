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

var testGenerators = (&Handler{}).GetGenerators()

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
	for algo, testGenerator := range testGenerators {
		if algo == "HS256" {
			// this is a symmetrical algorithm
			continue
		}

		ks, err := testGenerator.Generate("")
		if err != nil {
			t.Fatal(err)
		}

		for name, m := range managers {
			t.Run(fmt.Sprintf("case=%s/%s", algo, name), func(t *testing.T) {
				TestHelperManagerKey(m, algo, ks)(t)
			})
		}
	}
}

func TestManagerKeySet(t *testing.T) {
	for algo, testGenerator := range testGenerators {
		ks, err := testGenerator.Generate("")
		if err != nil {
			t.Fatal(err)
		}
		ks.Key("private")

		for name, m := range managers {
			t.Run(fmt.Sprintf("case=%s/%s", algo, name), func(t *testing.T) {
				TestHelperManagerKeySet(m, algo, ks)(t)
			})
		}
	}
}
