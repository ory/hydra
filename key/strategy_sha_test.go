package key

import (
	"testing"

	"github.com/ory-am/hydra/pkg"
	"github.com/stretchr/testify/assert"
)

func TestSHAStrategy(t *testing.T) {
	s := &SHAStrategy{}
	key, err := s.SymmetricKey("foo")
	pkg.RequireError(t, false, err)
	assert.Equal(t, "foo", key.ID)
	assert.True(t, len(key.Key) > 64, "%d: %s", len(key.Key), key.Key)
}
