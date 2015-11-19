package postgres

import (
	"database/sql"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/hydra/oauth/provider/storage"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var db *sql.DB
var store *Store

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
		log.Fatalf("Could not set up schemas: %v", err)
	}
	os.Exit(m.Run())
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
