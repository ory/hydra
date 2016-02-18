package rethinkdb

import (
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/ory-am/dockertest.v2"

	rdb "github.com/dancannon/gorethink"

	"github.com/ory-am/common/pkg"
	. "github.com/ory-am/hydra/oauth/connection"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var session *rdb.Session
var store *Store

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToRethinkDB(20, time.Second, func(url string) bool {
		rdbSession, err := rdb.Connect(rdb.ConnectOpts{
			Address:  url,
			Database: "hydra"})
		if err != nil {
			return false
		}

		_, err = rdb.DBCreate("hydra").RunWrite(rdbSession)
		if err != nil {
			return false
		}

		store = New(rdbSession)

		if err := store.CreateTables(); err != nil {
			return false
		}

		session = rdbSession

		return true
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(session, c)

	os.Exit(retCode)
}

func tearDown(session *rdb.Session, c dockertest.ContainerID) {
	defer session.Close()
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
