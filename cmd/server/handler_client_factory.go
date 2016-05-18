package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
)

func newClientManager(c *config.Config) client.Manager {
	ctx := c.Context()

	switch ctx.Connection.(type) {
	case *config.MemoryConnection:
		return &client.MemoryManager{
			Clients: map[string]*fosite.DefaultClient{},
			Hasher:  ctx.Hasher,
		}
	default:
		panic("Unknown connection type.")
	}
}

func newClientHandler(c *config.Config, router *httprouter.Router, manager client.Manager) *client.Handler {
	ctx := c.Context()
	h := &client.Handler{
		H: &herodot.JSON{},
		W: ctx.Warden, Manager: manager,
	}

	h.SetRoutes(router)
	return h
}
