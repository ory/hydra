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
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	. "github.com/ory/hydra/jwk"
	"github.com/ory/sqlcon/dockertest"
	"github.com/stretchr/testify/require"
)

var managers = map[string]Manager{
	"memory": new(MemoryManager),
}

var testGenerator = &RS256Generator{}

var encryptionKey, _ = RandomBytes(32)

func TestMain(m *testing.M) {
	runner := dockertest.Register()

	flag.Parse()
	if !testing.Short() {
		dockertest.Parallel([]func(){
			connectToPG,
			connectToMySQL,
		})
	}

	runner.Exit(m.Run())
}

func connectToPG() {
	db, err := dockertest.ConnectToTestPostgreSQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create schema: %v", err)
	}

	managers["postgres"] = s
}

func connectToMySQL() {
	db, err := dockertest.ConnectToTestMySQL()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create schema: %v", err)
	}

	managers["mysql"] = s
}

func TestManagerKey(t *testing.T) {
	ks, err := testGenerator.Generate("TestManagerKey", "sig")
	require.NoError(t, err)

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKey(m, ks, "TestManagerKey"))
	}
}

func TestManagerKeySet(t *testing.T) {
	ks, err := testGenerator.Generate("TestManagerKeySet", "sig")
	require.NoError(t, err)
	ks.Key("private")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), TestHelperManagerKeySet(m, ks, "TestManagerKeySet"))
	}
}
