package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	var _, err = Connect(
		ClientID("client-id"),
		ClientSecret("client-secret"),
		ClusterURL("https://localhost:4444"),
	)
	assert.NotNil(t, err)
}
