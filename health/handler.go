package health

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler struct {

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
	rw.WriteHeader(http.StatusNoContent)
}
