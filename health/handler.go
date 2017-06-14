package health

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/metrics"
)

type Handler struct {
	Metrics *metrics.MetricsManager
	H       *herodot.JSONWriter
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET("/health", h.Health)
}

// swagger:route GET /health health
//
// Check health status of instance
//
//     Responses:
//       204: emptyResponse
//       500: genericError
func (h *Handler) Health(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.Metrics.Lock()
	defer h.Metrics.Unlock()

	h.Metrics.Update()
	h.H.Write(rw, r, h.Metrics)
}
