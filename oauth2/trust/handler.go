package trust

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
	grantJWTBearerPath = "/trust/grants/jwt-bearer/issuers" // #nosec G101
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

// swagger:route POST /trust/grants/jwt-bearer/issuers admin trustJwtGrantIssuer
//
// Trust an OAuth2 JWT Bearer Grant Type Issuer
//
// Use this endpoint to establish a trust relationship for a JWT issuer
// to perform JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication
// and Authorization Grants [RFC7523](https://datatracker.ietf.org/doc/html/rfc7523).
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       201: trustedJwtGrantIssuer
//       400: genericError
//       409: genericError
//       500: genericError
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var grantRequest createGrantRequest

	if err := json.NewDecoder(r.Body).Decode(&grantRequest); err != nil {
		h.registry.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	if err := h.registry.GrantValidator().Validate(grantRequest); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	grant := Grant{
		ID:              uuid.New().String(),
		Issuer:          grantRequest.Issuer,
		Subject:         grantRequest.Subject,
		AllowAnySubject: grantRequest.AllowAnySubject,
		Scope:           grantRequest.Scope,
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

// swagger:route GET /trust/grants/jwt-bearer/issuers/{id} admin getTrustedJwtGrantIssuer
//
// Get a Trusted OAuth2 JWT Bearer Grant Type Issuer
//
// Use this endpoint to get a trusted JWT Bearer Grant Type Issuer. The ID is the one returned when you
// created the trust relationship.
///
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: trustedJwtGrantIssuer
//       404: genericError
//       500: genericError
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")

	grant, err := h.registry.GrantManager().GetConcreteGrant(r.Context(), id)
	if err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	h.registry.Writer().Write(w, r, grant)
}

// swagger:route DELETE /trust/grants/jwt-bearer/issuers/{id} admin deleteTrustedJwtGrantIssuer
//
// Delete a Trusted OAuth2 JWT Bearer Grant Type Issuer
//
// Use this endpoint to delete trusted JWT Bearer Grant Type Issuer. The ID is the one returned when you
// created the trust relationship.
//
// Once deleted, the associated issuer will no longer be able to perform the JSON Web Token (JWT) Profile
// for OAuth 2.0 Client Authentication and Authorization Grant.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       404: genericError
//       500: genericError
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")

	if err := h.registry.GrantManager().DeleteGrant(r.Context(), id); err != nil {
		h.registry.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /trust/grants/jwt-bearer/issuers admin listTrustedJwtGrantIssuers
//
// List Trusted OAuth2 JWT Bearer Grant Type Issuers
//
// Use this endpoint to list all trusted JWT Bearer Grant Type Issuers.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: trustedJwtGrantIssuers
//       500: genericError
func (h *Handler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, offset := pagination.Parse(r, 100, 0, 500)
	optionalIssuer := r.URL.Query().Get("issuer")

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
