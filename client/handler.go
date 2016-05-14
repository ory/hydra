package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/ladon"
)

type Handler struct {
	Manager Manager
	H       herodot.Herodot
	W       firewall.Firewall
}

const (
	ClientsHandlerPath = "/clients"
)

const (
	ClientsResource = "rn:hydra:clients"
	ClientResource  = "rn:hydra:clients:%s"
	Scope           = "hydra.clients"
)

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET(ClientsHandlerPath, h.GetAll)
	r.POST(ClientsHandlerPath, h.Create)
	r.GET(ClientsHandlerPath+"/:id", h.Get)
	r.DELETE(ClientsHandlerPath+"/:id", h.Delete)
}

func NewManager(c *config.Config) Manager {
	ctx := c.Context()

	switch ctx.Connection.(type) {
	case *config.MemoryConnection:
		return &MemoryManager{
			Clients: map[string]*fosite.DefaultClient{},
			Hasher:  ctx.Hasher,
		}
		break
	default:
		panic("Unknown connection type.")
	}
	return nil
}

func NewHandler(c *config.Config, router *httprouter.Router, manager Manager) *Handler {
	ctx := c.Context()
	h := &Handler{
		H:&herodot.JSON{},
		W: ctx.Warden,Manager:manager,
	}

	h.SetRoutes(router)
	return h
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var c fosite.DefaultClient
	var ctx = herodot.NewContext()

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: ClientsResource,
		Action:   "create",
		Context: ladon.Context{
			"owner": c.Owner,
		},
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-.,:;$%!&/()=?+*#<>"))
	if err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}
	c.Secret = []byte(string(secret))

	if err := h.Manager.CreateClient(&c); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.WriteCreated(ctx, w, r, ClientsHandlerPath+"/"+c.GetID(), &c)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: ClientsResource,
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	c, err := h.Manager.GetClients()
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, c)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	c, err := h.Manager.GetClient(id)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(ClientResource, id),
		Action:   "get",
		Context: ladon.Context{
			"owner": c.GetOwner(),
		},
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, c)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(ClientResource, id),
		Action:   "delete",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.DeleteClient(id); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
