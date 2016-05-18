package jwk

import (
	"testing"

	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	for _, c := range []struct {
		g     KeyGenerator
		check func(*jose.JsonWebKeySet)
	}{
		{
			g: &RS256Generator{},
			check: func(ks *jose.JsonWebKeySet) {
				assert.Len(t, ks, 2)
			},
		},
		{
			g: &ECDSA521Generator{},
			check: func(ks *jose.JsonWebKeySet) {
				assert.Len(t, ks, 2)
			},
		},
		{
			g: &ECDSA256Generator{},
			check: func(ks *jose.JsonWebKeySet) {
				assert.Len(t, ks, 2)
			},
		},
		{
			g: &HS256Generator{
				Length: 32,
			},
			check: func(ks *jose.JsonWebKeySet) {
				assert.Len(t, ks, 1)
			},
		},
	} {
		keys, err := c.g.Generate("foo")
		pkg.AssertError(t, false, err)
		if err != nil {
			c.check(keys)
		}
	}
}
