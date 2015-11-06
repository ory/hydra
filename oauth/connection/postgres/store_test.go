package postgres_test

import (
	"database/sql"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/dockertest"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/pborman/uuid"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/stretchr/testify/require"
	. "github.com/ory-am/hydra/oauth/connection"
	. "github.com/ory-am/hydra/oauth/connection/postgres"
	"log"
	"os"
	"testing"
	"time"
)

var store *Store
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	var c dockertest.ContainerID
	c, db, err = dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping: %s", err)
	}

	store = New(db)
	if err := store.CreateSchemas(); err != nil {
		log.Fatalf("Could not ping: %s", err)
	}
	os.Exit(m.Run())
}

var connection = &DefaultConnection{ID: uuid.New(), LocalSubject: "peter", RemoteSubject: "peterson", Provider: "google"}

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
