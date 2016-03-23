package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory-am/common/pkg"
	. "github.com/ory-am/hydra/endpoint/connection"
	. "github.com/ory-am/hydra/endpoint/connection/postgres"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ory-am/dockertest.v2"
)

var store *Store
var db *sql.DB

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	store = New(db)
	if err = store.CreateSchemas(); err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	db.Close()
	c.KillRemove()
}

var connection = &DefaultConnection{ID: uuid.New(), LocalSubject: "peter", RemoteSubject: "peterson", Provider: "google"}

func TestNotFound(t *testing.T) {
	_, err := store.Get("asdf")
	assert.Equal(t, pkg.ErrNotFound, err)
}

func TestCreateGetFindDelete(t *testing.T) {
	require.Nil(t, store.Create(connection))

	c, err := store.Get(connection.ID)
	require.Nil(t, err)
	require.Equal(t, connection, c)

	c, err = store.FindByRemoteSubject("google", "peterson")
	require.Nil(t, err)
	require.Equal(t, connection, c)

	cs, err := store.FindAllByLocalSubject("peter")
	require.Nil(t, err)
	require.Equal(t, connection, cs[0])

	require.Nil(t, store.Delete(connection.ID))
	_, err = store.Get(connection.ID)
	require.NotNil(t, err)
}

func TestCreateDuplicatesFails(t *testing.T) {
	require.Nil(t, store.Create(connection))
	require.NotNil(t, store.Create(connection))
	require.Nil(t, store.Delete(connection.ID))
}
