package server

import (
	"github.com/ory-am/hydra/jwk"
	"crypto/rsa"
	"crypto/ecdsa"
	"github.com/ory-am/hydra/pkg"
	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/config"
)

func (h *Handler) createRS256KeysIfNotExist(c *config.Config, set, lookup string) {
	ctx := c.Context()
	generator := jwk.RS256Generator{}

	if _, err := ctx.KeyManager.GetKey(set, lookup); errors.Is(err, pkg.ErrNotFound) {
		logrus.Infof("Key pair for signing %s is missing. Creating new one.", set)

		keys, err := generator.Generate("")
		pkg.Must(err, "Could not generate %s key: %s", set, err)

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
