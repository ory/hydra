package handler

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	. "github.com/ory-am/hydra/account"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/handler/middleware"
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
	r.Handle("/users", &hydcon.ContextAdapter{Ctx: context.Background(), Handler: h.Create(extractor)}).Methods("POST")
	r.Handle("/users/{id}", &hydcon.ContextAdapter{Ctx: context.Background(), Handler: h.Get(extractor)}).Methods("GET")
	r.Handle("/users/{id}", &hydcon.ContextAdapter{Ctx: context.Background(), Handler: h.Delete(extractor)}).Methods("DELETE")

	//r.Handle("/login", &hydcon.ContextAdapter{Ctx: context.Background(), Handler: h.Login(extractor)}).Methods("GET")
	//r.Handle("/logout", &hydcon.ContextAdapter{Ctx: context.Background(), Handler: h.Logout(extractor)}).Methods("GET")
}

func (h *Handler) Create(extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) hydcon.ContextHandler {
	return extractor(h.m.IsAuthenticated(h.m.IsAuthorized(hydcon.ContextHandlerFunc(
		func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
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

			js, err := json.Marshal(user)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.Write(js)
		}), "/users", "create"),
	))
}

func (h *Handler) Get(extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) hydcon.ContextHandler {
	return extractor(h.m.IsAuthenticated(hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		id, ok := mux.Vars(req)["id"]
		if !ok {
			http.Error(rw, "No id given.", http.StatusBadRequest)
			return
		}
		h.m.IsAuthorized(hydcon.ContextHandlerFunc(
			func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				user, err := h.s.Get(id)
				if err == ErrNotFound {
					http.Error(rw, err.Error(), http.StatusNotFound)
					return
				} else if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				js, err := json.Marshal(user)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Header().Set("Content-Type", "application/json")
				rw.Write(js)
			}), "/users/"+id, "get",
		).ServeHTTPContext(ctx, rw, req)
	})))
}

func (h *Handler) Delete(extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) hydcon.ContextHandler {
	return extractor(h.m.IsAuthenticated(hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		id, ok := mux.Vars(req)["id"]
		if !ok {
			http.Error(rw, "No id given.", http.StatusBadRequest)
			return
		}
		h.m.IsAuthorized(hydcon.ContextHandlerFunc(
			func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				if err := h.s.Delete(id); err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.WriteHeader(http.StatusAccepted)
			}), "/users/"+id, "delete",
		).ServeHTTPContext(ctx, rw, req)
	})))
}
