package connection_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/connection/postgres"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ory-am/dockertest.v2"
)

var connections = []connection.Connection{
	&connection.DefaultConnection{
		ID:            uuid.New(),
		LocalSubject:  "peter",
		RemoteSubject: "peterson",
		Provider:      "google",
	},
}

var managers = map[string]connection.Manager{}

var containers = []dockertest.ContainerID{}

func TestMain(m *testing.M) {
	setupPG()

	retCode := m.Run()
	for _, c := range containers {
		c.KillRemove()
	}

	os.Exit(retCode)
}

func setupPG() {
	var db *sql.DB
	var err error

	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	containers = append(containers, c)

	m := &postgres.Manager{DB: db}
	if err = m.CreateSchemas(); err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}

	managers["postgres"] = m
}

func TestNotFound(t *testing.T) {
	for _, store := range managers {
		_, err := store.Get("asdf")
		assert.EqualError(t, err, pkg.ErrNotFound.Error())
	}
}

func TestCreateDuplicatesFails(t *testing.T) {
	for _, store := range managers {
		require.Nil(t, store.Create(connections[0]))
		require.NotNil(t, store.Create(connections[0]))
		require.Nil(t, store.Delete(connections[0].GetID()))
	}
}

func TestCreateGetFindDelete(t *testing.T) {
	for _, store := range managers {
		for _, c := range connections {
			require.Nil(t, store.Create(c))

			res, err := store.Get(c.GetID())
			require.Nil(t, err)
			require.Equal(t, c, res)

			res, err = store.FindByRemoteSubject("google", "peterson")
			require.Nil(t, err)
			require.Equal(t, c, res)

			cs, err := store.FindAllByLocalSubject("peter")
			require.Nil(t, err)
			assert.Len(t, cs, 1)
			require.Equal(t, c, cs[0])

			require.Nil(t, store.Delete(c.GetID()))
			_, err = store.Get(c.GetID())
			require.NotNil(t, err)
		}
	}
}
