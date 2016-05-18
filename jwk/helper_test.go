package jwk

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestToX509PEMKeyPair(t *testing.T) {
	generator := &RS256Generator{}
	keys, err := generator.Generate("")
	key := keys.Key("private")[0]
	require.Nil(t, err, "%s", err)
	_, _, err = ToX509PEMKeyPair(key.Key)
	require.Nil(t, err, "%s", err)
}