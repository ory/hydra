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
package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
)

type testCase struct {
	name string
	b    BackendConnector
	u    string
}

var (
	tests []testCase = []testCase{
		{
			"memory",
			&MemoryBackend{},
			"memory",
		},
	}
	l                = logrus.New()
	hasher           = &fosite.BCrypt{WorkFactor: 8}
	encryptionKey, _ = jwk.RandomBytes(32)
	cipher           = &jwk.AEAD{Key: encryptionKey}
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		if uri := os.Getenv("TEST_DATABASE_POSTGRESQL"); uri != "" {
			tests = append(tests, testCase{"postgresql", &SQLBackend{}, uri})
		} else {
			log.Println("Did not find postgresql test database config, skipping backend connector test")
		}

		if uri := os.Getenv("TEST_DATABASE_MYSQL"); uri != "" {
			if !strings.HasPrefix(uri, "mysql") {
				uri = fmt.Sprintf("mysql://%s", uri)
			}
			tests = append(tests, testCase{"mysql", &SQLBackend{}, uri})
		} else {
			log.Println("Did not find mysql test database config, skipping backend connector test")
		}
	}

	os.Exit(m.Run())
}

func TestBackendConnectors(t *testing.T) {
	for _, tc := range tests {
		var cm client.Manager
		var fs pkg.FositeStorer

		t.Run(fmt.Sprintf("%s/Init", tc.name), func(t *testing.T) {
			if err := tc.b.Init(tc.u, l); err != nil {
				t.Fatalf("could not initialize backend due to error: %v", err)
			}
		})

		t.Run(fmt.Sprintf("%s/Ping", tc.name), func(t *testing.T) {
			if err := tc.b.Ping(); err != nil {
				t.Errorf("could not ping backend due to error: %v", err)
			}
		})

		t.Run(fmt.Sprintf("%s/NewClientManager", tc.name), func(t *testing.T) {
			if cm = tc.b.NewClientManager(hasher); cm == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/NewOAuth2Manager", tc.name), func(t *testing.T) {
			if fs = tc.b.NewOAuth2Manager(cm, time.Hour, "opaque"); fs == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/NewConsentManager", tc.name), func(t *testing.T) {
			if want := tc.b.NewConsentManager(cm, fs); want == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/NewJWKManager", tc.name), func(t *testing.T) {
			if want := tc.b.NewJWKManager(cipher); want == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/Prefixes", tc.name), func(t *testing.T) {
			prefixes := tc.b.Prefixes()
			for _, prefix := range prefixes {
				if strings.HasPrefix(tc.u, prefix) {
					return
				}
			}
			t.Errorf("did not find matching prefix for given backend uri")
		})
	}
}
