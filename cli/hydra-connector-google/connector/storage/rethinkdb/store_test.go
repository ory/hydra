package rethinkdb

import (
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/ory-am/dockertest.v2"

	rdb "github.com/dancannon/gorethink"

	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/endpoint/connector/storage"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
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

func TestGetNotFound(t *testing.T) {
	_, err := store.GetStateData("asdf")
	assert.Equal(t, pkg.ErrNotFound, err)
}

func TestCreateAndGetCases(t *testing.T) {
	var stateDataSet = map[string]storage.StateData{
		"valid": {
			ID:          uuid.New(),
			ClientID:    uuid.New(),
			RedirectURL: "http://localhost/",
			Scope:       "scope",
			State:       "state",
			Type:        "code",
			Provider:    "facebook",
			ExpiresAt:   time.Now().Add(time.Hour),
		},
		"invalid": {},
	}

	for k, c := range []struct {
		sd         storage.StateData
		passCreate bool
		passGet    bool
	}{
		{stateDataSet["valid"], true, true},
	} {
		assert.Equal(t, c.passCreate, store.SaveStateData(&c.sd) == nil, "Case %d", k)
		if !c.passCreate {
			continue
		}

		result, err := store.GetStateData(c.sd.ID)
		assert.Nil(t, err, "Case %d", k)

		assert.Equal(t, c.sd.IsExpired(), result.IsExpired(), "Case %d", k)
		c.sd.ExpiresAt = result.ExpiresAt
		assert.Equal(t, &c.sd, result, "Case %d", k)
	}
}
