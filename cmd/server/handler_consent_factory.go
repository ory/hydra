package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/oauth2"
)

func injectConsentManager(c *config.Config) {
	var ctx = c.Context()
	var manager oauth2.ConsentRequestManager

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		manager = oauth2.NewConsentRequestMemoryManager()
		break
	case *config.SQLConnection:
		manager = oauth2.NewConsentRequestSQLManager(con.GetDatabase())
		break
	case *config.PluginConnection:
		var err error
		if manager, err = con.NewConsentRequestManager(); err != nil {
			c.GetLogger().Fatalf("Could not load client manager plugin %s", err)
		}
		break
	default:
		panic("Unknown connection type.")
	}

	ctx.ConsentManager = manager

}

func newConsentHanlder(c *config.Config, router *httprouter.Router) *oauth2.ConsentSessionHandler {
	ctx := c.Context()
	h := &oauth2.ConsentSessionHandler{
		H: herodot.NewJSONWriter(c.GetLogger()),
		W: ctx.Warden, M: ctx.ConsentManager,
	}

	h.SetRoutes(router)
	return h
}
