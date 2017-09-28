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
	r.GET("/health/status", h.Health)
	r.GET("/health/metrics", h.Statistics)
}

// swagger:route GET /health/status health getInstanceStatus
//
// Check health status of this instance
//
// This endpoint returns `{ "status": "ok" }`. This status let's you know that the HTTP server is up and running. This
// status does currently not include checks whether the database connection is up and running. This endpoint does not
// require the `X-Forwarded-Proto` header when TLS termination is set.
//
//
// Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state,
// only to a single instance.
//
//     Responses:
//       200: healthStatus
//       500: genericError
func (h *Handler) Health(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rw.Write([]byte(`{"status": "ok"}`))
}

// swagger:route GET /health/metrics health getInstanceMetrics
//
// Show instance metrics (experimental)
//
// This endpoint returns an instance's metrics, such as average response time, status code distribution, hits per
// second and so on. The return values are currently not documented as this endpoint is still experimental.
//
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
