package key_test

import (
	"testing"

	"github.com/ory-am/hydra/key"
	"github.com/ory-am/hydra/pkg"
	"github.com/stretchr/testify/assert"
)

var managers = map[string]key.Manager{
	"memory": &key.MemoryManager{
		AsymmetricKeys: map[string]*key.AsymmetricKey{},
		SymmetricKeys:  map[string]*key.SymmetricKey{},
		Strategy: &key.DefaultKeyStrategy{
			AsymmetricKeyStrategy: &key.RSAPEMStrategy{},
			SymmetricKeyStrategy:  &key.SHAStrategy{},
		},
	},
}

func init() {

}

func TestAsymmetricManager(t *testing.T) {
	for k, m := range managers {
		_, err := m.GetAsymmetricKey("foo")
		pkg.AssertError(t, true, err, k)

		original, err := m.CreateAsymmetricKey("foo")
		pkg.AssertError(t, false, err, k)
		assert.NotEmpty(t, original.Private, k)
		assert.NotEmpty(t, original.Public, k)

		_, err = m.CreateAsymmetricKey("foo")
		pkg.AssertError(t, false, err, k)

		err = m.DeleteAsymmetricKey("foo")
		pkg.AssertError(t, false, err, k)

		_, err = m.GetAsymmetricKey("foo")
		pkg.AssertError(t, true, err, k)

	}
}

func TestSymmetricManager(t *testing.T) {
	for k, m := range managers {
		_, err := m.GetSymmetricKey("foo")
		pkg.AssertError(t, true, err, k)

		original, err := m.CreateSymmetricKey("foo")
		pkg.AssertError(t, false, err, k)
		assert.NotEmpty(t, original.Key, k)

		_, err = m.CreateSymmetricKey("foo")
		pkg.AssertError(t, false, err, k)

		err = m.DeleteSymmetricKey("foo")
		pkg.AssertError(t, false, err, k)

		_, err = m.GetSymmetricKey("foo")
		pkg.AssertError(t, true, err, k)

	}
}
