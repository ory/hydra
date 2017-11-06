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

package integration

import (
	"fmt"
	"log"
	"time"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

var resources []*dockertest.Resource
var pool *dockertest.Pool

func KillAll() {
	for _, resource := range resources {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Got an error while trying to purge resource: %s", err)
		}
	}
	resources = []*dockertest.Resource{}
}

func ConnectToMySQL() *sqlx.DB {
	if url := os.Getenv("TEST_DATABASE_MYSQL"); url != "" {
		log.Println("Found mysql test database config, skipping dockertest...")
		db, err := sqlx.Open("mysql", url)
		if err != nil {
			log.Fatalf("Could not connect to bootstrapped database: %s", err)
		}
		return db
	}

	var db *sqlx.DB
	var err error
	pool, err = dockertest.NewPool("")
	pool.MaxWait = time.Minute * 5
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resources = append(resources, resource)
	return db
}

func ConnectToPostgres() *sqlx.DB {
	if url := os.Getenv("TEST_DATABASE_POSTGRESQL"); url != "" {
		log.Println("Found postgresql test database config, skipping dockertest...")
		db, err := sqlx.Open("postgres", url)
		if err != nil {
			log.Fatalf("Could not connect to bootstrapped database: %s", err)
		}
		return db
	}

	var db *sqlx.DB
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=hydra"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/hydra?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resources = append(resources, resource)
	return db
}
