package jwt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map{"foo": "bar", "1": float64(2)}
	enc, err := m.Marshall()
	require.Nil(t, err)
	dec, err := Unmarshal(enc)
	require.Nil(t, err)
	assert.Equal(t, m, dec, "%v does not equal %v", m, dec)
	assert.Equal(t, m["foo"], dec["foo"])
}
