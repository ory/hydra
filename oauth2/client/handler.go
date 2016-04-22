package client

import (
	"encoding/json"
	"net/http"

	. "github.com/ory-am/common/pkg"
	"github.com/ory-am/common/rand/sequence"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"

	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/fosite/client"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/ladon"
	"fmt"
)

type Handler struct {
	m ClientManager
	h herodot.Herodot
	w warden.Warden
}

const (
	clientsPermission = "rn:hydra:clients"
	clientPermission = "rn:hydra:clients:%s"
)

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/clients", h.Create)
	r.GET("/clients/:id", h.Get)
	r.DELETE("/clients/:id", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var c client.SecureClient
	ctx := herodot.NewContext()

	if _, err := h.w.HTTPActionAllowed(r, &ladon.Request{
		Resource: clientsPermission,
		Action: "create",
	}, "hydra.clients"); err != nil {
		return h.h.WriteError(ctx, w, r, errors.New(err))
	}

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.h.WriteError(ctx, w, r, errors.New(err))
		return
	}

	secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))
	if err != nil {
		h.h.WriteError(ctx, w, r, errors.New(err))
		return
	}

	c.ID = uuid.New()
	c.Secret = secret
	if err := h.m.CreateClient(&c); err != nil {
		h.h.WriteError(ctx, w, r, errors.New(err))
		return
	}

	h.h.WriteCreated(ctx, w, r, "/clients/" + c.GetID(), c)
}

func (h *Handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, err := h.w.HTTPActionAllowed(r, &ladon.Request{
		Resource: fmt.Sprintf(clientPermission, id),
		Action: "get",
	}, "hydra.clients"); err != nil {
		return h.h.WriteError(ctx, w, r, errors.New(err))
	}

	c, err := h.m.GetClient(id)
	if err != nil {
		h.h.WriteError(ctx, w, r, errors.New(err))
		return
	}

	h.h.Write(ctx, w, r, c)
}

func (h *Handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, err := h.w.HTTPActionAllowed(r, &ladon.Request{
		Resource: fmt.Sprintf(clientPermission, id),
		Action: "delete",
	}, "hydra.clients"); err != nil {
		return h.h.WriteError(ctx, w, r, errors.New(err))
	}

	if err := h.m.RemoveClient(id); err != nil {
		h.h.WriteError(ctx, w, r, errors.New(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
