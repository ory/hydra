package policy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory/herodot"
	"github.com/ory/ladon"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"strconv"
)

const (
	endpoint         = "/policies"
	scope            = "hydra.policies"
	policyResource   = "rn:hydra:policies"
	policiesResource = "rn:hydra:policies:%s"
)

type Handler struct {
	Manager ladon.Manager
	H       herodot.Writer
	W       firewall.Firewall
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(endpoint, h.Create)
	r.GET(endpoint, h.Find)
	r.GET(endpoint+"/:id", h.Get)
	r.PUT(endpoint+"/:id", h.Update)
	r.DELETE(endpoint+"/:id", h.Delete)
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()
	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: policyResource,
		Action:   "find",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	val := r.URL.Query().Get("offset")
	if val == "" {
		val = "0"
	}

	offset, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	val = r.URL.Query().Get("limit")
	if val == "" {
		val = "500"
	}

	limit, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	policies, err := h.Manager.GetAll(offset, limit)
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	h.H.Write(w, r, policies)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p = ladon.DefaultPolicy{
		Conditions: ladon.Conditions{},
	}
	ctx := r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: policyResource,
		Action:   "create",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if p.ID == "" {
		p.ID = uuid.New()
	}

	if err := h.Manager.Create(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	h.H.WriteCreated(w, r, "/policies/"+p.ID, &p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(policiesResource, ps.ByName("id")),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	policy, err := h.Manager.Get(ps.ByName("id"))
	if err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}
	h.H.Write(w, r, policy)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(policiesResource, id),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.Delete(id); err != nil {
		h.H.WriteError(w, r, errors.New("Could not delete client"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	var p = ladon.DefaultPolicy{Conditions: ladon.Conditions{}}
	var ctx = r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(policiesResource, id),
		Action:   "update",
	}, scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if p.ID != id {
		h.H.WriteErrorCode(w, r, http.StatusBadRequest, errors.New("Payload ID does not match ID from URL"))
		return
	}

	if err := h.Manager.Delete(p.ID); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.Manager.Create(&p); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	h.H.Write(w, r, p)
}
