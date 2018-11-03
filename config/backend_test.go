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
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ory/dockertest"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	dockertestd "github.com/ory/sqlcon/dockertest"
)

var resources []*dockertest.Resource

type connectorFixture struct {
	name      string
	connector BackendConnector
	dsn       string
}

var (
	testConnectors = []connectorFixture{
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
		dockertestd.Parallel([]func(){
			bootstrapMySQL,
			bootstrapPostgres,
		})
	}

	s := m.Run()
	killAll()
	os.Exit(s)
}

func TestBackendConnectors(t *testing.T) {
	for _, tc := range testConnectors {
		var cm client.Manager
		var fs pkg.FositeStorer

		t.Run(fmt.Sprintf("%s/Init", tc.name), func(t *testing.T) {
			if err := tc.connector.Init(tc.dsn, l); err != nil {
				t.Fatalf("could not initialize backend due to error: %v", err)
			}
		})

		t.Run(fmt.Sprintf("%s/Ping", tc.name), func(t *testing.T) {
			if err := tc.connector.Ping(); err != nil {
				t.Errorf("could not ping backend due to error: %v", err)
			}
		})

		t.Run(fmt.Sprintf("%s/NewClientManager", tc.name), func(t *testing.T) {
			if cm = tc.connector.NewClientManager(hasher); cm == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/NewOAuth2Manager", tc.name), func(t *testing.T) {
			if fs = tc.connector.NewOAuth2Manager(cm, time.Hour, "opaque"); fs == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/NewConsentManager", tc.name), func(t *testing.T) {
			if want := tc.connector.NewConsentManager(cm, fs); want == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/NewJWKManager", tc.name), func(t *testing.T) {
			if want := tc.connector.NewJWKManager(cipher); want == nil {
				t.Errorf("expected non-nil result")
			}
		})

		t.Run(fmt.Sprintf("%s/Prefixes", tc.name), func(t *testing.T) {
			prefixes := tc.connector.Prefixes()
			for _, prefix := range prefixes {
				if strings.HasPrefix(tc.dsn, prefix) {
					return
				}
			}
			t.Errorf("did not find matching prefix for given backend uri")
		})
	}
}

func TestConnectorTracingOptions(t *testing.T) {
	for _, fixture := range testConnectors {
		if _, ok := fixture.connector.(*SQLBackend); !ok {
			// memory connector does not support tracing - skip
			continue
		}

		for _, testCase := range []struct {
			description string
			options     []ConnectorOptions
		}{
			{
				description: "WithTracing() option should result in spans being created on database interactions",
				options: []ConnectorOptions{
					WithTracing(),
					withAllowRootTraceSpans(),
					withRandomDriverName(), // Note: this option is being used because only one driver can be registered with the same name
				},
			},
			{
				description: "No spans should be created if tracing options have not been set",
				options:     []ConnectorOptions{},
			},
			{
				description: "No spans should be created if no trace exists in the supplied context when" +
					" withAllowRootTraceSpans() option has NOT been set",
				options: []ConnectorOptions{
					WithTracing(),
					withRandomDriverName(), // Note: this option is being used because only one driver can be registered with the same name
				},
			},
		} {
			t.Run(fmt.Sprintf("connector=%s - test scenario=%s", fixture.name, testCase.description), func(t *testing.T) {
				mockedTracer := mocktracer.New()
				defer mockedTracer.Reset()
				opentracing.SetGlobalTracer(mockedTracer)

				// using a fresh SQLBackend here so that the options are reset between tests
				// note: use of var is intentional here to highlight the fact that SQLBackend satisfies the BackendConnector interface
				var connector BackendConnector = &SQLBackend{}
				assert.NoError(t, connector.Init(fixture.dsn, l, testCase.options...))
				sqlBackend, ok := connector.(*SQLBackend)
				assert.True(t, ok)
				db := sqlBackend.db
				assert.NotNil(t, db)

				// notice how no parent span exists in the provided context to this query. This is useful for testing the
				// behaviour when WithAllowRootTraceSpans() is (un)set.
				db.QueryRowContext(context.TODO(), "SELECT NOW()")
				spans := mockedTracer.FinishedSpans()

				if sqlBackend.UseTracing && sqlBackend.allowRootTracingSpans {
					assert.NotEmpty(t, spans)
				} else {
					assert.Empty(t, spans)
				}
			})
		}
	}
}

func bootstrapMySQL() {
	if uu := os.Getenv("TEST_DATABASE_MYSQL"); uu != "" {
		log.Println("Found mysql test database config, skipping dockertest...")
		_, err := sqlx.Open("postgres", uu)
		if err != nil {
			log.Fatalf("Could not connect to bootstrapped database: %s", err)
		}
		u, _ := url.Parse("mysql://" + uu)
		testConnectors = append(testConnectors, connectorFixture{"mysql", &SQLBackend{}, u.String()})
		return
	}

	var db *sqlx.DB
	var err error
	var urls string

	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 5
	if err != nil {
		log.Fatalf("Could not Connect to docker: %s", err)
	}

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		urls = fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp"))
		db, err = sqlx.Open("mysql", urls)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not Connect to docker: %s", err)
	}

	u, _ := url.Parse("mysql://" + urls)
	resources = append(resources, resource)
	testConnectors = append(testConnectors, connectorFixture{"mysql", &SQLBackend{}, u.String()})
}

func bootstrapPostgres() {
	if uu := os.Getenv("TEST_DATABASE_POSTGRESQL"); uu != "" {
		log.Println("Found postgresql test database config, skipping dockertest...")
		_, err := sqlx.Open("postgres", uu)
		if err != nil {
			log.Fatalf("Could not connect to bootstrapped database: %s", err)
		}
		u, _ := url.Parse(uu)
		testConnectors = append(testConnectors, connectorFixture{"postgresql", &SQLBackend{}, u.String()})
		return
	}

	var db *sqlx.DB
	var err error
	var urls string

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not Connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=hydra"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		urls = fmt.Sprintf("postgres://postgres:secret@localhost:%s/hydra?sslmode=disable", resource.GetPort("5432/tcp"))
		db, err = sqlx.Open("postgres", urls)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not Connect to docker: %s", err)
	}

	u, _ := url.Parse(urls)
	resources = append(resources, resource)
	testConnectors = append(testConnectors, connectorFixture{"postgresql", &SQLBackend{}, u.String()})
}

func killAll() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not Connect to pool because %s", err)
	}

	for _, resource := range resources {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Got an error while trying to purge resource: %s", err)
		}
	}

	resources = []*dockertest.Resource{}
}
