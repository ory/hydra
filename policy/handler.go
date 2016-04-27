package policy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
)

type Handler struct {
	Manager ladon.Manager
	H       herodot.Herodot
	W       warden.Warden
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/policies", h.Create)
	r.GET("/policies", h.Find)
	r.GET("/policies/:id", h.Get)
	r.DELETE("/policies/:id", h.Delete)
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var subject = r.URL.Query().Get("subject")
	var ctx = herodot.NewContext()
	if subject == "" {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.New("Missing query parameter subject"))
	}

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: "rn:hydra:policies",
		Action:   "search",
	}, "hydra.policies.search"); err != nil {
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
		Resource: "rn:hydra:policies",
		Action:   "create",
	}, "hydra.policies.create"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}

	p.ID = uuid.New()
	if err := h.Manager.Create(&p); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}

	h.H.WriteCreated(ctx, w, r, "/policies/"+p.ID, p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := herodot.NewContext()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf("rn:hydra:policies:%s", ps.ByName("id")),
		Action:   "get",
	}, "hydra.policies.get"); err != nil {
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

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf("rn:hydra:policies:%s", ps.ByName("id")),
		Action:   "get",
	}, "hydra.policies.delete"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.Delete(ps.ByName("id")); err != nil {
		h.H.WriteError(ctx, w, r, errors.New("Could not delete client"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
