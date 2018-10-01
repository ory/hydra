package health

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/metrics"
)

type Handler struct {
	Metrics *metrics.MetricsManager
	H       *herodot.JSONWriter
	W       firewall.Firewall
	Config  *config.Config
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET("/health", h.Health)
	r.GET("/health/detailed", h.DetailedHealth)
	r.GET("/health/stats", h.Statistics)
}

// swagger:route GET /health health
//
// Check health status of instance
//
// This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) Health(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write(simpleStatus(h.Config))
}

// swagger:route GET /health/detailed health detailed
//
// Check health status of instance with detailed information including dependencies' health
//
// This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) DetailedHealth(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write(detailedStatus(h.Config))
}

// swagger:route GET /health/stats health getStatistics
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
