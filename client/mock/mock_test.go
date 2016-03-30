package mock

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientAlwaysTrue(t *testing.T) {
	c := NewAlwaysTrue()
	ok, err := c.IsActionAllowed(nil)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = c.IsRequestAllowed(nil, "", "", "")
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = c.IsAuthenticated("")
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestClientAlwaysFalse(t *testing.T) {
	c := NewAlwaysFalse()
	ok, err := c.IsActionAllowed(nil)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = c.IsRequestAllowed(nil, "", "", "")
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = c.IsAuthenticated("")
	assert.False(t, ok)
	assert.NotNil(t, err)
}
