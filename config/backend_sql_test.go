// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mysql    *url.URL
	postgres *url.URL
)
var resources []*dockertest.Resource

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		mysql = bootstrapMySQL()
		postgres = bootstrapPostgres()
	}

	s := m.Run()
	killAll()
	os.Exit(s)
}

func merge(u *url.URL, params map[string]string) *url.URL {
	b := new(url.URL)
	*b = *u
	for k, v := range params {
		b.Query().Add(k, v)
	}
	return b
}

func TestCleanQueryURL(t *testing.T) {
	a, _ := url.Parse("mysql://foo:bar@baz/db?max_conn_lifetime=1h&max_idle_conns=10&max_conns=10")
	b := cleanURLQuery(a)
	assert.NotEqual(t, a, b)
	assert.NotEqual(t, a.String(), b.String())
	assert.Equal(t, true, strings.Contains(a.String(), "max_conn_lifetime"))
	assert.Equal(t, false, strings.Contains(b.String(), "max_conn_lifetime"))
}

func TestSQLConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
		return
	}

	for _, tc := range []struct {
		s *SQLConnection
		d string
	}{
		{
			d: "mysql raw",
			s: &SQLConnection{
				URL: mysql,
			},
		},
		{
			d: "mysql max_conn_lifetime",
			s: &SQLConnection{
				URL: merge(mysql, map[string]string{"max_conn_lifetime": "1h"}),
			},
		},
		{
			d: "mysql max_conn_lifetime",
			s: &SQLConnection{
				URL: merge(mysql, map[string]string{"max_conn_lifetime": "1h", "max_idle_conns": "10", "max_conns": "10"}),
			},
		},
		{
			d: "pg raw",
			s: &SQLConnection{
				URL: postgres,
			},
		},
		{
			d: "pg max_conn_lifetime",
			s: &SQLConnection{
				URL: merge(postgres, map[string]string{"max_conn_lifetime": "1h"}),
			},
		},
		{
			d: "pg max_conn_lifetime",
			s: &SQLConnection{
				URL: merge(postgres, map[string]string{"max_conn_lifetime": "1h", "max_idle_conns": "10", "max_conns": "10"}),
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%s", tc.d), func(t *testing.T) {
			tc.s.L = logrus.New()
			db := tc.s.GetDatabase()
			require.Nil(t, db.Ping())
		})
	}
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

func bootstrapMySQL() *url.URL {
	if uu := os.Getenv("TEST_DATABASE_MYSQL"); uu != "" {
		log.Println("Found mysql test database config, skipping dockertest...")
		_, err := sqlx.Open("postgres", uu)
		if err != nil {
			log.Fatalf("Could not connect to bootstrapped database: %s", err)
		}
		u, _ := url.Parse("mysql://" + uu)
		return u
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

	resources = append(resources, resource)
	u, _ := url.Parse("mysql://" + urls)
	return u
}

func bootstrapPostgres() *url.URL {
	if uu := os.Getenv("TEST_DATABASE_POSTGRESQL"); uu != "" {
		log.Println("Found postgresql test database config, skipping dockertest...")
		_, err := sqlx.Open("postgres", uu)
		if err != nil {
			log.Fatalf("Could not connect to bootstrapped database: %s", err)
		}
		u, _ := url.Parse(uu)
		return u
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

	resources = append(resources, resource)
	u, _ := url.Parse(urls)
	return u
}
