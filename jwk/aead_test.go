package jwk

import (
	"testing"

	"crypto/rand"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
)

// RandomBytes returns n random bytes by reading from crypto/rand.Reader
func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.Wrap(err, "")
	}
	return bytes, nil
}

func TestAEAD(t *testing.T) {
	key, err := randomBytes(32)
	pkg.AssertError(t, false, err)

	a := &AEAD{
		Key: key,
	}

	for i := 0; i < 100; i++ {
		plain := []byte(uuid.New())
		ct, err := a.Encrypt(plain)
		pkg.AssertError(t, false, err)

		res, err := a.Decrypt(ct)
		pkg.AssertError(t, false, err)
		assert.Equal(t, plain, res)
	}
}
