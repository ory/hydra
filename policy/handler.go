package policy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/ory-am/hydra/config"
)

const (
	endpoint         = "/policies"
	scope            = "hydra.policies.search"
	policyResource   = "rn:hydra:policies"
	policiesResource = "rn:hydra:policies:%s"
)

type Handler struct {
	Manager ladon.Manager
	H       herodot.Herodot
	W       firewall.Firewall
}

func NewHandler(c *config.Config, router *httprouter.Router) *Handler {
	ctx := c.Context()

	h := &Handler{
		H: &herodot.JSON{},
		W: ctx.Warden,
		Manager: ctx.LadonManager,
	}
	h.SetRoutes(router)

	return h
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(endpoint, h.Create)
	r.GET(endpoint, h.Find)
	r.GET(endpoint+"/:id", h.Get)
	r.DELETE(endpoint+"/:id", h.Delete)
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var subject = r.URL.Query().Get("subject")
	var ctx = herodot.NewContext()
	if subject == "" {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.New("Missing query parameter subject"))
	}

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: policyResource,
		Action:   "search",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	policies, err := h.Manager.FindPoliciesForSubject(subject)
	if err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}
	h.H.Write(ctx, w, r, policies)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p ladon.DefaultPolicy
	ctx := herodot.NewContext()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: policyResource,
		Action:   "create",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}

	if p.ID == "" {
		p.ID = uuid.New()
	}
	if err := h.Manager.Create(&p); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}

	h.H.WriteCreated(ctx, w, r, "/policies/"+p.ID, p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := herodot.NewContext()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(policiesResource, ps.ByName("id")),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	policy, err := h.Manager.Get(ps.ByName("id"))
	if err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}
	h.H.Write(ctx, w, r, policy)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := herodot.NewContext()
	id := ps.ByName("id")

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(policiesResource, id),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.Delete(id); err != nil {
		h.H.WriteError(ctx, w, r, errors.New("Could not delete client"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
