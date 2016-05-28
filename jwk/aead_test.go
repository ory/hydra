package jwk

import (
	"testing"
	"github.com/ory-am/fosite/rand"
	"github.com/ory-am/hydra/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/pborman/uuid"
)

func TestAEAD(t *testing.T) {
	key, err := rand.RandomBytes(32)
	pkg.AssertError(t, false, err)

	a := &AEAD{
		Key: key,
	}

	for i := 0; i < 100; i++ {
		plain := []byte(uuid.New())
		ct, err := a.Encrypt(plain)
		pkg.AssertError(t, false, err)
		t.Logf("Ciphertext: %s", ct)

		res, err := a.Decrypt(ct)
		pkg.AssertError(t, false, err)
		assert.Equal(t, plain, res)
	}
}