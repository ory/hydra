// Package handler
//
// Defined permissions:
// * permission("rn:hydra:accounts") actions("create")
// * permission("rn:hydra:accounts:%s", id) actions("get", "delete")
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	. "github.com/ory-am/common/pkg"
	. "github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/middleware"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"net/http"
)

type Handler struct {
	s Storage
	m middleware.Middleware
}

func NewHandler(s Storage, m middleware.Middleware) *Handler {
	return &Handler{s, m}
}

func permission(id string) string {
	return fmt.Sprintf("rn:hydra:accounts:%s", id)
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h chd.ContextHandler) chd.ContextHandler) {
	r.Handle("/accounts", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
		h.m.IsAuthorized("rn:hydra:accounts", "create", nil),
	).ThenFunc(h.create)).Methods("POST")

	r.Handle("/accounts/{id}/password", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.updatePassword)).Methods("PUT")

	r.Handle("/accounts/{id}/data", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.updateData)).Methods("PUT")

	r.Handle("/accounts/{id}/username", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.updateUsername)).Methods("PUT")

	r.Handle("/accounts/{id}", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.get)).Methods("GET")

	r.Handle("/accounts/{id}", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.delete)).Methods("DELETE")
}

func (h *Handler) create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var p CreateAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		HttpError(rw, err, http.StatusBadRequest)
		return
	}

	// Force ID override
	p.ID = uuid.New()
	user, err := h.s.Create(p)
	if err != nil {
		WriteError(rw, err)
		return
	}

	WriteCreatedJSON(rw, fmt.Sprintf("/accounts/%s", user.GetID()), user)
}

func (h *Handler) updateUsername(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "put:data", middleware.NewEnv(req).Owner(id))(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			var p UpdateUsernameRequest
			if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
				HttpError(rw, err, http.StatusBadRequest)
				return
			}

			user, err := h.s.UpdateUsername(id, p)
			if err != nil {
				WriteError(rw, err)
				return
			}

			WriteJSON(rw, user)
		}),
	).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) updatePassword(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "put:password", middleware.NewEnv(req).Owner(id))(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			var p UpdatePasswordRequest
			if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
				HttpError(rw, err, http.StatusBadRequest)
				return
			}

			user, err := h.s.UpdatePassword(id, p)
			if err != nil {
				WriteError(rw, err)
				return
			}

			WriteJSON(rw, user)
		}),
	).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) updateData(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "put:data", middleware.NewEnv(req).Owner(id))(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			var p UpdateDataRequest
			if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
				HttpError(rw, err, http.StatusBadRequest)
				return
			}

			user, err := h.s.UpdateData(id, p)
			if err != nil {
				WriteError(rw, err)
				return
			}

			WriteJSON(rw, user)
		}),
	).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "get", middleware.NewEnv(req).Owner(id))(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			user, err := h.s.Get(id)
			if err != nil {
				WriteError(rw, err)
				return
			}

			WriteJSON(rw, user)
		}),
	).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "delete", middleware.NewEnv(req).Owner(id))(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.Delete(id); err != nil {
				HttpError(rw, err, http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusAccepted)
		}),
	).ServeHTTPContext(ctx, rw, req)
}
