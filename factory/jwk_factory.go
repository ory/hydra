package factory

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/jwk"
)

func NewJWKHandler(c *config.Config, router *httprouter.Router) *jwk.Handler {
	ctx := c.Context()
	h := &jwk.Handler{
		H: &herodot.JSON{},
		W: ctx.Warden,
	}
	h.SetRoutes(router)

	switch ctx.Connection.(type) {
	case *config.MemoryConnection:
		ctx.KeyManager = &jwk.MemoryManager{}
		h.Manager = ctx.KeyManager
		break
	default:
		panic("Unknown connection type.")
	}

	return h
}
