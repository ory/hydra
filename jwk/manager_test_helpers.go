package jwk

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
)

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func TestHelperManagerKey(m Manager, keys *jose.JsonWebKeySet) func(t *testing.T) {
	pub := keys.Key("public")
	priv := keys.Key("private")

	return func(t *testing.T) {
		_, err := m.GetKey("faz", "baz")
		assert.NotNil(t, err)

		err = m.AddKey("faz", First(priv))
		assert.Nil(t, err)

		got, err := m.GetKey("faz", "private")
		assert.Nil(t, err)
		assert.Equal(t, priv, got.Keys)

		err = m.AddKey("faz", First(pub))
		assert.Nil(t, err)

		got, err = m.GetKey("faz", "private")
		assert.Nil(t, err)
		assert.Equal(t, priv, got.Keys)

		got, err = m.GetKey("faz", "public")
		assert.Nil(t, err)
		assert.Equal(t, pub, got.Keys)

		err = m.DeleteKey("faz", "public")
		assert.Nil(t, err)

		_, err = m.GetKey("faz", "public")
		assert.NotNil(t, err)
	}
}

func TestHelperManagerKeySet(m Manager, keys *jose.JsonWebKeySet) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := m.GetKeySet("foo")
		pkg.AssertError(t, true, err)

		err = m.AddKeySet("bar", keys)
		assert.Nil(t, err)

		got, err := m.GetKeySet("bar")
		assert.Nil(t, err)
		assert.Equal(t, keys.Key("public"), got.Key("public"))
		assert.Equal(t, keys.Key("private"), got.Key("private"))

		err = m.DeleteKeySet("bar")
		assert.Nil(t, err)

		_, err = m.GetKeySet("bar")
		assert.NotNil(t, err)
	}
}
