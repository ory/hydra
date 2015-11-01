package handler

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/middleware"
	. "github.com/ory-am/hydra/oauth/connection"
	. "github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"net/http"
)

const permission = "rn:hydra:oauth2:subjects:%s:connections"

type handler struct {
	s Storage
	m *middleware.Middleware
}

type payload struct {
	ID string `json:"id,omitempty" `
}

func NewHandler(s Storage, m *middleware.Middleware) *handler {
	return &handler{s, m}
}

func (h *handler) SetRoutes(r *mux.Router, extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) {
	r.Handle("/oauth2/subjects/{subject}/connections", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Create)).Methods("POST")

	r.Handle("/oauth2/subjects/{subject}/connections", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Find)).Methods("GET")

	r.Handle("/oauth2/connections/{id}", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Get)).Methods("GET")

	r.Handle("/oauth2/connections/{id}", hydcon.NewContextAdapter(
		context.Background(),
		extractor,
		h.m.IsAuthenticated,
	).ThenFunc(h.Delete)).Methods("DELETE")
}

func (h *handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	subject, ok := mux.Vars(req)["subject"]
	if !ok {
		http.Error(rw, "No subject given.", http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(fmt.Sprintf(permission, subject), "create")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			var conn DefaultConnection
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&conn); err != nil {
				http.Error(rw, "Could not decode request: "+err.Error(), http.StatusBadRequest)
				return
			}

			if v, err := govalidator.ValidateStruct(conn); !v {
				if err != nil {
					http.Error(rw, err.Error(), http.StatusBadRequest)
					return
				}
				http.Error(rw, "Payload did not validate.", http.StatusBadRequest)
				return
			}

			conn.ID = uuid.New()
			if err := h.s.Create(&conn); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			WriteJSON(rw, &conn)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *handler) Find(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	subject, ok := mux.Vars(req)["subject"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	h.m.IsAuthorized(fmt.Sprintf(permission, subject), "get")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			conns, err := h.s.FindAllByLocalSubject(subject)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			WriteJSON(rw, conns)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	conn, err := h.s.Get(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	h.m.IsAuthorized(fmt.Sprintf(permission, conn.GetLocalSubject()), "get")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			WriteJSON(rw, conn)
		},
	)).ServeHTTPContext(ctx, rw, req)
}

func (h *handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(rw, "No id given.", http.StatusBadRequest)
		return
	}

	conn, err := h.s.Get(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	h.m.IsAuthorized(fmt.Sprintf(permission, conn.GetLocalSubject()), "delete")(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if err := h.s.Delete(conn.GetID()); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusAccepted)
		},
	)).ServeHTTPContext(ctx, rw, req)
}
