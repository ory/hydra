package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/health"
)

func newHealthHandler(c *config.Config, router *httprouter.Router) *health.Handler {
	h := &health.Handler{}
	h.SetRoutes(router)
	return h
}
