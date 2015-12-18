package postgres

import (
	"database/sql"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/hydra/hash"
	"github.com/pborman/uuid"
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

	store = New(&hash.BCrypt{10}, db)
	if err := store.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	os.Exit(m.Run())
}

func TestNotFound(t *testing.T) {
	_, err := store.Get("asdf")
	t.Log("Got error %s", err)
	assert.NotNil(t, err)
	assert.Equal(t, pkg.ErrNotFound, err)

	_, err = store.UpdateData("asdf", "{}")
	assert.NotNil(t, err)
	assert.Equal(t, pkg.ErrNotFound, err)
}

func TestCreateAndGetCases(t *testing.T) {
	a := uuid.New()
	b := uuid.New()
	for _, c := range []struct {
		data  []string
		extra string
		pass  bool
		find  bool
	}{
		{[]string{a, "1@bar", "secret"}, `{"foo": "bar"}`, true, true},
		{[]string{a, "1@foo", "secret"}, `{"foo": "bar"}`, false, true},
		{[]string{b, "1@bar", "secret"}, `{"foo": "bar"}`, false, false},
	} {
		result, err := store.Create(c.data[0], c.data[1], c.data[2], c.extra)
		if c.pass {
			assert.Nil(t, err)
			assert.Equal(t, c.data[0], result.GetID())
			assert.Equal(t, c.data[1], result.GetUsername())
			assert.NotEqual(t, c.data[2], result.GetPassword())
			assert.Equal(t, c.extra, result.GetData())

			result, err = store.Get(c.data[0])
			if c.find {
				assert.Nil(t, err)
				assert.Equal(t, c.data[0], result.GetID())
				assert.Equal(t, c.data[1], result.GetUsername())
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
	id := uuid.New()
	_, err := store.Create(id, "2@bar", "secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.Get(id)
	assert.Nil(t, err)

	err = store.Delete(id)
	assert.Nil(t, err)

	_, err = store.Get(id)
	assert.NotNil(t, err)
}

func TestUpdateUsername(t *testing.T) {
	id := uuid.New()
	_, err := store.Create(id, "3@bar", "secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.UpdateUsername(id, "3@foo", "wrong secret")
	assert.NotNil(t, err)

	_, err = store.UpdateUsername(id, "3@foo", "secret")
	assert.Nil(t, err)

	r, err := store.Get(id)
	assert.Nil(t, err)

	assert.Equal(t, id, r.GetID())
	assert.Equal(t, "3@foo", r.GetUsername())
	assert.NotEqual(t, "secret", r.GetPassword())
}

func TestUpdatePassword(t *testing.T) {
	id := uuid.New()
	account, err := store.Create(id, "4@bar", "old secret", `{"foo": "bar"}`)
	assert.Nil(t, err)

	_, err = store.UpdatePassword(id, "wrong old secret", "new secret")
	assert.NotNil(t, err)

	updatedAccount, err := store.UpdatePassword(id, "old secret", "new secret")
	assert.Nil(t, err)

	resultAccount, err := store.Get(id)
	assert.Nil(t, err)

	assert.Equal(t, updatedAccount.GetPassword(), resultAccount.GetPassword())
	assert.NotEqual(t, account.GetPassword(), resultAccount.GetPassword())
}

func TestAuthenticate(t *testing.T) {
	account, err := store.Create(uuid.New(), "5@bar", "secret", `{"foo": "bar"}`)
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
