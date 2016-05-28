package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/jwk"
	"golang.org/x/net/context"
	"github.com/square/go-jose"
	"github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
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
		break
	case *config.RethinkDBConnection:
		con.CreateTableIfNotExists("hydra_json_web_keys")
		m := &jwk.RethinkManager{
			Session: con.GetSession(),
			Keys: map[string]jose.JsonWebKeySet{},
			Table: r.Table("hydra_json_web_keys"),
		}
		if err := m.ColdStart(); err != nil {
			logrus.Fatalf("Could not fetch initial state: %s", err)
		}
		m.Watch(context.Background())
		ctx.KeyManager = m
		break
	default:
		logrus.Fatalf("Unknown connection type.")
	}

	h.Manager = ctx.KeyManager
	return h
}
