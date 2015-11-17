// Package handler
//
// Defined permissions:
// * permission("rn:hydra:policies") actions("create")
// * permission("rn:hydra:policies:%s", id) actions("get", "delete")
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	hctx "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/middleware"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"net/http"

	log "github.com/Sirupsen/logrus"
	. "github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/guard/operator"
	. "github.com/ory-am/ladon/policy"
)

type Handler struct {
	s Storage
	m *middleware.Middleware
	g Guarder
}

type payload struct {
	Resource   string            `json:"string"`
	Subject    string            `json:"subject"`
	Permission string            `json:"permission"`
	Context    *operator.Context `json:"context"`
}

func permission(id string) string {
	return fmt.Sprintf("rn:hydra:policies:%s", id)
}

func NewHandler(s Storage, m *middleware.Middleware, g Guarder) *Handler {
	return &Handler{s, m, g}
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h hctx.ContextHandler) hctx.ContextHandler) {
	r.Handle("/granted", hctx.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Granted)).Methods("POST")

	r.Handle("/policies", hctx.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
		h.m.IsAuthorized("rn:hydra:policies", "create", nil),
	).ThenFunc(h.Create)).Methods("POST")

	r.Handle("/policies/{id}", hctx.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Get)).Methods("GET")

	r.Handle("/policies/{id}", hctx.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Delete)).Methods("DELETE")
}

func (h *Handler) Granted(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var p payload
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	policies, err := h.s.FindPoliciesForSubject(p.Subject)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"resource":   p.Resource,
			"permission": p.Permission,
			"subject":    p.Subject,
			"context":    fmt.Sprintf("%s", p.Context),
		}).Warn("Could not fetch policies from store.")
		pkg.WriteJSON(rw, struct {
			Allowed bool   `json:"allowed"`
			Error   string `json:"error"`
		}{Allowed: false, Error: err.Error()})
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	allowed, err := h.g.IsGranted(p.Resource, p.Permission, p.Subject, policies, p.Context)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"resource":   p.Resource,
			"permission": p.Permission,
			"subject":    p.Subject,
			"policies":   fmt.Sprintf("%s", policies),
			"context":    fmt.Sprintf("%s", p.Context),
		}).Warn("Granted check failed.")
		pkg.WriteJSON(rw, struct {
			Allowed bool   `json:"allowed"`
			Error   string `json:"error"`
		}{Allowed: false, Error: err.Error()})
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"resource":   p.Resource,
		"permission": p.Permission,
		"subject":    p.Subject,
		"allowed":    allowed,
		"policies":   fmt.Sprintf("%s", policies),
		"context":    fmt.Sprintf("%s", p.Context),
	}).Info("Valid granted request.")
	pkg.WriteJSON(rw, struct {
		Allowed bool `json:"allowed"`
	}{Allowed: allowed})
}

func (h *Handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var p DefaultPolicy
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = uuid.New()
	if err := h.s.Create(&p); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	pkg.WriteJSON(rw, p)
}

func (h *Handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {

	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "get", nil)(hctx.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			policy, err := h.s.Get(id)
			if err != nil {
				http.NotFound(rw, req)
				return
			}
			pkg.WriteJSON(rw, policy)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "delete", nil)(hctx.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.Delete(id); err != nil {
				http.Error(rw, fmt.Sprintf("Could not retrieve client: %s", id), http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusAccepted)
		},
	)).ServeHTTPContext(ctx, rw, req)
}
