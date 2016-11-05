package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"
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
	case *config.RethinkDBConnection:
		con.CreateTableIfNotExists("hydra_clients")
		m := &client.RethinkManager{
			Session: con.GetSession(),
			Table:   r.Table("hydra_clients"),
			Hasher:  ctx.Hasher,
		}
		if err := m.ColdStart(); err != nil {
			logrus.Fatalf("Could not fetch initial state: %s", err)
		}
		m.Watch(context.Background())
		return m
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
