package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory/herodot"
	"github.com/ory-am/hydra/policy"
)

func newPolicyHandler(c *config.Config, router *httprouter.Router) *policy.Handler {
	ctx := c.Context()
	h := &policy.Handler{
		H:        herodot.NewJSONWriter(c.GetLogger()),
		W:       ctx.Warden,
		Manager: ctx.LadonManager,
	}
	h.SetRoutes(router)
	return h
}
