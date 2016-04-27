package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
)

type ClientHandler struct {
	Manager ClientManager
	H       herodot.Herodot
	W       warden.Warden
}

const (
	ClientsHandlerPath = "/clients"
)

const (
	ClientsResource = "rn:hydra:clients"
	ClientResource  = "rn:hydra:clients:%s"
	Scope = "hydra.clients"
)

func (h *ClientHandler) SetRoutes(r *httprouter.Router) {
	r.POST(ClientsHandlerPath, h.Create)
	r.GET(ClientsHandlerPath+"/:id", h.Get)
	r.DELETE(ClientsHandlerPath+"/:id", h.Delete)
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var c Client
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

	secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))
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

func (h *ClientHandler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(ClientResource, id),
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	c, err := h.Manager.GetClient(id)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, c)
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	w.WriteHeader(http.StatusAccepted)
}
