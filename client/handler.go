package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
)

type Handler struct {
	Manager Manager
	H       herodot.Writer
	W       firewall.Firewall
}

const (
	ClientsHandlerPath = "/clients"
)

const (
	ClientsResource = "rn:hydra:clients"
	ClientResource = "rn:hydra:clients:%s"
	Scope = "hydra.clients"
)

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET(ClientsHandlerPath, h.List)
	r.POST(ClientsHandlerPath, h.Create)
	r.GET(ClientsHandlerPath + "/:id", h.Get)
	r.PUT(ClientsHandlerPath + "/:id", h.Update)
	r.DELETE(ClientsHandlerPath + "/:id", h.Delete)
}

// swagger:parameters createOAuthClient updateOAuthClient
type createClientPayload struct {
	// in: body
	// required: true
	Client
}

// swagger:route POST /clients oauth2 clients createOAuthClient
//
// Updates an OAuth 2.0 Client. Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This
// endpoint should be well protected and only called by code you trust.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.clients
//
//     Responses:
//       200: oauthClient
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var c Client
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: ClientsResource,
		Action:   "create",
		Context: map[string]interface{}{
			"owner": c.Owner,
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if len(c.Secret) == 0 {
		secret, err := sequence.RuneSequence(12, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-.,:;$%!&/()=?+*#<>"))
		if err != nil {
			h.H.WriteError(w, r, errors.WithStack(err))
			return
		}
		c.Secret = string(secret)
	} else if len(c.Secret) < 6 {
		h.H.WriteError(w, r, errors.New("The client secret must be at least 6 characters long"))
	}

	secret := c.Secret
	if err := h.Manager.CreateClient(&c); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	c.Secret = secret
	h.H.WriteCreated(w, r, ClientsHandlerPath + "/" + c.GetID(), &c)
}

// swagger:route PUT /clients oauth2 clients updateOAuthClient
//
// Updates an OAuth 2.0 Client. Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This
// endpoint should be well protected and only called by code you trust.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.clients
//
//     Responses:
//       200: oauthClient
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c Client
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	o, err := h.Manager.GetConcreteClient(ps.ByName("id"))
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: ClientsResource,
		Action:   "update",
		Context: ladon.Context{
			"owner": o.Owner,
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if len(c.Secret) > 0 && len(c.Secret) < 6 {
		h.H.WriteError(w, r, errors.New("The client secret must be at least 6 characters long"))
	}

	c.ID = ps.ByName("id")
	if err := h.Manager.UpdateClient(&c); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	h.H.WriteCreated(w, r, ClientsHandlerPath + "/" + c.GetID(), &c)
}

// A list of clients.
// swagger:response clientsList
type listClientsResult struct {
	// in: body
	Clients []Client
}

// swagger:route GET /clients oauth2 clients listOAuthClients
//
// Fetches OAuth 2.0 Clients, never returns a client's secret.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.clients
//
//     Responses:
//       200: clientsList
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: ClientsResource,
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	c, err := h.Manager.GetClients()
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	for k, cc := range c {
		cc.Secret = ""
		c[k] = cc
	}

	h.H.Write(w, r, c)
}

// swagger:parameters getOAuthClient deleteOAuthClient
type queryClientPayload struct {
	// The id of the OAuth 2.0 Client.
	//
	// unique: true
	// in: path
	ID string `json:"id"`
}

// swagger:route GET /clients/{id} oauth2 clients getOAuthClient
//
// Fetches an OAuth 2.0 Client. Never returns the client's secret.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.clients
//
//     Responses:
//       200: oauthClient
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	c, err := h.Manager.GetConcreteClient(id)
	if err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ClientResource, id),
		Action:   "get",
		Context: ladon.Context{
			"owner": c.GetOwner(),
		},
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	c.Secret = ""
	h.H.Write(w, r, c)
}

// swagger:route DELETE /clients/{id} oauth2 clients deleteOAuthClient
//
// Deletes an OAuth 2.0 Client.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.clients
//
//     Responses:
//       204
//       401: genericError
//       403: genericError
//       500: genericError
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	var id = ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ClientResource, id),
		Action:   "delete",
	}, Scope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if err := h.Manager.DeleteClient(id); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
