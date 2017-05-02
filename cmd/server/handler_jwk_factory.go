package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory/herodot"
	"github.com/ory-am/hydra/jwk"
)

func injectJWKManager(c *config.Config) {
	ctx := c.Context()

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		ctx.KeyManager = &jwk.MemoryManager{}
		break
	case *config.SQLConnection:
		m := &jwk.SQLManager{
			DB: con.GetDatabase(),
			Cipher: &jwk.AEAD{
				Key: c.GetSystemSecret(),
			},
		}
		if err := m.CreateSchemas(); err != nil {
			c.GetLogger().Fatalf("Could not create jwk schema: %s", err)
		}
		ctx.KeyManager = m
		break
	default:
		c.GetLogger().Fatalf("Unknown connection type.")
	}
}

func newJWKHandler(c *config.Config, router *httprouter.Router) *jwk.Handler {
	ctx := c.Context()
	h := &jwk.Handler{
		H:        herodot.NewJSONWriter(c.GetLogger()),
		W:       ctx.Warden,
		Manager: ctx.KeyManager,
	}
	h.SetRoutes(router)
	return h
}
