package pkg

import "testing"
import "github.com/oleiade/reflections"
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/require"

func AssertObjectKeysEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		assert.Equal(t, c, d, "%s", k)
	}
}

func AssertObjectKeysNotEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		assert.NotEqual(t, c, d, "%s", k)
	}
}

func RequireObjectKeysEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		require.Equal(t, c, d, "%s", k)
	}
}
func RequireObjectKeysNotEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		require.NotEqual(t, c, d, "%s", k)
	}
}
