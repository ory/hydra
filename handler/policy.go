// Package handler
//
// Defined permissions:
// * permission("rn:hydra:policies") actions("create")
// * permission("rn:hydra:policies:%s", id) actions("get", "delete")
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/ory-am/common/pkg"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"github.com/ory-am/ladon"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/julienschmidt/httprouter"
)

type PolicyHandler struct {
	policyManager ladon.Manager
	h             herodot.Herodot
	w             warden.Warden
}

func (h *PolicyHandler) SetRoutes(r *httprouter.Router) {
	r.POST("/policies", h.Create)

	r.GET("/policies/:id", h.Get)
	r.DELETE("/policies/:id", h.Delete)
}

func (h *PolicyHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p ladon.DefaultPolicy
	ctx := herodot.NewContext()

	h.w.HTTPAuthorized(r, "hydra.policies")

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.h.WriteErrorCode(ctx, w, r, http.StatusInternalServerError, errors.New(err))
		return
	}

	p.ID = uuid.New()
	if err := h.policyManager.Create(&p); err != nil {
		h.h.WriteErrorCode(ctx, w, r, http.StatusInternalServerError, errors.New(err))
		return
	}

	h.h.WriteCreated(ctx, w, r, "/policies/" + p.ID, p)
}

func (h *PolicyHandler) Get(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := herodot.NewContext()

	policy, err := h.policyManager.Get(ps.ByName("id"))
	if err != nil {
		h.h.WriteErrorCode(ctx, rw, errors.New(err), http.StatusInternalServerError)
		return
	}
	h.h.Write(ctx, rw, req, policy)
}

func (h *PolicyHandler) Delete(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := herodot.NewContext()

	if err := h.policyManager.Delete(ps.ByName("id")); err != nil {
		h.h.WriteErrorCode(ctx, rw, errors.Errorf("Could not delete client"), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}
