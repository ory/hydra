package health

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/metrics"
)

type Handler struct {
	Metrics *metrics.MetricsManager
	H       *herodot.JSONWriter
	W       firewall.Firewall
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET("/health", h.Health)
	r.GET("/health/stats", h.Statistics)
}

// swagger:route GET /health health getHealthStatus
//
// Check health status of instance
//
// This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) Health(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write([]byte(`{"status": "ok"}`))
}

// swagger:route GET /health/stats health getInstanceStatistics
//
// Show instance statistics
//
// This endpoint returns information on the instance's health. It is currently not documented.
//
// The subject making the request needs to be assigned to a policy containing:
//
// ```
// {
//   "resources": ["rn:hydra:health:stats"],
//   "actions": ["get"],
//   "effect": "allow"
// }
// ```
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.health
//
//     Responses:
//       200: emptyResponse
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Statistics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:health:stats",
		Action:   "get",
	}, "hydra.health"); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.Metrics.Lock()
	defer h.Metrics.Unlock()

	h.Metrics.Update()
	h.H.Write(w, r, h.Metrics)
}
