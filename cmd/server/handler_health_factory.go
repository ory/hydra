package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/health"
)

func newHealthHandler(c *config.Config, router *httprouter.Router) *health.Handler {
	h := &health.Handler{
		Metrics: c.GetMetrics(),
		H:       herodot.NewJSONWriter(c.GetLogger()),
		W:       c.Context().Warden,
	}
	h.SetRoutes(router)
	return h
}
