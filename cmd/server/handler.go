package server

import (
	"github.com/ory-am/hydra/client"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/oauth2"
)

type Handler struct {
	Clients *client.Handler
	OAuth2 *oauth2.Handler
}

func (h *Handler) Listen(c *config.Config, router *httprouter.Router) {
	h.Clients = client.NewHandler(c, router)
	h.OAuth2 = oauth2.NewHandler(c, router, h.Clients.Manager)
}
