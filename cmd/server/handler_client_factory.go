package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory/herodot"
)

func newClientManager(c *config.Config) client.Manager {
	ctx := c.Context()

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		return &client.MemoryManager{
			Clients: map[string]client.Client{},
			Hasher:  ctx.Hasher,
		}
	case *config.SQLConnection:
		m := &client.SQLManager{
			DB:     con.GetDatabase(),
			Hasher: ctx.Hasher,
		}
		if err := m.CreateSchemas(); err != nil {
			logrus.Fatalf("Could not create client schema: %s", err)
		}
		return m
	default:
		panic("Unknown connection type.")
	}
}

func newClientHandler(c *config.Config, router *httprouter.Router, manager client.Manager) *client.Handler {
	ctx := c.Context()
	h := &client.Handler{
		H: herodot.NewJSONWriter(c.GetLogger()),
		W: ctx.Warden, Manager: manager,
	}

	h.SetRoutes(router)
	return h
}
