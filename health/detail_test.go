package health

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/hydra/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	mysql    *url.URL
	postgres *url.URL
)

var resources []*dockertest.Resource

func TestMain(m *testing.M) {
	mysql = bootstrapMySQL()
	postgres = bootstrapPostgres()

	s := m.Run()
	killAll()
	os.Exit(s)
}

func TestDbCheck(t *testing.T) {
	for _, tc := range []struct {
		s *config.SQLConnection
		d string
		n string
	}{
		{
			d: "mysql",
			s: &config.SQLConnection{
				URL: mysql,
			},
		},
		{
			d: "postgres",
			s: &config.SQLConnection{
				URL: postgres,
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%s", tc.d), func(t *testing.T) {
			tc.s.L = logrus.New()
			tc.s.GetDatabase()
			info := dbCheck(tc.s)
			assert.Equal(t, info.Name, "Database ("+tc.d+")")
			assert.Equal(t, info.Type, "internal")
			assert.Equal(t, info.State.Status, OK)
		})
	}

	s := &config.MemoryConnection{}
	info := dbCheck(s)
	assert.Equal(t, info.Name, "Database (memory)")
	assert.Equal(t, info.Type, "internal")
	assert.Equal(t, info.State.Status, OK)
	assert.Equal(t, int(info.ResponseTime), 0)
}

func TestDbCheckWithError(t *testing.T) {
	info := dbCheck(nil)
	assert.Equal(t, info.Type, "internal")
	assert.Equal(t, info.State.Status, CRIT)
	assert.Equal(t, info.State.Details, "No DB connection")
}

func TestDbCheckWithDBError(t *testing.T) {
	s := &config.SQLConnection{
		URL: &url.URL{},
		L:   logrus.New(),
	}
	info := dbCheck(s)
	assert.Equal(t, info.Name, "Database ()")
	assert.Equal(t, info.Type, "internal")
	assert.Equal(t, info.State.Status, "CRIT")
	assert.Equal(t, info.State.Details, "Database uninitialized")
}

func TestSimpleStatus(t *testing.T) {
	data := simpleStatus(new(config.Config))
	var d map[string]interface{}
	json.Unmarshal(data, &d)
	assert.Equal(t, d["status"], OK)
}

func TestDetailedStatusWithoutConnection(t *testing.T) {
	var c = new(config.Config)
	//Set to memory so we don't get an error when retrieving the context.
	c.DatabaseURL = "memory"
	c.Context().Connection = nil

	data := detailedStatus(c)
	var d map[string]interface{}
	json.Unmarshal(data, &d)
	assert.Equal(t, d["status"], CRIT)
	assert.Equal(t, d["name"], "Sand")

	dependent := d["dependencies"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, dependent["name"], "Database ()")
	assert.Equal(t, dependent["type"], "internal")

	state := dependent["state"].(map[string]interface{})
	assert.Equal(t, state["status"], CRIT)
	assert.Equal(t, state["details"], "No DB connection")
}

func TestDetailedStatusWithSQLConnection(t *testing.T) {
	var c = new(config.Config)
	//Set to memory so we don't get an error when retrieving the context.
	c.DatabaseURL = "memory"
	s := &config.SQLConnection{
		URL: mysql,
		L:   logrus.New(),
	}
	s.GetDatabase()
	c.Context().Connection = s
	data := detailedStatus(c)
	var d map[string]interface{}
	json.Unmarshal(data, &d)
	assert.Equal(t, d["status"], OK)
	assert.Equal(t, d["name"], "Sand")

	dependent := d["dependencies"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, dependent["name"], "Database (mysql)")
	assert.Equal(t, dependent["type"], "internal")

	state := dependent["state"].(map[string]interface{})
	assert.Equal(t, state["status"], OK)
	assert.Nil(t, state["details"])
}

func TestGetProjectWithData(t *testing.T) {
	os.Setenv("APPLICATION_LOG_LINKS", "http://log1.com https://log2.com")
	os.Setenv("APPLICATION_STATS_LINKS", "http://stats1.com https://stats2.com")

	info := getProject()
	assert.Equal(t, info.Logs, []string{"http://log1.com", "https://log2.com"})
	assert.Equal(t, info.Stats, []string{"http://stats1.com", "https://stats2.com"})
}

func TestGetProjectWithoutData(t *testing.T) {
	os.Setenv("APPLICATION_LOG_LINKS", "")
	os.Setenv("APPLICATION_STATS_LINKS", "")

	info := getProject()
	assert.Equal(t, info.Logs, []string{""})
	assert.Equal(t, info.Stats, []string{""})
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
