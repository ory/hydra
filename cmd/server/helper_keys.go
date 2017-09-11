package server

import (
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

func createRS256KeysIfNotExist(c *config.Config, set, kid, use string) {
	ctx := c.Context()
	generator := jwk.RS256Generator{}

	if _, err := ctx.KeyManager.GetKey(set, kid); errors.Cause(err) == pkg.ErrNotFound {
		c.GetLogger().Infof("Key pair for signing %s is missing. Creating new one.", set)

		keys, err := generator.Generate("")
		pkg.Must(err, "Could not generate %s key: %s", set, err)

		for i, k := range keys.Keys {
			k.Use = use
			keys.Keys[i] = k
		}
		err = ctx.KeyManager.AddKeySet(set, keys)
		pkg.Must(err, "Could not persist %s key: %s", set, err)
	}
}

func publicKey(key interface{}) interface{} {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
