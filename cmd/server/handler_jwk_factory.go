package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/jwk"
	"golang.org/x/net/context"
)

func newJWKHandler(c *config.Config, router *httprouter.Router) *jwk.Handler {
	ctx := c.Context()
	h := &jwk.Handler{
		H: &herodot.JSON{},
		W: ctx.Warden,
	}
	h.SetRoutes(router)

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		ctx.KeyManager = &jwk.MemoryManager{}
		h.Manager = ctx.KeyManager
		break
	case *config.RethinkDBConnection:
		con.CreateTableIfNotExists("hydra_policies")
		m := &jwk.RethinkManager{Session: con.GetSession()}
		m.ColdStart()
		m.Watch(context.Background())
		h.Manager = m
		break
	default:
		panic("Unknown connection type.")
	}

	return h
}
