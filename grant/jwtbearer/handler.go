package jwtbearer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/pagination"

	"github.com/ory/hydra/x"

	"github.com/julienschmidt/httprouter"
)

const (
	grantJWTBearerPath = "/grants/jwt-bearer"
)

type Handler struct {
	registry InternalRegistry
}

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{registry: r}
}

func (h *Handler) SetRoutes(admin *x.RouterAdmin) {
	admin.GET(grantJWTBearerPath+"/:id", h.Get)
	admin.GET(grantJWTBearerPath, h.List)

	admin.POST(grantJWTBearerPath, h.Create)

	admin.DELETE(grantJWTBearerPath+"/:id", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var grantRequest grantRequest

	if err := json.NewDecoder(r.Body).Decode(&grantRequest); err != nil {
		h.registry.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	if err := h.registry.GrantValidator().Validate(grantRequest); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	grant := Grant{
		ID:      uuid.New().String(),
		Issuer:  grantRequest.Issuer,
		Subject: grantRequest.Subject,
		Scope:   grantRequest.Scope,
		PublicKey: PublicKey{
			Set:   grantRequest.Issuer, // group all keys by issuer, so set=issuer
			KeyID: grantRequest.PublicKeyJWK.KeyID,
		},
		CreatedAt: time.Now().UTC().Round(time.Second),
		ExpiresAt: grantRequest.ExpiresAt.UTC().Round(time.Second),
	}

	if err := h.registry.GrantManager().CreateGrant(r.Context(), grant, grantRequest.PublicKeyJWK); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	h.registry.Writer().WriteCreated(w, r, grantJWTBearerPath+"/"+grant.ID, &grant)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")

	grant, err := h.registry.GrantManager().GetConcreteGrant(r.Context(), id)
	if err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	h.registry.Writer().Write(w, r, grant)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")

	if err := h.registry.GrantManager().DeleteGrant(r.Context(), id); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, offset := pagination.Parse(r, 100, 0, 500)
	var optionalIssuer = ps.ByName("issuer")

	grants, err := h.registry.GrantManager().GetGrants(r.Context(), limit, offset, optionalIssuer)
	if err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	n, err := h.registry.GrantManager().CountGrants(r.Context())
	if err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	pagination.Header(w, r.URL, n, limit, offset)

	if grants == nil {
		grants = []Grant{}
	}

	h.registry.Writer().Write(w, r, grants)
}
