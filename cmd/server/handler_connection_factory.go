package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/herodot"
)

func newConnectionHandler(c *config.Config, router *httprouter.Router) *connection.Handler {
	ctx := c.Context()

	h := &connection.Handler{}
	h.H = &herodot.JSON{}
	h.W = ctx.Warden
	h.SetRoutes(router)

	switch ctx.Connection.(type) {
	case *config.MemoryConnection:
		h.Manager = connection.NewMemoryManager()
		break
	default:
		panic("Unknown connection type.")
	}

	return h
}
