package handler

import (
	"github.com/gorilla/mux"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/middleware"
	"github.com/ory-am/osin-storage/storage"
	"golang.org/x/net/context"
	"net/http"
)

type Handler struct {
	s storage.Storage
	m *middleware.Middleware
}

func NewHandler(s storage.Storage, m *middleware.Middleware) *Handler {
	return &Handler{s, m}
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h hydcon.ContextHandler) hydcon.ContextHandler) {
}

func (h *Handler) Create(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
}

func (h *Handler) Get(ctx context.Context, rw http.ResponseWriter, req *http.Request) {

}

func (h *Handler) Delete(ctx context.Context, rw http.ResponseWriter, req *http.Request) {

}
