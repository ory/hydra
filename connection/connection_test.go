package connection_test

import (
	"testing"

	"github.com/ory-am/hydra/connection"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	c := &connection.DefaultConnection{
		ID:            uuid.New(),
		LocalSubject:  "peter",
		RemoteSubject: "peter@gmail.com",
		Provider:      "google",
	}

	assert.Equal(t, c.ID, c.GetID())
	assert.Equal(t, c.Provider, c.GetProvider())
	assert.Equal(t, c.LocalSubject, c.GetLocalSubject())
	assert.Equal(t, c.RemoteSubject, c.GetRemoteSubject())
}
