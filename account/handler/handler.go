package handler

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	. "github.com/ory-am/hydra/account"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/middleware"
	. "github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"net/http"
)

type Handler struct {
	s Storage
	m *middleware.Middleware
}

func NewHandler(s Storage, m *middleware.Middleware) *Handler {
	return &Handler{s, m}
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) {
	r.Handle("/users", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
		h.m.IsAuthorized("rn:hydra:users", "create"),
	).ThenFunc(h.Create),
	).Methods("POST")

	r.Handle("/users/{id}", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Get)).Methods("GET")

	r.Handle("/users/{id}", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Delete)).Methods("DELETE")
}

func (h *Handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	type Payload struct {
		Email    string `valid:"email,required" json:"email" `
		Password string `valid:"length(6|254),required" json:"password"`
		Data     string `valid:"optional,json", json:"data"`
	}

	var p Payload
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if v, err := govalidator.ValidateStruct(p); !v {
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(rw, "Payload did not validate.", http.StatusBadRequest)
		return
	}

	if p.Data == "" {
		p.Data = "{}"
	}

	user, err := h.s.Create(uuid.New(), p.Email, p.Password, p.Data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(rw, user)
}

func (h *Handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}
	h.m.IsAuthorized(fmt.Sprintf("rn:hydra:users:%s", id), "get")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			user, err := h.s.Get(id)
			if err == ErrNotFound {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			WriteJSON(rw, user)
		}),
	).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}
	h.m.IsAuthorized(fmt.Sprintf("rn:hydra:users:%s", id), "delete")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.Delete(id); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusAccepted)
		}),
	).ServeHTTPContext(ctx, rw, req)
}
