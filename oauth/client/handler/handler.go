package handler

import (
	"encoding/json"
	"fmt"
	"github.com/arekkas/osin"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/ory-am/common/rand/sequence"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/middleware"
	. "github.com/ory-am/hydra/pkg"
	"github.com/ory-am/osin-storage/storage"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"net/http"
)

type handler struct {
	s storage.Storage
	m *middleware.Middleware
}

type payload struct {
	ID           string `json:"id,omitempty" `
	Secret       string `json:"secret,omitempty"`
	RedirectURIs string `valid:"required", json:"redirectURIs"`
}

func NewHandler(s storage.Storage, m *middleware.Middleware) *handler {
	return &handler{s, m}
}

func (h *handler) SetRoutes(r *mux.Router, extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) {
	r.Handle("/clients", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
		h.m.IsAuthorized("rn:hydra:clients", "create"),
	).ThenFunc(h.Create)).Methods("POST")

	r.Handle("/clients/{id}", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Get)).Methods("GET")

	r.Handle("/clients/{id}", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Delete)).Methods("DELETE")
}

func (h *handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {

	var p payload
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

	secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &osin.DefaultClient{
		Id:          uuid.New(),
		Secret:      string(secret),
		RedirectUri: p.RedirectURIs,
		UserData:    "{}",
	}

	if err := h.s.CreateClient(client); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(rw, client)
}

func (h *handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(fmt.Sprintf("rn:hydra:clients:%s", id), "get")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			client, err := h.s.GetClient(id)
			if err != nil {
				http.Error(rw, fmt.Sprintf("Could not retrieve client: %s", id), http.StatusNotFound)
				return
			}
			WriteJSON(rw, client)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(fmt.Sprintf("rn:hydra:clients:%s", id), "delete")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.RemoveClient(id); err != nil {
				http.Error(rw, fmt.Sprintf("Could not retrieve client: %s", id), http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusAccepted)
		},
	)).ServeHTTPContext(ctx, rw, req)
}
