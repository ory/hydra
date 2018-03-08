/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

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
		if !testing.Short() {
			integration.BootParallel([]func(){
				connectToPG,
				connectToMySQL,
			})
		}
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
	ks, _ := testGenerator.Generate("TestManagerKey")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKey(m, ks, "TestManagerKey"))
	}
}

func TestManagerKeySet(t *testing.T) {
	ks, _ := testGenerator.Generate("TestManagerKeySet")
	ks.Key("private")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKeySet(m, ks, "TestManagerKeySet"))
	}
}
