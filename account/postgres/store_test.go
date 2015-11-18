package postgres

import (
	"database/sql"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/hydra/hash"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
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

	store = New(&hash.BCrypt{10}, db)
	if err := store.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	os.Exit(m.Run())
}

func TestCreateAndGetCases(t *testing.T) {
	for _, c := range []struct {
		data  []string
		extra string
		pass  bool
		find  bool
	}{
		{[]string{"1", "1@bar", "secret"}, `{"foo": "bar"}`, true, true},
		{[]string{"1", "1@foo", "secret"}, `{"foo": "bar"}`, false, true},
		{[]string{"2", "1@bar", "secret"}, `{"foo": "bar"}`, false, false},
	} {
		result, err := store.Create(c.data[0], c.data[1], c.data[2], c.extra)
		if c.pass {
			assert.Nil(t, err)
			assert.Equal(t, c.data[0], result.GetID())
			assert.Equal(t, c.data[1], result.GetEmail())
			assert.NotEqual(t, c.data[2], result.GetPassword())
			assert.Equal(t, c.extra, result.GetData())

			result, err = store.Get(c.data[0])
			if c.find {
				assert.Nil(t, err)
				assert.Equal(t, c.data[0], result.GetID())
				assert.Equal(t, c.data[1], result.GetEmail())
				assert.NotEqual(t, c.data[2], result.GetPassword())
				assert.Equal(t, c.extra, result.GetData())
			} else {
				assert.NotNil(t, err)
			}
		} else {
			assert.NotNil(t, err)
			_, err = store.Get(c.data[0])
			if c.find {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		}
	}
}

func TestDelete(t *testing.T) {
	_, err := store.Create("2", "2@bar", "secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.Get("2")
	assert.Nil(t, err)

	err = store.Delete("2")
	assert.Nil(t, err)

	_, err = store.Get("2")
	assert.NotNil(t, err)
}

func TestUpdateEmail(t *testing.T) {
	_, err := store.Create("3", "3@bar", "secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.UpdateEmail("3", "3@foo", "wrong secret")
	assert.NotNil(t, err)

	_, err = store.UpdateEmail("3", "3@foo", "secret")
	assert.Nil(t, err)

	r, err := store.Get("3")
	assert.Nil(t, err)

	assert.Equal(t, "3", r.GetID())
	assert.Equal(t, "3@foo", r.GetEmail())
	assert.NotEqual(t, "secret", r.GetPassword())
}

func TestUpdatePassword(t *testing.T) {
	account, err := store.Create("4", "4@bar", "old secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.UpdatePassword("4", "wrong old secret", "new secret")
	assert.NotNil(t, err)

	updatedAccount, err := store.UpdatePassword("4", "old secret", "new secret")
	assert.Nil(t, err)

	resultAccount, err := store.Get("4")
	assert.Nil(t, err)

	assert.Equal(t, updatedAccount.GetPassword(), resultAccount.GetPassword())
	assert.NotEqual(t, account.GetPassword(), resultAccount.GetPassword())
}

func TestAuthenticate(t *testing.T) {
	account, err := store.Create("5", "5@bar", "secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.Authenticate("5@bar", "wrong secret")
	assert.NotNil(t, err)
	_, err = store.Authenticate("doesnotexist@foo", "secret")
	assert.NotNil(t, err)
	_, err = store.Authenticate("", "")
	assert.NotNil(t, err)

	result, err := store.Authenticate("5@bar", "secret")
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(account, result), "Results do not match: (%v) does not equal ($v)", &account, &result)
}
