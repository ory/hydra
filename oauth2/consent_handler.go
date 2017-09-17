package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/pkg/errors"
)

const (
	ConsentRequestAccepted = "accepted"
	ConsentRequestRejected = "rejected"

	ConsentSessionPath = "/oauth2/consent/requests"

	ConsentResource = "rn:hydra:oauth2:consent:requests:%s"
	ConsentScope    = "hydra.consent"
)

type ConsentSessionHandler struct {
	H herodot.Writer
	M ConsentRequestManager
	W firewall.Firewall
}

func (h *ConsentSessionHandler) SetRoutes(r *httprouter.Router) {
	r.GET(ConsentSessionPath+"/:id", h.FetchConsentRequest)
	r.PATCH(ConsentSessionPath+"/:id/reject", h.RejectConsentRequestHandler)
	r.PATCH(ConsentSessionPath+"/:id/accept", h.AcceptConsentRequestHandler)
}

// swagger:route GET /.well-known/openid-configuration oauth2 openid-connect WellKnownHandler
//
// Server well known configuration
//
// For more information, please refer to https://openid.net/specs/openid-connect-discovery-1_0.html
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2:
//
//     Responses:
//       200: WellKnown
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) FetchConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := h.W.TokenAllowed(r.Context(), h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ConsentResource, ps.ByName("id")),
		Action:   "get",
	}, ConsentScope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if session, err := h.M.GetConsentRequest(ps.ByName("id")); err != nil {
		h.H.WriteError(w, r, err)
		return
	} else {
		h.H.Write(w, r, session)
	}
}

func (h *ConsentSessionHandler) RejectConsentRequestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := h.W.TokenAllowed(r.Context(), h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ConsentResource, ps.ByName("id")),
		Action:   "reject",
	}, ConsentScope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.M.RejectConsentRequest(ps.ByName("id")); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ConsentSessionHandler) AcceptConsentRequestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := h.W.TokenAllowed(r.Context(), h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ConsentResource, ps.ByName("id")),
		Action:   "accept",
	}, ConsentScope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	var payload AcceptConsentRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.M.AcceptConsentRequest(ps.ByName("id"), &payload); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
