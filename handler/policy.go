// Package handler
//
// Defined permissions:
// * permission("rn:hydra:policies") actions("create")
// * permission("rn:hydra:policies:%s", id) actions("get", "delete")
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	log "github.com/Sirupsen/logrus"
	"github.com/arekkas/flitt-api/vendor/github.com/ory-am/hydra/middleware"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	hctx "github.com/ory-am/common/handler"
	"github.com/ory-am/common/pkg"
	. "github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/guard/operator"
	. "github.com/ory-am/ladon/policy"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

type Handler struct {
	s Storage
	m middleware.Middleware
	o osin.Storage
	g Guarder
	j *jwt.JWT
}

type GrantedPayload struct {
	Resource   string            `json:"resource"`
	Token      string            `json:"token"`
	Permission string            `json:"permission"`
	Context    *operator.Context `json:"context"`
}

type CreatePolicyPayload struct {
	Description string              `json:"description"`
	Subjects    []string            `json:"subjects"`
	Effect      string              `json:"effect"`
	Resources   []string            `json:"resources"`
	Permissions []string            `json:"permissions"`
	Conditions  []*DefaultCondition `json:"conditions"`
}

func permission(id string) string {
	return fmt.Sprintf("rn:hydra:policies:%s", id)
}

func NewHandler(s Storage, m middleware.Middleware, g Guarder, j *jwt.JWT, o osin.Storage) *Handler {
	return &Handler{s: s, m: m, g: g, j: j, o: o}
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h hctx.ContextHandler) hctx.ContextHandler) {
	r.Handle("/guard/allowed", hctx.NewContextAdapter(
		context.Background(),
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
	auth, err := osin.CheckBasicAuth(req)
	if err != nil {
		pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	} else if auth == nil {
		pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	client, err := h.o.GetClient(auth.Username)
	if err != nil {
		pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	} else if client.GetSecret() != auth.Password {
		pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	var p GrantedPayload
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		pkg.HttpError(rw, errors.New(err), http.StatusBadRequest)
		return
	}

	token, err := h.j.VerifyToken([]byte(p.Token))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"resource":   p.Resource,
			"permission": p.Permission,
			"subject":    "",
			"context":    fmt.Sprintf("%s", p.Context),
		}).Warn("Token not valid.")
		pkg.WriteJSON(rw, struct {
			Allowed bool   `json:"allowed"`
			Error   string `json:"error"`
		}{Allowed: false, Error: err.Error()})
		return
	}

	subject, ok := token.Claims["sub"].(string)
	if !ok {
		err := errors.New("Bearer token is not valid.")
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"resource":   p.Resource,
			"permission": p.Permission,
			"subject":    "",
			"context":    fmt.Sprintf("%s", p.Context),
		}).Warn("Token does not claim a subject.")
		pkg.WriteJSON(rw, struct {
			Allowed bool   `json:"allowed"`
			Error   string `json:"error"`
		}{Allowed: false, Error: err.Error()})
		return
	}

	policies, err := h.s.FindPoliciesForSubject(subject)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"resource":   p.Resource,
			"permission": p.Permission,
			"subject":    subject,
			"context":    fmt.Sprintf("%s", p.Context),
		}).Warn("Could not fetch policies from store.")
		pkg.WriteJSON(rw, struct {
			Allowed bool   `json:"allowed"`
			Error   string `json:"error"`
		}{Allowed: false, Error: err.Error()})
		return
	}

	allowed, err := h.g.IsGranted(p.Resource, p.Permission, subject, policies, p.Context)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"resource":   p.Resource,
			"permission": p.Permission,
			"subject":    subject,
			"policies":   fmt.Sprintf("%s", policies),
			"context":    fmt.Sprintf("%s", p.Context),
		}).Warn("Granted check failed.")
		pkg.WriteJSON(rw, struct {
			Allowed bool   `json:"allowed"`
			Error   string `json:"error"`
		}{Allowed: false, Error: err.Error()})
		return
	}

	log.WithFields(log.Fields{
		"resource":   p.Resource,
		"permission": p.Permission,
		"subject":    subject,
		"allowed":    allowed,
		"policies":   fmt.Sprintf("%s", policies),
		"context":    fmt.Sprintf("%s", p.Context),
	}).Info("Got guard decision.")
	pkg.WriteJSON(rw, struct {
		Allowed bool `json:"allowed"`
	}{Allowed: allowed})
}

func (h *Handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var p DefaultPolicy
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		pkg.HttpError(rw, err, http.StatusBadRequest)
		return
	}

	p.ID = uuid.New()
	if err := h.s.Create(&p); err != nil {
		pkg.HttpError(rw, err, http.StatusInternalServerError)
		return
	}

	pkg.WriteCreatedJSON(rw, "/policies/"+p.ID, p)
}

func (h *Handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		pkg.HttpError(rw, errors.New("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "get", nil)(hctx.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			policy, err := h.s.Get(id)
			if err != nil {
				pkg.WriteError(rw, err)
				return
			}
			pkg.WriteJSON(rw, policy)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		pkg.HttpError(rw, errors.New("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "delete", nil)(hctx.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.Delete(id); err != nil {
				pkg.HttpError(rw, errors.Errorf("Could not retrieve client: %s", id), http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusAccepted)
		},
	)).ServeHTTPContext(ctx, rw, req)
}
