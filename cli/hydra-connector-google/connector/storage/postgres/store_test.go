package postgres

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/handler/connector/storage"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ory-am/dockertest.v2"
)

var db *sql.DB
var store *Store

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
