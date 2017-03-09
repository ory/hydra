package integration

import (
	r "gopkg.in/dancannon/gorethink.v2"
	"log"
	"time"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/ory-am/dockertest.v3"
	"gopkg.in/redis.v5"
)

var resources []*dockertest.Resource
var pool *dockertest.Pool

func KillAll() {
	for _, resource := range resources {
		pool.Purge(resource)
	}
	resources = []*dockertest.Resource{}
}

func ConnectToMySQL() *sqlx.DB {
	var db *sqlx.DB
	var err error
	pool, err = dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
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
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resources = append(resources, resource)
	return db
}

func ConnectToPostgres(database string) *sqlx.DB {
	var db *sqlx.DB
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + database})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), database))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resources = append(resources, resource)
	return db
}

func ConnectToRethinkDB(database string, tables ...string) *r.Session {
	var session *r.Session
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("rethinkdb", "2.3", []string{""})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		if session, err = r.Connect(r.ConnectOpts{Address: fmt.Sprintf("localhost:%s", resource.GetPort("28015/tcp")), Database: database}); err != nil {
			return err
		} else if _, err = r.DBCreate(database).RunWrite(session); err != nil {
			log.Printf("Database exists: %s", err)
			return err
		}

		for _, table := range tables {
			if _, err = r.TableCreate(table).RunWrite(session); err != nil {
				log.Printf("Could not create table: %s", err)
				return err
			}
		}

		time.Sleep(100 * time.Millisecond)
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resources = append(resources, resource)
	return session
}

func ConnectToRedis() *redis.Client {
	var db *redis.Client
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "3.2", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})

		return db.Ping().Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resources = append(resources, resource)
	return db
}
