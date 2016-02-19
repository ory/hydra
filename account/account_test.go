package account

import (
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountCases(t *testing.T) {
	for _, c := range []DefaultAccount{
		{ID: uuid.New(), Username: "foo@bar", Password: "secret", Data: ""},
		{ID: "2", Username: "baz@bar", Password: "top secret", Data: `{"foo": "bar"}`},
	} {
		assert.Equal(t, c.ID, c.GetID())
		assert.Equal(t, c.Username, c.GetUsername())
		assert.Equal(t, c.Password, c.GetPassword())
		assert.Equal(t, c.Data, c.GetData())
	}
}
