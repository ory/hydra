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
	//"reflect"
	"github.com/ory-am/hydra/account"
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
	_, err := store.Get(uuid.New())
	assert.Equal(t, pkg.ErrNotFound, err)
	_, err = store.UpdateData(uuid.New(), account.UpdateDataRequest{Data: "{}"})
	assert.Equal(t, pkg.ErrNotFound, err)
}

func TestCreateAndGetCases(t *testing.T) {
	a := uuid.New()
	b := uuid.New()
	for _, c := range []struct {
		p    account.CreateAccountRequest
		pass bool
		find bool
	}{
		{
			p: account.CreateAccountRequest{
				ID:       a,
				Username: "1@bar",
				Password: "secret",
				Data:     `{"foo": "bar"}`,
			},
			pass: true, find: true,
		},
		{
			p: account.CreateAccountRequest{
				ID:       a,
				Username: "1@bar",
				Password: "secret",
				Data:     `{"foo": "bar"}`,
			},
			pass: false, find: true,
		},
		{
			p: account.CreateAccountRequest{
				ID:       b,
				Username: "1@bar",
				Password: "secret",
				Data:     `{"foo": "bar"}`,
			},
			pass: false, find: false,
		},
		{
			p: account.CreateAccountRequest{
				ID:       uuid.New(),
				Username: uuid.New(),
				Password: "secret",
				Data:     "",
			},
			pass: true, find: true,
		},
	} {
		result, err := store.Create(c.p)
		if c.pass {
			assert.Nil(t, err)
			pkg.AssertObjectKeysEqual(t, c.p, result, "ID", "Username")

			result, err = store.Get(c.p.ID)
			assert.Equal(t, c.find, err == nil)
			if !c.find {
				continue
			}
			pkg.AssertObjectKeysEqual(t, c.p, result, "ID", "Username")
		} else {
			assert.NotNil(t, err)

			result, err = store.Get(c.p.ID)
			assert.Equal(t, c.find, err == nil)
			if !c.find {
				continue
			}
			pkg.AssertObjectKeysEqual(t, c.p, result, "ID", "Username")
		}
	}
}

func TestDelete(t *testing.T) {
	id := uuid.New()
	_, _ = store.Create(account.CreateAccountRequest{
		ID:       id,
		Username: uuid.New(),
		Password: "secret",
	})

	_, err := store.Get(id)
	assert.Nil(t, err)
	assert.Nil(t, store.Delete(id))
	_, err = store.Get(id)
	assert.NotNil(t, err)
}

func TestUpdateUsername(t *testing.T) {
	id := uuid.New()
	_, _ = store.Create(account.CreateAccountRequest{
		ID:       id,
		Username: uuid.New(),
		Password: "secret",
	})

	_, err := store.UpdateUsername(id, account.UpdateUsernameRequest{
		Username: uuid.New(),
		Password: "wrong secret",
	})
	assert.NotNil(t, err)

	r, err := store.UpdateUsername(id, account.UpdateUsernameRequest{
		Username: "3@foo",
		Password: "secret",
	})
	assert.Nil(t, err)
	assert.Equal(t, "3@foo", r.GetUsername())

	// Did it persist?
	r, err = store.Get(id)
	assert.Nil(t, err)
	assert.Equal(t, "3@foo", r.GetUsername())
}

func TestUpdatePassword(t *testing.T) {
	id := uuid.New()
	ac, _ := store.Create(account.CreateAccountRequest{
		ID:       id,
		Username: uuid.New(),
		Password: "secret",
	})

	_, err := store.UpdatePassword(id, account.UpdatePasswordRequest{
		CurrentPassword: "wrong old secret",
		NewPassword:     "new secret",
	})
	assert.NotNil(t, err)

	updatedAccount, err := store.UpdatePassword(id, account.UpdatePasswordRequest{
		CurrentPassword: "secret",
		NewPassword:     "new secret",
	})
	assert.Nil(t, err)

	resultAccount, err := store.Get(id)
	assert.Nil(t, err)

	assert.Equal(t, updatedAccount.GetPassword(), resultAccount.GetPassword())
	assert.NotEqual(t, ac.GetPassword(), resultAccount.GetPassword())
}

func TestAuthenticate(t *testing.T) {
	acc, _ := store.Create(account.CreateAccountRequest{
		ID:       uuid.New(),
		Username: "5@bar",
		Password: "secret",
	})

	_, err := store.Authenticate("5@bar", "wrong secret")
	assert.NotNil(t, err)
	_, err = store.Authenticate("doesnotexist@foo", "secret")
	assert.NotNil(t, err)
	_, err = store.Authenticate("", "")
	assert.NotNil(t, err)
	result, err := store.Authenticate("5@bar", "secret")
	assert.Nil(t, err)

	assert.Equal(t, acc, result)
}
