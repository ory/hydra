package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	. "github.com/ory-am/common/pkg"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/hydra/middleware"
	"github.com/ory-am/osin-storage/storage"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"

	"github.com/go-errors/errors"
)

type Handler struct {
	s storage.Storage
	m middleware.Middleware
}

func permission(id string) string {
	return fmt.Sprintf("rn:hydra:clients:%s", id)
}

func NewHandler(s storage.Storage, m middleware.Middleware) *Handler {
	return &Handler{s, m}
}

type payload struct {
	ID           string `json:"id,omitempty" `
	Secret       string `json:"secret,omitempty"`
	RedirectURIs string `valid:"required" json:"redirectURIs"`
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h chd.ContextHandler) chd.ContextHandler) {
	r.Handle("/clients", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
		h.m.IsAuthorized("rn:hydra:clients", "create", nil),
	).ThenFunc(h.Create)).Methods("POST")

	r.Handle("/clients/{id}", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Get)).Methods("GET")

	r.Handle("/clients/{id}", chd.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Delete)).Methods("DELETE")
}

func (h *Handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var p payload
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		HttpError(rw, err, http.StatusBadRequest)
		return
	}

	if v, err := govalidator.ValidateStruct(p); !v {
		if err != nil {
			HttpError(rw, err, http.StatusBadRequest)
			return
		}
		HttpError(rw, errors.New("Payload did not validate."), http.StatusBadRequest)
		return
	}

	secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))
	if err != nil {
		HttpError(rw, err, http.StatusInternalServerError)
		return
	}

	client := &osin.DefaultClient{
		Id:          uuid.New(),
		Secret:      string(secret),
		RedirectUri: p.RedirectURIs,
		UserData:    "{}",
	}

	if err := h.s.CreateClient(client); err != nil {
		WriteError(rw, err)
		return
	}

	WriteCreatedJSON(rw, "/clients/"+client.Id, client)
}

func (h *Handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.New("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "get", nil)(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			client, err := h.s.GetClient(id)
			if err != nil {
				WriteError(rw, err)
				return
			}
			WriteJSON(rw, client)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *Handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		HttpError(rw, errors.New("No id given."), http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(permission(id), "delete", nil)(chd.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.RemoveClient(id); err != nil {
				HttpError(rw, errors.Errorf("Could not retrieve client: %s", id), http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusAccepted)
		},
	)).ServeHTTPContext(ctx, rw, req)
}
