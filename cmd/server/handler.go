package server

import (
	"github.com/ory-am/hydra/client"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
)

type Handler struct {
	Clients *client.Handler
}

func (h *Handler) Listen(c *config.Config, router *httprouter.Router) {
	h.Clients.Listen(c, router)
}
