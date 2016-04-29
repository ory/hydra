package key

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/ory-am/hydra/pkg"
	"github.com/stretchr/testify/assert"
)

func TestRSAPEMStrategy(t *testing.T) {
	s := &RSAPEMStrategy{}
	key, err := s.AsymmetricKey("foo")
	pkg.RequireError(t, false, err)
	assert.Equal(t, "foo", key.ID)

	_, err = jwt.ParseRSAPublicKeyFromPEM(key.Public)
	pkg.RequireError(t, false, err)

	_, err = jwt.ParseRSAPrivateKeyFromPEM(key.Private)
	pkg.RequireError(t, false, err)
}
