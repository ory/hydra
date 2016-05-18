package jwk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustRSAPrivate(t *testing.T) {
	keys, err := new(RS256Generator).Generate("")
	assert.Nil(t, err)

	_, err = ToRSAPrivate(&keys.Key("private")[0])
	assert.Nil(t, err)

	MustRSAPrivate(&keys.Key("private")[0])

	_, err = ToRSAPublic(&keys.Key("public")[0])
	assert.Nil(t, err)
	MustRSAPublic(&keys.Key("public")[0])
}
